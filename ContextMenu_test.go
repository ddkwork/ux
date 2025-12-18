package ux

import (
	"fmt"
	"os"
	"testing"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/ddkwork/ux/widget/material"

	"gioui.org/app"
	"gioui.org/op"
	"github.com/ddkwork/golibrary/std/mylog"
)

func TestNewPopupMenu(t *testing.T) {
	t.Skip("finished")
	w := new(app.Window)
	m := NewContextMenu()
	rows := make([]layout.Widget, 0)
	for i := range 100 {
		rows = append(rows, func(gtx layout.Context) layout.Dimensions {
			rowClick := new(widget.Clickable)
			buttonStyle := material.Button(th, rowClick, "item"+fmt.Sprintf("%d", i))
			buttonStyle.Color = RowColor(i)
			return buttonStyle.Layout(gtx)
		})
	}

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			mylog.CheckIgnore(e.Err)
			os.Exit(0)
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			DarkBackground(gtx)
			m.LayoutTest(gtx, rows)
			e.Frame(gtx.Ops)
		}
	}
	app.Main()
}
