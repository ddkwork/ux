package main

import (
	"gioui.org/app"
	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/ux"
)

func main() {
	app.FileDropCallback(func(files []string) {
		mylog.Struct(files)
	})

	p := ux.NewPanel()
	p.AddChild(
		ux.NewInput("", "Input", "Hello, world!"),
		//ux.NewLogView(),
	)
	ux.Run("demo", p)
}
