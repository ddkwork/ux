package main

import (
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"os"

	terminal2 "github.com/ddkwork/ux/internal/terminal"
	"github.com/ddkwork/ux/internal/terminal/tint"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/mylog"
)

type TerminalWindow struct {
	Screen         *os.File
	screen         *terminal2.Screen
	quitChannel    chan any
	updatedChannel chan any
}

func (l TerminalWindow) Close() {
	l.quitChannel <- struct{}{}
}

func (l TerminalWindow) Open() error {
	w := &app.Window{}

	guiReady := make(chan any)
	var ops op.Ops

	button := new(widget.Clickable)
	settings := terminal2.NewConsoleSettings(terminal2.MaxSize(100, 30))

	go func() {
		w.Option(app.Size(unit.Dp(670), unit.Dp(524)))
		guiReady <- struct{}{}

		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				return

			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				if button.Clicked(gtx) {
					w.Perform(system.ActionClose)
				}
				terminal2.Console(l.screen, settings)(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}()

	// Wait for the GUI to be ready
	<-guiReady

	for {
		select {
		case <-l.quitChannel:
			w.Perform(system.ActionClose)

		case <-l.updatedChannel:
			w.Invalidate()
		}
	}
}

func NewTerminalWindow(size terminal2.Point) *TerminalWindow {
	updatedChannel := make(chan any)
	screen := terminal2.NewScreen(size, updatedChannel)

	r, w := mylog.Check3(os.Pipe())

	go func() {
		_ = mylog.Check2(io.Copy(screen, r))
	}()

	return &TerminalWindow{
		Screen:         w,
		screen:         screen,
		updatedChannel: updatedChannel,
		quitChannel:    make(chan any),
	}
}

func main() {
	w := NewTerminalWindow(terminal2.Point{
		X: 80,
		Y: 20,
	})
	os.Stdout = w.Screen
	os.Stderr = w.Screen

	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: "15:04:05",
			NoColor:    os.Getenv("NO_COLOR") == "1",
		}),
	))

	go func() {
		mylog.Check(w.Open())
	}()

	RESET := "\u001B[0m"
	BOLD := "\u001B[1m"
	FAINT := "\u001B[2m"

	go func() {
		//for i := 0; i < 200; i++ {
		//	fmt.Println(randomString(82))
		//}

		fmt.Println("ANSI Test")
		fmt.Println("=========")
		slog.Debug("This is not very important")
		slog.Info("Information message", "key", "value")
		slog.Warn("It's getting real")
		slog.Error("Oh no!")

		fmt.Println(BOLD + "This is bold" + RESET)
		fmt.Println(FAINT + "This is bold" + RESET)
		fmt.Println("\u001b[38;2;253;182;0mRgb code" + RESET)
		fmt.Println("\u001b[38;5;63m256 color code" + RESET)

		fmt.Println("")
		fmt.Println(randomString(200))
	}()

	print("Starting main")

	app.Main()
}

func randomString(n int) string {
	s := make([]rune, n)
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := range n {
		s[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(s)
}
