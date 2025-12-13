package sdk

import (
	"encoding/json"
	"fmt"
	"iter"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ddkwork/golibrary/std/stream"
	"github.com/google/uuid"
)

// FieldType represents the type of a field in the table.
type FieldType string

// Constants for FieldType.
const (
	FieldTypeSingleLineText FieldType = "text"
	FieldTypeNumber         FieldType = "number"
	FieldTypeSingleSelect   FieldType = "singleSelect"
	FieldTypeMultipleSelect FieldType = "multipleSelect"
	FieldTypeDateTime       FieldType = "dateTime"
	FieldTypeFormula        FieldType = "formula"
	FieldTypeAttachment     FieldType = "attachment"
	FieldTypeLink           FieldType = "link"
	FieldTypeUser           FieldType = "user"
	FieldTypePhone          FieldType = "phone"
	FieldTypeEmail          FieldType = "email"
	FieldTypeCheckbox       FieldType = "checkbox"
	FieldTypeURL            FieldType = "url"
	FieldTypeMultiLineText  FieldType = "multiLineText"
	FieldTypeCurrency       FieldType = "currency"
	FieldTypePercent        FieldType = "percent"
)

// Node represents a node in the tree table.
type Node struct {
	ID        string     // Unique identifier for the node (UUID)
	Type      string     // Node type (container nodes end with "_container")
	RowCells  []CellData // Row data (including formula columns)
	Children  []*Node    // Child nodes
	parent    *Node      // Parent node
	isOpen    bool       // Whether expanded (only for container nodes)
	GroupKey  string     // Grouping key
	RowNumber int        // Row number (for sorting)
}

// TreeTable represents the tree table structure.
type TreeTable struct {
	Root         *Node                        // Root node (virtual container)
	OriginalRoot *Node                        // Backup of the original root node
	Columns      []ColumnDefinition           // Header definitions (using ColumnDefinition)
	columnMap    map[string]*ColumnDefinition // Mapping from column name to definition
	SelectedNode *Node                        // Currently selected node
	once         sync.Once                    // One-time initialization marker

	// Callback functions
	OnRowSelected    func(n *Node)
	OnRowDoubleClick func(n *Node)
}

// ColumnDefinition defines a column in the table.
type ColumnDefinition struct {
	Name         string    // Column name (unique identifier)
	Type         FieldType // Data type
	Formula      string    // Column formula text (stores Go code!)
	Options      []string  // Options (for single/multiple select)
	IsDisabled   bool      // Whether editing is disabled
	Width        int       // Column width in pixels
	DefaultValue any       // Default value
	Values       []any     // Initial values for all cells in the column (for batch initialization)
}

// CellData represents a cell in the table.
type CellData struct {
	ColumnName string    // Associated column name (key! find column definition via this field)
	Value      any       // Cell value
	Type       FieldType // Data type (can inherit from column definition)
}

// ColumnConfig 列配置（简化版）
type ColumnConfig struct {
	Name     string    // 列名
	Type     FieldType // 数据类型
	Formula  string    // 公式（可选）
	Options  []string  // 选项（可选）
	Width    int       // 宽度（可选）
	Disabled bool      // 是否禁用（可选）
}

// TableData 表格数据（核心！）
type TableData struct {
	Columns []ColumnConfig // 列定义
	Rows    [][]any        // 行数据（每行对应一个数据记录）
}

// AsString returns the string representation of the cell value.
func (c *CellData) AsString() (string, bool) {
	v, ok := c.Value.(string)
	return v, ok
}

// AsInt returns the integer representation of the cell value.
func (c *CellData) AsInt() (int, bool) {
	switch v := c.Value.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	case int64:
		return int(v), true
	default:
		if i, err := strconv.Atoi(fmt.Sprintf("%v", c.Value)); err == nil {
			return i, true
		}
		return 0, false
	}
}

