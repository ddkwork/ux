package material

import (
	"embed"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/stream"
)

func ImageButton(th *Theme, button *widget.Clickable, fileName string, fs embed.FS, description string) IconButtonStyle {
	return IconButtonStyle{
		Background: th.Palette.ContrastBg,
		Color:      th.Palette.ContrastFg,
		Icon: &widget.Image{
			Src:      paint.NewImageOp(stream.LoadImage(fileName, fs)),
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
