package main

import (
	"fmt"
	"image"
	"image/color"
	"slices"
	"sort"
	"strings"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"gioui.org/x/richtext"
	"github.com/ddkwork/ux"
)

// 3. 表头排序
// 4. row右键菜单
// 5. 右键复制行或者列cells到剪切板,复制列是在表头的各自单元格右键复制,遍历得到
// 6. 单击row取出row结构体数据，并绘制选中行的背景色
// 7. 双击弹出from布局编辑row结构体View
// 8. 绘制列分割线
// 9. 可拖动列宽
// 10. 奇偶行斑马线背景色

type (
	ClickAction func(node *Node)
	TreeTable   struct {
		Children     []*Node
		root         *Node //？ how to use it?
		selectedNode *Node

		maxIndentWidth           int         // 层级列单元格最小宽度
		SelectionChangedCallback ClickAction // 行选中回调
		DoubleClickCallback      ClickAction // double click callback
		LongPressCallback        ClickAction // mobile long press callback

		widget.List

		headerButtons   []*widget.Clickable // 存储每列的 clickable 状态
		sortColumn      int                 // 当前排序的列索引
		sortAscending   bool                // 是否升序排序
		filterText      string
		filteredRows    []*Node
		onMenuItemClick func(tr *Node, item string)
	}
)

func NewTreeTable() *TreeTable {
	return &TreeTable{
		Children:                 nil,
		root:                     nil,
		selectedNode:             nil,
		SelectionChangedCallback: nil,
		List: widget.List{
			Scrollbar: widget.Scrollbar{},
			List: layout.List{
				Axis:        layout.Vertical,
				ScrollToEnd: false,
				Alignment:   0,
				Position:    layout.Position{},
			},
		},
	}
}

func (t *TreeTable) OnClick(fun ClickAction) *TreeTable {
	t.SelectionChangedCallback = fun
	return t
}

func (t *TreeTable) OnNodeDoubleClick(fun ClickAction) *TreeTable {
	t.DoubleClickCallback = fun
	return t
}

func (t *TreeTable) SetWidth(width int) *TreeTable {
	t.maxIndentWidth = width
	return t
}

func (t *TreeTable) SetRootRows(rootRows []*Node) *TreeTable {
	for _, row := range rootRows {
		t.expandNode(row) // 将每个节点展开
	}
	// 计算最大缩进宽度
	maxIndentWidth := t.calculateMaxIndentWidth(rootRows)
	t.maxIndentWidth = maxIndentWidth * 7
	t.Children = rootRows
	return t
}

// 递归函数，获取树形结构中的最大深度
func (t *TreeTable) calculateMaxDepth(node *Node) int {
	if node == nil {
		return 0
	}

	maxDepth := 0
	for _, child := range node.Children {
		childDepth := t.calculateMaxDepth(child)
		if childDepth > maxDepth {
			maxDepth = childDepth
		}
	}

	return maxDepth + 1 // 当前节点深度加1
}

