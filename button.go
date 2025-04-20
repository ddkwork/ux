package ux

import (
	"image"
	"image/color"
	"math"

	"github.com/ddkwork/ux/resources/colors"
	"github.com/ddkwork/ux/resources/icons"

	"github.com/ddkwork/ux/animationButton"
	"github.com/ddkwork/ux/widget/material"

	"gioui.org/font"
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

func NewButtonAnimation(button *widget.Clickable, icon any, text string, do func(gtx layout.Context)) *animationButton.Button {
	style := animationButton.ButtonStyle{
		Rounded:  animationButton.UniformRounded(unit.Dp(12)),
		TextSize: unit.Sp(12),
		Inset:    layout.UniformInset(unit.Dp(8)),
		// Font:        font.Font{},
		Icon:      icon,
		IconGap:   unit.Dp(1),
		Animation: animationButton.NewButtonAnimationDefault(),
		Border: widget.Border{
			Color:        colors.Grey200,
			CornerRadius: 16,
			Width:        0.5,
		},
		LoadingIcon: nil,
		Colors: animationButton.ButtonColors{
			TextColor:            th.Color.DefaultTextWhiteColor,
			BackgroundColor:      th.Color.InputFocusedBgColor,
			HoverBackgroundColor: &th.ContrastFg,
			HoverTextColor:       &th.Color.HoveredBorderBlueColor,
			BorderColor:          colors.White,
		},
	}
	if icon == nil {
		style.Border.CornerRadius = 13
	}
	return animationButton.NewButton(style, button, th, text, do)
}

func NewButtonAnimationDefault() animationButton.ButtonAnimation {
	return NewButtonAnimationScale(.98)
}

func NewButtonAnimationScale(v float32) animationButton.ButtonAnimation {
	return animationButton.NewButtonAnimationScale(v)
}

///////////////////////////////////////////////////////////////////////////////

type ButtonStyle struct {
	Text            string
	Icon            any
	IconPositionEnd bool
	// Color is the text color.
	color        color.NRGBA
	Font         font.Font
	TextSize     unit.Sp
	IconSize     unit.Dp
	Background   color.NRGBA
	CornerRadius unit.Dp
	Inset        layout.Inset
	Button       *widget.Clickable
	shaper       *text.Shaper
	border       widget.Border
}

type ButtonLayoutStyle struct {
	Background   color.NRGBA
	CornerRadius unit.Dp
	Button       *widget.Clickable
}

type IconButtonStyle struct {
	Background color.NRGBA
	// Color is the icons color.
	Color color.NRGBA
	Icon  *widget.Icon
	// Size is the icons size.
	Size        unit.Dp
	Inset       layout.Inset
	Button      *widget.Clickable
	Description string
}

func iconButtonSmall(button *widget.Clickable, icon any, txt string) ButtonStyle {
	style := Button(button, icon, txt)
	style.Inset = layout.Inset{}
	style.IconSize = defaultHierarchyIconSize
	style.border = widget.Border{}
	return style
}

func Button(button *widget.Clickable, icon any, text string) ButtonStyle {
	b := ButtonStyle{
		Text:            text,
		Icon:            icon,
		IconPositionEnd: false,
		color:           th.Fg,
		Font:            font.Font{},
		TextSize:        th.TextSize * 14.0 / 16.0,
		IconSize:        18,
		Background:      th.Bg,
		CornerRadius:    4,
		Inset:           layout.UniformInset(8),
		Button:          button,
		shaper:          th.Shaper,
		border: widget.Border{
			Color:        colors.Grey200,
			CornerRadius: 16,
			Width:        .5,
		},
	}
	b.Font.Typeface = th.Face
	return b
}

func (b ButtonStyle) Layout(gtx layout.Context) layout.Dimensions {
	return ButtonLayoutStyle{
		Background:   b.Background,
		CornerRadius: b.CornerRadius,
		Button:       b.Button,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		iconDims := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if b.Icon != nil {
				return layout.Inset{Right: unit.Dp(0)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Dp(b.IconSize)
					gtx.Constraints.Max.X = gtx.Dp(b.IconSize)
					return icons.Layout(gtx, b.Icon, b.color, b.IconSize)
				})
			}
			return layout.Dimensions{}
		})
		labelDims := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Label(th, b.TextSize, b.Text).Layout(gtx)
		})

		items := []layout.FlexChild{iconDims, labelDims}
		if b.Icon != nil && b.IconPositionEnd {
			items = []layout.FlexChild{labelDims, iconDims}
			b.Inset.Right = unit.Dp(5)
		}
		background := b.border.Color
		switch {
		case !gtx.Enabled():
			background = Disabled(b.Background)
		case b.Button.Hovered() || gtx.Focused(b.Button):
			background = colors.Grey400
		}
		b.border.Color = background
		if b.Icon == nil {
			b.border.CornerRadius = 13
		}
		return b.border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return b.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx, items...)
			})
		})
	})
}

