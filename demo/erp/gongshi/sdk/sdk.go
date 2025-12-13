package sdk

import (
	"encoding/json"
	"fmt"
	"iter"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ------------------------------ 字段类型定义（Airtable风格） ------------------------------
type FieldType string

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
)

// ------------------------------ 核心数据结构（支持公式） ------------------------------
type Node struct {
	ID        string     // 节点唯一ID（使用UUID）
	Type      string     // 节点类型（容器节点以"_container"结尾）
	RowCells  []CellData // 行数据（含公式列）
	Children  []*Node    // 子节点
	parent    *Node      // 父节点
	isOpen    bool       // 是否展开（仅容器节点有效）
	GroupKey  string     // 分组键
	RowNumber int        // 行号（用于排序）
}

type TreeTable struct {
	Root         *Node          // 根节点（虚拟容器）
	OriginalRoot *Node          // 原始根节点备份
	Columns      []CellData     // 表头定义（含公式列）
	columnMap    map[string]int // 列名到索引的映射
	SelectedNode *Node          // 当前选中节点
	once         sync.Once      // 一次性初始化标记

	// 回调函数
	OnRowSelected    func(n *Node)
	OnRowDoubleClick func(n *Node)
}

// CellData 单元格数据（增强类型安全）
type CellData struct {
	Name       string    // 列名（唯一标识）
	Value      any       // 单元格值（公式计算结果或手动输入值）
	Type       FieldType // 数据类型
	Formula    string    // 公式代码（Go代码片段）
	Options    []string  // 选项（用于单选/多选）
	IsDisabled bool      // 是否禁用编辑
	Width      int       // 列宽（像素）
	isHeader   bool      // 是否为表头单元格
}

// 类型安全的值获取方法
func (c *CellData) AsString() (string, bool) {
	v, ok := c.Value.(string)
	return v, ok
}

func (c *CellData) AsInt() (int, bool) {
	switch v := c.Value.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}

func (c *CellData) AsFloat() (float64, bool) {
	switch v := c.Value.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	default:
		return 0, false
	}
}

func (c *CellData) AsBool() (bool, bool) {
	v, ok := c.Value.(bool)
	return v, ok
}

// 判断是否为公式单元格
func (c *CellData) IsFormula() bool {
	return c.Type == FieldTypeFormula && c.Formula != ""
}

// ------------------------------ 节点方法（支持公式计算） ------------------------------
func NewNode(rowCells []CellData) *Node {
	return &Node{
		ID:        uuid.New().String(), // 使用UUID
		Type:      "node",
		RowCells:  rowCells,
		Children:  nil,
		parent:    nil,
		isOpen:    false,
		GroupKey:  "",
		RowNumber: 0,
	}
}

func NewContainerNode(typeKey string, rowCells []CellData) *Node {
	n := NewNode(rowCells)
	n.Type = typeKey + "_container"
	n.isOpen = true
	return n
}

// 克隆节点（深拷贝）
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

	// 复制行数据
	for i, cell := range n.RowCells {
		clone.RowCells[i] = cell
	}

	// 复制子节点
	for i, child := range n.Children {
		cloneChild := child.Clone()
		cloneChild.parent = clone
		clone.Children[i] = cloneChild
	}

	return clone
}

// 添加子节点
func (n *Node) AddChild(child *Node) {
	child.parent = n
	n.Children = append(n.Children, child)
}

// 批量添加子节点
func (n *Node) AddChildren(children []*Node) {
	for _, child := range children {
		n.AddChild(child)
	}
}

// 插入子节点到指定位置
func (n *Node) InsertChild(index int, child *Node) {
	if index < 0 || index > len(n.Children) {
		index = len(n.Children)
	}
	child.parent = n
	n.Children = append(n.Children[:index], append([]*Node{child}, n.Children[index:]...)...)
}

// 移除子节点
func (n *Node) RemoveChild(child *Node) {
	for i, c := range n.Children {
		if c.ID == child.ID {
			n.Children = append(n.Children[:i], n.Children[i+1:]...)
			return
		}
	}
}

// 判断是否为容器节点
func (n *Node) IsContainer() bool {
	return strings.HasSuffix(n.Type, "_container")
}

