package ux

//
//import (
//	"gioui.org/layout"
//	"gioui.org/unit"
//	"gioui.org/widget"
//	"github.com/ddkwork/ux/widget/material"
//)
//
//type LabeledInput struct {
//	Label          string
//	SpaceBetween   int
//	MinEditorWidth unit.Dp
//	MinLabelWidth  unit.Dp
//	Editor         *PatternEditor
//	Hint           string
//}
//
//func (l *LabeledInput) SetText(text string) {
//	l.Editor.SetText(text)
//}
//
//func (l *LabeledInput) Value() string {
//	return l.Editor.Value()
//}
//
//func (l *LabeledInput) SetHint(hint string) {
//	l.Hint = hint
//}
//
//func (l *LabeledInput) SetLabel(label string) {
//	l.Label = label
//}
//
//func (l *LabeledInput) SetOnChanged(f func(text string)) {
//	l.Editor.SetOnChanged(f)
//}
//
//func (l *LabeledInput) Layout(gtx layout.Context) layout.Dimensions {
//	return layout.Flex{
//		Axis:      layout.Horizontal,
//		Alignment: layout.Middle,
//	}.Layout(gtx,
//		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
//			return layout.Inset{Right: unit.Dp(l.SpaceBetween)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
//				gtx.Constraints.Min.X = gtx.Dp(l.MinLabelWidth)
//				return material.Label(th, th.TextSize, l.Label).Layout(gtx)
//			})
//		}),
//		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
//			gtx.Constraints.Min.X = gtx.Dp(l.MinEditorWidth)
//			return widget.Border{
//				// Color:        th.BorderColor,
//				Width:        unit.Dp(1),
//				CornerRadius: unit.Dp(4),
//			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
//				return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
//					l.Editor.SetHint(l.Hint)
//					return l.Editor.Layout(gtx)
//				})
//			})
//		}),
//	)
//}
