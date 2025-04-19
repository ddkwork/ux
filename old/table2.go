package main

import (
	"fmt"
	"image"

	"github.com/ddkwork/ux"
	"github.com/ddkwork/ux/resources/colors"

	"gioui.org/gesture"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/ux/widget/material"
)

type (
	Table2 struct {
		Columns      []Column2
		prevMetric   unit.Metric
		prevMaxWidth int

		SortOrder          sortOrder
		SortedBy           int
		drags              []tableDrag
		headers            []*widget.Clickable
		cells              []*widget.Clickable // 单元格单击事件
		ClickedColumnIndex int
		manualWidthSet     []bool // 新增状态标志数组，记录列是否被手动调整

		*widget.List
	}
	Column2 struct {
		Name      string
		Width     unit.Dp
		MinWidth  unit.Dp
		Alignment text.Alignment
	}
	tableDrag struct {
		drag           gesture.Drag
		hover          gesture.Hover
		startPos       float32
		shrinkNeighbor bool
	}
)

func NewTable2(columns []Column2) *Table2 {
	headers := make([]*widget.Clickable, len(columns))
	for i := range headers {
		headers[i] = &widget.Clickable{}
	}
	return &Table2{
		Columns:            columns,
		SortOrder:          0,
		SortedBy:           0,
		prevMetric:         unit.Metric{},
		prevMaxWidth:       0,
		drags:              nil,
		headers:            headers,
		cells:              nil,
		ClickedColumnIndex: -1,
		List: &widget.List{
			Scrollbar: widget.Scrollbar{},
			List: layout.List{
				Axis:        layout.Vertical,
				ScrollToEnd: false,
				Alignment:   0,
				Position:    layout.Position{},
			},
		},
		manualWidthSet: make([]bool, len(columns)),
	}
}

type (
	sortOrder uint8
	cellFn    func(gtx layout.Context, row, col int) layout.Dimensions
	rowFn     func(gtx layout.Context, row int) layout.Dimensions
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

// SetColumns SizeColumnToFit
func (t *Table2) SetColumns(gtx layout.Context, cols []Column2, cellData [][]string) {
	originalConstraints := gtx.Constraints // 保存原始约束
	for i := range cols {
		if t.manualWidthSet[i] { // 如果该列已手动调整
			continue // 跳过，保留用户手动调整的宽度
		}
		maxWidth := unit.Dp(0) // 当前最大宽度初始化为0
		for _, data := range cellData {
			if i < len(data) {
				currentWidth := ux.LabelWidth(gtx, data[i])
				// currentWidth += float32(gtx.Dp(material.Scrollbar(th, nil).Width())) + float32(len(cols)*gtx.Dp(defaultDividerWidth))
				if currentWidth > maxWidth {
					maxWidth = currentWidth // 更新最大宽度
				}
			}
		}
		cols[i].Width = max(maxWidth, ux.LabelWidth(gtx, cols[i].Name+" ⇧"))
	}
	gtx.Constraints = originalConstraints
	t.Columns = cols
}

func (t *Table2) Update(gtx layout.Context) {
	for _, btn := range t.headers {
		for {
			evt, ok := gtx.Source.Event(pointer.Filter{
				Target: btn,
				Kinds:  pointer.Press | pointer.Release | pointer.Drag,
			})
			if !ok {
				break
			}
			e, ok := evt.(pointer.Event)
			if !ok {
				break
			}
			switch e.Buttons {
			case pointer.ButtonPrimary:
				if e.Kind == pointer.Press {
				}
			case pointer.ButtonSecondary:

			}
		}
	}
}

func (t *Table2) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	t.resize(gtx)
	dims := w(gtx)
	dims.Size = gtx.Constraints.Constrain(dims.Size)
	return dims
}

func (t *Table2) ClickedColumn() (int, bool) {
	index := t.ClickedColumnIndex
	return index, index >= 0
}