// AsFloat returns the float64 representation of the cell value.
func (c *CellData) AsFloat() (float64, bool) {
	switch v := c.Value.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	default:
		if f, err := strconv.ParseFloat(fmt.Sprintf("%v", c.Value), 64); err == nil {
			return f, true
		}
		return 0, false
	}
}

// AsBool returns the boolean representation of the cell value.
func (c *CellData) AsBool() (bool, bool) {
	v, ok := c.Value.(bool)
	if ok {
		return v, true
	}

	// String representations of boolean
	if s, ok := c.Value.(string); ok {
		switch strings.ToLower(s) {
		case "true", "1", "yes", "是":
			return true, true
		case "false", "0", "no", "否":
			return false, true
		}
	}

	return false, false
}

// AsTime returns the time.Time representation of the cell value.
func (c *CellData) AsTime() (time.Time, bool) {
	switch v := c.Value.(type) {
	case time.Time:
		return v, true
	case string:
		// Try multiple time formats
		formats := []string{
			time.RFC3339,
			"2006-01-02",
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05Z",
			"01/02/2006",
			"2006/01/02",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, v); err == nil {
				return t, true
			}
		}
	}
	return time.Time{}, false
}

// IsFormula checks if the cell is a formula column.
func (c *CellData) IsFormula() bool {
	return c.Type == FieldTypeFormula
}

// IsSelect checks if the cell is a select type (single or multiple).
func (c *CellData) IsSelect() bool {
	return c.Type == FieldTypeSingleSelect || c.Type == FieldTypeMultipleSelect
}

// IsAttachment checks if the cell is an attachment type.
func (c *CellData) IsAttachment() bool {
	return c.Type == FieldTypeAttachment
}

// IsLink checks if the cell is a link type.
func (c *CellData) IsLink() bool {
	return c.Type == FieldTypeLink
}

// DetectValueType detects the field type of the cell value.
func (c *CellData) DetectValueType() FieldType {
	return inferType(c.Value)
}

// NewNode creates a new node with the given row cells.
func NewNode(rowCells []CellData) *Node {
	return &Node{
		ID:        uuid.New().String(),
		Type:      "node",
		RowCells:  rowCells,
		Children:  nil,
		parent:    nil,
		isOpen:    false,
		GroupKey:  "",
		RowNumber: 0,
	}
}

// NewContainerNode creates a new container node.
func NewContainerNode(typeKey string, rowCells []CellData) *Node {
	n := NewNode(rowCells)
	n.Type = typeKey + "_container"
	n.isOpen = true
	return n
}

// Clone creates a deep copy of the node.
func (n *Node) Clone() *Node {
	clone := &Node{
		ID:        uuid.New().String(),
		Type:      n.Type,
		RowCells:  make([]CellData, len(n.RowCells)),
		Children:  make([]*Node, len(n.Children)),
		isOpen:    n.isOpen,
		GroupKey:  n.GroupKey,
		RowNumber: n.RowNumber,
	}

	// Copy row data
	for i, cell := range n.RowCells {
		clone.RowCells[i] = cell
	}

	// Copy child nodes
	for i, child := range n.Children {
		cloneChild := child.Clone()
		cloneChild.parent = clone
		clone.Children[i] = cloneChild
	}

	return clone
}

// AddChild adds a child node.
func (n *Node) AddChild(child *Node) {
	child.parent = n
	n.Children = append(n.Children, child)
}

// AddChildren adds multiple child nodes.
func (n *Node) AddChildren(children []*Node) {
	for _, child := range children {
		n.AddChild(child)
	}
}

// InsertChild inserts a child node at the specified position.
func (n *Node) InsertChild(index int, child *Node) {
	if index < 0 || index > len(n.Children) {
		index = len(n.Children)
	}
	child.parent = n
	n.Children = append(n.Children[:index], append([]*Node{child}, n.Children[index:]...)...)
}

