package ux

import (
	"fmt"
	"github.com/ddkwork/golibrary/std/stream"
	"github.com/ddkwork/ux/languages"
	"image"

	"gioui.org/font"
	"gioui.org/gesture"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/oligo/gvcode/color"
	"github.com/oligo/gvcode/textstyle/decoration"
	"github.com/oligo/gvcode/textstyle/syntax"
	"github.com/oligo/gvcode/textview"
)

// RichTextLabel is a widget for laying out and drawing rich text.
type RichTextLabel struct {
	// Face defines the text style.
	Font font.Font
	// Text is the content displayed by the label.
	Text string
	// TextSize determines the size of the text glyphs.
	TextSize unit.Sp
	// Alignment specifies the text alignment.
	Alignment text.Alignment
	// WrapPolicy configures how displayed text will be broken into lines.
	WrapLine bool
	// LineHeight controls the distance between the baselines of lines of text.
	// If zero, a sensible default will be used.
	LineHeight unit.Sp
	// LineHeightScale applies a scaling factor to the LineHeight. If zero, a
	// sensible default will be used.
	LineHeightScale float32

	view    *textview.TextView
	cs      *syntax.ColorScheme
	shaper  *text.Shaper
	clicker gesture.Click
	dragger gesture.Drag
}

func (l *RichTextLabel) SetColorScheme(cs *syntax.ColorScheme) {
	if l.view == nil {
		l.view = textview.NewTextView()
	}
	l.view.SetColorScheme(cs)
	l.cs = cs
}

func (l *RichTextLabel) SetText(text string, textStyles []syntax.Token, decorations []decoration.Decoration) {
	if l.view == nil {
		l.view = textview.NewTextView()
	}

	l.view.SetText(text)
	l.Text = text
	if l.cs != nil {
		l.view.SetSyntaxTokens(textStyles...)
	}
	l.view.ClearDecorations("")
	if len(decorations) > 0 {
		l.view.AddDecorations(decorations...)
	}
}

func (l *RichTextLabel) stylingText(text string) ([]syntax.Token, []decoration.Decoration) {
	//if c.styledCode == text {
	//	return c.tokens
	//}

	// nolint:prealloc
	var tokens []syntax.Token
	var colors []decoration.Decoration

	offset := 0
	tokens_, style := languages.GetTokens(stream.NewBuffer(text), languages.NasmKind)

	//iterator, err := c.lexer.Tokenise(nil, text)
	//if err != nil {
	//	return tokens
	//}

	for _, token := range tokens_ {
		gtoken := syntax.Token{
			Start: offset,
			End:   offset + len([]rune(token.Value)),
			Scope: syntax.StyleScope(fmt.Sprintf("%s", token.Type)),
		}
		tokens = append(tokens, gtoken)
		offset = gtoken.End
		nrgba := chromaColorToNRGBA(style.Get(token.Type).Colour)
		colors = append(colors, decoration.Decoration{
			Source:        nil,
			Priority:      0,
			Start:         offset,
			End:           offset + len([]rune(token.Value)),
			Background:    &decoration.Background{Color: color.MakeColor(nrgba)},
			Underline:     nil,
			Squiggle:      nil,
			Strikethrough: nil,
			Border:        nil,
			Italic:        false,
			Bold:          false,
		})

	}

	//c.styledCode = text
	//c.tokens = tokens

	return tokens, colors
}

// Layout the label with the given shaper, font, size, text, and material, returning metadata about the shaped text.
func (l *RichTextLabel) Layout(gtx layout.Context) layout.Dimensions {
	l.Update(gtx)

	cs := gtx.Constraints
	if l.view == nil {
		l.view = textview.NewTextView()
		l.view.SetText(l.Text)
	}

	l.view.Alignment = l.Alignment
	l.view.SetWrapLine(l.WrapLine)
	l.view.TextSize = l.TextSize
	l.view.LineHeight = l.LineHeight
	l.view.LineHeightScale = l.LineHeightScale
	l.view.TabWidth = 4

	semantic.LabelOp(l.Text).Add(gtx.Ops)
	l.view.Layout(gtx, l.shaper)

	dims := l.view.Dimensions()
	dims.Size = cs.Constrain(dims.Size)

	defer clip.Rect(image.Rectangle{Max: dims.Size}).Push(gtx.Ops).Pop()
	pointer.CursorText.Add(gtx.Ops)
	event.Op(gtx.Ops, l)
	l.clicker.Add(gtx.Ops)
	l.dragger.Add(gtx.Ops)
	if len(l.Text) > 0 {
		l.view.PaintSelection(gtx, l.cs.SelectColor.Op(gtx.Ops))
		l.view.PaintText(gtx, l.cs.Foreground.Op(gtx.Ops))
	}

	return dims
}

func (l *RichTextLabel) Update(gtx layout.Context) {

}

func RichLabel(size unit.Sp, txt string) *RichTextLabel {
	l := &RichTextLabel{
		Font:     font.Font{Typeface: th.Face},
		Text:     txt,
		TextSize: size,
		WrapLine: true,
		shaper:   th.Shaper,
	}

	l.cs = &syntax.ColorScheme{}
	l.cs.Foreground = color.MakeColor(th.Palette.Fg)
	l.cs.Background = color.MakeColor(th.Palette.ContrastBg).MulAlpha(0x60)

	return l
}
