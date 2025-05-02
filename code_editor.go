package ux

import (
	"bytes"
	_ "embed"
	"image/color"

	"github.com/ddkwork/ux/languages"
	"github.com/ddkwork/ux/resources/images"

	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/ux/widget/material"
	"github.com/oligo/gvcode"
)

type CodeEditor struct {
	editor *gvcode.Editor
	code   string

	styledCode string
	styles     []*gvcode.TextStyle

	lexer     chroma.Lexer
	codeStyle *chroma.Style

	lang string

	onChange func(text string)

	font font.FontFace

	border widget.Border

	beatufier   widget.Clickable
	loadExample widget.Clickable

	onBeautify    func()
	onLoadExample func()

	vScrollbar      widget.Scrollbar
	vScrollbarStyle material.ScrollbarStyle
}

var (
	//go:embed resources/dracula.xml
	DraculaXML []byte
	Dracula    = mylog.Check2(chroma.NewXMLStyle(bytes.NewReader(DraculaXML)))

	//go:embed resources/fonts/consolas.ttf
	consolas []byte
)

func MustGetCodeEditorFont() font.FontFace {
	fontFaces := mylog.Check2(opentype.ParseCollection(consolas))
	return font.FontFace{
		Font: fontFaces[0].Font,
		Face: fontFaces[0].Face,
	}
}

func NewCodeEditor(code string, lang ...languages.LanguagesKind) *CodeEditor {
	l := "golang"
	if len(lang) > 0 {
		l = lang[0].String()
	}
	// 	editorFont := fonts.MustGetCodeEditorFont()
	c := &CodeEditor{
		editor: &gvcode.Editor{
			// Font:                  editorFont.Font,
			// TextSize:              unit.Sp(12),
			// LineHeightScale:       1,
			// WrapLine:              true,
			// ReadOnly:              false,
			// SoftTab:               true,
			// TabWidth:              4,
			// LineNumberGutter: 1,
			// TextMaterial:     rgbToOp(theme.TextColor),
			// SelectMaterial:        rgbToOp(theme.TextSelectionColor),
			// TextHighlightMaterial: rgbToOp(theme.TextSelectionColor),
		},
		code: code,
		font: font.FontFace{
			Font: font.Font{Typeface: "sourceSansPro"},
		},
		lang: l,
	}

	// c.editor.WithOptions(gvcode.WithShaperParams(c.font.Font, unit.Sp(12), text.Start, unit.Sp(16), 1))
	// c.editor.WithOptions(gvcode.WithTabWidth(4))
	// c.editor.WithOptions(gvcode.WithSoftTab(true))
	c.editor.WithOptions(gvcode.WrapLine(true))
	c.vScrollbarStyle = material.Scrollbar(th, &c.vScrollbar)
	c.border = widget.Border{
		Color:        rgb(0x6c6f76), // todo
		Width:        unit.Dp(1),
		CornerRadius: unit.Dp(4),
	}

	c.lexer = getLexer(l)
	// style := styles.Get("dracula")
	// if style == nil {
	//	style = styles.Fallback
	// }
	// c.codeStyle = style
	c.codeStyle = Dracula
	c.editor.SetText(code)
	return c
}

func getLexer(lang string) chroma.Lexer {
	lexer := lexers.Get(lang)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	return chroma.Coalesce(lexer)
}

func (c *CodeEditor) SetOnChanged(f func(text string)) {
	c.onChange = f
}

func (c *CodeEditor) SetOnBeautify(f func()) {
	c.onBeautify = f
}

func (c *CodeEditor) SetReadOnly(readOnly bool) {
	c.editor.WithOptions(gvcode.ReadOnlyMode(readOnly))
}

func (c *CodeEditor) SetOnLoadExample(f func()) {
	c.onLoadExample = f
}

func (c *CodeEditor) AppendText(text string) {
	c.editor.Insert(text)
}

func (c *CodeEditor) SetCode(code string) {
	c.editor.SetText(code)
	c.code = code
	c.editor.UpdateTextStyles(c.stylingText(c.editor.Text()))
}

func (c *CodeEditor) SetLanguage(lang string) {
	c.lang = lang
	c.lexer = getLexer(lang)
	c.editor.UpdateTextStyles(c.stylingText(c.editor.Text()))
}

func (c *CodeEditor) Code() string {
	return c.editor.Text()
}

