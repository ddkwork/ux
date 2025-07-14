package ux

import (
	"bytes"
	_ "embed"
	"fmt"
	"gioui.org/font/opentype"
	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/ux/languages"
	"github.com/ddkwork/ux/resources/images"
	"github.com/ddkwork/ux/widget/material"
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/oligo/gvcode"
	gvcolor "github.com/oligo/gvcode/color"
	"github.com/oligo/gvcode/textstyle/syntax"
)

type CodeEditor struct {
	editor *gvcode.Editor
	code   string

	styledCode string
	tokens     []syntax.Token

	lexer chroma.Lexer
	lang  languages.LanguagesKind

	onChange func(text string)

	font font.FontFace

	border widget.Border

	beatufier   widget.Clickable
	loadExample widget.Clickable

	withBeautify bool

	onLoadExample func()

	xScroll widget.Scrollbar
	yScroll widget.Scrollbar

	editorConfig EditorConfig
}
type EditorConfig struct {
	FontFamily        string `yaml:"fontFamily"`
	FontSize          int    `yaml:"FontSize"`
	Indentation       string `yaml:"indentation"` // spaces or tabs
	Theme             string `yaml:"theme"`       // theme name, e.q dracula, github, etc.
	TabWidth          int    `yaml:"tabWidth"`
	AutoCloseBrackets bool   `yaml:"autoCloseBrackets"`
	AutoCloseQuotes   bool   `yaml:"autoCloseQuotes"`
}

