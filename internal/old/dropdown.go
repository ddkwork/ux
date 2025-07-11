package main

import (
	"image"
	"image/color"

	"github.com/ddkwork/ux/resources/colors"
	"github.com/ddkwork/ux/resources/images"

	"github.com/ddkwork/ux/widget/material"
	"github.com/ddkwork/ux/x/component"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type DropDown struct {
	menuContextArea component.ContextArea
	menu            component.MenuState
	list            *widget.List

	MinWidth unit.Dp
	MaxWidth unit.Dp
	menuInit bool

	isOpen              bool
	selectedOptionIndex int
	lastSelectedIndex   int
	options             []*DropDownOption

	size image.Point

	borderWidth  unit.Dp
	cornerRadius unit.Dp

	onValueChange func(value string)
}

type DropDownOption struct {
	Text       string
	Value      string
	Identifier string
	clickable  widget.Clickable

	Icon      []byte
	IconColor color.NRGBA
	IconSize  unit.Dp

	isDivider bool
	isDefault bool
}

func NewDropDownOption(text string) *DropDownOption {
	return &DropDownOption{
		Text:      text,
		isDivider: false,
	}
}

func NewDropDownDivider() *DropDownOption {
	return &DropDownOption{
		isDivider: true,
	}
}

func (o *DropDownOption) WithIdentifier(identifier string) *DropDownOption {
	o.Identifier = identifier
	return o
}

func (o *DropDownOption) WithValue(value string) *DropDownOption {
	o.Value = value
	return o
}

func (o *DropDownOption) WithIcon(icon []byte, color color.NRGBA, size unit.Dp) *DropDownOption {
	o.Icon = icon
	o.IconColor = color
	o.IconSize = size
	return o
}

func (o *DropDownOption) DefaultSelected() *DropDownOption {
	o.isDefault = true
	return o
}

func (o *DropDownOption) GetText() string {
	if o == nil {
		return ""
	}

	return o.Text
}

func (o *DropDownOption) GetValue() string {
	if o == nil {
		return ""
	}

	return o.Value
}

func (c *DropDown) SetSelected(index int) {
	c.selectedOptionIndex = index
	c.lastSelectedIndex = index
}

func (c *DropDown) SetOnChanged(f func(value string)) {
	c.onValueChange = f
}

func (c *DropDown) SetSelectedByTitle(title string) {
	if len(c.options) == 0 {
		return
	}

	for i, opt := range c.options {
		if opt.Text == title {
			c.selectedOptionIndex = i
			c.lastSelectedIndex = i
			break
		}
	}
}

func (c *DropDown) SetSelectedByIdentifier(identifier string) {
	for i, opt := range c.options {
		if opt.Identifier == identifier {
			c.selectedOptionIndex = i
			c.lastSelectedIndex = i
			break
		}
	}
}

func (c *DropDown) SetSelectedByValue(value string) {
	for i, opt := range c.options {
		if opt.Value == value {
			c.selectedOptionIndex = i
			c.lastSelectedIndex = i
			break
		}
	}
}

func NewDropDown(options ...*DropDownOption) *DropDown {
	c := &DropDown{
		menuContextArea: component.ContextArea{
			Activation:       pointer.ButtonPrimary,
			AbsolutePosition: true,
		},
		list: &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		options:      options,
		borderWidth:  unit.Dp(1),
		cornerRadius: unit.Dp(4),

		menuInit: true,
	}

	return c
}

func NewDropDownWithoutBorder(options ...*DropDownOption) *DropDown {
	c := &DropDown{
		menuContextArea: component.ContextArea{
			Activation:       pointer.ButtonPrimary,
			AbsolutePosition: true,
		},
		list: &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		options: options,

		menuInit: true,
	}

	return c
}

func (c *DropDown) SelectedIndex() int {
	return c.selectedOptionIndex
}

func (c *DropDown) SetOptions(options ...*DropDownOption) {
	c.selectedOptionIndex = 0
	c.options = options
	if len(c.options) > 0 {
		c.menuInit = true
	}
}

func (c *DropDown) GetSelected() *DropDownOption {
	if len(c.options) == 0 {
		return nil
	}

	return c.options[c.selectedOptionIndex]
}

func (c *DropDown) box(gtx layout.Context, text string, maxWidth unit.Dp) layout.Dimensions {
	borderColor := th.BorderBlueColor // theme.BorderColor
	if c.isOpen {
		borderColor = th.InputActivatedBorderColor // theme.BorderColorFocused
	}

	border := widget.Border{
		Color:        borderColor,
		Width:        c.borderWidth,
		CornerRadius: c.cornerRadius,
	}

	if maxWidth == 0 {
		maxWidth = unit.Dp(gtx.Constraints.Max.X)
	}

	c.size.X = gtx.Dp(maxWidth)

	return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// calculate the minimum width of the box, considering images and padding
		gtx.Constraints.Min.X = gtx.Dp(maxWidth) - gtx.Dp(8)
		return layout.Inset{
			Top:    4,
			Bottom: 4,
			Left:   8,
			Right:  4,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return material.Label(th, th.TextSize, text).Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Dp(16)
					return images.Layout(gtx, images.NavigationExpandMoreIcon, th.Fg, 0)
				}),
			)
		})
	})
}

