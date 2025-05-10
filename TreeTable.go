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
	"github.com/ddkwork/ux/resources/colors"
	"github.com/ddkwork/ux/resources/images"
	"github.com/ddkwork/ux/widget/material"
	"github.com/google/go-cmp/cmp"
)

type (
	TreeTable[T any] struct {
		TableContext[T]                     // 实例化时传入的上下文
		Root                *Node[T]        // 根节点,保存数据到json只需要调用它即可
		header              tableHeader[T]  // 表头
		rootRows            []*Node[T]      // from root.children
		filteredRows        []*Node[T]      // 过滤后的行
		rootRowsWidget      []layout.Widget // from layout
		rows                [][]CellData    // 为了求最大列宽，把他当做全部节点展开的状态
		SelectedNode        *Node[T]        // 当前选中的节点
		contextMenu         *ContextMenu    // 行右键菜单,rootRows只需要一个
		columnCount         int             // 列数
		maxColumnCellWidths []unit.Dp       // 最宽的列label文本宽度for单元格
		maxColumnTextWidths []unit.Dp       // 最宽的列文本宽度for tui
		once                sync.Once       // 自动计算列宽一次
		// drag ........
		// inLayoutHeader               bool                     // for drag
		// columnResizeStart            unit.Dp                  //
		// columnResizeBase             unit.Dp                  //
		// columnResizeOverhead         unit.Dp                  //
		// preventUserColumnResize      bool                     //
		// awaitingSyncToModel          bool                     //
		// wasDragged                   bool                     //
		// dividerDrag                  bool                     //
	}
	TableContext[T any] struct {
		CustomContextMenuItems func(gtx layout.Context, n *Node[T]) iter.Seq[ContextMenuItem] // 通过SelectedNode传递给菜单的do取出元数据，比如删除文件,但是菜单是否绘制取决于当前渲染的行，所以要传递n给can
		MarshalRowCells        func(n *Node[T]) (cells []CellData)                            // 序列化节点元数据
		UnmarshalRowCells      func(n *Node[T], rows []CellData) T                            // 节点编辑后反序列化回节点
		RowSelectedCallback    func()                                                         // 行选中回调,通过SelectedNode传递给菜单
		RowDoubleClickCallback func()                                                         // double click callback,通过SelectedNode传递给菜单
		SetRootRowsCallBack    func()                                                         // 实例化所有节点回调,必要时调用root节点辅助操作
		JsonName               string                                                         // 保存序列化树形表格到文件的文件名
		IsDocument             bool                                                           // 是否生成markdown文档
		//	DragRemovedRowsCallback  func()
		//	DropOccurredCallback     func()
	}
	tableHeader[T any] struct {
		sortOrder          sortOrder    // 排序方式
		sortedBy           int          // 排序列索引
		columnCells        []CellData   // 每列的最大深度，宽度,也可以把它看做一行
		clickedColumnIndex int          // 被点击的列索引
		sortAscending      bool         // 升序还是降序
		contextMenu        *ContextMenu // 右键菜单，实现复制列数据到剪贴板
		node               *Node[T]     // 列宽同步传参需要
		// drags              []tableDrag  // 拖动表头列参数
		// manualWidthSet     []bool       // 新增状态标志数组，记录列是否被手动调整
	}
	CellData struct {
		Key              string      // 表头单元格文本或者节点编辑器每一行的key
		Value            string      // 单元格文本或者节点编辑器每一行的value
		Tooltip          string      // 单元格提示信息
		Icon             []byte      // 单元格图标，格式支持：*giosvg.Icon, *widget.Icon, *widget.Image, image.Image
		FgColor          color.NRGBA // 单元格前景色,着色渲染单元格
		IsNasm           bool        // 是否是nasm汇编代码,为表头提供不同的着色渲染样式
		Disabled         bool        // 是否显示表头或者body单元格，或者禁止编辑节点时候使用
		isHeader         bool        // 是否是表头单元格
		columID          int         // 列id,预计后期用于区域选中,每行的列数和表头是一样的，而表头单元格在init的时候填充了每列的id，那么body rows的列id只要在操作row节点的时候遍历一下表头的row node填充列id即可
		rowID            int         // 行id,预计后期用于区域选中,由layout遍历list渲染的index填充该节点的rowID
		widget.Clickable             // todo 单元格点击事件,如果换成编辑框，那么编辑框需要支持单机事件和双击编辑回车事件,以及RichText高亮单元格
		RichText                     // todo 单元格富文本
		width            unit.Dp
	}

	// tableDrag struct {
	// 	drag           gesture.Drag
	// 	hover          gesture.Hover
	// 	startPos       float32
	// 	shrinkNeighbor bool
	// }
)

func NewTreeTable[T any](data T) *TreeTable[T] {
	columnCells := InitHeader(data)
	columnCount := len(columnCells)
	n := NewNode(data)
	n.rowCells = columnCells
	root := newRoot(data)
	n.SetParent(root)
	return &TreeTable[T]{
		TableContext: TableContext[T]{},
		Root:         root,
		header: tableHeader[T]{
			sortOrder:          0,
			sortedBy:           0,
			columnCells:        columnCells,
			clickedColumnIndex: -1,
			sortAscending:      false,
			contextMenu:        NewContextMenu(),
			node:               n,
		},
		rootRows:       nil,
		filteredRows:   nil,
		rootRowsWidget: nil,
		rows:           nil,
		// columns:                      nil,
		SelectedNode:        nil,
		contextMenu:         NewContextMenu(),
		columnCount:         columnCount,
		maxColumnCellWidths: make([]unit.Dp, columnCount),
		maxColumnTextWidths: make([]unit.Dp, columnCount),
		once:                sync.Once{},
	}
}

func (t *TreeTable[T]) makeRootRowsWidget() {
	t.filteredRows = nil
	if t.rootRowsWidget == nil {
		t.rootRowsWidget = make([]layout.Widget, t.RootRowCount())
	}
	if len(t.rootRowsWidget) < t.RootRowCount() {
		diff := t.RootRowCount() - len(t.rootRowsWidget)
		sub := make([]layout.Widget, diff)
		t.rootRowsWidget = slices.Concat(t.rootRowsWidget, sub)
	}
	for i, row := range t.RootRows() {
		t.rootRowsWidget[i] = func(gtx layout.Context) layout.Dimensions {
			return t.RowFrame(gtx, row, i)
		}
	}
}

func (t *TreeTable[T]) RootRows() []*Node[T] {
	if t.filteredRows != nil {
		return t.filteredRows
	}
	t.rootRows = t.Root.Children
	return t.rootRows
}

func (t *TreeTable[T]) RootRowCount() int {
	return len(t.RootRows())
}

