package ux

import (
	"image"
	"image/color"
	"strings"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/x/component"
	"github.com/x-module/gioui-plugins/resource"
)

type SearchDropDown struct {
	SearchDropDown  *widget.Bool
	menuContextArea component.ContextArea
	menu            component.MenuState
	list            *widget.List

	minWidth unit.Dp
	width    unit.Dp

	isOpen              bool
	selectedOptionIndex int
	lastSelectedIndex   int
	options             []*SearchDropDownOption

	size image.Point

	borderWidth  unit.Dp
	cornerRadius unit.Dp

	onSelectedChange func(value string)

	SearchInput *Input
}

type SearchDropDownOption struct {
	Text       string
	Value      string
	Identifier string
	clickable  widget.Clickable

	Icon      *widget.Icon
	IconColor color.NRGBA

	isDefault bool
}

// SetWidth 设置width
func (c *SearchDropDown) SetWidth(width unit.Dp) {
	c.width = width
}

func newSearchDropDownOption(text string) *SearchDropDownOption {
	return &SearchDropDownOption{
		Text:  text,
		Value: text,
	}
}

func (o *SearchDropDownOption) WithIdentifier(identifier string) *SearchDropDownOption {
	o.Identifier = identifier
	return o
}

func (o *SearchDropDownOption) WithValue(value string) *SearchDropDownOption {
	o.Value = value
	return o
}

func (o *SearchDropDownOption) WithIcon(icon *widget.Icon, color color.NRGBA) *SearchDropDownOption {
	o.Icon = icon
	o.IconColor = color
	return o
}

func (o *SearchDropDownOption) DefaultSelected() *SearchDropDownOption {
	o.isDefault = true
	return o
}

func (o *SearchDropDownOption) GetText() string {
	if o == nil {
		return ""
	}

	return o.Text
}

func (o *SearchDropDownOption) GetValue() string {
	if o == nil {
		return ""
	}

	return o.Value
}

func (c *SearchDropDown) SetSelected(index int) {
	c.selectedOptionIndex = index
	c.lastSelectedIndex = index
}

func (c *SearchDropDown) SetOnChanged(f func(value string)) {
	c.onSelectedChange = f
}

func NewSearchDropDown() *SearchDropDown {
	c := &SearchDropDown{
		SearchDropDown: &widget.Bool{Value: true},
		menuContextArea: component.ContextArea{
			Activation:       pointer.ButtonPrimary,
			AbsolutePosition: true,
		},
		list: &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		// borderWidth:  unit.Dp(1),
		cornerRadius: unit.Dp(4),
		width:        th.Size.DefaultElementWidth,
		SearchInput:  NewInput("Search..."),
	}
	// c.SearchInput.SetBefore(func(gtx layout.Context) layout.Dimensions {
	// 	return resource.SearchIcon.Layout(gtx, th.Color.DefaultIconColor)
	// })
	c.SearchInput.SetIcon(resource.SearchIcon)
	size := ElementStyle{TextSize: unit.Sp(14), Height: unit.Dp(15), Inset: layout.UniformInset(unit.Dp(5)), IconSize: unit.Dp(18)}
	c.SearchInput.SetSize(size)
	return c
}

func NewSearchDropDownWithoutBorder() *SearchDropDown {
	c := &SearchDropDown{
		SearchDropDown: &widget.Bool{Value: true},
		menuContextArea: component.ContextArea{
			Activation:       pointer.ButtonPrimary,
			AbsolutePosition: true,
		},
		list: &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		borderWidth:  unit.Dp(0),
		cornerRadius: unit.Dp(4),
		// width:        th.Size.DefaultElementWidth,
		SearchInput: NewInput("Search..."),
	}
	return c
}

func (c *SearchDropDown) SelectedIndex() int {
	return c.selectedOptionIndex
}

func (c *SearchDropDown) SetOptions(options []*SearchDropDownOption) {
	c.selectedOptionIndex = 0
	c.options = options
	c.updateMenuItems()
}

func (c *SearchDropDown) GetSelected() *SearchDropDownOption {
	if len(c.options) == 0 {
		return nil
	}
	return c.options[c.selectedOptionIndex]
}

func (c *SearchDropDown) box1(gtx layout.Context, text string, maxWidth unit.Dp) layout.Dimensions {
	borderColor := th.Color.DropDownBorderColor
	if c.isOpen {
		borderColor = th.Color.DropDownHoveredBorderColor
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
	return c.SearchDropDown.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if c.SearchDropDown.Hovered() {
			border.Color = th.Color.DropDownHoveredBorderColor
		}
		return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			// calculate the minimum width of the box, considering icon and padding
			return layout.Background{}.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, gtx.Dp(th.Size.DefaultElementRadiusSize)).Push(gtx.Ops).Pop()
					paint.Fill(gtx.Ops, th.Color.DefaultBgGrayColor)
					return layout.Dimensions{Size: gtx.Constraints.Min}
				},
				func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Dp(maxWidth) - gtx.Dp(8)
					return layout.Inset{Top: unit.Dp(3), Bottom: unit.Dp(3)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						c.SearchInput.SetWidth(c.width).SetHasBorder(false)
						return c.SearchInput.Layout(gtx)
					})
				},
			)
		})
	})
}

