package main

import (
	"fmt"
	"strings"

	"gioui.org/unit"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"github.com/ddkwork/golibrary/stream/align"
	"github.com/ddkwork/ux"
)

// CellData 表示单元格的数据结构
type CellData struct {
	Text                   string
	Depth                  unit.Dp
	IsLastChild            bool
	maxColumnCellTextWidth unit.Dp
	maxColumnCellText      string
}

// Node 表示树节点
type Node struct {
	Header                      *Node
	RowCells                    []CellData
	Children                    []*Node
	Parent                      *Node
	maxLevelColumnCellTextWidth unit.Dp
	maxColumnCellWidth          unit.Dp
}

// Global constants
const (
	indentBase        = unit.Dp(1) // todo test
	columnHeaderCount = 3
)

// 初始化根节点
var root = NewRoot()

// 主程序入口
func main() {
	root.drawRows(root.MaxColumnCellTextWidths(ux.TransposeMatrix(collectRows(root))))
}

// collectRows 遍历树形结构并收集所有行数据，填充每个单元格的深度信息
func collectRows(node *Node) [][]CellData {
	var rows [][]CellData

	// 将当前节点的行单元格添加到 rows 中，同时设置深度
	rowCells := make([]CellData, len(node.RowCells))
	for i, cell := range node.RowCells {
		cellData := CellData{
			Text:        cell.Text,
			Depth:       node.Depth(),
			IsLastChild: false,
			// maxColumnCellTextWidth: 0,
			// maxColumnCellText:      "",
		}
		if i == len(node.RowCells)-1 {
			cellData.IsLastChild = true
		}
		rowCells[i] = cellData
	}
	rows = append(rows, rowCells)

	// 递归处理每个子节点
	for _, child := range node.Children {
		childRows := collectRows(child)
		rows = append(rows, childRows...)
	}
	return rows
}

func (n *Node) Depth() unit.Dp {
	if n.Parent != nil {
		return n.Parent.Depth() + 1
	}
	return 0 // 根节点深度为 0
}

func (n *Node) MaxDepth() unit.Dp {
	depth := n.Depth()
	for _, child := range n.Children {
		childDepth := child.Depth()
		if childDepth > depth {
			depth = childDepth
		}
	}
	return depth
}

func (n *Node) MaxColumnCellTextWidths(columns [][]CellData) []unit.Dp {
	columns = append(columns, n.Header.RowCells)
	columnWidths := make([]unit.Dp, columnHeaderCount)
	for rowIndex, column := range columns {
		for columnIndex, cell := range column {
			if columnIndex == columnHeaderCount { // todo bug  0- (columnHeaderCount-1) == columnIndex
				break
			}
			width := align.StringWidth(cell.Text) // 这里应该是导致非层级右边距太大的原因
			if width > columnWidths[columnIndex] {
				columnWidths[columnIndex] = width

				// 这个很重要，用于精确计算最大列单元格宽度
				columns[rowIndex][columnIndex].maxColumnCellTextWidth = width
				columns[rowIndex][columnIndex].maxColumnCellText = cell.Text
			}
		}
	}
	n.maxLevelColumnCellTextWidth = columns[0][0].maxColumnCellTextWidth // 层级列最大单元格文本宽度
	n.maxColumnCellWidth = n.MaxColumnCellWidth()                        // 层级列最大单元格宽度,限制宽度对齐
	return columnWidths
}

// NewRoot 创建并返回根节点
func NewRoot() *Node {
	rootNode := &Node{
		Header:   nil,
		RowCells: createHeaderCells(),
		Children: nil,
		Parent:   nil,
	}
	rootNode.Header = rootNode

	const topLevelRowsToMake = 100
	for i := 0; i < topLevelRowsToMake; i++ {
		row := createRow(i + 1)
		if i%10 == 3 {
			row.Children = createChildRows()
		}
		rootNode.AddChild(row)
	}
	return rootNode
}

// createHeaderCells 创建表头单元格
func createHeaderCells() []CellData {
	return []CellData{
		{Text: "Name"},
		{Text: "Value"},
		{Text: "Description"},
	}
}

// createRow 创建新行
func createRow(index int) *Node {
	return &Node{
		RowCells: []CellData{
			{Text: fmt.Sprintf("Row %d", index)},
			{Text: fmt.Sprintf("Value %d", index)},
			{Text: fmt.Sprintf("Some longer content for Row %d", index)},
		},
	}
}

// createChildRows 创建子行
func createChildRows() []*Node {
	childRows := make([]*Node, 5)
	for j := 0; j < 5; j++ {
		childRows[j] = &Node{
			RowCells: []CellData{
				{Text: fmt.Sprintf("Sub Row %d", j+1)},
				{Text: fmt.Sprintf("Sub Value %d", j+1)},
				{Text: fmt.Sprintf("Sub Description %d", j+1)},
			},
		}
	}
	return childRows
}

// AddChild 将子节点添加到当前节点的子节点列表中
func (n *Node) AddChild(child *Node) {
	child.Parent = n
	n.Children = append(n.Children, child)
}

