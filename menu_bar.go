package ux

import (
	"fmt"
	"image"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/x-module/gioui-plugins/utils"
)

type MenuBarAction func()

type MenuBar struct {
	menus []MenuBarItem
}

func NewMenuBar() *MenuBar {
	return &MenuBar{}
}

func (m *MenuBar) AddMenuBarItem(menus []MenuBarItem) *MenuBar {
	for key := range menus {
		menus[key].menuContextArea = component.ContextArea{
			Activation:       pointer.ButtonPrimary,
			AbsolutePosition: true,
		}
	}
	m.menus = menus
	return m
}

type MenuBarItem struct {
	Title           string
	menuContextArea component.ContextArea
	Items           []MenuBarItemElement
}

type MenuBarItemElement struct {
	Name   string
	click  widget.Clickable
	Action MenuBarAction
}

func (m *MenuBar) Layout(gtx layout.Context) layout.Dimensions {
	var items []layout.FlexChild
	for key := range m.menus {
		if len(items) > 0 {
			items = append(items, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Body1(th.Theme, "|").Layout(gtx)
					}),
				)
			}))
		}
		items = append(items, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Stack{}.Layout(gtx,
						layout.Stacked(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{Left: unit.Dp(3), Right: unit.Dp(3)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return material.Body2(th.Theme, m.menus[key].Title).Layout(gtx)
							})
						}),
						layout.Expanded(func(gtx layout.Context) layout.Dimensions {
							var subItems []layout.FlexChild
							for subKey := range m.menus[key].Items {
								subItems = append(subItems, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											return layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(5)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												return m.menus[key].Items[subKey].click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
													if m.menus[key].Items[subKey].click.Hovered() {
														fmt.Println("hovered")
														utils.DrawBackground(gtx, gtx.Constraints.Max, th.Color.MenuBarHoveredColor)
														gtx.Execute(op.InvalidateCmd{})
													}
													return material.Body2(th.Theme, m.menus[key].Items[subKey].Name).Layout(gtx)
												})
											})
										}),
									)
								}))
							}

							return m.menus[key].menuContextArea.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								offset := layout.Inset{
									Top: unit.Dp(20),
								}
								return offset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return widget.Border{
										Color:        th.Color.MenuBarBorderColor,
										Width:        unit.Dp(1),
										CornerRadius: th.Size.DefaultElementRadiusSize,
									}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return layout.Background{}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, gtx.Dp(th.Size.DefaultElementRadiusSize)).Push(gtx.Ops).Pop()
											paint.Fill(gtx.Ops, th.Color.MenuBarBgColor)
											return layout.Dimensions{Size: gtx.Constraints.Min}
										}, func(gtx layout.Context) layout.Dimensions {
											return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												return layout.Flex{Axis: layout.Vertical}.Layout(gtx, subItems...)
											})
										})
									})
								})
							})
						}),
					)
				}),
			)
		}))
	}
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, items...)
}

//
// return layout.Stack{}.Layout(gtx,
// layout.Stacked(func(gtx layout.Context) layout.Dimensions {
// 	return layout.Inset{Top: unit.Dp(3), Bottom: unit.Dp(3)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
// 		return Body1(m.theme, "▼").Layout(gtx)
// 	})
// }),
// )
