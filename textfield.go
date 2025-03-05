package ux

import (
	"image"
	"image/color"

	"github.com/ddkwork/ux/widget/material"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type TextFieldFuzz struct {
	textEditor widget.Editor
	Icon       *widget.Icon
	iconClick  widget.Clickable

	iconPositionStart bool

	Text        string
	Placeholder string

	size image.Point

	onIconClick  func()
	onTextChange func(text string)
	borderColor  color.NRGBA
}

func NewTextFieldFuzz(text, placeholder string) *TextFieldFuzz {
	t := &TextFieldFuzz{
		textEditor:  widget.Editor{},
		Text:        text,
		Placeholder: placeholder,
	}

	t.textEditor.SetText(text)
	t.textEditor.SingleLine = true
	return t
}

func (t *TextFieldFuzz) GetText() string {
	return t.textEditor.Text()
}

func (t *TextFieldFuzz) SetText(text string) {
	t.textEditor.SetText(text)
}

func (t *TextFieldFuzz) SetIcon(icon *widget.Icon, iconPositionStart bool) {
	t.Icon = icon
	t.iconPositionStart = iconPositionStart
}

func (t *TextFieldFuzz) SetMinWidth(width int) {
	t.size.X = width
}

func (t *TextFieldFuzz) SetBorderColor(color color.NRGBA) {
	t.borderColor = color
}

func (t *TextFieldFuzz) SetOnTextChange(f func(text string)) {
	t.onTextChange = f
}

func (t *TextFieldFuzz) SetOnIconClick(f func()) {
	t.onIconClick = f
}

func (t *TextFieldFuzz) Layout(gtx layout.Context) layout.Dimensions {
	for {
		event, ok := gtx.Event(key.FocusFilter{Target: t}, key.Filter{Name: key.NameEscape})
		if !ok {
			break
		}
		switch ev := event.(type) {
		case key.FocusEvent:
			gtx.Execute(key.FocusCmd{Tag: &t.textEditor})
		case key.Event:
			if ev.Name == key.NameEscape {
				gtx.Execute(key.FocusCmd{Tag: nil})
			}
		}
	}

	borderColor := th.InputActivatedBorderColor // theme.BorderColor  //todo
	if gtx.Source.Focused(&t.textEditor) {
		borderColor = th.FocusedBgColor // theme.BorderColorFocused
	}

	cornerRadius := unit.Dp(4)
	border := widget.Border{
		Color:        borderColor,
		Width:        unit.Dp(1),
		CornerRadius: cornerRadius,
	}

	leftPadding := unit.Dp(8)
	if t.Icon != nil && t.iconPositionStart {
		leftPadding = unit.Dp(0)
	}

	for {
		event, ok := t.textEditor.Update(gtx)
		if !ok {
			break
		}
		if _, ok := event.(widget.ChangeEvent); ok {
			if t.onTextChange != nil {
				t.onTextChange(t.textEditor.Text())
			}
		}
	}

	return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if t.size.X == 0 {
			t.size.X = gtx.Constraints.Min.X
		}

		gtx.Constraints.Min = t.size
		return layout.Inset{
			Top:    4,
			Bottom: 4,
			Left:   leftPadding,
			Right:  4,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			inputLayout := layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return material.Editor(th.Theme, &t.textEditor, t.Placeholder).Layout(gtx)
			})
			widgets := []layout.FlexChild{inputLayout}

			spacing := layout.SpaceBetween
			if t.Icon != nil {
				iconLayout := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					clk := &widget.Clickable{}
					if t.onIconClick != nil {
						clk = &t.iconClick
						if t.iconClick.Clicked(gtx) {
							t.onIconClick()
						}
					}
					b := Button(clk, t.Icon, "")
					b.Inset = layout.Inset{Left: unit.Dp(8), Right: unit.Dp(2), Top: unit.Dp(2), Bottom: unit.Dp(2)}
					return b.Layout(gtx)
				})

				if !t.iconPositionStart {
					widgets = []layout.FlexChild{inputLayout, iconLayout}
				} else {
					widgets = []layout.FlexChild{iconLayout, inputLayout}
					spacing = layout.SpaceEnd
				}
			}

			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: spacing}.Layout(gtx, widgets...)
		})
	})
}