// 计算所有根节点中的最大缩进宽度
func (t *TreeTable) calculateMaxIndentWidth(nodes []*Node) int {
	maxDepth := 0
	for _, node := range nodes {
		depth := t.calculateMaxDepth(node) // 递归获取每个节点的深度
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	return maxDepth * 8 // 假定每个层级产生8单位的缩进
}

func (t *TreeTable) expandNode(node *Node) {
	node.expanded = true // 设置节点为展开状态
	for _, child := range node.Children {
		t.expandNode(child) // 递归展开子节点
	}
}

type CallbackFun func(node *Node)

type ColumnInfo struct {
	ID           int
	CurrentWidth float32
	MiniWidth    float32
	MaxiWidth    float32
	Cell         string
}

type Node struct {
	RowCells          []ColumnInfo
	cells             []*widget.Clickable // 单元格单击事件
	Icon              *widget.Icon
	Children          []*Node
	expanded          bool
	selected          bool
	clickable         *widget.Clickable
	CellClickCallback CallbackFun
	parent            *Node

	MenuOptions   []string
	contextAreas  []component.ContextArea
	contextMenu   component.MenuState
	menuButtons   []*widget.Clickable
	menuList      widget.List
	menuInit      bool
	isMenuVisible bool
}

func SetParents(children []*Node, parent *Node) {
	for _, child := range children {
		child.parent = parent
		if len(child.Children) > 0 {
			SetParents(child.Children, child)
		}
	}
}

func (n *Node) Depth() int {
	count := 0
	p := n.parent
	for p != nil {
		count++
		p = p.parent
	}
	return count
}

func (n *Node) IsContainer() bool {
	return n.Children != nil && len(n.Children) > 0
}

var th = ux.ThemeDefault()

func (t *TreeTable) Layout(gtx layout.Context) layout.Dimensions {
	list := material.List(th.Theme, &t.List)
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.renderHeader(gtx) // 渲染表头
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return list.Layout(gtx, len(t.Children), func(gtx layout.Context, index int) layout.Dimensions {
				return t.renderNode(gtx, t.Children[index], index)
			})
		}),
	)
}

func (t *TreeTable) SortNodes() {
	if len(t.Children) == 0 || t.sortColumn >= len(t.Children[0].RowCells) {
		return // 如果没有子节点或者列索引无效，直接返回
	}

	sort.Slice(t.Children, func(i, j int) bool {
		cellI := t.Children[i].RowCells[t.sortColumn].Cell
		cellJ := t.Children[j].RowCells[t.sortColumn].Cell
		if t.sortAscending {
			return cellI < cellJ
		}
		return cellI > cellJ
	})
}

func (t *TreeTable) SortNodes2() {
	if len(t.Children) == 0 {
		return // 如果没有子节点，直接返回
	}

	sort.Slice(t.Children, func(i, j int) bool {
		cellI := t.Children[i].RowCells[t.sortColumn].Cell
		cellJ := t.Children[j].RowCells[t.sortColumn].Cell
		if t.sortAscending {
			return cellI < cellJ
		}
		return cellI > cellJ
	})
}

func (t *TreeTable) SortNodes1() {
	sort.Slice(t.Children, func(i, j int) bool {
		// 根据当前列索引和排序顺序进行比较
		if t.sortAscending {
			return t.Children[i].RowCells[t.sortColumn].Cell < t.Children[j].RowCells[t.sortColumn].Cell
		}
		return t.Children[i].RowCells[t.sortColumn].Cell > t.Children[j].RowCells[t.sortColumn].Cell
	})
}

func (t *TreeTable) Filter(text string) {
	t.filterText = text

	if text == "" {
		t.filteredRows = make([]*Node, 0)
		return
	}

	items := make([]*Node, 0)
	for i, item := range t.Children {
		if strings.Contains(item.RowCells[i].Cell, text) {
			items = append(items, item)
		}

		for i, child := range item.Children {
			if strings.Contains(child.RowCells[i].Cell, text) {
				items = append(items, child)
			}
		}
	}

	t.filteredRows = items
}

// 第一列和第二列缩进正常，第三列开始就有两个单元格对不齐，后面的就跟着对不齐了，找下b，改u哪g句代码
// todo 泛型绑定结构体
func (t *TreeTable) renderHeader(gtx layout.Context) layout.Dimensions { // 渲染表头
	return layout.Inset{}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// 使用 Flex 布局来显示表头
		return layout.Flex{
			Axis:    layout.Horizontal,
			Spacing: 0,
			// Alignment: layout.Middle,
			WeightSum: 0,
		}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				// gtx.Constraints.Min.X = gtx.Dp(t.maxIndentWidth) //todo not work
				return t.headerCell(gtx, "列 1", 1)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return t.headerCell(gtx, "depth", 2)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return t.headerCell(gtx, "列 3", 3)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return t.headerCell(gtx, "列 4", 4)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return t.headerCell(gtx, "列 5", 5)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return t.headerCell(gtx, "列 6", 6)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return t.headerCell(gtx, "列 7", 7)
			}),
		)
	})
}

