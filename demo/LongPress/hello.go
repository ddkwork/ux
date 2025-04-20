package main

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/ux"
	"github.com/ddkwork/ux/resources/icons"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
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

	keys := []string{
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
	flow := ux.NewFlow(5)
	for i, key := range keys {
		flow.AppendElem(i, ux.FlowElemButton{
			Title: key,
			Icon:  icons.IconMap.Values()[i],
			Do:    func(gtx layout.Context) { mylog.Info(key + " pressed") }, //run exe
			ContextMenuItems: []ux.ContextMenuItem{
				{
					Title:         "Balance",
					Icon:          icons.ActionAccountBalanceIcon,
					Can:           func() bool { return true },
					Do:            func() { mylog.Info("Balance item clicked") },
					AppendDivider: false,
					Clickable:     widget.Clickable{},
				},
				{
					Title:         "Account",
					Icon:          icons.ActionAccountBoxIcon,
					Can:           func() bool { return true },
					Do:            func() { mylog.Info("Account item clicked") },
					AppendDivider: false,
					Clickable:     widget.Clickable{},
				},
				{
					Title:         "Cart",
					Icon:          icons.ActionAddShoppingCartIcon,
					Can:           func() bool { return true },
					Do:            func() { mylog.Info("Cart item clicked") },
					AppendDivider: false,
					Clickable:     widget.Clickable{},
				},
			},
		})
	}

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			ux.BackgroundDark(gtx)
			flow.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}
