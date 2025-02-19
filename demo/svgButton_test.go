package main_test

import (
	"testing"

	"gioui.org/layout"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/ux"
)

func TestGoogleDummyButton_Layout(t *testing.T) {
	CircledChevronRightButton := ux.NewSVGButton("", ux.Svg2Icon([]byte(ux.CircledChevronRight)), func() {
		mylog.Info("svg button clicked")
	})
	CircledChevronDownButton := ux.NewSVGButton("", ux.Svg2Icon([]byte(ux.CircledChevronDown)), func() {
		mylog.Info("svg button clicked")
	})

	// SaveScreenshot(CircledChevronDownButton.Layout)
	ux.SaveScreenshot(func(gtx layout.Context) layout.Dimensions {
		list := layout.List{
			Axis:        layout.Vertical,
			ScrollToEnd: false,
			Alignment:   0,
			Position:    layout.Position{},
		}
		return list.Layout(gtx, 2, func(gtx layout.Context, index int) layout.Dimensions {
			switch index {
			case 0:
				return CircledChevronRightButton.Layout(gtx)
			}
			return CircledChevronDownButton.Layout(gtx)
		})
	})
}
