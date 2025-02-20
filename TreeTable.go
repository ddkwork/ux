package ux

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"io"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"gioui.org/gesture"
	"gioui.org/io/clipboard"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"github.com/ddkwork/golibrary/stream/align"
	"github.com/google/uuid"
)

type (
	TreeTable[T any] struct {
		Root         *Node[T] //？ how to use it?
		Children     []*Node[T]
		selectedNode *Node[T]
		filterText   string
		filteredRows []*Node[T]

		header            TableHeader[T]
		columnHeaderCount int
		inLayoutHeader    bool // for drag
		TableContext[T]
		widget.List

		maxColumnTextWidths []unit.Dp
		maxColumnTexts      []string

		Rows               [][]CellData
		Columns            [][]CellData
		maxColumnCellWidth unit.Dp
	}
	TableContext[T any] struct {
		ContextMenuItems       func(node *Node[T], gtx layout.Context) (items []ContextMenuItem)
		MarshalRow             func(node *Node[T]) (cells []CellData)
		UnmarshalRow           func(node *Node[T], values []string)
		RowSelectedCallback    func(node *Node[T]) // 行选中回调
		RowDoubleClickCallback func(node *Node[T]) // double click callback
		LongPressCallback      func(node *Node[T]) // mobile long press callback
		SetRootRowsCallBack    func(node *Node[T])
		JsonName               string
		IsDocument             bool
	}
	TableHeader[T any] struct {
		SortOrder          SortOrder
		SortedBy           int
		drags              []tableDrag
		ColumnCells        []CellData
		clickedColumnIndex int
		manualWidthSet     []bool // 新增状态标志数组，记录列是否被手动调整
		sortAscending      bool
		contextAreas       []*component.ContextArea
		contextMenu        *ContextMenu
	}
	CellData struct {
		ColumID int // 列id
		RowID   int // 行id

		Text     string  // 单元格文本
		maxDepth unit.Dp // 层级

		maxColumnTextWidth unit.Dp
		maxColumnText      string

		Current     unit.Dp // 正在使用的宽度
		Minimum     unit.Dp // 拖放表头列分隔条得到的最小宽度
		Maximum     unit.Dp // 拖放表头列分隔条得到的最大宽度
		AutoMinimum unit.Dp // 根据单元格内容预渲染自动计算最小宽度
		AutoMaximum unit.Dp // 根据单元格内容预渲染自动计算最大宽度
		Disabled    bool    // 是否显示表头或者body单元格
		Tooltip     string
		SvgBuffer   string
		ImageBuffer []byte
		FgColor     color.NRGBA
		IsNasm      bool
		IsHeader    bool
		widget.Clickable
		RichText
		leftIndent unit.Dp
	}
)

const ContainerKeyPostfix = "_container"

func NewTreeTable[T any](data T, ctx TableContext[T]) *TreeTable[T] {
	columnCells := initHeader(data)
	root := NewRoot(data)
	root.MarshalRow = ctx.MarshalRow
	return &TreeTable[T]{
		Root:         root,
		Children:     nil,
		selectedNode: nil,
		filterText:   "",
		filteredRows: nil,
		header: TableHeader[T]{
			SortOrder:          0,
			SortedBy:           0,
			drags:              make([]tableDrag, 0),
			ColumnCells:        columnCells,
			clickedColumnIndex: -1,
			manualWidthSet:     make([]bool, len(columnCells)),
		},
		columnHeaderCount: len(columnCells),
		inLayoutHeader:    false,
		TableContext:      ctx,
		List: widget.List{
			Scrollbar: widget.Scrollbar{},
			List: layout.List{
				Axis:        layout.Vertical,
				ScrollToEnd: false,
				Alignment:   0,
				Position:    layout.Position{},
			},
		},
		maxColumnTextWidths: nil,
		maxColumnTexts:      nil,
		Rows:                nil,
		Columns:             nil,
	}
}

func NewRoot[T any](data T) *Node[T] {
	return NewContainerNode("root", data)
}

func (n *Node[T]) IsRoot() bool {
	if n == nil { // todo bug
		return false
	}
	return n.parent == nil
}

func newNode[T any](typeKey string, isContainer bool, data T) *Node[T] {
	if isContainer {
		typeKey += ContainerKeyPostfix
	}
	n := &Node[T]{
		MarshalRow:               nil,
		MarshalColum:             nil,
		Icon:                     widget.Icon{},
		rowSelected:              false,
		rowClick:                 widget.Clickable{},
		CellClickedCallback:      nil,
		rowContextAreas:          nil,
		contextMenu:              nil,
		TableTheme:               DefaultTableTheme,
		ID:                       NewUUID(),
		Type:                     typeKey,
		parent:                   nil,
		Data:                     data,
		Children:                 nil,
		isOpen:                   isContainer,
		SelectionChangedCallback: nil,
		DoubleClickCallback:      nil,
		DragRemovedRowsCallback:  nil,
		DropOccurredCallback:     nil,
		filteredRows:             nil,
		columnResizeStart:        0,
		columnResizeBase:         0,
		columnResizeOverhead:     0,
		PreventUserColumnResize:  false,
		awaitingSizeColumnsToFit: false,
		awaitingSyncToModel:      false,
		wasDragged:               false,
		dividerDrag:              false,
		RowCells:                 nil,
	}
	n.wasDragged = false
	return n
}

type Node[T any] struct {
	// TableContext[T]
	MarshalRow   func(node *Node[T]) (cells []CellData) `json:"-"`
	MarshalColum func(node *Node[T]) (cells []CellData) `json:"-"`

	Icon                widget.Icon // 层级列图标，不是其他列的图标，其他列要看富文本是否支持渲染图标或者封装单元格渲染函数
	rowSelected         bool
	rowClick            widget.Clickable
	CellClickedCallback func(root *Node[T])

	rowContextAreas []*component.ContextArea
	contextMenu     *ContextMenu

	TableTheme `json:"-"`

	ID       uuid.UUID `json:"id"`
	Type     string    `json:"type"`
	parent   *Node[T]
	Data     T
	Children []*Node[T] `json:"children,omitempty"`
	isOpen   bool

	SelectionChangedCallback func()     `json:"-"`
	DoubleClickCallback      func()     `json:"-"`
	DragRemovedRowsCallback  func()     `json:"-"` // Called whenever a drag removes one or more rows from a model, but only if the source and destination tables were different.
	DropOccurredCallback     func()     `json:"-"` // Called whenever a drop occurs that modifies the model.
	filteredRows             []*Node[T] // todo move into treeTable

	columnResizeStart        unit.Dp
	columnResizeBase         unit.Dp
	columnResizeOverhead     unit.Dp
	PreventUserColumnResize  bool
	awaitingSizeColumnsToFit bool
	awaitingSyncToModel      bool
	wasDragged               bool
	dividerDrag              bool
	RowCells                 []CellData

	LongPressCallback func(node *Node[T]) // 长按回调
	pressStarted      time.Time           // 按压开始时间
	longPressed       bool                // 是否已经触发长按事件
}

func (n *Node[T]) UpdateTouch(gtx layout.Context) {
	// 检测触摸事件
	//for _, ev := range gtx.Events(n) {
	//	if e, ok := ev.(pointer.Event); ok {
	//		switch e.Type {
	//		case pointer.Press:
	//			n.pressStarted = time.Now() // 记录按压开始时间
	//			n.longPressed = false       // 重置长按状态
	//		case pointer.Release:
	//			if n.longPressed {
	//				// 如果已经触发了长按事件，不需要额外处理
	//				return
	//			}
	//			// 检查是否是点击事件
	//			if time.Since(n.pressStarted) < LongPressDuration {
	//				// 处理点击事件
	//				if n.rowClick.OnClicked() {
	//					n.isOpen = !n.isOpen
	//					if n.CellClickedCallback != nil {
	//						n.CellClickedCallback(n)
	//					}
	//					if t.RowSelectedCallback != nil {
	//						t.RowSelectedCallback(n)
	//					}
	//				}
	//			}
	//		}
	//	}
	//}

	// 检测长按事件
	if gtx.Now.Sub(n.pressStarted) > LongPressDuration && !n.longPressed {
		n.longPressed = true
		if n.LongPressCallback != nil {
			n.LongPressCallback(n)
		}
	}
}

