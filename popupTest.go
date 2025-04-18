package ux

import (
	"fmt"
	"gioui.org/example/component/icon"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/ddkwork/ux/widget/material"
	"github.com/ddkwork/ux/x/component"
	"image"
	"image/color"
	"strconv"
)

type PopupTest struct {
	RedButton, GreenButton, BlueButton       widget.Clickable
	BalanceButton, AccountButton, CartButton widget.Clickable
	component.ContextArea
	component.MenuState
	MenuInit bool
	widget.List
	RowClicks []widget.Clickable
	DrawRow   func(gtx layout.Context, index int) layout.Dimensions
}

func NewPopupTest(length int, DrawRow func(gtx layout.Context, index int) layout.Dimensions) *PopupTest {
	return &PopupTest{
		RedButton:     widget.Clickable{},
		GreenButton:   widget.Clickable{},
		BlueButton:    widget.Clickable{},
		BalanceButton: widget.Clickable{},
		AccountButton: widget.Clickable{},
		CartButton:    widget.Clickable{},
		ContextArea:   component.ContextArea{},
		MenuState:     component.MenuState{},
		MenuInit:      false,
		List: widget.List{
			Scrollbar: widget.Scrollbar{},
			List: layout.List{
				Axis:        layout.Vertical,
				ScrollToEnd: false,
				Alignment:   0,
				Position:    layout.Position{},
			},
		},
		RowClicks: make([]widget.Clickable, length),
		DrawRow:   DrawRow,
	}
}

func (p *PopupTest) Layout(gtx layout.Context) layout.Dimensions {
	if !p.MenuInit {
		p.MenuState = component.MenuState{
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
				component.MenuItem(th, &p.RedButton, "Red").Layout, //todo
				component.MenuItem(th, &p.GreenButton, "Green").Layout,
				component.MenuItem(th, &p.BlueButton, "Blue").Layout,
				func(gtx C) D {
					item := component.MenuItem(th, &p.BalanceButton, "Balance")
					item.Icon = icon.AccountBalanceIcon
					item.Hint = component.MenuHintText(th, "Hint")
					return item.Layout(gtx)
				},
				func(gtx C) D {
					item := component.MenuItem(th, &p.AccountButton, "Account")
					item.Icon = icon.AccountBoxIcon
					item.Hint = component.MenuHintText(th, "Hint")
					return item.Layout(gtx)
				},
				func(gtx C) D {
					item := component.MenuItem(th, &p.CartButton, "Cart")
					item.Icon = icon.CartIcon
					item.Hint = component.MenuHintText(th, "Hint")
					return item.Layout(gtx)
				},
			},
		}
	}
	if p.RedButton.Clicked(gtx) { //todo popup item callback
		println("Red button clicked")
	}
	if p.GreenButton.Clicked(gtx) {
		println("Green button clicked")
	}
	if p.BlueButton.Clicked(gtx) {
		println("Blue button clicked")
	}
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return material.List(th, &p.List).Layout(gtx, len(p.RowClicks), func(gtx layout.Context, index int) layout.Dimensions {
				rowClicks := &p.RowClicks[index]
				return rowClicks.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					if p.DrawRow == nil {
						if event, b := gtx.Event(pointer.Filter{Target: rowClicks, Kinds: pointer.Press | pointer.Release}); b {
							if e, ok := event.(pointer.Event); ok {
								if e.Kind == pointer.Press {
									switch {
									case e.Buttons.Contain(pointer.ButtonPrimary):
										println("Row selected (left click) " + strconv.Itoa(index)) //todo row selected callback
									case e.Buttons.Contain(pointer.ButtonSecondary):
										println("Row selected (right click)" + strconv.Itoa(index)) //todo row popup callback and return selected row index
									}
								}
							}
						}
						return Background{Color: RowColor(index)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							//todo Background not working
							return material.Button(th, rowClicks, "item"+fmt.Sprintf("%d", index)).Layout(gtx)
						})
					}
					return p.DrawRow(gtx, index)
				})
			})
		}),
		layout.Expanded(func(gtx C) D {
			return p.ContextArea.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min = image.Point{}
				return component.Menu(th, &p.MenuState).Layout(gtx) //所有行的item公用一个popup菜单而不是没行popup一个
			})
		}),
	)
}

func (p *PopupTest) drawContextArea(gtx C, th *material.Theme) D { //popup区域的背景色，位置，四角弧度
	return layout.Center.Layout(gtx, func(gtx C) D {
		menuStyle := component.Menu(th, &p.MenuState)
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