func (t *TreeTable[T]) Layout(gtx layout.Context) layout.Dimensions {
	t.once.Do(func() {
		if t.SetRootRowsCallBack != nil { // github.com/ddkwork/mitmproxy
			t.SetRootRowsCallBack()
		}
		if t.JsonName == "" {
			mylog.Check("JsonName is empty")
		}
		if t.MarshalRowCells == nil {
			t.MarshalRowCells = func(n *Node[T]) (cells []CellData) {
				return MarshalRow(n.Data, nil)
			}
		}
		if t.UnmarshalRowCells == nil {
			t.UnmarshalRowCells = func(n *Node[T], rows []CellData) T {
				return UnmarshalRow[T](rows, nil)
			}
		}
		mylog.CheckNil(t.RowSelectedCallback)

		// if t.FileDropCallback == nil {
		//	t.FileDropCallback = func(files []string) {
		//		if filepath.Ext(files[0]) == ".json" {
		//			mylog.Info("dropped file", files[0])
		//			table.ResetChildren()
		//			b := stream.NewBuffer(files[0])
		//			mylog.Check(json.Unmarshal(b.Bytes(), table))
		//		}
		//		mylog.Struct("todo", files)
		//	}
		// }
	})
	t.SizeColumnsToFit(gtx)
	t.makeRootRowsWidget()

	for _, n := range t.rootRows {
		t.contextMenu.Once.Do(func() {
			item := ContextMenuItem{}
			for _, kind := range CopyRowType.EnumTypes() {
				switch kind {
				case CopyRowType:
					item = ContextMenuItem{
						Title:     "",
						Icon:      images.SvgIconCopy,
						Can:       func() bool { return true },
						Do:        func() { t.SelectedNode.CopyRow(gtx, t.maxColumnTextWidths) },
						Clickable: widget.Clickable{},
					}
				case ConvertToContainerType:
					item = ContextMenuItem{
						Title: "",
						Icon:  images.SvgIconConvertToContainer,
						Can:   func() bool { return !n.Container() }, // n是当前渲染的行
						Do: func() {
							t.SelectedNode.SetType("ConvertToContainer" + ContainerKeyPostfix) // ? todo bug：这里是失败的，导致再次点击这里转换的节点后ConvertToNonContainer没有弹出来
							t.SelectedNode.ID = newID()
							t.SelectedNode.SetOpen(true)
							t.SelectedNode.Children = make([]*Node[T], 0)
						},
						Clickable: widget.Clickable{},
					}
				case ConvertToNonContainerType:
					item = ContextMenuItem{
						Title: "",
						Icon:  images.SvgIconConvertToNonContainer,
						Can:   func() bool { return n.Container() }, // n是当前渲染的行
						Do: func() {
							t.SelectedNode.SetType("")
							t.SelectedNode.ID = newID()
							for _, child := range t.SelectedNode.Children {
								child.parent = t.SelectedNode.parent
								child.ID = newID()
							}
							t.SelectedNode.ResetChildren()
						},
						AppendDivider: true,
						Clickable:     widget.Clickable{},
					}
				case NewType:
					item = ContextMenuItem{
						Title: "",
						Icon:  images.SvgIconCircledAdd,
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
						Icon:  images.SvgIconCircledVerticalEllipsis,
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
						Icon:      images.SvgIconTrash,
						Can:       func() bool { return true },
						Do:        func() { t.Remove(gtx) },
						Clickable: widget.Clickable{},
					}
				case DuplicateType:
					item = ContextMenuItem{
						Title: "",
						Icon:  images.SvgIconDuplicate,
						Can:   func() bool { return true },
						Do: func() {
							t.InsertAfter(gtx, t.SelectedNode.Clone())
						},
						Clickable: widget.Clickable{},
					}
				case EditType:
					item = ContextMenuItem{
						Title:         "",
						Icon:          images.SvgIconEdit,
						Can:           func() bool { return true },
						Do:            func() { t.Edit(gtx) },
						AppendDivider: true,
						Clickable:     widget.Clickable{},
					}
				case OpenAllType:
					item = ContextMenuItem{
						Title:     "",
						Icon:      images.SvgIconHierarchy,
						Can:       func() bool { return true },
						Do:        func() { t.OpenAll() },
						Clickable: widget.Clickable{},
					}
				case CloseAllType:
					item = ContextMenuItem{
						Title:     "",
						Icon:      images.SvgIconCircledVerticalEllipsis,
						Can:       func() bool { return true },
						Do:        func() { t.Root.CloseAll() },
						Clickable: widget.Clickable{},
					}
				case SaveDataType:
					item = ContextMenuItem{
						Title:     "",
						Icon:      images.SvgIconSaveContent,
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
			if items := t.CustomContextMenuItems(gtx, n); items != nil { // gtx用于调整列宽，触发增删改之后，n是正在渲染的节点，用于取出元数据控制菜单是否渲染
				for item := range items {
					if item.Can() {
						t.contextMenu.AddItem(item)
					}
				}
			}
		})
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return t.HeaderFrame(gtx) // 渲染表头	// t.inLayoutHeader = true
		// return t.layoutDrag(gtx, func(gtx layout.Context, row int) layout.Dimensions {
		// 	return t.HeaderFrame(gtx) // 渲染表头
		// })
	}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.contextMenu.Layout(gtx, t.rootRowsWidget)
			// ////////////////////////
			// t.inLayoutHeader = false
			// return t.layoutDrag(gtx, func(gtx layout.Context, row int) layout.Dimensions {
			//	return t.RowFrame(gtx, t.rootRows[index], index)
			// })
			// })
		}),
	)
}

func (t *TreeTable[T]) RowFrame(gtx layout.Context, n *Node[T], rowIndex int) layout.Dimensions {
	if n.rowCells == nil {
		n.rowCells = t.MarshalRowCells(n)
	}
	for i := range n.rowCells {
		n.rowCells[i].rowID = rowIndex                          // 行单元格id，暂时想到的应用场景是区域选中
		n.rowCells[i].columID = t.header.columnCells[i].columID // 列分隔符,更新层级列宽度
		n.rowCells[i].Key = t.header.columnCells[i].Key         // 为节点编辑的row布局kv对做准备

	}

	bgColor := t.processEvent(gtx, n) // 处理鼠标事件
	rowClick := &n.rowClick

	var rowCells []layout.FlexChild

	layoutHierarchyColumn := func(gtx layout.Context, cell CellData) layout.Dimensions {
		indent, width := t.cellWidth(gtx, n, &cell)
		cell.width = width
		gtx.Constraints.Min.X = int(width)
		gtx.Constraints.Max.X = int(width)
		return layout.Flex{
			Axis: layout.Horizontal,
			// Spacing:   0,
			// Alignment: layout.Start,
			// WeightSum: 0,
		}.Layout(gtx, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return rowClick.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// 绘制层级图标-----------------------------------------------------------------------------------------------------------------
				HierarchyInsert := layout.Inset{Left: indent, Top: -1} // 层级图标居中,行高调整后这里需要下移使得图标居中
				if !n.Container() {
					return HierarchyInsert.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Dimensions{}
					})
				}
				return HierarchyInsert.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					svg := images.SvgIconCircledChevronRight
					if n.isOpen {
						svg = images.SvgIconCircledChevronDown
					}
					return iconButtonSmall(new(widget.Clickable), svg, "").Layout(gtx)
					// return NewButton("", nil).SetRectIcon(true).SetIcon(svg).Layout(gtx)
				})
			})
		}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return rowClick.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return t.CellFrame(gtx, &cell)
				})
			}),
		)
	}

	rowCells = append(rowCells, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		hierarchyColumnCell := n.rowCells[HierarchyColumnID] // 层级列单元格
		return layoutHierarchyColumn(gtx, hierarchyColumnCell)
	}))

	// 绘制非层级列-----------------------------------------------------------------------------------------------------------------
	for i, cell := range n.rowCells {
		if i == HierarchyColumnID {
			continue
		}
		rowCells = append(rowCells, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return rowClick.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Clickable(gtx, &n.rowCells[i].Clickable, func(gtx layout.Context) layout.Dimensions {
					DrawColumnDivider(gtx, cell.columID) // 这里绘制的列分割线才没有虚线,gtx被破坏了? 永远不要移动这个位置
					if len(cell.Value) > 80 {
						cell.Value = cell.Value[:len(cell.Value)/2] + "..."
						// todo 更好的办法是让富文本编辑器做这个事情，对 maxline 。。。 看看代码编辑器扩建是如何实现这个的
						// 然后双击编辑行的时候从富文本取出完整行并换行显示，structView需要好好设计一下这个
						// 这个在抓包场景很那个，url列一般都长
					}
					_, cell.width = t.cellWidth(gtx, n, &cell)
					return t.CellFrame(gtx, &cell)
				})
			})
		}))
	}

	rows := []layout.FlexChild{ // 合成层级列和其他列的单元格为一行,并设置该行的背景和行高
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// mylog.Trace(rowIndex, bgColor)
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
				return t.RowFrame(gtx, child, rowIndex)
			}))
		}
	}
	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx, rows...)
}