func (t *TreeTable) headerCell(gtx layout.Context, title string, colIndex int) layout.Dimensions {
	colIndex--
	if t.headerButtons == nil {
		t.headerButtons = make([]*widget.Clickable, len(t.Children[0].RowCells))
	}
	if t.headerButtons[colIndex] == nil {
		t.headerButtons[colIndex] = new(widget.Clickable) // 动态初始化
	}
	clickable := t.headerButtons[colIndex]
	btn := material.Button(th.Theme, clickable, title)
	btn.Background = color.NRGBA{R: 64, G: 64, B: 64, A: 255}
	btn.Color = White
	btn.CornerRadius = unit.Dp(0)
	btn.Inset.Top = unit.Dp(4)
	btn.Inset.Bottom = unit.Dp(4)
	if clickable.Hovered() {
		btn.Background = color.NRGBA{R: 150, G: 150, B: 150, A: 255} // 鼠标悬停颜色
	}
	drawColumnDivider(gtx, colIndex+1, DividerFg)
	if clickable.Clicked(gtx) { // 使用传入的 gtx
		// 更新排序列
		if t.sortColumn == colIndex {
			t.sortAscending = !t.sortAscending // 切换升序/降序
		} else {
			t.sortColumn = colIndex // 更新排序列
			t.sortAscending = true  // 设为升序
		}
		t.SortNodes() // 调用排序函数
		println("Header cell clicked. Sorting on column:", colIndex)
	}
	return clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return btn.Layout(gtx) // 使用按钮的布局
	})
}

