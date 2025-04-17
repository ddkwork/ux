package main

import (
	"fmt"
	"gioui.org/example/component/icon"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/ddkwork/ux"
	"github.com/ddkwork/ux/widget/material"
	"github.com/ddkwork/ux/x/component"
	"image"
	"image/color"
	"strconv"
)

type MenuObj struct {
	redButton, greenButton, blueButton       widget.Clickable
	balanceButton, accountButton, cartButton widget.Clickable
	contextArea                              component.ContextArea
	menuState                                component.MenuState
	menuInit                                 bool
	widget.List
	rowClicks []widget.Clickable
}

func NewMenuObj() *MenuObj {
	m := &MenuObj{
		redButton:     widget.Clickable{},
		greenButton:   widget.Clickable{},
		blueButton:    widget.Clickable{},
		balanceButton: widget.Clickable{},
		accountButton: widget.Clickable{},
		cartButton:    widget.Clickable{},
		contextArea: component.ContextArea{
			LongPressDuration: 0,
			Activation:        0,
			AbsolutePosition:  false,
			PositionHint:      0,
		},
		menuState: component.MenuState{},
		menuInit:  false,
		List: widget.List{
			Scrollbar: widget.Scrollbar{},
			List: layout.List{
				Axis:        layout.Vertical, //todo bug
				ScrollToEnd: false,
				Alignment:   0,
				Position:    layout.Position{},
			},
		},
		rowClicks: make([]widget.Clickable, 100),
	}
	return m
}

type (
	C = layout.Context
	D = layout.Dimensions
)

func (p *MenuObj) Layout(gtx layout.Context) layout.Dimensions {
	if !p.menuInit {
		p.menuState = component.MenuState{
			Options: []func(gtx C) D{
				//func(gtx C) D {
				//	return layout.Inset{
				//		Left:  unit.Dp(16),
				//		Right: unit.Dp(16),
				//	}.Layout(gtx, material.Body1(th, "Menus support arbitrary widgets.\nThis is just a label!\nHere's a loader:").Layout)
				//},
				//component.Divider(th).Layout,
				//func(gtx C) D {
				//	return layout.Inset{
				//		Top:    unit.Dp(4),
				//		Bottom: unit.Dp(4),
				//		Left:   unit.Dp(16),
				//		Right:  unit.Dp(16),
				//	}.Layout(gtx, func(gtx C) D {
				//		gtx.Constraints.Max.X = gtx.Dp(unit.Dp(24))
				//		gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(24))
				//		return material.Loader(th).Layout(gtx)
				//	})
				//},
				//component.SubheadingDivider(th, "Colors").Layout,
				component.MenuItem(th.Theme, &p.redButton, "Red").Layout,
				component.MenuItem(th.Theme, &p.greenButton, "Green").Layout,
				component.MenuItem(th.Theme, &p.blueButton, "Blue").Layout,
				func(gtx C) D {
					item := component.MenuItem(th.Theme, &p.balanceButton, "Balance")
					item.Icon = icon.AccountBalanceIcon
					item.Hint = component.MenuHintText(th.Theme, "Hint")
					return item.Layout(gtx)
				},
				func(gtx C) D {
					item := component.MenuItem(th.Theme, &p.accountButton, "Account")
					item.Icon = icon.AccountBoxIcon
					item.Hint = component.MenuHintText(th.Theme, "Hint")
					return item.Layout(gtx)
				},
				func(gtx C) D {
					item := component.MenuItem(th.Theme, &p.cartButton, "Cart")
					item.Icon = icon.CartIcon
					item.Hint = component.MenuHintText(th.Theme, "Hint")
					return item.Layout(gtx)
				},
			},
		}
	}
	if p.redButton.Clicked(gtx) {
		println("Red button clicked")
	}
	if p.greenButton.Clicked(gtx) {
		println("Green button clicked")
	}
	if p.blueButton.Clicked(gtx) {
		println("Blue button clicked")
	}
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return material.List(th.Theme, &p.List).Layout(gtx, 100, func(gtx layout.Context, index int) layout.Dimensions {
				rowClicks := &p.rowClicks[index]
				return rowClicks.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					//gtx.Constraints.Min.X = 300 //todo
					if event, b := gtx.Event(pointer.Filter{Target: rowClicks, Kinds: pointer.Press | pointer.Release}); b {
						if e, ok := event.(pointer.Event); ok {
							if e.Kind == pointer.Press {
								switch {
								case e.Buttons.Contain(pointer.ButtonPrimary):
									println("Row selected (left click) " + strconv.Itoa(index))
								case e.Buttons.Contain(pointer.ButtonSecondary):
									println("Row selected (right click)" + strconv.Itoa(index))
								}
							}
						}
					}
					return material.Button(th.Theme, rowClicks, "item"+fmt.Sprintf("%d", index)).Layout(gtx)
				})
			})
		}),
		layout.Expanded(func(gtx C) D {
			return p.contextArea.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min = image.Point{}
				return component.Menu(th.Theme, &p.menuState).Layout(gtx)
			})
		}),
	)
}

func (p *MenuObj) drawContextArea(gtx C, th *material.Theme) D {
	return layout.Center.Layout(gtx, func(gtx C) D {
		menuStyle := component.Menu(th, &p.menuState)
		menuStyle.SurfaceStyle = component.SurfaceStyle{
			Theme: th,
			ShadowStyle: component.ShadowStyle{
				CornerRadius: 18,
				Elevation:    0,
			},
			Fill: color.NRGBA{R: 50, G: 50, B: 50, A: 255},
		}
		return menuStyle.Layout(gtx)
	})
}

var th = ux.NewTheme()
