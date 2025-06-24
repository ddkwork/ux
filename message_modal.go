package ux

import (
	"image/color"

	"github.com/ddkwork/ux/widget/material"
	"github.com/ddkwork/ux/x/component"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
)

const (
	MessageModalTypeInfo = "info"
	MessageModalTypeWarn = "warn"
	MessageModalTypeErr  = "err"
)

type MessageModal struct {
	Title string
	Body  string
	Type  string

	Visible bool

	options  []ModalOption
	onSubmit OnModalSubmit
}

type ModalOption struct {
	Text   string
	Button widget.Clickable
	Icon   []byte
}

type OnModalSubmit func(selectedOption string)

func NewMessageModal(title, body, modalType string, onSubmit OnModalSubmit, options ...ModalOption) *MessageModal {
	return &MessageModal{
		Title:    title,
		Body:     body,
		Type:     modalType,
		onSubmit: onSubmit,

		options: options,
	}
}

func (modal *MessageModal) Show() {
	modal.Visible = true
}

func (modal *MessageModal) Hide() {
	modal.Visible = false
}

func (modal *MessageModal) layout(gtx layout.Context) layout.Dimensions {
	borderColor := th.BorderBlueColor // theme.TableBorderColor //todo
	switch modal.Type {
	case MessageModalTypeErr:
		borderColor = color.NRGBA{R: 0xD1, G: 0x1E, B: 0x35, A: 0xFF}
	case MessageModalTypeInfo:
		borderColor = color.NRGBA{R: 0x1D, G: 0xBF, B: 0xEC, A: 0xFF}
	case MessageModalTypeWarn:
		borderColor = color.NRGBA{R: 0xFD, G: 0xB5, B: 0x0E, A: 0xFF}
	}

	border := widget.Border{
		Color:        borderColor,
		CornerRadius: unit.Dp(18),
		Width:        unit.Dp(1),
	}

	return layout.N.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(80)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.X = gtx.Dp(500)
				gtx.Constraints.Max.Y = gtx.Dp(180)

				return component.NewModalSheet(component.NewModal()).Layout(gtx, th, &component.VisibilityAnimation{}, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.Label(th, unit.Sp(14), modal.Title).Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.Body1(th, modal.Body).Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								count := len(modal.options)
								items := make([]layout.FlexChild, 0, count)
								for i := range modal.options {
									if modal.onSubmit != nil {
										if modal.options[i].Button.Clicked(gtx) {
											modal.onSubmit(modal.options[i].Text)
										}
									}

									items = append(items, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return Button(&modal.options[i].Button, nil, modal.options[i].Text).Layout(gtx)
									}),
										layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
									)
								}
								return layout.Inset{Top: unit.Dp(5)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{
										Axis:      layout.Horizontal,
										Alignment: layout.Middle,
										Spacing:   layout.SpaceStart,
									}.Layout(gtx,
										items...,
									)
								})
							}),
						)
					})
				})
			})
		})
	})
}

func (modal *MessageModal) Layout(gtx layout.Context) layout.Dimensions {
	if modal == nil || !modal.Visible {
		return layout.Dimensions{}
	}

	ops := op.Record(gtx.Ops)
	dims := modal.layout(gtx)
	defer op.Defer(gtx.Ops, ops.Stop())

	return dims
}
