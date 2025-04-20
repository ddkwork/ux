package ux

import (
	"fmt"
	"io"
	"strings"

	"github.com/ddkwork/ux/resources/icons"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/safemap"
	"github.com/ddkwork/ux/animationButton"
)

type IconView struct {
	clickMap    *safemap.M[string, *animationButton.Button]
	filterInput *Input // todo 调用appBar的搜索输入框
	keyWords    string
	filterMap   []layout.Widget
	flow        Flow
}

func NewIconView() *IconView {
	v := &IconView{
		clickMap:    new(safemap.M[string, *animationButton.Button]),
		filterInput: NewInput("请输入搜索关键字..."),
		keyWords:    "Edi",
		filterMap:   make([]layout.Widget, 0, icons.IconMap.Len()),
		flow: Flow{
			Num:       5,
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
			list: widget.List{
				Scrollbar: widget.Scrollbar{},
				List:      layout.List{},
			},
		},
	}
	v.filterInput.SetOnChanged(func(text string) {
		fmt.Println("change:", v.filterInput.GetText())
		v.keyWords = v.filterInput.GetText()
	})
	for _, name := range icons.IconMap.Keys() {
		v.clickMap.Set(name, NewButtonAnimation(name, icons.IconMap.GetMust(name), func(gtx layout.Context) { // todo 增加右键回调弹出菜单
			gtx.Execute(clipboard.WriteCmd{Data: io.NopCloser(strings.NewReader(name))})
		}))
	}
	return v
}

func (v *IconView) Layout(gtx layout.Context) layout.Dimensions {
	return v.flow.Layout(gtx, v.clickMap.Len(), func(gtx layout.Context, i int) layout.Dimensions {
		gtx.Constraints.Min.X = 400
		gtx.Constraints.Max.X = 400
		v.filter()
		if v.filterMap != nil {
			return layout.UniformInset(4).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return v.filterMap[i](gtx)
			})
		}
		return layout.UniformInset(4).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return v.clickMap.Values()[i].Layout(gtx)
		})
	})
}

func (v *IconView) filter() {
	for name := range icons.IconMap.Range() {
		if v.keyWords == "" || strings.Contains(strings.ToLower(name), strings.ToLower(v.keyWords)) {
			v.filterMap = append(v.filterMap, v.clickMap.GetMust(name).Layout)
		}
	}
}
