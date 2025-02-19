package ux

import (
	_ "embed"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"sync"

	"gioui.org/layout"
	"gioui.org/unit"

	"gioui.org/widget"
	"github.com/ddkwork/golibrary/mylog"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	IconSearch               = mylog.Check2(widget.NewIcon(icons.ActionSearch))
	IconDelete               = mylog.Check2(widget.NewIcon(icons.ActionDelete))
	IconExpand               = mylog.Check2(widget.NewIcon(icons.NavigationExpandMore))
	IconSave                 = mylog.Check2(widget.NewIcon(icons.ContentSave))
	IconMenu                 = mylog.Check2(widget.NewIcon(icons.NavigationMenu))
	IconCopy                 = mylog.Check2(widget.NewIcon(icons.ContentContentCopy))
	IconHome                 = mylog.Check2(widget.NewIcon(icons.ActionHome))
	IconHeart                = mylog.Check2(widget.NewIcon(icons.ActionFavorite))
	IconPlus                 = mylog.Check2(widget.NewIcon(icons.ContentAdd))
	IconEdit                 = mylog.Check2(widget.NewIcon(icons.ContentCreate))
	IconVisibility           = mylog.Check2(widget.NewIcon(icons.ActionVisibility))
	IconClose                = mylog.Check2(widget.NewIcon(icons.NavigationClose))
	IconArrowDropDown        = mylog.Check2(widget.NewIcon(icons.NavigationArrowDropDown))
	IconNaviLeft             = mylog.Check2(widget.NewIcon(icons.NavigationChevronLeft))
	IconNaviRight            = mylog.Check2(widget.NewIcon(icons.NavigationChevronRight))
	IconFileFolder           = mylog.Check2(widget.NewIcon(icons.FileFolder))
	IconUpload               = mylog.Check2(widget.NewIcon(icons.FileFileUpload))
	IconDownload             = mylog.Check2(widget.NewIcon(icons.FileFileDownload))
	IconRefresh              = mylog.Check2(widget.NewIcon(icons.NavigationRefresh))
	IconClean                = mylog.Check2(widget.NewIcon(icons.EditorFormatColorText))
	IconFileFolderOpen       = mylog.Check2(widget.NewIcon(icons.FileFolderOpen))
	IconFavorite             = mylog.Check2(widget.NewIcon(icons.ActionFavorite))
	IconAdd                  = mylog.Check2(widget.NewIcon(icons.ContentAdd))
	IconRemove               = mylog.Check2(widget.NewIcon(icons.ContentRemove))
	IconSettings             = mylog.Check2(widget.NewIcon(icons.ActionSettings))
	IconDone                 = mylog.Check2(widget.NewIcon(icons.ActionDone))
	IconForward              = mylog.Check2(widget.NewIcon(icons.NavigationChevronRight))
	IconDeleteForever        = mylog.Check2(widget.NewIcon(icons.ActionDeleteForever))
	IconStart                = mylog.Check2(widget.NewIcon(icons.AVPlayArrow))
	IconStop                 = mylog.Check2(widget.NewIcon(icons.AVStop))
	IconBack                 = mylog.Check2(widget.NewIcon(icons.NavigationArrowBack))
	IconVisibilityOff        = mylog.Check2(widget.NewIcon(icons.ActionVisibilityOff))
	IconNavArrowForward      = mylog.Check2(widget.NewIcon(icons.NavigationArrowForward))
	IconCircle               = mylog.Check2(widget.NewIcon(icons.ImageLens))
	IconNavRight             = mylog.Check2(widget.NewIcon(icons.NavigationChevronRight))
	IconNavExpandLess        = mylog.Check2(widget.NewIcon(icons.NavigationExpandLess))
	IconNavExpandMore        = mylog.Check2(widget.NewIcon(icons.NavigationExpandMore))
	IconActionCode           = mylog.Check2(widget.NewIcon(icons.ActionCode))
	IconActionUpdate         = mylog.Check2(widget.NewIcon(icons.ActionUpdate))
	IconActionHourGlassEmpty = mylog.Check2(widget.NewIcon(icons.ActionHourglassEmpty))
	IconRadioButtonUnchecked = mylog.Check2(widget.NewIcon(icons.ToggleRadioButtonUnchecked))
	IconRadioButtonChecked   = mylog.Check2(widget.NewIcon(icons.ToggleRadioButtonChecked))
	IconCheckCircle          = mylog.Check2(widget.NewIcon(icons.ActionCheckCircle))
	IconError                = mylog.Check2(widget.NewIcon(icons.AlertError))

	SwapHoriz      = mylog.Check2(widget.NewIcon(icons.ActionSwapHoriz))
	MenuIcon       = mylog.Check2(widget.NewIcon(icons.NavigationMenu))
	WorkspacesIcon = mylog.Check2(widget.NewIcon(icons.NavigationApps))
	FileFolderIcon = mylog.Check2(widget.NewIcon(icons.FileFolder))
	TunnelIcon     = mylog.Check2(widget.NewIcon(icons.ActionSwapVerticalCircle))
	ConsoleIcon    = mylog.Check2(widget.NewIcon(icons.HardwareDesktopMac))
	LogsIcon       = mylog.Check2(widget.NewIcon(icons.ActionSubject))
	SettingsIcon   = mylog.Check2(widget.NewIcon(icons.ActionSettings))
	DarkIcon       = mylog.Check2(widget.NewIcon(icons.ActionSubject))
)

