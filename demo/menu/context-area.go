// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"image"
	"time"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/widget"
)

// 简化约束类型以便编译
type (
	C = layout.Context
	D = layout.Dimensions
)

type ContextArea struct {
	lastUpdate       time.Time
	position         f32.Point
	dims             layout.Dimensions
	active           bool
	justActivated    bool
	justDismissed    bool
	Activation       pointer.Buttons
	AbsolutePosition bool
	PositionHint     layout.Direction

	// 保存菜单中的所有可点击项
	Clickables []*clickableWidget
}

// 可点击项定义
type clickableWidget struct {
	Area  image.Rectangle
	Click *widget.Clickable
}

func (r *ContextArea) Update(gtx C) {
	if gtx.Now == r.lastUpdate {
		return
	}
	r.lastUpdate = gtx.Now
	if r.Activation == 0 {
		r.Activation = pointer.ButtonSecondary
	}

	// 重置可点击项
	r.Clickables = nil

	// 1. 处理激活事件
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: r,
			Kinds:  pointer.Press,
		})
		if !ok {
			break
		}
		e, ok := ev.(pointer.Event)
		if !ok || !e.Buttons.Contain(r.Activation) {
			continue
		}

		r.active = true
		r.justActivated = true
		if !r.AbsolutePosition {
			r.position = e.Position
		}
	}

	// 2. 处理菜单外点击关闭
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: &r.active,
			Kinds:  pointer.Press,
		})
		if !ok {
			break
		}
		e, ok := ev.(pointer.Event)
		if !ok || e.Kind != pointer.Press {
			continue
		}

		// 检查点击位置是否在任何菜单项上
		clickedOnItem := false
		clickPos := image.Pt(int(e.Position.X), int(e.Position.Y))
		for _, item := range r.Clickables {
			if clickPos.In(item.Area) {
				clickedOnItem = true
				break
			}
		}

		// 如果点击了非菜单项区域，关闭菜单
		if !clickedOnItem {
			r.Dismiss()
		}
	}
}

func (r *ContextArea) Layout(gtx C, w layout.Widget) D {
	r.Update(gtx)
	dims := layout.Dimensions{Size: gtx.Constraints.Min}

	// 激活状态下显示菜单
	if r.active {
		// 录制菜单内容
		menuContent := op.Record(gtx.Ops)
		r.dims = w(gtx)
		menuCall := menuContent.Stop()

		// 计算菜单位置
		pos := r.adjustPosition(dims.Size)

		// 设置菜单位置偏移
		defer op.Offset(pos).Push(gtx.Ops).Pop()

		// 显示菜单内容
		menuCall.Add(gtx.Ops)

		// 创建菜单外遮挡层（用于关闭菜单）
		defer clip.Rect(image.Rect(-1e6, -1e6, 1e6, 1e6)).Push(gtx.Ops).Pop()
		event.Op(gtx.Ops, &r.active)
	}

	// 注册上下文区域事件
	defer clip.Rect(image.Rectangle{Max: dims.Size}).Push(gtx.Ops).Pop()
	event.Op(gtx.Ops, r)

	return dims
}

// 智能调整菜单位置
func (r *ContextArea) adjustPosition(areaSize image.Point) image.Point {
	pos := image.Pt(int(r.position.X), int(r.position.Y))

	// 水平防溢出
	if pos.X+r.dims.Size.X > areaSize.X {
		switch r.PositionHint {
		case layout.E, layout.NE, layout.SE:
			pos.X = areaSize.X - r.dims.Size.X
		case layout.W, layout.NW, layout.SW:
			pos.X = 0
		default: // 自动
			if newX := pos.X - r.dims.Size.X; newX >= 0 {
				pos.X = newX
			} else {
				pos.X = areaSize.X - r.dims.Size.X
			}
		}
	}

	// 垂直防溢出
	if pos.Y+r.dims.Size.Y > areaSize.Y {
		switch r.PositionHint {
		case layout.S, layout.SE, layout.SW:
			pos.Y = areaSize.Y - r.dims.Size.Y
		case layout.N, layout.NE, layout.NW:
			pos.Y = 0
		default: // 自动
			if newY := pos.Y - r.dims.Size.Y; newY >= 0 {
				pos.Y = newY
			} else {
				pos.Y = areaSize.Y - r.dims.Size.Y
			}
		}
	}

	return pos
}

// 添加可点击菜单项
func (r *ContextArea) AddClickable(area image.Rectangle, onClick *widget.Clickable) {
	r.Clickables = append(r.Clickables, &clickableWidget{
		Area:  area,
		Click: onClick,
	})
}

// 在菜单布局中处理点击项
func (r *ContextArea) HandleClicks(gtx C) {
	for _, item := range r.Clickables {
		if item.Click.Clicked(gtx) {
			r.Dismiss()
		}

		// 注册点击事件
		area := clip.Rect(item.Area).Push(gtx.Ops)
		event.Op(gtx.Ops, item)
		area.Pop()
	}
}

func (r *ContextArea) Dismiss() {
	r.active = false
	r.justDismissed = true
}

func (r *ContextArea) Active() bool {
	return r.active
}

func (r *ContextArea) Activated() bool {
	was := r.justActivated
	r.justActivated = false
	return was
}

func (r *ContextArea) Dismissed() bool {
	was := r.justDismissed
	r.justDismissed = false
	return was
}