func (e EditorConfig) Changed(other EditorConfig) bool {
	return e.FontFamily != other.FontFamily ||
		e.FontSize != other.FontSize ||
		e.Indentation != other.Indentation ||
		e.Theme != other.Theme ||
		e.TabWidth != other.TabWidth ||
		e.AutoCloseBrackets != other.AutoCloseBrackets ||
		e.AutoCloseQuotes != other.AutoCloseQuotes
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

// NewEditor is a helper function to setup a editor with the
// provided theme.
func NewEditor() *gvcode.Editor {
	editor := &gvcode.Editor{}

	colorScheme := syntax.ColorScheme{}
	colorScheme.Foreground = gvcolor.MakeColor(th.Fg)
	colorScheme.Background = gvcolor.MakeColor(th.Bg)
	colorScheme.SelectColor = gvcolor.MakeColor(th.ContrastBg).MulAlpha(0x60)
	colorScheme.LineColor = gvcolor.MakeColor(th.ContrastBg).MulAlpha(0x30)
	colorScheme.LineNumberColor = gvcolor.MakeColor(th.Fg).MulAlpha(0xb6)

	editor.WithOptions(
		gvcode.WrapLine(false),
		gvcode.WithFont(font.Font{Typeface: th.Face}),
		gvcode.WithTextSize(th.TextSize),
		gvcode.WithTextAlignment(text.Start),
		gvcode.WithLineHeight(0, 1.2),
		gvcode.WithTabWidth(4),
		gvcode.WithLineNumberGutterGap(unit.Dp(24)),
		gvcode.WithColorScheme(colorScheme),
	)

	return editor
}

const (
	IndentationSpaces = "spaces"
	IndentationTabs   = "tabs"
)

func NewCodeEditor(code string, lang languages.LanguagesKind) *CodeEditor {
	fff := MustGetCodeEditorFont()

	c := &CodeEditor{
		editor: NewEditor(),
		code:   code,
		font:   fff,
		lang:   lang,
		editorConfig: EditorConfig{
			FontFamily:        "JetBrains Mono",
			FontSize:          12,
			Indentation:       IndentationSpaces,
			Theme:             "dracula",
			TabWidth:          4,
			AutoCloseBrackets: true,
			AutoCloseQuotes:   true,
		},
	}

	c.lexer = getLexer(lang)
	c.setEditorOptions()

	//prefs.AddGlobalConfigChangeListener(func(old, updated domain.GlobalConfig) {
	//	if old.Spec.Editor.Changed(updated.Spec.Editor) {
	//		c.editorConfig = updated.Spec.Editor
	//		c.setEditorOptions()
	//	}
	//})

	c.border = widget.Border{
		Color:        th.BorderLightGrayColor,
		Width:        unit.Dp(1),
		CornerRadius: unit.Dp(4),
	}

	c.editor.SetText(code)
	return c
}

func (c *CodeEditor) setEditorOptions() {
	// color scheme
	colorScheme := syntax.ColorScheme{}
	colorScheme.Foreground = gvcolor.MakeColor(th.Fg)
	colorScheme.SelectColor = gvcolor.MakeColor(th.TextSelectionColor).MulAlpha(0x60)
	colorScheme.LineColor = gvcolor.MakeColor(th.ContrastBg).MulAlpha(0x30)
	colorScheme.LineNumberColor = gvcolor.MakeColor(th.ContrastFg).MulAlpha(0xb6)

	syntaxStyles := getColorStyles("dracula")
	for _, style := range syntaxStyles {
		colorScheme.AddStyle(style.scope, style.textStyle, style.color, style.background)
	}

	editorOptions := []gvcode.EditorOption{
		gvcode.WithFont(c.font.Font),
		gvcode.WithTextSize(unit.Sp(c.editorConfig.FontSize)),
		gvcode.WithTextAlignment(text.Start),
		gvcode.WithLineHeight(unit.Sp(16), 1),
		gvcode.WithTabWidth(c.editorConfig.TabWidth),
		gvcode.WithSoftTab(c.editorConfig.Indentation == IndentationSpaces),
		gvcode.WrapLine(true),
		gvcode.WithLineNumber(true),
		gvcode.WithColorScheme(colorScheme),
	}

	if !c.editorConfig.AutoCloseBrackets {
		editorOptions = append(editorOptions, gvcode.WithQuotePairs(map[rune]rune{}))
	}

	if !c.editorConfig.AutoCloseQuotes {
		editorOptions = append(editorOptions, gvcode.WithBracketPairs(map[rune]rune{}))
	}

	c.editor.WithOptions(editorOptions...)
}

func getLexer(lang languages.LanguagesKind) chroma.Lexer {
	lexer := lexers.Get(lang.String())
	if lexer == nil {
		lexer = lexers.Fallback
	}
	return chroma.Coalesce(lexer)
}

func (c *CodeEditor) WithBeautifier(enabled bool) {
	c.withBeautify = enabled
}

func (c *CodeEditor) SetOnChanged(f func(text string)) {
	c.onChange = f
}

func (c *CodeEditor) SetReadOnly(readOnly bool) {
	c.editor.WithOptions(gvcode.ReadOnlyMode(readOnly))
}

func (c *CodeEditor) SetOnLoadExample(f func()) {
	c.onLoadExample = f
}

func (c *CodeEditor) SetCode(code string) {
	c.editor.SetText(code)
	c.code = code
	c.editor.SetSyntaxTokens(c.stylingText(c.editor.Text())...)
}
func (c *CodeEditor) AppendText(text string) {
	c.editor.Insert(text)
}
func (c *CodeEditor) SetLanguage(lang languages.LanguagesKind) {
	c.lang = lang
	c.lexer = getLexer(lang)
	c.editor.SetSyntaxTokens(c.stylingText(c.editor.Text())...)
}

func (c *CodeEditor) Code() string {
	return c.editor.Text()
}

func (c *CodeEditor) Layout(gtx layout.Context) layout.Dimensions {
	if c.styledCode == "" {
		// First time styling
		c.editor.SetSyntaxTokens(c.stylingText(c.editor.Text())...)
	}

	scrollIndicatorColor := gvcolor.MakeColor(th.Fg).MulAlpha(0x30)

	if !c.editor.ReadOnly() {
		if ev, ok := c.editor.Update(gtx); ok {
			if _, ok := ev.(gvcode.ChangeEvent); ok {
				st := c.stylingText(c.editor.Text())
				c.tokens = st
				c.editor.SetSyntaxTokens(st...)
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

	xScrollDist := c.xScroll.ScrollDistance()
	yScrollDist := c.yScroll.ScrollDistance()
	if xScrollDist != 0.0 || yScrollDist != 0.0 {
		c.editor.Scroll(gtx, xScrollDist, yScrollDist)
	}

	flexH := layout.Flex{Axis: layout.Horizontal}

	if c.withBeautify {
		macro := op.Record(gtx.Ops)
		c.beautyButton(gtx)
		defer op.Defer(gtx.Ops, macro.Stop())
	}

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
					//btn.Color = theme.ButtonTextColor
					btn.Inset = layout.Inset{
						Top: unit.Dp(4), Bottom: unit.Dp(4),
						Left: unit.Dp(4), Right: unit.Dp(4),
					}

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
							dims := c.editor.Layout(gtx, th.Shaper)

							macro := op.Record(gtx.Ops)
							scrollbarDims := func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{
									Left: gtx.Metric.PxToDp(c.editor.GutterWidth()),
								}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									minX, maxX, _, _ := c.editor.ScrollRatio()
									bar := makeScrollbar(&c.xScroll, scrollIndicatorColor.NRGBA())
									return bar.Layout(gtx, layout.Horizontal, minX, maxX)
								})
							}(gtx)

							scrollbarOp := macro.Stop()
							defer op.Offset(image.Point{Y: dims.Size.Y - scrollbarDims.Size.Y}).Push(gtx.Ops).Pop()
							scrollbarOp.Add(gtx.Ops)
							return dims
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						_, _, minY, maxY := c.editor.ScrollRatio()
						bar := makeScrollbar(&c.yScroll, scrollIndicatorColor.NRGBA())
						return bar.Layout(gtx, layout.Vertical, minY, maxY)
					}),
				)
			})
		}),
	)
}

