// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
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
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops

	p := Page{
		redButton:          widget.Clickable{},
		greenButton:        widget.Clickable{},
		blueButton:         widget.Clickable{},
		balanceButton:      widget.Clickable{},
		accountButton:      widget.Clickable{},
		cartButton:         widget.Clickable{},
		leftFillColor:      color.NRGBA{},
		leftContextArea:    ContextArea{},
		leftMenu:           component.MenuState{},
		rightMenu:          component.MenuState{},
		menuInit:           false,
		menuDemoList:       widget.List{},
		menuDemoListStates: nil,
		List:               widget.List{},
	}

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			p.Layout(gtx, th)
			e.Frame(gtx.Ops)
			gtx.Ops.Reset()
			gtx.Reset()
		}
	}
}
