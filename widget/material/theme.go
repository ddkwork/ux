// SPDX-License-Identifier: Unlicense OR MIT

package material

import (
	"image/color"

	"gioui.org/layout"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/font"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

// Palette contains the minimal set of colors that a widget may need to
// draw itself.
type Palette struct {
	// Bg is the background color atop which content is currently being
	// drawn.
	Bg color.NRGBA

	// Fg is a color suitable for drawing on top of Bg.
	Fg color.NRGBA

	// ContrastBg is a color used to draw attention to active,
	// important, interactive widgets such as buttons.
	ContrastBg color.NRGBA

	// ContrastFg is a color suitable for content drawn on top of
	// ContrastBg.
	ContrastFg color.NRGBA
}

// Theme holds the general theme of an app or window. Different top-level
// windows should have different instances of Theme (with different Shapers;
// see the godoc for [text.Shaper]), though their other fields can be equal.
type Theme struct {
	Shaper *text.Shaper
	Palette
	TextSize unit.Sp
	Icon     struct {
		CheckBoxChecked   *widget.Icon
		CheckBoxUnchecked *widget.Icon
		RadioChecked      *widget.Icon
		RadioUnchecked    *widget.Icon
	}
	// Face selects the default typeface for text.
	Face font.Typeface

	// FingerSize is the minimum touch target size.
	FingerSize unit.Dp
	Color
	Size
}

// NewTheme constructs a theme (and underlying text shaper).
func NewTheme(isDark ...bool) *Theme {
	t := &Theme{Shaper: &text.Shaper{}}
	t.Palette = Palette{
		Fg:         rgb(0x000000),
		Bg:         rgb(0xffffff),
		ContrastBg: rgb(0x3f51b5),
		ContrastFg: rgb(0xffffff),
	}
	t.TextSize = 16

	t.Icon.CheckBoxChecked = mustIcon(widget.NewIcon(icons.ToggleCheckBox))
	t.Icon.CheckBoxUnchecked = mustIcon(widget.NewIcon(icons.ToggleCheckBoxOutlineBlank))
	t.Icon.RadioChecked = mustIcon(widget.NewIcon(icons.ToggleRadioButtonChecked))
	t.Icon.RadioUnchecked = mustIcon(widget.NewIcon(icons.ToggleRadioButtonUnchecked))

	// 38dp is on the lower end of possible finger size.
	t.FingerSize = 38

	t.Palette = Palette{
		Fg:         White,
		Bg:         BackgroundColor,
		ContrastBg: DeepPurpleA100,
		ContrastFg: Grey900,
	}
	if len(isDark) == 0 || !isDark[0] {
		t = t.dark()
	}
	t.TextSize = t.Size.DefaultTextSize

	return t
}

var ( // todo 如果场景多的话就来一个 colors的包，抓包什么的估计需要更多的自定义颜色，又要用户可更改
	// Black           = color.NRGBA{A: 0xff}                            // rgb(0, 0, 0)
	White = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff} // rgb(255, 255, 255)
	// DividerFg       = color.NRGBA{R: 88, G: 88, B: 88, A: 255}
	BackgroundColor = color.NRGBA{R: 54, G: 54, B: 54, A: 255}
	// ColorPink       = color.NRGBA{R: 166, G: 96, B: 192, A: 255}
	// ColorHeaderFg   = color.NRGBA{R: 70, G: 70, B: 70, A: 255}
	// ColorHeaderFg = color.NRGBA{R: 76, G: 76, B: 76, A: 255}
	// DeepPurple900   = color.NRGBA{R: 0x31, G: 0x1b, B: 0x92, A: 0xff} // rgb(49, 27, 146)
	DeepPurpleA100 = color.NRGBA{R: 0xb3, G: 0x88, B: 0xff, A: 0xff} // rgb(179, 136, 255)
	// DeepPurpleA200  = color.NRGBA{R: 0x7c, G: 0x4d, B: 0xff, A: 0xff} // rgb(124, 77, 255)
	// DeepPurpleA400  = color.NRGBA{R: 0x65, G: 0x1f, B: 0xff, A: 0xff} // rgb(101, 31, 255)
	// DeepPurpleA700  = color.NRGBA{R: 0x62, B: 0xea, A: 0xff}          // rgb(98, 0, 234)
	// Grey50          = color.NRGBA{R: 0xfa, G: 0xfa, B: 0xfa, A: 0xff} // rgb(250, 250, 250)
	// Grey100         = color.NRGBA{R: 0xf5, G: 0xf5, B: 0xf5, A: 0xff} // rgb(245, 245, 245)
	// Grey200         = color.NRGBA{R: 0xee, G: 0xee, B: 0xee, A: 0xff} // rgb(238, 238, 238)
	// Grey300         = color.NRGBA{R: 0xe0, G: 0xe0, B: 0xe0, A: 0xff} // rgb(224, 224, 224)
	// Grey400         = color.NRGBA{R: 0xbd, G: 0xbd, B: 0xbd, A: 0xff} // rgb(189, 189, 189)
	// Grey500         = color.NRGBA{R: 0x9e, G: 0x9e, B: 0x9e, A: 0xff} // rgb(158, 158, 158)
	// Grey600         = color.NRGBA{R: 0x75, G: 0x75, B: 0x75, A: 0xff} // rgb(117, 117, 117)
	// Grey700         = color.NRGBA{R: 0x61, G: 0x61, B: 0x61, A: 0xff} // rgb(97, 97, 97)
	// Grey800         = color.NRGBA{R: 0x42, G: 0x42, B: 0x42, A: 0xff} // rgb(66, 66, 66)
	Grey900 = color.NRGBA{R: 0x21, G: 0x21, B: 0x21, A: 0xff} // rgb(33, 33, 33)

)

