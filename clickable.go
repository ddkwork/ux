package ux

import (
	"image"
	"image/color"

	"gioui.org/io/input"
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

type Clickable struct {
	bgColor      color.NRGBA
	bgColorHover color.NRGBA
	clickable    widget.Clickable
	onClick      func()
	widget       layout.Widget
}

func NewClickable() *Clickable {
	return &Clickable{}
}

func (c *Clickable) SetBgColor(bgColor color.NRGBA) {
	c.bgColor = bgColor
}

func (c *Clickable) SetBgColorHover(bgColorHover color.NRGBA) {
	c.bgColorHover = bgColorHover
}

func (c *Clickable) SetWidget(widget layout.Widget) *Clickable {
	c.widget = widget
	return c
}

func (c *Clickable) SetOnClick(onClick func()) {
	c.onClick = onClick
}

func (c *Clickable) Layout(gtx layout.Context) layout.Dimensions {
	if c.bgColor == (color.NRGBA{}) {
		c.bgColor = th.Color.DefaultBgGrayColor
	}
	if c.bgColorHover == (color.NRGBA{}) {
		c.bgColorHover = th.Color.HoveredBorderBlueColor
	}
	for c.clickable.Clicked(gtx) {
		if c.onClick != nil {
			c.onClick()
		}
	}
	return c.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		semantic.Button.Add(gtx.Ops)
		return layout.Background{}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				rect := clip.UniformRRect(image.Rectangle{Max: image.Point{
					X: gtx.Constraints.Min.X,
					Y: gtx.Constraints.Min.Y,
				}}, gtx.Dp(th.Size.DefaultElementRadiusSize))
				defer rect.Push(gtx.Ops).Pop()
				if gtx.Source == (input.Source{}) {
					paint.Fill(gtx.Ops, Disabled(c.bgColorHover))
				} else if c.clickable.Hovered() {
					// paint.Fill(gtx.Ops, c.bgColorHover)
				}
				if gtx.Focused(c.clickable) {
					paint.Fill(gtx.Ops, c.bgColorHover)
				}
				return layout.Dimensions{Size: gtx.Constraints.Min}
			},
			c.widget,
		)
	})
}