//go:embed resources/images/CircledChevronDown.svg
var CircledChevronDown string

//go:embed resources/images/CircledChevronRight.svg
var CircledChevronRight string

type Icon struct {
	*widget.Icon
	Color color.NRGBA
	Size  unit.Dp
}

func (i Icon) Layout(gtx C) D {
	if i.Size <= 0 {
		i.Size = unit.Dp(18)
	}
	if i.Color == (color.NRGBA{}) {
		i.Color = WithAlpha(th.Palette.Fg, 0xb6)
	}

	iconSize := gtx.Dp(i.Size)
	gtx.Constraints = layout.Exact(image.Pt(iconSize, iconSize))

	return i.Icon.Layout(gtx, i.Color)
}

// MulAlpha applies the alpha to the color.
func MulAlpha(c color.NRGBA, alpha uint8) color.NRGBA {
	c.A = uint8(uint32(c.A) * uint32(alpha) / 0xFF)
	return c
}

// Disabled blends color towards the luminance and multiplies alpha.
// Blending towards luminance will desaturate the color.
// Multiplying alpha blends the color together more with the background.
func Disabled(c color.NRGBA) (d color.NRGBA) {
	const r = 80 // blend ratio
	lum := approxLuminance(c)
	d = mix(c, color.NRGBA{A: c.A, R: lum, G: lum, B: lum}, r)
	d = MulAlpha(d, 128+32)
	return
}

// Hovered blends dark colors towards white, and light colors towards
// black. It is approximate because it operates in non-linear sRGB space.
func Hovered(c color.NRGBA) (h color.NRGBA) {
	if c.A == 0 {
		// Provide a reasonable default for transparent widgets.
		return color.NRGBA{A: 0x44, R: 0x88, G: 0x88, B: 0x88}
	}
	const ratio = 0x20
	m := color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: c.A}
	if approxLuminance(c) > 128 {
		m = color.NRGBA{A: c.A}
	}
	return mix(m, c, ratio)
}

// mix mixes c1 and c2 weighted by (1 - a/256) and a/256 respectively.
func mix(c1, c2 color.NRGBA, a uint8) color.NRGBA {
	ai := int(a)
	return color.NRGBA{
		R: byte((int(c1.R)*ai + int(c2.R)*(256-ai)) / 256),
		G: byte((int(c1.G)*ai + int(c2.G)*(256-ai)) / 256),
		B: byte((int(c1.B)*ai + int(c2.B)*(256-ai)) / 256),
		A: byte((int(c1.A)*ai + int(c2.A)*(256-ai)) / 256),
	}
}

// approxLuminance is a fast approximate version of RGBA.Luminance.
func approxLuminance(c color.NRGBA) byte {
	const (
		r = 13933 // 0.2126 * 256 * 256
		g = 46871 // 0.7152 * 256 * 256
		b = 4732  // 0.0722 * 256 * 256
		t = r + g + b
	)
	return byte((r*int(c.R) + g*int(c.G) + b*int(c.B)) / t)
}

var (
	iconsOnce sync.Once
	icns      *Icons
)

