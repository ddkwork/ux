package ux

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"io"
	"iter"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"gioui.org/gesture"
	"gioui.org/io/key"

	"gioui.org/io/clipboard"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"github.com/ddkwork/golibrary/stream/align"
	"github.com/ddkwork/golibrary/stream/deepcopy"
	"github.com/ddkwork/golibrary/stream/uuid"
	"github.com/ddkwork/ux/icon"
	"github.com/ddkwork/ux/widget/material"
)

type (
	TreeTable[T any] struct {
		TableContext[T]                             // 实例化时传入的上下文
		Root                         *Node[T]       // 根节点,保存数据到json只需要调用它即可
		header                       tableHeader[T] // 表头
		rootRows                     []*Node[T]     // from root.children
		filteredRows                 []*Node[T]     // 过滤后的行
		SelectedNode                 *Node[T]       // 选中的节点,文件管理器外部自定义右键菜单增删改查文件需要通过它取出节点元数据结构体的文件路径字段，所以需要导出
		contextMenu                  *ContextMenu   // 行右键菜单,rootRows只需要一个
		columnCount                  int            // 列数
		maxDepth                     unit.Dp        // 最大层级深度,todo add Cache like maxHierarchyColumnWidthCache
		maxColumnCellWidths          []unit.Dp      // 最宽的列label文本宽度for单元格
		maxColumnTextWidths          []unit.Dp      // 最宽的列文本宽度for tui
		maxHierarchyColumnWidthCache unit.Dp
		// rows                    [][]CellData        // 矩阵置换参数，行转为列，增删改节点后重新生成它
		columns                 [][]CellData        // CopyColumn
		DragRemovedRowsCallback func(n *Node[T])    // Called whenever a drag removes one or more rows from a model, but only if the source and destination tables were different.
		DropOccurredCallback    func(n *Node[T])    // Called whenever a drop occurs that modifies the model.
		inLayoutHeader          bool                // for drag
		columnResizeStart       unit.Dp             //
		columnResizeBase        unit.Dp             //
		columnResizeOverhead    unit.Dp             //
		preventUserColumnResize bool                //
		awaitingSyncToModel     bool                //
		wasDragged              bool                //
		dividerDrag             bool                //
		LongPressCallback       func(node *Node[T]) `json:"-"` // 长按回调
		pressStarted            time.Time           // 按压开始时间
		longPressed             bool                // 是否已经触发长按事件
		widget.List                                 // 为rootRows渲染列表和滚动条
		once                    sync.Once           // 自动计算列宽一次
	}
	TableContext[T any] struct {
		ContextMenuItems       func(gtx layout.Context, n *Node[T]) (items []ContextMenuItem) // 通过SelectedNode传递给菜单的do取出元数据，比如删除文件,但是菜单是否绘制取决于当前渲染的行，所以要传递n给can
		MarshalRowCells        func(n *Node[T]) (cells []CellData)                            // 序列化节点元数据
		UnmarshalRowCells      func(n *Node[T], values []string)                              // 节点编辑后反序列化回节点
		RowSelectedCallback    func()                                                         // 行选中回调,通过SelectedNode传递给菜单
		RowDoubleClickCallback func()                                                         // double click callback,通过SelectedNode传递给菜单
		LongPressCallback      func()                                                         // mobile long press callback,通过SelectedNode传递给菜单
		SetRootRowsCallBack    func()                                                         // 实例化所有节点回调,必要时调用root节点辅助操作
		JsonName               string                                                         // 保存序列化树形表格到文件的文件名
		IsDocument             bool                                                           // 是否生成markdown文档
	}
	tableHeader[T any] struct {
		sortOrder          sortOrder    // 排序方式
		sortedBy           int          // 排序列索引
		drags              []tableDrag  // 拖动表头列参数
		columnCells        []CellData   // 每列的最大深度，宽度,也可以把它看做一行
		clickedColumnIndex int          // 被点击的列索引
		manualWidthSet     []bool       // 新增状态标志数组，记录列是否被手动调整
		sortAscending      bool         // 升序还是降序
		contextMenu        *ContextMenu // 右键菜单，实现复制列数据到剪贴板
	}

	//func Layout(gtx layout.Context, icon any, color color.NRGBA, size unit.Dp) layout.Dimensions {
	//	sizeDp := gtx.Dp(size)
	//	if icon != nil {
	//		gtx.Constraints.Min = image.Point{X: sizeDp}
	//		switch v := icon.(type) {
	//		case *giosvg.Icon:
	//			v.Layout(gtx)
	//		case *widget.Icon:
	//			v.Layout(gtx, color)
	//		case *widget.Image:
	//			v.Layout(gtx)
	//		case image.Image:
	//			icon := &widget.Image{
	//				Src:      paint.NewImageOp(v),
	//				Fit:      widget.Unscaled,
	//				Position: layout.Center,
	//				Scale:    1.0, // todo 测试按钮图标和层级图标
	//			}
	//			icon.Layout(gtx)
	//		}
	//	}
	//	return layout.Dimensions{
	//		Size: image.Point{X: sizeDp, Y: sizeDp},
	//	}
	//}

	CellData struct {
		Text             string      // 单元格文本
		Tooltip          string      // 单元格提示信息
		Icon             any         // 单元格图标，格式支持：*giosvg.Icon, *widget.Icon, *widget.Image, image.Image
		FgColor          color.NRGBA // 单元格前景色,着色渲染单元格
		IsNasm           bool        // 是否是nasm汇编代码,为表头提供不同的着色渲染样式
		Disabled         bool        // 是否显示表头或者body单元格，或者禁止编辑节点时候使用
		columID          int         // 列id,预计后期用于区域选中,每行的列数和表头是一样的，而表头单元格在init的时候填充了每列的id，那么body rows的列id只要在操作row节点的时候遍历一下表头的row node填充列id即可
		rowID            int         // 行id,预计后期用于区域选中,由layout遍历list渲染的index填充该节点的rowID
		widget.Clickable             // todo 单元格点击事件,如果换成编辑框，那么编辑框需要支持单机事件和双击编辑回车事件,以及RichText高亮单元格
		RichText                     // todo 单元格富文本
	}
)

// NewTreeTable 原理:通过flex+适当的上下文控制既可以完美绘制一个树形表格
// 包括表头，父结点，子节点，右键菜单，通通flex
func NewTreeTable[T any](data T) *TreeTable[T] {
	rowCells := InitHeader(data)
	columnCount := len(rowCells)
	return &TreeTable[T]{
		TableContext: TableContext[T]{},
		Root:         newRoot(data),
		header: tableHeader[T]{
			sortOrder:          0,
			sortedBy:           0,
			drags:              make([]tableDrag, 0),
			columnCells:        rowCells,
			clickedColumnIndex: -1,
			manualWidthSet:     make([]bool, columnCount),
		},
		rootRows:                     nil,
		filteredRows:                 nil,
		SelectedNode:                 nil,
		contextMenu:                  nil,
		columnCount:                  columnCount,
		maxDepth:                     0,
		maxColumnCellWidths:          nil,
		maxColumnTextWidths:          nil,
		maxHierarchyColumnWidthCache: 0,
		columns:                      nil,
		DragRemovedRowsCallback:      nil,
		DropOccurredCallback:         nil,
		inLayoutHeader:               false,
		columnResizeStart:            0,
		columnResizeBase:             0,
		columnResizeOverhead:         0,
		preventUserColumnResize:      false,
		awaitingSyncToModel:          false,
		wasDragged:                   false,
		dividerDrag:                  false,
		LongPressCallback:            nil,
		pressStarted:                 time.Time{},
		longPressed:                  false,
		List: widget.List{
			Scrollbar: widget.Scrollbar{},
			List: layout.List{
				Axis:        layout.Vertical,
				ScrollToEnd: false,
				Alignment:   0,
				Position:    layout.Position{},
			},
		},
		once: sync.Once{},
	}
}

