package ux

import (
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/styles"
	gvcolor "github.com/oligo/gvcode/color"
	"github.com/oligo/gvcode/textstyle/syntax"
)

type colorStyle struct {
	scope      syntax.StyleScope
	textStyle  syntax.TextStyle
	color      gvcolor.Color
	background gvcolor.Color
}

// registry holds the color styles for styles
var registry = make(map[string][]colorStyle)

func getColorStyles(name string) []colorStyle {
	if st, ok := registry[name]; ok {
		return st
	}

	style := styles.Get(name)
	if style == nil {
		style = styles.Fallback
	}

	out := make([]colorStyle, 0)
	for _, token := range style.Types() {
		if ok, styleColor := getColorStyle(token, style.Get(token)); ok {
			out = append(out, styleColor)
		} else {
			// If the token type is not recognized, we can skip it
			continue
		}
	}

	registry[name] = out
	return out
}

func getColorStyle(token chroma.TokenType, style chroma.StyleEntry) (bool, colorStyle) {
	styleStr := style.String()
	cc := strings.Split(styleStr, " ")
	if len(cc) < 1 {
		return false, colorStyle{}
	}

	var textStyle syntax.TextStyle = 0
	if strings.EqualFold(cc[0], "bold") {
		textStyle = syntax.Bold
	} else if strings.EqualFold(cc[0], "italic") {
		textStyle = syntax.Italic
	} else if strings.EqualFold(cc[0], "underline") {
		textStyle = syntax.Underline
	}

	var setColor gvcolor.Color
	if style.Colour.IsSet() {
		setColor = gvcolor.MakeColor(chromaColorToNRGBA(style.Colour))
	} else {
		setColor = gvcolor.MakeColor(th.Fg)
	}

	return true, colorStyle{
		scope:      syntax.StyleScope(fmt.Sprintf("%s", token)),
		textStyle:  textStyle,
		color:      setColor,
		background: gvcolor.MakeColor(th.Bg),
	}
}