// RemoveChild removes a child node.
func (n *Node) RemoveChild(child *Node) {
	for i, c := range n.Children {
		if c.ID == child.ID {
			n.Children = append(n.Children[:i], n.Children[i+1:]...)
			return
		}
	}
}

// IsContainer checks if the node is a container.
func (n *Node) IsContainer() bool {
	return strings.HasSuffix(n.Type, "_container")
}

// Depth returns the depth of the node in the tree.
func (n *Node) Depth() int {
	depth := 0
	for p := n.parent; p != nil; p = p.parent {
		depth++
	}
	return depth
}

// Walk iterates over the node and its descendants.
func (n *Node) Walk() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		if !yield(n) {
			return
		}
		for _, child := range n.Children {
			for node := range child.Walk() {
				if !yield(node) {
					return
				}
			}
		}
	}
}

// GetCell retrieves a cell by column name.
func (n *Node) GetCell(colName string, _ *TreeTable) *CellData {
	for i := range n.RowCells {
		if n.RowCells[i].ColumnName == colName {
			return &n.RowCells[i]
		}
	}
	return nil
}

// SetCellValue sets the value of a cell by column name.
func (n *Node) SetCellValue(colName string, value any, table *TreeTable) {
	for i := range n.RowCells {
		if n.RowCells[i].ColumnName == colName {
			cell := &n.RowCells[i]
			cell.Value = value
			if cell.Type == "" {
				cell.Type = inferType(value)
			}
			return
		}
	}
	// Column does not exist, add it
	colDef := table.GetColumnDefinition(colName)
	if colDef != nil {
		n.RowCells = append(n.RowCells, CellData{
			ColumnName: colName,
			Value:      value,
			Type:       colDef.Type,
		})
	} else {
		n.RowCells = append(n.RowCells, CellData{
			ColumnName: colName,
			Value:      value,
			Type:       inferType(value),
		})
	}
}

// Walk2 iterates over the node and its descendants with indices.
func (n *Node) Walk2() iter.Seq2[int, *Node] {
	return func(yield func(int, *Node) bool) {
		if !yield(0, n) {
			return
		}
		for i, child := range n.Children {
			for node := range child.Walk() {
				if !yield(i, node) {
					return
				}
			}
		}
	}
}

// NewTreeTable creates a new TreeTable instance.
func NewTreeTable() *TreeTable {
	return &TreeTable{
		Root:      NewContainerNode("root", nil),
		columnMap: make(map[string]*ColumnDefinition),
	}
}

// initColumnMap initializes the column name to definition mapping.
func (t *TreeTable) initColumnMap() {
	t.columnMap = make(map[string]*ColumnDefinition)
	for i := range t.Columns {
		col := &t.Columns[i]
		t.columnMap[col.Name] = col
	}
}

// GetColumnDefinition returns the column definition by name.
func (t *TreeTable) GetColumnDefinition(colName string) *ColumnDefinition {
	if col, ok := t.columnMap[colName]; ok {
		return col
	}
	return nil
}

// GetRow 通过行索引获取行节点
func (t *TreeTable) GetRow(rowIndex int) *Node {
	rows := t.AllRows()
	if rowIndex < 0 || rowIndex >= len(rows) {
		return nil
	}
	return rows[rowIndex]
}

// GetCellValue 通过行索引和列名获取单元格值
func (t *TreeTable) GetCellValue(rowIndex int, colName string) (any, bool) {
	row := t.GetRow(rowIndex)
	if row == nil {
		return nil, false
	}
	cell := row.GetCell(colName, t)
	if cell == nil {
		return nil, false
	}
	return cell.Value, true
}

// SetCellValue 通过行索引和列名设置单元格值
func (t *TreeTable) SetCellValue(rowIndex int, colName string, value any) bool {
	row := t.GetRow(rowIndex)
	if row == nil {
		return false
	}
	row.SetCellValue(colName, value, t)
	return true
}

