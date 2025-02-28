package main

import (
	_ "embed"
	"image"
	"os"

	"github.com/ddkwork/ux/widget/material"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"github.com/ddkwork/ux/giosvg"
)

//go:generate go run github.com/ddkwork/ux/giosvg/cmd/svggen -i "." -o "./school-bus.go" -pkg "main"

func init() {
	// os.Setenv("GIORENDERER", "forcecompute")
}

// Thanks to Freepik from Flaticon Licensed by Creative Commons 3.0 for the example icons shown below.

//go:embed  school-bus.svg
var bus []byte

func main() {
	data := bus
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Gio"))
		defer w.Perform(system.ActionClose)

		w.Option(app.Title("Gio"))
		defer w.Perform(system.ActionClose)

		vector, err := giosvg.NewVector(data)
		if err != nil {
			panic(err)
		}

		iconRuntime := giosvg.NewIcon(vector)
		iconGenerated := giosvg.NewIcon(VectorSchoolBus)

		w.Event()
		th := material.NewTheme()
		th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
		var ops op.Ops
		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				os.Exit(0)
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				gtx.Constraints.Max.X = gtx.Constraints.Max.X / 2
				gtx.Constraints.Min = gtx.Constraints.Max

				layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min = image.Point{} // Keep aspect ratio.
					return iconRuntime.Layout(gtx)
				})

				offset := op.Offset(image.Point{X: gtx.Constraints.Max.X}).Push(gtx.Ops)
				layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min = image.Point{} // Keep aspect ratio.
					return iconGenerated.Layout(gtx)
				})
				offset.Pop()
				gtx.Execute(op.InvalidateCmd{})
				e.Frame(&ops)
			}
		}
	}()
	app.Main()
}
