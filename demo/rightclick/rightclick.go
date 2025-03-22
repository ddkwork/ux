package main

import (
	"demo/rightclick/anchor"
	"gioui.org/io/event"
	"gioui.org/op/clip"
	"gioui.org/text"
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// RightClickArea wraps a widget and provides a right-click context menu
type RightClickArea struct {
	// Content is the actual right-clickable widget
	Content layout.Widget
	// Menu is the widget that should be rendered as a right-click context menu
	Menu layout.Widget
	*anchor.Anchor
	*Overlay
	leftPressed *pointer.ID
}

// LayoutUnderlay creates an invisible layer to listen for click events
// across the entire graphics context. It sizes itself to be the maximum
// size of the context, and should be anchored at the origin.
func (r *RightClickArea) LayoutUnderlay(gtx C) D {
	//defer op.Push(gtx.Ops).Pop()
	defer op.TransformOp{}.Push(gtx.Ops).Pop()
	pointer.PassOp{}.Push(gtx.Ops)
	//pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
	clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
	//pointer.InputOp{
	//	Tag:   &r.leftPressed,
	//	Types: pointer.Press | pointer.Release,
	//}.Add(gtx.Ops)

	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: &r.leftPressed,
			Kinds:  pointer.Press | pointer.Release,
		})
		if !ok {
			break
		}
		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}
		if e.Kind == pointer.Press {
			//r.Dismiss()//todo
			event.Op(gtx.Ops, r) //?
		}
	}

	return D{Size: gtx.Constraints.Max}
}

// CloseMenu cancels the display of the context menu.
func (r *RightClickArea) CloseMenu() {
	r.leftPressed = nil
	r.Anchor = nil
}

// Layout renders the clickable area and configures its overlay.
func (r *RightClickArea) Layout(gtx C) D {
	//defer op.Push(gtx.Ops).Pop()

	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: r,
			Kinds:  pointer.Press | pointer.Release | pointer.Drag,
		})
		if !ok {
			break
		}
		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}
		if e.Buttons.Contain(pointer.ButtonSecondary) {
			switch e.Kind {
			case pointer.Press, pointer.Drag:
				anchor := anchor.AnchorFrom(e.Position)
				r.Anchor = &anchor
				log.Print(anchor)
			case pointer.Cancel:
				r.Anchor = nil
			}
		}
	}

	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: &r.leftPressed,
			Kinds:  pointer.Press | pointer.Release | pointer.Drag,
		})
		if !ok {
			break
		}
		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}
		switch e.Kind {
		case pointer.Press, pointer.Drag:
			if e.Buttons.Contain(pointer.ButtonPrimary) {
				id := e.PointerID
				r.leftPressed = &id
			}
		case pointer.Release, pointer.Cancel:
			if r.leftPressed != nil && e.PointerID == *r.leftPressed {
				log.Println("left", e)
				r.Anchor = nil
				r.leftPressed = nil
			}
		}

	}
	if r.Anchor != nil {
		r.Overlay.LayoutAt(anchor.Anchor{}, r.LayoutUnderlay)
		r.Overlay.LayoutAt(*r.Anchor, r.Menu)
	}
	dims := r.Content(gtx)
	pointer.PassOp{}.Push(gtx.Ops)
	clip.Rect(image.Rectangle{Max: dims.Size}).Push(gtx.Ops)
	return dims
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
	var (
		th                    = material.NewTheme()
		btn, rBtn, gBtn, bBtn widget.Clickable
		ops                   op.Ops
		overlay               Overlay
		areaColor             = color.NRGBA{A: 255}
		rca                   = RightClickArea{
			Overlay: &overlay,
			Content: func(gtx C) D {
				btn := material.Button(th, &btn, "Reset")
				btn.Background = areaColor
				return btn.Layout(gtx)
			},
			Menu: func(gtx C) D {
				gtx.Constraints.Min = image.Point{}
				gtx.Constraints.Max.X = gtx.Dp(unit.Dp(200))
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return material.Button(th, &rBtn, "Redden").Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return material.Button(th, &bBtn, "Bluify").Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return material.Button(th, &gBtn, "Greenenate").Layout(gtx)
					}),
				)
			},
		}
	)
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			menuClicked := false
			if rBtn.Clicked(gtx) {
				menuClicked = true
				areaColor.R += 64
			}
			if gBtn.Clicked(gtx) {
				menuClicked = true
				areaColor.G += 64
			}
			if bBtn.Clicked(gtx) {
				menuClicked = true
				areaColor.B += 64
			}
			if menuClicked {
				rca.CloseMenu()
			}
			if btn.Clicked(gtx) {
				areaColor.R = 0
				areaColor.G = 0
				areaColor.B = 0
			}
			layout.Center.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return material.Body1(th, "Right-click or click this button:").Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Max.X /= 2
						gtx.Constraints.Max.Y /= 2
						gtx.Constraints.Min = gtx.Constraints.Max
						return rca.Layout(gtx)
					}),
				)
			})
			overlay.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}
