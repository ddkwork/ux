package images

import (
	"bytes"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/ux/internal/giosvg"
	_ "github.com/ddkwork/ux/resources/images/ico"
	"image"
	"image/color"
	"strings"
)

//type Type interface {
//	*widget.Icon | *widget.Image | *giosvg.Icon | []byte
//}

func Layout(gtx layout.Context, b []byte, color color.NRGBA, size unit.Dp) layout.Dimensions {
	if size == 0 {
		size = 18
	}
	sizeDp := gtx.Dp(size)
	if b == nil {
		return layout.Dimensions{
			Size: image.Point{X: sizeDp, Y: sizeDp},
		}
	}
	gtx.Constraints.Min = image.Point{X: sizeDp}
	if icon, err := widget.NewIcon(b); err == nil {
		return icon.Layout(gtx, color)
	}
	if v, err := giosvg.NewVector(b); err == nil {
		return giosvg.NewIcon(v).Layout(gtx)
	}
	img, _ := mylog.Check3(image.Decode(bytes.NewReader(b))) //_ "github.com/ddkwork/ux/resources/images/ico" 通过init注册ico解码器
	w := &widget.Image{
		Src:      paint.NewImageOp(img),
		Fit:      widget.Unscaled,
		Position: layout.Center,
		Scale:    1.0, // todo 测试按钮图标和层级图标
	}
	return w.Layout(gtx)
}

func IsSVG(b []byte) bool {
	if len(b) < 5 {
		return false
	}
	switch strings.ToLower(string(b[0:5])) {
	case "<!doc", "<?xml", "<svg ": //and extension is svg
		return true
	}
	return false
}
