package ux

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/ddkwork/ux/internal/animation"
	"github.com/ddkwork/ux/internal/animation/gween"
	"github.com/ddkwork/ux/internal/animation/gween/ease"
	"github.com/ddkwork/ux/resources/colors"
	"github.com/ddkwork/ux/resources/images"
	"github.com/ddkwork/ux/widget/material"
)

type ButtonAnimation struct {
	Text      string
	Style     ButtonAnimationStyle
	Clickable *widget.Clickable
	Label     *widget.Label
	Focused   bool
	Disabled  bool
	Loading   bool
	Show      bool
	Flex      bool

	animClickable    *widget.Clickable
	hoverSwitchState bool

	th *material.Theme
	do func(gtx layout.Context)
}

type ButtonAnimationStyle struct {
	Rounded     rounded
	TextSize    unit.Sp
	Inset       layout.Inset
	Font        font.Font
	Icon        []byte
	IconGap     unit.Dp
	Animation   buttonAnimationOption
	Border      widget.Border
	LoadingIcon *widget.Icon
	Colors      buttonColors
}

func NewButton(button *widget.Clickable, icon []byte, text string, do func(gtx layout.Context)) *ButtonAnimation {
	radius := unit.Dp(16)
	if icon == nil {
		radius = 13
	}
	return &ButtonAnimation{
		Text: text,
		Style: ButtonAnimationStyle{
			Rounded:  uniformRounded(unit.Dp(12)),
			TextSize: unit.Sp(12),
			Inset:    layout.UniformInset(unit.Dp(8)),
			// Font:        font.Font{},
			Icon:      icon,
			IconGap:   unit.Dp(1),
			Animation: newButtonAnimationScale(.98),
			Border: widget.Border{
				Color:        colors.Grey200,
				CornerRadius: radius,
				Width:        0.5,
			},
			LoadingIcon: nil,
			Colors: buttonColors{
				TextColor:            th.Color.DefaultTextWhiteColor,
				BackgroundColor:      th.Color.InputFocusedBgColor,
				HoverBackgroundColor: &th.ContrastFg,
				HoverTextColor:       &th.Color.HoveredBorderBlueColor,
				BorderColor:          colors.White,
			},
		},
		Clickable:        button,
		Label:            new(widget.Label),
		Focused:          false,
		Disabled:         false,
		Loading:          false,
		Show:             true,
		Flex:             false,
		animClickable:    new(widget.Clickable),
		hoverSwitchState: false,
		th:               th,
		do:               do,
	}
}

func newButtonAnimationScale(v float32) buttonAnimationOption {
	animationEnter := animation.NewAnimation(false,
		gween.NewSequence(
			gween.New(1, v, .1, ease.Linear),
		),
	)

	animationLeave := animation.NewAnimation(false,
		gween.NewSequence(
			gween.New(v, 1, .1, ease.Linear),
		),
	)

	animationClick := animation.NewAnimation(false,
		gween.NewSequence(
			gween.New(1, v, .1, ease.Linear),
			gween.New(v, 1, .4, ease.OutBounce),
		),
	)

	animationLoading := animation.NewAnimation(false,
		gween.NewSequence(
			gween.New(0, 1, 1, ease.Linear),
		),
	)
	animationLoading.Sequence.SetLoop(-1)

	return buttonAnimationOption{
		animationEnter:   animationEnter,
		transformEnter:   animation.TransformScaleCenter,
		animationLeave:   animationLeave,
		transformLeave:   animation.TransformScaleCenter,
		animationClick:   animationClick,
		transformClick:   animation.TransformScaleCenter,
		animationLoading: animationLoading,
		transformLoading: animation.TransformRotate,
	}
}

func (btn *ButtonAnimation) SetLoading(loading bool) {
	btn.Loading = loading
	btn.Disabled = loading

	animationLoading := btn.Style.Animation.animationLoading
	if loading {
		animationLoading.Reset().Start()
	} else {
		animationLoading.Reset()
	}
}

func (btn *ButtonAnimation) Clicked(gtx layout.Context) bool {
	if btn.Disabled {
		return false
	}
	return btn.Clickable.Clicked(gtx)
}

