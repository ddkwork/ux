package ux

import (
	"github.com/ddkwork/ux/giosvg"
	"image"
	"image/color"

	"github.com/ddkwork/golibrary/mylog"

	"gioui.org/io/input"
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/ddkwork/ux/animationButton"
)

func NewButtonAnimation(text string, icon *widget.Icon, callBack func()) *animationButton.Button {
	style := animationButton.ButtonStyle{
		Rounded:  animationButton.UniformRounded(unit.Dp(12)),
		TextSize: unit.Sp(12),
		Inset:    layout.UniformInset(unit.Dp(4)),
		// Font:        font.Font{},
		Icon:      icon,
		IconGap:   unit.Dp(1),
		Animation: animationButton.NewButtonAnimationDefault(),
		Border: widget.Border{
			Color:        White,
			CornerRadius: 14,
			Width:        0.5,
		},
		LoadingIcon: nil,
		Colors: animationButton.ButtonColors{
			TextColor:            th.Color.DefaultTextWhiteColor,
			BackgroundColor:      th.Color.InputFocusedBgColor,
			HoverBackgroundColor: &th.ContrastFg,
			HoverTextColor:       &th.Color.HoveredBorderBlueColor,
			BorderColor:          White,
		},
	}
	return animationButton.NewButton(style, th.Theme, text, callBack)
}

func NewButtonAnimationDefault() animationButton.ButtonAnimation {
	return NewButtonAnimationScale(.98)
}

func NewButtonAnimationScale(v float32) animationButton.ButtonAnimation {
	return animationButton.NewButtonAnimationScale(v)
}

type Button struct {
	icon     *widget.Icon
	svgIcon  *giosvg.Icon
	iconRect bool

	Axis layout.Axis
	*widget.Clickable
	callBack func()
	text     string
	Inset    layout.Inset
	spacer   unit.Dp // icon insets

	radius unit.Dp // border radius
	width  unit.Dp // border width
}

func NewButton(text string, callBack func()) *Button {
	m := &Button{
		icon:      nil,
		iconRect:  false,
		Axis:      0,
		Clickable: new(widget.Clickable),
		callBack:  callBack,
		text:      text,
		Inset: layout.Inset{
			Top:    5,
			Bottom: 5,
			Left:   17,
			Right:  17,
		},
		spacer: 1,
		radius: 16,
		width:  0.5,
	}
	return m
}

func (m *Button) SetHorizontal() *Button {
	m.Axis = layout.Horizontal
	return m
}

func (m *Button) SetVertical() *Button {
	m.Axis = layout.Vertical
	return m
}

func (m *Button) SetRectIcon(iconRect bool) *Button {
	m.iconRect = iconRect
	return m
}

func (m *Button) SetIcon(icon *widget.Icon) *Button {
	m.icon = icon
	return m
}

func (m *Button) SetSVGIcon(content string) *Button {
	m.svgIcon = Svg2Icon(content)
	return m
}

func Svg2Icon(b string) *giosvg.Icon {
	return giosvg.NewIcon(mylog.Check2(giosvg.NewVector([]byte(b))))
}

