package ux

import (
	"image"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/ddkwork/ux/resources/colors"
)

// Resize provides a draggable handle in between two widgets for resizing their area.
type Resize struct {
	// axis defines how the widgets and the handle are laid out.
	axis               layout.Axis // 水平拖动设置列宽，垂直拖动设置行高
	initialized        bool
	length             int
	totalHandlesLength int
	cells              []*Resizable // 表头是colum cell，行是row cell
	minLength          int          // 这里应该设置为树形表格每列的最大列宽
}

type (
	ResizeCallback func(index int, newWidth int)
	Resizable      struct {
		// ratio is only calculated during initialization, based on widget's natural size.
		//  It acts like minimum threshold ratio value beyond which widget size cannot be further reduced.
		ratio          float32
		Widget         layout.Widget
		DividerHandler layout.Widget
		Index          int            // 新增索引字段
		OnResize       ResizeCallback // 新增回调字段
		// dividerThickness int
		float
		resize *Resize
		prev   *Resizable
		next   *Resizable
	}
)

func NewResize(axis layout.Axis, onResize ResizeCallback, cells ...*Resizable) *Resize {
	r := &Resize{axis: axis, cells: cells}
	for i, rz := range cells {
		rz.Index = i // 设置索引
		rz.resize = r
		if rz.DividerHandler == nil {
			rz.DividerHandler = r.CustomResizeHandleBar
		}
		rz.OnResize = onResize // 设置回调函数
	}
	return r
}

// Layout displays w1 and w2 with handle in between.
//
// The widgets w1 and w2 must be able to gracefully resize their minimum and maximum dimensions
// in order for the resize to be smooth.
func (r *Resize) Layout(gtx layout.Context) layout.Dimensions {
	// Compute the first widget's max dragWidth/height.
	if len(r.cells) == 0 {
		return layout.Dimensions{}
	}
	if len(r.cells) == 1 {
		return r.cells[0].Widget(gtx)
	}

	if !r.initialized {
		r.init(gtx)
		r.initialized = true
	}

	// On Window Resize
	if r.length != r.axis.Convert(gtx.Constraints.Max).X {
		r.onWindowResize(gtx)
	}
	gtx.Constraints.Min = gtx.Constraints.Max

	flex := layout.Flex{Axis: r.axis}
	return flex.Layout(gtx,
		r.cells[0].Layout(gtx)...,
	)
}

func (r *Resize) init(gtx layout.Context) {
	r.length = r.axis.Convert(gtx.Constraints.Max).X
	if r.minLength == 0 {
		r.minLength = int(0.1 * float32(r.length))
	}
	allowedMinLength := r.length / len(r.cells)
	if r.minLength > allowedMinLength || r.minLength <= 0 {
		r.minLength = allowedMinLength
	}
	var totalRatio float32
	// Obtain the total ration to reset it between 0.0 - 1.00
	var totalHandlesLength int
	for i, cell := range r.cells {
		if cell.DividerHandler == nil {
			cell.DividerHandler = r.CustomResizeHandleBar
		}
		m := op.Record(gtx.Ops)
		d := cell.DividerHandler(gtx)
		m.Stop()
		totalHandlesLength += r.axis.Convert(d.Size).X
		m = op.Record(gtx.Ops)
		d = cell.Widget(gtx)
		m.Stop()
		cell.ratio = float32(r.axis.Convert(d.Size).X) / float32(r.length)
		totalRatio += cell.ratio
		var prevResizable *Resizable
		var nextResizable *Resizable
		if i != 0 {
			prevResizable = r.cells[i-1]
		}
		if i < len(r.cells)-1 {
			nextResizable = r.cells[i+1]
		}
		cell.prev = prevResizable
		cell.next = nextResizable
		cell.resize = r
		if i == len(r.cells)-1 {
			if totalRatio <= 1 {
				totalRatio -= cell.ratio
				cell.ratio = 1 - totalRatio
				totalRatio = 1
			}
		}
	}
	r.totalHandlesLength = totalHandlesLength
	// Reset the ratio between 0.0 - 1.00
	var currTotalRatio float32
	for _, cell := range r.cells {
		cell.ratio /= totalRatio // reset the total ratio
		currTotalRatio += cell.ratio
		cell.float.pos = int(float32(r.length) * currTotalRatio)
	}
	// mylog.Struct(totalRatio)
	// mylog.Struct(currTotalRatio)
}

