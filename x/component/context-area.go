// SPDX-License-Identifier: Unlicense OR MIT

package component

import (
	"image"
	"math"
	"runtime"
	"time"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
)

// ContextArea is a region of the UI that responds to certain
// pointer events by displaying a contextual widget. The contextual
// widget is overlaid using an op.DeferOp. The contextual widget
// can be dismissed by primary-clicking within or outside of it.
type ContextArea struct {
	lastUpdate    time.Time
	position      f32.Point
	dims          D
	active        bool
	startedActive bool
	justActivated bool
	justDismissed bool
	// Activation is the pointer Buttons within the context area
	// that trigger the presentation of the contextual widget. If this
	// is zero, it will default to pointer.ButtonSecondary.
	Activation pointer.Buttons
	// AbsolutePosition will position the contextual widget in the
	// relative to the position of the context area instead of relative
	// to the position of the click event that triggered the activation.
	// This is useful for controls (like button-activated menus) where
	// the contextual content should not be precisely attached to the
	// click position, but should instead be attached to the button.
	AbsolutePosition bool
	// PositionHint tells the ContextArea the closest edge/corner of the
	// window to where it is being used in the layout. This helps it to
	// position the contextual widget without it overflowing the edge of
	// the window.
	PositionHint      layout.Direction
	activationTimer   *time.Timer
	LongPressDuration time.Duration
}

func (r *ContextArea) Activate(p f32.Point) {
	r.active = true
	r.justActivated = true
	if !r.AbsolutePosition {
		r.position = p
	}
}

// Update performs event processing for the context area but does not lay it out.
// It is automatically invoked by Layout() if it has not already been called during
// a given frame.
func (r *ContextArea) Update(gtx layout.Context) {
	if gtx.Now == r.lastUpdate {
		return
	}
	r.lastUpdate = gtx.Now
	if r.Activation == 0 {
		r.Activation = pointer.ButtonSecondary
	}
	dismissTag := &r.dims

	r.startedActive = r.active
	// Summon the contextual widget if the area received a secondary click.
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: r,
			Kinds:  pointer.Press | pointer.Release | pointer.Scroll,
		})
		if !ok {
			break
		}
		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}
		if e.Kind == pointer.Scroll {
			//mylog.Success("scroll event", e.Scroll.X, e.Scroll.Y)
			r.Dismiss()
			break
		}
		if r.active {
			// Check whether we should dismiss menu.
			if e.Buttons.Contain(pointer.ButtonPrimary) {
				clickPos := e.Position.Sub(r.position)
				minPtr := f32.Point{}
				maxPtr := f32.Point{
					X: float32(r.dims.Size.X),
					Y: float32(r.dims.Size.Y),
				}
				if !(clickPos.X > minPtr.X && clickPos.Y > minPtr.Y && clickPos.X < maxPtr.X && clickPos.Y < maxPtr.Y) {
					r.Dismiss()
				}
			}
		}

		//todo  按系统精确和操作习惯精确设置，添加日志打印
		if (e.Buttons.Contain(pointer.ButtonPrimary) && e.Kind == pointer.Press) ||
			(e.Source == pointer.Touch && e.Kind == pointer.Press && runtime.GOOS != "android") {
			if r.activationTimer != nil {
				r.activationTimer.Stop()
			}
			if r.LongPressDuration == 0 {
				r.LongPressDuration = 500 * time.Millisecond
			}
			r.activationTimer = time.AfterFunc(r.LongPressDuration, func() {
				if runtime.GOOS == "android" {
					r.Activate(e.Position)
				}
			})
		}
		if e.Kind == pointer.Release {
			if r.activationTimer != nil {
				r.activationTimer.Stop()
			}
		}
		if e.Buttons.Contain(r.Activation) && e.Kind == pointer.Press {
			r.Activate(e.Position)
		}
	}

	// Dismiss the contextual widget if the user released a click within it.
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: dismissTag,
			Kinds:  pointer.Release,
		})
		if !ok {
			break
		}
		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}
		if e.Kind == pointer.Release {
			r.Dismiss()
		}
	}
}