func (t *TreeTable[T]) Layout(gtx layout.Context) layout.Dimensions {
	t.once.Do(func() {
		if t.SetRootRowsCallBack != nil { // mitmproxy
			t.SetRootRowsCallBack()
		}
		if t.JsonName == "" {
			mylog.Check("JsonName is empty")
		}
		mylog.CheckNil(t.UnmarshalRowCells)
		// mylog.CheckNil(ctx.SetRootRowsCallBack)//mitmproxy
		mylog.CheckNil(t.RowSelectedCallback)

		//if t.FileDropCallback == nil {
		//	t.FileDropCallback = func(files []string) {
		//		if filepath.Ext(files[0]) == ".json" {
		//			mylog.Info("dropped file", files[0])
		//			table.ResetChildren()
		//			b := stream.NewBuffer(files[0])
		//			mylog.Check(json.Unmarshal(b.Bytes(), table)) // todo test need a zero table?
		//			fnUpdate() t.SizeColumnsToFit(gtx, false)
		//		}
		//		mylog.Struct("todo", files)
		//	}
		//}
		t.SizeColumnsToFit(gtx)
	})
	t.rootRows = t.Root.Children // 增删改查不一定去计算列宽导致有些节点已经删除，是空指针导致panic，所有这里需要刷新一下rootRows，不知道这里是否会内存泄漏
	//list := material.List(th, &t.List)
	if t.contextMenu == nil {
		t.contextMenu = NewContextMenu(len(t.rootRows), nil)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.HeaderFrame(gtx) // 渲染表头
			//t.inLayoutHeader = true
			//return t.layoutDrag(gtx, func(gtx layout.Context, row int) layout.Dimensions {
			//	return t.HeaderFrame(gtx) // 渲染表头
			//})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			var rootRows []layout.Widget
			for i, n := range t.rootRows {
				rootRows = append(rootRows, func(gtx layout.Context) layout.Dimensions {
					t.contextMenu.Once.Do(func() {
						item := ContextMenuItem{}
						for _, kind := range CopyRowType.EnumTypes() {
							switch kind {
							case CopyRowType:
								item = ContextMenuItem{
									Title:     "",
									Icon:      SvgIconCopy,
									Can:       func() bool { return true },
									Do:        func() { t.SelectedNode.CopyRow(gtx, t.maxColumnTextWidths) },
									Clickable: widget.Clickable{},
								}
							case ConvertToContainerType:
								item = ContextMenuItem{
									Title: "",
									Icon:  SvgIconConvertToContainer,
									Can:   func() bool { return !n.Container() }, // n是当前渲染的行
									Do: func() {
										t.SelectedNode.SetType("ConvertToContainer" + ContainerKeyPostfix) //? todo bug：这里是失败的，导致再次点击这里转换的节点后ConvertToNonContainer没有弹出来
										t.SelectedNode.ID = newID()
										t.SelectedNode.SetOpen(true)
										t.SelectedNode.Children = make([]*Node[T], 0)
										t.updateMaxColumnCellWidth(gtx, t.SelectedNode)
									},
									Clickable: widget.Clickable{},
								}
							case ConvertToNonContainerType:
								item = ContextMenuItem{
									Title: "",
									Icon:  SvgIconConvertToNonContainer,
									Can:   func() bool { return n.Container() }, // n是当前渲染的行
									Do: func() {
										t.SelectedNode.SetType("")
										t.SelectedNode.ID = newID()
										for _, child := range t.SelectedNode.Children {
											child.parent = t.SelectedNode.parent
											child.ID = newID()
										}
										t.SelectedNode.ResetChildren()
										t.updateMaxColumnCellWidth(gtx, t.SelectedNode)
									},
									AppendDivider: true,
									Clickable:     widget.Clickable{},
								}
							case NewType:
								item = ContextMenuItem{
									Title: "",
									Icon:  SvgIconCircledAdd,
									Can:   func() bool { return true },
									Do: func() {
										var zero T
										t.InsertAfter(gtx, NewNode(zero))
									},
									Clickable: widget.Clickable{},
								}
							case NewContainerType:
								item = ContextMenuItem{
									Title: "",
									Icon:  SvgIconCircledVerticalEllipsis,
									Can:   func() bool { return true },
									Do: func() {
										var zero T // todo edit type?
										t.InsertAfter(gtx, NewContainerNode("NewContainerNode", zero))
									},
									Clickable: widget.Clickable{},
								}
							case DeleteType:
								item = ContextMenuItem{
									Title:     "",
									Icon:      SvgIconTrash,
									Can:       func() bool { return true },
									Do:        func() { t.Remove(gtx) },
									Clickable: widget.Clickable{},
								}
							case DuplicateType:
								item = ContextMenuItem{
									Title: "",
									Icon:  SvgIconDuplicate,
									Can:   func() bool { return true },
									Do: func() {
										t.InsertAfter(gtx, t.SelectedNode.Clone())
									},
									Clickable: widget.Clickable{},
								}
							case EditType:
								item = ContextMenuItem{
									Title:         "",
									Icon:          SvgIconEdit,
									Can:           func() bool { return true },
									Do:            func() { t.Edit(gtx) },
									AppendDivider: true,
									Clickable:     widget.Clickable{},
								}
							case OpenAllType:
								item = ContextMenuItem{
									Title:     "",
									Icon:      SvgIconHierarchy,
									Can:       func() bool { return true },
									Do:        func() { t.Root.OpenAll() },
									Clickable: widget.Clickable{},
								}
							case CloseAllType:
								item = ContextMenuItem{
									Title:     "",
									Icon:      SvgIconCircledVerticalEllipsis,
									Can:       func() bool { return true },
									Do:        func() { t.Root.CloseAll() },
									Clickable: widget.Clickable{},
								}
							case SaveDataType:
								item = ContextMenuItem{
									Title:     "",
									Icon:      SvgIconSaveContent,
									Can:       func() bool { return true },
									Do:        func() { t.SaveDate() },
									Clickable: widget.Clickable{},
								}
							}
							item.Title = kind.String()
							if item.Can() {
								t.contextMenu.AddItem(item)
							}
						}
						if items := t.ContextMenuItems(gtx, n); items != nil { // gtx用于调整列宽，触发增删改之后，n是正在渲染的节点，用于取出元数据控制菜单是否渲染
							for _, item := range items {
								if item.Can() {
									t.contextMenu.AddItem(item)
								}
							}
						}
					})
					return t.RowFrame2(gtx, n, i)
				})
			}
			t.contextMenu.DrawRow = func(gtx C, i int) D {
				return rootRows[i](gtx)
			}
			return t.contextMenu.Layout(gtx)

			//return list.Layout(gtx, len(t.rootRows), func(gtx layout.Context, index int) layout.Dimensions {
			//////////////////////////
			//t.inLayoutHeader = false
			//return t.layoutDrag(gtx, func(gtx layout.Context, row int) layout.Dimensions {
			//	return t.RowFrame(gtx, t.rootRows[index], index)
			//})
			//})
		}),
	)
}
func (t *TreeTable[T]) RowFrame2(gtx layout.Context, n *Node[T], rowIndex int) layout.Dimensions {
	n.rowCells = t.MarshalRowCells(n)
	for i := range n.rowCells {
		n.rowCells[i].rowID = rowIndex                          // 斑马线
		n.rowCells[i].columID = t.header.columnCells[i].columID // 列分隔符,更新层级列宽度
	}

	bgColor := t.processEvent(gtx, n, rowIndex) // 处理鼠标事件
	rowClick := &n.rowClick

	var rowCells []layout.FlexChild

	layoutHierarchyColumn := func(gtx layout.Context, cell CellData) layout.Dimensions {
		leftIndent := t.calcCurrentHierarchyColumnLeftIndent(gtx, n, cell)
		t.calcCurrentHierarchyColumnWidthAndSafeCheck(gtx, cell, leftIndent)

		maxColumnCellWidth := t.maxColumnCellWidths[hierarchyColumnID]
		gtx.Constraints.Min.X = int(maxColumnCellWidth)
		gtx.Constraints.Max.X = int(maxColumnCellWidth)

		return layout.Flex{
			Axis: layout.Horizontal,
			// Spacing:   0,
			// Alignment: layout.Start,
			// WeightSum: 0,
		}.Layout(gtx, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return rowClick.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// 绘制层级图标-----------------------------------------------------------------------------------------------------------------
				HierarchyInsert := layout.Inset{Left: leftIndent, Top: -1} // 层级图标居中,行高调整后这里需要下移使得图标居中
				if !n.Container() {
					return HierarchyInsert.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Dimensions{}
					})
				}
				return HierarchyInsert.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					svg := SvgIconCircledChevronRight
					if n.isOpen {
						svg = SvgIconCircledChevronDown
					}
					return iconButtonSmall(new(widget.Clickable), svg, "").Layout(gtx)
					// return NewButton("", nil).SetRectIcon(true).SetIcon(svg).Layout(gtx)
				})
			})
		}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return rowClick.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return t.CellFrame(gtx, cell, t.maxColumnCellWidths[hierarchyColumnID], false, layout.Inset{
						Top:    topPadding, // 层级列文本不知道什么原因往上飘了，top: 0的样子，所以往下挪一下让它居中
						Bottom: 0,
						Left:   leftPadding, // 绘制层级列文本,和层级图标聚拢在一起--------------------------------------
						Right:  0,
					})
				})
				//DrawColumnDivider(gtx, cell.columID) // 这里绘制的列分割线才没有虚线,gtx被破坏了? 永远不要移动这个位置
				//return layout.Stack{Alignment: layout.Center}.Layout(gtx,//似乎会影响层级列文本布局
				//	layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				//		if len(cell.Text) > 80 {
				//			cell.Text = cell.Text[:len(cell.Text)/2] + "..."
				//			// todo 更好的办法是让富文本编辑器做这个事情，对 maxline 。。。 看看代码编辑器扩建是如何实现这个的
				//			// 然后双击编辑行的时候从富文本取出完整行并换行显示，structView需要好好设计一下这个
				//			// 这个在抓包场景很那个，url列一般都长
				//		}
				//		return t.CellFrame(gtx, cell, t.maxColumnCellWidths[hierarchyColumnID], false, layout.Inset{
				//			Top:    topPadding, // 层级列文本不知道什么原因往上飘了，top: 0的样子，所以往下挪一下让它居中
				//			Bottom: 0,
				//			Left:   leftPadding, // 绘制层级列文本,和层级图标聚拢在一起--------------------------------------
				//			Right:  0,
				//		})
				//	}),
				//	//t.drawContextMenu(n, hierarchyColumnID), //todo bug
				//)
			}),
		)
	}

	rowCells = append(rowCells, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		hierarchyColumnCell := n.rowCells[hierarchyColumnID] // 层级列单元格
		return layoutHierarchyColumn(gtx, hierarchyColumnCell)
	}))

	// 绘制非层级列-----------------------------------------------------------------------------------------------------------------
	for i, cell := range n.rowCells {
		if i == hierarchyColumnID {
			continue
		}
		rowCells = append(rowCells, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return rowClick.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Clickable(gtx, &n.rowCells[i].Clickable, func(gtx layout.Context) layout.Dimensions {
					DrawColumnDivider(gtx, cell.columID) // 这里绘制的列分割线才没有虚线,gtx被破坏了? 永远不要移动这个位置
					if len(cell.Text) > 80 {
						cell.Text = cell.Text[:len(cell.Text)/2] + "..."
						// todo 更好的办法是让富文本编辑器做这个事情，对 maxline 。。。 看看代码编辑器扩建是如何实现这个的
						// 然后双击编辑行的时候从富文本取出完整行并换行显示，structView需要好好设计一下这个
						// 这个在抓包场景很那个，url列一般都长
					}
					return t.CellFrame(gtx, cell, t.maxColumnCellWidths[i], false, layout.Inset{
						Top:    0,
						Bottom: 0,
						Left:   leftPadding,
						// Right:  rightPadding,
					})
				})
			})
		}))
	}

	rows := []layout.FlexChild{ // 合成层级列和其他列的单元格为一行,并设置该行的背景和行高
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return Background{bgColor}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.Y = gtx.Dp(defaultRowHeight) // 主题的字体大小也会影响行高，这里设置最小行高为14dp
				gtx.Constraints.Max.Y = gtx.Dp(defaultRowHeight) // 限制行高以避免列分割线呈现虚线视觉
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, rowCells...)
			})

		}),
	}
	if n.CanHaveChildren() && n.isOpen { // 如果是容器节点则递归填充孩子节点形成多行
		for _, child := range n.Children {
			rows = append(rows, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				rowIndex++
				return t.RowFrame2(gtx, child, rowIndex)
			}))
		}
	}
	// 把全部行垂直居中排列，rowClick点击后根据点击状态显示了这里填充了多少行，展开节点后看到的行就是这里来的
	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx, rows...)
}

