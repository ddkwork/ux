package ux

import (
	"fmt"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/safemap"
	"io"
	"strings"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/ddkwork/ux/animationButton"
	"github.com/ddkwork/ux/widget/material"
)

type IconView struct {
	//*widget.List
	clickMap *safemap.M[string, *animationButton.Button]
	filter   *Input
	keyWords string
	elements []layout.Widget
	flow     *Flow
}

func NewIconView() *IconView {
	i := &IconView{
		//List: &widget.List{
		//	Scrollbar: widget.Scrollbar{},
		//	List: layout.List{
		//		Axis:        layout.Vertical,
		//		ScrollToEnd: false,
		//		Alignment:   0,
		//		Position:    layout.Position{},
		//	},
		//},
		clickMap: new(safemap.M[string, *animationButton.Button]),
		filter:   NewInput("请输入搜索关键字..."),
		keyWords: "Edi",
		elements: make([]layout.Widget, 0, IconMap.Len()),
		flow: &Flow{
			Num:       5,
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
			list: &widget.List{
				Scrollbar: widget.Scrollbar{},
				List:      layout.List{},
			},
		},
	}
	i.filter.SetOnChanged(func(text string) {
		fmt.Println("change:", i.filter.GetText())
		i.keyWords = i.filter.GetText()
	})
	for _, name := range IconMap.Keys() {
		i.clickMap.Set(name, NewButtonAnimation(name, IconMap.GetMust(name), func(gtx layout.Context) {
			gtx.Execute(clipboard.WriteCmd{Data: io.NopCloser(strings.NewReader(name))})
		}))
	}
	return i
}

func (i *IconView) Layout(gtx layout.Context) layout.Dimensions {
	return i.flow.Layout(gtx, i.clickMap.Len(), func(gtx layout.Context, index int) layout.Dimensions {
		gtx.Constraints.Min.X = 400
		gtx.Constraints.Max.X = 400
		return layout.UniformInset(4).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return i.clickMap.Values()[index].Layout(gtx)
		})
	})
	//return outlay.FlowWrap{//卡顿
	//	Axis:      layout.Horizontal,
	//	Alignment: layout.Middle,
	//}.Layout(gtx, i.clickMap.Len(), func(gtx layout.Context, index int) layout.Dimensions {
	//	gtx.Constraints.Min.X = 400
	//	gtx.Constraints.Max.X = 400
	//	return i.clickMap.Values()[index].Layout(gtx)
	//})

	//var children []layout.Widget
	//for _, button := range i.clickMap.Values() {
	//	children = append(children, func(gtx layout.Context) layout.Dimensions {
	//		gtx.Constraints.Min.X = 400
	//		gtx.Constraints.Max.X = 400
	//		return button.Layout(gtx)
	//	})
	//}
	//return outlay.RigidRows{ //卡顿
	//	Axis:         layout.Horizontal,
	//	Alignment:    layout.Middle,
	//	Spacing:      80,
	//	CrossSpacing: 80,
	//	CrossAlign:   80,
	//}.Layout(gtx, children...)

	//i.getElements() // todo not work
	//return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
	//	layout.Rigid(func(gtx layout.Context) layout.Dimensions {
	//		return i.filter.Layout(gtx)
	//	}),
	//	layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
	//	layout.Rigid(func(gtx layout.Context) layout.Dimensions {
	//		return material.List(th.Theme, i.List).Layout(gtx, len(i.elements), func(gtx layout.Context, index int) layout.Dimensions {
	//			return i.elements[index](gtx)
	//		})
	//	}),
	//)
}

func (i *IconView) getElements() {
	for name, v := range IconMap.Range() {
		if i.keyWords == "" || strings.Contains(strings.ToLower(name), strings.ToLower(i.keyWords)) {
			// fmt.Println("keywords:", keyWords, "name:", name)
			i.elements = append(i.elements, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Body1(th.Theme, name).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Spacer{Width: unit.Dp(10)}.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						//gtx.Constraints.Min.X = 80
						return v.Layout(gtx, th.Color.WarningColor)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return i.clickMap.GetMust(name).Layout(gtx)
					}),
				)
			})
		}
	}
}
