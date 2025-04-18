package icon

import (
	"image"
	"image/color"

	"gioui.org/op/paint"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/ddkwork/ux/giosvg"
)

func Layout(gtx layout.Context, icon any, color color.NRGBA, size unit.Dp) layout.Dimensions {
	sizeDp := gtx.Dp(size)
	if icon != nil {
		gtx.Constraints.Min = image.Point{X: sizeDp}
		switch v := icon.(type) {
		case *giosvg.Icon:
			v.Layout(gtx)
		case *widget.Icon:
			v.Layout(gtx, color)
		case *widget.Image:
			v.Layout(gtx)
		case image.Image:
			img := &widget.Image{
				Src:      paint.NewImageOp(v),
				Fit:      widget.Unscaled,
				Position: layout.Center,
				Scale:    1.0, // todo 测试按钮图标和层级图标
			}
			img.Layout(gtx)
		}
	}
	return layout.Dimensions{
		Size: image.Point{X: sizeDp, Y: sizeDp},
	}
}
