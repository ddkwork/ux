package ux

import (
	"image"
	"image/color"

	"github.com/ddkwork/ux/x/colorpicker"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type ColorPicker struct {
	picker     *colorpicker.State
	muxState   *colorpicker.MuxState
	background color.NRGBA
}

func NewColorPicker() *ColorPicker {
	current := color.NRGBA{R: 255, G: 128, B: 75, A: 255}
	picker := colorpicker.State{}
	picker.SetColor(current)
	muxState := colorpicker.NewMuxState(
		[]colorpicker.MuxOption{
			{
				Label: "current",
				Value: &current,
			},
			{
				Label: "background",
				Value: &th.Palette.Bg,
			},
			{
				Label: "foreground",
				Value: &th.Palette.Fg,
			},
		}...)
	return &ColorPicker{
		picker:     &picker,
		muxState:   &muxState,
		background: *muxState.Color(),
	}
}

func (c *ColorPicker) Layout(gtx layout.Context) layout.Dimensions {
	if c.muxState.Update(gtx) {
		c.background = *c.muxState.Color()
	}
	current := color.NRGBA{R: 255, G: 128, B: 75, A: 255}
	if c.picker.Update(gtx) {
		current = c.picker.Color()
		c.muxState.Options["current"] = &current
		c.background = *c.muxState.Color()
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return colorpicker.PickerStyle{
				Label:         "Current",
				Theme:         th.Theme,
				State:         c.picker,
				MonospaceFace: "Go Mono",
			}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return colorpicker.Mux(th.Theme, c.muxState, "Display Right:").Layout(gtx)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					size := gtx.Constraints.Max
					paint.FillShape(gtx.Ops, c.background, clip.Rect(image.Rectangle{Max: size}).Op())
					return D{Size: size}
				}),
			)
		}),
	)
}