// ----------------------------------------------------------

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

const LongPressDuration = 500 * time.Millisecond // 自定义长按持续时间

func (t *TreeTable[T]) SetRootRows(rootRows []*Node[T]) *TreeTable[T] {
	for _, row := range rootRows {
		t.expandNode(row)
		// 设置长按回调
		row.LongPressCallback = func(node *Node[T]) {
			// 长按时执行的操作
			t.header.clickedColumnIndex = -1 // 重置点击列索引（如果需要）
			t.selectedNode = node            // 设置选中节点
			// 显示上下文菜单
			//t.header.contextMenu.Show(gtx, func(gtx layout.Context) layout.Dimensions {
			//	return t.drawContextArea(gtx, &t.header.contextMenu.MenuState)
			//})
		}
		// row.UpdateTouch(gtx) // 初始化触摸事件处理
	}
	t.Children = rootRows
	return t
}

func (n *Node[T]) SetParents(children []*Node[T], parent *Node[T]) {
	for _, child := range children {
		child.parent = parent
		if len(child.Children) > 0 {
			n.SetParents(child.Children, child)
		}
	}
}

// 递归函数，获取树形结构中的最大深度
func (t *TreeTable[T]) MaxDepth() unit.Dp {
	maxDepth := unit.Dp(1)
	t.Root.Walk(func(node *Node[T]) {
		childDepth := node.Depth()
		if childDepth > maxDepth {
			maxDepth = childDepth
		}
	})
	return maxDepth
}

func (t *TreeTable[T]) SizeColumnsToFit(gtx layout.Context, isTui bool) {
	originalConstraints := gtx.Constraints // 保存原始约束
	fmtRow := func(row []CellData, id int) {
		return
		b := stream.NewBuffer("")
		for _, data := range row {
			b.WriteString(fmt.Sprintf("%-20s", data.Text))
		}
		// mylog.Warning("rowCells "+fmt.Sprint(id), b.String())
	}
	fmtColumn := func(column []CellData) {
		return
		b := stream.NewBuffer("")
		for _, data := range column {
			b.WriteStringLn(data.Text)
		}
		// mylog.Json("-------> col: "+fmt.Sprint(column[0].ColumID), b.String())
	}

	t.Rows = make([][]CellData, 0, len(t.Children)) // 用于存储所有行
	t.Root.Walk(func(node *Node[T]) {
		rowCells := t.MarshalRow(node)
		t.Rows = append(t.Rows, rowCells)
	})
	t.Rows = slices.Insert(t.Rows, 0, t.header.ColumnCells) // 插入表头行
	for row, rowCells := range t.Rows {
		fmtRow(rowCells, row)
	}

	t.Columns = TransposeMatrix(t.Rows)

	t.maxColumnTextWidths = make([]unit.Dp, len(t.header.ColumnCells))
	t.maxColumnTexts = make([]string, len(t.header.ColumnCells))
	for i, column := range t.Columns {
		fmtColumn(column)
		if t.header.manualWidthSet[i] { // 如果该列已手动调整
			continue // 跳过，保留用户手动调整的宽度
		}
		t.maxColumnTextWidths[i] = 0
		t.maxColumnTexts[i] = ""
		for _, data := range column {
			if len(data.Text) > len(t.maxColumnTexts[i]) {
				t.maxColumnTexts[i] = data.Text
			}
			if isTui {
				t.maxColumnTextWidths[i] = max(t.maxColumnTextWidths[i], align.StringWidth[unit.Dp](data.Text))
			} else {
				t.maxColumnTextWidths[i] = max(t.maxColumnTextWidths[i], CalculateTextWidth(gtx, data.Text))
			}
		}
	}

	maxDepth := t.MaxDepth()
	for i, maxWidth := range t.maxColumnTextWidths {
		t.header.ColumnCells[i].Minimum = maxWidth
		t.header.ColumnCells[i].Current = maxWidth
		t.header.ColumnCells[i].maxDepth = maxDepth
	}
	gtx.Constraints = originalConstraints

	t.Children = t.Root.Children
}

// TransposeMatrix 把行切片矩阵置换为列切片,用于计算最大列宽的参数
func TransposeMatrix[T any](rows [][]T) (columns [][]T) {
	if len(rows) == 0 {
		return [][]T{}
	}
	columns = make([][]T, len(rows[0]))
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

func (t *TreeTable[T]) Layout(gtx layout.Context) layout.Dimensions {
	t.SizeColumnsToFit(gtx, false)
	list := material.List(th.Theme, &t.List)
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.HeaderFrame(gtx) // 渲染表头
			t.inLayoutHeader = true
			return t.layoutDrag(gtx, func(gtx layout.Context, row int) layout.Dimensions {
				return t.HeaderFrame(gtx) // 渲染表头
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return list.Layout(gtx, len(t.Children), func(gtx layout.Context, index int) layout.Dimensions {
				node := t.Children[index]
				node.UpdateTouch(gtx) // 更新触摸事件处理逻辑
				return t.RowFrame(gtx, node, index)
				return t.RowFrame(gtx, t.Children[index], index)
				//t.inLayoutHeader = false
				//return t.layoutDrag(gtx, func(gtx layout.Context, row int) layout.Dimensions {
				//	return t.RowFrame(gtx, t.Children[index], index)
				//})
			})
		}),
	)
}

func initHeader(data any) (Columns []CellData) {
	fields := stream.ReflectVisibleFields(data)
	Columns = make([]CellData, 0)
	for i, field := range fields {
		if field.Tag.Get("table") != "" { // 中文表头简短
			field.Name = field.Tag.Get("table")
		}
		Columns = append(Columns, CellData{
			ColumID:     i,
			Text:        field.Name,
			Current:     20,
			Minimum:     20,
			Maximum:     10000,
			AutoMinimum: 0,
			AutoMaximum: 0,
			Disabled:    false,
			Tooltip:     "",
			SvgBuffer:   "",
			ImageBuffer: nil,
			FgColor:     color.NRGBA{},
			IsNasm:      false,
			Clickable:   widget.Clickable{},
			RichText:    RichText{},
		})
	}
	return
}

