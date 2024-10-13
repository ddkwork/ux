package ux

import (
	"bytes"
	"fmt"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

type Image struct {
	src     string
	imageOp paint.ImageOp
}

func NewImage(src string) *Image {
	image := &Image{
		src: src,
	}
	data, err := image.LoadImage(src)
	if err != nil {
		panic(err)
	}
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
	file, err := os.ReadFile(fmt.Sprintf("%s", fileName))
	if err != nil {
		return nil, err
	}

	// 获取fileName后缀
	temp := strings.Split(fileName, ".")
	suffix := temp[len(temp)-1]

	var img image.Image
	if suffix == "png" {
		img, err = png.Decode(bytes.NewReader(file))
	} else if suffix == "jpg" {
		img, err = jpeg.Decode(bytes.NewReader(file))
	}
	if err != nil {
		return nil, err
	}
	return img, nil
}
