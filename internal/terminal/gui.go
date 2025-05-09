package terminal

import (
	_ "embed"
	"image"
	"image/color"
	"strings"

	"github.com/ddkwork/ux"

	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/x/styledtext"
)

func Console(screen *Screen, settings *ConsoleSettings) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		// Figure out character height
		gtx = settings.update(screen, gtx)

		return bordered(gtx, 1, color.NRGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 128,
		}, func(gtx layout.Context) layout.Dimensions {
			evt, got := gtx.Source.Event(pointer.Filter{
				Target: settings.scrollTag,
				Kinds:  pointer.Scroll,
				ScrollX: pointer.ScrollRange{ // todo test
					Min: -100,
					Max: 100,
				},
				ScrollY: pointer.ScrollRange{
					Min: -100,
					Max: 100,
				},
			})
			if got {
				switch evt := evt.(type) {
				case pointer.Event:
					screen.scrollTop = min(
						max(
							screen.scrollTop+int(evt.Scroll.Y),
							0,
						),
						screen.top,
					)
				}
			}

			defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, screen.defaults.BgColor)

			// Declare the tag
			event.Op(gtx.Ops, settings.scrollTag)

			return layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top:    settings.paddingY,
						Right:  settings.paddingX,
						Bottom: settings.paddingY,
						Left:   settings.paddingX,
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						t := styledtext.Text(ux.NewTheme().Shaper, createSpansFrom(screen)...)
						t.WrapPolicy = styledtext.WrapGraphemes
						return t.Layout(gtx, nil)
					})
				}),

				// Render the scrollbar
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					// scrollbar area
					offset := gtx.Dp(8)
					area := image.Rectangle{
						Min: image.Point{X: gtx.Constraints.Max.X - int(float32(offset)*1.7), Y: offset},
						Max: image.Point{X: gtx.Constraints.Max.X - offset, Y: gtx.Constraints.Max.Y - offset},
					}

					l := float32(screen.Size.Y) / float32(screen.top+screen.Size.Y)

					if l < 1.0 {
						total := float32(area.Dy())
						height := total * l

						offsetTop := (total - height) * float32(screen.scrollTop) / float32(screen.top)

						bar := area
						bar.Max.Y = area.Min.Y + int(height)
						bar = bar.Add(image.Point{Y: int(offsetTop)})

						defer clip.RRect{
							Rect: bar,
							SE:   gtx.Dp(3),
							SW:   gtx.Dp(3),
							NW:   gtx.Dp(3),
							NE:   gtx.Dp(3),
						}.Push(gtx.Ops).Pop()

						paint.Fill(gtx.Ops, color.NRGBA{
							R: 255,
							G: 255,
							B: 255,
							A: 20,
						})
					}

					return layout.Dimensions{
						Size:     image.Point{},
						Baseline: 0,
					}
				}),
			)
		})
	}
}

func bordered(gtx layout.Context, width unit.Dp, c color.NRGBA, f layout.Widget) layout.Dimensions {
	defer clip.Rect{
		Max: gtx.Constraints.Max,
	}.Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, c)
	return layout.UniformInset(width).Layout(gtx, f)
}

func createSpansFrom(screen *Screen) []styledtext.SpanStyle {
	var spans []styledtext.SpanStyle

	for _, line := range screen.VisibleLines() {
		x := 0
		for _, run := range line.runs {
			if run.start > x {
				spans = append(spans, styledtext.SpanStyle{Content: strings.Repeat(" ", run.start-x)})
			}

			f := screen.defaults.Font
			if run.style.Bold {
				f = screen.defaults.BoldFont
			}

			spans = append(spans, styledtext.SpanStyle{
				Content: run.text,
				Size:    screen.defaults.FontSize,
				Color:   run.style.FgColor(),
				Font:    f,
			})

			x = run.end()
		}
		spans = append(spans, styledtext.SpanStyle{Content: "\n"})
	}

	return spans
}
