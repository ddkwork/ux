package ux

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Modal struct {
	visible       bool
	content       layout.Widget
	closeButton   *Button
	clickerWidget *Clickable
	title         string
	height        int
	width         int
}

func NewModal() *Modal {
	m := &Modal{
		height:        300,
		width:         500,
		title:         "modal",
		clickerWidget: NewClickable(),
	}
	m.closeButton = NewButton("", func() {
		m.visible = false
	}).SetIcon(IconClose)
	return m
}

func (m *Modal) SetWidth(width int) *Modal {
	m.width = width
	return m
}

func (m *Modal) SetHeight(height int) *Modal {
	m.height = height
	return m
}

func (m *Modal) Visible() bool {
	return m.visible
}

func (m *Modal) SetTitle(title string) *Modal {
	m.title = title
	return m
}

func (m *Modal) SetContent(content layout.Widget) {
	m.content = content
	m.visible = true
}

func (m *Modal) Close() {
	m.visible = false
}

func (m *Modal) Layout(gtx layout.Context) layout.Dimensions {
	if !m.visible {
		return layout.Dimensions{}
	}
	if m.visible {
		// 绘制全屏半透明遮罩层
		paint.Fill(gtx.Ops, th.Color.DefaultMaskBgColor)
	}
	width := gtx.Dp(unit.Dp(m.width))
	height := gtx.Dp(unit.Dp(m.height))
	return m.clickerWidget.SetWidget(func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// Set the size of the Modal
				gtx.Constraints = layout.Exact(image.Point{X: width, Y: height})
				rc := clip.RRect{
					Rect: image.Rectangle{Max: image.Point{
						X: gtx.Constraints.Min.X,
						Y: gtx.Constraints.Min.Y,
					}},
					NW: 20, NE: 20, SE: 20, SW: 20,
				}
				paint.FillShape(gtx.Ops, th.Color.DefaultContentBgGrayColor, rc.Op(gtx.Ops))
				// Center the text inside the Modal
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{Left: 0, Right: 10, Bottom: 10, Top: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return layout.Stack{Alignment: layout.N}.Layout(gtx,
									layout.Stacked(func(gtx layout.Context) layout.Dimensions {
										label := material.Label(th.Theme, unit.Sp(16), m.title)
										label.Color = th.Color.DefaultTextWhiteColor
										return label.Layout(gtx)
									}),
									layout.Expanded(func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{Left: unit.Dp(m.width - 30)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return m.closeButton.Layout(gtx)
										})
									}),
								)
							})
						})
					}),
					DrawLineFlex(th.Color.DefaultLineColor, unit.Dp(1), unit.Dp(m.width)),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{Left: 30, Right: 30, Bottom: 30, Top: 30}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return m.content(gtx)
							})
						})
					}),
				)
			})
		})
	}).Layout(gtx)
}

func DrawLineFlex(background color.NRGBA, height, width unit.Dp) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return DrawLine(gtx, background, height, width)
	})
}

func DrawLine(gtx layout.Context, background color.NRGBA, height, width unit.Dp) layout.Dimensions {
	w, h := gtx.Dp(width), gtx.Dp(height)
	tabRect := image.Rect(0, 0, w, h)
	paint.FillShape(gtx.Ops, background, clip.Rect(tabRect).Op())
	return layout.Dimensions{Size: image.Pt(w, h)}
}
