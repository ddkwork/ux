package sdk

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