func (b ButtonLayoutStyle) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	min := gtx.Constraints.Min
	return b.Button.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		semantic.Button.Add(gtx.Ops)
		return layout.Background{}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				rr := gtx.Dp(b.CornerRadius)
				defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, rr).Push(gtx.Ops).Pop()
				//background := b.Background
				//switch {
				//case !gtx.Enabled():
				//	background = Disabled(b.Background)
				//case b.Button.Hovered() || gtx.Focused(b.Button):
				//	background = Hovered(b.Background)
				//}
				//paint.Fill(gtx.Ops, background)
				for _, c := range b.Button.History() {
					drawInk(gtx, c)
				}
				return layout.Dimensions{Size: gtx.Constraints.Min}
			},
			func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min = min
				return layout.Center.Layout(gtx, w)
			},
		)
	})
}

func drawInk(gtx layout.Context, c widget.Press) {
	// duration is the number of seconds for the
	// completed animation: expand while fading in, then
	// out.
	const (
		expandDuration = float32(0.5)
		fadeDuration   = float32(0.9)
	)

	now := gtx.Now

	t := float32(now.Sub(c.Start).Seconds())

	end := c.End
	if end.IsZero() {
		// If the press hasn't ended, don't fade-out.
		end = now
	}

	endt := float32(end.Sub(c.Start).Seconds())

	// Compute the fade-in/out position in [0;1].
	var alphat float32
	{
		var haste float32
		if c.Cancelled {
			// If the press was cancelled before the inkwell
			// was fully faded in, fast forward the animation
			// to match the fade-out.
			if h := 0.5 - endt/fadeDuration; h > 0 {
				haste = h
			}
		}
		// Fade in.
		half1 := min(t/fadeDuration+haste, 0.5)

		// Fade out.
		half2 := float32(now.Sub(end).Seconds())
		half2 /= fadeDuration
		half2 += haste
		if half2 > 0.5 {
			// Too old.
			return
		}

		alphat = half1 + half2
	}

	// Compute the expand position in [0;1].
	sizet := t
	if c.Cancelled {
		// Freeze expansion of cancelled presses.
		sizet = endt
	}
	sizet /= expandDuration

	// Animate only ended presses, and presses that are fading in.
	if !c.End.IsZero() || sizet <= 1.0 {
		gtx.Execute(op.InvalidateCmd{})
	}

	if sizet > 1.0 {
		sizet = 1.0
	}

	if alphat > .5 {
		// Start fadeout after half the animation.
		alphat = 1.0 - alphat
	}
	// Twice the speed to attain fully faded in at 0.5.
	t2 := alphat * 2
	// Beziér ease-in curve.
	alphaBezier := t2 * t2 * (3.0 - 2.0*t2)
	sizeBezier := sizet * sizet * (3.0 - 2.0*sizet)
	size := gtx.Constraints.Min.X
	if h := gtx.Constraints.Min.Y; h > size {
		size = h
	}
	// Cover the entire constraints min rectangle and
	// apply curve values to size and color.
	size = int(float32(size) * 2 * float32(math.Sqrt(2)) * sizeBezier)
	alpha := 0.7 * alphaBezier
	const col = 0.8
	ba, bc := byte(alpha*0xff), byte(col*0xff)
	rgba := MulAlpha(color.NRGBA{A: 0xff, R: bc, G: bc, B: bc}, ba)
	ink := paint.ColorOp{Color: rgba}
	ink.Add(gtx.Ops)
	rr := size / 2
	defer op.Offset(c.Position.Add(image.Point{
		X: -rr,
		Y: -rr,
	})).Push(gtx.Ops).Pop()
	defer clip.UniformRRect(image.Rectangle{Max: image.Pt(size, size)}, rr).Push(gtx.Ops).Pop()
	paint.PaintOp{}.Add(gtx.Ops)
}