func (m *Button) Layout(gtx layout.Context) layout.Dimensions {
	if m.callBack != nil {
		if m.Clicked(gtx) {
			m.callBack()
		}
	}

	// textOnly := m.icon == nil && m.svgIcon == nil && m.text != ""
	// iconOnly := m.icon != nil || m.svgIcon != nil && m.text == ""
	// svgOnly := m.svgIcon != nil && m.text == ""
	// iconAndText := m.icon != nil || m.svgIcon != nil && m.text != ""
	// svgAndText := m.svgIcon != nil && m.text != ""

	if m.icon == nil && m.svgIcon == nil && m.text != "" { // 只有文字
		// btn := material.Button(th.Theme, m.Clickable, m.text)
		// btn.Inset = layout.UniformInset(2) // todo test

		return material.ButtonLayoutStyle{
			Background:   th.Color.InputFocusedBgColor,
			CornerRadius: m.radius,
			Button:       m.Clickable,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return widget.Border{
				Color:        Grey200,
				Width:        m.width,
				CornerRadius: m.radius,
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return m.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					label := material.Body1(th.Theme, m.text)
					label.Color = th.Color.DefaultTextWhiteColor
					if m.Hovered() {
						label.TextSize++
					}
					return label.Layout(gtx)
				})
			})
		})
	}

	if m.text == "" && m.icon != nil || m.svgIcon != nil { // 树形层级图标，没有文字
		if m.iconRect { // 带图标的编辑框，图标背景色和按钮背景色一致
			return material.Clickable(gtx, m.Clickable, func(gtx C) D {
				sz := gtx.Dp(defaultIconSize)
				size := image.Pt(sz, sz)
				gtx.Constraints.Min = size
				gtx.Constraints.Max = size
				background := func(gtx C) D {
					defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Max}, sz/2).Push(gtx.Ops).Pop() // 使用UniformRRect绘制圆形背景
					var backgroundColor color.NRGBA
					if m.Hovered() || gtx.Focused(m) {
						backgroundColor = th.Fg // 悬停或聚焦时的颜色
					} else {
						backgroundColor = th.Color.InputFocusedBgColor // 默认颜色
					}
					paint.Fill(gtx.Ops, backgroundColor)
					return layout.Dimensions{Size: gtx.Constraints.Min}
				}
				if m.icon != nil {
					return layout.Background{}.Layout(gtx, background, func(gtx layout.Context) layout.Dimensions {
						return m.icon.Layout(gtx, th.Theme.Bg)
					})
				}
				return layout.UniformInset(1).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return m.svgIcon.Layout(gtx)
				})
				return layout.Background{}.Layout(gtx, background, func(gtx layout.Context) layout.Dimensions {
					return m.svgIcon.Layout(gtx)
				})
			})
		}
		btn := material.IconButton(th.Theme, m.Clickable, m.icon, m.text)
		return btn.Layout(gtx)
	}

	// 图标和文字都有,两个都画
	return material.ButtonLayout(th.Theme, m.Clickable).Layout(gtx, func(gtx C) D {
		top := m.Inset.Top - 2
		bottom := m.Inset.Bottom - 2
		if top < 0 {
			top = 0
		}
		if bottom < 0 {
			bottom = 0
		}
		return layout.Inset{Top: top, Bottom: bottom, Left: m.Inset.Left, Right: m.Inset.Right}.Layout(gtx, func(gtx C) D {
			iconAndLabel := layout.Flex{Axis: m.Axis, Alignment: layout.Middle}
			layIcon := layout.Rigid(func(gtx C) D {
				var d layout.Dimensions
				if m.icon != nil {
					d = m.icon.Layout(gtx, th.Theme.Fg)
				} else {
					d = m.svgIcon.Layout(gtx) // todo theme check
				}

				if m.Axis == layout.Horizontal {
					return layout.Inset{Right: m.spacer}.Layout(gtx, func(gtx C) D {
						return d
					})
				}
				return layout.Inset{Bottom: m.spacer}.Layout(gtx, func(gtx C) D {
					return d
				})
			})

			layLabel := layout.Rigid(func(gtx C) D {
				l := material.Body1(th.Theme, m.text)
				l.Color = th.Theme.Palette.ContrastFg
				if m.Axis == layout.Horizontal {
					return layout.Inset{Left: m.spacer}.Layout(gtx, func(gtx C) D {
						return l.Layout(gtx)
					})
				}
				return layout.Inset{Top: m.spacer}.Layout(gtx, func(gtx C) D {
					return l.Layout(gtx)
				})
			})

			return iconAndLabel.Layout(gtx, layIcon, layLabel)
		})
	})
}

