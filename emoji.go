package ux

import (
	"math"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/ddkwork/ux/widget/material"
)

func layoutEmoji(gtx layout.Context) layout.Dimensions {
	var sel widget.Selectable
	message := "🥳🧁🍰🎁🎂🎈🎺🎉🎊\n📧〽️🧿🌶️🔋\n😂❤️😍🤣😊\n🥺🙏💕😭😘\n👍😅👏"
	var customTruncator widget.Bool
	var maxLines widget.Float
	maxLines.Value = 0

	const (
		minLinesRange = 1
		maxLinesRange = 5
	)

	inset := layout.UniformInset(5)

	l := material.H4(th.Theme, message)
	if customTruncator.Value {
		l.Truncator = "cont..."
	} else {
		l.Truncator = ""
	}
	l.MaxLines = minLinesRange + int(math.Round(float64(maxLines.Value)*(maxLinesRange-minLinesRange)))
	l.State = &sel
	return inset.Layout(gtx, l.Layout)
}
