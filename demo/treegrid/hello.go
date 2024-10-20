package main

import (
	"os"

	"gioui.org/layout"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/ux"

	"gioui.org/app"
	"gioui.org/op"
)

func main() {
	go func() {
		w := new(app.Window)
		mylog.Check(loop(w))
		os.Exit(0)
	}()
	app.Main()
}

type (
	C = layout.Context
	D = layout.Dimensions
)

func loop(w *app.Window) error {
	// th := material.NewTheme()
	// th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	t := NewTreeTable()
	t.SetRootRows(TestRootRows)
	SetParents(TestRootRows, nil)
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			ux.BackgroundDark(gtx)
			t.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}
