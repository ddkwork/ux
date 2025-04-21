package main

import (
	"bytes"
	"image/png"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/gorpher/gowin32"
)

func ExtractIcon2Png(filename string) []byte {
	img, e := gowin32.ExtractPrivateExtractIcons(filename, 128, 128)
	if e != nil {
		return installerJpg //The specified resource type cannot be found in the image file. 有些cli程序是没有图标的
	}
	b := new(bytes.Buffer)
	mylog.Check(png.Encode(b, img))
	return b.Bytes()
}