// 获取节点深度
func (n *Node) Depth() int {
	depth := 0
	for p := n.parent; p != nil; p = p.parent {
		depth++
	}
	return depth
}

// 使用iter迭代器遍历节点
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

// 获取指定列名的单元格（自动计算公式列值）
func (n *Node) GetCell(colName string, table *TreeTable) *CellData {
	for i := range n.RowCells {
		if n.RowCells[i].Name == colName {
			cell := &n.RowCells[i]
			// 如果是公式列且值未计算，则执行公式计算
			if cell.IsFormula() {
				// table.calculateFormulaCell(n, cell)
			}
			return cell
		}
	}
	return nil
}

// 设置单元格值（允许设置公式列）
func (n *Node) SetCellValue(colName string, value any, table *TreeTable) {
	for i := range n.RowCells {
		if n.RowCells[i].Name == colName {
			cell := &n.RowCells[i]
			cell.Value = value
			if cell.Type == "" {
				cell.Type = inferType(value)
			}
			return
		}
	}
	// 列不存在则新增
	colDef := table.GetColumn(colName)
	if colDef != nil {
		n.RowCells = append(n.RowCells, CellData{
			Name:  colName,
			Value: value,
			Type:  colDef.Type,
		})
	} else {
		n.RowCells = append(n.RowCells, CellData{
			Name:  colName,
			Value: value,
			Type:  inferType(value),
		})
	}
}

// 从字符串值探测数据类型
func DetectTypeFromString(s string) FieldType {
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

// 推断值类型（使用FieldType）
func inferType(v any) FieldType {
	switch val := v.(type) {
	case string:
		return DetectTypeFromString(val)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return FieldTypeNumber
	case float32, float64:
		return FieldTypeNumber
	case bool:
		return FieldTypeCheckbox
	case time.Time:
		return FieldTypeDateTime
	default:
		// 尝试将其他类型转换为字符串再检测
		return DetectTypeFromString(fmt.Sprintf("%v", v))
	}
}

// 检测单元格值的类型
func (c *CellData) DetectValueType() FieldType {
	return inferType(c.Value)
}

// 在节点级别检测列类型
func (n *Node) DetectCellType(colName string) FieldType {
	cell := n.GetCell(colName, nil)
	if cell == nil {
		return FieldTypeSingleLineText
	}
	return cell.DetectValueType()
}

// ------------------------------ 行列增删改查方法 ------------------------------
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

// 获取所有数据节点（直接从根节点的子节点开始遍历）
func (t *TreeTable) dataNodes() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		// 遍历根节点的所有直接子节点
		for _, child := range t.Root.Children {
			// 使用栈实现深度优先遍历
			stack := []*Node{child}
			for len(stack) > 0 {
				// 弹出栈顶节点
				n := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				// 处理当前节点
				if !yield(n) {
					return
				}

				// 将子节点逆序入栈（保证从左到右的顺序）
				for i := len(n.Children) - 1; i >= 0; i-- {
					stack = append(stack, n.Children[i])
				}
			}
		}
	}
}

// 获取所有数据节点的索引迭代器
func (t *TreeTable) dataNodesIndexed() iter.Seq2[int, *Node] {
	return func(yield func(int, *Node) bool) {
		idx := 0
		for node := range t.dataNodes() {
			if !yield(idx, node) {
				return
			}
			idx++
		}
	}
}

// 获取所有数据节点的切片
func (t *TreeTable) dataNodesSlice() []*Node {
	var nodes []*Node
	for node := range t.dataNodes() {
		nodes = append(nodes, node)
	}
	return nodes
}

// 获取所有行节点（深度优先遍历，跳过根节点）
func (t *TreeTable) AllRows() []*Node {
	return t.dataNodesSlice()
}

// 获取所有行节点（迭代器版本，跳过根节点）
func (t *TreeTable) AllRows2() iter.Seq2[int, *Node] {
	return t.dataNodesIndexed()
}

// 获取行数
func (t *TreeTable) RowCount() int {
	count := 0
	for range t.dataNodes() {
		count++
	}
	return count
}