func (t *TreeTable[T]) HeaderFrame(gtx layout.Context) layout.Dimensions {
	var cols []layout.FlexChild
	elems := make([]*Resizable, 0)
	for i, cell := range t.header.ColumnCells {
		if cell.Disabled {
			continue
		}
		clickable := &t.header.ColumnCells[i].Clickable
		if clickable.Clicked(gtx) {
			t.header.clickedColumnIndex = i
			if t.header.SortedBy == t.header.clickedColumnIndex {
				t.header.sortAscending = !t.header.sortAscending // 切换升序/降序

				switch t.header.SortOrder {
				case SortNone:
					t.header.SortOrder = SortAscending
				case SortAscending:
					t.header.SortOrder = SortDescending
				case SortDescending:
					t.header.SortOrder = SortAscending
				}

			} else {
				t.header.SortedBy = t.header.clickedColumnIndex // 更新排序列
				t.header.sortAscending = true                   // 设为升序
				t.header.SortOrder = SortAscending
			}

			mylog.Info("clickedColumnIndex", t.header.clickedColumnIndex)
			if t.header.clickedColumnIndex > -1 && t.header.SortedBy > -1 {
				for i := range t.header.ColumnCells {
					t.header.ColumnCells[i].Text = strings.TrimSuffix(t.header.ColumnCells[i].Text, " ⇩")
					t.header.ColumnCells[i].Text = strings.TrimSuffix(t.header.ColumnCells[i].Text, " ⇧")
				}
				switch t.header.SortOrder {
				case SortNone:
				case SortAscending:
					t.header.ColumnCells[i].Text += " ⇧"
				case SortDescending:
					t.header.ColumnCells[i].Text += " ⇩"
				}
				t.SortNodes()
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
				switch e.Buttons {
				case pointer.ButtonPrimary:
					/*
						var LongPressDuration time.Duration = 250 * time.Millisecond

								case gesture.KindPress:
							i.pressStarted = gtx.Now

							if !i.longPressed && i.pressing && gtx.Now.Sub(i.pressStarted) > LongPressDuration {
								i.longPressed = true
								return Event{Type: LongPress}, true
							}

						所以合理的方案是patch官方的contextAreas和gtx的input source代码，支持长按事件
					*/
					//todo模拟长按事件并给上下文区域增加长按事件来支持移动平台，body也要模拟，以及代码编辑器也要支持右键和长按弹出复制选中文本到剪切板的功能
				case pointer.ButtonSecondary:
					if e.Kind == pointer.Press {
						t.header.clickedColumnIndex = i
					}
				}
			}
		}

		cols = append(cols, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return Background{Color: ColorHeaderFg}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// 实现右键复制列到剪切板功能
				return layout.Stack{}.Layout(gtx,
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return material.Clickable(gtx, clickable, func(gtx layout.Context) layout.Dimensions {
							t.header.ColumnCells[i].IsHeader = true
							t.header.ColumnCells[0].Minimum += calculateMaxColumnCellWidth(t.header.ColumnCells[0])
							cellFrame := t.Root.CellFrame(gtx, t.header.ColumnCells[i])
							elems = append(elems, &Resizable{Widget: func(gtx layout.Context) layout.Dimensions {
								return cellFrame
							}})
							return cellFrame
						})
					}),
					layout.Expanded(func(gtx layout.Context) layout.Dimensions {
						if t.header.contextAreas == nil {
							t.header.contextAreas = make([]*component.ContextArea, len(t.header.ColumnCells))
						}
						contextArea := t.header.contextAreas[i]
						if contextArea == nil {
							contextArea = &component.ContextArea{
								/*
									var LongPressDuration time.Duration = 250 * time.Millisecond

											case gesture.KindPress:
										i.pressStarted = gtx.Now

										if !i.longPressed && i.pressing && gtx.Now.Sub(i.pressStarted) > LongPressDuration {
											i.longPressed = true
											return Event{Type: LongPress}, true
										}

									所以合理的方案是patch官方的contextAreas和gtx的input source代码，支持长按事件
								*/
								Activation:       pointer.ButtonSecondary,
								AbsolutePosition: true,
								PositionHint:     0,
							}
							t.header.contextAreas[i] = contextArea
						}
						if t.header.contextMenu == nil {
							t.header.contextMenu = NewContextMenu()
							t.header.contextMenu.AddItem(ContextMenuItem{
								Title:         "CopyColumn",
								Icon:          IconCopy,
								Can:           func() bool { return true },
								Do:            func() { t.CopyColumn(gtx) },
								AppendDivider: true,
								Clickable:     widget.Clickable{},
							})
						}
						t.header.contextMenu.OnClicked(gtx)
						return contextArea.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return t.drawContextArea(gtx, &t.header.contextMenu.MenuState)
						})
					}),
				)
			})
		}))
	}
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, WeightSum: 0}.Layout(gtx, cols...)

	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, WeightSum: 0}.Layout(gtx, cols...)
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			resizeWidget = NewResizeWidget(layout.Horizontal, elems...)
			return resizeWidget.Layout(gtx)
		}),
	)
}

var resizeWidget *Resize

