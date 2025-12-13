package sdk

import "encoding/json"

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
