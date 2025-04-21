package main

import (
	_ "embed"
	"gioui.org/app"
	"gioui.org/layout"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"github.com/ddkwork/golibrary/stream/desktop"
	"github.com/ddkwork/ux"
	"github.com/ddkwork/ux/resources/images"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (

	//go:embed appicon.png
	VStartPng []byte

	//go:embed  installer.png
	installerJpg []byte
)

func main() {
	w := new(app.Window)
	w.Option(app.Title("VStart"))

	panel := ux.NewPanel(w)
	flow := ux.NewFlow(8)
	panel.AddChild(flow)

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
		case ".exe": //msi invalid argument , not support
			if stream.IsWindows() {
				path = filepath.ToSlash(path)

				oldPng := path[:len(path)-len(filepath.Ext(path))] + ".png"
				if stream.IsFilePath(oldPng) {
					mylog.Check(os.Remove(oldPng))
				}

				png := ExtractIcon2Png(path)
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

	ux.Run(panel)
}
