package ux

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/std/safemap"
	"github.com/ddkwork/ux/widget/material"
)

// inf is an infinite axis constraint.
const inf = 1e6

// Flow lays out at most rowElemCount filterMap along the main axis.
// The number of cross axis filterMap depend on the total number of filterMap.
type Flow struct {
	rowElemCount int
	axis         layout.Axis
	alignment    layout.Alignment
	list         widget.List
	elems        *safemap.M[*ButtonAnimation, *ContextMenu]
}

func NewFlow(rowElemCount int, elems *safemap.M[*ButtonAnimation, *ContextMenu]) *Flow {
	if rowElemCount == 0 {
		rowElemCount = 5
	}
	return &Flow{
		rowElemCount: rowElemCount,
		axis:         layout.Horizontal,
		alignment:    layout.Middle,
		list:         widget.List{},
		elems:        elems,
	}
}

func (g *Flow) Layout(gtx layout.Context) layout.Dimensions {
	if g.rowElemCount == 0 {
		return layout.Dimensions{Size: gtx.Constraints.Min}
	}
	if g.axis == g.list.Axis {
		if g.axis == layout.Horizontal {
			g.list.Axis = layout.Vertical
		} else {
			g.list.Axis = layout.Horizontal
		}
		g.list.Alignment = g.alignment
	}
	csMax := gtx.Constraints.Max
	return material.List(th, &g.list).Layout(gtx, (g.elems.Len()+g.rowElemCount-1)/g.rowElemCount, func(gtx layout.Context, index int) layout.Dimensions {
		if g.axis == layout.Horizontal {
			gtx.Constraints.Max.Y = inf
		} else {
			gtx.Constraints.Max.X = inf
		}
		gtx.Constraints.Min = image.Point{}
		var mainMax, crossMax int
		left := axisMain(g.axis, csMax)
		i := index * g.rowElemCount
		n := min(g.elems.Len(), i+g.rowElemCount)
		for ; i < n; i++ {
			gtx.Constraints.Min.X = 300 // todo 定义常量 导航按钮宽度
			gtx.Constraints.Max.X = gtx.Constraints.Min.X

			dims := g.elems.Values()[i].Layout(gtx, []layout.Widget{
				func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(2).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return g.elems.Keys()[i].Layout(gtx)
					})
				},
			})

			main := axisMain(g.axis, dims.Size)
			crossMax = max(crossMax, axisCross(g.axis, dims.Size))
			left -= main
			if left <= 0 {
				mainMax = axisMain(g.axis, csMax)
				break
			}
			pt := axisPoint(g.axis, main, 0)
			op.Offset(pt).Add(gtx.Ops)
			mainMax += main
		}
		return layout.Dimensions{Size: axisPoint(g.axis, mainMax, crossMax)}
	})
}

func axisPoint(a layout.Axis, main, cross int) image.Point {
	if a == layout.Horizontal {
		return image.Point{X: main, Y: cross}
	} else {
		return image.Point{X: cross, Y: main}
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
