package sdk

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"github.com/ddkwork/golibrary/std/stream"
)

func (t *TreeTable) IsRowSelected() bool { return t.SelectedNode != nil }

// AllRows returns all row nodes (depth-first traversal, skipping root).
func (t *TreeTable) AllRows() []*Node {
	return t.dataNodesSlice()
}

// updateRowNumbers updates the RowNumber field for all nodes.
func (t *TreeTable) updateRowNumbers() {
	rowNum := 0
	for node := range t.dataNodes() {
		node.RowNumber = rowNum
		rowNum++
	}
}
func (t *TreeTable) updateRowNumber(parent *Node, currentRowIndex int) int {
	parent.RowNumber = currentRowIndex
	currentRowIndex++
	if parent.CanHaveChildren() {
		for _, child := range parent.Children {
			currentRowIndex = t.updateRowNumber(child, currentRowIndex)
		}
	}
	return currentRowIndex
}

// GetRow 通过行索引获取行节点
func (t *TreeTable) GetRow(rowIndex int) *Node {
	rows := t.AllRows()
	if rowIndex < 0 || rowIndex >= len(rows) {
		return nil
	}
	return rows[rowIndex]
}

// RowCount returns the number of rows.
func (t *TreeTable) RowCount() int {
	count := 0
	for range t.dataNodes() {
		count++
	}
	return count
}

// CountTableRows returns the number of table rows, including all descendants, whether open or not.
func CountTableRows(rows []*Node) int {
	count := len(rows)
	for _, row := range rows {
		if row.CanHaveChildren() {
			count += CountTableRows(row.Children)
		}
	}
	return count
}
func (n *Node) CopyRow(gtx layout.Context, widths []int) string {
	g := stream.NewGeneratedFile()
	g.WriteString("var rowData = []string{ ")
	cells := n.RowCells
	for i, cell := range cells {
		g.WriteString(fmt.Sprintf("%-*s", widths[i], strconv.Quote(cell.AsString())) + ",")
	}
	g.P("}")
	gtx.Execute(clipboard.WriteCmd{Data: io.NopCloser(strings.NewReader(g.String()))})
	return g.String()
}
