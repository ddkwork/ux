// SPDX-License-Identifier: Unlicense OR MIT

package component

import (
	"image"
	"image/color"

	"github.com/ddkwork/ux/resources/images"

	"github.com/ddkwork/ux/widget/material"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"gioui.org/x/outlay"
)

// SurfaceStyle defines the visual aspects of a material design surface
// with (optionally) rounded corners and a drop shadow.
type SurfaceStyle struct {
	*material.Theme
	// The CornerRadius and Elevation fields of the embedded shadow
	// style also define the corner radius and elevation of the card.
	ShadowStyle
	// Theme background color will be used if empty.
	Fill color.NRGBA
}

// Surface creates a Surface style for the provided theme with sensible default
// elevation and rounded corners.
func Surface(th *material.Theme) SurfaceStyle {
	return SurfaceStyle{
		Theme:       th,
		ShadowStyle: Shadow(unit.Dp(4), unit.Dp(4)),
	}
}

// Layout renders the SurfaceStyle, taking the dimensions of the surface from
// gtx.Constraints.Min.
func (c SurfaceStyle) Layout(gtx C, w layout.Widget) D {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			c.ShadowStyle.Layout(gtx)
			surface := clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, gtx.Dp(c.ShadowStyle.CornerRadius))
			var fill color.NRGBA
			if empty := (color.NRGBA{}); c.Fill == empty {
				fill = c.Theme.Bg
			} else {
				fill = c.Fill
			}
			paint.FillShape(gtx.Ops, fill, surface.Op(gtx.Ops))
			return D{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(w),
	)
}

// DividerStyle defines the presentation of a material divider, as specified
// here: https://material.io/components/dividers
type DividerStyle struct {
	Thickness unit.Dp
	Fill      color.NRGBA
	layout.Inset

	Subheading      material.LabelStyle
	SubheadingInset layout.Inset
}

// Divider creates a simple full-bleed divider.
func Divider(th *material.Theme) DividerStyle {
	return DividerStyle{
		Thickness: unit.Dp(1),
		Fill:      WithAlpha(th.Fg, 0x60),
		Inset: layout.Inset{
			Top:    unit.Dp(8),
			Bottom: unit.Dp(8),
		},
	}
}

// SubheadingDivider creates a full-bleed divider with a subheading.
func SubheadingDivider(th *material.Theme, subheading string) DividerStyle {
	return DividerStyle{
		Thickness: unit.Dp(1),
		Fill:      WithAlpha(th.Fg, 0x60),
		Inset: layout.Inset{
			Top:    unit.Dp(8),
			Bottom: unit.Dp(4),
		},
		Subheading: DividerSubheadingText(th, subheading),
		SubheadingInset: layout.Inset{
			Left:   unit.Dp(8),
			Bottom: unit.Dp(8),
		},
	}
}

// Layout renders the divider. If gtx.Constraints.Min.X is zero, it will
// have zero size and render nothing.
func (d DividerStyle) Layout(gtx layout.Context) layout.Dimensions {
	if gtx.Constraints.Min.X == 0 {
		return D{}
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return d.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				weight := gtx.Dp(d.Thickness)
				line := image.Rectangle{Max: image.Pt(gtx.Constraints.Min.X, weight)}
				paint.FillShape(gtx.Ops, d.Fill, clip.Rect(line).Op())
				return D{Size: line.Max}
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if d.Subheading == (material.LabelStyle{}) {
				return D{}
			}
			return d.SubheadingInset.Layout(gtx, d.Subheading.Layout)
		}),
	)
}

// MenuItemStyle defines the presentation of a Menu element that has a label
// and optionally an images and a hint text.
type MenuItemStyle struct {
	State      *widget.Clickable
	HoverColor color.NRGBA

	LabelInset outlay.Inset
	Label      material.LabelStyle

	Icon      []byte
	IconSize  unit.Dp
	IconInset outlay.Inset
	IconColor color.NRGBA

	Hint      material.LabelStyle
	HintInset outlay.Inset
}