func (t *TreeTable[T]) layoutDrag(gtx layout.Context, w RowFn) layout.Dimensions {
	var (
		cols          = len(t.Columns)        // 获取列的数量
		dividers      = cols                  // 列的分隔数
		tallestHeight = gtx.Constraints.Min.Y // 初始化最高的行高

		dividerWidth                   = gtx.Dp(DefaultDividerWidth)                   // 获取分割线的宽度
		dividerMargin                  = gtx.Dp(DefaultDividerMargin)                  // 获取分割线的边距
		dividerHandleMinVerticalMargin = gtx.Dp(DefaultDividerHandleMinVerticalMargin) // 获取分隔处理器的最小垂直边距
		dividerHandleMaxHeight         = gtx.Dp(DefaultDividerHandleMaxHeight)         // 获取分隔处理器的最大高度
		dividerHandleWidth             = gtx.Dp(DefaultDividerHandleWidth)             // 获取分隔处理器的宽度
		dividerHandleRadius            = gtx.Dp(DefaultDividerHandleRadius)            // 获取分隔处理器的圆角半径

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
		col := &t.Columns[i][0]       // 获取当前列
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
			col.Current += delta                              // 更新列的宽度
			if drag.shrinkNeighbor && i != len(t.Columns)-1 { // 如果需要收缩相邻列且不是最后一列
				nextCol := &t.Columns[i+1][0] // 获取下一个列
				nextCol.Current -= delta      // 更新下一个列的宽度
				if col.Current < minWidth {   // 如果当前列宽度小于最小宽度
					d := minWidth - col.Current // 计算需要增加的宽度
					col.Current = minWidth      // 将当前列宽度设为最小宽度
					nextCol.Current -= d        // 更新下一个列的宽度
				}
				if nextCol.Current < minWidth { // 如果下一个列宽度小于最小宽度
					d := minWidth - nextCol.Current // 计算需要增加的宽度
					nextCol.Current = minWidth      // 将下一个列宽度设为最小宽度
					col.Current -= d                // 更新当前列宽度
				}
			} else {
				// 如果不需要收缩
				if col.Current < minWidth { // 如果当前列宽度小于最小宽度
					col.Current = minWidth // 将当前列宽度设为最小宽度
				}
			}

			if col.Current < col.Minimum { // 如果当前列宽度小于其最小宽度
				col.Minimum = col.Current // 更新列的最小宽度为当前宽度
			}

			var total unit.Dp               // 初始化总宽度
			for _, col := range t.Columns { // 遍历所有列计算总宽度
				total += col[0].Current // 累加当前列的宽度
			}
			total += unit.Dp(len(t.Columns) * gtx.Dp(DefaultDividerWidth)) // 加上所有分隔符的总宽度
			if total < unit.Dp(gtx.Constraints.Min.X) {                    // 如果总宽度小于最小约束宽度
				t.Columns[len(t.Columns)-1][0].Current += unit.Dp(gtx.Constraints.Min.X) - total // 调整最后一列的宽度以适应
			}
		}
	}

	for { // 开始绘制列
		// First draw all columns, leaving gaps for the drag handlers
		var (
			start             = 0             // 初始化当前位置
			origTallestHeight = tallestHeight // 记录最初的高度
		)
		r := op.Record(gtx.Ops)    // 记录当前操作集合
		totalWidth := 0            // 初始化总宽度
		for i := range t.Columns { // 遍历所有列
			colWidth := int(t.Columns[i][0].Current) // 获取当前列的宽度
			totalWidth += colWidth                   // 更新总宽度
		}
		extra := gtx.Constraints.Min.X - len(t.Columns)*gtx.Dp(DefaultDividerWidth) - totalWidth // 计算多余宽度
		colExtra := extra                                                                        // 将多余宽度赋值给列额外宽度

		for i := range t.Columns { // 绘制所有列
			colWidth := int(t.Columns[i][0].Current) // 获取当前列宽度
			if colExtra > 0 {                        // 如果有多余宽度
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
				drag     = &t.header.drags[i]           // 获取当前拖动对象
				colWidth = int(t.Columns[i][0].Current) // 获取当前列宽度
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

func NewNode[T any](data T) (child *Node[T]) {
	return newNode("", false, data)
}

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

func (t *TreeTable[T]) SortNodes() {
	if len(t.Children) == 0 || t.header.SortedBy >= len(t.Children[0].RowCells) {
		return // 如果没有子节点或者列索引无效，直接返回
	}
	sort.Slice(t.Children, func(i, j int) bool {
		if t.Children[i].RowCells == nil { // why? module do not need this
			t.Children[i].RowCells = t.MarshalRow(t.Children[i])
		}
		if t.Children[j].RowCells == nil {
			t.Children[j].RowCells = t.MarshalRow(t.Children[j])
		}
		cellI := t.Children[i].RowCells[t.header.SortedBy].Text
		cellJ := t.Children[j].RowCells[t.header.SortedBy].Text
		if t.header.sortAscending {
			return cellI < cellJ
		}
		return cellI > cellJ
	})
}

const (
	HierarchyIndent = unit.Dp(8)
	iconWidth       = unit.Dp(12)
)

func calculateMaxColumnCellWidth(c CellData) unit.Dp { // 计算最大列单元格宽度
	return c.maxDepth*HierarchyIndent + // 最大深度的左缩进
		iconWidth + // 图标宽度,不管深度是多少，每一行都只会有一个层级图标
		c.maxColumnTextWidth + unit.Dp(8*2) + 20 + // 左右padding,20是sort图标的宽度或者容器节点求和的文本宽度
		DividerWidth // 列分隔条宽度
}

func RowColor(rowIndex int) color.NRGBA { // 奇偶行背景色
	bgColor := color.NRGBA{R: 57, G: 57, B: 57, A: 255}
	if rowIndex%2 != 0 {
		bgColor = color.NRGBA{R: 45, G: 45, B: 45, A: 255}
	}
	return bgColor
}

var modal = NewModal()

func (t *TreeTable[T]) RowFrame(gtx layout.Context, node *Node[T], rowIndex int) layout.Dimensions {
	node.RowCells = t.MarshalRow(node)
	for i := range node.RowCells { // 对齐表头和数据列
		node.RowCells[i].maxColumnTextWidth = t.maxColumnTextWidths[i]
		node.RowCells[i].leftIndent = node.Depth() * HierarchyIndent
		node.RowCells[i].RowID = rowIndex
		node.RowCells[i].Minimum = t.header.ColumnCells[i].Minimum
		node.RowCells[i].Maximum = t.header.ColumnCells[i].Minimum
		node.RowCells[i].Current = t.header.ColumnCells[i].Minimum
		node.RowCells[i].maxDepth = t.header.ColumnCells[i].maxDepth
		node.RowCells[i].ColumID = t.header.ColumnCells[i].ColumID
	}
	rowClick := &node.rowClick

	evt, ok := gtx.Source.Event(pointer.Filter{
		Target: rowClick,
		Kinds:  pointer.Press | pointer.Release,
	})
	if ok {
		e, ok := evt.(pointer.Event)
		if ok {
			if e.Kind == pointer.Press { // 长按应该是touch类型而不是press类型?
				t.selectedNode = node
			}
			switch e.Buttons {
			case pointer.ButtonPrimary, pointer.ButtonSecondary:
				if e.Kind == pointer.Press {
					t.selectedNode = node
				}
			}
		}
	}

	click, ok := rowClick.Update(gtx)
	if ok {
		switch click.NumClicks {
		case 1:
			node.isOpen = !node.isOpen // 切换展开状态
			// todo bug 左键,右键，长按按下选中设置选中节点,但是目前只有左键按下才会被设置
			t.selectedNode = node // 记录被点击的节点,todo 右击也需要填充它,但是在右键菜单中干这个事情似乎时机不对,上下文区域需要支持手动激活方法
			if node.CellClickedCallback != nil {
				node.CellClickedCallback(node) // 单元格点击回调
			}
			if t.RowSelectedCallback != nil {
				t.RowSelectedCallback(node) // 行选中回调
			}

		case 2:
			modal.SetTitle("edit row")
			modal.SetContent(func(gtx layout.Context) layout.Dimensions {
				editNode := NewStructView(node.Data, func() (elems []CellData) {
					return t.MarshalRow(t.selectedNode)
				})
				return editNode.Layout(gtx)
			})

			//if t.RowDoubleClickCallback != nil { // 行双击回调
			//	go t.RowDoubleClickCallback(node)
			//	//gtx.Execute(op.InvalidateCmd{})
			//}
		}
	}
	//if node.LenChildren()%2 == 1 {
	//	rowIndex--
	//}
	bgColor := RowColor(rowIndex)

	if rowClick.Hovered() { // 设置悬停背景色
		bgColor = th.Color.TreeHoveredBgColor
	}
	if t.selectedNode == node { // 设置选中背景色
		bgColor = color.NRGBA{
			R: 255,
			G: 186,
			B: 44,
			A: 91,
		}
		// bgColor = Orange300
	}

	var rowCells []layout.FlexChild

	layoutHierarchyColumn := func(gtx layout.Context, cell CellData) layout.Dimensions {
		c := node.RowCells[0]
		c.leftIndent = node.Depth() * HierarchyIndent
		if !node.CanHaveChildren() {
			c.leftIndent += iconWidth
		}
		if node.parent.IsRoot() {
			c.leftIndent = HierarchyIndent
			if !node.CanHaveChildren() {
				c.leftIndent = HierarchyIndent + iconWidth // 根节点HierarchyIndent + 图标宽度 + 左padding
			}
		}

		// 自适应列宽，这在动态插入节点的情况下可能影响性能
		maxColumnCellWidth := calculateMaxColumnCellWidth(c)
		gtx.Constraints.Min.X = int(maxColumnCellWidth)
		gtx.Constraints.Max.X = int(maxColumnCellWidth)

		return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return rowClick.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					// 绘制层级图标-----------------------------------------------------------------------------------------------------------------
					HierarchyInsert := layout.Inset{Left: c.leftIndent, Top: 0} // 层级图标居中,行高调整后这里需要下移使得图标居中
					if !node.CanHaveChildren() {
						return HierarchyInsert.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Dimensions{}
						})
					}
					return HierarchyInsert.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						svg := CircledChevronRight
						if node.isOpen {
							svg = CircledChevronDown
						}
						return NewButton("", nil).SetRectIcon(true).SetSVGIcon(svg).Layout(gtx)
					})
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				// 绘制层级列文本,和层级图标聚拢在一起-----------------------------------------------------------------------------------------------------------------
				return rowClick.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return node.CellFrame(gtx, c)
				})
			}),
		)
	}

	rowCells = append(rowCells, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layoutHierarchyColumn(gtx, node.RowCells[0])
	}))

	// 绘制非层级列-----------------------------------------------------------------------------------------------------------------
	for i, cell := range node.RowCells[1:] {
		rowCells = append(rowCells, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return rowClick.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Clickable(gtx, &node.RowCells[i].Clickable, func(gtx layout.Context) layout.Dimensions {
					return layout.Stack{Alignment: layout.Center}.Layout(gtx, // 层级列就懒得弹了，copy这个逻辑就行了，要弹的话，长按不支持有点纠结移动平台
						layout.Stacked(func(gtx layout.Context) layout.Dimensions {
							if len(cell.Text) > 80 {
								cell.Text = cell.Text[:len(cell.Text)/2] + "..."
								// todo 这里更新前面已经渲染过的行不方便，所以要在layout或者实例化的时候提前处理
								// 更好的办法是让富文本编辑器做这个事情，对 maxline 。。。 看看代码编辑器扩建是如何实现这个的
								// 然后双击编辑行的时候从富文本取出完整行并换行显示，structView需要好好设计一下这个
								// 这个在抓包场景很那个，url列一般都长
							}
							return node.CellFrame(gtx, cell)
						}),
						layout.Expanded(func(gtx layout.Context) layout.Dimensions {
							if modal.Visible() {
								return modal.Layout(gtx)
							}
							return layout.Dimensions{}
						}),
						layout.Expanded(func(gtx layout.Context) layout.Dimensions {
							if node.rowContextAreas == nil {
								node.rowContextAreas = make([]*component.ContextArea, len(node.RowCells))
							}
							contextArea := node.rowContextAreas[i]
							if contextArea == nil {
								contextArea = &component.ContextArea{
									/*
										var LongPressDuration time.Duration = 250 * time.Millisecond

												case gesture.KindPress:
											i.pressStarted = gtx.Now

											if !i.longPressed && i.pressing && gtx.Now.Sub(i.pressStarted) > LongPressDuration {
												i.longPressed = true
												return Event{Type: LongPress}, true
											}

										所以合理的方案是patch官方的contextAreas和gtx的input source代码，支持长按事件
									*/
									Activation: pointer.ButtonSecondary,
									// todo 根据gioview的作者提示，安卓上需要过滤长按手势事件实现如下:
									// 计算pointer Press到Release的持续时长就可以了，Gio在处理触摸事件和鼠标事件是统一的，
									// 安卓应该也是一致的处理方式，只是event Source变成了Touch。
									// 需要制作一个过滤touch事件的apk测试
									AbsolutePosition: true,
									PositionHint:     0,
								}
								node.rowContextAreas[i] = contextArea
							}
							if node.contextMenu == nil {
								node.contextMenu = NewContextMenu()
								item := ContextMenuItem{}
								for _, kind := range CopyRowType.EnumTypes() {
									switch kind {
									case CopyRowType:
										item = ContextMenuItem{
											Title:     "",
											Icon:      IconCopy,
											Can:       func() bool { return true },
											Do:        func() { node.CopyRow(gtx) },
											Clickable: widget.Clickable{},
										}
									case ConvertToContainerType:
										item = ContextMenuItem{
											Title: "",
											Icon:  IconClean,
											Can: func() bool {
												return node.CanHaveChildren()
											},
											Do: func() {
												mylog.Info("convert to container")
											},
											Clickable: widget.Clickable{},
										}
									case ConvertToNonContainerType:
										item = ContextMenuItem{
											Title: "",
											Icon:  IconActionCode,
											Can: func() bool {
												return node.CanHaveChildren()
											},
											Do: func() {
												mylog.Info("convert to non-container")
											},
											AppendDivider: true,
											Clickable:     widget.Clickable{},
										}
									case NewType:
										item = ContextMenuItem{
											Title: "",
											Icon:  IconArrowDropDown,
											Can: func() bool {
												return true
											},
											Do: func() {
												mylog.CheckNil(t.selectedNode)
												var zero T
												clone := NewNode(zero) // todo 为什么生成了容器节点？
												clone.SetParent(t.selectedNode)

												// clone := t.selectedNode.Clone()
												// clone.Data = zero //

												index := t.selectedNode.RowToIndex() + 1
												switch {
												case t.selectedNode.CanHaveChildren(), t.selectedNode.IsRoot():
													// t.selectedNode.AddChild(clone) //todo 应该插入到选中的孩子下标的后一个，这样是插入到最后一个去了
													t.selectedNode.Children = slices.Insert(t.selectedNode.Children, index, clone)
												default:
													// t.selectedNode.parent.AddChild(clone)
													t.selectedNode.Children = slices.Insert(t.selectedNode.Children, index, clone)
												}
												// 这里应该取已选中的节点，但是这里取右键按下事件并给选中节点赋值，然而右键菜单会因事件执激活菜单失败，弹不出菜单。
												t.SizeColumnsToFit(gtx, false) // 非得排序才能刷新成功新增的节点
												t.List.Update(gtx, layout.Vertical, 0, 0)

												// todo open editor window,把双击编辑节点的代码提到单独的函数，然后调用它
											},
											Clickable: widget.Clickable{},
										}
									case NewContainerType:
										item = ContextMenuItem{
											Title: "",
											Icon:  IconAdd,
											Can: func() bool {
												return true
											},
											Do: func() {
												mylog.Info("new container")
											},
											Clickable: widget.Clickable{},
										}
									case DeleteType:
										item = ContextMenuItem{
											Title:     "",
											Icon:      IconDelete,
											Can:       func() bool { return true },
											Do:        func() { node.RemoveFromParent() },
											Clickable: widget.Clickable{},
										}
									case DuplicateType:
										item = ContextMenuItem{
											Title: "",
											Icon:  IconActionUpdate,
											Can:   func() bool { return true },
											Do: func() {
												mylog.Info("duplicate")
											},
											Clickable: widget.Clickable{},
										}
									case EditType:
										item = ContextMenuItem{
											Title: "",
											Icon:  IconEdit,
											Can: func() bool {
												return true
											},
											Do: func() {
												mylog.Info("edit")
											},
											AppendDivider: true,
											Clickable:     widget.Clickable{},
										}
									case OpenAllType:
										item = ContextMenuItem{
											Title:     "",
											Icon:      IconFileFolderOpen,
											Can:       func() bool { return true },
											Do:        func() { node.OpenAll() },
											Clickable: widget.Clickable{},
										}
									case CloseAllType:
										item = ContextMenuItem{
											Title:     "",
											Icon:      IconClose,
											Can:       func() bool { return true },
											Do:        func() { node.CloseAll() },
											Clickable: widget.Clickable{},
										}
									}
									item.Title = kind.String()
									node.contextMenu.AddItem(item)
								}
								if items := t.ContextMenuItems(node, gtx); items != nil {
									for _, item := range items {
										node.contextMenu.AddItem(item)
									}
								}
							}
							node.contextMenu.OnClicked(gtx)
							return contextArea.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return t.drawContextArea(gtx, &node.contextMenu.MenuState)
							})
						}),
					)
				})
			})
		}))
	}

	rows := []layout.FlexChild{ // 合成层级列和其他列的单元格为一行,并设置该行的背景和行高
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return Background{bgColor}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(19)) //
				gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(19)) //限制行高以避免列分割线呈现虚线视觉
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, rowCells...)
			})
		}),
	}
	if node.CanHaveChildren() && node.isOpen { // 如果是容器节点则递归填充孩子节点形成多行
		for _, child := range node.Children {
			rows = append(rows, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				rowIndex++
				return t.RowFrame(gtx, child, rowIndex)
			}))
		}
	}
	// 把全部行垂直居中排列，rowClick点击后根据点击状态显示了这里填充了多少行，展开节点后看到的行就是这里来的
	return layout.Flex{Axis: layout.Vertical, Spacing: 0, Alignment: layout.Middle, WeightSum: 0}.Layout(gtx, rows...)
}

