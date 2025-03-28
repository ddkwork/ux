package main

import (
	"fmt"
	"image/color"

	"github.com/ddkwork/ux/widget/material"
	"github.com/ddkwork/ux/x/component"

	"github.com/ddkwork/ux"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type MenuObj struct {
	balanceButton, accountButton, cartButton widget.Clickable
	menuState                                component.MenuState
	menuInit                                 bool
	menuDemoList                             widget.List
	contextAreas                             []*component.ContextArea
	widget.List
}

func NewMenuObj() *MenuObj {
	m := &MenuObj{}
	return m
}

type (
	C = layout.Context
	D = layout.Dimensions
)

func (p *MenuObj) Layout(gtx layout.Context) layout.Dimensions {
	// paint.Fill(gtx.Ops, color.NRGBA(Grey800))
	if !p.menuInit {
		p.menuState = component.MenuState{
			Options: []func(gtx C) D{
				func(gtx C) D {
					item := component.MenuItem(ux.ThemeDefault().Theme, &p.balanceButton, "Balance")
					item.Icon = ux.ContentCreateIcon
					item.Hint = component.MenuHintText(ux.ThemeDefault().Theme, "Hint")
					return item.Layout(gtx)
				},
				func(gtx C) D {
					item := component.MenuItem(ux.ThemeDefault().Theme, &p.accountButton, "Account")
					item.Icon = ux.ActionCodeIcon
					item.Hint = component.MenuHintText(ux.ThemeDefault().Theme, "Hint")
					return item.Layout(gtx)
				},
				func(gtx C) D {
					item := component.MenuItem(ux.ThemeDefault().Theme, &p.cartButton, "Cart")
					item.Icon = ux.ActionSearchIcon
					item.Hint = component.MenuHintText(ux.ThemeDefault().Theme, "Hint")
					return item.Layout(gtx)
				},
			},
		}
	}
	if p.balanceButton.Clicked(gtx) {
		println("balance")
	}
	if p.accountButton.Clicked(gtx) {
		println("account")
	}
	if p.cartButton.Clicked(gtx) {
		println("cart")
	}
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(.5, func(gtx C) D {
			p.menuDemoList.Axis = layout.Vertical
			return material.List(ux.ThemeDefault().Theme, &p.menuDemoList).Layout(gtx, 30, func(gtx C, index int) D {
				if len(p.contextAreas) < index+1 {
					p.contextAreas = append(p.contextAreas, &component.ContextArea{})
				}
				contextArea := p.contextAreas[index]
				return layout.Stack{}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.UniformInset(unit.Dp(8)).Layout(gtx, material.Body1(ux.ThemeDefault().Theme, fmt.Sprintf("Item %d", index)).Layout)
					}),
					layout.Expanded(func(gtx C) D {
						return contextArea.Layout(gtx, func(gtx C) D {
							gtx.Constraints.Min.X = 0
							return p.drawContextArea(gtx, ux.ThemeDefault().Theme)
							return component.Menu(ux.ThemeDefault().Theme, &p.menuState).Layout(gtx)
						})
					}),
				)
			})
		}),
	)
}

func (p *MenuObj) drawContextArea(gtx C, th *material.Theme) D {
	return layout.Center.Layout(gtx, func(gtx C) D { // 重置min x y 到0，并根据max x y 计算弹出菜单的合适大小
		// mylog.Struct("todo",gtx.Constraints)
		menuStyle := component.Menu(th, &p.menuState)
		menuStyle.SurfaceStyle = component.SurfaceStyle{
			Theme: th,
			ShadowStyle: component.ShadowStyle{
				CornerRadius: 18, // 弹出菜单的椭圆角度
				Elevation:    0,
				// AmbientColor:  color.NRGBA(Blue400),
				// PenumbraColor: color.NRGBA(Blue400),
				// UmbraColor:    color.NRGBA(Blue400),
			},
			Fill: color.NRGBA{R: 50, G: 50, B: 50, A: 255}, // 弹出菜单的背景色
		}
		return menuStyle.Layout(gtx)
	})
}
