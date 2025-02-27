package main_test

import (
	"testing"

	"gioui.org/layout"
	"github.com/ddkwork/ux"
)

func TestButton_Layout(t *testing.T) {
	ux.SaveScreenshot(ux.NewButton("杨凯隐", nil).Layout)
}

func TestNewButtonAnimation(t *testing.T) {
	ux.SaveScreenshot(func(gtx layout.Context) layout.Dimensions {
		animation := ux.NewButtonAnimation("xxxxxxxxxxxxxxxxxxxxxxxxxx", ux.IconAdd, func(gtx layout.Context) {
			print("xxxxx")
		})
		animation.SetLoading(true)
		return animation.Layout(gtx)
	})
}