func (t *TreeTable[T]) RowFrame(gtx layout.Context, n *Node[T], rowIndex int) layout.Dimensions {
	n.rowCells = t.MarshalRowCells(n)
	for i := range n.rowCells {
		n.rowCells[i].rowID = rowIndex                          // 斑马线
		n.rowCells[i].columID = t.header.columnCells[i].columID // 列分隔符,更新层级列宽度
	}

	// bgColor := t.processEvent(gtx, n, rowIndex) // 处理鼠标事件
	bgColor := RowColor(rowIndex)
	rowClick := &n.rowClick

	var rowCells []*Resizable

	layoutHierarchyColumn := func(gtx layout.Context, cell CellData) layout.Dimensions {
		leftIndent := t.calcCurrentHierarchyColumnLeftIndent(gtx, n, cell)
		t.calcCurrentHierarchyColumnWidthAndSafeCheck(gtx, cell, leftIndent)

		maxColumnCellWidth := t.maxColumnCellWidths[hierarchyColumnID]
		gtx.Constraints.Min.X = int(maxColumnCellWidth)
		gtx.Constraints.Max.X = int(maxColumnCellWidth)

		return layout.Flex{
			Axis: layout.Horizontal,
			// Spacing:   0,
			// Alignment: layout.Start,
			// WeightSum: 0,
		}.Layout(gtx, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return rowClick.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// 绘制层级图标-----------------------------------------------------------------------------------------------------------------
				HierarchyInsert := layout.Inset{Left: leftIndent, Top: -1} // 层级图标居中,行高调整后这里需要下移使得图标居中
				if !n.Container() {
					return HierarchyInsert.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Dimensions{}
					})
				}
				return HierarchyInsert.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					svg := SvgIconCircledChevronRight
					if n.isOpen {
						svg = SvgIconCircledChevronDown
					}

					return iconButtonSmall(new(widget.Clickable), svg, "").Layout(gtx)
					// return NewButton("", nil).SetRectIcon(true).SetIcon(svg).Layout(gtx)
				})
			})
		}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return rowClick.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return t.CellFrame(gtx, cell, t.maxColumnCellWidths[hierarchyColumnID], false, layout.Inset{
						Top:    topPadding, // 层级列文本不知道什么原因往上飘了，top: 0的样子，所以往下挪一下让它居中
						Bottom: 0,
						Left:   leftPadding, // 绘制层级列文本,和层级图标聚拢在一起--------------------------------------
						Right:  0,
					})
				})
				//DrawColumnDivider(gtx, cell.columID) // 这里绘制的列分割线才没有虚线,gtx被破坏了? 永远不要移动这个位置
				//return layout.Stack{Alignment: layout.Center}.Layout(gtx,//似乎会影响层级列文本布局
				//	layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				//		if len(cell.Text) > 80 {
				//			cell.Text = cell.Text[:len(cell.Text)/2] + "..."
				//			// todo 更好的办法是让富文本编辑器做这个事情，对 maxline 。。。 看看代码编辑器扩建是如何实现这个的
				//			// 然后双击编辑行的时候从富文本取出完整行并换行显示，structView需要好好设计一下这个
				//			// 这个在抓包场景很那个，url列一般都长
				//		}
				//		return t.CellFrame(gtx, cell, t.maxColumnCellWidths[hierarchyColumnID], false, layout.Inset{
				//			Top:    topPadding, // 层级列文本不知道什么原因往上飘了，top: 0的样子，所以往下挪一下让它居中
				//			Bottom: 0,
				//			Left:   leftPadding, // 绘制层级列文本,和层级图标聚拢在一起--------------------------------------
				//			Right:  0,
				//		})
				//	}),
				//	//t.drawContextMenu(n, hierarchyColumnID), //todo bug
				//)
			}),
		)
	}

	rowCells = append(rowCells, &Resizable{
		ratio: 0,
		Widget: func(gtx layout.Context) layout.Dimensions {
			hierarchyColumnCell := n.rowCells[hierarchyColumnID] // 层级列单元格
			return layoutHierarchyColumn(gtx, hierarchyColumnCell)
		},
		DividerHandler: nil,
		float:          float{},
		resize:         nil,
		prev:           nil,
		next:           nil,
	})

	// 绘制非层级列-----------------------------------------------------------------------------------------------------------------
	for i, cell := range n.rowCells {
		if i == hierarchyColumnID {
			continue
		}
		rowCells = append(rowCells, &Resizable{
			ratio: 0,
			Widget: func(gtx layout.Context) layout.Dimensions {
				return rowClick.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return material.Clickable(gtx, &n.rowCells[i].Clickable, func(gtx layout.Context) layout.Dimensions {
						DrawColumnDivider(gtx, cell.columID) // 这里绘制的列分割线才没有虚线,gtx被破坏了? 永远不要移动这个位置
						return layout.Stack{Alignment: layout.Center}.Layout(gtx,
							layout.Stacked(func(gtx layout.Context) layout.Dimensions {
								if len(cell.Text) > 80 {
									cell.Text = cell.Text[:len(cell.Text)/2] + "..."
									// todo 更好的办法是让富文本编辑器做这个事情，对 maxline 。。。 看看代码编辑器扩建是如何实现这个的
									// 然后双击编辑行的时候从富文本取出完整行并换行显示，structView需要好好设计一下这个
									// 这个在抓包场景很那个，url列一般都长
								}
								return t.CellFrame(gtx, cell, t.maxColumnCellWidths[i], false, layout.Inset{
									Top:    0,
									Bottom: 0,
									Left:   leftPadding,
									// Right:  rightPadding,
								})
							}),
							// t.drawContextMenu(n, i),//todo
						)
					})
				})
			},
			DividerHandler: nil,
			float:          float{},
			resize:         nil,
			prev:           nil,
			next:           nil,
		})
	}

	rows := []layout.FlexChild{ // 合成层级列和其他列的单元格为一行,并设置该行的背景和行高
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return Background{bgColor}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.Y = gtx.Dp(defaultRowHeight) // 主题的字体大小也会影响行高，这里设置最小行高为14dp
				gtx.Constraints.Max.Y = gtx.Dp(defaultRowHeight) // 限制行高以避免列分割线呈现虚线视觉
				return NewResizeWidget(layout.Horizontal, func(index int, newWidth int) {
					fmt.Printf("列 %d 新宽度: %dpx\n", index, newWidth)
					// 这里可以更新表格列宽或执行其他操作
					t.maxColumnCellWidths[index] = unit.Dp(newWidth)
				}, rowCells...).Layout(gtx)
			})
		}),
	}
	if n.CanHaveChildren() && n.isOpen { // 如果是容器节点则递归填充孩子节点形成多行
		for _, child := range n.Children {
			rows = append(rows, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				rowIndex++
				return t.RowFrame(gtx, child, rowIndex)
			}))
		}
	}

	// 把全部行垂直居中排列，rowClick点击后根据点击状态显示了这里填充了多少行，展开节点后看到的行就是这里来的
	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx, rows...)
}

func (t *TreeTable[T]) processEvent(gtx layout.Context, n *Node[T], rowIndex int) color.NRGBA {
	rowClick := &n.rowClick
	evt, ok := gtx.Source.Event(pointer.Filter{
		Target: rowClick,
		Kinds:  pointer.Press, //| pointer.Release | pointer.Drag,
	})
	if ok {
		e, ok := evt.(pointer.Event)
		if ok {
			switch {
			case e.Kind == pointer.Press: // 左键，右键，双击
				t.SelectedNode = n
				// bgColor = Orange300
			case e.Source == pointer.Touch: // todo检查是否长按并测试apk
				t.SelectedNode = n
			}
		}
	}

	click, ok := rowClick.Update(gtx)
	if ok {
		switch click.NumClicks {
		case 1:
			n.isOpen = !n.isOpen // 切换展开状态
			//if n.cellClickedCallback != nil {
			//	n.cellClickedCallback(n) // 单元格点击回调
			//}
			if t.RowSelectedCallback != nil {
				t.RowSelectedCallback() // 行选中回调
			}

		case 2:
			t.Edit(gtx)
			//if t.RowDoubleClickCallback != nil { // 行双击回调
			//	go t.RowDoubleClickCallback(n)
			//	//gtx.Execute(op.InvalidateCmd{})
			//}
		}
	}

	bgColor := RowColor(rowIndex)
	switch {
	case t.SelectedNode == n: // 设置选中背景色,这个需要在第一的位置,在选中另外的节点时间段内这个条件成立，能保持背景色不变
		bgColor = color.NRGBA{R: 200, G: 145, B: 50, A: 255}
	case rowClick.Hovered(): // 设置悬停背景色
		// https://rgbacolorpicker.com/
		bgColor = th.Color.TreeHoveredBgColor
	default:
		//todo 如果children的最后一个节点是黑色，len Chidren是奇数，那么root的node父级的父级的背景色需要设置为白色,bug
		//if n.LenChildren()%2 == 1 && bgColor == rowBlackColor {
		//	bgColor = rowWhiteColor
		//}
	}
	return bgColor
}

func (t *TreeTable[T]) IsRowSelected() bool { return t.SelectedNode != nil }