// 获取列数
func (t *TreeTable) ColCount() int {
	return len(t.Columns)
}

// 获取列索引
func (t *TreeTable) ColIndex(colName string) int {
	if idx, ok := t.columnMap[colName]; ok {
		return idx
	}
	return -1
}

// 获取列名
func (t *TreeTable) ColName(colIndex int) string {
	if colIndex < 0 || colIndex >= len(t.Columns) {
		return ""
	}
	return t.Columns[colIndex].Name
}

// 初始化列映射
func (t *TreeTable) initColumnMap() {
	t.columnMap = make(map[string]int)
	for i, col := range t.Columns {
		t.columnMap[col.Name] = i
	}
}

// 获取默认值
func getDefaultValue(ft FieldType) any {
	switch ft {
	case FieldTypeNumber:
		return 0
	case FieldTypeCheckbox:
		return false
	case FieldTypeDateTime:
		return time.Now().Format(time.RFC3339)
	case FieldTypeSingleSelect, FieldTypeMultipleSelect:
		return ""
	default:
		return ""
	}
}

// 添加新列（增强版）
func (t *TreeTable) AddColumn(col CellData, index int) {
	if index < 0 || index > len(t.Columns) {
		index = len(t.Columns)
	}

	// 检查是否已存在同名列
	if _, exists := t.columnMap[col.Name]; exists {
		return // 或更新现有列
	}

	// 插入新列
	t.Columns = append(t.Columns[:index], append([]CellData{col}, t.Columns[index:]...)...)
	t.initColumnMap()

	// 为所有行添加新列的单元格
	for node := range t.dataNodes() {
		node.SetCellValue(col.Name, getDefaultValue(col.Type), t)
	}
}

// 删除列（增强版）
func (t *TreeTable) DeleteColumn(colName string) bool {
	idx := t.ColIndex(colName)
	if idx == -1 {
		return false
	}

	// 从列定义中删除
	t.Columns = append(t.Columns[:idx], t.Columns[idx+1:]...)
	t.initColumnMap()

	// 从所有行中删除该列的单元格
	for node := range t.dataNodes() {
		for i := len(node.RowCells) - 1; i >= 0; i-- {
			if node.RowCells[i].Name == colName {
				node.RowCells = append(node.RowCells[:i], node.RowCells[i+1:]...)
			}
		}
	}
	return true
}

// 重命名列
func (t *TreeTable) RenameColumn(oldName, newName string) bool {
	idx := t.ColIndex(oldName)
	if idx == -1 {
		return false
	}

	// 更新列定义
	t.Columns[idx].Name = newName

	// 更新列映射
	delete(t.columnMap, oldName)
	t.columnMap[newName] = idx

	// 更新所有行中的单元格名称
	for node := range t.dataNodes() {
		for i, cell := range node.RowCells {
			if cell.Name == oldName {
				node.RowCells[i].Name = newName
				break
			}
		}
	}
	return true
}

// 更新列属性（增强版）
func (t *TreeTable) UpdateColumn(colName string, updateFunc func(*CellData)) bool {
	idx := t.ColIndex(colName)
	if idx == -1 {
		return false
	}

	// 应用更新函数
	updateFunc(&t.Columns[idx])

	// 更新所有行中的单元格
	for node := range t.dataNodes() {
		for i := range node.RowCells {
			if node.RowCells[i].Name == colName {
				updateFunc(&node.RowCells[i])
				break
			}
		}
	}
	return true
}

// 获取列定义
func (t *TreeTable) GetColumn(colName string) *CellData {
	idx := t.ColIndex(colName)
	if idx == -1 {
		return nil
	}
	return &t.Columns[idx]
}