//	func SaveButtonLayout(gtx DefaultDraw.Context, theme *material.Theme, rowClick *widget.Clickable) DefaultDraw.Dimensions {
//		border := widget.Border{
//			Color:        th.Bg, //todo: change color to th.Fg
//			Width:        unit.Dp(1),
//			CornerRadius: unit.Dp(4),
//		}
//
//		return DefaultDraw.Inset{Left: unit.Dp(15)}.Layout(gtx, func(gtx DefaultDraw.Context) DefaultDraw.Dimensions {
//			return border.Layout(gtx, func(gtx DefaultDraw.Context) DefaultDraw.Dimensions {
//				return Clickable(gtx, rowClick, func(gtx DefaultDraw.Context) DefaultDraw.Dimensions {
//					return DefaultDraw.Inset{Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, func(gtx DefaultDraw.Context) DefaultDraw.Dimensions {
//						return DefaultDraw.Flex{Axis: DefaultDraw.Vertical, Alignment: DefaultDraw.Middle}.Layout(gtx,
//							DefaultDraw.Rigid(func(gtx DefaultDraw.Context) DefaultDraw.Dimensions {
//								return DefaultDraw.Flex{Axis: DefaultDraw.Horizontal, Alignment: DefaultDraw.Middle}.Layout(gtx,
//									DefaultDraw.Rigid(func(gtx DefaultDraw.Context) DefaultDraw.Dimensions {
//										gtx.Constraints.Max.X = gtx.Dp(16)
//										return IconSave.Layout(gtx, th.Palette.ContrastFg)
//									}),
//									DefaultDraw.Rigid(func(gtx DefaultDraw.Context) DefaultDraw.Dimensions {
//										return material.Body1(theme, "Save").Layout(gtx)
//									}),
//								)
//							}),
//						)
//					})
//				})
//			})
//		})
//	}
func IconButton(icon *widget.Icon, button *widget.Clickable, description string) material.IconButtonStyle {
	return material.IconButtonStyle{
		Background:  th.Palette.Bg,
		Color:       WithAlpha(th.Palette.Fg, 0xb6),
		Icon:        icon,
		Size:        18,
		Inset:       layout.UniformInset(4),
		Button:      button,
		Description: description,
	}
}

type IconButton2 struct {
	Icon                 *widget.Icon
	Size                 unit.Dp
	Color                color.NRGBA
	BackgroundColor      color.NRGBA
	BackgroundColorHover color.NRGBA

	SkipFocus bool
	Clickable *widget.Clickable

	OnClick func()
}

func (ib *IconButton2) Layout(gtx layout.Context) layout.Dimensions {
	if ib.BackgroundColor == (color.NRGBA{}) {
		ib.BackgroundColor = th.Color.DefaultIconColor
	}

	if ib.BackgroundColorHover == (color.NRGBA{}) {
		ib.BackgroundColorHover = Hovered(ib.BackgroundColor)
	}

	for ib.Clickable.Clicked(gtx) {
		if ib.OnClick != nil {
			ib.OnClick()
		}
	}

	return ib.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		semantic.Button.Add(gtx.Ops)
		return layout.Background{}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, 4).Push(gtx.Ops).Pop()
				background := ib.BackgroundColor
				if gtx.Source == (input.Source{}) {
					background = Disabled(ib.BackgroundColor)
				} else if ib.Clickable.Hovered() || (gtx.Focused(ib.Clickable) && !ib.SkipFocus) {
					background = ib.BackgroundColorHover
				}
				paint.Fill(gtx.Ops, background)
				return layout.Dimensions{Size: gtx.Constraints.Min}
			},
			func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Dp(ib.Size)
				return ib.Icon.Layout(gtx, ib.Color)
			},
		)
	})
}

type NavButton struct {
	btn          widget.Clickable
	cornerRadius unit.Dp
	borderWidth  unit.Dp
	borderColor  color.NRGBA
	background   color.NRGBA
	text         string
}

func NewNavButton(text string) *NavButton {
	return &NavButton{
		btn:          widget.Clickable{},
		cornerRadius: 16,
		borderWidth:  1,
		borderColor:  Grey200,
		background:   color.NRGBA{},
		text:         text,
	}
}

func (btn *NavButton) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return material.ButtonLayoutStyle{
		Background:   btn.background,
		CornerRadius: btn.cornerRadius,
		Button:       &btn.btn,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return widget.Border{
			Color:        btn.borderColor,
			Width:        btn.borderWidth,
			CornerRadius: btn.cornerRadius,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    8,
				Bottom: 8,
				Left:   20,
				Right:  20,
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.Body1(th, btn.text)
				return label.Layout(gtx)
			})
		})
	})
}
