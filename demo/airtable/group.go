package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TreeNode 表示树形表格中的节点
type TreeNode struct {
	Level    int        // 层级深度（0表示根节点，1表示一级子节点等）
	Text     string     // 节点文本
	Children []TreeNode // 子节点
	Expanded bool       // 是否展开
	Data     []string   // 其他列的数据
}

// TreeTable 自定义树形表格组件
type TreeTable struct {
	*tview.Box
	Root         TreeNode // 根节点
	Columns      []string // 列标题
	ColumnWidths []int    // 计算后的列宽
	Selected     int      // 选中的行索引
	Offset       int      // 滚动偏移
}

// NewTreeTable 创建新的树形表格
func NewTreeTable() *TreeTable {
	return &TreeTable{
		Box:      tview.NewBox(),
		Columns:  []string{"Name", "Department", "Age"},
		Selected: -1,
	}
}

// SetRoot 设置根节点
func (tt *TreeTable) SetRoot(root TreeNode) *TreeTable {
	tt.Root = root
	tt.calculateColumnWidths()
	return tt
}

// calculateColumnWidths 计算各列的最佳宽度
func (tt *TreeTable) calculateColumnWidths() {
	// 初始化列宽，至少为列标题的宽度
	tt.ColumnWidths = make([]int, len(tt.Columns))
	for i, col := range tt.Columns {
		tt.ColumnWidths[i] = len(col)
	}

	// 遍历所有节点，更新列宽
	var traverse func(node TreeNode)
	traverse = func(node TreeNode) {
		// 第一列的宽度需要考虑缩进
		indentedWidth := node.Level*2 + len(node.Text)
		if indentedWidth > tt.ColumnWidths[0] {
			tt.ColumnWidths[0] = indentedWidth
		}

		// 其他列的宽度
		for i, value := range node.Data {
			if i+1 < len(tt.ColumnWidths) && len(value) > tt.ColumnWidths[i+1] {
				tt.ColumnWidths[i+1] = len(value)
			}
		}

		// 递归处理子节点
		if node.Expanded {
			for _, child := range node.Children {
				traverse(child)
			}
		}
	}

	traverse(tt.Root)

	// 添加一些边距
	for i := range tt.ColumnWidths {
		tt.ColumnWidths[i] += 2
	}
}

// Draw 绘制树形表格
func (tt *TreeTable) Draw(screen tcell.Screen) {
	tt.Box.Draw(screen)
	x, y, width, height := tt.GetInnerRect()

	// 绘制列标题
	colX := x
	for i, col := range tt.Columns {
		colWidth := tt.ColumnWidths[i]
		if colX+colWidth > x+width {
			colWidth = x + width - colX
		}
		if colWidth > 0 {
			tview.Print(screen, col, colX, y, colWidth, tview.AlignLeft, tcell.ColorYellow)
		}
		colX += colWidth
	}

	// 绘制分隔线
	y++
	for i := x; i < x+width; i++ {
		screen.SetContent(i, y, tcell.RuneHLine, nil, tcell.StyleDefault.Foreground(tcell.ColorGray))
	}
	y++

	// 绘制数据行
	visibleNodes := tt.getVisibleNodes()
	startIdx := max(0, min(tt.Offset, len(visibleNodes)-height+2))
	endIdx := min(startIdx+height-2, len(visibleNodes))

	for i := startIdx; i < endIdx; i++ {
		node := visibleNodes[i]
		rowY := y + i - startIdx

		// 高亮选中的行
		style := tcell.StyleDefault
		if i == tt.Selected {
			style = style.Reverse(true)
		}

		// 绘制每一列
		colX := x
		for colIdx := 0; colIdx < len(tt.Columns); colIdx++ {
			colWidth := tt.ColumnWidths[colIdx]
			if colX+colWidth > x+width {
				colWidth = x + width - colX
			}

			var text string
			if colIdx == 0 {
				// 第一列显示树形结构
				indent := strings.Repeat("  ", node.Level)
				prefix := "├─ "
				if len(node.Children) > 0 {
					if node.Expanded {
						prefix = "▼ "
					} else {
						prefix = "► "
					}
				} else {
					prefix = "• "
				}
				text = indent + prefix + node.Text
			} else if colIdx-1 < len(node.Data) {
				// 其他列显示数据
				text = node.Data[colIdx-1]
			}

			if colWidth > 0 {
				tview.Print(screen, text, colX, rowY, colWidth, tview.AlignLeft, style.GetUnderlineColor())
			}
			colX += colWidth
		}
	}
}

