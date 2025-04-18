package ux

import (
	"bytes"
	"gioui.org/app"
	_ "gioui.org/app/permission/networkstate"
	_ "gioui.org/app/permission/storage"
	"gioui.org/gpu/headless"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/safemap"
	"github.com/ddkwork/ux/widget/material"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
)

var th = material.NewTheme()

func NewTheme() *material.Theme {
	return th
}

type Widget interface {
	Layout(gtx layout.Context) layout.Dimensions
}

var ZeroWidget = func(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{}
}

type Panel struct { // 使用泛型而不是接口，这样返回的每个控件结构体字段无需断言，并且有类型约束，是安全的
	layout.Flex
	children          []layout.FlexChild
	w                 *app.Window
	dumpCanvas        bool
	imageAsBackground bool
}

func (p *Panel) Empty() bool {
	return len(p.children) == 0
}

func NewHPanel(w *app.Window) *Panel {
	panel := NewPanel(w)
	panel.Axis = layout.Horizontal
	return panel
}

func NewPanel(w *app.Window) *Panel {
	return &Panel{
		Flex: layout.Flex{
			Axis:      layout.Vertical,
			Spacing:   0,
			Alignment: 0,
			WeightSum: 0,
		},
		children: make([]layout.FlexChild, 0),
		w:        w,
	}
}

func (p *Panel) AddChildCallback(childCallback func(gtx layout.Context) layout.Dimensions) {
	p.children = append(p.children, layout.Rigid(childCallback))
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
	BackgroundDark(gtx)
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

type AppBar struct {
	Search   *Input
	ToolBars []*TipIconButton
	About    *TipIconButton
}

func InitAppBar(panel *Panel, toolBars []*TipIconButton, speechTxt string) *AppBar {
	search := NewInput("请输入搜索关键字...").SetIcon(ActionSearchIcon).SetRadius(16)
	panel.AddChildFlexed(1, search) // todo 太多之后apk需要管理溢出

	for _, toolbar := range toolBars {
		panel.AddChild(toolbar)
	}

	about := NewTooltipButton(AlertErrorIcon, "about", func() { // todo ico make
		if mylog.IsAndroid() {
			mylog.Info("android not support about window")
			return
		}
		About(NewWindow("about"), speechTxt)
	})
	panel.AddChild(about)
	return &AppBar{
		Search:   search,
		ToolBars: toolBars,
		About:    about,
	}
}

func BackgroundDark(gtx layout.Context) {
	rect := clip.Rect{Max: gtx.Constraints.Max}
	paint.FillShape(gtx.Ops, BackgroundColor, rect.Op())
}

func drawImageBackground(gtx layout.Context) {
	data := mylog.Check2(os.ReadFile("asset/background.png"))
	img := mylog.Check2(png.Decode(bytes.NewReader(data)))
	dst := image.NewRGBA(image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y))
	draw.Draw(dst, dst.Bounds(), img, image.Point{}, draw.Over)
	paint.NewImageOp(dst).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

// //////////////////////////////////

var ModalCallbacks = new(safemap.M[string, func()])

func NewWindow(title string) *app.Window {
	w := new(app.Window)
	w.Option(
		app.Title(title),
		app.Size(1200, 600),
		// app.Decorated(false),
	)
	w.Perform(system.ActionCenter)
	// mylog.Check(android_background_service.Start()) // todo fix xml
	return w
}

func Run(p *Panel) {
	mylog.Call(func() {
		/*
				var et event.Event
			if m.plugin != nil {
				et = m.plugin(m.win)
			} else {
				et = m.win.Event()
			}
			m.Explorer.ListenEvents(et)
		*/

		//var (
		//	deco  widget.Decorations
		//	title string
		//)

		var ops op.Ops
		go func() {
			for {
				e := p.w.Event()
				switch e := e.(type) {
				case app.DestroyEvent:
					mylog.Check(e.Err)
					os.Exit(0)
				// case app.ConfigEvent:
				//	deco.Maximized = e.Config.Mode == app.Maximized
				//	title = e.Config.Title
				case app.FrameEvent:
					gtx := app.NewContext(&ops, e)
					for _, v := range ModalCallbacks.Range() {
						v()
					}
					p.Layout(gtx)

					//p.w.Perform(deco.Update(gtx))
					//decorationsStyle := material.Decorations(th, &deco, ^system.Action(0), title)
					//decorationsStyle.Background = color.NRGBA{
					//	R: 44,
					//	G: 44,
					//	B: 44,
					//	A: 255,
					//}
					//decorationsStyle.Foreground = th.Fg
					//decorationsStyle.Layout(gtx)

					e.Frame(gtx.Ops)
				}
			}
		}()
		app.Main()
	})
}

func SaveScreenshot(callback Widget) {
	const scale = 1.5
	size := image.Point{X: 1200 * scale, Y: 600 * scale}
	w := mylog.Check2(headless.NewWindow(size.X, size.Y))
	gtx := layout.Context{
		Ops: new(op.Ops),
		Metric: unit.Metric{
			PxPerDp: scale,
			PxPerSp: scale,
		},
		Constraints: layout.Exact(size),
	}
	BackgroundDark(gtx)
	callback.Layout(gtx)
	mylog.Check(w.Frame(gtx.Ops))
	img := image.NewRGBA(image.Rectangle{Max: size})
	mylog.Check(w.Screenshot(img))
	var buf bytes.Buffer
	mylog.Check(png.Encode(&buf, img))
	mylog.Check(os.WriteFile(filepath.Join(DataDir(), "canvas.png"), buf.Bytes(), 0o666))
}

func DataDir() string {
	switch {
	case mylog.IsAndroid():
		return mylog.Check2(app.DataDir())
	case mylog.IsTermux():
		return "/data/data/com.termux/files/usr" // todo choose another dir
	default: // windows,linux
		return "."
	}
}

type (
	C = layout.Context
	D = layout.Dimensions
)

type Rect struct {
	Color color.NRGBA
	Size  image.Point
	Radii int
}

func (r Rect) Layout(gtx C) D {
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
	//richText := NewRichText()
	//richText.AddSpan(richtext.SpanStyle{
	//	// Font:        font.Font{},
	//	Size:        unit.Sp(12),
	//	Color:       White,
	//	Content:     text,
	//	Interactive: false,
	//})
	//recording := Record(gtx, func(gtx layout.Context) layout.Dimensions {
	//	gtx.Constraints.Min.X = 0
	//	return richText.Layout(gtx)
	//})
	//// fmt.Printf("Calculated dragWidth: %v\n", unit.Dp(recording.Dimensions.Size.X))
	//return unit.Dp(recording.Dimensions.Size.X)
	text += "  ⇧" // 为排序图标留位置,不要修改这里，稳定了
	body := material.Body1(th, text)
	body.MaxLines = 1
	recording := Record(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = 0
		return body.Layout(gtx)
	})
	return unit.Dp(recording.Dimensions.Size.X)
}

func MaxLabelWidth(gtx layout.Context, keys []string) unit.Dp {
	// originalConstraints := gtx.Constraints
	maxWidth := unit.Dp(0)
	for _, data := range keys {
		currentWidth := LabelWidth(gtx, data) // 可以使用max []unit.Dp，但是多了一层make []unit.Dp，浪费内存
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

func Record(gtx layout.Context, w layout.Widget) Recording { // 应用场景:计算单元格宽度求平均宽度
	m := op.Record(gtx.Ops)
	dims := w(gtx)
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
