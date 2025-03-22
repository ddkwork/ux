package main

import (
	"demo/rightclick/anchor"
	"gioui.org/layout"
	"gioui.org/op"
	"image"
)

type Overlay struct {
	items []overlayItem
}

type overlayItem struct {
	anchor.Anchor
	layout.Widget
}

func (o *Overlay) LayoutAt(anchor anchor.Anchor, widget layout.Widget) {
	o.items = append(o.items, overlayItem{Anchor: anchor, Widget: widget})
}

func (o *Overlay) Layout(gtx layout.Context) layout.Dimensions {
	for _, item := range o.items {
		macro := op.Record(gtx.Ops)
		dims := item.Widget(gtx)
		call := macro.Stop()

		offset := item.Anchor.OffsetWithin(layout.FPt(dims.Size), layout.FPt(gtx.Constraints.Max))
		func(item overlayItem) {
			defer op.TransformOp{}.Push(gtx.Ops).Pop()
			//defer op.Push(gtx.Ops).Pop()
			op.Offset(image.Point{
				X: int(offset.X),
				Y: int(offset.Y),
			}).Add(gtx.Ops)
			call.Add(gtx.Ops)
		}(item)
	}
	o.items = o.items[:0]
	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}