type Icons struct {
	ContentAdd            *widget.Icon
	ActionDelete          *widget.Icon
	CheckBoxBlank         *widget.Icon
	CheckBoxIndeterminate *widget.Icon
	ArrowLeft             *widget.Icon
	ArrowRight            *widget.Icon
	ArrowUp               *widget.Icon
	ArrowDown             *widget.Icon
}

func GetIcons() *Icons {
	iconsOnce.Do(func() {
		icns = new(Icons)
		icns.ContentAdd, _ = widget.NewIcon(icons.ContentAdd)
		icns.ActionDelete, _ = widget.NewIcon(icons.ActionDelete)
		icns.CheckBoxBlank, _ = widget.NewIcon(icons.ToggleCheckBoxOutlineBlank)
		icns.CheckBoxIndeterminate, _ = widget.NewIcon(icons.ToggleIndeterminateCheckBox)
		icns.ArrowLeft, _ = widget.NewIcon(icons.HardwareKeyboardArrowLeft)
		icns.ArrowRight, _ = widget.NewIcon(icons.HardwareKeyboardArrowRight)
		icns.ArrowUp, _ = widget.NewIcon(icons.HardwareKeyboardArrowUp)
		icns.ArrowDown, _ = widget.NewIcon(icons.HardwareKeyboardArrowDown)
	})
	return icns
}

//func (i *Icons) Elems()[]*widget.Icon  {
//
//}

var DeleteIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionDelete)
	return icon
}()

var CircleIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ImageLens)
	return icon
}()

var SaveIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ContentSave)
	return icon
}()

var CopyIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ContentContentCopy)
	return icon
}()

var SearchIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionSearch)
	return icon
}()

var HomeIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionHome)
	return icon
}()

var OtherIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionHelp)
	return icon
}()

var HeartIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionFavorite)
	return icon
}()

var PlusIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ContentAdd)
	return icon
}()

var EditIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ContentCreate)
	return icon
}()

var VisibilityIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionVisibility)
	return icon
}()

var CloseIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationClose)
	return icon
}()

var ForwardIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationChevronRight)
	return icon
}()

var ExpandIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationExpandMore)
	return icon
}()

var UploadIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.FileFileUpload)
	return icon
}()

var MoreVertIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationMoreVert)
	return icon
}()

var SendIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ContentSend)
	return icon
}()

var RefreshIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationRefresh)
	return icon
}()

var CleanIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.EditorFormatColorText)
	return icon
}()

var ActionVisibilityIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionVisibility)
	return icon
}()

var ActionVisibilityOffIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionVisibilityOff)
	return icon
}()

var ActionPermIdentityIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionPermIdentity)
	return icon
}()

var EditorFunctionsIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.EditorFunctions)
	return icon
}()

var EditorBorderAllIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.EditorBorderAll)
	return icon
}()

var MapsDirectionsRunIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.MapsDirectionsRun)
	return icon
}()

var ActionZoomInIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionZoomIn)
	return icon
}()

var NavigationSubdirectoryArrowRightIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationSubdirectoryArrowRight)
	return icon
}()

var ActionInfoOutlineIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionInfoOutline)
	return icon
}()

var ActionStarRateIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionStarRate)
	return icon
}()

var NavigationRefreshIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationRefresh)
	return icon
}()

var AlertErrorOutlineIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.AlertErrorOutline)
	return icon
}()

var AlertErrorIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.AlertError)
	return icon
}()

var ActionCheckCircleIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionCheckCircle)
	return icon
}()

var ActionHighlightOffIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionHighlightOff)
	return icon
}()

var NavigationCancelIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationCancel)
	return icon
}()

var ActionCloseIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ContentClear)
	return icon
}()

var ActionPointIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.AVFiberManualRecord)
	return icon
}()

var ActionMinIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ContentRemove)
	// icon, _ := widget.NewIcon(icons.NotificationDoNotDisturbOn)
	return icon
}()

var ActionFullIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionOpenWith)
	return icon
}()

var ArrowRightIcon *widget.Icon = func() *widget.Icon {
	//icon, _ := widget.NewIcon(icons.HardwareKeyboardArrowRight)
	icon, _ := widget.NewIcon(icons.NavigationChevronRight)
	return icon
}()

var ArrowDownIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.HardwareKeyboardArrowDown)
	return icon
}()

var ArrowUpIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.HardwareKeyboardArrowUp)
	return icon
}()