func (c *CodeEditor) beautyButton(gtx layout.Context) layout.Dimensions {
	if c.beatufier.Clicked(gtx) {
		//c.SetCode(BeautifyCode(c.lang, c.code))
		c.SetCode(c.code)
		if c.onChange != nil {
			c.onChange(c.editor.Text())
		}
	}

	return layout.Inset{Bottom: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.SE.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			btn := Button(&c.beatufier, images.ContentClearIcon, "Beautify")
			//btn.Color = theme.ButtonTextColor
			btn.Inset = layout.Inset{
				Top: 4, Bottom: 4,
				Left: 4, Right: 4,
			}
			return btn.Layout(gtx)
		})
	})
}

func makeScrollbar(scroll *widget.Scrollbar, color color.NRGBA) material.ScrollbarStyle {
	bar := material.Scrollbar(th, scroll)
	bar.Indicator.Color = color
	bar.Indicator.CornerRadius = unit.Dp(0)
	bar.Indicator.MinorWidth = unit.Dp(8)
	bar.Track.MajorPadding = unit.Dp(0)
	bar.Track.MinorPadding = unit.Dp(1)
	return bar
}

func (c *CodeEditor) stylingText(text string) []syntax.Token {
	if c.styledCode == text {
		return c.tokens
	}

	// nolint:prealloc
	var tokens []syntax.Token

	offset := 0

	iterator, err := c.lexer.Tokenise(nil, text)
	if err != nil {
		return tokens
	}

	for _, token := range iterator.Tokens() {
		gtoken := syntax.Token{
			Start: offset,
			End:   offset + len([]rune(token.Value)),
			Scope: syntax.StyleScope(fmt.Sprintf("%s", token.Type)),
		}
		tokens = append(tokens, gtoken)
		offset = gtoken.End
	}

	c.styledCode = text
	c.tokens = tokens

	return tokens
}

func chromaColorToNRGBA(textColor chroma.Colour) color.NRGBA {
	return color.NRGBA{
		R: textColor.Red(),
		G: textColor.Green(),
		B: textColor.Blue(),
		A: 0xff,
	}
}
