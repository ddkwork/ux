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
	ux.Run("demo", ux.NewLogView())
}
