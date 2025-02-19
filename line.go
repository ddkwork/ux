package ux

import (
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type Line struct {
	width      float32 // 线条宽度
	dashLength float32 // 选线中间的虚线长度
	sap        float32
	color      color.NRGBA
}

func NewLine() *Line {
	return &Line{
		width:      2,
		sap:        5,
		dashLength: 5,
		color:      th.Color.DefaultLineColor,
	}
}

// Color 设置color
func (l *Line) Color(color color.NRGBA) *Line {
	l.color = color
	return l
}

func (l *Line) Width(width float32) *Line {
	l.width = width
	return l
}

func (l *Line) Sap(sap float32) *Line {
	l.sap = sap
	return l
}

func (l *Line) DashLength(dashLength float32) *Line {
	l.dashLength = dashLength
	return l
}

func (l *Line) getPoint(a, b f32.Point, t float32) f32.Point {
	return f32.Pt(
		a.X+(b.X-a.X)*t,
		a.Y+(b.Y-a.Y)*t,
	)
}

// DashedLine 虚线
func (l *Line) DashedLine(gtx layout.Context, start, end f32.Point) *Line {
	ops := gtx.Ops
	var path clip.Path
	path.Begin(ops)
	dx := end.X - start.X
	dy := end.Y - start.Y
	length := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	// Draw dashes
	for i := float32(0); i < length; i += l.dashLength + l.sap {
		dashEnd := i + l.dashLength
		if dashEnd > length {
			dashEnd = length
		}
		// Calculate start and end points of the current dash
		t0 := i / length
		t1 := dashEnd / length
		p0 := l.getPoint(start, end, t0)
		p1 := l.getPoint(start, end, t1)
		path.MoveTo(p0)
		path.LineTo(p1)
	}
	// Create a stroke style and fill the path with the specified color
	stroke := clip.Stroke{
		Path:  path.End(),
		Width: l.width,
	}.Op()
	paint.FillShape(ops, l.color, stroke)
	return l
}

// Line 普通直线
func (l *Line) Line(gtx layout.Context, start f32.Point, end f32.Point) *Line {
	var path clip.Path
	path.Begin(gtx.Ops)
	path.MoveTo(start) // 起始点
	path.LineTo(end)   // 结束点
	paint.FillShape(gtx.Ops, l.color, clip.Stroke{Path: path.End(), Width: l.width}.Op())
	return l
}

func (l *Line) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{}
}
