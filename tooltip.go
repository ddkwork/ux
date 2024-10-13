package ux

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

type (
	TipIconButton struct {
		*widget.Clickable
		material.IconButtonStyle
		component.Tooltip
		component.TipArea
		callback func()
	}
)

func NewTooltipButton(icon *widget.Icon, tip string, callback func()) *TipIconButton {
	t := component.PlatformTooltip(th.Theme, tip)
	t.Bg = Yellow100
	t.Text.Color = th.Color.ButtonTextBlackColor
	t.CornerRadius = 14
	t.Text.MaxLines = 3 //todo newlines
	clickable := &widget.Clickable{}
	iconButtonStyle := material.IconButton(th.Theme, clickable, icon, "")
	iconButtonStyle.Color = th.Fg
	iconButtonStyle.Background = th.Color.InputFocusedBgColor
	iconButtonStyle.Inset = layout.UniformInset(unit.Dp(6))
	b := &TipIconButton{
		Clickable:       clickable,
		IconButtonStyle: iconButtonStyle,
		Tooltip:         t,
		TipArea:         component.TipArea{},
		callback:        callback,
	}
	return b
}

func (t *TipIconButton) Layout(gtx C) D {
	if t.Clickable.Clicked(gtx) {
		if t.callback != nil {
			t.callback()
		}
	}
	//return t.TipArea.Layout(gtx, t.Tooltip, t.IconButtonStyle.Layout)
	return layout.Inset{Top: 4}.Layout(gtx, func(gtx C) D {
		return t.TipArea.Layout(gtx, t.Tooltip, t.IconButtonStyle.Layout)
	})
}
