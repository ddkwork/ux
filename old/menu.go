/**
 * Created by Goland
 * @file   menu.go
 * @author 李锦 <lijin@cavemanstudio.net>
 * @date   2024/8/7 14:20
 * @desc   menu.go
 */

package main

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/ddkwork/ux"
	"github.com/ddkwork/ux/widget/material"
)

type Menu struct {
	list      widget.List
	menuItems []MenuItem
	textSize  unit.Sp
	// clickedIndex int
}

func NewMenu() *Menu {
	return &Menu{
		list: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		textSize: unit.Sp(13),
	}
}

type MenuItem struct {
	Icon     *widget.Icon
	Text     string
	SubMenu  []MenuItem
	click    widget.Clickable
	selected bool
	list     widget.List
}

// SetMenuItemOptions 设置menuItemOptions
func (m *Menu) SetMenuItemOptions(menuItems []MenuItem) {
	for i := range menuItems {
		menuItems[i].list = widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		}
	}
	m.menuItems = menuItems
}

func (m *MenuItem) layout(gtx layout.Context, parent *Menu, haasLine bool) layout.Dimensions {
	if m.click.Clicked(gtx) {
		for i := range parent.menuItems {
			parent.menuItems[i].selected = false
			for j := range parent.menuItems[i].SubMenu {
				parent.menuItems[i].SubMenu[j].selected = false
			}
		}
		m.selected = true
		for i, menu := range parent.menuItems {
			for _, subMenu := range menu.SubMenu {
				if subMenu.selected {
					parent.menuItems[i].selected = true
				}
			}
		}
	}

	textColor := th.Color.MenuItemTextColor
	if (m.click.Hovered() || m.selected) && m.SubMenu == nil {
		textColor = th.Color.MenuItemTextSelectedColor
	}

	return m.click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if m.Icon != nil {
							gtx.Constraints.Max.X = gtx.Dp(unit.Dp(17))
							return m.Icon.Layout(gtx, textColor)
						}
						return layout.Dimensions{}
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: unit.Dp(5)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							label := material.Label(th, unit.Sp(13), m.Text)
							label.Color = textColor
							return label.Layout(gtx)
						})
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if m.selected && m.SubMenu != nil {
					return layout.Inset{Top: unit.Dp(5)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return m.list.Layout(gtx, len(m.SubMenu), func(gtx layout.Context, index int) layout.Dimensions {
							gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(25))
							gtx.Constraints.Min.X = gtx.Dp(unit.Dp(70))
							return layout.Inset{Left: unit.Dp(30)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return m.SubMenu[index].layout(gtx, parent, false)
							})
						})
					})
				}
				return layout.Dimensions{}
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if haasLine {
					return layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return ux.NewLine().Color(th.Color.DefaultLineColor).Line(gtx, f32.Pt(0, 0), f32.Pt(float32(gtx.Constraints.Max.X), 0)).Layout(gtx)
					})
				} else {
					return layout.Dimensions{}
				}
			}),
		)
	})
}

func (m *Menu) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(15), Right: unit.Dp(5)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return m.list.Layout(gtx, len(m.menuItems), func(gtx layout.Context, index int) layout.Dimensions {
			gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(30))
			gtx.Constraints.Min.X = gtx.Dp(unit.Dp(100))
			return m.menuItems[index].layout(gtx, m, true)
		})
	})
}
