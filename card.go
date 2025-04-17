package ux

import (
	"image"
	"image/color"

	"github.com/ddkwork/ux/widget/material"
	"github.com/ddkwork/ux/x/component"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/x/outlay"
)

type (
	Card struct {
		Name      string
		Flag      image.Image
		btn       widget.Clickable
		container component.SurfaceStyle
		shadow    component.ShadowStyle

		IsSearchedFor     bool
		IsActiveContinent bool
		Selected          bool

		menu            component.MenuState
		ctxArea         component.ContextArea
		isMenuTriggered bool

		// menu options
		selectBtn     widget.Clickable
		deselectBtn   widget.Clickable
		copyToClipBtn widget.Clickable
	}
)

func (c *Card) LayCard(gtx C) D {
	size := image.Pt(gtx.Dp(unit.Dp(float32(200))), gtx.Dp(unit.Dp(float32(250))))

	c.container.Theme = th
	c.container.Theme.Bg = BackgroundColor
	c.container.Elevation = unit.Dp(5)
	c.shadow.CornerRadius = unit.Dp(18)
	c.shadow.Elevation = unit.Dp(8)
	c.shadow.AmbientColor = ColorPink // color.NRGBA{A: 0x10}
	c.shadow.PenumbraColor = color.NRGBA{A: 0x20}
	c.shadow.UmbraColor = color.NRGBA{A: 0x30}

	if !c.isMenuTriggered {
		lbl := "Select"
		btn := &c.selectBtn
		if c.Selected {
			lbl = "Deselect"
			btn = &c.deselectBtn
		}
		var item component.MenuItemStyle
		item.LabelInset = outlay.Inset{
			Top: unit.Dp(5),
			// Right:  unit.Dp(5),
			Bottom: unit.Dp(5),
			// Left:   unit.Dp(5),
		}
		item = component.MenuItem(th, btn, lbl)

		c.menu = component.MenuState{
			Options: []func(gtx C) D{
				item.Layout,
				component.MenuItem(th, &c.copyToClipBtn, "Copy as JSON").Layout,
			},
		}
	}

	if c.Selected {
		c.shadow.AmbientColor = color.NRGBA{G: 255, A: 85}
		c.shadow.PenumbraColor = color.NRGBA{G: 255, A: 170}
		c.shadow.UmbraColor = color.NRGBA{G: 255, A: 255}
	}
	if c.btn.Hovered() {
		if c.Selected {
			c.shadow.AmbientColor = color.NRGBA{R: 255, G: 107, B: 108, A: 85}
			c.shadow.PenumbraColor = color.NRGBA{R: 255, G: 107, B: 108, A: 170}
			c.shadow.UmbraColor = color.NRGBA{R: 255, G: 107, B: 108, A: 255}
		} else {
			c.shadow.AmbientColor = color.NRGBA{R: 233, G: 255, B: 219, A: 85}
			c.shadow.PenumbraColor = color.NRGBA{R: 233, G: 255, B: 219, A: 170}
			c.shadow.UmbraColor = color.NRGBA{R: 233, G: 255, B: 219, A: 255}
		}
	}

	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			c.container.ShadowStyle = c.shadow

			return c.container.Layout(gtx, func(gtx C) D {
				return material.Clickable(gtx, &c.btn, func(gtx C) D {
					gtx.Constraints = layout.Exact(gtx.Constraints.Constrain(size))
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(gtx,

							// country name
							layout.Rigid(func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Flexed(1, func(gtx C) D {
										return layout.Center.Layout(gtx, func(gtx C) D {
											return material.Body2(th, c.Name).Layout(gtx)
										})
									}),
								)
							}),

							// country flag
							layout.Rigid(func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Flexed(1, func(gtx C) D {
										return layout.Center.Layout(gtx, func(gtx C) D {
											var flag D
											if c.Flag == nil {
												flag = material.Loader(th).Layout(gtx)
											} else {
												flag = widget.Image{
													Src: paint.NewImageOp(c.Flag),
													Fit: widget.ScaleDown,
												}.Layout(gtx)
											}
											return flag
										})
									}))
							}))
					})
				})
			})
		}),
		layout.Expanded(func(gtx C) D {
			return c.ctxArea.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min = image.Point{}
				return component.Menu(th, &c.menu).Layout(gtx)
			})
		}))
}