func (t Theme) WithPalette(p Palette) Theme {
	t.Palette = p
	return t
}

func mustIcon(ic *widget.Icon, err error) *widget.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}

func rgb(c uint32) color.NRGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

// ////////////////////////////////////////////////////
func (t *Theme) dark() *Theme {
	t.Color = Color{
		DefaultWindowBgGrayColor:     color.NRGBA{R: 32, G: 34, B: 36, A: 255},
		DefaultContentBgGrayColor:    color.NRGBA{R: 24, G: 24, B: 28, A: 255},
		CardBgColor:                  color.NRGBA{R: 24, G: 24, B: 28, A: 255},
		DefaultBgGrayColor:           color.NRGBA{R: 53, G: 54, B: 56, A: 255},
		DefaultTextWhiteColor:        color.NRGBA{R: 223, G: 223, B: 224, A: 255},
		DefaultLinkColor:             color.NRGBA{R: 107, G: 155, B: 250, A: 255},
		DefaultBorderGrayColor:       color.NRGBA{R: 53, G: 54, B: 56, A: 255},
		DefaultBorderBlueColor:       color.NRGBA{R: 127, G: 231, B: 196, A: 255},
		DefaultLineColor:             color.NRGBA{R: 43, G: 45, B: 49, A: 255},
		DefaultMaskBgColor:           color.NRGBA{R: 10, G: 10, B: 12, A: 230},
		DefaultIconColor:             color.NRGBA{R: 136, G: 136, B: 137, A: 255},
		TableHeaderBgColor:           color.NRGBA{R: 24, G: 24, B: 28, A: 255},
		InputInactiveBorderColor:     color.NRGBA{R: 81, G: 82, B: 89, A: 255},
		InputActiveBorderColor:       color.NRGBA{R: 189, G: 189, B: 189, A: 255},
		InputHoveredBorderColor:      color.NRGBA{R: 189, G: 189, B: 189, A: 255},
		InputFocusedBorderColor:      color.NRGBA{R: 189, G: 189, B: 189, A: 255},
		InputFocusedBgColor:          color.NRGBA{R: 53, G: 54, B: 56, A: 255},
		InputActivatedBorderColor:    color.NRGBA{R: 53, G: 54, B: 56, A: 255},
		ButtonBorderColor:            color.NRGBA{R: 76, G: 76, B: 79, A: 255},
		ButtonDefaultTextColor:       color.NRGBA{R: 216, G: 216, B: 217, A: 255},
		ButtonTertiaryBgColor:        color.NRGBA{R: 24, G: 24, B: 28, A: 255},
		ButtonTertiaryTextColor:      color.NRGBA{R: 149, G: 149, B: 150, A: 255},
		ButtonTextBlackColor:         color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		ButtonDefaultColor:           color.NRGBA{R: 24, G: 24, B: 28, A: 255},
		ButtonTertiaryColor:          color.NRGBA{R: 24, G: 24, B: 28, A: 255},
		WhiteColor:                   color.NRGBA{R: 202, G: 202, B: 203, A: 255},
		GreenColor:                   color.NRGBA{R: 101, G: 231, B: 188, A: 255},
		ErrorColor:                   color.NRGBA{R: 232, G: 127, B: 127, A: 255},
		WarningColor:                 color.NRGBA{R: 242, G: 201, B: 126, A: 255},
		SuccessColor:                 color.NRGBA{R: 99, G: 226, B: 184, A: 255},
		BlueColor:                    color.NRGBA{R: 68, G: 137, B: 245, A: 255},
		InfoColor:                    color.NRGBA{R: 113, G: 192, B: 231, A: 255},
		PrimaryColor:                 color.NRGBA{R: 99, G: 226, B: 184, A: 255},
		SwitchTabHoverTextColor:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		SwitchTabSelectedTextColor:   color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		SwitchTabSelectedLineColor:   color.NRGBA{R: 70, G: 137, B: 245, A: 255},
		RadioSelectBgColor:           color.NRGBA{R: 201, G: 204, B: 207, A: 255},
		MenuBarBgColor:               color.NRGBA{R: 39, G: 39, B: 42, A: 255},
		MenuBarBorderColor:           color.NRGBA{R: 80, G: 80, B: 81, A: 255},
		MenuBarHoveredColor:          color.NRGBA{R: 19, G: 87, B: 191, A: 255},
		BorderBlueColor:              color.NRGBA{R: 127, G: 231, B: 196, A: 255},
		BorderLightGrayColor:         color.NRGBA{R: 65, G: 65, B: 68, A: 255},
		HoveredBorderBlueColor:       color.NRGBA{R: 127, G: 231, B: 196, A: 255},
		FocusedBorderBlueColor:       color.NRGBA{R: 127, G: 231, B: 196, A: 255},
		ActivatedBorderBlueColor:     color.NRGBA{R: 127, G: 231, B: 196, A: 255},
		FocusedBgColor:               color.NRGBA{R: 33, G: 50, B: 46, A: 255},
		TextSelectionColor:           color.NRGBA{R: 92, G: 136, B: 177, A: 255},
		HintTextColor:                color.NRGBA{R: 136, G: 136, B: 137, A: 255},
		DropDownBorderColor:          color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		DropDownSelectedItemBgColor:  color.NRGBA{R: 81, G: 82, B: 89, A: 255},
		DropDownHoveredBorderColor:   color.NRGBA{R: 189, G: 189, B: 189, A: 255},
		DropDownBgGrayColor:          color.NRGBA{R: 72, G: 72, B: 77, A: 255},
		DropDownItemHoveredGrayColor: color.NRGBA{R: 90, G: 90, B: 96, A: 255},
		ActionTipsBgGrayColor:        color.NRGBA{A: 255, R: 48, G: 48, B: 51},
		ProgressBarColor:             color.NRGBA{R: 127, G: 200, B: 235, A: 255},
		MenuHoveredBgColor:           color.NRGBA{R: 45, G: 45, B: 48, A: 255},
		MenuSelectedBgColor:          color.NRGBA{R: 35, G: 54, B: 51, A: 255},
		LogTextWhiteColor:            color.NRGBA{R: 202, G: 202, B: 203, A: 255},
		NotificationBgColor:          color.NRGBA{R: 72, G: 72, B: 77, A: 255},
		NotificationTextWhiteColor:   color.NRGBA{R: 219, G: 219, B: 220, A: 255},
		ModalBgGrayColor:             color.NRGBA{R: 44, G: 44, B: 50, A: 255},
		DropdownMenuBgColor:          color.NRGBA{}, // todo
		DropdownTextColor:            color.NRGBA{}, // todo
		NoticeInfoColor:              color.NRGBA{R: 108, G: 184, B: 221, A: 255},
		NoticeSuccessColor:           color.NRGBA{R: 101, G: 231, B: 188, A: 255},
		NoticeWaringColor:            color.NRGBA{R: 242, G: 201, B: 126, A: 255},
		NoticeErrorColor:             color.NRGBA{R: 231, G: 127, B: 127, A: 255},
		JsonStartEndColor:            color.NRGBA{R: 194, G: 196, B: 202, A: 255},
		JsonKeyColor:                 color.NRGBA{R: 159, G: 101, B: 150, A: 255},
		JsonStringColor:              color.NRGBA{R: 105, G: 168, B: 114, A: 255},
		JsonNumberColor:              color.NRGBA{R: 41, G: 159, B: 171, A: 255},
		JsonBoolColor:                color.NRGBA{R: 161, G: 112, B: 88, A: 255},
		JsonNullColor:                color.NRGBA{R: 170, G: 118, B: 93, A: 255},
		MenuItemTextColor:            color.NRGBA{R: 150, G: 150, B: 150, A: 255},
		MenuItemTextSelectedColor:    color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		CloseIconColor:               color.NRGBA{R: 255, G: 95, B: 86, A: 255},
		MinIconColor:                 color.NRGBA{R: 255, G: 188, B: 45, A: 255},
		FullIconColor:                color.NRGBA{R: 43, G: 200, B: 64, A: 255},
		TreeIconColor:                color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		// TreeHoveredBgColor:          color.NRGBA{R: 59, G: 60, B: 61, A: 255},
		TreeHoveredBgColor:          color.NRGBA{R: 84, G: 84, B: 84, A: 255},
		TreeClickedBgColor:          color.NRGBA{R: 87, G: 87, B: 87, A: 255},
		MarkdownMarkColor:           color.NRGBA{R: 255, G: 255, B: 0, A: 255},
		MarkdownDefaultColor:        color.NRGBA{R: 223, G: 223, B: 224, A: 255},
		MarkdownHeaderColor:         color.NRGBA{R: 102, G: 204, B: 204, A: 255},
		MarkdownBlockquoteBgColorL1: color.NRGBA{R: 48, G: 49, B: 53, A: 255},
		MarkdownBlockquoteBgColorL2: color.NRGBA{R: 64, G: 66, B: 70, A: 255},
		MarkdownBlockquoteBgColorL3: color.NRGBA{R: 78, G: 81, B: 86, A: 255},
		MarkdownBlockquoteBgColorL4: color.NRGBA{R: 91, G: 95, B: 100, A: 255},
		MarkdownBlockquoteBgColorL5: color.NRGBA{R: 103, G: 107, B: 113, A: 255},
		MarkdownBlockquoteBgColorL6: color.NRGBA{R: 113, G: 118, B: 124, A: 255},
		MarkdownBlockquoteBgColorL7: color.NRGBA{R: 122, G: 128, B: 134, A: 255},
	}
	t.Size = Size{
		Tiny:                     ElementStyle{TextSize: unit.Sp(9), Height: unit.Dp(10), Inset: layout.UniformInset(unit.Dp(4)), IconSize: unit.Dp(14)},
		Small:                    ElementStyle{TextSize: unit.Sp(12), Height: unit.Dp(15), Inset: layout.UniformInset(unit.Dp(6)), IconSize: unit.Dp(18)},
		Medium:                   ElementStyle{TextSize: unit.Sp(12), Height: unit.Dp(17), Inset: layout.UniformInset(unit.Dp(6)), IconSize: unit.Dp(20)},
		Large:                    ElementStyle{TextSize: unit.Sp(20), Height: unit.Dp(25), Inset: layout.UniformInset(unit.Dp(10)), IconSize: unit.Dp(30)},
		DefaultElementWidth:      500,
		DefaultTextSize:          12,
		DropdownTextSize:         12,
		DefaultIconSize:          20,
		DefaultElementRadiusSize: 4,
		DefaultWidgetRadiusSize:  8,
		MarkdownPointSize:        14,
	}
	return t
}

