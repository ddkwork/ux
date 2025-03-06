package ux

import (
	"fmt"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/safemap"
	"io"
	"strings"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"github.com/ddkwork/ux/animationButton"
)

type IconView struct {
	clickMap    *safemap.M[string, *animationButton.Button]
	filterInput *Input //todo 调用appBar的搜索输入框
	keyWords    string
	filterMap   []layout.Widget
	flow        *Flow
}

func NewIconView() *IconView {
	i := &IconView{
		clickMap:    new(safemap.M[string, *animationButton.Button]),
		filterInput: NewInput("请输入搜索关键字..."),
		keyWords:    "Edi",
		filterMap:   make([]layout.Widget, 0, IconMap.Len()),
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
	i.filterInput.SetOnChanged(func(text string) {
		fmt.Println("change:", i.filterInput.GetText())
		i.keyWords = i.filterInput.GetText()
	})
	for _, name := range IconMap.Keys() {
		i.clickMap.Set(name, NewButtonAnimation(name, IconMap.GetMust(name), func(gtx layout.Context) { //todo 增加右键回调弹出菜单
			gtx.Execute(clipboard.WriteCmd{Data: io.NopCloser(strings.NewReader(name))})
		}))
	}
	return i
}

func (i *IconView) Layout(gtx layout.Context) layout.Dimensions {
	return i.flow.Layout(gtx, i.clickMap.Len(), func(gtx layout.Context, index int) layout.Dimensions {
		gtx.Constraints.Min.X = 400
		gtx.Constraints.Max.X = 400
		i.filter()
		if i.filterMap != nil {
			return layout.UniformInset(4).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return i.filterMap[index](gtx)
			})
		}
		return layout.UniformInset(4).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return i.clickMap.Values()[index].Layout(gtx)
		})
	})
}

func (i *IconView) filter() {
	for name := range IconMap.Range() {
		if i.keyWords == "" || strings.Contains(strings.ToLower(name), strings.ToLower(i.keyWords)) {
			i.filterMap = append(i.filterMap, i.clickMap.GetMust(name).Layout)
		}
	}
}