// 批量检测列类型
func (t *TreeTable) DetectColumnTypes() map[string]FieldType {
	typeMap := make(map[string]FieldType)

	for _, col := range t.Columns {
		// 收集该列所有非空值
		values := make([]any, 0)
		for node := range t.dataNodes() {
			if cell := node.GetCell(col.Name, t); cell != nil && cell.Value != nil {
				values = append(values, cell.Value)
			}
		}

		// 如果有值，检测最常见的类型
		if len(values) > 0 {
			typeCounts := make(map[FieldType]int)
			for _, val := range values {
				ft := inferType(val)
				typeCounts[ft]++
			}

			// 找出出现频率最高的类型
			maxCount := 0
			dominantType := FieldTypeSingleLineText
			for ft, count := range typeCounts {
				if count > maxCount {
					maxCount = count
					dominantType = ft
				}
			}

			typeMap[col.Name] = dominantType
		} else {
			// 没有数据时，使用列定义的类型或默认类型
			if col.Type != "" {
				typeMap[col.Name] = col.Type
			} else {
				typeMap[col.Name] = FieldTypeSingleLineText
			}
		}
	}

	return typeMap
}

// 自动检测并更新列类型
func (t *TreeTable) AutoDetectAndUpdateTypes() {
	typeMap := t.DetectColumnTypes()
	for colName, detectedType := range typeMap {
		if currentType := t.GetColumn(colName).Type; currentType != detectedType {
			t.UpdateColumn(colName, func(c *CellData) {
				c.Type = detectedType
			})
		}
	}
}

// 添加新行（使用表头定义）
func (t *TreeTable) AddRow(values map[string]any, parentID string, position int) *Node {
	cells := make([]CellData, 0, len(t.Columns))

	// 根据表头创建单元格
	for _, col := range t.Columns {
		value := values[col.Name]
		if value == nil {
			value = getDefaultValue(col.Type)
		}

		cells = append(cells, CellData{
			Name:  col.Name,
			Value: value,
			Type:  col.Type,
		})
	}

	return t.addRowWithCells(cells, parentID, position)
}

// 内部方法：使用预定义单元格添加行
func (t *TreeTable) addRowWithCells(cells []CellData, parentID string, position int) *Node {
	var parent *Node
	if parentID == "" {
		parent = t.Root
	} else {
		for node := range t.dataNodes() {
			if node.ID == parentID {
				parent = node
				break
			}
		}
		if parent == nil {
			parent = t.Root
		}
	}

	newNode := NewNode(cells)

	if position < 0 || position > len(parent.Children) {
		parent.AddChild(newNode)
	} else {
		parent.InsertChild(position, newNode)
	}

	return newNode
}

// 插入行
func (t *TreeTable) InsertRow(index int, cells []CellData) bool {
	rows := t.AllRows()
	if index < 0 || index > len(rows) {
		return false
	}

	// 找到插入位置对应的节点
	var targetNode *Node
	var parent *Node
	var posInParent int

	if index == len(rows) {
		// 插入到最后
		lastRow := rows[len(rows)-1]
		parent = lastRow.parent
		if parent == nil {
			parent = t.Root
		}
		posInParent = len(parent.Children)
	} else {
		targetNode = rows[index]
		parent = targetNode.parent
		if parent == nil {
			parent = t.Root
		}

		// 查找在父节点中的位置
		for i, child := range parent.Children {
			if child.ID == targetNode.ID {
				posInParent = i
				break
			}
		}
	}

	// 创建新节点并插入
	newNode := NewNode(cells)
	parent.InsertChild(posInParent, newNode)
	return true
}

// 删除行
func (t *TreeTable) DeleteRow(rowIndex int) bool {
	rows := t.AllRows()
	if rowIndex < 0 || rowIndex >= len(rows) {
		return false
	}

	nodeToDelete := rows[rowIndex]
	parent := nodeToDelete.parent
	if parent == nil {
		return false
	}

	parent.RemoveChild(nodeToDelete)
	return true
}

// 删除行（按ID）
func (t *TreeTable) DeleteRowByID(nodeID string) bool {
	for node := range t.dataNodes() {
		if node.ID == nodeID {
			parent := node.parent
			if parent == nil {
				return false
			}
			parent.RemoveChild(node)
			return true
		}
	}
	return false
}

