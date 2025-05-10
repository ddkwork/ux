package ux

import (
	"bytes"
	"image"
	"image/png"
	"iter"
	"os"
	"path/filepath"

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
	"gioui.org/x/explorer"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/safemap"
	"github.com/ddkwork/ux/resources/colors"
	"github.com/ddkwork/ux/resources/images"
	"github.com/ddkwork/ux/widget/material"
	"golang.org/x/image/draw"
)

func NewTheme() *material.Theme {
	return th
}

var (
	th      = material.NewTheme()
	explore *explorer.Explorer
)

func Run(title string, widget Widget) {
	mylog.Call(func() {
		w := new(app.Window)
		w.Option(
			app.Title(title),
			app.Size(1200, 600),
			// app.Decorated(false),
		)
		w.Perform(system.ActionCenter)
		explore = explorer.NewExplorer(w)
		// mylog.Check(android_background_service.Start()) // todo fix xml

		var ops op.Ops
		go func() {
			for {
				e := w.Event()
				explore.ListenEvents(e)
				switch e := e.(type) {
				case app.DestroyEvent:
					mylog.Check(e.Err)
					os.Exit(0)

				case app.FrameEvent:
					gtx := app.NewContext(&ops, e)
					DarkBackground(gtx)
					for _, v := range ModalCallbacks.Range() {
						v()
					}
					mylog.Call(func() { widget.Layout(gtx) })

					e.Frame(gtx.Ops)
				}
			}
		}()
		app.Main()
	})
}

var ModalCallbacks = new(safemap.M[string, func()]) // todo use stack

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
	DarkBackground(gtx)
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

type AppBar struct {
	Search *Input
	About  *TipIconButton
}

func InitAppBar(hPanel *Panel, toolBars iter.Seq[*TipIconButton], speechTxt string) *AppBar {
	search := NewInput("请输入搜索关键字...").SetIcon(images.ActionSearchIcon).SetRadius(16)
	hPanel.AddChildFlexed(1, search) // todo 太多之后apk需要管理溢出

	for toolbar := range toolBars {
		hPanel.AddChild(toolbar)
	}

	about := NewTooltipButton(images.AlertErrorIcon, "about", func() {
		if mylog.IsAndroid() {
			mylog.Info("android not support about window")
			return
		}
		About(speechTxt)
	})
	hPanel.AddChild(about)
	return &AppBar{
		Search: search,
		About:  about,
	}
}

func DarkBackground(gtx layout.Context) {
	rect := clip.Rect{Max: gtx.Constraints.Max}
	paint.FillShape(gtx.Ops, colors.BackgroundColor, rect.Op())
}

func drawImageBackground(gtx layout.Context) {
	data := mylog.Check2(os.ReadFile("asset/background.png"))
	img := mylog.Check2(png.Decode(bytes.NewReader(data)))
	dst := image.NewRGBA(image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y))
	draw.Draw(dst, dst.Bounds(), img, image.Point{}, draw.Over)
	paint.NewImageOp(dst).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