func (c *SearchDropDown) box(gtx layout.Context, text string, maxWidth unit.Dp) layout.Dimensions {
	c.SearchInput.SetWidth(c.width).SetHasBorder(false)
	return c.SearchInput.Layout(gtx)
}

func (c *SearchDropDown) SetSize(size image.Point) {
	c.size = size
}

func (c *SearchDropDown) SetMinWidth(minWidth unit.Dp) {
	c.minWidth = minWidth
}

//
//// update
//func (c *SearchDropDown) update(gtx layout.Context) {
//	c.SearchInput.SetonChanged(func(gtx layout.Context) {
//		fmt.Println("current text:", c.SearchInput.GetText())
//		c.updateMenuItems(c.SearchInput.GetText())
//		gtx.Execute(op.InvalidateCmd{})
//	})
//}

// Layout the SearchDropDown.
func (c *SearchDropDown) Layout(gtx layout.Context) layout.Dimensions {
	// c.update(gtx)
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
		if c.onSelectedChange != nil {
			c.onSelectedChange(c.options[c.selectedOptionIndex].Value)
		}
		c.lastSelectedIndex = c.selectedOptionIndex
	}

	if c.minWidth == 0 {
		c.minWidth = unit.Dp(150)
	}

	text := ""
	if c.selectedOptionIndex >= 0 && c.selectedOptionIndex < len(c.options) {
		text = c.options[c.selectedOptionIndex].Text
	}
	box := c.box1(gtx, text, c.width)
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// c.SearchInput.SetWidth(c.width).SetHasBorder(false)
			// return c.SearchInput.Layout(gtx)
			return box
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return c.menuContextArea.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				offset := layout.Inset{
					Top:  unit.Dp(float32(box.Size.Y)/gtx.Metric.PxPerDp + 1),
					Left: unit.Dp(0),
				}
				return offset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Dp(c.minWidth)
					if c.width != 0 {
						gtx.Constraints.Max.X = gtx.Dp(c.width)
					}
					if len(c.menu.Options) > 0 {
						m := component.Menu(th.Theme, &c.menu)
						m.SurfaceStyle.Fill = th.Color.DropDownBgGrayColor
						return m.Layout(gtx)
					}
					return layout.Dimensions{}
				})
			})
		}),
	)
}

// updateMenuItems creates or updates menu items based on options and calculates minWidth.
func (c *SearchDropDown) updateMenuItems() {
	c.menu.Options = c.menu.Options[:0]
	for i, opt := range c.options {
		current := opt
		c.menu.Options = append(c.menu.Options, func(gtx layout.Context) layout.Dimensions {
			itm := component.MenuItem(th.Theme, &current.clickable, current.Text)

			itm.HoverColor = th.Color.DropDownItemHoveredGrayColor
			if current.Icon != nil {
				itm.Icon = current.Icon
				itm.IconColor = current.IconColor
			}
			itm.Label.TextSize = th.Size.DropdownTextSize
			if c.GetSelected().Text == current.Text {
				itm.Label.Color = th.Color.DropDownSelectedItemBgColor
				itm.Label.Font.Weight = font.Bold
			} else {
				itm.Label.Color = th.Color.DefaultTextWhiteColor
			}

			return itm.Layout(gtx)

			bgColor := RowColor(i)
			return Background{bgColor}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return itm.Layout(gtx)
			})
		})
	}
}

// updateMenuItems creates or updates menu items based on options and calculates minWidth.
func (c *SearchDropDown) SetMenuItems(key string) {
	c.menu.Options = c.menu.Options[:0]
	for _, opt := range c.options {
		if !strings.Contains(opt.Text, key) || key == "" {
			continue
		}
		current := opt
		c.menu.Options = append(c.menu.Options, func(gtx layout.Context) layout.Dimensions {
			itm := component.MenuItem(th.Theme, &current.clickable, current.Text)
			itm.HoverColor = th.Color.DropDownItemHoveredGrayColor

			if current.Icon != nil {
				itm.Icon = current.Icon
				itm.IconColor = current.IconColor
			}
			itm.Label.TextSize = th.Size.DropdownTextSize
			if c.GetSelected().Text == current.Text {
				itm.Label.Color = th.Color.DropDownSelectedItemBgColor
				itm.Label.Font.Weight = font.Bold
			} else {
				itm.Label.Color = th.Color.DefaultTextWhiteColor
			}
			return itm.Layout(gtx)
		})
	}
}
