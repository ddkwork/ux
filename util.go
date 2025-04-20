package ux

import (
	_ "embed"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/ddkwork/ux/resources/colors"
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

func (i Icon) Layout(gtx layout.Context) layout.Dimensions {
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

func rgb(c uint32) color.NRGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func drawColumnDivider(gtx layout.Context, col int, color color.NRGBA) { //绘制列分隔条,todo最后一列没绘制到
	if col > 0 {
		dividerWidth := 1
		tallestHeight := gtx.Constraints.Min.Y
		stack3 := clip.Rect{Max: image.Pt(dividerWidth, tallestHeight)}.Push(gtx.Ops)
		paint.Fill(gtx.Ops, color)
		stack3.Pop()
	}
}

func HighlightRow(gtx layout.Context) { // 高亮选中行为蓝色
	paint.FillShape(gtx.Ops, color.NRGBA(colors.Blue400), clip.Rect{Max: gtx.Constraints.Max}.Op())
}

func DrawCrosswalk(gtx layout.Context, row int) { // 绘制斑马线
	if row%2 == 0 {
		paint.FillShape(gtx.Ops, rowWhiteColor, clip.Rect{Max: gtx.Constraints.Max}.Op())
	} else {
		paint.FillShape(gtx.Ops, rowBlackColor, clip.Rect{Max: gtx.Constraints.Max}.Op())
	}
}
