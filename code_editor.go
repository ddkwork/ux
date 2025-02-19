package ux

import (
	"bytes"
	_ "embed"
	"gioui.org/font/opentype"
	"github.com/ddkwork/golibrary/mylog"
	"image/color"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	giovieweditor "github.com/oligo/gioview/editor"
)

const (
	CodeLanguageJSON       = "JSON"
	CodeLanguageYAML       = "YAML"
	CodeLanguageXML        = "XML"
	CodeLanguagePython     = "Python"
	CodeLanguageGolang     = "Golang"
	CodeLanguageJava       = "Java"
	CodeLanguageJavaScript = "JavaScript"
	CodeLanguageRuby       = "Ruby"
	CodeLanguageShell      = "Shell"
	CodeLanguageDotNet     = "Shell"
	CodeLanguageProperties = "properties"
)

type CodeEditor struct {
	editor *giovieweditor.Editor
	code   string

	styledCode string
	styles     []*giovieweditor.TextStyle

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

	editorConf      *giovieweditor.EditorConf
	vScrollbar      widget.Scrollbar
	vScrollbarStyle material.ScrollbarStyle
}

//go:embed resources/dracula.xml
var DraculaXML []byte
var Dracula *chroma.Style

//go:embed resources/fonts/consolas.ttf
var consolas []byte

func init() {
	Dracula = mylog.Check2(chroma.NewXMLStyle(bytes.NewReader(DraculaXML)))
}

func MustGetCodeEditorFont() font.FontFace {
	fontFaces := mylog.Check2(opentype.ParseCollection(consolas))
	return font.FontFace{
		Font: fontFaces[0].Font,
		Face: fontFaces[0].Face,
	}
}

func NewCodeEditor(code string, lang string) *CodeEditor {
	code = strings.ReplaceAll(code, "\t", "    ")
	editorFont := MustGetCodeEditorFont()
	shaper := text.NewShaper(text.WithCollection([]font.FontFace{editorFont}))

	c := &CodeEditor{
		editor: new(giovieweditor.Editor),
		code:   code,
		font:   MustGetCodeEditorFont(),
		lang:   lang,

		editorConf: &giovieweditor.EditorConf{
			Shaper:          shaper,
			TextColor:       th.Fg,
			Bg:              th.Bg,
			SelectionColor:  Blue500, //todo
			TypeFace:        editorFont.Font.Typeface,
			TextSize:        unit.Sp(14),
			LineHeightScale: 1,
			ShowLineNum:     true,
			LineNumPadding:  unit.Dp(10),
			LineHighlightColor: color.NRGBA{
				R: 0xbb,
				G: 0xbb,
				B: 0xbb,
				A: 0x33,
			},
		},
	}

	c.vScrollbarStyle = material.Scrollbar(th.Theme, &c.vScrollbar)
	c.border = widget.Border{
		Color:        rgb(0x6c6f76), //todo
		Width:        unit.Dp(1),
		CornerRadius: unit.Dp(4),
	}

	c.lexer = getLexer(lang)

	style := styles.Get("dracula")
	if style == nil {
		style = styles.Fallback
	}

	c.codeStyle = style //todo delete this
	c.codeStyle = Dracula

	c.editor.WrapPolicy = text.WrapGraphemes
	c.editor.SetText(code, false)

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
	c.editor.ReadOnly = readOnly
}

func (c *CodeEditor) SetOnLoadExample(f func()) {
	c.onLoadExample = f
}

func (c *CodeEditor) AppendText(text string) {
	c.editor.Insert(text)
}

func (c *CodeEditor) SetCode(code string) {
	c.editor.SetText(code, false)
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

	if !c.editor.ReadOnly {
		if ev, ok := c.editor.Update(gtx); ok {
			if _, ok := ev.(giovieweditor.ChangeEvent); ok {
				c.editor.UpdateTextStyles(c.stylingText(c.editor.Text()))
				if c.onChange != nil {
					c.onChange(c.editor.Text())
					c.code = c.editor.Text()
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

					//btn := Button(th.Material(), &c.loadExample, RefreshIcon, IconPositionStart, "Load Example")
					btn := NewNavButton("Load Example")

					//btn.Color = theme.ButtonTextColor
					//btn.Inset = layout.Inset{
					//	Top: unit.Dp(4), Bottom: unit.Dp(4),
					//	Left: unit.Dp(4), Right: unit.Dp(4),
					//}

					return btn.Layout(gtx, th.Theme)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if c.onBeautify == nil {
						return layout.Dimensions{}
					}
					btn := NewNavButton("Beautify")

					//btn := Button(theme.Material(), &c.beatufier, CleanIcon, IconPositionStart, "Beautify")
					//btn.Color = theme.ButtonTextColor
					//btn.Inset = layout.Inset{
					//	Top: 4, Bottom: 4,
					//	Left: 4, Right: 4,
					//}

					return btn.Layout(gtx, th.Theme)
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

func (c *CodeEditor) editorStyle(gtx layout.Context, hint string) layout.Dimensions {
	editorDims := giovieweditor.NewEditor(c.editor, c.editorConf, hint).Layout(gtx)

	layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		viewportStart, viewportEnd := c.editor.ViewPortRatio()
		return c.vScrollbarStyle.Layout(gtx, layout.Vertical, viewportStart, viewportEnd)
	})

	if delta := c.vScrollbar.ScrollDistance(); delta != 0 {
		c.editor.ScrollByRatio(gtx, delta)
	}

	return editorDims
}

func (c *CodeEditor) stylingText(text string) []*giovieweditor.TextStyle {
	if c.styledCode == text {
		return c.styles
	}

	// nolint:prealloc
	var textStyles []*giovieweditor.TextStyle

	offset := 0

	iterator, err := c.lexer.Tokenise(nil, text)
	if err != nil {
		return textStyles
	}

	for _, token := range iterator.Tokens() {
		entry := c.codeStyle.Get(token.Type)

		textStyle := &giovieweditor.TextStyle{
			Start: offset,
			End:   offset + len([]rune(token.Value)),
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
