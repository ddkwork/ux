package sdk

import (
	"image/color"

	"github.com/ddkwork/golibrary/std/stream/uuid"
)

// TableModel provides access to the root nodes of the table's data underlying model.
type TableModel[T TableRowConstraint] interface {
	RootRowCount() int
	RootRows() []T
	SetRootRows(rows []T)
}

// TableRowData provides information about a single row of data.
type TableRowData interface {
	ID() uuid.ID
	Parent() *Node
	SetParent(parent *Node)
	CanHaveChildren() bool
	Children() []*Node
	SetChildren(children []*Node)
	CellDataForSort(col int) string
	ColumnCell(row, col int, foreground, background color.Color, selected, indirectlySelected, focused bool) any
	IsOpen() bool
	SetOpen(open bool)
}

// TableRowConstraint defines the constraints required of the data type used for data rows in tables.
type TableRowConstraint interface {
	//comparable
	TableRowData
}

type SimpleTableModel struct {
	roots []*Node
}

func (m *SimpleTableModel) RootRowCount() int {
	return len(m.roots)
}

func (m *SimpleTableModel) RootRows() []*Node {
	return m.roots
}

func (m *SimpleTableModel) SetRootRows(rows []*Node) {
	m.roots = rows
}