func (t *TreeTable) renderNode(gtx layout.Context, node *Node, rowIndex int) layout.Dimensions {
	if !node.menuInit {
		node.menuInit = true
		node.contextMenu = component.MenuState{
			Options: func() []func(gtx layout.Context) layout.Dimensions {
				out := make([]func(gtx layout.Context) layout.Dimensions, 0, len(node.MenuOptions))
				node.menuButtons = make([]*widget.Clickable, 0, len(node.MenuOptions))
				for i, opt := range node.MenuOptions {
					node.menuButtons = append(node.menuButtons, new(widget.Clickable))
					if opt == "-" {
						out = append(out, component.Divider(th.Theme).Layout)
						continue
					}
					out = append(out, component.MenuItem(th.Theme, node.menuButtons[i], opt).Layout)
				}
				return out
			}(),
		}
	}

	for i := range node.menuButtons {
		if node.menuButtons[i].Clicked(gtx) {
			if t.onMenuItemClick != nil {
				t.onMenuItemClick(node, node.MenuOptions[i])
			}
		}
	}

	// 确定行背景颜色，偶数行和奇数行不同
	bgColor := ux.RowColor(rowIndex)
	isContainer := len(node.Children) > 0

	if node.clickable == nil {
		node.clickable = &widget.Clickable{}
	}

	for { // todo bug
		break
		click, ok := node.clickable.Update(gtx)
		if !ok {
			break
		}
		switch click.NumClicks {
		case 1:
			//if node.CellClickedCallback != nil { //todo rename
			//	go node.CellClickedCallback(node)
			//	gtx.Execute(op.InvalidateCmd{})
			//}

			if node.clickable.Clicked(gtx) {
				node.expanded = !node.expanded // 切换展开状态
				t.selectedNode = node          // 记录被点击的节点
				if node.CellClickCallback != nil {
					node.CellClickCallback(node) // 调用节点的点击回调
				}
				if t.SelectionChangedCallback != nil {
					t.SelectionChangedCallback(node) // 处理树的全局点击事件
				}
			}

		case 2:
			if t.DoubleClickCallback != nil { // todo rename
				go t.DoubleClickCallback(node)
				gtx.Execute(op.InvalidateCmd{})
			}
		default:
			if node.Children == nil {
				continue
			}
		}
	}

	// 仅在isMenuVisible为真时绘制菜单
	if node.isMenuVisible {
		t.drawContextArea(gtx, node)
	}

	if node.clickable.Clicked(gtx) {
		node.expanded = !node.expanded // 切换展开状态
		t.selectedNode = node          // 记录被点击的节点
		if node.CellClickCallback != nil {
			node.CellClickCallback(node) // 调用节点的点击回调
		}
		if t.SelectionChangedCallback != nil {
			t.SelectionChangedCallback(node) // 处理树的全局点击事件
		}
	}

	if node.clickable.Hovered() {
		bgColor = th.Color.InputFocusedBgColor // TreeHoveredBgColor
	}

	const baseIndent = 8
	HierarchyInsert := layout.Inset{Left: unit.Dp((node.Depth() + 1) * baseIndent), Top: 1}

	var rowCells []layout.FlexChild
	// 绘制层级图标，虽然非容器节点没有图标，但是需要绘制深度+1的空白图标占位来缩进层级
	icon := ArrowRightIcon
	if node.expanded {
		icon = ArrowDownIcon
	}
	rowCells = append(rowCells, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return node.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			if !isContainer {
				return HierarchyInsert.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{}
				})
			}
			return HierarchyInsert.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return reDrawIcon(gtx, icon)
			})
		})
	}))

	// 绘制层级列文本
	HierarchyIndent := 0 // 层级图标和层级文本聚拢,视觉这样才是ok的
	if !isContainer {
		HierarchyIndent = (node.parent.Depth())*baseIndent - iconSize // 层级列非容器节点的文本和父节点文本对齐算法,即图标左侧的占用宽度
	}
	// 绘制层级列（固定宽度）
	rowCells = append(rowCells, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Left: unit.Dp(HierarchyIndent)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Dp(unit.Dp(t.maxIndentWidth)) // 限制层级列最小宽度
			return node.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				richText := ux.NewRichText()
				richText.AddSpan(richtext.SpanStyle{
					// Font:        font.Font{},
					Size:        unit.Sp(12),
					Color:       Orange100,
					Content:     node.RowCells[0].Cell,
					Interactive: false,
				})
				return layoutInsetCell(gtx, richText.Layout)
			})
		})
	}))

	// 绘制其他列
	for i, cell := range node.RowCells[1:] {
		noHierarchyColumIndent := -((node.Depth()) * baseIndent) + iconSize // 非层级列容器节点单元格负缩进,todo handle DividerWidth?
		if i > 0 && node.Depth() < t.calculateMaxDepth(node) {
			// todo remove this condition?
			// 另一种方案是copy grid的取单元格平均坐标宽度的代码,然后强制宽度为平均宽度,这样就能保证单元格对齐。不过这似乎有点复杂，需要加入几个update函数
			// 此外，要对齐表头也只能移动坐标偏移这个办法了，所以先搞右键菜单和排序bug修复，以及unison的n叉树的增删改查功能，然后再考虑表头对齐的问题
			noHierarchyColumIndent -= i + 1
			if node.Depth() == 0 && i > 1 {
				noHierarchyColumIndent -= i + 2
			}
			if noHierarchyColumIndent < -t.maxIndentWidth {
				noHierarchyColumIndent = -t.maxIndentWidth // 确保不会过度缩进
			}
		}
		if !isContainer { // 非层级列是正常的对齐算法
			noHierarchyColumIndent = -(node.parent.Depth())*baseIndent + iconSize // 非层级列非容器节点单元格负缩进算法，和层级列非容器节点是对称的
		}
		rowCells = append(rowCells, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(noHierarchyColumIndent)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions { // 添加缩进
				if node.cells == nil {
					node.cells = make([]*widget.Clickable, len(node.RowCells))
				}
				clickable := node.cells[i]
				if clickable == nil {
					clickable = &widget.Clickable{}
					node.cells[i] = clickable
				}
				return node.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return material.Clickable(gtx, clickable, func(gtx layout.Context) layout.Dimensions {
						drawColumnDivider(gtx, i+1, DividerFg) // 为每列绘制列分隔条
						richText := ux.NewRichText()
						richText.AddSpan(richtext.SpanStyle{
							// Font:        font.Font{},
							Size:        unit.Sp(12),
							Color:       Yellow200,
							Content:     cell.Cell,
							Interactive: false,
						})
						insetCell := layoutInsetCell(gtx, richText.Layout)

						// contextMenu todo
						if node.MenuOptions == nil {
							node.MenuOptions = []string{
								"add",
								"delete",
								"edit",
								"copy",
								"cut",
								"paste",
								"select all",
								"invert selection",
								"expand all",
								"collapse all",
								"refresh",
								"help",
							}
						}

						// 先检查上下文菜单只在右键单击时显示
						if !node.isMenuVisible {
							return insetCell
						}

						return layout.Stack{}.Layout(gtx,
							layout.Stacked(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{}.Layout(gtx,
									layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
										node.menuList.Axis = layout.Vertical
										return material.List(th.Theme, &node.menuList).Layout(gtx, len(node.MenuOptions), func(gtx layout.Context, index int) layout.Dimensions {
											if len(node.contextAreas) < index+1 {
												node.contextAreas = append(node.contextAreas, component.ContextArea{})
											}
											// state := &node.contextAreas[index]
											return layout.Stack{}.Layout(gtx,
												layout.Stacked(func(gtx layout.Context) layout.Dimensions {
													gtx.Constraints.Min.X = gtx.Constraints.Max.X
													return layout.UniformInset(unit.Dp(8)).Layout(gtx, material.Body1(th.Theme, node.MenuOptions[index]).Layout)
												}),
												layout.Expanded(func(gtx layout.Context) layout.Dimensions {
													//return state.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
													//	gtx.Constraints.Min.X = 0
													//	return component.Menu(th.Theme, &node.contextMenu).Layout(gtx)
													//})

													evt, ok := gtx.Source.Event(pointer.Filter{
														Target: node,
														Kinds:  pointer.Press | pointer.Release,
													})
													if !ok {
														return D{}
													}
													e, ok := evt.(pointer.Event)
													if !ok {
														return D{}
													}
													if e.Buttons == pointer.ButtonSecondary {
														if e.Kind == pointer.Press {
															// cellMenu.active = true
														} else if e.Kind == pointer.Release {
															// cellMenu.active = false
														}
														// if cellMenu.active {
														return node.contextAreas[index].Layout(gtx, func(gtx C) D {
															// m.rowIdx = row
															// m.colIdx = col
															gtx.Constraints.Max.X = 500
															gtx.Constraints.Max.Y = 1400
															return t.drawContextArea(gtx, node)
														})
														//}
													}
													return D{}
												}),
											)
										})
									}),
								)
							}),
						)
					})
				})
			})
		}))
	}

	//component.Resize{
	//	Axis:  0,
	//	Ratio: 0,
	//}.Layout(gtx,
	//	func(gtx layout.Context) layout.Dimensions { //当前列
	//
	//	},
	//	func(gtx layout.Context) layout.Dimensions { //下一列
	//
	//	},
	//	func(gtx layout.Context) layout.Dimensions { //分隔条
	//		rect := image.Rectangle{
	//			Max: image.Point{
	//				X: gtx.Dp(unit.Dp(4)),
	//				Y: gtx.Constraints.Max.Y,
	//			},
	//		}
	//		paint.FillShape(gtx.Ops, color.NRGBA{A: 200}, clip.Rect(rect).Op())
	//		return D{Size: rect.Max}
	//	},
	//)

	// 创建行背景和行高
	items := []layout.FlexChild{
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Background{}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				if t.selectedNode == node {
					bgColor = ColorPink
				}
				paint.FillShape(gtx.Ops, bgColor, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Min.Y)}.Op())
				return layout.Dimensions{Size: gtx.Constraints.Min}
			}, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(22)) // 行高
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, rowCells...)
			})
		}),
	}

	// 递归渲染子节点
	if node.expanded && isContainer {
		for _, child := range node.Children {
			items = append(items, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				rowIndex++
				return t.renderNode(gtx, child, rowIndex)
			}))
		}
	}
	return layout.Flex{
		Axis:      layout.Vertical,
		Spacing:   0,
		Alignment: layout.Middle,
		WeightSum: 0,
	}.Layout(gtx, items...)
}

