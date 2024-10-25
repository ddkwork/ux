package ux

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type RadioButton struct {
	radioButton *widget.Bool
	key         string
	label       string
	group       *widget.Enum
	iconSize    unit.Dp
	textSize    unit.Sp
}

// NewRadioButton returns a RadioButton with a label. The key specifies
// the value for the EnumTypes.
func NewRadioButton(group *widget.Enum, key, label string) *RadioButton {
	r := &RadioButton{
		radioButton: &widget.Bool{Value: true},
		group:       group,
		key:         key,
		label:       label,
		iconSize:    th.Size.DefaultIconSize,
		textSize:    th.Size.DefaultTextSize,
	}
	return r
}

func (r *RadioButton) SetSize(size ElementStyle) {
	r.iconSize = size.IconSize
	r.textSize = size.TextSize
}

// Layout updates enum and displays the radio button.
func (r *RadioButton) Layout(gtx layout.Context) layout.Dimensions {
	iconColor := th.Color.BorderLightGrayColor
	if r.group.Value == r.key {
		iconColor = th.Color.RadioSelectBgColor
	}
	return r.radioButton.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if r.radioButton.Hovered() {
			iconColor = th.Color.RadioSelectBgColor
		}
		rb := material.RadioButton(th.Theme, r.group, r.key, r.label)
		rb.IconColor = iconColor
		rb.Color = th.Color.DefaultTextWhiteColor
		rb.Size = r.iconSize
		rb.TextSize = r.textSize
		return rb.Layout(gtx)
	})
}