// getVisibleNodes 获取所有可见的节点（平铺列表）
func (tt *TreeTable) getVisibleNodes() []TreeNode {
	var result []TreeNode
	var traverse func(node TreeNode)

	traverse = func(node TreeNode) {
		result = append(result, node)
		if node.Expanded {
			for _, child := range node.Children {
				traverse(child)
			}
		}
	}

	traverse(tt.Root)
	return result
}

// InputHandler 处理输入事件
func (tt *TreeTable) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return tt.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		visibleNodes := tt.getVisibleNodes()

		switch event.Key() {
		case tcell.KeyUp:
			if tt.Selected > 0 {
				tt.Selected--
			}
		case tcell.KeyDown:
			if tt.Selected < len(visibleNodes)-1 {
				tt.Selected++
			}
		case tcell.KeyEnter:
			if tt.Selected >= 0 && tt.Selected < len(visibleNodes) {
				node := visibleNodes[tt.Selected]
				if len(node.Children) > 0 {
					// 切换展开/折叠状态
					tt.toggleNode(&tt.Root, tt.Selected)
					tt.calculateColumnWidths()
				}
			}
		}
	})
}

// toggleNode 切换节点的展开/折叠状态
func (tt *TreeTable) toggleNode(root *TreeNode, index int) bool {
	var traverse func(node *TreeNode, current *int) bool
	traverse = func(node *TreeNode, current *int) bool {
		if *current == index {
			node.Expanded = !node.Expanded
			return true
		}
		*current++

		if node.Expanded {
			for i := range node.Children {
				if found := traverse(&node.Children[i], current); found {
					return true
				}
			}
		}
		return false
	}

	current := 0
	return traverse(root, &current)
}

// Focus 使组件获得焦点
func (tt *TreeTable) Focus(delegate func(p tview.Primitive)) {
	tt.Box.Focus(delegate)
}

// HasFocus 检查组件是否有焦点
func (tt *TreeTable) HasFocus() bool {
	return tt.Box.HasFocus()
}

func main() {
	app := tview.NewApplication()

	// 创建示例数据
	root := TreeNode{
		Level:    0,
		Text:     "Departments",
		Expanded: true,
		Children: []TreeNode{
			{
				Level:    1,
				Text:     "Engineering",
				Expanded: true,
				Data:     []string{"Engineering", "3 employees"},
				Children: []TreeNode{
					{
						Level:    2,
						Text:     "Alice",
						Data:     []string{"Engineering", "30"},
						Children: nil,
					},
					{
						Level:    2,
						Text:     "Bob",
						Data:     []string{"Engineering", "28"},
						Children: nil,
					},
					{
						Level:    2,
						Text:     "Diana",
						Data:     []string{"Engineering", "32"},
						Children: nil,
					},
				},
			},
			{
				Level:    1,
				Text:     "Marketing",
				Expanded: true,
				Data:     []string{"Marketing", "2 employees"},
				Children: []TreeNode{
					{
						Level:    2,
						Text:     "Charlie",
						Data:     []string{"Marketing", "35"},
						Children: nil,
					},
					{
						Level:    2,
						Text:     "Eve",
						Data:     []string{"Marketing", "29"},
						Children: nil,
					},
				},
			},
		},
	}

	// 创建树形表格
	treeTable := NewTreeTable()
	treeTable.SetRoot(root).
		SetBorder(true).
		SetTitle("Employee Tree Table")

	if err := app.SetRoot(treeTable, true).Run(); err != nil {
		panic(err)
	}
}
