package ux

import (
	"bytes"
	_ "embed"
	"image/color"
	"strings"

	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/ddkwork/golibrary/mylog"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	giovieweditor "github.com/oligo/gioview/editor"
)

const (
	CodeLanguageJSON   = "JSON"
	CodeLanguageYAML   = "YAML"
	CodeLanguageXML    = "XML"
	CodeLanguagePython = "Python"
	CodeLanguageGO     = "GO"
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
}

//go:embed dracula.xml
var DraculaXML []byte
var Dracula *chroma.Style

func init() {
	Dracula = mylog.Check2(chroma.NewXMLStyle(bytes.NewReader(DraculaXML)))
}

//go:embed consolas.ttf
var consolas []byte

func MustGetCodeEditorFont() font.FontFace {
	fontFaces := mylog.Check2(opentype.ParseCollection(consolas))
	return font.FontFace{
		Font: fontFaces[0].Font,
		Face: fontFaces[0].Face,
	}
}

func NewCodeEditor(code string, lang string) *CodeEditor {
	code = strings.ReplaceAll(code, "\t", "    ")
	c := &CodeEditor{
		editor:        new(giovieweditor.Editor),
		code:          code,
		styledCode:    "",
		styles:        nil,
		lexer:         nil,
		codeStyle:     nil,
		font:          MustGetCodeEditorFont(),
		lang:          lang,
		onChange:      nil,
		border:        widget.Border{},
		beatufier:     widget.Clickable{},
		loadExample:   widget.Clickable{},
		onBeautify:    nil,
		onLoadExample: nil,
	}
	c.border = widget.Border{
		Color:        th.Bg,
		Width:        unit.Dp(1),
		CornerRadius: unit.Dp(4),
	}
	c.lexer = getLexer(lang)
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

func (c *CodeEditor) SetOnLoadExample(f func()) {
	c.onLoadExample = f
}

func (c *CodeEditor) AppendText(text string) {
	c.editor.Insert(text)
}

func (c *CodeEditor) SetCode(code string) {
	code = strings.ReplaceAll(code, "\t", "    ")
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

	if e, ok := c.editor.Update(gtx); ok {
		switch e.(type) {
		case giovieweditor.ChangeEvent:
			c.editor.UpdateTextStyles(c.stylingText(c.editor.Text()))
			if c.onChange != nil {
				c.onChange(c.editor.Text())
				c.code = c.editor.Text()
			}
			// case key.Event://todo ctrl+x +d use insert method
		}
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

					// btn := Button(theme.Theme(), &c.loadExample, RefreshIcon, IconPositionStart, "Load Example")
					// btn.Color = theme.ButtonTextColor
					btn := NewNavButton("Load Example")
					//btn.Inset = DefaultDraw.Inset{
					//	Top: unit.Dp(4), Bottom: unit.Dp(4),
					//	Left: unit.Dp(4), Right: unit.Dp(4),
					//}

					if c.loadExample.Clicked(gtx) {
						c.onLoadExample()
					}

					return btn.Layout(gtx, th.Theme)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if c.onBeautify == nil {
						return layout.Dimensions{}
					}

					// btn := Button(theme.Theme(), &c.beatufier, CleanIcon, IconPositionStart, "Beautify")
					// btn.Color = theme.ButtonTextColor
					btn := NewNavButton("Beautify")
					//btn.Inset = DefaultDraw.Inset{
					//	Top: 4, Bottom: 4,
					//	Left: 4, Right: 4,
					//}

					if c.beatufier.Clicked(gtx) {
						c.onBeautify()
					}

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
							editorConf := &giovieweditor.EditorConf{
								Shaper:         th.Theme.Shaper,
								TextColor:      th.Theme.Fg,
								Bg:             BackgroundColor,
								SelectionColor: Blue500, // todo
								LineHighlightColor: color.NRGBA{
									R: 0xbb,
									G: 0xbb,
									B: 0xbb,
									A: 0x33,
								},
								TypeFace:        c.font.Font.Typeface,
								TextSize:        unit.Sp(14),
								Weight:          0,
								LineHeight:      0,
								LineHeightScale: 1.2,
								ColorScheme:     "",
								ShowLineNum:     true,
								LineNumPadding:  unit.Dp(10),
							}
							//lineInfos := mylog.Check2(c.editor.VisibleLines())
							//c.ScrollBar.Layout(gtx, 5, len(lineInfos), &DefaultDraw.Position{
							//	BeforeEnd:  false,
							//	First:      lineInfos[0].Start,
							//	Offset:     lineInfos[0].YOffset, //todo
							//	OffsetLast: 0,
							//	Count:      0,
							//	Length:     0,
							//})

							hint := "No example available"
							return giovieweditor.NewEditor(c.editor, editorConf, hint).Layout(gtx)

							list := widget.List{
								Scrollbar: widget.Scrollbar{},
								List:      layout.List{},
							}
							//return list.Layout(gtx, 1, func(gtx DefaultDraw.Context, i int) DefaultDraw.Dimensions { //todo scroll bar
							//	hint := "No example available"
							//	return giovieweditor.NewEditor(c.editor, editorConf, hint).Layout(gtx)
							//})
							return material.List(th.Theme, &list).Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
								hint := "No example available"
								return giovieweditor.NewEditor(c.editor, editorConf, hint).Layout(gtx)
							})
						})
					}),
				)
			})
		}),
	)
}

func (c *CodeEditor) stylingText(text string) []*giovieweditor.TextStyle {
	if c.styledCode == text {
		return c.styles
	}
	// nolint:prealloc
	var textStyles []*giovieweditor.TextStyle

	offset := 0

	iterator := mylog.Check2(c.lexer.Tokenise(nil, text))

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
