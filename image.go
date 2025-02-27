package ux

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/stream"
)

type Image struct {
	src     string
	imageOp paint.ImageOp
}

func NewImage(src string) *Image {
	return &Image{
		src:     src,
		imageOp: paint.NewImageOp(stream.LoadImage(src)),
	}
}

func (i *Image) Layout(gtx layout.Context) layout.Dimensions {
	return widget.Image{
		Src:      i.imageOp,
		Fit:      widget.Unscaled,
		Position: layout.Center,
		Scale:    1.0,
	}.Layout(gtx)
}