func (btn *ButtonAnimation) Layout(gtx layout.Context) layout.Dimensions {
	if !btn.Show {
		return layout.Dimensions{} // todo test
	}
	if btn.Clicked(gtx) {
		if btn.do != nil {
			btn.do(gtx)
		}
	}

	animationEnter := btn.Style.Animation.animationEnter
	transformEnter := btn.Style.Animation.transformEnter
	animationLeave := btn.Style.Animation.animationLeave
	transformLeave := btn.Style.Animation.transformLeave
	animationClick := btn.Style.Animation.animationClick
	transformClick := btn.Style.Animation.transformClick
	animationLoading := btn.Style.Animation.animationLoading
	transformLoading := btn.Style.Animation.transformLoading

	clickable := btn.Clickable
	animClickable := btn.animClickable
	style := btn.Style
	colors := btn.Style.Colors

	return clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return animClickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			semantic.Button.Add(gtx.Ops)

			{
				if animationEnter != nil {
					state := animationEnter.Update(gtx)
					if state.Active {
						defer transformEnter(gtx, state.Value).Push(gtx.Ops).Pop()
					}
				}
			}

			{
				if animationLeave != nil {
					state := animationLeave.Update(gtx)
					if state.Active {
						defer transformLeave(gtx, state.Value).Push(gtx.Ops).Pop()
					}
				}
			}

			{
				if animationClick != nil {
					state := animationClick.Update(gtx)
					if state.Active {
						defer transformClick(gtx, state.Value).Push(gtx.Ops).Pop()
					}
				}
			}

			backgroundColor := colors.BackgroundColor
			textColor := colors.TextColor

			if !btn.Disabled {
				if animClickable.Hovered() {
					pointer.CursorPointer.Add(gtx.Ops)
					if colors.HoverBackgroundColor != nil {
						backgroundColor = *colors.HoverBackgroundColor
					}

					if colors.HoverTextColor != nil {
						textColor = *colors.HoverTextColor
					}
				}

				if animClickable.Hovered() && !btn.hoverSwitchState {
					btn.hoverSwitchState = true

					if animationEnter != nil {
						animationEnter.Start()
					}

					if animationLeave != nil {
						animationLeave.Reset()
					}

					gtx.Execute(op.InvalidateCmd{})
				}

				if !animClickable.Hovered() && btn.hoverSwitchState {
					btn.hoverSwitchState = false

					if animationLeave != nil {
						animationLeave.Start()
					}

					if animationEnter != nil {
						animationEnter.Reset()
					}

					gtx.Execute(op.InvalidateCmd{})
				}

				if animClickable.Clicked(gtx) {
					if animationClick != nil {
						animationClick.Reset().Start()
					}
				}
			} // else {
			// animationLeave.Reset()
			// animationEnter.Reset()
			// animationClick.Reset()
			// }

			c := op.Record(gtx.Ops)
			style.Border.Color = colors.BorderColor
			dims := style.Border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return style.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					var iconWidget layout.Widget
					if style.Icon != nil {
						iconWidget = func(gtx layout.Context) layout.Dimensions {
							// icon := style.Icon
							// if style.LoadingIcon != nil && btn.Loading {
							//	icon = style.LoadingIcon
							// }

							var dims layout.Dimensions
							r := op.Record(gtx.Ops)
							if btn.Flex {
								dims = layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return images.Layout(gtx, style.Icon, textColor, 24)
								})
							} else {
								dims = images.Layout(gtx, style.Icon, textColor, 24)
							}
							c := r.Stop()

							{
								gtx.Constraints.Min = dims.Size

								if animationLoading != nil {
									state := animationLoading.Update(gtx)
									if state.Active {
										defer transformLoading(gtx, state.Value).Push(gtx.Ops).Pop()
									}
								}
							}

							c.Add(gtx.Ops)
							return dims
						}
					}

					if btn.Text != "" {
						var childs []layout.FlexChild

						if iconWidget != nil {
							childs = append(childs,
								layout.Rigid(iconWidget),
								layout.Rigid(layout.Spacer{Width: style.IconGap}.Layout),
							)
						}

						childs = append(childs,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								paint.ColorOp{Color: textColor}.Add(gtx.Ops)
								label := material.Label(th, style.TextSize, btn.Text)
								label.Alignment = text.Middle
								return label.Layout(gtx)
								return btn.Label.Layout(gtx, btn.th.Shaper, style.Font,
									style.TextSize, btn.Text, op.CallOp{})
							}),
						)

						return layout.Flex{
							// Axis:      layout.Horizontal,//todo support vertical
							Axis:      layout.Vertical,
							Alignment: layout.Middle,
						}.Layout(gtx, childs...)
					} else {
						return iconWidget(gtx)
					}
				})
			})
			m := c.Stop()

			if btn.Flex {
				dims = layout.Dimensions{Size: gtx.Constraints.Max}
			}

			bounds := image.Rectangle{Max: dims.Size}
			paint.FillShape(gtx.Ops, backgroundColor,
				clip.RRect{
					Rect: bounds,
					SE:   gtx.Dp(style.Rounded.SE),
					SW:   gtx.Dp(style.Rounded.SW),
					NE:   gtx.Dp(style.Rounded.NE),
					NW:   gtx.Dp(style.Rounded.NW),
				}.Op(gtx.Ops),
			)

			m.Add(gtx.Ops)
			return dims
		})
	})
}

// /////////////////////////////////////

type buttonAnimationOption struct {
	animationEnter   *animation.Animation
	transformEnter   animation.TransformFunc
	animationLeave   *animation.Animation
	transformLeave   animation.TransformFunc
	animationClick   *animation.Animation
	transformClick   animation.TransformFunc
	animationLoading *animation.Animation
	transformLoading animation.TransformFunc
}

type buttonColors struct {
	TextColor            color.NRGBA
	BackgroundColor      color.NRGBA
	HoverBackgroundColor *color.NRGBA
	HoverTextColor       *color.NRGBA
	BorderColor          color.NRGBA
}
type rounded struct {
	NW unit.Dp
	NE unit.Dp
	SW unit.Dp
	SE unit.Dp
}

func uniformRounded(r unit.Dp) rounded {
	return rounded{
		NW: r, NE: r, SW: r, SE: r,
	}
}