func (t *TreeTable[T]) updateRowNumber(parent *Node[T], currentRowIndex int) int {
	parent.RowNumber = currentRowIndex
	currentRowIndex++
	if parent.CanHaveChildren() {
		for _, child := range parent.Children {
			currentRowIndex = t.updateRowNumber(child, currentRowIndex)
		}
	}
	return currentRowIndex
}

func CountTableRows[T any](rows []*Node[T]) int {
	count := len(rows)
	for _, row := range rows {
		if row.CanHaveChildren() {
			count += CountTableRows(row.Children)
		}
	}
	return count
}

func (t *TreeTable[T]) processEvent(gtx layout.Context, n *Node[T]) color.NRGBA {
	rowClick := &n.rowClick
	evt, ok := gtx.Source.Event(pointer.Filter{
		Target: rowClick,
		Kinds:  pointer.Press, // | pointer.Release | pointer.Drag,
	})
	if ok {
		e, ok := evt.(pointer.Event)
		if ok {
			switch {
			case e.Kind == pointer.Press: // 左键，右键，双击
				t.SelectedNode = n
				// bgColor = Orange300
			case e.Source == pointer.Touch:
				t.SelectedNode = n
			}
		}
	}

	click, ok := rowClick.Update(gtx)
	if ok {
		switch click.NumClicks {
		case 1:
			n.isOpen = !n.isOpen // 切换展开状态
			if t.RowSelectedCallback != nil {
				t.RowSelectedCallback() // 行选中回调
			}

		case 2:
			t.Edit(gtx)
			// if t.RowDoubleClickCallback != nil { // 行双击回调
			//	go t.RowDoubleClickCallback(n)
			//	//gtx.Execute(op.InvalidateCmd{})
			// }
		}
	}
	bgColor := RowColor(n.RowNumber)
	switch {
	case t.SelectedNode == n: // 设置选中背景色,这个需要在第一的位置,在选中另外的节点时间段内这个条件成立，能保持背景色不变
		bgColor = color.NRGBA{R: 200, G: 145, B: 50, A: 255}
	case rowClick.Hovered(): // 设置悬停背景色
		// https://rgbacolorpicker.com/
		bgColor = th.Color.TreeHoveredBgColor
	default:
	}
	return bgColor
}

func (t *TreeTable[T]) IsRowSelected() bool { return t.SelectedNode != nil }

func (t *TreeTable[T]) CellFrame(gtx layout.Context, cell *CellData) layout.Dimensions {
	// 固定单元格宽度为计算好的每列最大label宽度实现表头和body对齐,因为表头和body都调用这个函数渲染单元格，只有限制min和max才能每列保证表头单元格和body单元格具有相等的宽度，从而实现表头和body对齐
	if cell.width == 0 { // mitmproxy start
		// panic("cell width is 0")
		return layout.Dimensions{}
	}
	gtx.Constraints.Min.X = int(cell.width)
	gtx.Constraints.Max.X = int(cell.width)
	DrawColumnDivider(gtx, cell.columID) // 为每列绘制列分隔条

	if cell.FgColor == (color.NRGBA{}) {
		cell.FgColor = colors.White
	}
	// richText := NewRichText()
	//
	// if cell.IsNasm {
	// 	tokens, style := languages.GetTokens(stream.NewBuffer(cell.Value), languages.NasmKind)
	// 	for _, token := range tokens {
	// 		colour := style.Get(token.Type).Colour
	// 		color := color.NRGBA{
	// 			R: colour.Red(),
	// 			G: colour.Green(),
	// 			B: colour.Blue(),
	// 			A: 255,
	// 		}
	// 		richText.AddSpan(richtext.SpanStyle{
	// 			// Font:        font.Font{},
	// 			Size:        unit.Sp(12),
	// 			Color:       color,
	// 			Content:     cell.Value,
	// 			Interactive: false,
	// 		})
	// 	}
	// 	return inset.Layout(gtx, richText.Layout)
	// }
	// richText.AddSpan(richtext.SpanStyle{
	// 	// Font:        font.Font{},
	// 	Size:        unit.Sp(12),
	// 	Color:       cell.FgColor,
	// 	Content:     cell.Value,
	// 	Interactive: false,
	// })
	//
	// return inset.Layout(gtx, richText.Layout)

	return layout.Flex{
		Axis:      layout.Horizontal,
		Spacing:   0,
		Alignment: 0,
		WeightSum: 0,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if cell.Icon == nil { // 格式支持：*giosvg.Icon, *widget.Icon, *widget.Image, image.Image
				return layout.Dimensions{}
			}
			size := image.Pt(gtx.Dp(defaultHierarchyColumnIconSize), gtx.Dp(defaultHierarchyColumnIconSize))
			gtx.Constraints.Min = size
			gtx.Constraints.Max = size
			return layout.Inset{
				Top:    topPadding,
				Bottom: 0,
				Left:   leftPadding,
				// Right:  rightPadding,
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return images.Layout(gtx, cell.Icon, color.NRGBA{}, defaultHierarchyColumnIconSize)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			top := topPadding
			bottom := bottomPadding
			if stream.IsAndroid() {
				top /= 2
			}
			if cell.isHeader {
				return layout.Inset{
					Top:    top,
					Bottom: bottom,
					Left:   leftPadding,
					Right:  0,
				}.Layout(gtx, material.Body2(th, cell.Key).Layout)
			}
			return layout.Inset{
				Top:    top,
				Bottom: 0,
				Left:   leftPadding,
				Right:  0,
			}.Layout(gtx, material.Body2(th, cell.Value).Layout)
		}),
	)
}

func InitHeader(data any) (columnCells []CellData) {
	fields := stream.ReflectVisibleFields(data)
	columnCells = make([]CellData, 0)
	for i, field := range fields {
		if field.Tag.Get("table") != "" { // 中文表头简短
			field.Name = field.Tag.Get("table")
		}
		columnCells = append(columnCells, CellData{
			Key:       field.Name,
			Value:     field.Name,
			Tooltip:   "",
			Icon:      nil,
			FgColor:   color.NRGBA{},
			IsNasm:    false,
			Disabled:  false,
			isHeader:  true,
			columID:   i,
			rowID:     0,
			Clickable: widget.Clickable{},
			RichText:  RichText{},
		})
	}
	return
}

