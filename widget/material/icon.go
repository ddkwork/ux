package material

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream/ico"
	"golang.org/x/image/bmp"
)

func LoadImage(fileName string, data []byte) image.Image { // 一般而言，我们使用embed
	var b []byte
	if data != nil {
		b = data
	} else {
		b = mylog.Check2(os.ReadFile(fileName))
	}
	var img image.Image
	switch filepath.Ext(fileName) {
	case ".png":
		img = mylog.Check2(png.Decode(bytes.NewReader(b)))
	case ".jpg", ".jpeg":
		img = mylog.Check2(jpeg.Decode(bytes.NewReader(b)))
	case ".gif":
		img = mylog.Check2(gif.Decode(bytes.NewReader(b)))
	case ".ico":
		img = mylog.Check2(ico.Decode(bytes.NewReader(b)))
	case ".bmp":
		img = mylog.Check2(bmp.Decode(bytes.NewReader(b)))
	// case ".svg": //svg的话是giosvg直接解码元数据实现layout方法渲染
	default:
		mylog.Check("unsupported image format")
	}
	return img
}

func ImageButton(th *Theme, button *widget.Clickable, fileName string, data []byte, description string) IconButtonStyle {
	return IconButtonStyle{
		Background: th.Palette.ContrastBg,
		Color:      th.Palette.ContrastFg,
		Icon: &widget.Image{
			Src:      paint.NewImageOp(LoadImage(fileName, data)),
			Fit:      widget.Unscaled,
			Position: layout.Center,
			Scale:    1.0, // todo 测试按钮图标和层级图标
		},
		Size:        24,
		Inset:       layout.UniformInset(12),
		Button:      button,
		Description: description,
	}
}