func (t *Table2) SortByClickedColumn() (int, bool) {
	if col, ok := t.ClickedColumn(); ok {
		if col == t.SortedBy {
			switch t.SortOrder {
			case sortNone:
				t.SortOrder = sortAscending
			case sortAscending:
				t.SortOrder = sortDescending
			case sortDescending:
				t.SortOrder = sortAscending
			default:
				panic(fmt.Sprintf("unhandled case %v", t.SortOrder))
			}
		} else {
			t.SortedBy = col
			t.SortOrder = sortAscending
		}
		return col, true
	}
	return 0, false
}

func (t *Table2) resize(gtx layout.Context) {
	// 检查当前约束是否与之前记录一致，如果一致则不需要调整
	if gtx.Constraints.Max.X == t.prevMaxWidth && gtx.Metric == t.prevMetric {
		return
	}

	var (
		oldAvailable = unit.Dp(t.prevMaxWidth - t.prevMetric.Dp(material.Scrollbar(th, nil).Width()) - len(t.Columns)*t.prevMetric.Dp(defaultDividerWidth)) // 之前的可用宽度
		available    = unit.Dp(gtx.Constraints.Max.X - gtx.Dp(material.Scrollbar(th, nil).Width()) - len(t.Columns)*gtx.Dp(defaultDividerWidth))            // 当前可用宽度
	)

	// 避免负值和零导致计算错误
	if oldAvailable < 1 {
		oldAvailable = 1
	}
	if available < 1 {
		available = 1
	}

	defer func() {
		// 更新前一个约束和度量
		t.prevMaxWidth = gtx.Constraints.Max.X
		t.prevMetric = gtx.Metric
	}()

	if available > oldAvailable { // 如果当前可用空间大于之前的可用空间
		var totalWidth unit.Dp
		for i := range t.Columns { // 计算所有列的总宽度
			totalWidth += t.Columns[i].Width
		}
		// 如果总宽度超出可用空间，跳过更新
		if totalWidth > (available) {
			return
		}
	}

	var (
		dividerWidth       = gtx.Dp(defaultDividerWidth)       // 获取分隔符的宽度
		dividerMargin      = gtx.Dp(defaultDividerMargin)      // 获取分隔符的边距
		dividerHandleWidth = gtx.Dp(defaultDividerHandleWidth) // 获取分隔处理器的宽度

		globalMinWidth = unit.Dp(dividerWidth + dividerMargin + dividerHandleWidth) // 计算全局最小宽度
	)

	// 遍历每一列，更新列宽
	for i := range t.Columns {
		col := &t.Columns[i]
		// 计算新的列宽度
		r := col.Width / (oldAvailable)
		col.Width = max(col.MinWidth, globalMinWidth, r*(available)) // 确保更新后的宽度不小于最小宽度
	}
}

type TableRowStyle struct {
	Table  *Table2
	Header bool
}

func TableRow(tbl *Table2, hdr bool) *TableRowStyle {
	return &TableRowStyle{
		Table:  tbl,
		Header: hdr,
	}
}

