package ux

import (
	"gioui.org/app"
	"gioui.org/op"
	"github.com/ddkwork/golibrary/mylog"
	"os"
	"testing"
)

func TestNewPopupTest(t *testing.T) {
	w := new(app.Window)
	p := NewPopupTest(100)
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			mylog.CheckIgnore(e.Err)
			os.Exit(0)
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			BackgroundDark(gtx)
			p.Layout(gtx, nil)
			e.Frame(gtx.Ops)
		}
	}
	app.Main()
}
