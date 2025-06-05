package ux

import (
	"bytes"
	"image"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/std/mylog"
	_ "github.com/ddkwork/ux/resources/images/ico"
)

type Image struct {
	// src     string
	imageOp paint.ImageOp
}

// NewImage 任务管理器或者音速启动,可以编码为任何图片格式，对于加载进gio，只需要它返回的image.Image接口就可以了
// img, e := gowin32.ExtractPrivateExtractIcons(filename, 128, 128)
func NewImage(b []byte) *Image {
	img, _ := mylog.Check3(image.Decode(bytes.NewReader(b))) //_ "github.com/ddkwork/ux/resources/images/ico" 通过init注册ico解码器
	return &Image{
		imageOp: paint.NewImageOp(img),
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
