package ux

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type Card_ struct {
	radius  int
	padding unit.Dp
	bgColor color.NRGBA
	content Widget
}

func NewCard(content Widget) *Card_ {
	return &Card_{
		radius:  15,
		padding: unit.Dp(20),
		bgColor: th.Color.CardBgColor,
		content: content,
	}
}

func (c *Card_) SetRadius(radius int) *Card_ {
	c.radius = radius
	return c
}

func (c *Card_) SetBgColor(bgColor color.NRGBA) *Card_ {
	c.bgColor = bgColor
	return c
}

func (c *Card_) SetPadding(padding unit.Dp) *Card_ {
	c.padding = padding
	return c
}

func fill(gtx layout.Context, color color.NRGBA) layout.Dimensions {
	cs := gtx.Constraints
	d := image.Point{X: cs.Max.X, Y: cs.Min.Y}
	track := image.Rectangle{
		Max: d,
	}
	defer clip.Rect(track).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, color)
	return layout.Dimensions{Size: d}
}

func (c *Card_) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			rect := clip.UniformRRect(image.Rectangle{Max: image.Point{
				X: gtx.Constraints.Max.X,
				Y: gtx.Constraints.Min.Y,
			}}, c.radius)
			defer rect.Push(gtx.Ops).Pop()
			return fill(gtx, c.bgColor)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(c.padding).Layout(gtx, layout.Widget(c.content))
		}),
	)
}