func (t *TreeTable[T]) CellFrame(gtx layout.Context, data CellData, width unit.Dp, isHeader bool, inset layout.Inset) layout.Dimensions {
	// 固定单元格宽度为计算好的每列最大label宽度实现表头和body对齐,因为表头和body都调用这个函数渲染单元格，只有限制min和max才能每列保证表头单元格和body单元格具有相等的宽度，从而实现表头和body对齐
	gtx.Constraints.Min.X = int(width)
	gtx.Constraints.Max.X = int(width)
	DrawColumnDivider(gtx, data.columID) // 为每列绘制列分隔条

	if data.FgColor == (color.NRGBA{}) {
		data.FgColor = White
	}
	//richText := NewRichText()
	//richText.AddSpan(richtext.SpanStyle{
	//	// Font:        font.Font{},
	//	Size:        unit.Sp(12),
	//	Color:       data.FgColor,
	//	Content:     data.Text,
	//	Interactive: false,
	//})
	if isHeader { // 加高表头高度
		inset.Top = 8
		inset.Bottom = 8
	}
	return layout.Flex{
		Axis:      layout.Horizontal,
		Spacing:   0,
		Alignment: 0,
		WeightSum: 0,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if data.Icon == nil { // 格式支持：*giosvg.Icon, *widget.Icon, *widget.Image, image.Image
				return layout.Dimensions{}
			}
			size := image.Pt(gtx.Dp(defaultHierarchyColumnIconSize), gtx.Dp(defaultHierarchyColumnIconSize))
			gtx.Constraints.Min = size
			gtx.Constraints.Max = size
			return layout.Inset{
				Top:    topPadding + 1,
				Bottom: 0,
				Left:   leftPadding,
				Right:  0,
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return icon.Layout(gtx, data.Icon, color.NRGBA{}, defaultHierarchyColumnIconSize)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return inset.Layout(gtx, material.Body2(th, data.Text).Layout)
		}),
	)
	// return inset.Layout(gtx, richText.Layout)
}

func InitHeader(data any) (rowCells []CellData) {
	fields := stream.ReflectVisibleFields(data)
	rowCells = make([]CellData, 0)
	for i, field := range fields {
		if field.Tag.Get("table") != "" { // 中文表头简短
			field.Name = field.Tag.Get("table")
		}
		rowCells = append(rowCells, CellData{
			Text:      field.Name,
			Tooltip:   "",
			FgColor:   color.NRGBA{},
			IsNasm:    false,
			Disabled:  false,
			columID:   i,
			rowID:     0,
			Clickable: widget.Clickable{},
			RichText:  RichText{},
		})
	}
	return
}

func (t *TreeTable[T]) SizeColumnsToFit(gtx layout.Context) { // 增删改查中，只有增改+实例化，共3次需要调用这个函数，增改需要更新缓存的每列最大列宽即可，所以这个函数理论上只需要执行一次，这样性能最好
	originalConstraints := gtx.Constraints         // 保存原始约束
	rows := make([][]CellData, 0, len(t.rootRows)) // 用于存储所有行,如果不这么做的话，节点增删改查就不会实时刷新
	for _, node := range t.Root.Walk() {
		rows = append(rows, t.MarshalRowCells(node))
	}
	rows = slices.Insert(rows, 0, t.header.columnCells) // 插入表头行
	t.columns = TransposeMatrix(rows)                   // 如果不这么做的话，节点增删改查就不会实时刷新,为了提高性能需要手动刷新节点和宽度
	t.maxColumnCellWidths = make([]unit.Dp, t.columnCount)
	t.maxColumnTextWidths = make([]unit.Dp, t.columnCount)
	for i, column := range t.columns {
		for _, data := range column {
			t.maxColumnTextWidths[i] = max(t.maxColumnTextWidths[i], align.StringWidth[unit.Dp](data.Text))
			if gtx != (layout.Context{}) {
				t.maxColumnCellWidths[i] = max(t.maxColumnCellWidths[i], LabelWidth(gtx, data.Text))
			}
		}
	}
	// todo 增删改查max一下更新层级列宽度，其他列的话也max一下，小于对应列的最大宽度就不用更新了，这样性能会好一点
	t.updateMaxHierarchyColumnCellWidth()
	gtx.Constraints = originalConstraints
}

func (t *TreeTable[T]) updateMaxHierarchyColumnCellWidth() { // 计算层级列最大列单元格宽度,表头成列和body层级列的最大宽度公用这个函数，且计算也统一执行SizeColumnsToFit函数,这样表头和body都调用同一个层级列最大宽度，稳稳的对齐了
	t.maxDepth = t.Root.MaxDepth() // todo新增容器节点需要刷新
	if t.maxHierarchyColumnWidthCache == 0 {
		t.maxHierarchyColumnWidthCache = t.maxColumnCellWidths[hierarchyColumnID]
	}

	// leftPadding //有左填充,rootRows不应该紧靠左边 todo 如果层级列宽度检查失败的话考虑处理这个，加上去就可以了,另外是层级列的图标应按安排两个:层级图标和层级列文本图标，最后是表头单元格的排序svg,如果以后要渲染这个的话
	maxWidth := t.maxDepth*HierarchyIndent + // 最大深度的左缩进
		defaultHierarchyIconSize + // 层级图标宽度
		defaultHierarchyColumnIconSize + // 层级列图标宽度
		leftPadding + // 左侧padding
		t.maxHierarchyColumnWidthCache + // 最长的单元格文本宽度
		// rightPadding + // 右侧padding
		DividerWidth // 列分隔条宽度

	t.maxColumnCellWidths[hierarchyColumnID] = max(maxWidth, t.maxHierarchyColumnWidthCache)
}

const (
	hierarchyColumnID = 0
	leftPadding       = unit.Dp(4) // 单元格左侧填充，和列分割线的间距
	// rightPadding          = unit.Dp(4)   // 单元格右侧填充，和列分割线的间距
	topPadding                     = unit.Dp(4)   // 单元格上侧填充，似单元格间文本居中
	bottomPadding                  = unit.Dp(4)   // 单元格下侧填充，似单元格间文本居中
	minColumnWidth                 = unit.Dp(100) // 最小列宽
	defaultHierarchyIconSize       = unit.Dp(12)  // 层级图标默认大小
	defaultHierarchyColumnIconSize = unit.Dp(16)  // 层级列图标默认大小
	HierarchyIndent                = unit.Dp(16)  // 层级缩进
	defaultRowHeight               = unit.Dp(22)  // 行高
	defaultHeaderHeight            = unit.Dp(24)  // 表头行高
	defaultHeaderFontSize          = unit.Sp(12)  // 表头字体大小
	defaultRowFontSize             = unit.Sp(12)  // 行字体大小
)

func (t *TreeTable[T]) SaveDate() { // todo 支持apk数据目录 app.dataDir()
	t.JsonName = strings.TrimSuffix(t.JsonName, ".json")
	stream.MarshalJsonToFile(t.Root, filepath.Join("cache", t.JsonName+".json"))
	stream.WriteTruncate(filepath.Join("cache", t.JsonName+".txt"), t.Document()) // 调用t.Format()
	if t.IsDocument {
		b := stream.NewBuffer("")
		b.WriteStringLn("# " + t.JsonName + " document table")
		b.WriteStringLn("```text")
		b.WriteStringLn(t.Document())
		b.WriteStringLn("```")
		stream.WriteTruncate("README2.md", b.String())
	}
}

// TransposeMatrix 把行切片矩阵置换为列切片,用于计算最大列宽的参数
func TransposeMatrix[T any](rows [][]T) (columns [][]T) {
	if len(rows) == 0 {
		return [][]T{}
	}
	numColumns := len(rows[0])
	columns = make([][]T, numColumns)
	for i := range columns {
		columns[i] = make([]T, len(rows))
	}
	for i, row := range rows {
		for j := range row {
			columns[j][i] = row[j]
		}
	}
	return
}

func (t *TreeTable[T]) HeaderFrame(gtx layout.Context) layout.Dimensions {
	// 实现右键复制列到剪切板功能
	if t.header.contextMenu == nil {
		t.header.contextMenu = NewContextMenu(1, nil)
		t.header.contextMenu.AddItem(ContextMenuItem{
			Title:         "CopyColumn",
			Icon:          SvgIconCopy,
			Can:           func() bool { return true },
			Do:            func() { t.CopyColumn(gtx) },
			AppendDivider: true,
			Clickable:     widget.Clickable{},
		})
	}
	t.header.contextMenu.DrawRow = func(gtx layout.Context, index int) layout.Dimensions {
		var cols []layout.FlexChild
		elems := make([]*Resizable, 0)
		for i, cell := range t.header.columnCells {
			if cell.Disabled {
				continue
			}
			clickable := &t.header.columnCells[i].Clickable
			if clickable.Clicked(gtx) {
				t.header.clickedColumnIndex = i
				if t.header.sortedBy == t.header.clickedColumnIndex {
					t.header.sortAscending = !t.header.sortAscending // 切换升序/降序

					switch t.header.sortOrder {
					case sortNone:
						t.header.sortOrder = sortAscending
					case sortAscending:
						t.header.sortOrder = sortDescending
					case sortDescending:
						t.header.sortOrder = sortAscending
					}

				} else {
					t.header.sortedBy = t.header.clickedColumnIndex // 更新排序列
					t.header.sortAscending = true                   // 设为升序
					t.header.sortOrder = sortAscending
				}

				mylog.Info("clickedColumnIndex", t.header.clickedColumnIndex)
				// mylog.Info(t.header.columnCells[i].Text, LabelWidth(gtx, t.header.columnCells[i].Text))
				if t.header.clickedColumnIndex > -1 && t.header.sortedBy > -1 {
					for i := range t.header.columnCells {
						t.header.columnCells[i].Text = strings.TrimSuffix(t.header.columnCells[i].Text, " ⇩")
						t.header.columnCells[i].Text = strings.TrimSuffix(t.header.columnCells[i].Text, " ⇧")
					}
					switch t.header.sortOrder {
					case sortNone:
					case sortAscending:
						t.header.columnCells[i].Text += " ⇧"
					case sortDescending:
						t.header.columnCells[i].Text += " ⇩"
					}
					t.Sort()
					// t.header.clickedColumnIndex = -1 //重置点击列索引
				}
			}
			// 拦截右击事件并在事件中赋值命中的列id
			evt, ok := gtx.Source.Event(pointer.Filter{
				Target: clickable,
				Kinds:  pointer.Press | pointer.Release,
			})
			if ok {
				e, ok := evt.(pointer.Event)
				if ok {
					if e.Kind == pointer.Press || e.Source == pointer.Touch {
						t.header.clickedColumnIndex = i
					}

					switch e.Buttons {
					case pointer.ButtonPrimary:
					case pointer.ButtonSecondary:
						if e.Kind == pointer.Press {
							t.header.clickedColumnIndex = i
						}
					}
				}
			}
			cols = append(cols, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return Background{Color: ColorHeaderFg}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return material.Clickable(gtx, clickable, func(gtx layout.Context) layout.Dimensions {
						cellFrame := t.CellFrame(gtx, t.header.columnCells[i], t.maxColumnCellWidths[i], true, layout.Inset{
							Top:    8,
							Bottom: 8,
							Left:   leftPadding,
							// Right:  rightPadding,
						})
						elems = append(elems, &Resizable{Widget: func(gtx layout.Context) layout.Dimensions {
							return cellFrame
						}})
						return cellFrame
					})
				})
			}))
		}
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, WeightSum: 0}.Layout(gtx, cols...)
	}
	return t.header.contextMenu.Layout(gtx)

	//return layout.Stack{}.Layout(gtx,
	//	layout.Stacked(func(gtx layout.Context) layout.Dimensions {
	//		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, WeightSum: 0}.Layout(gtx, cols...)
	//	}),
	//	layout.Expanded(func(gtx layout.Context) layout.Dimensions {
	//		resizeWidget = NewResizeWidget(layout.Horizontal, elems...)
	//		return resizeWidget.Layout(gtx)
	//	}),
	//)
}

