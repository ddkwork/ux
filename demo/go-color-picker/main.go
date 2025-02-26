package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

func main() {
	fmt.Println("Please mouse left click to pick the color of current mouse position.")
	fmt.Println("To pick the color of mouse positions, please enter ctrl+Shift+c.")
	fmt.Println("To quit checking mouse event, please enter ctrl+Shift+q.")

	// Stop event
	EventHook(hook.KeyDown, []string{"q", "shift", "ctrl"}, func(e hook.Event) {
		EventEnd()
	})

	// Pick the color of event
	EventHook(hook.KeyDown, []string{"c", "shift", "ctrl"}, func(e hook.Event) {
		posx, posy := robotgo.GetMousePos()
		color := robotgo.GetPixelColor(posx, posy)
		robotgo.WriteAll(color)

		fmt.Println("Pick the color: ", color)
	})

	// Get color of current mouse position
	EventHook(hook.MouseDown, []string{}, func(e hook.Event) {
		posx, posy := robotgo.GetMousePos()
		color := robotgo.GetPixelColor(posx, posy)

		fmt.Println("pos-> ", posx, posy)
		fmt.Println("color-> ", color)
	})

	s := EventStart()
	<-EventProcess(s)
}
