package main

import (
	"log"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/ux/hashicon"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(
			app.Size(740, 500))
		mylog.Check(run(w))

		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	columnCount := 6
	rowCount := 4

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			grid(gtx, rowCount, columnCount)
			e.Frame(gtx.Ops)
		}
	}
}

func grid(gtx layout.Context, rowCount int, columnCount int) layout.Dimensions {
	var rows []layout.FlexChild
	var i int

	for y := 0; y < rowCount; y++ {

		var cols []layout.FlexChild
		for x := 0; x < columnCount; x++ {
			cols = append(cols,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					hash := strconv.FormatInt(int64(i), 10)
					dims := hashicon.Hashicon{Config: hashicon.DefaultConfig}.Layout(gtx, 100, hash)
					i++
					return dims
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(20)}.Layout),
			)
		}

		rows = append(rows,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{}.Layout(gtx, cols...)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
		)
	}

	return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, rows...)
	})
}
