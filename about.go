package ux

import (
	"fmt"
	"image"
	"image/color"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/ddkwork/ux/widget/material"
)

// A []string to hold the speech as a list of paragraphs
var paragraphList []string

// Colors
type colorMode struct {
	background color.NRGBA
	foreground color.NRGBA
	focusbar   color.NRGBA
}

func readText(data string) []string {
	text := strings.Split(data, "\n")
	for i := 1; i <= 10; i++ {
		text = append(text, "")
	}
	return text
}

// The main draw function
func About(data string) error {
	w := new(app.Window)
	w.Option(
		app.Title("about"),
		app.Size(1200, 600),
		// app.Decorated(false),
	)
	w.Perform(system.ActionCenter)
	paragraphList = readText(data)
	// y-position for text
	var scrollY unit.Dp = 0

	// y-position for red focusBar
	var focusBarY unit.Dp = 170

	// maxIndentWidth of text area
	var textWidth unit.Dp = 550

	// fontSize
	var fontSize unit.Sp = 35

	// Are we auto scrolling?
	var autoscroll bool = false
	var autospeed unit.Dp = 1

	// th defines the material design style
	th := material.NewTheme()

	// ops are the operations from the UI
	var ops op.Ops

	// Define a tag for input routing
	tag := "My Input Routing Tag - which could be this silly string, or an int/float/address, or anything else"

	// Colors
	colorDark := colorMode{
		background: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
		foreground: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		focusbar:   color.NRGBA{R: 0xff, G: 0x00, B: 0x00, A: 0x33},
	}

	colorLight := colorMode{
		background: color.NRGBA{R: 0xff, G: 0xfe, B: 0xe0, A: 0xff},
		foreground: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
		focusbar:   color.NRGBA{R: 0xff, A: 0x66},
	}

	// Define a color to start with. We like dark
	myColor := colorDark

	for {
		// listen for events in the w
		switch winE := w.Event().(type) {

		// Should we draw a new frame?
		case app.FrameEvent:
			gtx := app.NewContext(&ops, winE)

			// ---------- Handle input ----------
			// Time to deal with inputs since last frame.

			// Scrolled a mouse wheel?
			for {
				ev, ok := gtx.Event(
					pointer.Filter{
						Target:  tag,
						Kinds:   pointer.Scroll,
						ScrollY: pointer.ScrollRange{Min: -1, Max: +1},
					},
				)
				if !ok {
					break
				}
				fmt.Printf("SCROLL: %+v\n", ev)
				scrollY = max(scrollY+unit.Dp(ev.(pointer.Event).Scroll.Y*float32(fontSize)), 0)
			}

			// Pressed a mouse button?
			for {
				ev, ok := gtx.Event(
					pointer.Filter{
						Target: tag,
						Kinds:  pointer.Press,
					},
				)
				if !ok {
					break
				}
				fmt.Printf("PRESS : %+v\n", ev)
				// Start / stop
				autoscroll = !autoscroll
			}

			// Pressed a key?
			for {
				ev, ok := gtx.Event(
					key.Filter{Name: key.NameSpace},
					key.Filter{Optional: key.ModShift, Name: "U"},
					key.Filter{Optional: key.ModShift, Name: "D"},
					key.Filter{Optional: key.ModShift, Name: "J"},
					key.Filter{Optional: key.ModShift, Name: "K"},
					key.Filter{Optional: key.ModShift, Name: key.NameUpArrow},
					key.Filter{Optional: key.ModShift, Name: key.NameDownArrow},
					key.Filter{Optional: key.ModShift, Name: key.NamePageUp},
					key.Filter{Optional: key.ModShift, Name: key.NamePageDown},
					key.Filter{Optional: key.ModShift, Name: "F"},
					key.Filter{Optional: key.ModShift, Name: "S"},
					key.Filter{Optional: key.ModShift, Name: "+"},
					key.Filter{Optional: key.ModShift, Name: "-"},
					key.Filter{Optional: key.ModShift, Name: "W"},
					key.Filter{Optional: key.ModShift, Name: "N"},
					key.Filter{Optional: key.ModShift, Name: "C"},
				)
				if !ok {
					break
				}
				fmt.Printf("KEY   : %+v\n", ev)
				if ev.(key.Event).State == key.Press {
					name := ev.(key.Event).Name
					mod := ev.(key.Event).Modifiers

					// Set stepsize
					var stepSize unit.Dp = 1
					if mod == key.ModShift {
						stepSize = 5
					}

					// Start / stop
					if name == key.NameSpace {
						autoscroll = !autoscroll
						if autoscroll && autospeed <= 0 {
							autospeed = stepSize
						}
					}

					// Move the focusBar Up
					if name == "U" {
						focusBarY = focusBarY - stepSize
					}

					// Move the focusBar Down
					if name == "D" {
						focusBarY = focusBarY + stepSize
					}

					// List up
					if name == "K" || name == key.NameUpArrow {
						scrollY = scrollY - stepSize*4
					}
					if name == key.NamePageUp {
						scrollY = scrollY - stepSize*100
					}
					if scrollY < 0 {
						scrollY = 0
					}

					// List down
					if name == "J" || name == key.NameDownArrow || name == key.NamePageDown {
						scrollY = scrollY + stepSize*4
					}
					if name == key.NamePageDown {
						scrollY = scrollY + stepSize*100
					}

					// Faster scrollspeed
					if name == "F" {
						autoscroll = true
						autospeed += stepSize
					}

					// Slower scrollspeed
					if name == "S" {
						if autospeed > 0 {
							autospeed -= stepSize
						}
						if autospeed <= 0 {
							autospeed = 0
							autoscroll = false
						}
					}

					// To increase the fontsize
					if name == "+" {
						fontSize = fontSize + unit.Sp(stepSize)
					}

					// To decrease the fontsize
					if name == "-" {
						fontSize = fontSize - unit.Sp(stepSize)
					}

					// Widen text to be displayed
					if name == "W" {
						textWidth = textWidth + stepSize*10
					}
					// Narrow text to be displayed
					if name == "N" {
						textWidth = textWidth - stepSize*10
					}

					// Swhich Colormode
					if name == "C" {
						if myColor == colorDark {
							myColor = colorLight
						} else {
							myColor = colorDark
						}
					}
				}
			}

			// ---------- LAYOUT ----------
			// First we DefaultDraw the user interface.
			// Let's start with a background color
			paint.Fill(&ops, myColor.background)

			// ---------- THE SCROLLING TEXT ----------
			// First, check if we should autoscroll
			// That's done by increasing the value of scrollY
			if autoscroll {
				if autospeed < 0 {
					autospeed = 0
				}
				scrollY = scrollY + autospeed
				// Invalidate 50 times per second
				inv := op.InvalidateCmd{At: gtx.Now.Add(time.Second / 50)}
				gtx.Execute(inv)
			}
			// Then we use scrollY to control the distance from the top of the screen to the first element.
			// We visualize the text using a list where each paragraph is a separate item.
			vizList := layout.List{
				Axis: layout.Vertical,
				Position: layout.Position{
					Offset: int(scrollY),
				},
			}

			// ---------- MARGINS ----------
			// Margins
			var marginWidth unit.Dp
			marginWidth = (unit.Dp(gtx.Constraints.Max.X) - textWidth) / 3
			margins := layout.Inset{
				Left:   marginWidth,
				Right:  marginWidth,
				Top:    unit.Dp(0),
				Bottom: unit.Dp(0),
			}

			// ---------- LIST WITHIN MARGINS ----------
			// 1) First the margins ...
			margins.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					// 2) ... then the list inside those margins ...
					return vizList.Layout(gtx, len(paragraphList),
						// 3) ... where each paragraph is a separate item
						func(gtx layout.Context, index int) layout.Dimensions {
							// One label per paragraph
							paragraph := material.Label(th, unit.Sp(fontSize), paragraphList[index])
							// The text is centered
							paragraph.Alignment = text.Middle
							// Set color
							paragraph.Color = myColor.foreground
							// Return the laid out paragraph
							return paragraph.Layout(gtx)
						},
					)
				},
			)

			// ---------- THE FOCUS BAR ----------
			// DefaultDraw the transparent red focus bar.
			focusBar := clip.Rect{
				Min: image.Pt(0, int(focusBarY)),
				Max: image.Pt(gtx.Constraints.Max.X, int(focusBarY)+int(fontSize*1.5)),
			}.Push(&ops)
			paint.ColorOp{Color: myColor.focusbar}.Add(&ops)
			paint.PaintOp{}.Add(&ops)
			focusBar.Pop()

			// ---------- REGISTERING EVENTS ----------
			// registering events here work
			event.Op(&ops, tag)

			// ---------- FINALIZE ----------
			// Frame completes the FrameEvent by drawing the graphical operations from ops into the w.
			winE.Frame(&ops)

			// Should we shut down?
		case app.DestroyEvent:
			return winE.Err
		}
	}
}