// drawRows 方法中的层级列绘制部分
func (n *Node) drawRows1(rows [][]CellData, maxColumnCellTextWidths []int) {
	buf := stream.NewBuffer("")

	// 绘制表头
	buf.WriteStringLn("─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────")
	buf.WriteString("│")
	buf.WriteString("       ")
	// writeRow(buf, n.Header.RowCells, maxColumnCellTextWidths) // 传入 maxColumnCellTextWidths
	buf.NewLine()
	buf.WriteStringLn("─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────")

	// 绘制行数据
	const (
		indent          = "│   "
		childPrefix     = "├───"
		lastChildPrefix = "└───"
	)
	for _, row := range rows[1:] {
		// 只处理第一列作为层级列
		b := stream.NewBuffer("")
		levelIndent := strings.Repeat(" ", int(row[0].Depth))
		// b.WriteString("│")         // 往缓冲区添加一条分隔线
		b.WriteString(levelIndent) // 增加层级缩进
		// if row[0].Depth > 0 {
		b.WriteString("├───")
		//} else {
		indent := maxColumnCellTextWidths[0] - b.Len()
		b.WriteString(strings.Repeat(" ", indent))
		//}
		//mylog.Trace("层级列实际宽度", b.Len())
		//mylog.Trace("预期层级列宽度", maxColumnCellTextWidths[0])
		if b.Len() != maxColumnCellTextWidths[0] {
			// panic("层级列宽度校验失败b.Len() != maxColumnCellTextWidths[0]")
		}
		buf.WriteString(b.String())
		buf.Indent(1)

		// 绘制非层级列
		// writeRow(buf, row[1:], maxColumnCellTextWidths[1:]) //绘制非层级列
		buf.NewLine()
	}

	mylog.Json("RootRows", buf.String())
}

// 设置常量
const (
	indent          = "│   "
	childPrefix     = "├───"
	lastChildPrefix = "└───"
)

func (n *Node) MaxColumnCellWidth() unit.Dp {
	HierarchyIndent := unit.Dp(1)
	DividerWidth := align.StringWidth(" │ ")    // todo test
	iconWidth := align.StringWidth(childPrefix) // todo test
	return n.MaxDepth()*HierarchyIndent +       // 最大深度的左缩进
		iconWidth + // 图标宽度,不管深度是多少，每一行都只会有一个层级图标
		n.maxLevelColumnCellTextWidth + 5 + //(8 * 2) + 20 + // 左右padding,20是sort图标的宽度或者容器节点求和的文本宽度
		DividerWidth // 列分隔条宽度
}

func (n *Node) drawHeader(maxColumnCellTextWidths []unit.Dp) *stream.Buffer {
	buf := stream.NewBuffer("")

	all := n.maxColumnCellWidth
	for _, width := range maxColumnCellTextWidths[1:] {
		all += width
	}
	all += align.StringWidth("│")*unit.Dp(len(maxColumnCellTextWidths)) + 1 // 最后一个分隔符的宽度

	buf.WriteStringLn("┌─" + strings.Repeat("─", int(all)))
	buf.WriteString("│")

	// 计算每个单元格的左边距
	for i, cell := range n.Header.RowCells {
		paddedText := fmt.Sprintf("%-*s", int(maxColumnCellTextWidths[i]), cell.Text) // 左对齐填充

		// 添加左边距，仅在首列进行处理，依据列宽计算
		if i == 0 {
			buf.WriteString(strings.Repeat(" ", int(n.maxColumnCellWidth-maxColumnCellTextWidths[i]-1))) // -1是分隔符的空间
		}

		buf.WriteString(paddedText)
		if i < len(n.Header.RowCells)-1 {
			buf.WriteString(" │ ") // 在每个单元格之间添加分隔符
		}
	}

	buf.NewLine()
	buf.WriteStringLn("├─" + strings.Repeat("─", int(all)))
	return buf
}

func (n *Node) drawRows(maxColumnCellTextWidths []unit.Dp) {
	buf := n.drawHeader(maxColumnCellTextWidths)
	n.printChildren(buf, n.Children, "", maxColumnCellTextWidths) // 传入子节点打印函数
	mylog.Json("RootRows", buf.String())
}

func (n *Node) printChildren(out *stream.Buffer, children []*Node, parentPrefix string, maxColumnCellTextWidths []unit.Dp) {
	lastIdx := len(children) - 1
	for i, child := range children {
		prefix := "├──" //"├───"
		if i == lastIdx {
			prefix = "╰──" //"└───"
		}

		// 打印层级列
		leftIndent := parentPrefix + prefix
		rightIndent := ""
		if align.StringWidth(leftIndent)+align.StringWidth(child.RowCells[0].Text) < n.maxColumnCellWidth {
			rightIndent = strings.Repeat(".", int(n.maxColumnCellWidth-align.StringWidth(leftIndent)-align.StringWidth(child.RowCells[0].Text)))
		}

		// 打印层级列内容
		fmt.Fprint(out, leftIndent, child.RowCells[0].Text+rightIndent)

		// 打印非层级列，确保 "│" 在每行的开头
		for j := 1; j < columnHeaderCount; j++ {
			if j < len(child.RowCells) {
				paddedText := fmt.Sprintf("%-*s", int(maxColumnCellTextWidths[j]), child.RowCells[j].Text)
				// 确保非层级列前加上"│"
				fmt.Fprintf(out, " │ %s", paddedText)
			}
		}
		out.NewLine()

		// 递归打印子节点，确保目录符号在前
		if len(child.Children) > 0 {
			indent := strings.Repeat(" ", int(child.Depth()*indentBase)) // 根据 indentBase 确定缩进
			// 确保新前缀的目录符号在最前面
			newParentPrefix := "│" + parentPrefix + indent
			n.printChildren(out, child.Children, newParentPrefix, maxColumnCellTextWidths)
		}
	}
}
