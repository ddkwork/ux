// SPDX-License-Identifier: Unlicense OR MIT

package component

import (
	"image"
	"log"
	"math"
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
	LongPressDuration time.Duration //长按激活的时间
	activationTimer   *time.Timer   //用于检测长按的定时器
	lastUpdate        time.Time     //上次更新的时间，用于避免重复处理事件
	position          f32.Point     //上下文菜单的位置
	dims              D             //上下文菜单的尺寸
	active            bool          //表示上下文菜单是否显示
	startedActive     bool          //用于记录是否在当前帧开始时菜单是激活状态
	justActivated     bool          //标记菜单是否刚刚被激活
	justDismissed     bool          //标记菜单是否刚刚被隐藏
	// Activation is the pointer Buttons within the context area
	// that trigger the presentation of the contextual widget. If this
	// is zero, it will default to pointer.ButtonSecondary.
	Activation pointer.Buttons //触发菜单显示的指针按钮（默认为右键）
	// AbsolutePosition will position the contextual widget in the
	// relative to the position of the context area instead of relative
	// to the position of the click event that triggered the activation.
	// This is useful for controls (like button-activated menus) where
	// the contextual content should not be precisely attached to the
	// click position, but should instead be attached to the button.
	AbsolutePosition bool //是否将菜单固定在 ContextArea 的位置，而不是点击位置,下拉菜单，弹出菜单，顶部菜单栏等需要固定位置，表格右键菜单则需要点哪弹哪，不要固定位置，应该根据右击按下的位置计算相对坐标给弹出区域
	// PositionHint tells the ContextArea the closest edge/corner of the
	// window to where it is being used in the layout. This helps it to
	// position the contextual widget without it overflowing the edge of
	// the window.
	PositionHint layout.Direction //提示 ContextArea 在窗口中的位置，用于调整菜单显示以避免溢出窗口边缘,todo bug
}

func (r *ContextArea) Activate(p f32.Point) {
	r.active = true
	r.justActivated = true
	r.position = p

	if !r.AbsolutePosition {
		r.position = p
	}
}

