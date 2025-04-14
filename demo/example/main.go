package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"regexp"
	"strings"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/oligo/gvcode"
	"github.com/oligo/gvcode/addons/completion"
	wg "github.com/oligo/gvcode/widget"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type EditorApp struct {
	window *app.Window
	th     *material.Theme
	state  *gvcode.Editor
	popup  completion.CompletionPopup
}

const (
	syntaxPattern = "package|import|type|func|struct|for|var|switch|case|if|else"
)

func (ed *EditorApp) run() error {

	var ops op.Ops
	for {
		e := ed.window.Event()

		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx C) D {
				return ed.layout(gtx, ed.th)
			})
			e.Frame(gtx.Ops)
		}
	}
}

func (ed *EditorApp) layout(gtx C, th *material.Theme) D {
	for {
		evt, ok := ed.state.Update(gtx)
		if !ok {
			break
		}

		switch evt.(type) {
		case gvcode.ChangeEvent:
			styles := HightlightTextByPattern(ed.state.Text(), syntaxPattern)
			ed.state.UpdateTextStyles(styles)
		}
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			lb := material.Label(th, th.TextSize, "gvcode editor")
			lb.Alignment = text.Middle
			return lb.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
		layout.Flexed(1, func(gtx C) D {
			borderColor := th.Fg
			borderColor.A = 0xb6
			return widget.Border{
				Color: borderColor, Width: unit.Dp(1),
			}.Layout(gtx, func(gtx C) D {
				return layout.Inset{
					Top:    unit.Dp(6),
					Bottom: unit.Dp(6),
					Left:   unit.Dp(24),
					Right:  unit.Dp(24),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					es := wg.NewEditor(th, ed.state)
					es.Font.Typeface = "monospace"
					es.Font.Weight = font.SemiBold
					es.TextSize = unit.Sp(12)
					es.LineHeightScale = 1.5
					es.TextHighlightColor = color.NRGBA{R: 120, G: 120, B: 120, A: 200}

					return es.Layout(gtx)
				})
			})
		}),
		layout.Rigid(func(gtx C) D {
			line, col := ed.state.CaretPos()
			lb := material.Label(th, th.TextSize*0.8, fmt.Sprintf("Line:%d, Col:%d", line+1, col+1))
			lb.Alignment = text.End
			return lb.Layout(gtx)
		}),
	)

}

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	th := material.NewTheme()
	th.ContrastBg = color.NRGBA{R: 40, G: 204, B: 187, A: 255}

	editorApp := EditorApp{
		window: &app.Window{},
		th:     th,
	}
	editorApp.window.Option(app.Title("Basic Example"))

	gvcode.SetDebug(false)
	editorApp.state = &gvcode.Editor{}
	editorApp.state.WithOptions(
		gvcode.WrapLine(true),
	)

	var quotePairs = map[rune]rune{
		'\'': '\'',
		'"':  '"',
		'`':  '`',
		'“':  '”',
	}

	// Bracket pairs
	var bracketPairs = map[rune]rune{
		'(': ')',
		'{': '}',
		'[': ']',
		'<': '>',
	}

	thisFile, _ := os.ReadFile("./main.go")
	editorApp.state.SetText(string(thisFile))
	// editorApp.state.SetHighlights([]editor.TextRange{{Start: 0, End: 5}})
	styles := HightlightTextByPattern(editorApp.state.Text(), syntaxPattern)
	editorApp.state.UpdateTextStyles(styles)

	// Setting up auto-completion.
	cm := &completion.DefaultCompletion{Editor: editorApp.state}
	// set completion triggers
	cm.SetTriggers(
		gvcode.AutoTrigger{},
		gvcode.KeyTrigger{Name: "P", Modifiers: key.ModShortcut})
	// set the completion algorithms
	cm.SetCompletors(&goCompletor{})
	// set popup widget to let user navigate the candidates.
	editorApp.popup = *completion.NewCompletionPopup(editorApp.state, cm)
	cm.SetPopup(func(gtx layout.Context, items []gvcode.CompletionCandicate) layout.Dimensions {
		editorApp.popup.TextSize = unit.Sp(12)
		editorApp.popup.Size = image.Point{
			X: gtx.Dp(unit.Dp(400)),
			Y: gtx.Dp(unit.Dp(200)),
		}

		return editorApp.popup.Layout(gtx, th, items)
	})

	editorApp.state.WithOptions(
		gvcode.WithSoftTab(true),
		gvcode.WithQuotePairs(quotePairs),
		gvcode.WithBracketPairs(bracketPairs),
		gvcode.WithAutoCompletion(cm),
	)

	go func() {
		err := editorApp.run()
		if err != nil {
			os.Exit(1)
		}

		os.Exit(0)
	}()

	app.Main()

}

func HightlightTextByPattern(text string, pattern string) []*gvcode.TextStyle {
	var styles []*gvcode.TextStyle

	re := regexp.MustCompile(pattern)
	matches := re.FindAllIndex([]byte(text), -1)
	for _, match := range matches {
		styles = append(styles, &gvcode.TextStyle{
			TextRange: gvcode.TextRange{
				Start: match[0],
				End:   match[1],
			},
			Color:      rgbaToOp(color.NRGBA{R: 255, A: 255}),
			Background: rgbaToOp(color.NRGBA{R: 215, G: 215, B: 215, A: 250}),
		})
	}

	return styles
}

func rgbaToOp(textColor color.NRGBA) op.CallOp {
	ops := new(op.Ops)

	m := op.Record(ops)
	paint.ColorOp{Color: textColor}.Add(ops)
	return m.Stop()
}

var golangKeywords = []string{
	"break",
	"default",
	"func",
	"interface",
	"select",
	"case",
	"defer", "go", "map", "struct",
	"chan", "else", "goto", "package", "switch",
	"const", "fallthrough", "if", "range", "type",
	"continue", "for", "import", "return", "var",
}

type goCompletor struct {
}

func (c *goCompletor) Suggest(ctx gvcode.CompletionContext) []gvcode.CompletionCandicate {
	candicates := make([]gvcode.CompletionCandicate, 0)
	for _, kw := range golangKeywords {
		if strings.Contains(kw, ctx.Input) {
			candicates = append(candicates, gvcode.CompletionCandicate{
				Label:       kw,
				InsertText:  kw,
				Description: kw,
				Kind:        "text",
			})
		}
	}

	return candicates
}