func (row *TableRowStyle) Layout(gtx layout.Context, w rowFn) layout.Dimensions {
	var (
		cols          = len(row.Table.Columns) // 获取列的数量
		dividers      = cols                   // 列的分隔数
		tallestHeight = gtx.Constraints.Min.Y  // 初始化最高的行高

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

	if len(row.Table.drags) < dividers { // 如果拖动数组没有足够的长度
		row.Table.drags = make([]tableDrag, dividers) // 初始化拖动对象数组
	}

	// OPT(dh): we don't need to do this for each row, only once per table
	for i := range row.Table.drags { // 遍历每个拖动对象
		drag := &row.Table.drags[i]   // 获取当前拖动对象
		col := &row.Table.Columns[i]  // 获取当前列
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
				row.Table.manualWidthSet[i] = true // 在这里标记列宽已手动调整
				// There may be multiple drag events in a single frame. We mustn't apply all of them or we'll
				// drag too far. Only react to the final event.
				delta = unit.Dp(ev.Position.X - drag.startPos) // 计算当前拖动的偏移量
			}
		}
		if delta != 0 { // 如果存在拖动偏移量
			col.Width += delta                                        // 更新列的宽度
			if drag.shrinkNeighbor && i != len(row.Table.Columns)-1 { // 如果需要收缩相邻列且不是最后一列
				nextCol := &row.Table.Columns[i+1] // 获取下一个列
				nextCol.Width -= delta             // 更新下一个列的宽度
				if col.Width < minWidth {          // 如果当前列宽度小于最小宽度
					d := minWidth - col.Width // 计算需要增加的宽度
					col.Width = minWidth      // 将当前列宽度设为最小宽度
					nextCol.Width -= d        // 更新下一个列的宽度
				}
				if nextCol.Width < minWidth { // 如果下一个列宽度小于最小宽度
					d := minWidth - nextCol.Width // 计算需要增加的宽度
					nextCol.Width = minWidth      // 将下一个列宽度设为最小宽度
					col.Width -= d                // 更新当前列宽度
				}
			} else {
				// 如果不需要收缩
				if col.Width < minWidth { // 如果当前列宽度小于最小宽度
					col.Width = minWidth // 将当前列宽度设为最小宽度
				}
			}

			if col.Width < col.MinWidth { // 如果当前列宽度小于其最小宽度
				col.MinWidth = col.Width // 更新列的最小宽度为当前宽度
			}

			var total unit.Dp                       // 初始化总宽度
			for _, col := range row.Table.Columns { // 遍历所有列计算总宽度
				total += col.Width // 累加当前列的宽度
			}
			total += unit.Dp(len(row.Table.Columns) * gtx.Dp(defaultDividerWidth)) // 加上所有分隔符的总宽度
			if total < unit.Dp(gtx.Constraints.Min.X) {                            // 如果总宽度小于最小约束宽度
				row.Table.Columns[len(row.Table.Columns)-1].Width += unit.Dp(gtx.Constraints.Min.X) - total // 调整最后一列的宽度以适应
			}
		}
	}

	for { // 开始绘制列
		// First draw all columns, leaving gaps for the drag handlers
		var (
			start             = 0             // 初始化当前位置
			origTallestHeight = tallestHeight // 记录最初的高度
		)
		r := op.Record(gtx.Ops)            // 记录当前操作集合
		totalWidth := 0                    // 初始化总宽度
		for i := range row.Table.Columns { // 遍历所有列
			colWidth := int(row.Table.Columns[i].Width) // 获取当前列的宽度
			totalWidth += colWidth                      // 更新总宽度
		}
		extra := gtx.Constraints.Min.X - len(row.Table.Columns)*gtx.Dp(defaultDividerWidth) - totalWidth // 计算多余宽度
		colExtra := extra                                                                                // 将多余宽度赋值给列额外宽度

		for i := range row.Table.Columns { // 绘制所有列
			colWidth := int(row.Table.Columns[i].Width) // 获取当前列宽度
			if colExtra > 0 {                           // 如果有多余宽度
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
		for i := range row.Table.drags { // 遍历每个拖动对象
			var (
				drag     = &row.Table.drags[i]             // 获取当前拖动对象
				colWidth = int(row.Table.Columns[i].Width) // 获取当前列宽度
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

			if row.Header { // 如果当前行是表头
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

					paint.FillShape(gtx.Ops, colors.Red200, handleLeft.Op(gtx.Ops))     // 填充左侧形状
					paint.FillShape(gtx.Ops, colors.Yellow100, handleRight.Op(gtx.Ops)) // 填充右侧形状
				}

				// Draw the vertical bar
				// stack3 := clip.Rect{Max: image.Pt(dividerWidth, tallestHeight)}.Push(gtx.Ops)
				// Fill( gtx.Ops, win.Theme.Palette.Table2.Divider) // 如果有需要，在此处可绘制分割线
				// stack3.Pop()
			}
			// 为表头和每列绘制列分隔条
			stack3 := clip.Rect{Max: image.Pt(dividerWidth, tallestHeight)}.Push(gtx.Ops) // 绘制分隔条的矩形区域
			paint.Fill(gtx.Ops, colors.DividerFg)                                         // 填充分隔条的颜色
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

type TableHeaderRowStyle struct {
	Table *Table2
}

func TableHeaderRow(tbl *Table2) *TableHeaderRowStyle {
	return &TableHeaderRowStyle{Table: tbl}
}

func (row *TableHeaderRowStyle) Layout(gtx layout.Context) layout.Dimensions {
	return TableRow(row.Table, true).Layout(gtx, func(gtx layout.Context, col int) layout.Dimensions {
		var (
			cell   = &row.Table.Columns[col]
			button = row.Table.headers[col]
		)
		if button.Clicked(gtx) {
			mylog.Trace("header clicked", col)
			row.Table.ClickedColumnIndex = col
		}
		gtx.Constraints.Min.Y = gtx.Dp(20)                                                       // 限制行高
		paint.FillShape(gtx.Ops, colors.ColorHeaderFg, clip.Rect{Max: gtx.Constraints.Max}.Op()) // 表头背景色
		return layout.Inset{
			// Top: defaultHeaderPadding,
			// Top:    5,
			// Bottom: 5,
			// Left:   3,
			// Right:  0,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			var s string
			// OPT(dh): avoid allocations for string building by precomputing and storing the column clickables.
			if row.Table.SortedBy == col {
				switch row.Table.SortOrder {
				case sortNone:
					s = cell.Name
				case sortAscending:
					s = "⇧" + cell.Name
				case sortDescending:
					s = "⇩" + cell.Name
				default:
					panic(fmt.Sprintf("unhandled case %v", row.Table.SortOrder))
				}
			} else {
				s = cell.Name
			}
			return material.Clickable(gtx, button, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:    0,
					Bottom: 0,
					Left:   0,
					Right:  0,
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.Y = gtx.Dp(20) // 限制行高
					body1 := material.Body1(th, s)
					body1.MaxLines = 1
					return body1.Layout(gtx)
				})
			})
		})
	})
}