// Update performs event processing for the context area but does not lay it out.
// It is automatically invoked by Layout() if it has not already been called during
// a given frame.
func (r *ContextArea) Update(gtx C) {
	if gtx.Now == r.lastUpdate {
		return
	}
	r.lastUpdate = gtx.Now
	if r.Activation == 0 {
		r.Activation = pointer.ButtonSecondary
	}
	suppressionTag := &r.active
	dismissTag := &r.dims

	r.startedActive = r.active
	// 处理指针事件
	// Summon the contextual widget if the area recieved a secondary click.
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: r,
			Kinds:  pointer.Press | pointer.Release,
		})
		if !ok {
			break
		}
		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}
		if r.active {
			// 检查是否应关闭菜单
			// Check whether we should dismiss menu.
			if e.Buttons.Contain(pointer.ButtonPrimary) {
				clickPos := e.Position.Sub(r.position)
				min := f32.Point{}
				max := f32.Point{
					X: float32(r.dims.Size.X),
					Y: float32(r.dims.Size.Y),
				}
				if !(clickPos.X > min.X && clickPos.Y > min.Y && clickPos.X < max.X && clickPos.Y < max.Y) {
					r.Dismiss()
				}
			}
		}
		if (e.Buttons.Contain(pointer.ButtonPrimary) && e.Kind == pointer.Press) ||
			(e.Source == pointer.Touch && e.Kind == pointer.Press) {
			if r.activationTimer != nil {
				r.activationTimer.Stop()
			}
			if r.LongPressDuration == 0 {
				r.LongPressDuration = 500 * time.Millisecond
			}
			r.activationTimer = time.AfterFunc(r.LongPressDuration, func() {
				r.Activate(e.Position)
			})
		}
		if e.Kind == pointer.Release {
			if r.activationTimer != nil {
				r.activationTimer.Stop()
			}
		}
		if e.Buttons.Contain(r.Activation) && e.Kind == pointer.Press {
			r.Show()
			r.Activate(e.Position)
		}
	}

	// 处理外部点击事件以关闭菜单
	// Dismiss the contextual widget if the user clicked outside of it.
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: suppressionTag,
			Kinds:  pointer.Press,
		})
		if !ok {
			break
		}
		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}
		if e.Kind == pointer.Press {
			r.Dismiss()
		}
	}

	// 处理内部释放事件以关闭菜单
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
	suppressionTag := &r.active

	var contextual op.CallOp
	if r.active || r.startedActive {
		// Render if the layout started as active to ensure that widgets
		// within the contextual content get to update their state in reponse
		// to the event that dismissed the contextual widget.
		contextual = func() op.CallOp {
			macro := op.Record(gtx.Ops)
			r.dims = w(gtx)
			//dims = r.dims
			return macro.Stop()
		}()
	}

	if r.active {
		// 调整菜单位置以避免溢出窗口边缘
		if int(r.position.X)+r.dims.Size.X > gtx.Constraints.Max.X {
			if newX := int(r.position.X) - r.dims.Size.X; newX < 0 {
				switch r.PositionHint {
				case layout.E, layout.NE, layout.SE:
					r.position.X = float32(gtx.Constraints.Max.X - r.dims.Size.X)
				case layout.W, layout.NW, layout.SW:
					r.position.X = 0
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
					r.position.Y = 0 // 确保菜单顶部不超出窗口顶部
				default:
					r.position.Y = float32(gtx.Constraints.Max.Y - r.dims.Size.Y)
				}
			} else {
				r.position.Y = float32(newY)
			}
		}
		// 确保菜单顶部不超出窗口最小Y限制
		if r.position.Y < float32(gtx.Constraints.Min.Y) {
			r.position.Y = float32(gtx.Constraints.Min.Y)
		}
		log.Println(r.position)

		// 创建透明的遮罩层以阻止对菜单下方内容的输入
		// Lay out a transparent scrim to block input to things beneath the
		// contextual widget.
		suppressionScrim := func() op.CallOp {
			macro2 := op.Record(gtx.Ops)
			pr := clip.Rect(image.Rectangle{Min: image.Point{-1e6, -1e6}, Max: image.Point{1e6, 1e6}})
			stack := pr.Push(gtx.Ops)
			event.Op(gtx.Ops, suppressionTag)
			stack.Pop()
			return macro2.Stop()
		}()
		op.Defer(gtx.Ops, suppressionScrim)

		// 布局上下文菜单
		// Lay out the contextual widget itself.
		pos := image.Point{
			X: int(math.Round(float64(r.position.X))),
			Y: int(math.Round(float64(r.position.Y))),
		}
		macro := op.Record(gtx.Ops)
		op.Offset(pos).Add(gtx.Ops)
		contextual.Add(gtx.Ops)

		// 创建遮罩层以检测完成与菜单的交互
		// Lay out a scrim on top of the contextual widget to detect
		// completed interactions with it (that should dismiss it).
		pt := pointer.PassOp{}.Push(gtx.Ops)
		stack := clip.Rect(image.Rectangle{Max: r.dims.Size}).Push(gtx.Ops)
		event.Op(gtx.Ops, &r.dims)

		stack.Pop()
		pt.Pop()
		contextual = macro.Stop()
		op.Defer(gtx.Ops, contextual)
	}

	// 捕获上下文区域的指针事件
	// Capture pointer events in the contextual area.
	defer pointer.PassOp{}.Push(gtx.Ops).Pop()
	defer clip.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Push(gtx.Ops).Pop()
	event.Op(gtx.Ops, r)

	return D{Size: gtx.Constraints.Min}
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
// the last call to Activated.
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
	r.active = true
	r.justActivated = true
}
