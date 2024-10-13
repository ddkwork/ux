package ux

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateTextWidth(t *testing.T) {
	SaveScreenshot(func(gtx layout.Context) layout.Dimensions {
		width := CalculateTextWidth(gtx, "Sub Row 1 (9)")
		assert.Equal(t, unit.Dp(131), width) //but saved is 216
		return layout.Dimensions{}
	})
}
