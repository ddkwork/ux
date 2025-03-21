package main

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/x/outlay"
	"github.com/ddkwork/ux"
	"github.com/ddkwork/ux/widget/material"
	"github.com/ddkwork/ux/x/component"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
)

// https://github.com/hkontrol/hkapp
func main() {
	go func() {
		w := new(app.Window)
		app.Title("hkontroller")
		app.Size(unit.Dp(400), unit.Dp(600))
		// w := new(apWindow)
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

	var tagSearchBtn widget.Clickable
	var tagRemoveBtn widget.Clickable
	tagCtxMenu := component.MenuState{
		Options: []func(gtx C) D{
			func(gtx C) D {
				item := component.MenuItem(th, &tagSearchBtn, "Search")
				item.Icon = ux.ActionVisibilityIcon
				item.Hint = component.MenuHintText(th, "")
				return item.Layout(gtx)
			},
			func(gtx C) D {
				item := component.MenuItem(th, &tagRemoveBtn, "Remove")
				item.Icon = ux.ContentCreateIcon
				item.Hint = component.MenuHintText(th, "")
				return item.Layout(gtx)
			},
		},
	}

	selectedAccTags := []string{
		"xxxxxx",
		"yyyyyy",
		"zzzzzz",
		"aaaaaa",
		"bbbbb",
		"ccccc",
		"ddddd",
		"eeeee",
		"ffffff",
		"gggggg",
		"hhhhhh",
		"iiiiii",
		"jjjjjj",
		"kkkkkk",
		"llllll",
		"mmmmmm",
		"nnnnnn",
		"oooooo",
		"pppppp",
		"qqqqqq",
		"rrrrrr",
		"ssssss",
		"tttttt",
		"vvvvvv",
		"wwwwww",
	}
	tagClickables := make([]widget.Clickable, len(selectedAccTags))
	tagCtxAreas := make([]component.ContextArea, len(selectedAccTags))
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			outlay.FlowWrap{
				Axis:      0,
				Alignment: 0,
			}.Layout(gtx, len(selectedAccTags), func(gtx layout.Context, i int) layout.Dimensions {
				state := &tagCtxAreas[i]

				return layout.Stack{}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx C) D {
							return tagClickables[i].Layout(gtx, func(gtx C) D {
								if tagClickables[i].Clicked(gtx) {
									println("clicked", selectedAccTags[i])
								}
								return material.Button(th, &tagClickables[i], selectedAccTags[i]).Layout(gtx)
							})
						})
					}),
					layout.Expanded(func(gtx C) D {
						return state.Layout(gtx, func(gtx C) D {
							gtx.Constraints.Min.X = 0
							return component.Menu(th, &tagCtxMenu).Layout(gtx)
						})
					}),
				)
			})
			e.Frame(gtx.Ops)
		}
	}
}

type (
	C = layout.Context
	D = layout.Dimensions
)