// -----------------------------------------------------------------------------------------------------------------

func (t *TreeTable[T]) drawContextArea(gtx C, menuState *component.MenuState) D {
	return layout.Center.Layout(gtx, func(gtx C) D { // 重置min x y 到0，并根据max x y 计算弹出菜单的合适大小
		// mylog.Struct("todo",gtx.Constraints)
		gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(4000)) //当行高限制后，这里需要取消限制，理想值是取表格高度或者屏幕高度，其次是增加滚动条或者树形右键菜单
		menuStyle := component.Menu(th.Theme, menuState)
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

//-----------------------------------------------------------------------------------------------------------------

const DividerWidth = unit.Dp(1)

// 分隔线绘制函数
func DrawColumnDivider(gtx layout.Context, col int) {
	if col > 0 { // 层级列不要绘制分隔线
		tallestHeight := gtx.Dp(unit.Dp(gtx.Constraints.Max.Y))
		stack3 := clip.Rect{Max: image.Pt(int(DividerWidth), tallestHeight)}.Push(gtx.Ops)
		paint.Fill(gtx.Ops, DividerFg)
		stack3.Pop()
	}
}

//---------------------------------------泛型n叉树实现------------------------------------------

func (n *Node[T]) AddChildByData(data T) { n.AddChild(NewNode(data)) }

func (n *Node[T]) AddChildrenByDatas(datas ...T) {
	for _, data := range datas {
		n.AddChild(NewNode(data))
	}
}

func (n *Node[T]) AddContainerByData(typeKey string, data T) (newContainer *Node[T]) { // 我们需要返回新的容器节点用于递归填充它的孩子节点，用例是explorer文件资源管理器
	newContainer = NewContainerNode(typeKey, data)
	n.AddChild(newContainer)
	return newContainer
}

func (n *Node[T]) Sum() string {
	// container column 0 key is empty string
	key := n.Type
	key = strings.TrimSuffix(key, ContainerKeyPostfix)
	if n.LenChildren() == 0 {
		return key
	}
	key += " (" + fmt.Sprint(n.LenChildren()) + ")"
	return key
}

func NewUUID() uuid.UUID {
	return mylog.Check2(uuid.NewRandom())
}

func (n *Node[T]) UUID() uuid.UUID {
	return n.ID
}

func (n *Node[T]) Container() bool {
	return strings.HasSuffix(n.Type, ContainerKeyPostfix)
}

func (n *Node[T]) kind(base string) string {
	if n.Container() {
		return base + " Container"
	}
	return base
}

func (n *Node[T]) GetType() string {
	return n.Type
}

func (n *Node[T]) SetType(typeKey string) {
	n.Type = typeKey
}

func (n *Node[T]) IsOpen() bool {
	return n.isOpen && n.Container()
}

func (n *Node[T]) SetOpen(open bool) {
	n.isOpen = open && n.Container()
}

func (n *Node[T]) Parent() *Node[T] {
	return n.parent
}

func (n *Node[T]) SetParent(parent *Node[T]) {
	n.parent = parent
}

func (n *Node[T]) clearUnusedFields() {
	if !n.Container() {
		n.Children = nil
		n.isOpen = false
	}
}

func (n *Node[T]) CanHaveChildren() bool {
	return n.HasChildren()
}

func (n *Node[T]) HasChildren() bool {
	return n.Container() && len(n.Children) > 0
}

func (n *Node[T]) SetChildren(children []*Node[T]) {
	n.Children = children
}

func (n *Node[T]) CellDataForSort(col int) string {
	return n.MarshalRow(n)[col].Text
}

func (n *Node[T]) AddChild(child *Node[T]) {
	child.parent = n
	n.Children = append(n.Children, child)
}

func (n *Node[T]) CellFrame(gtx layout.Context, data CellData) layout.Dimensions {
	// 固定单元格宽度为计算好的每列最大宽度
	gtx.Constraints.Min.X = int(data.Minimum)
	gtx.Constraints.Max.X = int(data.Minimum)

	DrawColumnDivider(gtx, data.ColumID) // 为每列绘制列分隔条

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
	inset := layout.Inset{
		Top:    0,
		Bottom: 0,
		Left:   8 / 2,
		Right:  8 / 2,
	}
	if data.IsHeader { // 加高表头高度
		inset.Top = 2
		inset.Bottom = 2
	}
	return inset.Layout(gtx, material.Body2(th.Theme, data.Text).Layout)
	// return inset.Layout(gtx, richText.Layout)
}

func (n *Node[T]) RootRows() []*Node[T] {
	if n.filteredRows != nil {
		return n.filteredRows
	}
	return n.Children
}

func (n *Node[T]) SetRootRows(rows []*Node[T]) {
	n.filteredRows = nil
	n.Children = rows
	n.SyncToModel()
}

func (n *Node[T]) RowToIndex() int {
	for row, data := range n.Children {
		if data.ID == n.ID {
			return row
		}
	}
	return -1
}

func (n *Node[T]) SyncToModel() {
	//rowCount := 0
	//roots := n.RootRows()
	//if n.filteredRows != nil {
	//	rowCount = len(n.filteredRows)
	//} else {
	//	for _, row := range roots {
	//	}
	//}
}

type Point struct {
	X, Y unit.Dp
}

func (t *TreeTable[T]) ScrollRowIntoView(row int) {
	t.List.ScrollTo(row)
}

// -----------------------------------------------------------------------------

func (t *TreeTable[T]) CopyColumn(gtx layout.Context) string {
	if t.header.clickedColumnIndex < 0 {
		gtx.Execute(clipboard.WriteCmd{Data: io.NopCloser(strings.NewReader("t.header.clickedColumnIndex < 0 "))})
		return "t.header.clickedColumnIndex < 0 "
	}
	b := stream.NewBuffer("var columnData = []string{")
	b.NewLine()
	b.WriteString(strconv.Quote(t.header.ColumnCells[t.header.clickedColumnIndex].Text))
	b.WriteStringLn(",")
	cellData := t.Columns[t.header.clickedColumnIndex]
	for _, datum := range cellData {
		b.WriteString(strconv.Quote(datum.Text))
		b.WriteStringLn(",")
	}
	b.WriteStringLn("}")
	gtx.Execute(clipboard.WriteCmd{Data: io.NopCloser(strings.NewReader(b.String()))})
	return b.String()
}

func (n *Node[T]) CopyRow(gtx layout.Context) string {
	b := stream.NewBuffer("var rowData = []string{")
	cells := n.RowCells
	for i, cell := range cells {
		b.WriteString(strconv.Quote(cell.Text))
		if i < len(cells)-1 {
			b.WriteString(",")
		}
	}
	b.WriteStringLn("}")
	gtx.Execute(clipboard.WriteCmd{Data: io.NopCloser(strings.NewReader(b.String()))})
	return b.String()
}

func (n *Node[T]) IsFiltered() bool {
	return n.filteredRows != nil
}

func (n *Node[T]) ApplyFilter(filter func(row *Node[T]) bool) {
	if filter == nil {
		if n.filteredRows == nil {
			return
		}
		n.filteredRows = nil
	} else {
		n.filteredRows = make([]*Node[T], 0)
		for _, row := range n.RootRows() {
			n.applyFilter(row, filter)
		}
	}
	n.SyncToModel()
	//if n.header != nil && n.header.HasSort() {
	//	n.header.ApplySort()
	//}
}

func (n *Node[T]) applyFilter(row *Node[T], filter func(row *Node[T]) bool) {
	if !filter(row) {
		n.filteredRows = append(n.filteredRows, row)
	}
	if row.CanHaveChildren() {
		for _, child := range row.Children {
			n.applyFilter(child, filter)
		}
	}
}

func (t *TreeTable[T]) Filter(text string) {
	t.filterText = text

	if text == "" {
		t.filteredRows = make([]*Node[T], 0)
		return
	}

	items := make([]*Node[T], 0)
	//for i, item := range t.Children {
	//	if strings.Contains(item.RowCells[i].Cell, text) {
	//		items = append(items, item)
	//	}
	//
	//	for i, child := range item.Children {
	//		if strings.Contains(child.RowCells[i].Cell, text) {
	//			items = append(items, child)
	//		}
	//	}
	//}

	t.filteredRows = items
}

func (n *Node[T]) ApplyFilter_(tag string) {
	n.filteredRows = make([]*Node[T], 0)
	// var node *Node[T]
	// node = n.Root()

	n.WalkContainer(func(node *Node[T]) {
		if node.Container() {
			cells := n.MarshalRow(node)
			for _, cell := range cells {
				if strings.EqualFold(cell.Text, tag) {
					n.filteredRows = append(n.filteredRows, node) // 先过滤所有容器节点
				}
			}
		}
	})

	for i, row := range n.filteredRows {
		children := make([]*Node[T], 0)
		row.Walk(func(node *Node[T]) {
			cells := n.MarshalRow(node)
			for _, cell := range cells {
				if strings.EqualFold(cell.Text, tag) {
					children = append(children, node) // 过滤子节点
				}
			}
		})
		n.filteredRows[i].SetChildren(children)
	}

	n.SetChildren(n.filteredRows)
}

func (n *Node[T]) OpenAll() {
	n.WalkContainer(func(node *Node[T]) {
		if node.Container() {
			node.SetOpen(true)
		}
	})
}

func (t *TreeTable[T]) expandNode(node *Node[T]) {
	node.isOpen = true // 设置节点为展开状态
	for _, child := range node.Children {
		t.expandNode(child) // 递归展开子节点
	}
}

func (n *Node[T]) CloseAll() {
	n.WalkContainer(func(node *Node[T]) {
		if node.Container() {
			node.SetOpen(false)
		}
	})
}

func (n *Node[T]) DiscloseRow(row *Node[T], delaySync bool) bool { // todo merge CloseAll and DiscloseRow
	modified := false
	p := row.Parent()
	var zero *Node[T]
	for p != zero {
		if !p.IsOpen() {
			p.SetOpen(true)
			modified = true
		}
		p = p.Parent()
	}
	if modified {
		n.SyncToModel() // todo 加入是否处于过滤状态的字段，以及重设rootRows后重新layout刷新
	}
	return modified
}

//-----------------------------------------------------------------------------

func CountTableRows[T any](rows []*Node[T]) int { // 计算整个表的总行数
	count := len(rows)
	for _, row := range rows {
		if row.CanHaveChildren() {
			count += CountTableRows(row.Children)
		}
	}
	return count
}

func RowContainsRow[T any](ancestor, descendant *Node[T]) bool { // todo use walk and  rename to Contains
	var zero *Node[T]
	for descendant != zero && descendant != ancestor {
		descendant = descendant.Parent()
	}
	return descendant == ancestor
}

func (n *Node[T]) RemoveFromParent() {
	mylog.CheckNil(n.parent)
	n.parent.Remove(n.ID)
}

func (n *Node[T]) Remove(id uuid.UUID) { // todo add remove callback
	if n.ID == id {
		n.parent.Remove(id)
		return
	}
	for i, child := range n.Children {
		if child.ID == id {
			n.Children = slices.Delete(n.Children, i, i+1)
			return
		}
	}
	mylog.Check("Node not found in parent")
}

func (n *Node[T]) Find(id uuid.UUID) *Node[T] {
	if n.ID == id {
		return n
	}
	for _, child := range n.Children {
		found := child.Find(id)
		if found != nil {
			return found
		}
	}
	return nil
}

func (n *Node[T]) Sort(cmp func(a T, b T) bool) { // todo merge
	sort.SliceStable(n.Children, func(i, j int) bool {
		return cmp(n.Children[i].Data, n.Children[j].Data)
	})
	for _, child := range n.Children {
		child.Sort(cmp)
	}
}

func (n *Node[T]) Walk(callback func(node *Node[T])) {
	callback(n)
	for _, child := range n.Children {
		child.Walk(callback)
	}
}

func (n *Node[T]) WalkQueue(callback func(node *Node[T])) {
	queue := []*Node[T]{n}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		callback(node)
		for _, child := range node.Children {
			queue = append(queue, child)
		}
	}
}

func (n *Node[T]) Containers() []*Node[T] {
	containers := make([]*Node[T], 0)
	for _, child := range n.Children {
		if child.Container() {
			containers = append(containers, child)
		}
	}
	return containers
}

func (n *Node[T]) WalkContainer(callback func(node *Node[T])) {
	callback(n) // always walk root here
	containers := make([]*Node[T], 0)
	for _, child := range n.Children {
		if child.Container() {
			containers = append(containers, child)
		}
	}
	for _, container := range containers {
		container.Walk(callback)
	}
}

func (t *TreeTable[T]) FormatHeader(maxColumnCellTextWidths []unit.Dp) *stream.Buffer {
	buf := stream.NewBuffer("")

	all := t.maxColumnCellWidth
	for _, width := range maxColumnCellTextWidths {
		all += width
	}
	all += align.StringWidth[unit.Dp]("│")*unit.Dp(len(maxColumnCellTextWidths)) + 4 //?

	buf.WriteStringLn("┌─" + strings.Repeat("─", int(all)))
	buf.WriteString("│")

	// 计算每个单元格的左边距
	for i, cell := range t.header.ColumnCells {
		paddedText := fmt.Sprintf("%-*s", int(maxColumnCellTextWidths[i]), cell.Text) // 左对齐填充

		// 添加左边距，仅在首列进行处理，依据列宽计算
		if i == 0 {
			buf.WriteString(strings.Repeat(" ", int(t.maxColumnCellWidth-maxColumnCellTextWidths[i]-1))) // -1是分隔符的空间
		}

		buf.WriteString(paddedText)
		if i < len(t.header.ColumnCells)-1 {
			buf.WriteString(" │ ") // 在每个单元格之间添加分隔符
		}
	}

	buf.NewLine()
	buf.WriteStringLn("├─" + strings.Repeat("─", int(all)))
	return buf
}

func (t *TreeTable[T]) FormatChildren(out *stream.Buffer, children []*Node[T]) {
	for i, child := range children {
		child.RowCells = t.MarshalRow(child)
		HierarchyColumBuf := stream.NewBuffer("")
		for j, cell := range child.RowCells {
			if j == 0 {
				HierarchyColumBuf.WriteString("│")
				HierarchyColumBuf.WriteString(strings.Repeat(" ", int(child.Depth()*indentBase)))
				if i == len(children)-1 {
					HierarchyColumBuf.WriteString("╰──") //"└───"
				} else {
					HierarchyColumBuf.WriteString("├──")
				}
				HierarchyColumBuf.WriteString(cell.Text)
				if align.StringWidth[unit.Dp](HierarchyColumBuf.String()) < t.maxColumnCellWidth {
					HierarchyColumBuf.WriteString(strings.Repeat(" ", int(t.maxColumnCellWidth-align.StringWidth[unit.Dp](HierarchyColumBuf.String()))))
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

const (
	indent          = "│   "
	childPrefix     = "├───"
	lastChildPrefix = "└───"
	indentBase      = unit.Dp(3)
)

func (t *TreeTable[T]) MaxColumnCellWidth() unit.Dp {
	HierarchyIndent := unit.Dp(1)
	DividerWidth := align.StringWidth[unit.Dp](" │ ")
	iconWidth := align.StringWidth[unit.Dp](childPrefix)
	return t.MaxDepth()*HierarchyIndent + // 最大深度的左缩进
		iconWidth + // 图标宽度,不管深度是多少，每一行都只会有一个层级图标
		t.maxColumnTextWidths[0] + 5 + //(8 * 2) + 20 + // 左右padding,20是sort图标的宽度或者容器节点求和的文本宽度
		DividerWidth // 列分隔条宽度
}

func (t *TreeTable[T]) Format() *stream.Buffer {
	t.SizeColumnsToFit(layout.Context{}, true)
	t.maxColumnCellWidth = t.MaxColumnCellWidth()
	buf := t.FormatHeader(t.maxColumnTextWidths)
	t.FormatChildren(buf, t.Children) // 传入子节点打印函数
	mylog.Json("RootRows", buf.String())
	return buf
}

func (t *TreeTable[T]) String() string {
	return t.Format().String()
}

func (t *TreeTable[T]) Document() string {
	s := stream.NewBuffer("")
	// s.WriteStringLn("// interface or method name here")
	// s.WriteStringLn("/*")
	lines := t.Format().ToLines()
	for _, line := range lines {
		s.WriteStringLn("  " + line)
	}
	// s.WriteStringLn("*/")
	return s.String()
}

func (n *Node[T]) Depth() unit.Dp {
	if n.parent != nil {
		return n.parent.Depth() + 1
	}
	return 1
}

func (n *Node[T]) LenChildren() int {
	return len(n.Children)
}

func (n *Node[T]) LastChild() (lastChild *Node[T]) {
	if n.IsRoot() {
		return n.Children[len(n.Children)-1]
	}
	return n.parent.Children[len(n.parent.Children)-1]
}

func (n *Node[T]) IsLastChild() bool {
	return n.LastChild() == n
}

func (n *Node[T]) ResetChildren() {
	n.Children = nil
	n.filteredRows = nil
}

func (n *Node[T]) CopyFrom(from *Node[T]) *Node[T] {
	*n = *from
	return n
}

func (n *Node[T]) ApplyTo(to *Node[T]) *Node[T] {
	*to = *n
	return n
}

func (n *Node[T]) Clone() (newNode *Node[T]) {
	defer n.SyncToModel()
	if n.Container() {
		return NewContainerNode(n.Type, n.Data)
	}
	return NewNode(n.Data)
}

// ---------------------- todo delete
//func NewTable3[T any](data T, ctx TableContext[T]) {
//	if ctx.JsonName == "" {
//		mylog.Check("JsonName is empty")
//	}
//	ctx.JsonName = strings.TrimSuffix(ctx.JsonName, ".json")
//
//	mylog.CheckNil(ctx.UnmarshalRow)
//	// mylog.CheckNil(ctx.SetRootRowsCallBack)//mitmproxy
//	mylog.CheckNil(ctx.RowSelectedCallback)
//
//	//table, header = newTable(data, ctx)
//	//fnUpdate := func() {
//	//	table.SetRootRows(table.Children)
//	//	table.SizeColumnsToFit(true)
//	//	stream.MarshalJsonToFile(table, filepath.Join("cache", ctx.JsonName+".json"))
//	//	stream.WriteTruncate(filepath.Join("cache", ctx.JsonName+".txt"), table.Document())
//	//	if ctx.IsDocument {
//	//		b := stream.NewBuffer("")
//	//		b.WriteStringLn("# " + ctx.JsonName + " document table")
//	//		b.WriteStringLn("```text")
//	//		b.WriteStringLn(table.Document())
//	//		b.WriteStringLn("```")
//	//		stream.WriteTruncate("README.md", b.String())
//	//	}
//	//}
//
//	//if ctx.SetRootRowsCallBack != nil { // mitmproxy
//	//	ctx.SetRootRowsCallBack(table)
//	//}
//	//table.RowSelectedCallback = func() { ctx.RowSelectedCallback(table) }
//	//
//	//app.FileDropCallback(func(files []string) {
//	//	if filepath.Ext(files[0]) == ".json" {
//	//		mylog.Info("dropped file", files[0])
//	//		table.ResetChildren()
//	//		b := stream.NewBuffer(files[0])
//	//		mylog.Check(json.Unmarshal(b.Bytes(), table)) // todo test need a zero table?
//	//		fnUpdate()
//	//	}
//	//	mylog.Struct("todo",files)
//	//})
//	//
//	//table.RowDoubleClickCallback = func() {
//	//	rows := table.SelectedRows(false)
//	//	for _, row := range rows {
//	//		mylog.Struct("todo",row.RowCells)
//	//		// todo icon edit
//	//		//app.RunWithIco("edit row #"+fmt.Sprint(i), rowPngBuffer, func(w *Window) {
//	//		//	content := w.Content()
//	//		//	nodeEditor, RowPanel := NewStructView(row.Cell, func(data T) (values []CellData) {
//	//		//		return table.MarshalRow(row)
//	//		//	})
//	//		//	content.AddChild(nodeEditor)
//	//		//	content.AddChild(RowPanel)
//	//		//	panel := NewButtonsPanel(
//	//		//		[]string{
//	//		//			"apply", "cancel",
//	//		//		},
//	//		//		func() {
//	//		//			ctx.UnmarshalRow(row, nodeEditor.getFieldValues())
//	//		//			nodeEditor.Update(row.Cell)
//	//		//			table.SyncToModel()
//	//		//			stream.MarshalJsonToFile(table.Children, ctx.JsonName+".json")
//	//		//			// w.Dispose()
//	//		//		},
//	//		//		func() {
//	//		//			w.Dispose()
//	//		//		},
//	//		//	)
//	//		//	RowPanel.AddChild(panel)
//	//		//	RowPanel.AddChild(NewVSpacer())
//	//		//})
//	//	}
//	//}
//	//fnUpdate()
//	return
//}
//func newTable[T any](data T, ctx TableContext[T]) {
//	//for i, column := range root.Columns {
//	//	text := NewText(column.Cell, &TextDecoration{
//	//		Font:           LabelFont,
//	//		Background:     nil,
//	//		Foreground:     nil,
//	//		BaselineOffset: 0,
//	//		Underline:      false,
//	//		StrikeThrough:  false,
//	//	})
//	//	root.Columns[i].Minimum = text.Current() + root.Padding.Left + root.Padding.Right
//	//}
//	//
//	//root.KeyDownCallback = func(keyCode KeyCode, mod Modifiers, repeat bool) bool {
//	//	if mod == 0 && (keyCode == KeyBackspace || keyCode == KeyDelete) {
//	//		root.PerformCmd(root, DeleteItemID)
//	//		return true
//	//	}
//	//	return root.DefaultKeyDown(keyCode, mod, repeat) // todo add delete,move to ctx menu,exporter need delete file or dir
//	//}
//	//
//	//root.InstallDragSupport(nil, "dragKey", "singularName", "pluralName")
//	//InstallDropSupport[T, any](root, "dragKey", func(from, to *Node[T]) bool { return from == to }, nil, nil)
//	//
//	//header = NewTableHeader[T](root)
//	//for _, column := range root.Columns {
//	//	columnHeader := NewTableColumnHeader[T](column.Cell, "")
//	//	columnHeader.MouseDownCallback = func(where Point, button, clickCount int, mod Modifiers) bool {
//	//		return true
//	//	}
//	//	NewContextMenuItems(columnHeader, columnHeader.CellsLabel.MouseDownCallback,
//	//		ContextMenuItem{
//	//			Title: "copy column",
//	//			id:    0,
//	//			Can: func(a any) bool {
//	//				return true
//	//			},
//	//			Do: func(a any) { root.CopyColumn() },
//	//		},
//	//	).Install()
//	//	header.ColumnHeaders = append(header.ColumnHeaders, columnHeader)
//	//}
//	//
//	//return root, header
//}

// func addWrappedText(parent *Node[T], ink color.NRGBA, font font.Face, data CellData) {
//
//	if data.IsNasm {
//		tokens, style := languages.GetTokens(stream.NewBuffer(data.Content), languages.NasmKind)
//		rowPanel := NewPanel()
//		rowPanel.SetLayout(&FlexLayout{Columns: len(tokens)})
//		parent.AddChild(rowPanel)
//		for _, token := range tokens {
//			colour := style.GetMust(token.Type).Colour
//			keys := NewRichLabel()
//			color := RGB(
//				int(colour.Red()),
//				int(colour.Green()),
//				int(colour.Blue()),
//			)
//			keys.Content = NewText(token.Value, &TextDecoration{
//				Font:           font,
//				Foreground:     color,
//				Background:     ContentColor,
//				BaselineOffset: 0,
//				Underline:      false,
//				StrikeThrough:  false,
//			})
//			keys.OnBackgroundInk = ink
//			rowPanel.AddChild(keys)
//		}
//		return
//	}
//
// decoration := &TextDecoration{Font: font}
// var lines []*Content
//
//	if data.MaxWidth > 0 {
//		lines = NewTextWrappedLines(data.Content, decoration, data.MaxWidth)
//	} else {
//
//		lines = NewTextLines(data.Content, decoration)
//	}
//
//	for _, line := range lines {
//		keys := NewLabel()
//		keys.Content = line.String()
//		keys.Font = font
//		keys.LabelTheme.OnBackgroundInk = ink
//		if data.FgColor != 0 {
//			keys.LabelTheme.OnBackgroundInk = data.FgColor
//		}
//		if data.Disabled {
//			keys.LabelTheme.OnBackgroundInk = DarkRed
//		}
//		size := LabelFont.Size() + 7
//		if data.ImageBuffer != nil {
//			keys.Drawable = &SizedDrawable{
//				Drawable: nil,
//				Size:     NewSize(size, size),
//			}
//		}
//		if data.SvgBuffer != "" {
//			keys.Drawable = &DrawableSVG{
//				SVG:  nil,
//				Size: NewSize(size, size),
//			}
//		}
//		// LabelStyle(keys)
//		parent.AddChild(keys)
//	}
//
// }

//const iconSize = unit.Dp(12)
//
//var svgButtonTest *GoogleDummyButton
//
//func (t *TreeTable[T]) ColumnCell(row, col int, foreground, background color.NRGBA, selected, indirectlySelected, focused bool) layout.Dimensions {
//	return layout.Dimensions{}
//}

// var editNode *StructView
