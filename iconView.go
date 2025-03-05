package ux

import (
	"fmt"
	"io"
	"strings"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/ddkwork/ux/animationButton"
	"github.com/ddkwork/ux/widget/material"
)

type IconView struct {
	*widget.List
	clickMap map[string]*animationButton.Button
	filter   *Input
	keyWords string
	elements []layout.Widget
}

func NewIconView() *IconView {
	i := &IconView{
		List: &widget.List{
			Scrollbar: widget.Scrollbar{},
			List: layout.List{
				Axis:        layout.Vertical,
				ScrollToEnd: false,
				Alignment:   0,
				Position:    layout.Position{},
			},
		},
		clickMap: make(map[string]*animationButton.Button),
		filter:   NewInput("请输入搜索关键字..."),
		keyWords: "Edi",
		elements: make([]layout.Widget, 0, IconMap.Len()),
	}
	i.filter.SetOnChanged(func(text string) {
		fmt.Println("change:", i.filter.GetText())
		i.keyWords = i.filter.GetText()
	})
	for _, name := range IconMap.Keys() {
		i.clickMap[name] = NewButtonAnimation(name, ContentContentCopyIcon, func(gtx layout.Context) {
			gtx.Execute(clipboard.WriteCmd{Data: io.NopCloser(strings.NewReader(name))})
		})
	}
	return i
}

func (i *IconView) Layout(gtx layout.Context) layout.Dimensions {
	i.getElements() // todo not work
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return i.filter.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.List(th.Theme, i.List).Layout(gtx, len(i.elements), func(gtx layout.Context, index int) layout.Dimensions {
				return i.elements[index](gtx)
			})
		}),
	)
}

func (i *IconView) getElements() {
	for name, v := range IconMap.Range() {
		if i.keyWords == "" || strings.Contains(strings.ToLower(name), strings.ToLower(i.keyWords)) {
			// fmt.Println("keywords:", keyWords, "name:", name)
			i.elements = append(i.elements, func(gtx layout.Context) layout.Dimensions {
				//if clickMap[name].OnClicked(gtx) {
				//	copyResponse(gtx, name)
				//}
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Body1(th.Theme, name).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Spacer{Width: unit.Dp(10)}.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = 80
						return v.Layout(gtx, th.Color.WarningColor)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return i.clickMap[name].Layout(gtx)
					}),
				)
			})
		}
	}
}