// dataNodes returns a sequence of all data nodes (children of the root).
func (t *TreeTable) dataNodes() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		for _, child := range t.Root.Children {
			stack := []*Node{child}
			for len(stack) > 0 {
				n := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				if !yield(n) {
					return
				}

				for i := len(n.Children) - 1; i >= 0; i-- {
					stack = append(stack, n.Children[i])
				}
			}
		}
	}
}

// dataNodesSlice returns a slice of all data nodes.
func (t *TreeTable) dataNodesSlice() []*Node {
	var nodes []*Node
	for node := range t.dataNodes() {
		nodes = append(nodes, node)
	}
	return nodes
}

// AllRows returns all row nodes (depth-first traversal, skipping root).
func (t *TreeTable) AllRows() []*Node {
	return t.dataNodesSlice()
}

// RowCount returns the number of rows.
func (t *TreeTable) RowCount() int {
	count := 0
	for range t.dataNodes() {
		count++
	}
	return count
}

// ColCount returns the number of columns.
func (t *TreeTable) ColCount() int {
	return len(t.Columns)
}

// ColIndex returns the index of a column by name.
func (t *TreeTable) ColIndex(colName string) int {
	for i, col := range t.Columns {
		if col.Name == colName {
			return i
		}
	}
	return -1
}

// ColName returns the name of a column by index.
func (t *TreeTable) ColName(colIndex int) string {
	if colIndex < 0 || colIndex >= len(t.Columns) {
		return ""
	}
	return t.Columns[colIndex].Name
}

// getDefaultValue returns the default value for a field type.
func getDefaultValue(ft FieldType) any {
	switch ft {
	case FieldTypeNumber, FieldTypeCurrency, FieldTypePercent:
		return 0.0
	case FieldTypeCheckbox:
		return false
	case FieldTypeDateTime:
		return time.Now()
	case FieldTypeSingleSelect, FieldTypeMultipleSelect:
		return ""
	case FieldTypeURL:
		return "https://"
	case FieldTypeEmail:
		return "example@domain.com"
	case FieldTypePhone:
		return "+1234567890"
	default:
		return ""
	}
}

// AddColumn adds a new column.
func (t *TreeTable) AddColumn(col ColumnDefinition, index int) {
	if index < 0 || index > len(t.Columns) {
		index = len(t.Columns)
	}

	// Check if column with same name already exists
	if _, exists := t.columnMap[col.Name]; exists {
		return
	}

	// Insert new column
	t.Columns = append(t.Columns[:index], append([]ColumnDefinition{col}, t.Columns[index:]...)...)
	t.initColumnMap()

	// Add new cell to all rows
	for node := range t.dataNodes() {
		node.SetCellValue(col.Name, getDefaultValue(col.Type), t)
	}
}

// DeleteColumn deletes a column.
func (t *TreeTable) DeleteColumn(colName string) bool {
	idx := -1
	for i, col := range t.Columns {
		if col.Name == colName {
			idx = i
			break
		}
	}
	if idx == -1 {
		return false
	}

	// Remove from column definitions
	t.Columns = append(t.Columns[:idx], t.Columns[idx+1:]...)
	t.initColumnMap()

	// Remove cell from all rows
	for node := range t.dataNodes() {
		for i := len(node.RowCells) - 1; i >= 0; i-- {
			if node.RowCells[i].ColumnName == colName {
				node.RowCells = append(node.RowCells[:i], node.RowCells[i+1:]...)
			}
		}
	}
	return true
}

// RenameColumn renames a column.
func (t *TreeTable) RenameColumn(oldName, newName string) bool {
	idx := -1
	for i, col := range t.Columns {
		if col.Name == oldName {
			idx = i
			break
		}
	}
	if idx == -1 {
		return false
	}

	// Update column definition
	t.Columns[idx].Name = newName

	// Update column mapping
	delete(t.columnMap, oldName)
	t.columnMap[newName] = &t.Columns[idx]

	// Update cell names in all rows
	for node := range t.dataNodes() {
		for i, cell := range node.RowCells {
			if cell.ColumnName == oldName {
				node.RowCells[i].ColumnName = newName
				break
			}
		}
	}
	return true
}

