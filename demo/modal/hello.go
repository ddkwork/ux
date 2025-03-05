// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/ddkwork/ux"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

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

var (
	ExplorerViewID = ux.NewViewID("FileExplorerView")
	v              = ux.ModalView{
		View: &data{
			finished: false,
			edit:     ux.NewInput("Hello, Gio"),
			closeBtn: new(widget.Clickable),
		},
		Padding: layout.Inset{},
		// Background: th.Bg,
		MaxWidth:  unit.Dp(760),
		MaxHeight: 0.7,
		Radius:    unit.Dp(8),
		Halted:    false,
	}
)

func loop(w *app.Window) error {
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			if v.Finished() {
				v.Anim().ToggleVisibility(gtx.Now)
				gtx.Execute(op.InvalidateCmd{})
				// w.Perform(system.ActionClose)
				w.Invalidate()
				continue
			}
			v.ShowUp(gtx)
			v.Layout(gtx)

			e.Frame(gtx.Ops)
		}
	}
}

type data struct {
	finished bool
	edit     *ux.Input
	closeBtn *widget.Clickable
}

func (d *data) Actions() []ux.ViewAction {
	icon1, _ := widget.NewIcon(icons.Action3DRotation)
	return []ux.ViewAction{
		{
			Name: "xxx",
			Icon: icon1,
			OnClicked: func(gtx ux.C) {
				fmt.Println("xxx")
			},
		},
	}
}

func (d *data) Layout(gtx layout.Context) layout.Dimensions {
	for d.closeBtn.Clicked(gtx) {
		d.finished = true
	}
	if d.finished {
		return layout.Dimensions{}
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ux.Button(d.closeBtn, ux.NavigationCloseIcon, "Close").Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return d.edit.Layout(gtx)
		}),
	)
}

func (d *data) ID() ux.ViewID {
	return ExplorerViewID
}

func (d *data) Location() url.URL {
	return url.URL{}
}

func (d *data) Title() string {
	return "Explorer"
}

func (d *data) OnFinish() {
	d.finished = true
}

func (d *data) Finished() bool {
	return d.finished
}
