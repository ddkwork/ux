package ux

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/ddkwork/ux/widget/material"
	"github.com/ddkwork/ux/x/component"
)

type InputModal struct {
	textField *TextFieldFuzz
	addBtn    widget.Clickable
	closeBtn  widget.Clickable

	Title string

	onClose func()
	onAdd   func(text string)
}

func NewInputModal(title, placeholder string) *InputModal {
	ed := NewTextFieldFuzz("", placeholder)
	ed.SetIcon(FileFolderIcon, true)
	return &InputModal{
		textField: ed,
		Title:     title,
	}
}

func (i *InputModal) SetOnClose(f func()) {
	i.onClose = f
}

func (i *InputModal) SetOnAdd(f func(text string)) {
	i.onAdd = f
}

func (i *InputModal) SetText(text string) {
	i.textField.SetText(text)
}

func (i *InputModal) layout(gtx layout.Context) layout.Dimensions {
	if i.onClose != nil && i.closeBtn.Clicked(gtx) {
		i.onClose()
	}

	if i.onAdd != nil && i.addBtn.Clicked(gtx) {
		i.onAdd(i.textField.GetText())
	}

	border := widget.Border{
		// Color:        theme.TableBorderColor,//todo
		CornerRadius: unit.Dp(4),
		Width:        unit.Dp(1),
	}

	return layout.N.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(80)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.X = gtx.Dp(500)
				gtx.Constraints.Max.Y = gtx.Dp(180)

				return component.NewModalSheet(component.NewModal()).Layout(gtx, th.Theme, &component.VisibilityAnimation{}, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.Label(th.Theme, unit.Sp(14), i.Title).Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return i.textField.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceStart}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										closeBtn := Button(&i.closeBtn, NavigationCloseIcon, "Close")
										// closeBtn.Color = theme.ButtonTextColor//todo
										return closeBtn.Layout(gtx)
									}),
									layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										addBtn := Button(&i.addBtn, ContentAddIcon, "Add")
										// addBtn.Color = theme.ButtonTextColor
										// addBtn.Background = theme.SendButtonBgColor
										return addBtn.Layout(gtx)
									}),
								)
							}),
						)
					})
				})
			})
		})
	})
}

func (i *InputModal) Layout(gtx layout.Context) layout.Dimensions {
	ops := op.Record(gtx.Ops)
	dims := i.layout(gtx)
	defer op.Defer(gtx.Ops, ops.Stop())
	return dims
}
