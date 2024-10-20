package ux

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/x-module/gioui-plugins/resource"
)

func Nav(gtx layout.Context) layout.Dimensions {
	var visibilityAnimation component.VisibilityAnimation

	nav := component.NewNav("Hello", "--subtitle")
	nav.AddNavItem(component.NavItem{
		Name: "aaaaaaaaa",
		Icon: resource.PlusIcon,
	})
	nav.AddNavItem(component.NavItem{
		Name: "bbbb",
		Icon: resource.PlusIcon,
	})
	nav.AddNavItem(component.NavItem{
		Name: "cccc",
		Icon: resource.PlusIcon,
	})
	resize := component.Resize{
		Axis:  layout.Horizontal,
		Ratio: 0.2,
	}

	return resize.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			// return nav.Layout(gtx, th.Material(), &visibilityAnimation)
			return nav.LayoutContents(gtx, th.Theme, &visibilityAnimation)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.H6(th.Theme, "Hello").Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Body1(th.Theme, "Subtitle").Layout(gtx)
					}),
				)
			})
		},
		func(gtx layout.Context) layout.Dimensions {
			rect := image.Rectangle{
				Max: image.Point{
					X: gtx.Dp(unit.Dp(4)),
					Y: gtx.Constraints.Max.Y,
				},
			}
			paint.FillShape(gtx.Ops, color.NRGBA{A: 200}, clip.Rect(rect).Op())
			return layout.Dimensions{Size: rect.Max}
		},
	)
}