var (
	rowWhiteColor = color.NRGBA{R: 49, G: 49, B: 49, A: 255} // 白色
	rowBlackColor = color.NRGBA{R: 43, G: 43, B: 43, A: 255} // 黑色
)

func RowColor(rowIndex int) color.NRGBA { // 奇偶行背景色
	if rowIndex%2 != 0 {
		return rowWhiteColor
	}
	return rowBlackColor
}

const DividerWidth = unit.Dp(1)

// DrawColumnDivider 分隔线绘制函数
func DrawColumnDivider(gtx layout.Context, col int) {
	if col > 0 { // 层级列不要绘制分隔线
		tallestHeight := gtx.Dp(unit.Dp(gtx.Constraints.Max.Y))
		stack3 := clip.Rect{Max: image.Pt(int(DividerWidth), tallestHeight)}.Push(gtx.Ops)
		paint.Fill(gtx.Ops, DividerFg)
		stack3.Pop()
	}
}

func (t *TreeTable[T]) CopyColumn(gtx layout.Context) string {
	if t.header.clickedColumnIndex < 0 {
		gtx.Execute(clipboard.WriteCmd{Data: io.NopCloser(strings.NewReader("t.header.clickedColumnIndex < 0 "))})
		return "t.header.clickedColumnIndex < 0 "
	}
	g := stream.NewGeneratedFile()
	g.P("var columnData = []string{")
	g.P(strconv.Quote(t.header.columnCells[t.header.clickedColumnIndex].Text), ",")
	cellData := t.columns[t.header.clickedColumnIndex]
	for _, datum := range cellData {
		g.P(strconv.Quote(datum.Text), ",")
	}
	g.P("}")
	gtx.Execute(clipboard.WriteCmd{Data: io.NopCloser(strings.NewReader(g.Format()))})
	return g.String()
}

type TableDragData[T any] struct {
	Table *Node[T]
	Rows  []*Node[T]
}

var DefaultTableTheme = TableTheme{
	HierarchyIndent:   16,
	MinimumRowHeight:  16,
	ColumnResizeSlop:  4,
	ShowRowDivider:    true,
	ShowColumnDivider: true,
}

type TableTheme struct {
	Padding           layout.Inset
	HierarchyColumnID int
	HierarchyIndent   unit.Dp
	MinimumRowHeight  unit.Dp
	ColumnResizeSlop  unit.Dp
	ShowRowDivider    bool
	ShowColumnDivider bool
}

func (t *TreeTable[T]) layoutDrag(gtx layout.Context, w rowFn) layout.Dimensions {
	columns := t.maxColumnCellWidths
	var (
		cols          = t.columnCount         // 获取列的数量
		dividers      = cols                  // 列的分隔数
		tallestHeight = gtx.Constraints.Min.Y // 初始化最高的行高

		dividerWidth                   = gtx.Dp(defaultDividerWidth)                   // 获取分割线的宽度
		dividerMargin                  = gtx.Dp(defaultDividerMargin)                  // 获取分割线的边距
		dividerHandleMinVerticalMargin = gtx.Dp(defaultDividerHandleMinVerticalMargin) // 获取分隔处理器的最小垂直边距
		dividerHandleMaxHeight         = gtx.Dp(defaultDividerHandleMaxHeight)         // 获取分隔处理器的最大高度
		dividerHandleWidth             = gtx.Dp(defaultDividerHandleWidth)             // 获取分隔处理器的宽度
		dividerHandleRadius            = gtx.Dp(defaultDividerHandleRadius)            // 获取分隔处理器的圆角半径

		minWidth = unit.Dp(dividerWidth + dividerMargin + dividerHandleWidth) // 计算分隔符的最小宽度
	)
	if cols == 0 { // 如果列数为0，直接返回最小约束大小
		return layout.Dimensions{Size: gtx.Constraints.Min}
	}

	if len(t.header.drags) < dividers { // 如果拖动数组没有足够的长度
		t.header.drags = make([]tableDrag, dividers) // 初始化拖动对象数组
	}

	// OPT(dh): we don't need to do this for each t, only once per table
	for i := range t.header.drags { // 遍历每个拖动对象
		drag := &t.header.drags[i]    // 获取当前拖动对象
		width := columns[i]           // 获取当前列宽度
		drag.hover.Update(gtx.Source) // 更新当前列的悬停状态
		// OPT(dh): Events allocates
		var delta unit.Dp // 初始化偏移量

		for { // 循环处理拖动事件
			ev, ok := drag.drag.Update(gtx.Metric, gtx.Source, gesture.Horizontal) // 更新拖动状态
			if !ok {                                                               // 如果事件不合法，退出循环
				break
			}
			switch ev.Kind { // 根据事件类型进行处理
			case pointer.Press: // 按下事件
				drag.startPos = ev.Position.X                             // 记录拖动的起始位置
				drag.shrinkNeighbor = !ev.Modifiers.Contain(key.ModShift) // 判断是否收缩相邻列
			case pointer.Drag: // 拖动事件
				t.header.manualWidthSet[i] = true // 在这里标记列宽已手动调整
				// There may be multiple drag events in a single frame. We mustn't apply all of them or we'll
				// drag too far. Only react to the final event.
				delta = unit.Dp(ev.Position.X - drag.startPos) // 计算当前拖动的偏移量
			}
		}
		if delta != 0 { // 如果存在拖动偏移量
			width += delta                                  // 更新列的宽度
			if drag.shrinkNeighbor && i != len(columns)-1 { // 如果需要收缩相邻列且不是最后一列
				nextCol := columns[i+1] // 获取下一个列
				nextCol -= delta        // 更新下一个列的宽度
				if width < minWidth {   // 如果当前列宽度小于最小宽度
					d := minWidth - width // 计算需要增加的宽度
					width = minWidth      // 将当前列宽度设为最小宽度
					nextCol -= d          // 更新下一个列的宽度
				}
				if nextCol < minWidth { // 如果下一个列宽度小于最小宽度
					d := minWidth - nextCol // 计算需要增加的宽度
					nextCol = minWidth      // 将下一个列宽度设为最小宽度
					width -= d              // 更新当前列宽度
				}
			} else {
				// 如果不需要收缩
				if width < minWidth { // 如果当前列宽度小于最小宽度
					width = minWidth // 将当前列宽度设为最小宽度
				}
			}

			//if width < width.autoWidth { // 如果当前列宽度小于其最小宽度
			//	width = width.autoWidth // 更新列的最小宽度为当前宽度
			//}
			width = max(width, minWidth) // 确保列的宽度不小于最小宽度

			var total unit.Dp             // 初始化总宽度
			for _, col := range columns { // 遍历所有列计算总宽度
				total += col // 累加当前列的宽度
			}
			total += unit.Dp(len(columns) * gtx.Dp(defaultDividerWidth)) // 加上所有分隔符的总宽度
			if total < unit.Dp(gtx.Constraints.Min.X) {                  // 如果总宽度小于最小约束宽度
				columns[len(columns)-1] += unit.Dp(gtx.Constraints.Min.X) - total // 调整最后一列的宽度以适应
			}
		}
	}

	for { // 开始绘制列
		// First draw all columns, leaving gaps for the drag handlers
		var (
			start             = 0             // 初始化当前位置
			origTallestHeight = tallestHeight // 记录最初的高度
		)
		r := op.Record(gtx.Ops)  // 记录当前操作集合
		totalWidth := 0          // 初始化总宽度
		for i := range columns { // 遍历所有列
			colWidth := int(columns[i]) // 获取当前列的宽度
			totalWidth += colWidth      // 更新总宽度
		}
		extra := gtx.Constraints.Min.X - len(columns)*gtx.Dp(defaultDividerWidth) - totalWidth // 计算多余宽度
		colExtra := extra                                                                      // 将多余宽度赋值给列额外宽度

		for i := range columns { // 绘制所有列
			colWidth := int(columns[i]) // 获取当前列宽度
			if colExtra > 0 {           // 如果有多余宽度
				colWidth++ // 当前列宽度加一
				colExtra-- // 多余宽度减一
			}

			gtx := gtx                            // 更新 gtx 上下文
			gtx.Constraints.Min.X = colWidth      // 设置当前列的最小宽度
			gtx.Constraints.Max.X = colWidth      // 设置当前列的最大宽度
			gtx.Constraints.Min.Y = tallestHeight // 设置当前列的最小高度

			stack := op.Offset(image.Pt(start, 0)).Push(gtx.Ops) // 设置当前列绘制的偏移量

			dims := w(gtx, i)                                // 绘制当前列，返回所需尺寸
			dims.Size = gtx.Constraints.Constrain(dims.Size) // 应用约束限制
			tallestHeight = dims.Size.Y                      // 更新当前行的高度
			if i == 0 && tallestHeight > origTallestHeight { // 如果当前行的高度大于初始高度
				origTallestHeight = tallestHeight // 更新初始高度
			}

			start += colWidth + dividerWidth // 更新绘制起始位置
			stack.Pop()                      // 弹出当前堆栈
		}
		call := r.Stop() // 停止记录操作集合

		if tallestHeight > origTallestHeight { // 如果当前高度大于最初高度
			continue // 重新绘制当前行
		}

		call.Add(gtx.Ops) // 将操作添加到上下文中

		// Then draw the drag handlers. The handlers overdraw the columns when hovered.
		var (
			dividerHandleHeight    = min(tallestHeight-2*dividerHandleMinVerticalMargin, dividerHandleMaxHeight) // 获取分隔处理器的高度
			dividerHandleTopMargin = (tallestHeight - dividerHandleHeight) / 2                                   // 计算顶部边距
			dividerStart           = 0                                                                           // 初始化分隔起始位置
			dividerExtra           = extra                                                                       // 设置额外宽度
		)
		for i := range t.header.drags { // 遍历每个拖动对象
			var (
				drag     = &t.header.drags[i] // 获取当前拖动对象
				colWidth = int(columns[i])    // 获取当前列宽度
			)
			dividerStart += colWidth // 更新分隔符的起始位置
			if dividerExtra > 0 {    // 如果还有多余的宽度
				dividerStart++ // 分隔符起始位置加一
				dividerExtra-- // 多余宽度减一
			}

			// We add the drag handler slightly outside the drawn divider, to make it easier to press.
			//
			// We use op.Offset instead of folding dividerStart into the clip.Rect because we want to set the
			// origin of the drag coordinates.
			stack := op.Offset(image.Pt(dividerStart, 0)).Push(gtx.Ops) // set origin for drag coordinates
			stack2 := clip.Rect{
				Min: image.Pt(-dividerMargin-dividerHandleWidth, 0),                         // 设置分隔符的最小边界
				Max: image.Pt(dividerWidth+dividerMargin+dividerHandleWidth, tallestHeight), // 设置分隔符的最大边界
			}.Push(gtx.Ops) // 在上下文中推入分隔符的裁剪矩形

			if t.inLayoutHeader { // 如果当前行是表头 todo test
				drag.hover.Update(gtx.Source)        // 更新悬停状态
				drag.drag.Add(gtx.Ops)               // 添加拖动操作到上下文中
				drag.hover.Add(gtx.Ops)              // 添加悬停操作到上下文中
				pointer.CursorColResize.Add(gtx.Ops) // 设置拖动光标样式为调整列宽

				// Draw the left and right extensions when hovered.
				if drag.hover.Update(gtx.Source) || drag.drag.Dragging() { // 如果悬停或者拖动
					handleShape := clip.UniformRRect(
						image.Rect(
							0,
							dividerHandleTopMargin,
							dividerHandleWidth,
							dividerHandleTopMargin+dividerHandleHeight),
						dividerHandleRadius,
					) // 设定分隔处理器的形状
					handleLeft := handleShape
					handleLeft.Rect = handleShape.Rect.Add(image.Pt(-(dividerMargin + dividerHandleWidth), 0)) // 为左边形状添加偏移
					handleRight := handleShape
					handleRight.Rect = handleRight.Rect.Add(image.Pt(dividerWidth+dividerMargin, 0)) // 为右边形状添加偏移

					paint.FillShape(gtx.Ops, Red200, handleLeft.Op(gtx.Ops))     // 填充左侧形状
					paint.FillShape(gtx.Ops, Yellow100, handleRight.Op(gtx.Ops)) // 填充右侧形状
				}

				// Draw the vertical bar
				// stack3 := clip.Rect{Max: image.Pt(dividerWidth, tallestHeight)}.Push(gtx.Ops)
				// Fill( gtx.Ops, win.Theme.Palette.Table2.Divider) // 如果有需要，在此处可绘制分割线
				// stack3.Pop()
			}
			// 为表头和每列绘制列分隔条
			stack3 := clip.Rect{Max: image.Pt(dividerWidth, tallestHeight)}.Push(gtx.Ops) // 绘制分隔条的矩形区域
			paint.Fill(gtx.Ops, DividerFg)                                                // 填充分隔条的颜色
			stack3.Pop()                                                                  // 弹出分隔条的绘制堆栈

			dividerStart += dividerWidth // 更新分隔符的起始位置
			stack2.Pop()                 // 弹出分隔符的绘制堆栈
			stack.Pop()                  // 弹出当前列的绘制堆栈
		}

		return layout.Dimensions{
			Size: image.Pt(start, tallestHeight), // 返回绘制完成后的整体宽度和高度
		}
	}
}

