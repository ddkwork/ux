package main

import (
	"fmt"
	"github.com/ddkwork/ux"
	"image"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// Anchor is an opaque reference to a global coordinate position.
// It can be provided to methods in this package as a reference
// to a global coordinate.
type Anchor struct {
	point f32.Point
}

// AnchorFrom wraps an f32.Point within an Anchor, preventing the
// coordinates within from being used in any way other than determining
// an offset using the OffsetWithin method.
func AnchorFrom(point f32.Point) Anchor {
	return Anchor{point}
}

// String is provided for debugging purposes.
func (a Anchor) String() string {
	return fmt.Sprintf("anchor (%f,%f)", a.point.X, a.point.Y)
}

// OffsetWithin returns an offset that will allow a widget of size contentSize
// to be rendered within the provided bounds. The offset is as close as possible
// to the coordinates wrapped within the
func (a Anchor) OffsetWithin(contentSize, bounds f32.Point) f32.Point {
	offset := a.point
	if contentSize.X+a.point.X > bounds.X {
		offset.X = bounds.X - contentSize.X
	}
	if contentSize.Y+a.point.Y > bounds.Y {
		offset.Y = bounds.Y - contentSize.Y
	}
	return offset
}

type Overlay struct {
	items []overlayItem
}

type overlayItem struct {
	Anchor
	layout.Widget
}

func (o *Overlay) LayoutAt(anchor Anchor, widget layout.Widget) {
	o.items = append(o.items, overlayItem{Anchor: anchor, Widget: widget})
}

func (o *Overlay) Layout(gtx layout.Context) layout.Dimensions {
	for _, item := range o.items {
		macro := op.Record(gtx.Ops)
		dims := item.Widget(gtx)
		call := macro.Stop()
		offset := item.OffsetWithin(layout.FPt(dims.Size), layout.FPt(gtx.Constraints.Max))
		func(item overlayItem) {
			defer op.TransformOp{}.Push(gtx.Ops).Pop()
			op.Offset(image.Point{X: int(offset.X), Y: int(offset.Y)}).Add(gtx.Ops)
			call.Add(gtx.Ops)
		}(item)
	}
	o.items = o.items[:0]
	return layout.Dimensions{Size: gtx.Constraints.Max}
}

// RightClickArea wraps a widget and provides a right-click context menu
type RightClickArea struct {
	// Content is the actual right-clickable widget
	Content layout.Widget
	// Menu is the widget that should be rendered as a right-click context menu
	Menu layout.Widget
	*Anchor
	*Overlay
	leftPressed *pointer.ID
}

// LayoutUnderlay creates an invisible layer to listen for click events
// across the entire graphics context. It sizes itself to be the maximum
// size of the context, and should be anchored at the origin.
func (r *RightClickArea) LayoutUnderlay(gtx C) D {
	pt := pointer.PassOp{}.Push(gtx.Ops)
	stack := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
	event.Op(gtx.Ops, r)
	stack.Pop()
	pt.Pop()
	return D{Size: gtx.Constraints.Max}
}

// CloseMenu cancels the display of the context menu.
func (r *RightClickArea) CloseMenu() {
	r.leftPressed = nil
	r.Anchor = nil
}

// Layout renders the clickable area and configures its overlay.
func (r *RightClickArea) Layout(gtx C) D {
	event.Op(gtx.Ops, r)
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
				anchor := AnchorFrom(e.Position)
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
				r.leftPressed = &e.PointerID
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
		r.Overlay.LayoutAt(Anchor{}, r.LayoutUnderlay)
		r.Overlay.LayoutAt(*r.Anchor, r.Menu)
	}
	dims := r.Content(gtx)
	pointer.PassOp{}.Push(gtx.Ops).Pop()
	clip.Rect(image.Rectangle{Max: dims.Size}).Push(gtx.Ops).Pop()
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
		th  = material.NewTheme()
		ops op.Ops
	)
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	rows := make([]RightClickArea, 0, 30)
	for range 30 {
		var (
			overlay               Overlay
			btn, rBtn, gBtn, bBtn widget.Clickable
		)
		rows = append(rows,
			RightClickArea{
				Overlay: &overlay,
				Content: func(gtx C) D {
					return material.Button(th, &btn, "Reset").Layout(gtx)
				},
				Menu: func(gtx C) D {
					gtx.Constraints.Min = image.Point{}
					gtx.Constraints.Max.X = gtx.Dp(200)
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
			},
		)
	}

	list := widget.List{
		Scrollbar: widget.Scrollbar{},
		List: layout.List{
			Axis:        layout.Vertical,
			ScrollToEnd: false,
			Alignment:   0,
			Position:    layout.Position{},
		},
	}

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			ux.BackgroundDark(gtx)
			material.List(th, &list).Layout(gtx, len(rows), func(gtx layout.Context, index int) layout.Dimensions {
				item := rows[index]
				return layout.Stack{}.Layout(gtx,
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return item.Layout(gtx)
					}),
					layout.Expanded(func(gtx layout.Context) layout.Dimensions {
						return item.Overlay.Layout(gtx)
					}),
				)
			})
			e.Frame(gtx.Ops)
		}
	}
}