func (c *DropDown) SetSize(size image.Point) {
	c.size = size
}

// Layout the DropDown.
func (c *DropDown) Layout(gtx layout.Context) layout.Dimensions {
	c.isOpen = c.menuContextArea.Active()

	for i, opt := range c.options {
		if opt.isDefault {
			c.selectedOptionIndex = i
		}

		for opt.clickable.Clicked(gtx) {
			c.isOpen = false
			c.selectedOptionIndex = i
		}
	}

	if c.selectedOptionIndex != c.lastSelectedIndex {
		if c.onValueChange != nil {
			go c.onValueChange(c.options[c.selectedOptionIndex].Value)
		}
		c.lastSelectedIndex = c.selectedOptionIndex
	}

	// Update menu items only if options change
	if c.menuInit {
		c.menuInit = false
		c.updateMenuItems()
	}

	if c.MinWidth == 0 {
		c.MinWidth = unit.Dp(150)
	}

	text := ""
	if c.selectedOptionIndex >= 0 && c.selectedOptionIndex < len(c.options) {
		text = c.options[c.selectedOptionIndex].Text
	}

	box := c.box(gtx, text, c.MaxWidth)
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return box
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return c.menuContextArea.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				offset := layout.Inset{
					Top:  unit.Dp(float32(box.Size.Y)/gtx.Metric.PxPerDp + 1),
					Left: unit.Dp(4),
				}
				return offset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Dp(c.MinWidth)
					if c.MaxWidth != 0 {
						gtx.Constraints.Max.X = gtx.Dp(c.MaxWidth)
					}
					m := component.Menu(th, &c.menu)
					m.SurfaceStyle.Fill = th.FocusedBgColor // theme.DropDownMenuBgColor
					return m.Layout(gtx)
				})
			})
		}),
	)
}

// updateMenuItems creates or updates menu items based on options and calculates minWidth.
func (c *DropDown) updateMenuItems() {
	c.menu.Options = c.menu.Options[:0]
	for _, opt := range c.options {
		c.menu.Options = append(c.menu.Options, func(gtx layout.Context) layout.Dimensions {
			if opt.isDivider {
				dv := component.Divider(th)
				dv.Fill = th.FullIconColor // theme.BorderColor
				return dv.Layout(gtx)
			}

			itm := component.MenuItem(th, &opt.clickable, opt.Text)
			if opt.Icon != nil {
				itm.Icon = opt.Icon
				itm.IconColor = opt.IconColor
				itm.IconSize = opt.IconSize
			}

			itm.Label.Color = colors.White
			return itm.Layout(gtx)
		})
	}
}
