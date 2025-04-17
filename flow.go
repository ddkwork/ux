package ux

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"github.com/ddkwork/ux/widget/material"
)

// inf is an infinite axis constraint.
const inf = 1e6

// FlowElement lays out the ith element of a Grid.
type FlowElement func(gtx layout.Context, i int) layout.Dimensions

// Flow lays out at most Num filterMap along the main axis.
// The number of cross axis filterMap depend on the total number of filterMap.
type Flow struct {
	Num       int
	Axis      layout.Axis
	Alignment layout.Alignment
	list      *widget.List
}

type wrapData struct {
	dims layout.Dimensions
	call op.CallOp
}

func (g *Flow) Layout(gtx layout.Context, num int, el FlowElement) layout.Dimensions {
	if g.Num == 0 {
		return layout.Dimensions{Size: gtx.Constraints.Min}
	}
	if g.Axis == g.list.Axis {
		if g.Axis == layout.Horizontal {
			g.list.Axis = layout.Vertical
		} else {
			g.list.Axis = layout.Horizontal
		}
		g.list.Alignment = g.Alignment
	}
	csMax := gtx.Constraints.Max
	return material.List(th, g.list).Layout(gtx, (num+g.Num-1)/g.Num, func(gtx layout.Context, idx int) layout.Dimensions {
		if g.Axis == layout.Horizontal {
			gtx.Constraints.Max.Y = inf
		} else {
			gtx.Constraints.Max.X = inf
		}
		gtx.Constraints.Min = image.Point{}
		var mainMax, crossMax int
		left := axisMain(g.Axis, csMax)
		i := idx * g.Num
		n := min(num, i+g.Num)
		for ; i < n; i++ {
			dims := el(gtx, i)
			main := axisMain(g.Axis, dims.Size)
			crossMax = max(crossMax, axisCross(g.Axis, dims.Size))
			left -= main
			if left <= 0 {
				mainMax = axisMain(g.Axis, csMax)
				break
			}
			pt := axisPoint(g.Axis, main, 0)
			op.Offset(pt).Add(gtx.Ops)
			mainMax += main
		}
		return layout.Dimensions{Size: axisPoint(g.Axis, mainMax, crossMax)}
	})
}

func axisPoint(a layout.Axis, main, cross int) image.Point {
	if a == layout.Horizontal {
		return image.Point{main, cross}
	} else {
		return image.Point{cross, main}
	}
}

func axisMain(a layout.Axis, sz image.Point) int {
	if a == layout.Horizontal {
		return sz.X
	} else {
		return sz.Y
	}
}

func axisCross(a layout.Axis, sz image.Point) int {
	if a == layout.Horizontal {
		return sz.Y
	} else {
		return sz.X
	}
}
