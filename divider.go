package ux

import (
	"image"
	"image/color"

	"github.com/ddkwork/ux/resources/colors"

	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"gioui.org/layout"
	"gioui.org/unit"
)

type DividerStyle struct {
	Thickness unit.Dp
	Fill      color.NRGBA
	Inset     layout.Inset
	Axis      layout.Axis
}

func (d *DividerStyle) Layout(gtx layout.Context) layout.Dimensions {
	if (d.Axis == layout.Horizontal && gtx.Constraints.Min.X == 0) ||
		(d.Axis == layout.Vertical && gtx.Constraints.Min.Y == 0) {
		return D{}
	}

	if d.Fill == (color.NRGBA{}) {
		d.Fill = colors.DividerFg
	}

	return d.Inset.Layout(gtx, func(gtx C) D {
		weight := gtx.Dp(d.Thickness)
		line := image.Rectangle{Max: image.Pt(weight, gtx.Constraints.Min.Y)}
		if d.Axis == layout.Horizontal {
			line = image.Rectangle{Max: image.Pt(gtx.Constraints.Min.X, weight)}
		}
		paint.FillShape(gtx.Ops, d.Fill, clip.Rect(line).Op())
		return D{Size: line.Max}
	})
}

func Divider(axis layout.Axis, thickness unit.Dp) *DividerStyle {
	return &DividerStyle{
		Thickness: thickness,
		Axis:      axis,
	}
}