// 移动行
func (t *TreeTable) MoveRow(fromIndex, toIndex int) bool {
	rows := t.AllRows()
	if fromIndex < 0 || fromIndex >= len(rows) || toIndex < 0 || toIndex >= len(rows) {
		return false
	}

	fromNode := rows[fromIndex]
	toNode := rows[toIndex]

	// 不能移动到自己的子树中
	if isDescendant(fromNode, toNode) {
		return false
	}

	// 从原位置移除
	fromParent := fromNode.parent
	if fromParent == nil {
		return false
	}
	fromParent.RemoveChild(fromNode)

	// 插入到新位置
	toParent := toNode.parent
	if toParent == nil {
		toParent = t.Root
	}

	// 查找在父节点中的位置
	pos := 0
	for i, child := range toParent.Children {
		if child.ID == toNode.ID {
			pos = i
			break
		}
	}

	if toIndex < fromIndex {
		toParent.InsertChild(pos, fromNode)
	} else {
		toParent.InsertChild(pos+1, fromNode)
	}

	return true
}

// 检查是否是后代节点
func isDescendant(ancestor, descendant *Node) bool {
	for node := descendant.parent; node != nil; node = node.parent {
		if node == ancestor {
			return true
		}
	}
	return false
}

// 复制行
func (t *TreeTable) CopyRow(rowIndex int) *Node {
	rows := t.AllRows()
	if rowIndex < 0 || rowIndex >= len(rows) {
		return nil
	}

	original := rows[rowIndex]
	cloned := original.Clone()

	// 添加到原位置之后
	parent := original.parent
	if parent == nil {
		parent = t.Root
	}

	// 查找在父节点中的位置
	pos := 0
	for i, child := range parent.Children {
		if child.ID == original.ID {
			pos = i + 1
			break
		}
	}

	parent.InsertChild(pos, cloned)
	return cloned
}

// 获取单元格值（通过行索引和列名）
func (t *TreeTable) GetCellValue(rowIndex int, colName string) (any, bool) {
	rows := t.AllRows()
	if rowIndex < 0 || rowIndex >= len(rows) {
		return nil, false
	}

	cell := rows[rowIndex].GetCell(colName, t)
	if cell == nil {
		return nil, false
	}
	return cell.Value, true
}

// 设置单元格值（通过行索引和列名）
func (t *TreeTable) SetCellValue(rowIndex int, colName string, value any) bool {
	rows := t.AllRows()
	if rowIndex < 0 || rowIndex >= len(rows) {
		return false
	}

	rows[rowIndex].SetCellValue(colName, value, t)
	return true
}

// 获取整行数据
func (t *TreeTable) GetRow(rowIndex int) []CellData {
	rows := t.AllRows()
	if rowIndex < 0 || rowIndex >= len(rows) {
		return nil
	}
	return rows[rowIndex].RowCells
}

// 设置整行数据
func (t *TreeTable) SetRow(rowIndex int, cells []CellData) bool {
	rows := t.AllRows()
	if rowIndex < 0 || rowIndex >= len(rows) {
		return false
	}

	rows[rowIndex].RowCells = cells
	return true
}

// 排序行
func (t *TreeTable) SortRows(colName string, ascending bool) {
	rows := t.AllRows()

	sort.Slice(rows, func(i, j int) bool {
		valI, okI := t.GetCellValue(i, colName)
		valJ, okJ := t.GetCellValue(j, colName)

		if !okI || !okJ {
			return false
		}

		// 尝试数值比较
		if numI, ok := ToFloat(valI); ok {
			if numJ, ok := ToFloat(valJ); ok {
				if ascending {
					return numI < numJ
				}
				return numI > numJ
			}
		}

		// 字符串比较
		strI := fmt.Sprintf("%v", valI)
		strJ := fmt.Sprintf("%v", valJ)

		if ascending {
			return strI < strJ
		}
		return strI > strJ
	})

	// 重建树结构（保持父子关系）
	t.rebuildTreeFromSortedRows(rows)
}

// 辅助函数：转换为float64
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

// 从排序后的行重建树
func (t *TreeTable) rebuildTreeFromSortedRows(sortedRows []*Node) {
	// 创建ID到节点的映射
	idMap := make(map[string]*Node)
	for _, row := range sortedRows {
		idMap[row.ID] = row
	}

	// 重建父子关系
	for _, row := range sortedRows {
		// 保存原始子节点
		children := row.Children
		row.Children = nil

		// 重新添加子节点（按顺序）
		for _, child := range children {
			if childNode, ok := idMap[child.ID]; ok {
				row.AddChild(childNode)
			}
		}
	}
}

