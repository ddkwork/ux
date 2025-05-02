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
	"github.com/ddkwork/golibrary/mylog"

	"github.com/ddkwork/golibrary/safemap"
	"github.com/ddkwork/golibrary/stream"
)

func TestTreeTable_ContextMenuItem(t1 *testing.T) {
	m := safemap.NewOrdered[string, string](func(yield func(string, string) bool) {
		yield("CopyRow", "CopyRow")
		yield("ConvertToContainer", "ConvertToContainer")
		yield("ConvertToNonContainer", "ConvertToNonContainer")
		yield("New", "New")
		yield("NewContainer", "NewContainer")
		yield("Delete", "Delete")
		yield("Duplicate", "Duplicate")
		yield("Edit", "Edit")
		yield("OpenAll", "OpenAll")
		yield("CloseAll", "CloseAll")
		yield("SaveData", "SaveData")
	})
	stream.NewGeneratedFile().EnumTypes("ContextMenuItem", m)
}

func TestNewPopupMenu(t *testing.T) {
	t.Skip("finished")
	w := new(app.Window)
	m := NewContextMenu()
	for i := range 100 {
		m.CreatItem(func(gtx layout.Context) layout.Dimensions {
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
			BackgroundDark(gtx)
			m.LayoutTest(gtx)
			e.Frame(gtx.Ops)
		}
	}
	app.Main()
}