type ElementStyle struct {
	TextSize  unit.Sp
	Height    unit.Dp
	Inset     layout.Inset
	IconSize  unit.Dp
	TextColor color.NRGBA
}

type Size struct {
	Tiny                     ElementStyle
	Small                    ElementStyle
	Medium                   ElementStyle
	Large                    ElementStyle
	DefaultElementWidth      unit.Dp
	DefaultTextSize          unit.Sp
	DropdownTextSize         unit.Sp
	DefaultIconSize          unit.Dp
	DefaultElementRadiusSize unit.Dp
	DefaultWidgetRadiusSize  unit.Dp
	MarkdownPointSize        unit.Sp
}
type Color struct { // todo更细的分类和嵌套，方便维护
	DefaultWindowBgGrayColor     color.NRGBA
	DefaultContentBgGrayColor    color.NRGBA
	CardBgColor                  color.NRGBA
	DefaultBgGrayColor           color.NRGBA
	DefaultTextWhiteColor        color.NRGBA
	DefaultLinkColor             color.NRGBA
	DefaultBorderGrayColor       color.NRGBA
	DefaultBorderBlueColor       color.NRGBA
	DefaultLineColor             color.NRGBA
	DefaultMaskBgColor           color.NRGBA
	DefaultIconColor             color.NRGBA
	TableHeaderBgColor           color.NRGBA
	InputInactiveBorderColor     color.NRGBA
	InputActiveBorderColor       color.NRGBA
	InputHoveredBorderColor      color.NRGBA
	InputFocusedBorderColor      color.NRGBA
	InputFocusedBgColor          color.NRGBA
	InputActivatedBorderColor    color.NRGBA
	ButtonBorderColor            color.NRGBA
	ButtonDefaultTextColor       color.NRGBA
	ButtonTertiaryBgColor        color.NRGBA
	ButtonTertiaryTextColor      color.NRGBA
	ButtonTextBlackColor         color.NRGBA
	ButtonDefaultColor           color.NRGBA
	ButtonTertiaryColor          color.NRGBA
	WhiteColor                   color.NRGBA
	GreenColor                   color.NRGBA
	ErrorColor                   color.NRGBA
	WarningColor                 color.NRGBA
	SuccessColor                 color.NRGBA
	BlueColor                    color.NRGBA
	InfoColor                    color.NRGBA
	PrimaryColor                 color.NRGBA
	SwitchTabHoverTextColor      color.NRGBA
	SwitchTabSelectedTextColor   color.NRGBA
	SwitchTabSelectedLineColor   color.NRGBA
	RadioSelectBgColor           color.NRGBA
	MenuBarBgColor               color.NRGBA
	MenuBarBorderColor           color.NRGBA
	MenuBarHoveredColor          color.NRGBA
	BorderBlueColor              color.NRGBA
	BorderLightGrayColor         color.NRGBA
	HoveredBorderBlueColor       color.NRGBA
	FocusedBorderBlueColor       color.NRGBA
	ActivatedBorderBlueColor     color.NRGBA
	FocusedBgColor               color.NRGBA
	TextSelectionColor           color.NRGBA
	HintTextColor                color.NRGBA
	DropDownBorderColor          color.NRGBA
	DropDownSelectedItemBgColor  color.NRGBA
	DropDownHoveredBorderColor   color.NRGBA
	DropDownBgGrayColor          color.NRGBA
	DropDownItemHoveredGrayColor color.NRGBA
	ActionTipsBgGrayColor        color.NRGBA
	ProgressBarColor             color.NRGBA
	MenuHoveredBgColor           color.NRGBA
	MenuSelectedBgColor          color.NRGBA
	LogTextWhiteColor            color.NRGBA
	NotificationBgColor          color.NRGBA
	NotificationTextWhiteColor   color.NRGBA
	ModalBgGrayColor             color.NRGBA
	DropdownMenuBgColor          color.NRGBA
	DropdownTextColor            color.NRGBA
	NoticeInfoColor              color.NRGBA
	NoticeSuccessColor           color.NRGBA
	NoticeWaringColor            color.NRGBA
	NoticeErrorColor             color.NRGBA
	JsonStartEndColor            color.NRGBA
	JsonKeyColor                 color.NRGBA
	JsonStringColor              color.NRGBA
	JsonNumberColor              color.NRGBA
	JsonBoolColor                color.NRGBA
	JsonNullColor                color.NRGBA
	MenuItemTextColor            color.NRGBA
	MenuItemTextSelectedColor    color.NRGBA
	CloseIconColor               color.NRGBA
	MinIconColor                 color.NRGBA
	FullIconColor                color.NRGBA
	TreeIconColor                color.NRGBA
	TreeHoveredBgColor           color.NRGBA
	TreeClickedBgColor           color.NRGBA
	MarkdownMarkColor            color.NRGBA
	MarkdownDefaultColor         color.NRGBA
	MarkdownHeaderColor          color.NRGBA
	MarkdownBlockquoteBgColorL1  color.NRGBA
	MarkdownBlockquoteBgColorL2  color.NRGBA
	MarkdownBlockquoteBgColorL3  color.NRGBA
	MarkdownBlockquoteBgColorL4  color.NRGBA
	MarkdownBlockquoteBgColorL5  color.NRGBA
	MarkdownBlockquoteBgColorL6  color.NRGBA
	MarkdownBlockquoteBgColorL7  color.NRGBA
}