type Point struct {
	X, Y unit.Dp
}

func (t *TreeTable[T]) ScrollRowIntoView(row int) {
	t.List.ScrollTo(row)
}

const LongPressDuration = 500 * time.Millisecond // 自定义长按持续时间
func (t *TreeTable[T]) RootRows() []*Node[T] {
	if t.filteredRows != nil {
		return t.filteredRows
	}
	return t.rootRows
}

func newID() uuid.ID { return uuid.New('n') }

func (t *TreeTable[T]) IsFiltered() bool { return t.filteredRows != nil }
func (t *TreeTable[T]) Filter(text string) {
	if text == "" {
		t.Root.Children = t.rootRows
		t.Root.OpenAll()
		return
	}
	t.filteredRows = make([]*Node[T], 0)
	for _, node := range t.Root.WalkContainer() { // todo bug 需要改回之前的回调模式？需要调试，编辑节点模态窗口bug
		cells := t.MarshalRowCells(node)
		for _, cell := range cells {
			// mylog.Info(cell.Text, text)
			if strings.EqualFold(cell.Text, text) { // 忽略大小写的情况下相等,支持unicode
				t.filteredRows = append(t.filteredRows, node) // 先过滤所有容器节点
			}
		}
	}
	//for i, row := range t.filteredRows {
	//	break
	//	children := make([]*Node[T], 0)
	//	for _, node := range row.Walk() { // todo bug
	//		columnCells := t.MarshalRowCells(node)
	//		for _, cell := range columnCells {
	//			mylog.Info(cell.Text, text)
	//			if strings.EqualFold(cell.Text, text) {
	//				children = append(children, node) // 过滤子节点
	//			}
	//		}
	//	}
	//	t.filteredRows[i].Children = children
	//}
	if len(t.filteredRows) == 0 {
		return
	}
	// todo 检查layout部分是否调用filteredRows以及filteredRows的大小是否是0，清空过滤后恢复原始的rootRows
	t.Root.Children = t.filteredRows
	// t.rootRows = t.filteredRows
	t.Root.OpenAll()
}

func (t *TreeTable[T]) Sort() {
	if len(t.rootRows) == 0 || t.header.sortedBy >= len(t.rootRows[0].rowCells) {
		return // 如果没有子节点或者列索引无效，直接返回
	}
	sort.Slice(t.rootRows, func(i, j int) bool {
		if t.rootRows[i].rowCells == nil { // why? module do not need this
			t.rootRows[i].rowCells = t.MarshalRowCells(t.rootRows[i])
		}
		if t.rootRows[j].rowCells == nil {
			t.rootRows[j].rowCells = t.MarshalRowCells(t.rootRows[j])
		}
		cellI := t.rootRows[i].rowCells[t.header.sortedBy].Text
		cellJ := t.rootRows[j].rowCells[t.header.sortedBy].Text
		if t.header.sortAscending {
			return cellI < cellJ
		}
		return cellI > cellJ
	})
}

//---------------------------------------泛型n叉树实现------------------------------------------

type Node[T any] struct {
	ID       uuid.ID    // 节点唯一标识符
	Type     string     // 容器节点类型，用于区分不同类型的容器节点
	Data     T          // 元数据
	Children []*Node[T] // 子节点,json只保存以上四个导出的字段
	parent   *Node[T]   // 父节点
	isOpen   bool       // 是否展开
	rowCells []CellData // 行单元格数据
	// rowSelected         bool                     // 行是否被选中
	rowClick widget.Clickable // 行点击按钮绘制,每个节点或者每个单元格对应一个 todo删除，因为popup已经有了
	// cellClickedCallback func(n *Node[T])         // 单元格点击回调，todo 删除?
}

const ContainerKeyPostfix = "_container"

func newRoot[T any](data T) *Node[T] { return NewContainerNode("root", data) }
func (n *Node[T]) IsRoot() bool      { return n.parent == nil }

func NewNode[T any](data T) (child *Node[T]) { return newNode("", false, data) }
func NewContainerNode[T any](typeKey string, data T) (container *Node[T]) {
	n := newNode(typeKey, true, data)
	n.Children = make([]*Node[T], 0)
	return n
}

func NewContainerNodes[T any](typeKeys []string, objects ...T) (containerNodes []*Node[T]) {
	containerNodes = make([]*Node[T], 0)
	var data T // it is zero value
	for i, key := range typeKeys {
		if len(objects) > 0 {
			data = objects[i]
		}
		containerNodes = append(containerNodes, NewContainerNode(key, data))
	}
	return
}

func newNode[T any](typeKey string, isContainer bool, data T) *Node[T] {
	if isContainer {
		typeKey += ContainerKeyPostfix
	}
	n := &Node[T]{
		ID:       newID(),
		Type:     typeKey,
		parent:   nil,
		Data:     data,
		Children: nil,
		isOpen:   isContainer,
		rowCells: nil,
	}
	return n
}