// Layout renders the context area and -- if the area is activated by an
// appropriate gesture -- also the provided widget overlaid using an op.DeferOp.
func (r *ContextArea) Layout(gtx C, w layout.Widget) D {
	r.Update(gtx)
	dismissTag := &r.dims
	dims := D{Size: gtx.Constraints.Min}

	var contextual op.CallOp
	if r.active || r.startedActive {
		// Render if the layout started as active to ensure that widgets
		// within the contextual content get to update their state in repose
		// to the event that dismissed the contextual widget.
		contextual = func() op.CallOp {
			macro := op.Record(gtx.Ops)
			r.dims = w(gtx)
			return macro.Stop()
		}()
	}

	if r.active {
		if int(r.position.X)+r.dims.Size.X > gtx.Constraints.Max.X {
			if newX := int(r.position.X) - r.dims.Size.X; newX < 0 {
				switch r.PositionHint {
				case layout.E, layout.NE, layout.SE:
					r.position.X = float32(gtx.Constraints.Max.X - r.dims.Size.X)
				case layout.W, layout.NW, layout.SW:
					r.position.X = 0
				default:
					panic("unhandled default case")
				}
			} else {
				r.position.X = float32(newX)
			}
		}
		if int(r.position.Y)+r.dims.Size.Y > gtx.Constraints.Max.Y {
			if newY := int(r.position.Y) - r.dims.Size.Y; newY < 0 {
				switch r.PositionHint {
				case layout.S, layout.SE, layout.SW:
					r.position.Y = float32(gtx.Constraints.Max.Y - r.dims.Size.Y)
				case layout.N, layout.NE, layout.NW:
					r.position.Y = 0
				default:
					r.position.Y = float32(gtx.Constraints.Max.Y - r.dims.Size.Y)
				}
			} else {
				r.position.Y = float32(newY)
			}
		}
		// log.Println(r.position)

		// Lay out the contextual widget itself.
		pos := image.Point{
			X: int(math.Round(float64(r.position.X))),
			Y: int(math.Round(float64(r.position.Y))),
		}
		macro := op.Record(gtx.Ops)
		op.Offset(pos).Add(gtx.Ops)
		contextual.Add(gtx.Ops)

		// Lay out a scrim on top of the contextual widget to detect
		// completed interactions with it (that should dismiss it).
		pt := pointer.PassOp{}.Push(gtx.Ops)
		stack := clip.Rect(image.Rectangle{Max: r.dims.Size}).Push(gtx.Ops)
		event.Op(gtx.Ops, dismissTag)

		stack.Pop()
		pt.Pop()
		contextual = macro.Stop()
		op.Defer(gtx.Ops, contextual)
	}

	// Capture pointer events in the contextual area.
	defer pointer.PassOp{}.Push(gtx.Ops).Pop()
	defer clip.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Push(gtx.Ops).Pop()
	event.Op(gtx.Ops, r)

	return dims
}

// Dismiss sets the ContextArea to not be active.
func (r *ContextArea) Dismiss() {
	r.active = false
	r.justDismissed = true
}

// Active returns whether the ContextArea is currently active (whether
// it is currently displaying overlaid content or not).
func (r *ContextArea) Active() bool {
	return r.active
}

// Activated returns true if the context area has become active since
// the last call to Activate.
func (r *ContextArea) Activated() bool {
	defer func() {
		r.justActivated = false
	}()
	return r.justActivated
}

// Dismissed returns true if the context area has been dismissed since
// the last call to Dismissed.
func (r *ContextArea) Dismissed() bool {
	defer func() {
		r.justDismissed = false
	}()
	return r.justDismissed
}

func (r *ContextArea) Show() {
	r.Activate(r.position)
}
