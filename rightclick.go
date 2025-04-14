package ux

import (
	"fmt"
	"github.com/ddkwork/ux/widget/material"
	"image"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"
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

func (o *Overlay) Layout(gtx layout.Context) layout.Dimensions { //stack遮罩layout的自定义坐标布局实现
	for _, item := range o.items {
		macro := op.Record(gtx.Ops)
		dims := item.Widget(gtx)
		call := macro.Stop()
		offset := item.OffsetWithin(layout.FPt(dims.Size), layout.FPt(gtx.Constraints.Max))
		//func(item overlayItem) {
		stack := op.Offset(image.Point{X: int(offset.X), Y: int(offset.Y)}).Push(gtx.Ops)
		call.Add(gtx.Ops)
		stack.Pop()
		//}(item)
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

type PopMenu struct {
	Overlay
	itemProvider *ContextMenu
	content      layout.Widget //row or button or label
	*RightClickArea
}

func (p *PopMenu) Layout(gtx layout.Context) layout.Dimensions {
	//event.Op(gtx.Ops, p)

	//行单双击实践
	//if p.rowButton.Clicked(gtx) {
	//	println("row selected")
	//}
	p.RightClickArea.Layout(gtx) //在哪弹出菜单，通过调用content布局，也就是点击行触发rowButton的点击事件，获点击的行，然后rowButton的layout内布局表格ro的渲染，所以理论上可以无需渲染rowButton出来任何东西，调试模式可以渲染一个带标题的按钮
	// todo 弹出菜单限制在他的区域内？似乎不合理,并且可能会导致溢出,
	//  但是不这样限制的话应该关联点击了哪一行，如果不在当前行弹出菜单也是有点奇怪
	//  原版似乎是限制区域的,这似乎合理，得实现这个事件行为

	//右键菜单的item事件
	for _, item := range p.itemProvider.Items {
		if item.Clickable.Clicked(gtx) {
			println(item.Title) //todo call item callback
			p.CloseMenu()
		}
	}
	return p.Overlay.Layout(gtx) //弹出菜单
}

func NewPopMenu(itemProvider *ContextMenu, drawRow layout.Widget) *PopMenu {
	p := &PopMenu{
		Overlay:      Overlay{},
		itemProvider: itemProvider,
		content: func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = gtx.Constraints.Max
			return drawRow(gtx) //todo bug
		},
		RightClickArea: nil,
	}
	r := &RightClickArea{
		Content: func(gtx layout.Context) layout.Dimensions { return p.content(gtx) },
		Menu: func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = image.Point{}
			gtx.Constraints.Max.X = gtx.Dp(unit.Dp(200)) //todo 斑马线，分隔条，圆角，长按支持apk
			var children []layout.FlexChild
			for i, item := range p.itemProvider.Items {
				children = append(children, layout.Rigid(func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X //弹出菜单的item统一宽度，太窄不美观
					//gtx.Constraints.Min.Y = 3000                  //todo bug                  //gtx.Constraints.Max.Y
					//return material.Button(th.Theme, &item.Clickable, item.Title).Layout(gtx)
					return Background{Color: RowColor(i + 1)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions { //not work
						return material.Button(th.Theme, &item.Clickable, item.Title).Layout(gtx)
					})
				}))
			}
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
		},
		Anchor:      nil,
		Overlay:     &p.Overlay,
		leftPressed: nil,
	}
	p.RightClickArea = r
	return p
}

func loop(w *app.Window) error {
	popMenu := NewPopMenu(NewContextMenu(), nil)
	popMenu.itemProvider.AddItem(ContextMenuItem{
		Title:         "item1",
		Icon:          nil,
		Can:           nil,
		Do:            nil,
		AppendDivider: false,
		Clickable:     widget.Clickable{},
	})

	popMenu.itemProvider.AddItem(ContextMenuItem{
		Title:         "item2",
		Icon:          nil,
		Can:           nil,
		Do:            nil,
		AppendDivider: false,
		Clickable:     widget.Clickable{},
	})

	popMenu.itemProvider.AddItem(ContextMenuItem{
		Title:         "item3",
		Icon:          nil,
		Can:           nil,
		Do:            nil,
		AppendDivider: false,
		Clickable:     widget.Clickable{},
	})

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			BackgroundDark(gtx)
			layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions { return popMenu.Layout(gtx) })
			e.Frame(gtx.Ops)
		}
	}
}