//	func (t *Theme ) darkNaive() *Theme {
//		t.Color.DefaultWindowBgGrayColor = color.NRGBA{R: 17, G: 15, B: 20, A: 255}
//		t.Color.DefaultContentBgGrayColor = color.NRGBA{R: 24, G: 24, B: 28, A: 255}
//
//		t.Color.DefaultBgGrayColor = color.NRGBA{R: 53, G: 54, B: 56, A: 255}
//		t.Color.DefaultTextWhiteColor = color.NRGBA{R: 223, G: 223, B: 224, A: 255}
//		t.Color.DefaultBorderGrayColor = color.NRGBA{R: 53, G: 54, B: 56, A: 255}
//		t.Color.DefaultBorderBlueColor = color.NRGBA{R: 127, G: 231, B: 196, A: 255}
//
//		t.Color.DefaultLineColor = color.NRGBA{R: 44, G: 44, B: 47, A: 255}
//		t.Color.DefaultMaskBgColor = color.NRGBA{R: 10, G: 10, B: 12, A: 230}
//
//		t.Color.DefaultIconColor = color.NRGBA{R: 136, G: 136, B: 137, A: 255}
//		t.Color.BorderBlueColor = color.NRGBA{R: 127, G: 231, B: 196, A: 255}
//		t.Color.BorderLightGrayColor = color.NRGBA{R: 65, G: 65, B: 68, A: 255}
//		t.Color.HoveredBorderBlueColor = color.NRGBA{R: 127, G: 231, B: 196, A: 255}
//		t.Color.FocusedBorderBlueColor = color.NRGBA{R: 127, G: 231, B: 196, A: 255}
//		t.Color.ActivatedBorderBlueColor = color.NRGBA{R: 127, G: 231, B: 196, A: 255}
//		t.Color.FocusedBgColor = color.NRGBA{R: 33, G: 50, B: 46, A: 255}
//		t.Color.TextSelectionColor = color.NRGBA{R: 92, G: 136, B: 177, A: 255}
//		t.Color.HintTextColor = color.NRGBA{R: 136, G: 136, B: 137, A: 255}
//
//		t.Color.DropDownHoveredBorderColor = color.NRGBA{R: 127, G: 231, B: 196, A: 255}
//		t.Color.DropDownBgGrayColor = color.NRGBA{R: 72, G: 72, B: 77, A: 255}
//		t.Color.DropDownItemHoveredGrayColor = color.NRGBA{R: 90, G: 90, B: 96, A: 255}
//
//		t.Color.InputInactiveBorderColor = color.NRGBA{R: 53, G: 54, B: 56, A: 255}
//		t.Color.InputHoveredBorderColor = color.NRGBA{R: 127, G: 231, B: 196, A: 255}
//		t.Color.InputActiveBorderColor = color.NRGBA{R: 53, G: 54, B: 56, A: 255}
//		t.Color.InputFocusedBorderColor = color.NRGBA{R: 33, G: 50, B: 46, A: 255}
//		t.Color.InputActivatedBorderColor = color.NRGBA{R: 53, G: 54, B: 56, A: 255}
//
//		t.Color.ButtonBorderColor = color.NRGBA{R: 76, G: 76, B: 79, A: 255}
//		t.Color.ButtonTertiaryBgColor = color.NRGBA{R: 24, G: 24, B: 28, A: 255}
//		t.Color.ButtonTertiaryTextColor = color.NRGBA{R: 149, G: 149, B: 150, A: 255}
//		t.Color.ButtonDefaultTextColor = color.NRGBA{R: 216, G: 216, B: 217, A: 255}
//		t.Color.ButtonTextBlackColor = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
//		t.Color.WhiteColor = color.NRGBA{R: 202, G: 202, B: 203, A: 255}
//		t.Color.GreenColor = color.NRGBA{R: 101, G: 231, B: 188, A: 255}
//		t.Color.ErrorColor = color.NRGBA{R: 232, G: 127, B: 127, A: 255}
//		t.Color.WarningColor = color.NRGBA{R: 242, G: 201, B: 126, A: 255}
//		t.Color.SuccessColor = color.NRGBA{R: 99, G: 226, B: 184, A: 255}
//		t.Color.InfoColor = color.NRGBA{R: 113, G: 192, B: 231, A: 255}
//		t.Color.PrimaryColor = color.NRGBA{R: 99, G: 226, B: 184, A: 255}
//		t.Color.ButtonDefaultColor = color.NRGBA{R: 24, G: 24, B: 28, A: 255}
//		t.Color.ButtonTertiaryColor = color.NRGBA{R: 24, G: 24, B: 28, A: 255}
//
//
//		t.Color.ActionTipsBgGrayColor = color.NRGBA{A: 255, R: 48, G: 48, B: 51}
//		t.Color.ProgressBarColor = color.NRGBA{R: 127, G: 200, B: 235, A: 255}
//
//		t.Color.MenuHoveredBgColor = color.NRGBA{R: 45, G: 45, B: 48, A: 255}
//		t.Color.MenuSelectedBgColor = color.NRGBA{R: 35, G: 54, B: 51, A: 255}
//		t.Color.LogTextWhiteColor = color.NRGBA{R: 202, G: 202, B: 203, A: 255}
//
//		t.Color.NotificationBgColor = color.NRGBA{R: 72, G: 72, B: 77, A: 255}
//		t.Color.NotificationTextWhiteColor = color.NRGBA{R: 219, G: 219, B: 220, A: 255}
//		t.Color.ModalBgGrayColor = color.NRGBA{R: 44, G: 44, B: 50, A: 255}
//
//		t.Color.DropdownMenuBgColor = color.NRGBA{R: 72, G: 72, B: 77, A: 255}
//		t.Color.DropdownTextColor = color.NRGBA{R: 212, G: 212, B: 213, A: 255}
//
//		t.Color.NoticeInfoColor = color.NRGBA{R: 108, G: 184, B: 221, A: 255}
//		t.Color.NoticeSuccessColor = color.NRGBA{R: 101, G: 231, B: 188, A: 255}
//		t.Color.NoticeWaringColor = color.NRGBA{R: 242, G: 201, B: 126, A: 255}
//		t.Color.NoticeErrorColor = color.NRGBA{R: 231, G: 127, B: 127, A: 255}
//
//		t.Color.JsonStartEndColor = color.NRGBA{R: 194, G: 196, B: 202, A: 255}
//		t.Color.JsonKeyColor = color.NRGBA{R: 159, G: 101, B: 150, A: 255}
//		t.Color.JsonStringColor = color.NRGBA{R: 105, G: 168, B: 114, A: 255}
//		t.Color.JsonNumberColor = color.NRGBA{R: 41, G: 159, B: 171, A: 255}
//		t.Color.JsonBoolColor = color.NRGBA{R: 161, G: 112, B: 88, A: 255}
//		t.Color.JsonNullColor = color.NRGBA{R: 170, G: 118, B: 93, A: 255}
//
//		t.Size.Tiny = ElementStyle{TextSize: unit.Sp(9), Height: unit.Dp(10), Inset: layout.UniformInset(unit.Dp(4)), IconSize: unit.Dp(14)}
//		t.Size.Small = ElementStyle{TextSize: unit.Sp(12), Height: unit.Dp(15), Inset: layout.UniformInset(unit.Dp(6)), IconSize: unit.Dp(18)}
//		t.Size.Medium = ElementStyle{TextSize: unit.Sp(14), Height: unit.Dp(20), Inset: layout.UniformInset(unit.Dp(8)), IconSize: unit.Dp(24)}
//		t.Size.Large = ElementStyle{TextSize: unit.Sp(20), Height: unit.Dp(25), Inset: layout.UniformInset(unit.Dp(10)), IconSize: unit.Dp(30)}
//
//		t.Size.DefaultElementWidth = unit.Dp(500)
//		t.Size.DefaultTextSize = unit.Sp(14)
//		t.Size.DropdownTextSize = unit.Sp(13)
//		t.Size.DefaultIconSize = unit.Dp(20)
//		t.Size.DefaultElementRadiusSize = unit.Dp(4)
//		t.Size.DefaultWidgetRadiusSize = unit.Dp(8)
//
//		return t
//	}