type TableSimpleRowStyle struct {
	Table *Table2
}

func TableSimpleRow(tbl *Table2) TableSimpleRowStyle {
	return TableSimpleRowStyle{Table: tbl}
}

func (row TableSimpleRowStyle) Layout(gtx layout.Context, rowIdx int, cellFn cellFn) layout.Dimensions {
	c := ux.RowColor(rowIdx)

	if row.Table.cells == nil {
		row.Table.cells = make([]*widget.Clickable, len(row.Table.Columns))
		for i := range row.Table.cells {
			row.Table.cells[i] = &widget.Clickable{}
		}
	}
	if rowIdx >= len(row.Table.cells) {
		row.Table.cells = make([]*widget.Clickable, rowIdx*len(row.Table.Columns))
		for i := range row.Table.cells {
			row.Table.cells[i] = &widget.Clickable{}
		}
	}
	hover := row.Table.cells[rowIdx]
	update, b := hover.Update(gtx)
	if b {
		if update.NumClicks == 1 {
			// c = th.ContrastBg
			// c = ColorPink
			// HighlightRow(gtx)
			c = colors.Red400
			// paint.FillShape(gtx.Ops, ColorPink, clip.Rect{Max: gtx.Constraints.Max}.Op())
		}
		if row.Table.cells[rowIdx].Hovered() {
			c = th.ContrastFg
			// HighlightRow(gtx)
		}
	}
	return ux.Background{Color: c}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return TableRow(row.Table, false).Layout(gtx, func(gtx layout.Context, col int) layout.Dimensions {
			defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
			gtx.Constraints.Min.Y = 0
			return row.Table.cells[rowIdx].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Clickable(gtx, row.Table.cells[rowIdx], func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top:    5,
						Bottom: 5,
						Left:   5,
						Right:  0,
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						// gtx.Constraints.Min.Y = gtx.Dp(27) //限制行高
						return cellFn(gtx, rowIdx, col)
					})
				})
			})
		})
	})
}

func SimpleTable(gtx layout.Context, tbl *Table2, rows int, cellFn cellFn) layout.Dimensions {
	return tbl.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return TableHeaderRow(tbl).Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.List(th, tbl.List).Layout(gtx, rows, func(gtx layout.Context, row int) layout.Dimensions {
					return TableSimpleRow(tbl).Layout(gtx, row, cellFn)
				})
			}),
		)
	})
}