// SizeColumnsToFit todo not explort for now
func (t *TreeTable[T]) SizeColumnsToFit(gtx layout.Context) { // 增删改查中，只有增改+实例化，共3次需要调用这个函数，增改需要更新缓存的每列最大列宽即可，所以这个函数理论上只需要执行一次，这样性能最好
	if t.RootRowCount() == 0 {
		return // mitmproxy start
	}
	originalConstraints := gtx.Constraints // 保存原始约束\
	if t.rows == nil {
		t.updateRowNumber(t.Root, 0)
		t.rows = make([][]CellData, t.Root.LastChild().RowNumber)
	}
	count := t.Root.LastChild().RowNumber
	if len(t.rows) < count {
		t.updateRowNumber(t.Root, 0)
		subLen := count - len(t.rows)
		sub := make([][]CellData, subLen)
		t.rows = slices.Concat(t.rows, sub)
		t.updateRowNumber(t.Root, 0)
		count = t.Root.LastChild().RowNumber
	}
	for i := range count { // todo  1437	   1256535 ns/op
		for _, node := range t.Root.Walk() {
			if t.rows[i] == nil { // this will not safe for edit node data later
				t.rows[i] = t.MarshalRowCells(node)
			}
		}
	}

	// t.columns = TransposeMatrix(t.rows) // 如果不这么做的话，节点增删改查就不会实时刷新,为了提高性能需要手动刷新节点和宽度
	for i, data := range TransposeMatrix(t.rows) { // todo 1561	    731960 ns/op
		if data.isHeader {
			t.maxColumnTextWidths[i] = max(t.maxColumnTextWidths[i], align.StringWidth[unit.Dp](data.Key), align.StringWidth[unit.Dp](t.header.columnCells[i].Key))
			if gtx != (layout.Context{}) {
				t.maxColumnCellWidths[i] = max(t.maxColumnCellWidths[i], LabelWidth(gtx, data.Key), LabelWidth(gtx, t.header.columnCells[i].Key))
			}
		} else {
			t.maxColumnTextWidths[i] = max(t.maxColumnTextWidths[i], align.StringWidth[unit.Dp](data.Value), align.StringWidth[unit.Dp](t.header.columnCells[i].Value))
			if gtx != (layout.Context{}) {
				t.maxColumnCellWidths[i] = max(t.maxColumnCellWidths[i], LabelWidth(gtx, data.Value), LabelWidth(gtx, t.header.columnCells[i].Value))
			}
		}
	}
	gtx.Constraints = originalConstraints
}

func (t *TreeTable[T]) SaveDate() {
	go func() {
		t.JsonName = strings.TrimSuffix(t.JsonName, ".json")
		t.JsonName += ".json"
		if stream.IsAndroid() {
			go func() {
				f := mylog.Check2(explore.CreateFile(t.JsonName))
				mylog.Check2(f.Write([]byte(t.String())))
				mylog.Check(f.Close())
			}()
			return
		}

		stream.MarshalJsonToFile(t.Root, filepath.Join("cache", t.JsonName))
		stream.WriteTruncate(filepath.Join("cache", t.JsonName+".txt"), t.Document()) // 调用t.Format()
		if t.IsDocument {
			b := stream.NewBuffer("")
			b.WriteStringLn("# " + t.JsonName + " document table")
			b.WriteStringLn("```text")
			b.WriteStringLn(t.Document())
			b.WriteStringLn("```")
			stream.WriteTruncate("README2.md", b.String())
		}
	}()
}

// TransposeMatrix 把行切片矩阵置换为列切片,用于计算最大列宽的参数
// https://github.com/BaseMax/SparseMatrixLinkedListGo
func TransposeMatrix[T any](rootRows [][]T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for _, row := range rootRows {
			for j, v := range row {
				if !yield(j, v) { // j 是列索引，作为转置后的行键
					return
				}
			}
		}
	}
}

func (t *TreeTable[T]) HeaderFrame(gtx layout.Context) layout.Dimensions {
	// 实现右键复制列到剪切板功能
	t.header.contextMenu.Once.Do(func() {
		t.header.contextMenu.AddItem(ContextMenuItem{
			Title: "CopyColumn",
			Icon:  images.SvgIconCopy,
			Can:   func() bool { return true },
			Do: func() {
				t.CopyColumn(gtx)
				t.header.clickedColumnIndex = -1 // 重置点击列索引
			},
			AppendDivider: true,
			Clickable:     widget.Clickable{},
		})
	})
	return t.header.contextMenu.Layout(gtx, []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			var cols []layout.FlexChild
			elems := make([]*Resizable, 0)
			for i, cell := range t.header.columnCells {
				_, cell.width = t.cellWidth(gtx, t.header.node, &cell)
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
					for j := range t.header.columnCells {
						t.header.columnCells[j].Key = strings.TrimSuffix(t.header.columnCells[j].Key, " ⇩")
						t.header.columnCells[j].Key = strings.TrimSuffix(t.header.columnCells[j].Key, " ⇧")
						t.header.columnCells[j].Value = strings.TrimSuffix(t.header.columnCells[j].Value, " ⇩")
						t.header.columnCells[j].Value = strings.TrimSuffix(t.header.columnCells[j].Value, " ⇧")
					}
					// mylog.Info(t.header.columnCells[j].Value, LabelWidth(gtx, t.header.columnCells[j].Value))
					if t.header.clickedColumnIndex > -1 && t.header.sortedBy > -1 {
						switch t.header.sortOrder {
						case sortNone:
						case sortAscending:
							t.header.columnCells[i].Key += " ⇧"
						case sortDescending:
							t.header.columnCells[i].Key += " ⇩"
						}
						t.Sort()
					}
				}
				// 拦截右击事件并在事件中赋值命中的列id
				for {
					evt, ok := gtx.Source.Event(pointer.Filter{Target: clickable, Kinds: pointer.Press | pointer.Release})
					if !ok {
						break
					}
					e, ok := evt.(pointer.Event)
					if !ok {
						break
					}
					if e.Kind == pointer.Press || e.Source == pointer.Touch {
						t.header.clickedColumnIndex = i
					}
					switch e.Buttons {
					case pointer.ButtonPrimary:

					case pointer.ButtonSecondary:
						if e.Kind == pointer.Press {
							t.header.clickedColumnIndex = i
						}
					default:
					}
				}

				cols = append(cols, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return Background{Color: colors.ColorHeaderFg}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return material.Clickable(gtx, clickable, func(gtx layout.Context) layout.Dimensions {
							cellFrame := t.CellFrame(gtx, &cell)
							elems = append(elems, &Resizable{Widget: func(gtx layout.Context) layout.Dimensions {
								return cellFrame
							}})
							return cellFrame
						})
					})
				}))
			}
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, WeightSum: 0}.Layout(gtx, cols...)
		},
	})
}