func (n *Node[T]) AddChildByData(data T) { n.AddChild(NewNode(data)) }

func (n *Node[T]) AddChildrenByData(datas ...T) {
	for _, data := range datas {
		n.AddChild(NewNode(data))
	}
}

func (n *Node[T]) AddContainerByData(typeKey string, data T) (newContainer *Node[T]) { // 我们需要返回新的容器节点用于递归填充它的孩子节点，用例是explorer文件资源管理器
	newContainer = NewContainerNode(typeKey, data)
	n.AddChild(newContainer)
	return newContainer
}

func (n *Node[T]) SumChildren() string {
	// container column 0 key is empty string
	k := n.Type
	k = strings.TrimSuffix(k, ContainerKeyPostfix)
	if n.LenChildren() == 0 {
		return k
	}
	k += " (" + fmt.Sprint(n.LenChildren()) + ")"
	return k
}

func (n *Node[T]) UUID() uuid.ID   { return n.ID }
func (n *Node[T]) Container() bool { return strings.HasSuffix(n.Type, ContainerKeyPostfix) }

//func (n *Node[T]) kind(base string) string {
//	if n.Container() {
//		return base + " Container"
//	}
//	return base
//}

func (n *Node[T]) GetType() string           { return n.Type }
func (n *Node[T]) SetType(typeKey string)    { n.Type = typeKey }
func (n *Node[T]) IsOpen() bool              { return n.isOpen && n.Container() }
func (n *Node[T]) SetOpen(open bool)         { n.isOpen = open && n.Container() }
func (n *Node[T]) Parent() *Node[T]          { return n.parent }
func (n *Node[T]) SetParent(parent *Node[T]) { n.parent = parent }

//	func (n *Node[T]) clearUnusedFields() {
//		if !n.Container() {
//			n.Children = nil
//			n.isOpen = false
//		}
//	}
func (n *Node[T]) ResetChildren()                 { n.Children = nil }
func (n *Node[T]) CanHaveChildren() bool          { return n.HasChildren() }
func (n *Node[T]) HasChildren() bool              { return n.Container() && len(n.Children) > 0 }
func (n *Node[T]) CellDataForSort(col int) string { return n.rowCells[col].Text }
func (n *Node[T]) AddChild(child *Node[T]) {
	child.parent = n
	n.Children = append(n.Children, child)
}

func (n *Node[T]) InsertAfter(after *Node[T]) {
	after.parent = n.parent
	n.parent.Children = slices.Insert(n.parent.Children, n.Index()+1, after)
}

func (n *Node[T]) Index() int {
	return slices.Index(n.parent.Children, n)
	//for i, child := range n.parent.Children {
	//	if n.ID == child.ID {
	//		return i
	//	}
	//}
	//panic("not found index") // 永远不可能选中root，所以可以放心panic，root不显示，只显示它的children作为rootRows
}

func (n *Node[T]) Remove() {
	for i, child := range n.parent.Walk() {
		if child.ID == n.ID {
			n.parent.Children = slices.Delete(n.parent.Children, i, i+1)
			break
		}
	}
}

// 通知单元格节点列宽更新事件:
// 增 AddChild(实例化阶段，once内那一次就够了), InsertAfter (DuplicateType) SetChildren
// 删 Remove,需要gtx来重新调整列宽
// 改 EditType,双击或者右键触发
// 查 Find,filter 无需调整列宽
// 过滤,无需调整列宽
// 排序,无需调整列宽

func (t *TreeTable[T]) calcCurrentHierarchyColumnLeftIndent(gtx layout.Context, n *Node[T], cell CellData) unit.Dp {
	leftIndent := unit.Dp(0)
	switch { // 按视觉要求渲染
	case n.parent.IsRoot():
		if n.Container() {
			leftIndent = leftPadding // 有左填充,rootRows不应该紧靠左边，其余的node按视觉要求渲染
		} else {
			leftIndent = leftPadding + defaultHierarchyIconSize
		}
	default:
		depth := n.Depth() - 1
		if n.Container() {
			leftIndent = depth*HierarchyIndent/2 + defaultHierarchyIconSize
		} else {
			leftIndent = leftPadding + depth*HierarchyIndent + defaultHierarchyIconSize
		}
	}
	return leftIndent
}

// 预绘制层级列并返回宽度,然后检查实际绘制的body层级列宽度是否大于最大限制宽度，如果大于则panic，否则更新最大限制宽度,解决放哪是检查层级列单元格是否绘制了额外的图标或者表头层级列加入了排序图标绘制等问题
func (t *TreeTable[T]) calcCurrentHierarchyColumnWidthAndSafeCheck(gtx layout.Context, cell CellData, leftIndent unit.Dp) unit.Dp {
	// 预先计算单元格宽度,todo层级列按需要渲染另外的图标，rootrows深度的左侧还有一个leftPadding没有计算进去，不管了，先这样，反正对齐了
	currentWidth := leftIndent + leftPadding + defaultHierarchyColumnIconSize + LabelWidth(gtx, cell.Text) + DividerWidth // 预计算body层级列实际宽度,todo 对currentWidth进行渲染单元测试
	if currentWidth > t.maxColumnCellWidths[hierarchyColumnID] {                                                          // 这个检查是必要的，否则会导致单元格宽度超出限制导致对不齐表头层级列
		mylog.Todo(currentWidth)
		mylog.Todo(t.maxColumnCellWidths[hierarchyColumnID])
		// mylog.Info("currentWidth", currentWidth)
		// mylog.Info("t.maxColumnCellWidths[hierarchyColumnID]", t.maxColumnCellWidths[hierarchyColumnID])
		// panic("预渲染列宽大于限制宽度")//对于更新节点单元格文本变长的情况，不要panic
	}
	return currentWidth
}

// 预渲染非层级列获取单元格宽度给 updateMaxColumnCellWidth 使用
func (t *TreeTable[T]) calcCurrentNonHierarchyColumnWidthAndSafeCheck(gtx layout.Context, cell CellData) unit.Dp {
	return leftPadding + LabelWidth(gtx, cell.Text) // 每列右侧有空间，不用右填充了，层级列也有
}

// 这避免了重复矩阵置换来提高更新列宽性能，避免增删改查卡顿
func (t *TreeTable[T]) updateMaxColumnCellWidth(gtx layout.Context, n *Node[T]) { // for InsertAfter and edit node late

	if n.Container() {
		t.updateMaxHierarchyColumnCellWidth() // 先更新最大深度再更改宽度，否则层级列的调整不正确
	}

	cells := t.MarshalRowCells(n) // row的cells,取出每个单元格和对应的列max一下确定是否重新调整列宽，这会提高性能
	for columnID, cell := range cells {
		switch columnID {
		case hierarchyColumnID: // 预渲染层级列
			hierarchyColumnWidth := t.calcCurrentHierarchyColumnWidthAndSafeCheck(gtx, cell, t.calcCurrentHierarchyColumnLeftIndent(gtx, n, cell))
			t.maxColumnCellWidths[hierarchyColumnID] = max(t.maxColumnCellWidths[hierarchyColumnID], hierarchyColumnWidth)
		default: ////预渲染其他列
			t.maxColumnCellWidths[columnID] = max(t.maxColumnCellWidths[columnID], t.calcCurrentNonHierarchyColumnWidthAndSafeCheck(gtx, cell))
		}
	}
}

func (t *TreeTable[T]) InsertAfter(gtx layout.Context, after *Node[T]) {
	t.SelectedNode.InsertAfter(after)
	t.updateMaxColumnCellWidth(gtx, after)
}

func (t *TreeTable[T]) Remove(gtx layout.Context) {
	t.SelectedNode.Remove()
	if t.SelectedNode.Container() {
		t.updateMaxHierarchyColumnCellWidth()
	}
}

// /////////////////////////////////////////////////////////

func (t *TreeTable[T]) Edit(gtx layout.Context) { // 编辑节点不会对最大深度有影响
	defer t.updateMaxColumnCellWidth(gtx, t.SelectedNode)
	ModalCallbacks.Reset()
	editor := NewStructView("edit row", t.SelectedNode.Data, // todo merge StructView
		func(a any) []string {
			rowCells := t.MarshalRowCells(t.SelectedNode)
			var rows []string
			for _, cell := range rowCells {
				rows = append(rows, cell.Text)
			}
			return rows
		},
		func(cells []string) any {
			t.UnmarshalRowCells(t.SelectedNode, cells)
			return t.SelectedNode.Data
		},
	)
	editor.Modal = true
	editor.SetOnApply(func() { // todo bug ,debug it
		t.rootRows = t.Root.Children
		t.updateMaxHierarchyColumnCellWidth()
		mylog.Todo("save json data ?")
	})
	ModalCallbacks.Set("node editor", func() {
		if editor != nil && editor.Visible {
			editor.Layout(gtx)
		}
	})
	// todo 其实泽丽不用更新 t.updateMaxHierarchyColumnCellWidth()
	// 但如果编辑的时候把层级列的单元格文本变长就需要，递归一下maxDepth应该不会牺牲多大的性能
	// 如果编辑节点点击应用更新出现卡顿的话，判断下层级列是否被编辑来跳过执行updateMaxHierarchyColumnCellWidth提高性能
}

func (n *Node[T]) Find() (found *Node[T]) {
	for _, child := range n.parent.Walk() {
		if child.ID == n.ID {
			found = child
			break
		}
	}
	return
}