// ------------------------------ 树形表格核心方法（含公式列支持） ------------------------------
func NewTreeTable() *TreeTable {
	table := &TreeTable{}

	// 默认表头（使用FieldType）
	defaultColumns := []CellData{
		{Name: "姓名", Type: FieldTypeSingleLineText, Width: 120},
		{Name: "出生年份", Type: FieldTypeNumber, Width: 100},
		{Name: "年龄", Type: FieldTypeFormula,
			Formula: `return 2024 - ctx["出生年份"].(int)`, Width: 80},
		{Name: "女工日结", Type: FieldTypeNumber, Width: 100},
		{Name: "计算结果", Type: FieldTypeFormula,
			Formula: `/* 公式逻辑 */`, Width: 120},
		{Name: "入职日期", Type: FieldTypeDateTime, Width: 120},
		{Name: "状态", Type: FieldTypeSingleSelect,
			Options: []string{"在职", "离职"}, Width: 80},
	}

	table.Columns = defaultColumns
	table.initColumnMap()

	// 创建根节点（虚拟容器）
	root := NewContainerNode("root", nil)
	table.Root = root
	table.OriginalRoot = root.Clone()

	// 添加示例数据（直接作为根节点的子节点）
	group1 := NewContainerNode("department", []CellData{
		{Name: "姓名", Value: "技术部", Type: FieldTypeSingleLineText},
	})
	root.AddChild(group1)

	emp1 := NewNode([]CellData{
		{Name: "姓名", Value: "张三", Type: FieldTypeSingleLineText},
		{Name: "出生年份", Value: 1990, Type: FieldTypeNumber},
		{Name: "女工日结", Value: 200.0, Type: FieldTypeNumber},
		{Name: "入职日期", Value: "2020-01-15", Type: FieldTypeDateTime},
	})
	group1.AddChild(emp1)

	root.AddChild(group1)

	// 设置回调
	table.OnRowSelected = func(n *Node) {}
	table.OnRowDoubleClick = func(n *Node) {}

	return table
}

// SumIf 方法（优化版）
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

// ------------------------------ Markdown渲染（显示公式计算结果） ------------------------------
func (t *TreeTable) ToMarkdown() string {
	var sb strings.Builder
	sb.WriteString("# 树形表格结构（含公式列）\n\n")
	sb.WriteString("| 层级 | 类型 |")
	for _, col := range t.Columns {
		sb.WriteString(fmt.Sprintf(" %s |", col.Name))
	}
	sb.WriteString("\n|------|------|")
	for range t.Columns {
		sb.WriteString("------|")
	}
	sb.WriteString("\n")

	// 使用新的遍历方法
	//idx := 0
	for node := range t.dataNodes() {
		// 计算相对深度（相对于根节点）
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

// ------------------------------ 数据导入/导出（支持公式列） ------------------------------
// 导出为JSON（含公式定义）
func (t *TreeTable) ToJSON() ([]byte, error) {
	type exportData struct {
		Columns []CellData `json:"columns"`
		Root    *Node      `json:"root"`
	}
	return json.MarshalIndent(exportData{t.Columns, t.Root}, "", "  ")
}

// 从JSON导入（恢复公式列）
func FromJSON(data []byte) (*TreeTable, error) {
	var d struct {
		Columns []CellData `json:"columns"`
		Root    *Node      `json:"root"`
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

// 并行处理所有行
func (t *TreeTable) ProcessRowsConcurrently(processFunc func(*Node)) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, runtime.NumCPU()) // 限制并发数

	for node := range t.dataNodes() {
		wg.Add(1)
		sem <- struct{}{} // 获取信号量

		go func(n *Node) {
			defer wg.Done()
			defer func() { <-sem }() // 释放信号量
			processFunc(n)
		}(node)
	}

	wg.Wait()
}

// 查找符合条件的节点
func (t *TreeTable) FindNodes(predicate func(*Node) bool) []*Node {
	var results []*Node
	for node := range t.dataNodes() {
		if predicate(node) {
			results = append(results, node)
		}
	}
	return results
}
