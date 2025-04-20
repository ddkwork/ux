package main

import (
	_ "embed"
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"github.com/ddkwork/golibrary/stream/desktop"
	"github.com/ddkwork/ux"
	"github.com/ddkwork/ux/resources/images"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:embed VStart.png
var VStartPng []byte

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
	flow := ux.NewFlow(5)
	i := 0
	mylog.Check(filepath.WalkDir("d:\\app", func(path string, info fs.DirEntry, err error) error {
		switch {
		case strings.Contains(path, "RECYCLE.BIN"):
			return err
		case info.IsDir():
			return err
		}
		ext := filepath.Ext(path)
		switch ext {
		case ".exe":
			if stream.IsWindows() {
				path = filepath.ToSlash(path)
				png := ExtractIcon2Png(path)
				if png == nil {
					png = VStartPng
				}
				flow.AppendElem(i, ux.FlowElemButton{
					Title: stream.AlignString(stream.BaseName(path), 5),
					Icon:  png,
					Do: func(gtx layout.Context) {
						stream.RunCommandArgs("start", path)
					},
					ContextMenuItems: []ux.ContextMenuItem{
						{
							Title: "open dir",
							Icon:  images.FileFolderOpenIcon,
							Can:   func() bool { return true },
							Do: func() {
								go desktop.Open(filepath.Dir(path))
								// go stream.RunCommandArgs("cd ", filepath.Dir(button.Tooltip.String()), "start", button.Tooltip.String())
							},
						},
					},
				})
			}
			i++
		}
		return err
	}))

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			ux.BackgroundDark(gtx)
			flow.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}
