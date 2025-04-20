package main

import (
	"github.com/ddkwork/ux"
	"github.com/ddkwork/ux/resources/icons"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
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
	menus := make([]*ux.ContextMenu, len(selectedAccTags))

	flow := ux.Flow{
		Num:       5,
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			ux.BackgroundDark(gtx)

			flow.Layout(gtx, len(selectedAccTags), func(gtx layout.Context, i int) layout.Dimensions {
				return layout.UniformInset(2).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = 200 //todo into flow
					gtx.Constraints.Max.X = 200
					return tagClickables[i].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						if tagClickables[i].Clicked(gtx) {
							println("clicked", selectedAccTags[i])
						}
						menus[i] = ux.NewContextMenuWithRootRows(func(gtx layout.Context) layout.Dimensions {
							return ux.Button(&tagClickables[i], icons.IconMap.Values()[i], selectedAccTags[i]).Layout(gtx)
						})
						return menus[i].LayoutTest(gtx)
					})
				})
			})
			e.Frame(gtx.Ops)
		}
	}
}
