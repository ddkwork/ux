package ux

import (
	"image"
	"image/color"
	"sync"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/ddkwork/ux/widget/material"
)

type Panel struct { // 使用泛型而不是接口，这样返回的每个控件结构体字段无需断言，并且有类型约束，是安全的
	layout.Flex
	children          []layout.FlexChild
	dumpCanvas        bool
	imageAsBackground bool
}

func (p *Panel) Empty() bool {
	return len(p.children) == 0
}

func NewHPanel() *Panel {
	panel := NewPanel()
	panel.Axis = layout.Horizontal
	return panel
}

func NewPanel() *Panel { // 90% is Vertical
	return &Panel{
		Flex: layout.Flex{
			Axis:      layout.Vertical,
			Spacing:   0,
			Alignment: 0,
			WeightSum: 0,
		},
		children: make([]layout.FlexChild, 0),
	}
}

func (p *Panel) AddChildCallback(childCallback Widget) {
	p.children = append(p.children, layout.Rigid(childCallback.Layout))
}

func (p *Panel) AddChild(child ...Widget) {
	for _, c := range child {
		p.children = append(p.children, layout.Rigid(c.Layout))
	}
}

func (p *Panel) AddChildFlexed(weight float32, child Widget) {
	p.children = append(p.children, layout.Flexed(weight, child.Layout))
}

func (p *Panel) Layout(gtx layout.Context) layout.Dimensions {
	if p.dumpCanvas {
		SaveScreenshot(p)
	}
	if p.imageAsBackground {
		drawImageBackground(gtx)
	}
	if p.Empty() {
		return p.Flex.Layout(gtx)
	}
	return p.Flex.Layout(gtx, p.children...)
}

func (p *Panel) SetDumpCanvas(dumpCanvas bool) {
	p.dumpCanvas = dumpCanvas
}

func (p *Panel) SetImageAsBackground(imageAsBackground bool) {
	p.imageAsBackground = imageAsBackground
}

type Rect struct {
	Color color.NRGBA
	Size  image.Point
	Radii int
}

func (r Rect) Layout(gtx layout.Context) layout.Dimensions {
	paint.FillShape(
		gtx.Ops,
		r.Color,
		clip.UniformRRect(
			image.Rectangle{
				Max: r.Size,
			},
			r.Radii,
		).Op(gtx.Ops))
	return layout.Dimensions{Size: r.Size}
}

func WithAlpha(c color.NRGBA, a uint8) color.NRGBA {
	return color.NRGBA{
		R: c.R,
		G: c.G,
		B: c.B,
		A: a,
	}
}

func LabelWidth(gtx layout.Context, text string) unit.Dp {
	// fmt.Printf("Calculating text dragWidth for: %s\n", text)
	// fmt.Printf("Current Min.X: %v\n", gtx.Constraints.Min.X)
	// richText := NewRichText()
	// richText.AddSpan(richtext.SpanStyle{
	//	// Font:        font.Font{},
	//	Size:        unit.Sp(12),
	//	Color:       White,
	//	Content:     text,
	//	Interactive: false,
	// })
	// recording := Record(gtx, func(gtx layout.Context) layout.Dimensions {
	//	gtx.Constraints.Min.X = 0
	//	return richText.Layout(gtx)
	// })
	// // fmt.Printf("Calculated dragWidth: %v\n", unit.Dp(recording.Dimensions.Size.X))
	// return unit.Dp(recording.Dimensions.Size.X)
	body := material.Body1(th, text)
	body.MaxLines = 1
	recording := Record(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = 0
		return body.Layout(gtx)
	})
	return unit.Dp(recording.Dimensions.Size.X)
}

func MaxLabelWidth(gtx layout.Context, rows []CellData) unit.Dp {
	// originalConstraints := gtx.Constraints
	maxWidth := unit.Dp(0)
	for _, data := range rows {
		currentWidth := LabelWidth(gtx, data.Name) // 可以使用max []unit.Dp，但是多了一层make []unit.Dp，浪费内存
		if currentWidth > maxWidth {
			maxWidth = currentWidth
		}
	}
	// gtx.Constraints = originalConstraints
	return maxWidth
}

type Recording struct {
	Call       op.CallOp
	Dimensions layout.Dimensions
}

func (r Recording) Layout(gtx layout.Context) layout.Dimensions {
	defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
	r.Call.Add(gtx.Ops)
	return r.Dimensions
}

var lock = &sync.RWMutex{}

func Record(gtx layout.Context, w layout.Widget) Recording { // 应用场景:计算单元格宽度求平均宽度
	m := op.Record(gtx.Ops)
	lock.Lock()
	dims := w(gtx)
	lock.Unlock()
	c := m.Stop()
	return Recording{c, dims}
}

type Background struct {
	Color color.NRGBA
}

func (b Background) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := w(gtx)
	call := macro.Stop()
	paint.FillShape(gtx.Ops, b.Color, clip.Rect{Max: dims.Size}.Op())
	call.Add(gtx.Ops)
	return dims
}