// UpdateColumn updates a column's attributes.
func (t *TreeTable) UpdateColumn(colName string, updateFunc func(*ColumnDefinition)) bool {
	colDef := t.GetColumnDefinition(colName)
	if colDef == nil {
		return false
	}

	// Apply update function
	updateFunc(colDef)

	// Update cells in all rows
	for node := range t.dataNodes() {
		for i := range node.RowCells {
			if node.RowCells[i].ColumnName == colName {
				// Update cell type if needed
				if node.RowCells[i].Type == "" {
					node.RowCells[i].Type = colDef.Type
				}
				break
			}
		}
	}
	return true
}

// SetRootRows sets the root rows using column definitions.
func (t *TreeTable) SetRootRows(columns []ColumnDefinition) {
	// Create new root container node
	t.Root = NewContainerNode("root", nil)
	t.OriginalRoot = t.Root.Clone()

	// Set column definitions
	t.Columns = make([]ColumnDefinition, len(columns))
	copy(t.Columns, columns)
	t.initColumnMap()

	// Determine row count (longest column)
	rowCount := 0
	for _, col := range columns {
		if len(col.Values) > rowCount {
			rowCount = len(col.Values)
		}
	}

	// Add rows
	for rowIdx := 0; rowIdx < rowCount; rowIdx++ {
		var cells []CellData
		for _, col := range columns {
			value := nilAny
			if rowIdx < len(col.Values) {
				value = col.Values[rowIdx]
			} else {
				value = col.DefaultValue
			}

			cells = append(cells, CellData{
				ColumnName: col.Name,
				Value:      value,
				Type:       col.Type,
			})
		}
		t.Root.AddChild(NewNode(cells))
	}
}

// LoadTableData 一键加载表格数据（超直观！）
func (t *TreeTable) LoadTableData(data TableData) {
	// 清空现有数据
	t.Root = NewContainerNode("root", nil)
	t.Columns = make([]ColumnDefinition, len(data.Columns))
	t.columnMap = make(map[string]*ColumnDefinition)

	// 设置列定义
	for i, colCfg := range data.Columns {
		t.Columns[i] = ColumnDefinition{
			Name:         colCfg.Name,
			Type:         colCfg.Type,
			Formula:      colCfg.Formula,
			Options:      colCfg.Options,
			Width:        colCfg.Width,
			IsDisabled:   colCfg.Disabled,
			DefaultValue: getDefaultValue(colCfg.Type),
			Values:       nil,
		}
		t.columnMap[colCfg.Name] = &t.Columns[i]
	}

	// 添加行数据
	for _, rowData := range data.Rows {
		var cells []CellData
		for colIdx, cellValue := range rowData {
			if colIdx < len(data.Columns) {
				colName := data.Columns[colIdx].Name
				cells = append(cells, CellData{
					ColumnName: colName,
					Value:      cellValue,
					Type:       data.Columns[colIdx].Type,
				})
			}
		}
		t.Root.AddChild(NewNode(cells))
	}
	stream.WriteTruncate("tmp/1.md", t.ToMarkdown())

}

