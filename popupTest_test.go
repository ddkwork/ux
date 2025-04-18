package ux

import (
	"gioui.org/app"
	"gioui.org/op"
	"os"
	"testing"
)

func TestNewPopupTest(t *testing.T) {
	go func() {
		w := new(app.Window)
		p := NewPopupTest(100)
		var ops op.Ops
		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				panic(e.Err)
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				BackgroundDark(gtx)
				p.Layout(gtx, nil)
				e.Frame(gtx.Ops)
			}
		}
		os.Exit(0)
	}()
	app.Main()
}
