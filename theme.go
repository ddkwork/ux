package ux

import (
	"image/color"

	"github.com/ddkwork/ux/widget/material"

	"gioui.org/layout"
	"gioui.org/unit"
)

var th = NewTheme()

//var theme = Theme()
//
//func Theme() *material.Theme {
//	t := material.NewTheme()
//	t.Palette = material.Palette{
//Fg:         White,
//Bg:         BackgroundColor,
//ContrastBg: DeepPurpleA100,
//ContrastFg: Grey900,
//	}
//	t.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
//	return t
//}

type ElementStyle struct {
	TextSize  unit.Sp
	Height    unit.Dp
	Inset     layout.Inset
	IconSize  unit.Dp
	TextColor color.NRGBA
}

type Size struct {
	Tiny   ElementStyle
	Small  ElementStyle
	Medium ElementStyle
	Large  ElementStyle

	DefaultElementWidth      unit.Dp
	DefaultTextSize          unit.Sp
	DropdownTextSize         unit.Sp
	DefaultIconSize          unit.Dp
	DefaultElementRadiusSize unit.Dp
	DefaultWidgetRadiusSize  unit.Dp

	MarkdownPointSize unit.Sp
}
type Color struct {
	DefaultWindowBgGrayColor  color.NRGBA
	DefaultContentBgGrayColor color.NRGBA
	CardBgColor               color.NRGBA
	DefaultBgGrayColor        color.NRGBA
	DefaultTextWhiteColor     color.NRGBA
	DefaultLinkColor          color.NRGBA
	DefaultBorderGrayColor    color.NRGBA
	DefaultBorderBlueColor    color.NRGBA
	DefaultLineColor          color.NRGBA
	DefaultMaskBgColor        color.NRGBA
	DefaultIconColor          color.NRGBA

	TableHeaderBgColor color.NRGBA

	InputInactiveBorderColor  color.NRGBA
	InputActiveBorderColor    color.NRGBA
	InputHoveredBorderColor   color.NRGBA
	InputFocusedBorderColor   color.NRGBA
	InputFocusedBgColor       color.NRGBA
	InputActivatedBorderColor color.NRGBA

	ButtonBorderColor       color.NRGBA
	ButtonDefaultTextColor  color.NRGBA
	ButtonTertiaryBgColor   color.NRGBA
	ButtonTertiaryTextColor color.NRGBA
	ButtonTextBlackColor    color.NRGBA
	ButtonDefaultColor      color.NRGBA
	ButtonTertiaryColor     color.NRGBA

	WhiteColor   color.NRGBA
	GreenColor   color.NRGBA
	ErrorColor   color.NRGBA
	WarningColor color.NRGBA
	SuccessColor color.NRGBA
	BlueColor    color.NRGBA
	InfoColor    color.NRGBA
	PrimaryColor color.NRGBA

	SwitchTabHoverTextColor    color.NRGBA
	SwitchTabSelectedTextColor color.NRGBA
	SwitchTabSelectedLineColor color.NRGBA

	RadioSelectBgColor color.NRGBA

	MenuBarBgColor      color.NRGBA
	MenuBarBorderColor  color.NRGBA
	MenuBarHoveredColor color.NRGBA

	BorderBlueColor          color.NRGBA
	BorderLightGrayColor     color.NRGBA
	HoveredBorderBlueColor   color.NRGBA
	FocusedBorderBlueColor   color.NRGBA
	ActivatedBorderBlueColor color.NRGBA
	FocusedBgColor           color.NRGBA
	TextSelectionColor       color.NRGBA
	HintTextColor            color.NRGBA

	DropDownBorderColor          color.NRGBA
	DropDownSelectedItemBgColor  color.NRGBA
	DropDownHoveredBorderColor   color.NRGBA
	DropDownBgGrayColor          color.NRGBA
	DropDownItemHoveredGrayColor color.NRGBA

	ActionTipsBgGrayColor color.NRGBA
	ProgressBarColor      color.NRGBA

	MenuHoveredBgColor  color.NRGBA
	MenuSelectedBgColor color.NRGBA
	LogTextWhiteColor   color.NRGBA

	NotificationBgColor        color.NRGBA
	NotificationTextWhiteColor color.NRGBA
	ModalBgGrayColor           color.NRGBA

	DropdownMenuBgColor color.NRGBA
	DropdownTextColor   color.NRGBA

	NoticeInfoColor    color.NRGBA
	NoticeSuccessColor color.NRGBA
	NoticeWaringColor  color.NRGBA
	NoticeErrorColor   color.NRGBA

	JsonStartEndColor color.NRGBA
	JsonKeyColor      color.NRGBA
	JsonStringColor   color.NRGBA
	JsonNumberColor   color.NRGBA
	JsonBoolColor     color.NRGBA
	JsonNullColor     color.NRGBA

	MenuItemTextColor         color.NRGBA
	MenuItemTextSelectedColor color.NRGBA

	CloseIconColor color.NRGBA
	MinIconColor   color.NRGBA
	FullIconColor  color.NRGBA

	TreeIconColor      color.NRGBA
	TreeHoveredBgColor color.NRGBA
	TreeClickedBgColor color.NRGBA

	MarkdownMarkColor           color.NRGBA
	MarkdownDefaultColor        color.NRGBA
	MarkdownHeaderColor         color.NRGBA
	MarkdownBlockquoteBgColorL1 color.NRGBA
	MarkdownBlockquoteBgColorL2 color.NRGBA
	MarkdownBlockquoteBgColorL3 color.NRGBA
	MarkdownBlockquoteBgColorL4 color.NRGBA
	MarkdownBlockquoteBgColorL5 color.NRGBA
	MarkdownBlockquoteBgColorL6 color.NRGBA
	MarkdownBlockquoteBgColorL7 color.NRGBA
}

