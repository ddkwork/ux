package main

import (
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/ux"
)

func main() {
	go func() {
		w := new(app.Window)
		mylog.Check(loop(w))
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	c := ux.Calendar{}
	c.Inset = layout.UniformInset(unit.Dp(16))
	c.FirstDayOfWeek = time.Monday

	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			ux.BackgroundDark(gtx)
			c.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}