// SortByColumn sorts all rows by the specified column.
func (t *TreeTable) SortByColumn(colName string, ascending bool) {
	rows := t.AllRows()

	// Sort rows based on column values
	sort.Slice(rows, func(i, j int) bool {
		cellI := rows[i].GetCell(colName, t)
		cellJ := rows[j].GetCell(colName, t)

		if cellI == nil && cellJ == nil {
			return false
		}
		if cellI == nil {
			return !ascending
		}
		if cellJ == nil {
			return ascending
		}

		// Compare values based on type
		switch cellI.Type {
		case FieldTypeNumber, FieldTypeCurrency, FieldTypePercent:
			valI, _ := cellI.AsFloat()
			valJ, _ := cellJ.AsFloat()
			if ascending {
				return valI < valJ
			}
			return valI > valJ

		case FieldTypeDateTime:
			valI, okI := cellI.AsTime()
			valJ, okJ := cellJ.AsTime()
			if okI && okJ {
				if ascending {
					return valI.Before(valJ)
				}
				return valI.After(valJ)
			}
			// Fallback to string comparison
			strI := fmt.Sprintf("%v", cellI.Value)
			strJ := fmt.Sprintf("%v", cellJ.Value)
			if ascending {
				return strI < strJ
			}
			return strI > strJ

		case FieldTypeCheckbox:
			valI, _ := cellI.AsBool()
			valJ, _ := cellJ.AsBool()
			if ascending {
				return !valI && valJ
			}
			return valI && !valJ

		default: // Text types
			strI := fmt.Sprintf("%v", cellI.Value)
			strJ := fmt.Sprintf("%v", cellJ.Value)
			if ascending {
				return strI < strJ
			}
			return strI > strJ
		}
	})

	// Rebuild the tree with sorted rows
	t.Root.Children = make([]*Node, 0, len(rows))
	for _, row := range rows {
		row.parent = nil
		t.Root.AddChild(row)
	}

	// Update row numbers
	t.updateRowNumbers()
}

// updateRowNumbers updates the RowNumber field for all nodes.
func (t *TreeTable) updateRowNumbers() {
	rowNum := 0
	for node := range t.dataNodes() {
		node.RowNumber = rowNum
		rowNum++
	}
}

// SumIf sums values in a column where another column matches a value.
func (t *TreeTable) SumIf(filterColumn, filterValue, sumColumn string) float64 {
	total := 0.0
	for node := range t.dataNodes() {
		filterCell := node.GetCell(filterColumn, t)
		sumCell := node.GetCell(sumColumn, t)

		if filterCell != nil && sumCell != nil {
			if fmt.Sprint(filterCell.Value) == filterValue {
				if val, ok := ToFloat(sumCell.Value); ok {
					total += val
				}
			}
		}
	}
	return total
}

// ToFloat converts a value to float64.
func ToFloat(val any) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	default:
		if f, err := strconv.ParseFloat(fmt.Sprintf("%v", val), 64); err == nil {
			return f, true
		}
		return 0, false
	}
}

// GroupBy 按指定列分组
func (t *TreeTable) GroupBy(columnName string) error {
	// 获取所有行
	allRows := t.AllRows()
	if len(allRows) == 0 {
		return nil
	}

	// 按分组列的值排序，便于分组
	sort.Slice(allRows, func(i, j int) bool {
		cellI := allRows[i].GetCell(columnName, t)
		cellJ := allRows[j].GetCell(columnName, t)

		if cellI == nil && cellJ == nil {
			return false
		}
		if cellI == nil {
			return true
		}
		if cellJ == nil {
			return false
		}

		return fmt.Sprintf("%v", cellI.Value) < fmt.Sprintf("%v", cellJ.Value)
	})

	// 创建新的根容器
	newRoot := NewContainerNode("root", nil)
	newRoot.isOpen = true

	// 分组处理
	currentGroup := ""
	var currentGroupContainer *Node

	for _, row := range allRows {
		cell := row.GetCell(columnName, t)
		groupValue := "未分组"
		if cell != nil {
			groupValue = fmt.Sprintf("%v", cell.Value)
		}

		// 如果分组值改变，创建新的分组容器
		if groupValue != currentGroup {
			currentGroup = groupValue
			currentGroupContainer = NewContainerNode("group", []CellData{
				{ColumnName: columnName, Value: groupValue, Type: FieldTypeSingleLineText},
			})
			currentGroupContainer.GroupKey = groupValue
			currentGroupContainer.isOpen = true
			newRoot.AddChild(currentGroupContainer)
		}

		// 将行添加到当前分组
		row.parent = nil
		currentGroupContainer.AddChild(row)
	}

	// 更新根节点
	t.Root = newRoot
	t.OriginalRoot = newRoot.Clone()

	return nil
}