type Theme struct {
	*material.Theme
	Color
	Size
}

func NewTheme(isDark ...bool) *Theme {
	t := &Theme{
		Theme: material.NewTheme(),
	}
	t.Theme.Palette = material.Palette{
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
func (t *Theme) dark() *Theme {
	t.Color.DefaultWindowBgGrayColor = color.NRGBA{R: 32, G: 34, B: 36, A: 255}
	t.Color.DefaultContentBgGrayColor = color.NRGBA{R: 24, G: 24, B: 28, A: 255}

	t.Color.DefaultBgGrayColor = color.NRGBA{R: 53, G: 54, B: 56, A: 255}
	t.Color.DefaultTextWhiteColor = color.NRGBA{R: 223, G: 223, B: 224, A: 255}
	t.Color.DefaultLinkColor = color.NRGBA{R: 107, G: 155, B: 250, A: 255}
	t.Color.DefaultBorderGrayColor = color.NRGBA{R: 53, G: 54, B: 56, A: 255}
	t.Color.DefaultBorderBlueColor = color.NRGBA{R: 127, G: 231, B: 196, A: 255}
	t.Color.DefaultLineColor = color.NRGBA{R: 43, G: 45, B: 49, A: 255}
	t.Color.DefaultMaskBgColor = color.NRGBA{R: 10, G: 10, B: 12, A: 230}
	t.Color.DefaultIconColor = color.NRGBA{R: 136, G: 136, B: 137, A: 255}

	t.Color.TableHeaderBgColor = color.NRGBA{R: 24, G: 24, B: 28, A: 255}

	t.Color.MenuBarBgColor = color.NRGBA{R: 39, G: 39, B: 42, A: 255}
	t.Color.MenuBarBorderColor = color.NRGBA{R: 80, G: 80, B: 81, A: 255}
	t.Color.MenuBarHoveredColor = color.NRGBA{R: 19, G: 87, B: 191, A: 255}

	t.Color.BorderBlueColor = color.NRGBA{R: 127, G: 231, B: 196, A: 255}
	t.Color.BorderLightGrayColor = color.NRGBA{R: 65, G: 65, B: 68, A: 255}
	t.Color.HoveredBorderBlueColor = color.NRGBA{R: 127, G: 231, B: 196, A: 255}
	t.Color.FocusedBorderBlueColor = color.NRGBA{R: 127, G: 231, B: 196, A: 255}
	t.Color.ActivatedBorderBlueColor = color.NRGBA{R: 127, G: 231, B: 196, A: 255}
	t.Color.FocusedBgColor = color.NRGBA{R: 33, G: 50, B: 46, A: 255}
	t.Color.TextSelectionColor = color.NRGBA{R: 92, G: 136, B: 177, A: 255}
	t.Color.HintTextColor = color.NRGBA{R: 136, G: 136, B: 137, A: 255}

	t.Color.InputInactiveBorderColor = color.NRGBA{R: 81, G: 82, B: 89, A: 255}
	t.Color.InputHoveredBorderColor = color.NRGBA{R: 189, G: 189, B: 189, A: 255}
	t.Color.InputActiveBorderColor = color.NRGBA{R: 189, G: 189, B: 189, A: 255}
	t.Color.InputFocusedBorderColor = color.NRGBA{R: 189, G: 189, B: 189, A: 255}
	t.Color.InputFocusedBgColor = color.NRGBA{R: 53, G: 54, B: 56, A: 255}
	t.Color.InputActivatedBorderColor = color.NRGBA{R: 53, G: 54, B: 56, A: 255}

	t.Color.DropDownSelectedItemBgColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	t.Color.DropDownBorderColor = color.NRGBA{R: 81, G: 82, B: 89, A: 255}
	t.Color.DropDownHoveredBorderColor = color.NRGBA{R: 189, G: 189, B: 189, A: 255}
	t.Color.DropDownBgGrayColor = color.NRGBA{R: 72, G: 72, B: 77, A: 255}
	t.Color.DropDownItemHoveredGrayColor = color.NRGBA{R: 90, G: 90, B: 96, A: 255}
	t.Color.DropdownMenuBgColor = color.NRGBA{R: 72, G: 72, B: 77, A: 255}
	t.Color.DropdownTextColor = color.NRGBA{R: 212, G: 212, B: 213, A: 255}

	t.Color.ButtonBorderColor = color.NRGBA{R: 76, G: 76, B: 79, A: 255}
	t.Color.ButtonTertiaryBgColor = color.NRGBA{R: 24, G: 24, B: 28, A: 255}
	t.Color.ButtonTertiaryTextColor = color.NRGBA{R: 149, G: 149, B: 150, A: 255}
	t.Color.ButtonDefaultTextColor = color.NRGBA{R: 216, G: 216, B: 217, A: 255}
	t.Color.ButtonTextBlackColor = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	t.Color.WhiteColor = color.NRGBA{R: 202, G: 202, B: 203, A: 255}
	t.Color.GreenColor = color.NRGBA{R: 101, G: 231, B: 188, A: 255}
	t.Color.ErrorColor = color.NRGBA{R: 232, G: 127, B: 127, A: 255}
	t.Color.WarningColor = color.NRGBA{R: 242, G: 201, B: 126, A: 255}
	t.Color.SuccessColor = color.NRGBA{R: 99, G: 226, B: 184, A: 255}
	t.Color.InfoColor = color.NRGBA{R: 113, G: 192, B: 231, A: 255}
	t.Color.BlueColor = color.NRGBA{R: 68, G: 137, B: 245, A: 255}
	t.Color.PrimaryColor = color.NRGBA{R: 99, G: 226, B: 184, A: 255}
	t.Color.ButtonDefaultColor = color.NRGBA{R: 24, G: 24, B: 28, A: 255}
	t.Color.ButtonTertiaryColor = color.NRGBA{R: 24, G: 24, B: 28, A: 255}

	t.Color.CardBgColor = color.NRGBA{R: 24, G: 24, B: 28, A: 255}

	t.Color.SwitchTabHoverTextColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	t.Color.SwitchTabSelectedTextColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	t.Color.SwitchTabSelectedLineColor = color.NRGBA{R: 70, G: 137, B: 245, A: 255}

	t.Color.RadioSelectBgColor = color.NRGBA{R: 201, G: 204, B: 207, A: 255}

	t.Color.ActionTipsBgGrayColor = color.NRGBA{A: 255, R: 48, G: 48, B: 51}
	t.Color.ProgressBarColor = color.NRGBA{R: 127, G: 200, B: 235, A: 255}

	t.Color.MenuHoveredBgColor = color.NRGBA{R: 45, G: 45, B: 48, A: 255}
	t.Color.MenuSelectedBgColor = color.NRGBA{R: 35, G: 54, B: 51, A: 255}
	t.Color.LogTextWhiteColor = color.NRGBA{R: 202, G: 202, B: 203, A: 255}

	t.Color.NotificationBgColor = color.NRGBA{R: 72, G: 72, B: 77, A: 255}
	t.Color.NotificationTextWhiteColor = color.NRGBA{R: 219, G: 219, B: 220, A: 255}
	t.Color.ModalBgGrayColor = color.NRGBA{R: 44, G: 44, B: 50, A: 255}

	t.Color.NoticeInfoColor = color.NRGBA{R: 108, G: 184, B: 221, A: 255}
	t.Color.NoticeSuccessColor = color.NRGBA{R: 101, G: 231, B: 188, A: 255}
	t.Color.NoticeWaringColor = color.NRGBA{R: 242, G: 201, B: 126, A: 255}
	t.Color.NoticeErrorColor = color.NRGBA{R: 231, G: 127, B: 127, A: 255}

	t.Color.JsonStartEndColor = color.NRGBA{R: 194, G: 196, B: 202, A: 255}
	t.Color.JsonKeyColor = color.NRGBA{R: 159, G: 101, B: 150, A: 255}
	t.Color.JsonStringColor = color.NRGBA{R: 105, G: 168, B: 114, A: 255}
	t.Color.JsonNumberColor = color.NRGBA{R: 41, G: 159, B: 171, A: 255}
	t.Color.JsonBoolColor = color.NRGBA{R: 161, G: 112, B: 88, A: 255}
	t.Color.JsonNullColor = color.NRGBA{R: 170, G: 118, B: 93, A: 255}

	t.Color.MenuItemTextColor = color.NRGBA{R: 150, G: 150, B: 150, A: 255}
	t.Color.MenuItemTextSelectedColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}

	t.Color.CloseIconColor = color.NRGBA{R: 255, G: 95, B: 86, A: 255}
	t.Color.MinIconColor = color.NRGBA{R: 255, G: 188, B: 45, A: 255}
	t.Color.FullIconColor = color.NRGBA{R: 43, G: 200, B: 64, A: 255}

	t.Color.TreeIconColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	t.Color.TreeHoveredBgColor = color.NRGBA{R: 59, G: 60, B: 61, A: 255}
	t.Color.TreeHoveredBgColor = color.NRGBA{R: 84, G: 84, B: 84, A: 255}
	t.Color.TreeClickedBgColor = color.NRGBA{R: 87, G: 87, B: 87, A: 255}

	t.Color.MarkdownMarkColor = color.NRGBA{R: 255, G: 255, B: 0, A: 255}
	t.Color.MarkdownDefaultColor = color.NRGBA{R: 223, G: 223, B: 224, A: 255}
	t.Color.MarkdownHeaderColor = color.NRGBA{R: 102, G: 204, B: 204, A: 255}

	t.Color.MarkdownBlockquoteBgColorL1 = color.NRGBA{R: 48, G: 49, B: 53, A: 255}
	t.Color.MarkdownBlockquoteBgColorL2 = color.NRGBA{R: 64, G: 66, B: 70, A: 255}
	t.Color.MarkdownBlockquoteBgColorL3 = color.NRGBA{R: 78, G: 81, B: 86, A: 255}
	t.Color.MarkdownBlockquoteBgColorL4 = color.NRGBA{R: 91, G: 95, B: 100, A: 255}
	t.Color.MarkdownBlockquoteBgColorL5 = color.NRGBA{R: 103, G: 107, B: 113, A: 255}
	t.Color.MarkdownBlockquoteBgColorL6 = color.NRGBA{R: 113, G: 118, B: 124, A: 255}
	t.Color.MarkdownBlockquoteBgColorL7 = color.NRGBA{R: 122, G: 128, B: 134, A: 255}

	t.Size.Tiny = ElementStyle{TextSize: unit.Sp(9), Height: unit.Dp(10), Inset: layout.UniformInset(unit.Dp(4)), IconSize: unit.Dp(14)}
	t.Size.Small = ElementStyle{TextSize: unit.Sp(12), Height: unit.Dp(15), Inset: layout.UniformInset(unit.Dp(6)), IconSize: unit.Dp(18)}
	t.Size.Medium = ElementStyle{TextSize: unit.Sp(12), Height: unit.Dp(17), Inset: layout.UniformInset(unit.Dp(6)), IconSize: unit.Dp(20)}
	t.Size.Large = ElementStyle{TextSize: unit.Sp(20), Height: unit.Dp(25), Inset: layout.UniformInset(unit.Dp(10)), IconSize: unit.Dp(30)}

	t.Size.DefaultElementWidth = unit.Dp(500)
	t.Size.DefaultTextSize = unit.Sp(12)
	t.Size.DropdownTextSize = unit.Sp(12)
	t.Size.DefaultIconSize = unit.Dp(20)
	t.Size.DefaultElementRadiusSize = unit.Dp(4)
	t.Size.DefaultWidgetRadiusSize = unit.Dp(8)

	t.Size.MarkdownPointSize = unit.Sp(14)

	return t
}
