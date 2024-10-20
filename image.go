package ux

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/mylog"
)

type Image struct {
	src     string
	imageOp paint.ImageOp
}

func NewImage(src string) *Image {
	image := &Image{
		src: src,
	}
	data := mylog.Check2(image.LoadImage(src))

	image.imageOp = paint.NewImageOp(data)
	return image
}

func (i *Image) Layout(gtx layout.Context) layout.Dimensions {
	return widget.Image{
		Src:      i.imageOp,
		Fit:      widget.Unscaled,
		Position: layout.Center,
		Scale:    1.0,
	}.Layout(gtx)
}

func (i *Image) LoadImage(fileName string) (image.Image, error) {
	file := mylog.Check2(os.ReadFile(fmt.Sprintf("%s", fileName)))

	// 获取fileName后缀
	temp := strings.Split(fileName, ".")
	suffix := temp[len(temp)-1]

	var img image.Image
	if suffix == "png" {
		img = mylog.Check2(png.Decode(bytes.NewReader(file)))
	} else if suffix == "jpg" {
		img = mylog.Check2(jpeg.Decode(bytes.NewReader(file)))
	}

	return img, nil
}
