package sdk

import (
	"sync"

	"github.com/ddkwork/golibrary/std/mylog"
)

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
	filteredRows     []*Node
	groupedRows      []*Node
	rootRows         []*Node
}

// NewTreeTable creates a new TreeTable instance.
func NewTreeTable() *TreeTable {
	return &TreeTable{
		Root:      NewContainerNode("root", nil),
		columnMap: make(map[string]*ColumnDefinition),
	}
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
			var nilAny any = nil
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

func (t *TreeTable) RootRows() []*Node {
	if t.groupedRows != nil {
		return t.groupedRows
	}
	if t.filteredRows != nil {
		return t.filteredRows
	}
	if t.Root == nil {
		mylog.CheckNil(t.OriginalRoot)
		t.Root = t.OriginalRoot
	}
	mylog.CheckNil(t.Root)
	t.rootRows = t.Root.Children
	return t.rootRows
}

func (t *TreeTable) RootRowCount() int {
	return len(t.RootRows())
}
