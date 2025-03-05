// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"log"
	"os"

	"gioui.org/widget"
	"github.com/ddkwork/ux"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
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

func loop(w *app.Window) error {
	var modal *ux.MessageModal
	modal = ux.NewMessageModal("Error", "err.Error()", ux.MessageModalTypeErr, func(selectedOption string) {
		modal.Hide()
	}, ux.ModalOption{
		Text:   "close",
		Button: widget.Clickable{},
		Icon:   ux.NavigationCloseIcon,
	},
		ux.ModalOption{
			Text:   "submit",
			Button: widget.Clickable{},
			Icon:   ux.NavigationSubdirectoryArrowRightIcon,
		},
	)
	modal.Show()

	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			ux.BackgroundDark(gtx)
			modal.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}