// Aggregate 对分组进行聚合计算
func (t *TreeTable) Aggregate(groupColumn, targetColumn, aggType string) map[string]float64 {
	result := make(map[string]float64)

	// 遍历所有分组容器
	for _, node := range t.Root.Children {
		if node.IsContainer() && strings.HasPrefix(node.Type, "group") {
			groupKey := node.GroupKey
			if groupKey == "" {
				if cell := node.GetCell(groupColumn, t); cell != nil {
					groupKey = fmt.Sprintf("%v", cell.Value)
				} else {
					groupKey = "未分组"
				}
			}

			var aggregateValue float64
			count := 0

			// 遍历分组内的所有行
			for row := range node.Walk() {
				if row.IsContainer() {
					continue
				}

				if cell := row.GetCell(targetColumn, t); cell != nil {
					if val, ok := ToFloat(cell.Value); ok {
						switch aggType {
						case "sum":
							aggregateValue += val
						case "avg":
							aggregateValue += val
							count++
						case "max":
							if count == 0 || val > aggregateValue {
								aggregateValue = val
							}
							count = 1
						case "min":
							if count == 0 || val < aggregateValue {
								aggregateValue = val
							}
							count = 1
						}
					}
				}
			}

			// 处理平均值
			if aggType == "avg" && count > 0 {
				aggregateValue /= float64(count)
			}

			if aggType == "count" {
				rowCount := 0
				for range node.Walk() {
					if !node.IsContainer() {
						rowCount++
					}
				}
				result[groupKey] = float64(rowCount)
			} else {
				result[groupKey] = aggregateValue
			}
		}
	}

	return result
}

// GetGroups 获取所有分组
func (t *TreeTable) GetGroups() []*Node {
	var groups []*Node
	for _, child := range t.Root.Children {
		if child.IsContainer() && strings.HasPrefix(child.Type, "group") {
			groups = append(groups, child)
		}
	}
	return groups
}

// ExpandAllGroups 展开所有分组
func (t *TreeTable) ExpandAllGroups() {
	for node := range t.Root.Walk() {
		if node.IsContainer() {
			node.isOpen = true
		}
	}
}

// CollapseAllGroups 折叠所有分组
func (t *TreeTable) CollapseAllGroups() {
	for node := range t.Root.Walk() {
		if node.IsContainer() && !strings.HasPrefix(node.Type, "root") {
			node.isOpen = false
		}
	}
}

// Ungroup 取消分组，回到平面结构
func (t *TreeTable) Ungroup() {
	newRoot := NewContainerNode("root", nil)
	newRoot.isOpen = true

	// 将所有行提取到根节点
	for _, node := range t.Root.Children {
		for row := range node.Walk() {
			if !row.IsContainer() {
				row.parent = nil
				newRoot.AddChild(row)
			}
		}
	}

	t.Root = newRoot
	t.OriginalRoot = newRoot.Clone()
}

