package ux

import (
	_ "embed"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"

	"gioui.org/layout"
	"gioui.org/unit"

	"gioui.org/widget"
)

type Icon struct {
	*widget.Icon
	Color color.NRGBA
	Size  unit.Dp
}

func (i Icon) Layout(gtx C) D {
	if i.Size <= 0 {
		i.Size = unit.Dp(18)
	}
	if i.Color == (color.NRGBA{}) {
		i.Color = WithAlpha(th.Palette.Fg, 0xb6)
	}

	iconSize := gtx.Dp(i.Size)
	gtx.Constraints = layout.Exact(image.Pt(iconSize, iconSize))

	return i.Icon.Layout(gtx, i.Color)
}

// MulAlpha applies the alpha to the color.
func MulAlpha(c color.NRGBA, alpha uint8) color.NRGBA {
	c.A = uint8(uint32(c.A) * uint32(alpha) / 0xFF)
	return c
}

// Disabled blends color towards the luminance and multiplies alpha.
// Blending towards luminance will desaturate the color.
// Multiplying alpha blends the color together more with the background.
func Disabled(c color.NRGBA) (d color.NRGBA) {
	const r = 80 // blend ratio
	lum := approxLuminance(c)
	d = mix(c, color.NRGBA{A: c.A, R: lum, G: lum, B: lum}, r)
	d = MulAlpha(d, 128+32)
	return
}

// Hovered blends dark colors towards white, and light colors towards
// black. It is approximate because it operates in non-linear sRGB space.
func Hovered(c color.NRGBA) (h color.NRGBA) {
	if c.A == 0 {
		// Provide a reasonable default for transparent widgets.
		return color.NRGBA{A: 0x44, R: 0x88, G: 0x88, B: 0x88}
	}
	const ratio = 0x20
	m := color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: c.A}
	if approxLuminance(c) > 128 {
		m = color.NRGBA{A: c.A}
	}
	return mix(m, c, ratio)
}