func (t *TreeTable[T]) cellWidth(gtx layout.Context, n *Node[T], cell *CellData) (indent, width unit.Dp) {
	defer func() {
		t.header.columnCells[cell.columID].width = max(t.header.columnCells[cell.columID].width, width)
		cell.width = t.header.columnCells[cell.columID].width
		gtx.Execute(op.InvalidateCmd{})
	}()
	v := cell.Value
	if cell.isHeader {
		v = cell.Key
	}
	v += "  ⇧"
	switch cell.columID {
	case HierarchyColumnID:
		leftIndent := HierarchyIndent*(n.Depth()-1) + leftPadding
		if n.Parent().IsRoot() && !n.CanHaveChildren() {
			leftIndent = HierarchyIndent * (n.Depth())
		} else {
			if !n.CanHaveChildren() {
				leftIndent -= leftPadding
				if n.parent.CanHaveChildren() {
					leftIndent += defaultHierarchyColumnIconSize
				}
			}
		}
		if n.CanHaveChildren() {
			if !strings.HasSuffix(cell.Value, ")") {
				panic("you should call n.SumChildren method")
			}
			// v = n.SumChildren() + "  ⇧"
		}
		labelWidth := LabelWidth(gtx, v)
		current := leftIndent + labelWidth + DividerWidth
		maxWidth := (t.MaxDepth()-1)*HierarchyIndent + // 最大深度的左缩进
			defaultHierarchyIconSize + // 层级图标宽度
			defaultHierarchyColumnIconSize + // 层级列图标宽度
			leftPadding + // 左侧padding
			labelWidth + // 最长的单元格文本宽度
			rightPadding + // 右侧padding
			DividerWidth // 列分隔条宽度

		t.maxColumnCellWidths[HierarchyColumnID] = max(maxWidth, current, t.header.columnCells[HierarchyColumnID].width)
		t.header.columnCells[HierarchyColumnID].width = t.maxColumnCellWidths[HierarchyColumnID]

		if current > t.maxColumnCellWidths[HierarchyColumnID] {
			mylog.Todo(cmp.Diff(t.maxColumnCellWidths[HierarchyColumnID], current)) // 预渲染列宽大于限制宽度
		}
		gtx.Execute(op.InvalidateCmd{})
		return leftIndent, t.maxColumnCellWidths[HierarchyColumnID]
	default:
		labelWidth := LabelWidth(gtx, v)
		w := leftPadding + labelWidth + DividerWidth
		t.maxColumnCellWidths[cell.columID] = max(t.maxColumnCellWidths[cell.columID], w) // todo  runtime error: index out of range [6] with length 5 in jsonTreeTable.go:112
		gtx.Execute(op.InvalidateCmd{})
		return 0, t.maxColumnCellWidths[cell.columID]
	}
}

const (
	HierarchyColumnID              = 0
	leftPadding                    = unit.Dp(4)  // 单元格左侧填充，和列分割线的间距
	rightPadding                   = unit.Dp(4)  // 单元格右侧填充，和列分割线的间距
	topPadding                     = unit.Dp(4)  // 单元格上侧填充，似单元格间文本居中
	bottomPadding                  = unit.Dp(4)  // 单元格下侧填充，似单元格间文本居中
	defaultHierarchyIconSize       = unit.Dp(12) // 层级图标默认大小
	defaultHierarchyColumnIconSize = unit.Dp(16) // 层级列图标默认大小
	HierarchyIndent                = unit.Dp(16) // 层级缩进
	defaultRowHeight               = unit.Dp(22) // 行高
	// defaultHeaderHeight            = unit.Dp(24)  // 表头行高
	// defaultHeaderFontSize          = unit.Sp(12)  // 表头字体大小
	// defaultRowFontSize             = unit.Sp(12)  // 行字体大小
)

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
		paint.Fill(gtx.Ops, colors.DividerFg)
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
	g.P(strconv.Quote(t.header.columnCells[t.header.clickedColumnIndex].Value), ",")
	for i, datum := range TransposeMatrix(t.rows) {
		if i == t.header.clickedColumnIndex {
			g.P(strconv.Quote(datum.Value), ",")
		}
	}
	g.P("}")
	gtx.Execute(clipboard.WriteCmd{Data: io.NopCloser(strings.NewReader(g.Format()))})
	return g.String()
}

func (t *TreeTable[T]) ScrollRowIntoView(row int) {
	t.contextMenu.list.ScrollTo(row)
}

func newID() uuid.ID { return uuid.New('n') }

func (t *TreeTable[T]) IsFiltered() bool { return t.filteredRows != nil }
func (t *TreeTable[T]) Filter(text string) {
	if text == "" {
		t.filteredRows = nil
		t.Root.Children = t.rootRows
		t.OpenAll()
		return
	}
	t.filteredRows = make([]*Node[T], 0)          // todo root row is not container case handle
	for _, node := range t.Root.WalkContainer() { // todo bug 需要改回之前的回调模式？需要调试，编辑节点模态窗口bug
		if node.rowCells == nil {
			node.rowCells = t.MarshalRowCells(node)
		}
		for _, cell := range node.rowCells {
			// mylog.Info(cell.Value, text)
			if strings.EqualFold(cell.Value, text) { // 忽略大小写的情况下相等,支持unicode
				t.filteredRows = append(t.filteredRows, node) // 先过滤所有容器节点
			}
		}
	}
	// for i, row := range t.filteredRows {
	//	break
	//	children := make([]*Node[T], 0)
	//	for _, node := range row.Walk() { // todo bug
	//		columnCells := t.MarshalRowCells(node)
	//		for _, cell := range columnCells {
	//			mylog.Info(cell.Value, text)
	//			if strings.EqualFold(cell.Value, text) {
	//				children = append(children, node) // 过滤子节点
	//			}
	//		}
	//	}
	//	t.filteredRows[i].Children = children
	// }
	if len(t.filteredRows) == 0 {
		return
	}
	t.OpenAll()
	/*
		func (t *Table[T]) ApplyFilter(filter func(row T) bool) {
			if filter == nil {
				if t.filteredRows == nil {
					return
				}
				t.filteredRows = nil
			} else {
				t.filteredRows = make([]T, 0)
				for _, row := range t.Model.RootRows() {
					t.applyFilter(row, filter)
				}
			}
			if t.header != nil && t.header.HasSort() {
				t.header.ApplySort()
			}
		}

		func (t *Table[T]) applyFilter(row T, filter func(row T) bool) {
			if !filter(row) {
				t.filteredRows = append(t.filteredRows, row)
			}
			if row.CanHaveChildren() {
				for _, child := range row.Children() {
					t.applyFilter(child, filter)
				}
			}
		}
	*/
}

func (t *TreeTable[T]) OpenAll() {
	if t.filteredRows != nil {
		for _, row := range t.filteredRows {
			row.OpenAll()
		}
		return
	}
	t.Root.OpenAll()
}

func (t *TreeTable[T]) CloseAll() {
	if t.filteredRows != nil {
		for _, row := range t.filteredRows {
			row.CloseAll()
		}
		return
	}
	t.Root.CloseAll()
}

func (t *TreeTable[T]) Sort() {
	if len(t.rootRows) == 0 || t.header.sortedBy >= len(t.rootRows[0].rowCells) {
		return // 如果没有子节点或者列索引无效，直接返回
	}
	sort.Slice(t.rootRows, func(i, j int) bool {
		if t.rootRows[i].rowCells == nil {
			t.rootRows[i].rowCells = t.MarshalRowCells(t.rootRows[i])
		}
		if t.rootRows[j].rowCells == nil {
			t.rootRows[j].rowCells = t.MarshalRowCells(t.rootRows[j])
		}
		cellI := t.rootRows[i].rowCells[t.header.sortedBy].Value
		cellJ := t.rootRows[j].rowCells[t.header.sortedBy].Value
		if t.header.sortAscending {
			return cellI < cellJ
		}
		return cellI > cellJ
	})
}

