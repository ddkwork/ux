package ux

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/ux/authlayout"
	"github.com/inkeliz/giosvg"
)

func NewSVGButton(title string, icon *giosvg.Icon, callback func()) *GoogleDummyButton {
	DefaultDarkGoogleButtonStyle.Text = title
	return &GoogleDummyButton{
		ButtonStyle: DefaultDarkGoogleButtonStyle,
		Pointer:     &authlayout.Pointer{},
		icon:        icon,
		callback:    callback,
	}
}

func Svg2Icon(b []byte) *giosvg.Icon {
	return giosvg.NewIcon(mylog.Check2(giosvg.NewVector(b)))
}

type GoogleDummyButton struct {
	authlayout.ButtonStyle
	*authlayout.Pointer
	icon     *giosvg.Icon
	callback func()
}

func (g *GoogleDummyButton) Layout(gtx layout.Context) layout.Dimensions {
	return g.LayoutText(gtx, g.Text)
}

func (g *GoogleDummyButton) LayoutText(gtx layout.Context, text string) layout.Dimensions {
	if g.icon == nil {
		// g.icon = giosvg.NewIcon(internal.VectorGoogleLogo)//todo
	}
	if g.Clicked(gtx) {
		if g.callback != nil {
			g.callback()
		}
	}
	return g.LayoutText_(gtx, g.icon, g.Pointer, text, 0, gtx.Dp(24))
}

var DefaultDarkGoogleButtonStyle = authlayout.ButtonStyle{
	Text:                "Continue with Google",
	TextSize:            unit.Dp(16),
	TextFont:            font.Font{},
	TextShaper:          th.Shaper,
	TextColor:           color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	TextAlignment:       layout.Middle,
	IconAlignment:       layout.Start,
	BackgroundColor:     th.Bg,
	BackgroundIconColor: th.Color.DefaultIconColor,
	Format:              authlayout.FormatRounded,
}