func (c *CodeEditor) Layout(gtx layout.Context) layout.Dimensions {
	if c.styledCode == "" {
		// First time styling
		c.editor.UpdateTextStyles(c.stylingText(c.editor.Text()))
	}

	if !c.editor.ReadOnly() {
		if ev, ok := c.editor.Update(gtx); ok {
			if _, ok := ev.(gvcode.ChangeEvent); ok {
				st := c.stylingText(c.editor.Text())
				c.styles = st
				c.editor.UpdateTextStyles(st)
				if c.onChange != nil {
					c.onChange(c.editor.Text())
					c.code = c.editor.Text()
				}
			}
		}

		filters := []event.Filter{
			key.FocusFilter{Target: c.editor},
			key.Filter{Focus: c.editor, Name: "D", Required: key.ModShortcut},
		}
		for {
			ke, ok := gtx.Event(filters...)
			if !ok {
				break
			}
			switch e := ke.(type) {
			case key.Event:
				if !gtx.Focused(c.editor) || e.State != key.Press {
					break
				}
				if e.Modifiers.Contain(key.ModShortcut) {
					switch e.Name {
					case "D":
						// c.editor.DuplicateLine()
						// c.editor.DuplicateLine()
						// func (e *Editor) DuplicateLine() {
						//	e.initBuffer()
						//	if e.text.SelectionLen() == 0 {
						//		e.scratch = e.text.SelectedLine(e.scratch)
						//		if len(e.scratch) > 0 && e.scratch[len(e.scratch)-1] != '\n' {
						//			e.scratch = append(e.scratch, '\n')
						//		}
						//	} else {
						//		e.scratch = e.text.SelectedText(e.scratch)
						//	}
						//	e.Insert(string(e.scratch))
						// }
						c.editor.UpdateTextStyles(c.stylingText(c.editor.Text()))
					}
				}
			}
		}
	}

	if c.loadExample.Clicked(gtx) {
		c.onLoadExample()
	}

	if c.beatufier.Clicked(gtx) {
		c.onBeautify()
	}

	flexH := layout.Flex{Axis: layout.Horizontal}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if c.onLoadExample == nil {
						return layout.Dimensions{}
					}

					btn := Button(&c.loadExample, images.NavigationRefreshIcon, "Load Example")
					// btn := NewNavButton("Load Example")
					// btn.Color = theme.ButtonTextColor
					// btn.Inset = layout.Inset{
					//	Top: unit.Dp(4), Bottom: unit.Dp(4),
					//	Left: unit.Dp(4), Right: unit.Dp(4),
					// }

					return btn.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if c.onBeautify == nil {
						return layout.Dimensions{}
					}

					btn := Button(&c.beatufier, images.EditorFormatColorTextIcon, "Beautify")
					// btn.Color = theme.ButtonTextColor//todo
					btn.Inset = layout.Inset{
						Top: 4, Bottom: 4,
						Left: 4, Right: 4,
					}
					// btn := NewNavButton("Load Example")
					return btn.Layout(gtx)
				}),
			)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return c.border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return flexH.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top:    unit.Dp(4),
							Bottom: unit.Dp(4),
							Left:   unit.Dp(8),
							Right:  unit.Dp(4),
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							hint := "No example available"
							return c.editorStyle(gtx, hint)
						})
					}),
				)
			})
		}),
	)
}

func (c *CodeEditor) editorStyle(gtx layout.Context, _ string) layout.Dimensions {
	es := NewEditor(c.editor)
	es.Font.Typeface = "sourceSansPro"
	es.SelectionColor = th.TextSelectionColor
	editorDims := es.Layout(gtx)

	layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		viewportStart, viewportEnd := c.editor.ViewPortRatio()
		return c.vScrollbarStyle.Layout(gtx, layout.Vertical, viewportStart, viewportEnd)
	})

	if delta := c.vScrollbar.ScrollDistance(); delta != 0 {
		c.editor.ScrollByRatio(gtx, delta)
	}

	return editorDims
}

func (c *CodeEditor) stylingText(text string) []*gvcode.TextStyle {
	if c.styledCode == text {
		return c.styles
	}

	// nolint:prealloc
	var textStyles []*gvcode.TextStyle

	offset := 0

	iterator := mylog.Check2(c.lexer.Tokenise(nil, text))

	for _, token := range iterator.Tokens() {
		entry := c.codeStyle.Get(token.Type)

		textStyle := &gvcode.TextStyle{
			TextRange: gvcode.TextRange{
				Start: offset,
				End:   offset + len([]rune(token.Value)),
			},
			Color: rgbToOp(th.Fg),
			// Background: rgbToOp(c.theme.Bg),
		}

		if entry.Colour.IsSet() {
			textStyle.Color = chromaColorToOp(entry.Colour)
		}

		textStyles = append(textStyles, textStyle)
		offset = textStyle.End
	}

	c.styledCode = text
	c.styles = textStyles

	return textStyles
}

func chromaColorToOp(textColor chroma.Colour) op.CallOp {
	ops := new(op.Ops)

	m := op.Record(ops)
	paint.ColorOp{Color: color.NRGBA{
		R: textColor.Red(),
		G: textColor.Green(),
		B: textColor.Blue(),
		A: 0xff,
	}}.Add(ops)
	return m.Stop()
}

func rgbToOp(color color.NRGBA) op.CallOp {
	ops := new(op.Ops)

	m := op.Record(ops)
	paint.ColorOp{Color: color}.Add(ops)
	return m.Stop()
}