// CopyRow todo 看起来tui的列宽计算需要细细优化
// var rowData = []string{ "Row 4 (5)"  ,"GET" ,"example.com","/api/v4/resource"  ,"application/json","1593"       ,"OK"  ,"获取资源4"       ,"process4.exe"   ,"1m48s",}
// var rowData = []string{ "Sub Sub Row2","GET" ,"example.com","/api/v4/resource1-2","application/json","106"        ,"OK"  ,"获取资源4-1-2"   ,"process4-1-2.exe","7s"   ,}
// var rowData = []string{ "Sub Row3"   ,"GET" ,"example.com","/api/v4/resource3" ,"application/json","106"        ,"OK"  ,"获取资源4-3"     ,"process4-3.exe" ,"7s"   ,}
// var rowData = []string{ "Row5"       ,"GET" ,"example.com","/api/v5/resource"  ,"application/json","104"        ,"OK"  ,"获取资源5"       ,"process5.exe"   ,"5s"   ,}
// var rowData = []string{ "Row11"      ,"GET" ,"example.com","/api/v11/resource" ,"application/json","110"        ,"OK"  ,"获取资源11"      ,"process11.exe"  ,"11s"  ,}
func (n *Node[T]) CopyRow(gtx layout.Context, widths []unit.Dp) string {
	g := stream.NewGeneratedFile()
	g.WriteString("var rowData = []string{ ")
	cells := n.rowCells
	for i, cell := range cells {
		g.WriteString(fmt.Sprintf("%-*s", int(widths[i]), strconv.Quote(cell.Text)) + ",")
	}
	g.P("}")
	gtx.Execute(clipboard.WriteCmd{Data: io.NopCloser(strings.NewReader(g.String()))})
	return g.String()
}

func CountTableRows[T any](rows []*Node[T]) int { // 计算整个表的总行数,可用于设计表格虚拟化滚动来提高性能，但是太复杂了，既要考虑rootrows的行数，又要考虑节点展开后的子节点的行数
	count := len(rows)
	for _, row := range rows {
		if row.CanHaveChildren() {
			count += CountTableRows(row.Children)
		}
	}
	return count
}

func (n *Node[T]) Walk() iter.Seq2[int, *Node[T]] {
	return func(yield func(int, *Node[T]) bool) {
		if !yield(0, n) {
			return
		}
		for i, child := range n.Children {
			if !yield(i, child) { // 迭代索引是为了insert和remove时定位
				break
			}
			if child.CanHaveChildren() {
				// 函数式编程,Walk 方法返回的是一个函数。这个返回的函数接受一个参数（也是一个函数），这个参数就是 yield
				child.Walk()(yield) // 迭代子节点的子节点
			}
		}
	}
}

func (n *Node[T]) Containers() iter.Seq2[int, *Node[T]] {
	return func(yield func(int, *Node[T]) bool) {
		for i, child := range n.Children {
			if child.Container() { // 迭代当前节点下的所有容器节点
				if !yield(i, child) {
					break
				}
			}
		}
	}
}

func (n *Node[T]) WalkContainer() iter.Seq2[int, *Node[T]] {
	return func(yield func(int, *Node[T]) bool) {
		if n.Container() {
			if !yield(0, n) {
				return
			}
		}
		for i, container := range n.Containers() {
			if !yield(i, container) {
				break
			}
			for _, child := range container.Children {
				if child.CanHaveChildren() {
					child.WalkContainer()(yield)
				}
			}
		}
	}
}

//func (n *Node[T]) WalkQueue() iter.Seq[*Node[T]] { // todo 删除，性能应该不行，和Walk结果一样的，理论上
//	return func(yield func(*Node[T]) bool) {
//		queue := []*Node[T]{n}
//		for len(queue) > 0 {
//			node := queue[0]
//			queue = queue[1:]
//			if !yield(node) {
//				break
//			}
//			for _, child := range node.Children {
//				queue = append(queue, child) // 这里将子节点添加到队列
//				if child.CanHaveChildren() { // 如果子节点是一个容器，递归地添加它的子节点
//					for _, subChild := range child.Children {
//						queue = append(queue, subChild)
//					}
//				}
//			}
//		}
//	}
//}

func (n *Node[T]) MaxDepth() unit.Dp { // todo 对root进行10万行的压力测试
	maxDepth := unit.Dp(0)
	for _, node := range n.Walk() {
		childDepth := node.Depth()
		if childDepth > maxDepth {
			maxDepth = childDepth
		}
	}
	return maxDepth
}

func (n *Node[T]) Depth() unit.Dp {
	count := unit.Dp(0)
	p := n.parent
	for p != nil {
		count++
		p = p.parent
	}
	return count
}

func (n *Node[T]) LenChildren() int { return len(n.Children) }
func (n *Node[T]) LastChild() (lastChild *Node[T]) {
	if n.IsRoot() && n.CanHaveChildren() {
		return n.Children[n.LenChildren()-1]
	}
	return n.parent.Children[n.parent.LenChildren()-1]
}
func (n *Node[T]) IsLastChild() bool { return n.LastChild() == n }
func (n *Node[T]) CopyFrom(from *Node[T]) *Node[T] {
	*n = *from
	return n
}

func (n *Node[T]) ApplyTo(to *Node[T]) *Node[T] {
	*to = *n
	return n
}

func (n *Node[T]) OpenAll() {
	for _, node := range n.WalkContainer() {
		node.SetOpen(true)
	}
}

func (n *Node[T]) CloseAll() {
	for _, node := range n.WalkContainer() {
		node.SetOpen(false)
	}
}

func (n *Node[T]) Clone() (to *Node[T]) {
	to = deepcopy.Copy(n)
	to.parent = n.parent
	to.ID = newID()
	if n.CanHaveChildren() {
		n.setParents(to.Children, to, true)
	}
	to.OpenAll()
	return
}

func (n *Node[T]) SetParents(children []*Node[T], parent *Node[T]) {
	n.setParents(children, parent, false)
}

func (n *Node[T]) setParents(children []*Node[T], parent *Node[T], isNewID bool) {
	for _, child := range children {
		child.parent = parent
		if isNewID {
			child.ID = newID()
		}
		if child.CanHaveChildren() {
			n.setParents(child.Children, child, isNewID)
		}
	}
}

// -------------------------tui
const (
	indent          = "│   "
	childPrefix     = "├───"
	lastChildPrefix = "└───"
	indentBase      = unit.Dp(3)
)

func (t *TreeTable[T]) maxColumnCellTextWidth() unit.Dp {
	HierarchyIndent := unit.Dp(1)
	DividerWidth := align.StringWidth[unit.Dp](" │ ")
	iconWidth := align.StringWidth[unit.Dp](childPrefix)
	return t.Root.MaxDepth()*HierarchyIndent + // 最大深度的左缩进
		iconWidth + // 图标宽度,不管深度是多少，每一行都只会有一个层级图标
		t.maxColumnTextWidths[0] + 5 + //(8 * 2) + 20 + // 左右padding,20是sort图标的宽度或者容器节点求和的文本宽度
		DividerWidth // 列分隔条宽度
}

func (t *TreeTable[T]) Format() *stream.Buffer {
	buf := t.FormatHeader(t.maxColumnTextWidths)
	t.FormatChildren(buf, t.rootRows) // 传入子节点打印函数
	// mylog.Json("RootRows", buf.String())
	return buf
}

func (t *TreeTable[T]) String() string { return t.Format().String() }

func (t *TreeTable[T]) Document() string {
	s := stream.NewBuffer("")
	// s.WriteStringLn("// interface or method name here")
	// s.WriteStringLn("/*")
	for line := range t.Format().ToLines() {
		s.WriteStringLn("  " + line)
	}
	// s.WriteStringLn("*/")
	return s.String()
}

func (t *TreeTable[T]) FormatHeader(maxColumnCellTextWidths []unit.Dp) *stream.Buffer {
	buf := stream.NewBuffer("")
	all := t.maxColumnCellTextWidth()
	for _, width := range maxColumnCellTextWidths {
		all += width
	}
	all += align.StringWidth[unit.Dp]("│")*unit.Dp(len(maxColumnCellTextWidths)) + 4 //?
	buf.WriteStringLn("┌─" + strings.Repeat("─", int(all)))
	buf.WriteString("│")

	// 计算每个单元格的左边距
	for i, cell := range t.header.columnCells {
		paddedText := fmt.Sprintf("%-*s", int(maxColumnCellTextWidths[i]), cell.Text) // 左对齐填充

		// 添加左边距，仅在首列进行处理，依据列宽计算
		if i == hierarchyColumnID {
			buf.WriteString(strings.Repeat(" ", int(t.maxColumnCellTextWidth()-maxColumnCellTextWidths[i]-1))) // -1是分隔符的空间
		}

		buf.WriteString(paddedText)
		if i < t.columnCount-1 {
			buf.WriteString(" │ ") // 在每个单元格之间添加分隔符
		}
	}

	buf.NewLine()
	buf.WriteStringLn("├─" + strings.Repeat("─", int(all)))
	return buf
}

func (t *TreeTable[T]) FormatChildren(out *stream.Buffer, children []*Node[T]) {
	for i, child := range children {
		child.rowCells = t.MarshalRowCells(child)
		HierarchyColumBuf := stream.NewBuffer("")
		for j, cell := range child.rowCells {
			if j == hierarchyColumnID {
				HierarchyColumBuf.WriteString("│")
				HierarchyColumBuf.WriteString(strings.Repeat(" ", int(child.Depth()*indentBase)))
				if i == len(children)-1 {
					HierarchyColumBuf.WriteString("╰──") //"└───"
				} else {
					HierarchyColumBuf.WriteString("├──")
				}
				HierarchyColumBuf.WriteString(cell.Text)
				if align.StringWidth[unit.Dp](HierarchyColumBuf.String()) < t.maxColumnCellTextWidth() {
					HierarchyColumBuf.WriteString(strings.Repeat(" ", int(t.maxColumnCellTextWidth()-align.StringWidth[unit.Dp](HierarchyColumBuf.String()))))
				}
				HierarchyColumBuf.WriteString(" │ ")
				out.WriteString(HierarchyColumBuf.String())
				HierarchyColumBuf.Reset()
				continue
			}
			out.WriteString(cell.Text)
			if align.StringWidth[unit.Dp](cell.Text) < t.maxColumnTextWidths[j] {
				out.WriteString(strings.Repeat(" ", int(t.maxColumnTextWidths[j]-align.StringWidth[unit.Dp](cell.Text))))
			}
			out.WriteString(" │ ")
		}
		out.NewLine()
		if len(child.Children) > 0 {
			t.FormatChildren(out, child.Children)
		}
	}
}
