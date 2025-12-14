package sdk

import (
	"github.com/ddkwork/golibrary/std/stream"
	"github.com/ddkwork/ux/demo/erp/gongshi/sdk/field"
)

// ColumnDefinition defines a column in the table.
type ColumnDefinition struct {
	Name         string          // Column name (unique identifier)
	Type         field.FieldType // Data type
	Formula      string          // Column formula text (stores Go code!)
	Options      []string        // Options (for single/multiple select)
	IsDisabled   bool            // Whether editing is disabled
	Width        int             // Column width in pixels
	DefaultValue any             // Default value
	Values       []any           // Initial values for all cells in the column (for batch initialization)
}

// ColumnConfig 列配置（简化版）
type ColumnConfig struct {
	Name     string          // 列名
	Type     field.FieldType // 数据类型
	Formula  string          // 公式（可选）
	Options  []string        // 选项（可选）
	Width    int             // 宽度（可选）
	Disabled bool            // 是否禁用（可选）
}

func (t *TreeTable) CopyColumn() string {
	//if t.header.clickedColumnIndex < 0 {
	//	gtx.Execute(clipboard.WriteCmd{Data: io.NopCloser(strings.NewReader("t.header.clickedColumnIndex < 0 "))})
	//	return "t.header.clickedColumnIndex < 0 "
	//}
	g := stream.NewGeneratedFile()
	g.P("var columnData = []string{")
	//g.P(strconv.Quote(t.header.columnCells[t.header.clickedColumnIndex].Value), ",")
	//for i, datum := range TransposeMatrix(t.rows) {
	//	if i == t.header.clickedColumnIndex {
	//		g.P(strconv.Quote(datum.Value), ",")
	//	}
	//}
	g.P("}")
	//gtx.Execute(clipboard.WriteCmd{Data: io.NopCloser(strings.NewReader(g.Format()))})
	return g.String()
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
	for node := range t.DataNodes() {
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
	for node := range t.DataNodes() {
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
	for node := range t.DataNodes() {
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
	for node := range t.DataNodes() {
		for i := range node.RowCells {
			if node.RowCells[i].ColumnName == colName {
				// Update cell type if needed
				if node.RowCells[i].Type.Valid() {
					node.RowCells[i].Type = colDef.Type
				}
				break
			}
		}
	}
	return true
}
