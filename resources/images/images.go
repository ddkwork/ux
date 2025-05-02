package images

import (
	"bytes"
	"image"
	"image/color"
	"strings"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/ux/internal/giosvg"
	_ "github.com/ddkwork/ux/resources/images/ico"
)

//type Type interface {
//	*widget.Icon | *widget.Image | *giosvg.Icon | []byte
//}

func Layout(gtx layout.Context, b []byte, color color.NRGBA, size unit.Dp) layout.Dimensions {
	const defaultIconSize = unit.Dp(24)
	sz := gtx.Constraints.Min.X
	if sz == 0 {
		sz = gtx.Dp(defaultIconSize)
	}
	// sizeDp := gtx.Dp(size)
	sizeDp := gtx.Constraints.Constrain(image.Pt(sz, sz))
	if b == nil {
		return layout.Dimensions{
			Size: sizeDp,
		}
	}
	gtx.Constraints.Min = sizeDp
	gtx.Constraints.Max = sizeDp
	if icon, err := widget.NewIcon(b); err == nil {
		return icon.Layout(gtx, color)
	}
	if v, err := giosvg.NewVector(b); err == nil {
		return giosvg.NewIcon(v).Layout(gtx)
	}
	img, _ := mylog.Check3(image.Decode(bytes.NewReader(b))) //_ "github.com/ddkwork/ux/resources/images/ico" 通过init注册ico解码器
	w := &widget.Image{
		Src:      paint.NewImageOp(img),
		Fit:      widget.ScaleDown,
		Position: layout.Center,
		Scale:    0,
	}
	return w.Layout(gtx)
}

func IsSVG(b []byte) bool {
	if len(b) < 5 {
		return false
	}
	switch strings.ToLower(string(b[0:5])) {
	case "<!doc", "<?xml", "<svg ": // and extension is svg
		return true
	}
	return false
}

// Unscaled（不缩放）：这个值表示不改变小部件的原始尺寸。
// Contain（包含）：按照小部件的原始宽高比尽可能大地放大小部件，使其不会被裁剪（整个小部件仍然可见）。
// Cover（覆盖）：按照小部件的原始宽高比尽可能大地缩小小部件，使其能够覆盖整个布局区域（部分小部件可能会被裁剪）。
// ScaleDown（缩小）：如果小部件的尺寸超过布局区域，则按照原始宽高比缩小小部件，使其完全位于布局区域内；如果小部件尺寸小于或等于布局区域，则保持不变。
// Fill（填充）：拉伸或缩小小部件以适应布局区域，不保持小部件的原始宽高比。
