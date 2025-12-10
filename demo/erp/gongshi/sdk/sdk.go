package sdk

import "fmt"

type MemoryTable struct {
	rows []map[string]interface{}
}

func NewMemoryTable() *MemoryTable {
	return &MemoryTable{rows: make([]map[string]interface{}, 0)}
}

func (t *MemoryTable) AddRow(row map[string]interface{}) {
	if _, ok := row["计算结果"]; !ok {
		row["计算结果"] = 0.0
	}
	t.rows = append(t.rows, row)
}

func (t *MemoryTable) RowCount() int {
	return len(t.rows)
}

func (t *MemoryTable) GetRow(i int) map[string]interface{} {
	if i < 0 || i >= len(t.rows) {
		return nil
	}
	row := make(map[string]interface{})
	for k, v := range t.rows[i] {
		row[k] = v
	}
	return row
}

func (t *MemoryTable) SetValue(i int, col string, val interface{}) error {
	if i < 0 || i >= len(t.rows) {
		return fmt.Errorf("索引越界")
	}
	t.rows[i][col] = val
	return nil
}

func (t *MemoryTable) SumIf(field string, crit interface{}, sumField string) float64 {
	var sum float64
	for _, r := range t.rows {
		if v, ok := r[field]; ok && fmt.Sprint(v) == fmt.Sprint(crit) {
			if sv, ok := r[sumField].(float64); ok {
				sum += sv
			} else if iv, ok := r[sumField].(int); ok {
				sum += float64(iv)
			}
		}
	}
	return sum
}