// MenuItem constructs a default MenuItemStyle based on the theme, state, and label.
func MenuItem(th *material.Theme, state *widget.Clickable, label string) MenuItemStyle {
	return MenuItemStyle{
		State: state,
		LabelInset: outlay.Inset{
			Start:  unit.Dp(16),
			End:    unit.Dp(16),
			Top:    unit.Dp(8),
			Bottom: unit.Dp(8),
		},
		IconSize: unit.Dp(24),
		IconInset: outlay.Inset{
			Start: unit.Dp(16),
		},
		IconColor: th.Fg,
		HintInset: outlay.Inset{
			End: unit.Dp(16),
		},
		Label:      material.Body1(th, label),
		HoverColor: WithAlpha(th.ContrastBg, 0x30),
	}
}

// Layout renders the MenuItemStyle. If gtx.Constraints.Min.X is zero, it will render
// itself to be as compact as possible horizontally.
func (m MenuItemStyle) Layout(gtx layout.Context) layout.Dimensions {
	min := gtx.Constraints.Min.X
	compact := min == 0
	return material.Clickable(gtx, m.State, func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				area := image.Rectangle{
					Max: gtx.Constraints.Min,
				}
				if m.State.Hovered() {
					paint.FillShape(gtx.Ops, m.HoverColor, clip.Rect(area).Op())
				}
				return D{Size: area.Max}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = min
				return outlay.Flex{
					Alignment: layout.Middle,
				}.Layout(gtx,
					outlay.Rigid(func(gtx layout.Context) layout.Dimensions {
						if m.Icon == nil {
							return D{}
						}
						return m.IconInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							iconSize := gtx.Dp(m.IconSize)
							gtx.Constraints = layout.Exact(image.Point{X: iconSize, Y: iconSize})
							return images.Layout(gtx, m.Icon, m.IconColor, m.IconSize)
						})
					}),
					outlay.Rigid(func(gtx layout.Context) layout.Dimensions {
						return m.LabelInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return m.Label.Layout(gtx)
						})
					}),
					outlay.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						if compact {
							return D{}
						}
						return D{Size: gtx.Constraints.Min}
					}),
					outlay.Rigid(func(gtx layout.Context) layout.Dimensions {
						if empty := (material.LabelStyle{}); m.Hint == empty {
							return D{}
						}
						return m.HintInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return m.Hint.Layout(gtx)
						})
					}),
				)
			}),
		)
	})
}

// MenuHintText returns a LabelStyle suitable for use as hint text in a
// MenuItemStyle.
func MenuHintText(th *material.Theme, label string) material.LabelStyle {
	l := material.Body1(th, label)
	l.Color = WithAlpha(l.Color, 0xaa)
	return l
}

// DividerSubheadingText returns a LabelStyle suitable for use as a subheading
// in a divider.
func DividerSubheadingText(th *material.Theme, label string) material.LabelStyle {
	l := material.Body2(th, label)
	l.Color = WithAlpha(l.Color, 0xaa)
	return l
}

// MenuState holds the state of a menu material design component
// across frames.
type MenuState struct {
	OptionList layout.List
	Options    []func(gtx layout.Context) layout.Dimensions
}

// MenuStyle defines the presentation of a material design menu component.
type MenuStyle struct {
	*MenuState
	*material.Theme
	// Inset applied around the rendered contents of the state's Options field.
	layout.Inset
	SurfaceStyle
}

// Menu constructs a menu with the provided state and a default Surface behind
// it.
func Menu(th *material.Theme, state *MenuState) MenuStyle {
	m := MenuStyle{
		Theme:        th,
		MenuState:    state,
		SurfaceStyle: Surface(th),
		Inset: layout.Inset{
			Top:    unit.Dp(8),
			Bottom: unit.Dp(8),
		},
	}
	m.OptionList.Axis = layout.Vertical
	return m
}

// Layout renders the menu.
func (m MenuStyle) Layout(gtx layout.Context) layout.Dimensions {
	var fakeOps op.Ops
	originalOps := gtx.Ops
	gtx.Ops = &fakeOps
	maxWidth := 0
	for _, w := range m.Options {
		dims := w(gtx)
		if dims.Size.X > maxWidth {
			maxWidth = dims.Size.X
		}
	}
	gtx.Ops = originalOps
	return m.SurfaceStyle.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return m.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return m.OptionList.Layout(gtx, len(m.Options), func(gtx C, index int) D {
				gtx.Constraints.Min.X = maxWidth
				gtx.Constraints.Max.X = maxWidth
				return m.Options[index](gtx)
			})
		})
	})
}