func SetDepth(nodes []*Node, depth int) {
	for _, node := range nodes {
		node.RowCells = slices.Insert(node.RowCells, 1, ColumnInfo{Cell: fmt.Sprintf(" (Depth: %d)", depth)})
		SetDepth(node.Children, depth+1)
	}
}

// 在初始化时调用这个函数
func init() {
	SetDepth(TestRootRows, 0)
}

func (t *TreeTable) drawContextArea(gtx C, node *Node) D {
	return layout.Center.Layout(gtx, func(gtx C) D { // 重置min x y 到0，并根据max x y 计算弹出菜单的合适大小
		// mylog.Struct("todo",gtx.Constraints)
		menuStyle := component.Menu(th.Theme, &node.contextMenu)
		menuStyle.SurfaceStyle = component.SurfaceStyle{
			Theme: th.Theme,
			ShadowStyle: component.ShadowStyle{
				CornerRadius: 18, // 弹出菜单的椭圆角度
				Elevation:    0,
				// AmbientColor:  color.NRGBA(colornames.Blue400),
				// PenumbraColor: color.NRGBA(colornames.Blue400),
				// UmbraColor:    color.NRGBA(colornames.Blue400),
			},
			Fill: color.NRGBA{R: 50, G: 50, B: 50, A: 255}, // 弹出菜单的背景色
		}
		return menuStyle.Layout(gtx)
	})
}