type (
	sortOrder uint8
	// cellFn    func(gtx layout.Context, row, col int) layout.Dimensions
	rowFn func(gtx layout.Context, row int) layout.Dimensions
)

const (
	sortNone sortOrder = iota
	sortAscending
	sortDescending
)

const (
	defaultDividerWidth                   unit.Dp = 1
	defaultDividerMargin                  unit.Dp = 1
	defaultDividerHandleMinVerticalMargin unit.Dp = 2
	defaultDividerHandleMaxHeight         unit.Dp = 12
	defaultDividerHandleWidth             unit.Dp = 3
	defaultDividerHandleRadius            unit.Dp = 2
	defaultHeaderPadding                  unit.Dp = 5
	defaultHeaderBorder                   unit.Dp = 1
)

// ---------------------------------------泛型n叉树实现------------------------------------------

type Node[T any] struct {
	ID        uuid.ID          // 节点唯一标识符
	Type      string           // 容器节点类型，用于区分不同类型的容器节点
	Data      T                // 元数据
	Children  []*Node[T]       // 子节点,json只保存以上四个导出的字段
	parent    *Node[T]         // 父节点
	isOpen    bool             // 是否展开
	rowCells  []CellData       // 行单元格数据
	rowClick  widget.Clickable // 行点击按钮绘制,每个节点对应一个
	RowNumber int              // 展开所有节点状态下的下标，用于交替行背景色的条件
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

func (t *TreeTable[T]) AddChildByData(parent *Node[T], data T) {
	parent.AddChildByData(data)
}

func (t *TreeTable[T]) AddChildrenByData(parent *Node[T], data ...T) {
	parent.AddChildrenByData(data...)
}

func (t *TreeTable[T]) AddContainerByData(parent *Node[T], typeKey string, data T) {
	parent.AddContainerByData(typeKey, data)
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

func (n *Node[T]) GetType() string           { return n.Type }
func (n *Node[T]) SetType(typeKey string)    { n.Type = typeKey }
func (n *Node[T]) IsOpen() bool              { return n.isOpen && n.Container() }
func (n *Node[T]) SetOpen(open bool)         { n.isOpen = open && n.Container() }
func (n *Node[T]) Parent() *Node[T]          { return n.parent }
func (n *Node[T]) SetParent(parent *Node[T]) { n.parent = parent }

func (n *Node[T]) ResetChildren()        { n.Children = nil }
func (n *Node[T]) CanHaveChildren() bool { return n.HasChildren() }
func (n *Node[T]) HasChildren() bool     { return n.Container() && len(n.Children) > 0 }
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
	// for i, child := range n.parent.Children {
	//	if n.ID == child.ID {
	//		return i
	//	}
	// }
	// panic("not found index") // 永远不可能选中root，所以可以放心panic，root不显示，只显示它的children作为rootRows
}

func (n *Node[T]) Remove() {
	for i, child := range n.parent.Walk() {
		if child.ID == n.ID {
			n.parent.Children = slices.Delete(n.parent.Children, i, i+1)
			break
		}
	}
}

func (t *TreeTable[T]) InsertAfter(gtx layout.Context, after *Node[T]) {
	t.SelectedNode.InsertAfter(after)
}

func (t *TreeTable[T]) Remove(gtx layout.Context) {
	t.SelectedNode.Remove()
}

func (t *TreeTable[T]) Edit(gtx layout.Context) { // 编辑节点不会对最大深度有影响
	editor := NewStructView("edit row", t.SelectedNode.Data, func(key string, field any) (value string) {
		return "" // todo 这里需要实现编辑节点的功能
	})
	editor.Modal = true
	editor.SetOnApply(func() { // todo bug ,debug it
		// t.UnmarshalRowCells[T](t.SelectedNode, t.SelectedNode.rowCells)
		t.SelectedNode.Data = t.UnmarshalRowCells(t.SelectedNode, editor.Rows) // todo test
		mylog.Todo("save json data ?")
	})
	PopupCallback = editor
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
		g.WriteString(fmt.Sprintf("%-*s", int(widths[i]), strconv.Quote(cell.Value)) + ",")
	}
	g.P("}")
	gtx.Execute(clipboard.WriteCmd{Data: io.NopCloser(strings.NewReader(g.String()))})
	return g.String()
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

// func (n *Node[T]) WalkQueue() iter.Seq[*Node[T]] { // 性能应该不行
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
// }

func (t *TreeTable[T]) MaxDepth() unit.Dp {
	maxDepth := unit.Dp(0)
	for _, node := range t.Root.Walk() {
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
	to = deepcopy.Clone(n)
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

func (n *Node[T]) SetChildren(children []*Node[T]) {
	n.Children = children
}

// -------------------------tui
const (
	indent          = "│   "
	childPrefix     = "├───"
	lastChildPrefix = "└───"
	indentBase      = unit.Dp(3)
)

func (t *TreeTable[T]) maxColumnCellTextWidth() unit.Dp {
	hierarchyIndent := unit.Dp(1)
	dividerWidth := align.StringWidth[unit.Dp](" │ ")
	iconWidth := align.StringWidth[unit.Dp](childPrefix)
	return t.MaxDepth()*hierarchyIndent + // 最大深度的左缩进
		iconWidth + // 图标宽度,不管深度是多少，每一行都只会有一个层级图标
		t.maxColumnTextWidths[0] + 5 + // (8 * 2) + 20 + // 左右padding,20是sort图标的宽度或者容器节点求和的文本宽度
		dividerWidth // 列分隔条宽度
}

func (t *TreeTable[T]) Format() *stream.Buffer {
	buf := t.FormatHeader(t.maxColumnTextWidths)
	t.FormatChildren(buf, t.rootRows) // 传入子节点打印函数
	// mylog.Json("DrawRowCallback", buf.String())
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
	b := stream.NewBuffer("")
	all := t.maxColumnCellTextWidth()
	for _, width := range maxColumnCellTextWidths {
		all += width
	}
	all += align.StringWidth[unit.Dp]("│")*unit.Dp(len(maxColumnCellTextWidths)) + 4 // ?
	b.WriteStringLn("┌─" + strings.Repeat("─", int(all)))
	b.WriteString("│")

	// 计算每个单元格的左边距
	for i, cell := range t.header.columnCells {
		paddedText := fmt.Sprintf("%-*s", int(maxColumnCellTextWidths[i]), cell.Value) // 左对齐填充

		// 添加左边距，仅在首列进行处理，依据列宽计算
		if i == HierarchyColumnID {
			b.WriteString(strings.Repeat(" ", int(t.maxColumnCellTextWidth()-maxColumnCellTextWidths[i]-1))) // -1是分隔符的空间
		}

		b.WriteString(paddedText)
		if i < t.columnCount-1 {
			b.WriteString(" │ ") // 在每个单元格之间添加分隔符
		}
	}

	b.NewLine()
	b.WriteStringLn("├─" + strings.Repeat("─", int(all)))
	return b
}

func (t *TreeTable[T]) FormatChildren(b *stream.Buffer, children []*Node[T]) {
	for i, child := range children {
		if child.rowCells == nil {
			child.rowCells = t.MarshalRowCells(child)
		}
		HierarchyColumBuf := stream.NewBuffer("")
		for j, cell := range child.rowCells {
			if j == HierarchyColumnID {
				HierarchyColumBuf.WriteString("│")
				HierarchyColumBuf.WriteString(strings.Repeat(" ", int(child.Depth()*indentBase)))
				if i == len(children)-1 {
					HierarchyColumBuf.WriteString("╰──") // "└───"
				} else {
					HierarchyColumBuf.WriteString("├──")
				}
				HierarchyColumBuf.WriteString(cell.Value)
				if align.StringWidth[unit.Dp](HierarchyColumBuf.String()) < t.maxColumnCellTextWidth() {
					HierarchyColumBuf.WriteString(strings.Repeat(" ", int(t.maxColumnCellTextWidth()-align.StringWidth[unit.Dp](HierarchyColumBuf.String()))))
				}
				HierarchyColumBuf.WriteString(" │ ")
				b.WriteString(HierarchyColumBuf.String())
				HierarchyColumBuf.Reset()
				continue
			}
			b.WriteString(cell.Value)
			if align.StringWidth[unit.Dp](cell.Value) < t.maxColumnTextWidths[j] {
				b.WriteString(strings.Repeat(" ", int(t.maxColumnTextWidths[j]-align.StringWidth[unit.Dp](cell.Value))))
			}
			b.WriteString(" │ ")
		}
		b.NewLine()
		if len(child.Children) > 0 {
			t.FormatChildren(b, child.Children)
		}
	}
}

// -------------------------todo

// type TableDragData[T any] struct {
// 	Table *Node[T]
// 	Rows  []*Node[T]
// }

// var DefaultTableTheme = TableTheme{
// 	HierarchyIndent:   16,
// 	MinimumRowHeight:  16,
// 	ColumnResizeSlop:  4,
// 	ShowRowDivider:    true,
// 	ShowColumnDivider: true,
// }
//
// type TableTheme struct {
// 	Padding           layout.Inset
// 	HierarchyColumnID int
// 	HierarchyIndent   unit.Dp
// 	MinimumRowHeight  unit.Dp
// 	ColumnResizeSlop  unit.Dp
// 	ShowRowDivider    bool
// 	ShowColumnDivider bool
// }

// type Point struct {
// 	X, Y unit.Dp
// }

// func (t *TreeTable[T]) layoutDrag(gtx layout.Context, w rowFn) layout.Dimensions {
// 	columns := t.maxColumnCellWidths
// 	var (
// 		cols          = t.columnCount         // 获取列的数量
// 		dividers      = cols                  // 列的分隔数
// 		tallestHeight = gtx.Constraints.Min.Y // 初始化最高的行高
//
// 		dividerWidth                   = gtx.Dp(defaultDividerWidth)                   // 获取分割线的宽度
// 		dividerMargin                  = gtx.Dp(defaultDividerMargin)                  // 获取分割线的边距
// 		dividerHandleMinVerticalMargin = gtx.Dp(defaultDividerHandleMinVerticalMargin) // 获取分隔处理器的最小垂直边距
// 		dividerHandleMaxHeight         = gtx.Dp(defaultDividerHandleMaxHeight)         // 获取分隔处理器的最大高度
// 		dividerHandleWidth             = gtx.Dp(defaultDividerHandleWidth)             // 获取分隔处理器的宽度
// 		dividerHandleRadius            = gtx.Dp(defaultDividerHandleRadius)            // 获取分隔处理器的圆角半径
//
// 		minWidth = unit.Dp(dividerWidth + dividerMargin + dividerHandleWidth) // 计算分隔符的最小宽度
// 	)
// 	if cols == 0 { // 如果列数为0，直接返回最小约束大小
// 		return layout.Dimensions{Size: gtx.Constraints.Min}
// 	}
//
// 	if len(t.header.drags) < dividers { // 如果拖动数组没有足够的长度
// 		t.header.drags = make([]tableDrag, dividers) // 初始化拖动对象数组
// 	}
//
// 	// OPT(dh): we don't need to do this for each t, only once per table
// 	for i := range t.header.drags { // 遍历每个拖动对象
// 		drag := &t.header.drags[i]    // 获取当前拖动对象
// 		width := columns[i]           // 获取当前列宽度
// 		drag.hover.Update(gtx.Source) // 更新当前列的悬停状态
// 		// OPT(dh): Events allocates
// 		var delta unit.Dp // 初始化偏移量
//
// 		for { // 循环处理拖动事件
// 			ev, ok := drag.drag.Update(gtx.Metric, gtx.Source, gesture.Horizontal) // 更新拖动状态
// 			if !ok {                                                               // 如果事件不合法，退出循环
// 				break
// 			}
// 			switch ev.Kind { // 根据事件类型进行处理
// 			case pointer.Press: // 按下事件
// 				drag.startPos = ev.Position.X                             // 记录拖动的起始位置
// 				drag.shrinkNeighbor = !ev.Modifiers.Contain(key.ModShift) // 判断是否收缩相邻列
// 			case pointer.Drag: // 拖动事件
// 				t.header.manualWidthSet[i] = true // 在这里标记列宽已手动调整
// 				// There may be multiple drag events in a single frame. We mustn't apply all of them or we'll
// 				// drag too far. Only react to the final event.
// 				delta = unit.Dp(ev.Position.X - drag.startPos) // 计算当前拖动的偏移量
// 			}
// 		}
// 		if delta != 0 { // 如果存在拖动偏移量
// 			width += delta                                  // 更新列的宽度
// 			if drag.shrinkNeighbor && i != len(columns)-1 { // 如果需要收缩相邻列且不是最后一列
// 				nextCol := columns[i+1] // 获取下一个列
// 				nextCol -= delta        // 更新下一个列的宽度
// 				if width < minWidth {   // 如果当前列宽度小于最小宽度
// 					d := minWidth - width // 计算需要增加的宽度
// 					width = minWidth      // 将当前列宽度设为最小宽度
// 					nextCol -= d          // 更新下一个列的宽度
// 				}
// 				if nextCol < minWidth { // 如果下一个列宽度小于最小宽度
// 					d := minWidth - nextCol // 计算需要增加的宽度
// 					nextCol = minWidth      // 将下一个列宽度设为最小宽度
// 					width -= d              // 更新当前列宽度
// 				}
// 			} else {
// 				// 如果不需要收缩
// 				if width < minWidth { // 如果当前列宽度小于最小宽度
// 					width = minWidth // 将当前列宽度设为最小宽度
// 				}
// 			}
//
// 			// if width < width.autoWidth { // 如果当前列宽度小于其最小宽度
// 			//	width = width.autoWidth // 更新列的最小宽度为当前宽度
// 			// }
// 			// width = max(width, minWidth) // 确保列的宽度不小于最小宽度
//
// 			var total unit.Dp             // 初始化总宽度
// 			for _, col := range columns { // 遍历所有列计算总宽度
// 				total += col // 累加当前列的宽度
// 			}
// 			total += unit.Dp(len(columns) * gtx.Dp(defaultDividerWidth)) // 加上所有分隔符的总宽度
// 			if total < unit.Dp(gtx.Constraints.Min.X) {                  // 如果总宽度小于最小约束宽度
// 				columns[len(columns)-1] += unit.Dp(gtx.Constraints.Min.X) - total // 调整最后一列的宽度以适应
// 			}
// 		}
// 	}
//
// 	for { // 开始绘制列
// 		// First draw all columns, leaving gaps for the drag handlers
// 		var (
// 			start             = 0             // 初始化当前位置
// 			origTallestHeight = tallestHeight // 记录最初的高度
// 		)
// 		r := op.Record(gtx.Ops)  // 记录当前操作集合
// 		totalWidth := 0          // 初始化总宽度
// 		for i := range columns { // 遍历所有列
// 			colWidth := int(columns[i]) // 获取当前列的宽度
// 			totalWidth += colWidth      // 更新总宽度
// 		}
// 		extra := gtx.Constraints.Min.X - len(columns)*gtx.Dp(defaultDividerWidth) - totalWidth // 计算多余宽度
// 		colExtra := extra                                                                      // 将多余宽度赋值给列额外宽度
//
// 		for i := range columns { // 绘制所有列
// 			colWidth := int(columns[i]) // 获取当前列宽度
// 			if colExtra > 0 {           // 如果有多余宽度
// 				colWidth++ // 当前列宽度加一
// 				colExtra-- // 多余宽度减一
// 			}
//
// 			gtx := gtx                            // 更新 gtx 上下文
// 			gtx.Constraints.Min.X = colWidth      // 设置当前列的最小宽度
// 			gtx.Constraints.Max.X = colWidth      // 设置当前列的最大宽度
// 			gtx.Constraints.Min.Y = tallestHeight // 设置当前列的最小高度
//
// 			stack := op.Offset(image.Pt(start, 0)).Push(gtx.Ops) // 设置当前列绘制的偏移量
//
// 			dims := w(gtx, i)                                // 绘制当前列，返回所需尺寸
// 			dims.Size = gtx.Constraints.Constrain(dims.Size) // 应用约束限制
// 			tallestHeight = dims.Size.Y                      // 更新当前行的高度
// 			if i == 0 && tallestHeight > origTallestHeight { // 如果当前行的高度大于初始高度
// 				origTallestHeight = tallestHeight // 更新初始高度
// 			}
//
// 			start += colWidth + dividerWidth // 更新绘制起始位置
// 			stack.Pop()                      // 弹出当前堆栈
// 		}
// 		call := r.Stop() // 停止记录操作集合
//
// 		if tallestHeight > origTallestHeight { // 如果当前高度大于最初高度
// 			continue // 重新绘制当前行
// 		}
//
// 		call.Add(gtx.Ops) // 将操作添加到上下文中
//
// 		// Then draw the drag handlers. The handlers overdraw the columns when hovered.
// 		var (
// 			dividerHandleHeight    = min(tallestHeight-2*dividerHandleMinVerticalMargin, dividerHandleMaxHeight) // 获取分隔处理器的高度
// 			dividerHandleTopMargin = (tallestHeight - dividerHandleHeight) / 2                                   // 计算顶部边距
// 			dividerStart           = 0                                                                           // 初始化分隔起始位置
// 			dividerExtra           = extra                                                                       // 设置额外宽度
// 		)
// 		for i := range t.header.drags { // 遍历每个拖动对象
// 			var (
// 				drag     = &t.header.drags[i] // 获取当前拖动对象
// 				colWidth = int(columns[i])    // 获取当前列宽度
// 			)
// 			dividerStart += colWidth // 更新分隔符的起始位置
// 			if dividerExtra > 0 {    // 如果还有多余的宽度
// 				dividerStart++ // 分隔符起始位置加一
// 				dividerExtra-- // 多余宽度减一
// 			}
//
// 			// We add the drag handler slightly outside the drawn divider, to make it easier to press.
// 			//
// 			// We use op.Offset instead of folding dividerStart into the clip.Rect because we want to set the
// 			// origin of the drag coordinates.
// 			stack := op.Offset(image.Pt(dividerStart, 0)).Push(gtx.Ops) // set origin for drag coordinates
// 			stack2 := clip.Rect{
// 				Min: image.Pt(-dividerMargin-dividerHandleWidth, 0),                         // 设置分隔符的最小边界
// 				Max: image.Pt(dividerWidth+dividerMargin+dividerHandleWidth, tallestHeight), // 设置分隔符的最大边界
// 			}.Push(gtx.Ops) // 在上下文中推入分隔符的裁剪矩形
//
// 			if t.inLayoutHeader { // 如果当前行是表头 todo test
// 				drag.hover.Update(gtx.Source)        // 更新悬停状态
// 				drag.drag.Add(gtx.Ops)               // 添加拖动操作到上下文中
// 				drag.hover.Add(gtx.Ops)              // 添加悬停操作到上下文中
// 				pointer.CursorColResize.Add(gtx.Ops) // 设置拖动光标样式为调整列宽
//
// 				// Draw the left and right extensions when hovered.
// 				if drag.hover.Update(gtx.Source) || drag.drag.Dragging() { // 如果悬停或者拖动
// 					handleShape := clip.UniformRRect(
// 						image.Rect(
// 							0,
// 							dividerHandleTopMargin,
// 							dividerHandleWidth,
// 							dividerHandleTopMargin+dividerHandleHeight),
// 						dividerHandleRadius,
// 					) // 设定分隔处理器的形状
// 					handleLeft := handleShape
// 					handleLeft.Rect = handleShape.Rect.Add(image.Pt(-(dividerMargin + dividerHandleWidth), 0)) // 为左边形状添加偏移
// 					handleRight := handleShape
// 					handleRight.Rect = handleRight.Rect.Add(image.Pt(dividerWidth+dividerMargin, 0)) // 为右边形状添加偏移
//
// 					paint.FillShape(gtx.Ops, colors.Red200, handleLeft.Op(gtx.Ops))     // 填充左侧形状
// 					paint.FillShape(gtx.Ops, colors.Yellow100, handleRight.Op(gtx.Ops)) // 填充右侧形状
// 				}
//
// 				// Draw the vertical bar
// 				// stack3 := clip.Rect{Max: image.Pt(dividerWidth, tallestHeight)}.Push(gtx.Ops)
// 				// Fill( gtx.Ops, win.Theme.Palette.Table2.Divider) // 如果有需要，在此处可绘制分割线
// 				// stack3.Pop()
// 			}
// 			// 为表头和每列绘制列分隔条
// 			stack3 := clip.Rect{Max: image.Pt(dividerWidth, tallestHeight)}.Push(gtx.Ops) // 绘制分隔条的矩形区域
// 			paint.Fill(gtx.Ops, colors.DividerFg)                                         // 填充分隔条的颜色
// 			stack3.Pop()                                                                  // 弹出分隔条的绘制堆栈
//
// 			dividerStart += dividerWidth // 更新分隔符的起始位置
// 			stack2.Pop()                 // 弹出分隔符的绘制堆栈
// 			stack.Pop()                  // 弹出当前列的绘制堆栈
// 		}
//
// 		return layout.Dimensions{
// 			Size: image.Pt(start, tallestHeight), // 返回绘制完成后的整体宽度和高度
// 		}
// 	}
// }

// v0.4
