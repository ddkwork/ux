package sdk

import (
	"fmt"
	"sort"

	"github.com/ddkwork/ux/demo/erp/gongshi/sdk/field"
)

// SortByColumn sorts all rows by the specified column.
func (t *TreeTable) SortByColumn(colName string, ascending bool) {
	rows := t.AllRows()

	// Sort rows based on column values
	sort.Slice(rows, func(i, j int) bool {
		cellI := rows[i].GetCell(colName)
		cellJ := rows[j].GetCell(colName)

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
		case field.NumberType, field.CurrencyType, field.PercentType:
			valI := cellI.AsFloat()
			valJ := cellJ.AsFloat()
			if ascending {
				return valI < valJ
			}
			return valI > valJ

		case field.DateTimeType:
			valI := cellI.AsTime()
			valJ := cellJ.AsTime()
			if ascending {
				return valI.Before(valJ)
			}
			return valI.After(valJ)
			// Fallback to string comparison
			//strI := fmt.Sprintf("%v", cellI.Value)
			//strJ := fmt.Sprintf("%v", cellJ.Value)
			//if ascending {
			//	return strI < strJ
			//}
			//return strI > strJ

		case field.CheckboxType:
			valI := cellI.AsBool()
			valJ := cellJ.AsBool()
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
