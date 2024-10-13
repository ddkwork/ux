package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/ddkwork/ux"
	"log"
	"os"
)

type CustomView struct {
	Title string
}

func (c *CustomView) Layout(gtx ux.Gtx) layout.Dimensions {
	return func(gtx ux.Gtx) ux.Dim {
		return material.Body1(ux.ThemeDefault().Theme, c.Title).Layout(gtx)
	}(gtx)
}

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
	//resizer := ux.Resize{}
	cust1 := CustomView{Title: "Widget One Widget One Widget One Widget One Widget One Widget One Widget One Widget One"}
	cust2 := CustomView{Title: "Widget Two Widget Two Widget Two Widget Two Widget Two Widget Two Widget Two Widget Two "}
	cust3 := CustomView{Title: "Widget Three Widget Three Widget Three Widget Three Widget Three Widget Three Widget Three"}
	cust4 := CustomView{Title: "Widget Four Widget Four Widget Four Widget Four Widget Four Widget Four Widget Four "}
	r1 := ux.Resizable{Widget: cust1.Layout}
	r2 := ux.Resizable{Widget: cust2.Layout}
	r3 := ux.Resizable{Widget: cust3.Layout}
	r4 := ux.Resizable{Widget: cust4.Layout}

	resizables := []*ux.Resizable{&r1, &r2, &r3, &r4}
	resizer := ux.NewResizeWidget(layout.Horizontal, resizables...)

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
			resizer.Layout(gtx)
			//resizer.Layout(gtx, cust2.Layout, nil)
			//resizer.Layout(gtx, cust3.Layout, nil)
			//resizer.Layout(gtx, cust4.Layout, nil)
			e.Frame(gtx.Ops)
		}
	}
}
