package main

import (
	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/ddkwork/ux"
	"github.com/ddkwork/ux/widget/material"
	"image"
	"image/color"
	"os"
)

func main() {
	w := new(app.Window)
	var ops op.Ops
	page := new(Page)
	go func() {
		for {
			switch e := w.Event().(type) {
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				ux.BackgroundDark(gtx)
				page.Layout(gtx)
				e.Frame(gtx.Ops)
			case app.DestroyEvent:
				os.Exit(0)
			}
		}
	}()
	app.Main()
}

type Page struct {
	dialog      layout.Widget
	dialogCover widget.Clickable
	submit      widget.Clickable
	content     widget.Clickable
}

func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	// if click on "Open Dialog:
	th := ux.NewTheme()
	if p.submit.Clicked(gtx) {
		// Your dialog:
		p.dialog = func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X /= 5
			gtx.Constraints.Max.Y = gtx.Dp(200)

			defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Max}, 24).Push(gtx.Ops).Pop()
			paint.ColorOp{Color: color.NRGBA{R: 255, G: 255, B: 255, A: 255}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			// Click blocker:
			//pointer.InputOp{Tag: &p, Types: pointer.Press | pointer.Release}.Add(gtx.Ops)
			//gtx.Event(pointer.Filter{
			//	Target: &p,
			//	Kinds:  pointer.Press | pointer.Release,
			//})
			event.Op(gtx.Ops, &p)

			gtx.Constraints.Min = gtx.Constraints.Max
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.Y = 0
				paint.ColorOp{Color: color.NRGBA{0, 0, 0, 255}}.Add(gtx.Ops)
				return material.Button(th.Theme, &p.content, "Text on dialog").Layout(gtx)
			})
		}
	}

	// if click on the background, close the dialog
	if p.dialogCover.Clicked(gtx) {
		p.dialog = nil // Remove the dialog
	}

	// overlay dialog
	if p.dialog != nil {
		// get size of the dialog  widget
		dialogRec := op.Record(gtx.Ops)
		dims := p.dialog(gtx)
		dialogCall := dialogRec.Stop()

		// add padding and background opacity
		overlay := op.Record(gtx.Ops)

		p.dialogCover.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			paint.ColorOp{Color: color.NRGBA{0, 0, 0, 128}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: gtx.Constraints.Max}
		})

		offset := op.Offset(image.Point{
			X: (gtx.Constraints.Max.X - dims.Size.X) / 2,
			Y: (gtx.Constraints.Max.Y - dims.Size.Y) / 2,
		}).Push(gtx.Ops)
		dialogCall.Add(gtx.Ops)
		offset.Pop()

		overlayCall := overlay.Stop()

		op.Defer(gtx.Ops, overlayCall)
	}

	// main page
	{
		//gtx := gtx // copy
		gtx.Constraints.Max.X /= 2
		gtx.Constraints.Min = gtx.Constraints.Max

		// left
		layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return material.H5(th.Theme, "Some sidebar text here").Layout(gtx)
		})

		// right
		right := op.Offset(image.Pt(gtx.Constraints.Max.X, 0)).Push(gtx.Ops)
		layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return p.submit.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.Y = gtx.Dp(30)
						gtx.Constraints.Min.Y = gtx.Constraints.Max.Y

						defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
						paint.ColorOp{Color: color.NRGBA{255, 0, 0, 255}}.Add(gtx.Ops)
						paint.PaintOp{}.Add(gtx.Ops)

						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.Y = 0
							paint.ColorOp{Color: color.NRGBA{0, 0, 0, 255}}.Add(gtx.Ops)
							//return widget.Label{}.Layout(gtx, shaper, text.Font{}, 14, "Open Dialog")
							return material.Button(th.Theme, &p.submit, "Open Dialog").Layout(gtx)
						})
					})
				}),
			)
		})
		right.Pop()
	}
	return layout.Dimensions{Size: gtx.Constraints.Max}
}
