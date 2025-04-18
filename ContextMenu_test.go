package ux

import (
	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/mylog"
	"os"
	"testing"

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
	w := new(app.Window)
	p := NewContextMenu(100, nil)
	p.AddItem(ContextMenuItem{
		Title:         "item1",
		Icon:          nil,
		Can:           func() bool { return false },
		Do:            func() { mylog.Info("item1 clicked") },
		AppendDivider: false,
		Clickable:     widget.Clickable{},
	})
	p.AddItem(ContextMenuItem{
		Title: "item2",
		Icon:  nil,
		Can:   func() bool { return true },
		Do:    func() { mylog.Info("item2 clicked") },
	})

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			mylog.CheckIgnore(e.Err)
			os.Exit(0)
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			BackgroundDark(gtx)
			p.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
	app.Main()
}
