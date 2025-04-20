package ux

import (
	"fmt"
	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/ux/resources/icons"
	"io"
	"strings"
)

type IconView struct {
	filterInput *Input // todo 调用appBar的搜索输入框
	keyWords    string
	filterMap   []layout.Widget
	flow        *Flow
	buttons     []widget.Clickable
}

func NewIconView() *IconView {
	v := &IconView{
		filterInput: NewInput("请输入搜索关键字..."),
		keyWords:    "Edi",
		filterMap:   make([]layout.Widget, 0, icons.IconMap.Len()),
		flow:        NewFlow(5),
		buttons:     make([]widget.Clickable, icons.IconMap.Len()),
	}
	v.filterInput.SetOnChanged(func(text string) {
		fmt.Println("change:", v.filterInput.GetText())
		v.keyWords = v.filterInput.GetText()
	})
	for i, name := range icons.IconMap.Keys() {
		v.flow.AppendElem(i, FlowElemButton{
			Title: name,
			Icon:  icons.IconMap.GetMust(name),
			Do: func(gtx layout.Context) {
				gtx.Execute(clipboard.WriteCmd{Data: io.NopCloser(strings.NewReader(name))})
			},
			ContextMenuItems: []ContextMenuItem{
				{
					Title:         "Balance",
					Icon:          icons.ActionAccountBalanceIcon,
					Can:           func() bool { return true },
					Do:            func() { mylog.Info("Balance item clicked") },
					AppendDivider: false,
					Clickable:     widget.Clickable{},
				},
				{
					Title:         "Account",
					Icon:          icons.ActionAccountBoxIcon,
					Can:           func() bool { return true },
					Do:            func() { mylog.Info("Account item clicked") },
					AppendDivider: false,
					Clickable:     widget.Clickable{},
				},
				{
					Title:         "Cart",
					Icon:          icons.ActionAddShoppingCartIcon,
					Can:           func() bool { return true },
					Do:            func() { mylog.Info("Cart item clicked") },
					AppendDivider: false,
					Clickable:     widget.Clickable{},
				},
			},
		})
	}
	return v
}

func (v *IconView) Layout(gtx layout.Context) layout.Dimensions {
	//v.filter()//todo need layout filterInput
	return v.flow.Layout(gtx)
}

func (v *IconView) filter() {
	i := 0
	for name := range icons.IconMap.Range() {
		i++
		if i > len(v.buttons)-1 {
			break
		}
		if v.keyWords == "" || strings.Contains(strings.ToLower(name), strings.ToLower(v.keyWords)) {
			v.flow.buttons[i].Show = true //todo set
			continue
		}
		v.flow.buttons[i].Show = false
	}
}
