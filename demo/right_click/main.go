package main

import (
	"github.com/ddkwork/ux"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
)

func main() {
	go func() {
		w := new(app.Window)
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	p := New()
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			ux.BackgroundDark(gtx)
			p.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}
