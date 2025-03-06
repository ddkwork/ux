package ux

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/x/outlay"
	"github.com/ddkwork/ux/widget/material"
	"github.com/ddkwork/ux/x/component"
)

type InputModal struct {
	Title    string
	widget   layout.Widget
	applyBtn widget.Clickable
	closeBtn widget.Clickable
	onApply  func()
	onClose  func()
}

func NewInputModal(title string, w layout.Widget) *InputModal {
	return &InputModal{
		widget: w,
		Title:  title,
	}
}

func (i *InputModal) SetOnClose(f func()) { i.onClose = f }
func (i *InputModal) SetOnApply(f func()) { i.onApply = f }
func (i *InputModal) Layout(gtx layout.Context) layout.Dimensions {
	ops := op.Record(gtx.Ops)
	dims := i.layout(gtx)
	defer op.Defer(gtx.Ops, ops.Stop())
	return dims
}

func (i *InputModal) layout(gtx layout.Context) layout.Dimensions {
	if i.closeBtn.Clicked(gtx) {
		i.onClose()
	}

	if i.onApply != nil && i.applyBtn.Clicked(gtx) {
		i.onApply()
	}

	border := widget.Border{
		Color:        th.BorderBlueColor,
		CornerRadius: unit.Dp(4),
		Width:        unit.Dp(1),
	}

	return layout.N.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(80)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				//gtx.Constraints.Max.X = gtx.Dp(500)
				//gtx.Constraints.Max.Y = gtx.Dp(180)

				return component.NewModalSheet(component.NewModal()).Layout(gtx, th.Theme, &component.VisibilityAnimation{}, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.Label(th.Theme, unit.Sp(14), i.Title).Layout(gtx)
							}),
							outlay.EmptyRigidVertical(20),
							//layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
							layout.Rigid(i.widget),
							outlay.EmptyRigidVertical(20),
							//layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceStart}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										closeBtn := Button(&i.closeBtn, NavigationCloseIcon, "Close")
										// closeBtn.Color = theme.ButtonTextColor//todo
										return closeBtn.Layout(gtx)
									}),
									layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										applyBtn := Button(&i.applyBtn, ActionAssignmentTurnedInIcon, "Apply")
										// applyBtn.Color = theme.ButtonTextColor
										// applyBtn.Background = theme.SendButtonBgColor
										return applyBtn.Layout(gtx)
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