func layoutInsetCell(gtx layout.Context, cell layout.Widget) layout.Dimensions {
	return layout.Inset{
		Top:    4, // 文本居中，drawColumnDivider需要设置tallestHeight := gtx.Dp(unit.Dp(32))增加高度避免虚线
		Bottom: 0,
		Left:   8,
		Right:  8,
	}.Layout(gtx, cell)
}

const iconSize = 12

func reDrawIcon(gtx layout.Context, icon *widget.Icon) layout.Dimensions {
	size := gtx.Dp(iconSize)
	point := image.Pt(size, size)
	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			radius := size / 2
			return component.Rect{Color: White, Size: point, Radii: radius}.Layout(gtx)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = image.Point{X: size}
			icon.Layout(gtx, Black)
			return layout.Dimensions{
				Size: image.Pt(size, size),
			}
		}),
	)
}

const DividerWidth = 1

// 分隔线绘制函数
func drawColumnDivider(gtx layout.Context, col int, color color.NRGBA) {
	if col > 0 {
		dividerWidth := DividerWidth
		// 使用默认的行高作为分隔线的高度
		tallestHeight := gtx.Dp(unit.Dp(32)) // 或其他合适的高度
		// 或者这里可以从单元格的内容获取最大高度
		if gtx.Constraints.Min.Y > tallestHeight {
			tallestHeight = gtx.Constraints.Min.Y
		}
		stack3 := clip.Rect{Max: image.Pt(dividerWidth, tallestHeight)}.Push(gtx.Ops)
		paint.Fill(gtx.Ops, color)
		stack3.Pop()
	}
}