func (r *Resize) onWindowResize(gtx layout.Context) {
	currMinLength := r.minLength
	prevLength := r.length
	r.minLength = (currMinLength / prevLength) * r.length
	r.length = r.axis.Convert(gtx.Constraints.Max).X
	for _, cell := range r.cells {
		cell.float.pos = int((float32(cell.float.pos) / float32(prevLength)) * float32(r.length))
	}
}

type float struct {
	pos  int // position in pixels of the handle
	drag gesture.Drag
}

func (r *Resizable) Layout(gtx layout.Context) []layout.FlexChild {
	m := op.Record(gtx.Ops)
	dims := r.handleDrag(gtx)
	c := m.Stop()
	children := []layout.FlexChild{
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			prePos := 0
			if r.prev != nil {
				prePos = r.prev.pos
			}
			gtx.Constraints.Max = image.Point{X: r.pos - prePos, Y: gtx.Constraints.Max.Y}
			if r.resize.axis == layout.Vertical {
				gtx.Constraints.Max = r.resize.axis.Convert(gtx.Constraints.Max)
			}
			d := r.Widget(gtx)
			d.Size = gtx.Constraints.Max
			return d
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			c.Add(gtx.Ops)
			return dims
		}),
	}
	if r.next != nil {
		children = append(children, r.next.Layout(gtx)...)
	}
	return children
}

func (r *Resizable) handleDrag(gtx layout.Context) layout.Dimensions {
	if r.next == nil {
		return layout.Dimensions{}
	}
	gtx.Constraints.Min = image.Point{}
	dims := r.DividerHandler(gtx)

	var de *pointer.Event

	e, ok := r.float.drag.Update(gtx.Metric, gtx.Source, gesture.Axis(r.resize.axis))
	if ok {
		if e.Kind == pointer.Drag {
			de = &e
		}
	}

	var posDifference float32
	if de != nil {

		// 记录拖动前的宽度用于比较
		prevWidth := r.pos
		if r.prev != nil {
			prevWidth -= r.prev.pos
		}

		posDifference = de.Position.X
		if r.resize.axis == layout.Vertical {
			posDifference = de.Position.Y
		}

		if posDifference < 0 {
			for curr := r; curr != nil; curr = curr.prev {
				curr.float.pos += int(posDifference)
				minPos := r.resize.minLength
				if curr.prev != nil {
					minPos = curr.prev.pos + curr.resize.minLength
				}
				if curr.float.pos < minPos {
					curr.float.pos = minPos
				} else {
					break
				}
			}
		}
		if posDifference > 0 {
			for curr := r; curr != nil; curr = curr.next {
				curr.float.pos += int(posDifference)
				maxPos := r.resize.length
				if curr.next != nil {
					maxPos = curr.next.pos - curr.resize.minLength
				}
				if curr.float.pos > maxPos {
					curr.float.pos = maxPos
				} else {
					break
				}
			}
		}
		// 计算新宽度并触发回调
		newWidth := r.pos
		if r.prev != nil {
			newWidth -= r.prev.pos
		}

		// 当宽度变化且回调函数存在时触发
		if newWidth != prevWidth && r.OnResize != nil {
			r.OnResize(r.Index, newWidth)
		}
	}

	rect := image.Rectangle{Max: dims.Size}
	defer clip.Rect(rect).Push(gtx.Ops).Pop()
	r.float.drag.Add(gtx.Ops)
	cursor := pointer.CursorRowResize
	if r.resize.axis == layout.Horizontal {
		cursor = pointer.CursorColResize
	}
	cursor.Add(gtx.Ops)

	return layout.Dimensions{Size: dims.Size}
}

func (r *Resize) CustomResizeHandleBar(gtx layout.Context) layout.Dimensions {
	x := gtx.Dp(unit.Dp(4))
	y := gtx.Constraints.Max.Y
	if r.axis == layout.Vertical {
		x = gtx.Constraints.Max.X
		y = gtx.Dp(unit.Dp(4))
	}
	rect := image.Rectangle{
		Max: image.Point{
			X: x,
			Y: y,
		},
	}
	paint.FillShape(gtx.Ops, colors.DividerFg, clip.Rect(rect).Op())
	return layout.Dimensions{Size: rect.Max}
}