// ToMarkdown exports the table to Markdown format.
func (t *TreeTable) ToMarkdown() string {
	var sb strings.Builder
	sb.WriteString("# Tree Table Structure\n\n")
	sb.WriteString("| Level | Type |")
	for _, col := range t.Columns {
		sb.WriteString(fmt.Sprintf(" %s |", col.Name))
	}
	sb.WriteString("\n|-------|------|")
	for range t.Columns {
		sb.WriteString("-------|")
	}
	sb.WriteString("\n")

	for node := range t.dataNodes() {
		relativeDepth := node.Depth() - 1
		indent := strings.Repeat("&nbsp;&nbsp;&nbsp;", relativeDepth)

		icon := "📄"
		if node.IsContainer() {
			if node.isOpen {
				icon = "📂"
			} else {
				icon = "📁"
			}
		}

		sb.WriteString(fmt.Sprintf("| %s%s | %s |", indent, icon, node.Type))

		for _, col := range t.Columns {
			cell := node.GetCell(col.Name, t)
			value := "-"
			if cell != nil {
				value = fmt.Sprintf("%v", cell.Value)
			}
			sb.WriteString(fmt.Sprintf(" %s |", value))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// ToJSON exports the table to JSON.
func (t *TreeTable) ToJSON() ([]byte, error) {
	type exportData struct {
		Columns []ColumnDefinition `json:"columns"`
		Root    *Node              `json:"root"`
	}
	return json.MarshalIndent(exportData{t.Columns, t.Root}, "", "  ")
}

// FromJSON imports the table from JSON.
func FromJSON(data []byte) (*TreeTable, error) {
	var d struct {
		Columns []ColumnDefinition `json:"columns"`
		Root    *Node              `json:"root"`
	}
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, err
	}

	table := &TreeTable{
		Root:         d.Root,
		OriginalRoot: d.Root.Clone(),
		Columns:      d.Columns,
	}
	table.initColumnMap()

	return table, nil
}

// 从字符串值探测数据类型
func detectTypeFromString(s string) FieldType {
	// 尝试解析为布尔值
	if strings.EqualFold(s, "true") || strings.EqualFold(s, "false") ||
		s == "1" || s == "0" || s == "是" || s == "否" {
		return FieldTypeCheckbox
	}

	// 尝试解析为整数
	if _, err := strconv.Atoi(s); err == nil {
		return FieldTypeNumber
	}

	// 尝试解析为浮点数
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return FieldTypeNumber
	}

	// 尝试解析为日期时间 (RFC3339格式)
	if _, err := time.Parse(time.RFC3339, s); err == nil {
		return FieldTypeDateTime
	}

	// 尝试解析为简单日期格式 (YYYY-MM-DD)
	if _, err := time.Parse("2006-01-02", s); err == nil {
		return FieldTypeDateTime
	}

	// 尝试解析为时间格式 (HH:MM:SS)
	if _, err := time.Parse("15:04:05", s); err == nil {
		return FieldTypeDateTime
	}

	// 尝试解析为日期时间组合 (YYYY-MM-DD HH:MM:SS)
	if _, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return FieldTypeDateTime
	}

	// 检查是否是URL
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return FieldTypeURL
	}

	// 检查是否是电子邮件
	if strings.Contains(s, "@") && strings.Contains(s, ".") {
		return FieldTypeEmail
	}

	// 检查是否是电话号码 (简单验证)
	if len(s) >= 7 && len(s) <= 15 {
		allDigits := true
		for _, r := range s {
			if r < '0' || r > '9' {
				if r != '+' && r != '-' && r != '(' && r != ')' && r != ' ' {
					allDigits = false
					break
				}
			}
		}
		if allDigits {
			return FieldTypePhone
		}
	}

	// 检查是否是多行文本 (包含换行符)
	if strings.Contains(s, "\n") {
		return FieldTypeMultiLineText
	}

	// 默认作为单行文本
	return FieldTypeSingleLineText
}

func inferType(v any) FieldType {
	if v == nil {
		return FieldTypeSingleLineText
	}

	switch val := v.(type) {
	case string:
		return detectTypeFromString(val)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return FieldTypeNumber
	case float32, float64:
		return FieldTypeNumber
	case bool:
		return FieldTypeCheckbox
	case time.Time:
		return FieldTypeDateTime
	default:
		return detectTypeFromString(fmt.Sprintf("%v", v))
	}
}

var nilAny any = nil
