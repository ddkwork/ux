package ux

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func LayoutErrorLabel(gtx C, e error) D {
	if e != nil {
		return layout.Inset{
			Top:    unit.Dp(10),
			Bottom: unit.Dp(10),
			Left:   unit.Dp(15),
			Right:  unit.Dp(15),
		}.Layout(gtx, func(gtx C) D {
			label := material.Label(th.Theme, th.TextSize*0.8, e.Error())
			label.Color = color.NRGBA{R: 255, A: 255}
			label.Alignment = text.Middle
			return label.Layout(gtx)
		})
	} else {
		return layout.Dimensions{}
	}
}
