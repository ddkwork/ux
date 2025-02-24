package terminal

import (
	"fmt"
	"github.com/ddkwork/ux/terminal/tint"
	"io"
	"log/slog"
	"math/rand"
	"os"

	"github.com/ddkwork/golibrary/mylog"
)

func Demo() (*Screen, *ConsoleSettings) {
	settings := NewConsoleSettings(MaxSize(0, 0))
	screen := NewScreen(Point{X: 20, Y: 40}, nil)

	r, w := mylog.Check3(os.Pipe())
	go func() {
		mylog.Check2(io.Copy(screen, r))
	}()
	os.Stdout = w
	os.Stderr = w

	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: "15:04:05",
			NoColor:    os.Getenv("NO_COLOR") == "1",
		}),
	))

	RESET := "\u001B[0m"
	BOLD := "\u001B[1m"
	FAINT := "\u001B[2m"

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

	mylog.Success("terminal test success")
	mylog.Info("terminal test info")
	mylog.Warning("terminal test warn")
	return screen, settings
}

func randomString(n int) string {
	s := make([]rune, n)
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := 0; i < n; i++ {
		s[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(s)
}
