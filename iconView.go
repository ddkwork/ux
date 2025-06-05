package ux

import (
	"fmt"
	"io"
	"strings"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/safemap"
	"github.com/ddkwork/ux/resources/images"
)

type IconView struct {
	filterInput *Input // todo 调用appBar的搜索输入框
	keyWords    string
	flow        *Flow
	buttons     []*widget.Clickable
	elems       *safemap.M[*ButtonAnimation, *ContextMenu]
}

func NewIconView() *IconView {
	view := &IconView{
		filterInput: NewInput("请输入搜索关键字..."),
		keyWords:    "Edi",
		flow:        nil,
		buttons:     make([]*widget.Clickable, images.IconMap.Len()),
		elems:       nil,
	}
	view.filterInput.SetOnChanged(func(text string) {
		fmt.Println("change:", view.filterInput.GetText())
		view.keyWords = view.filterInput.GetText()
	})
	m := new(safemap.M[*ButtonAnimation, *ContextMenu])
	for i, name := range images.IconMap.Keys() {
		view.buttons[i] = &widget.Clickable{}
		k := NewButton(view.buttons[i], images.IconMap.GetMust(name), name, func(gtx layout.Context) {
			gtx.Execute(clipboard.WriteCmd{Data: io.NopCloser(strings.NewReader(name))})
		})
		v := NewContextMenu()
		v.Once.Do(func() {
			v.AddItem(ContextMenuItem{
				Title:         "Balance",
				Icon:          images.ActionAccountBalanceIcon,
				Can:           func() bool { return true },
				Do:            func() { mylog.Info("Balance item clicked") },
				AppendDivider: false,
				Clickable:     widget.Clickable{},
			})
			v.AddItem(ContextMenuItem{
				Title:         "Account",
				Icon:          images.ActionAccountBoxIcon,
				Can:           func() bool { return true },
				Do:            func() { mylog.Info("Account item clicked") },
				AppendDivider: false,
				Clickable:     widget.Clickable{},
			})
			v.AddItem(ContextMenuItem{
				Title:         "Cart",
				Icon:          images.ActionAddShoppingCartIcon,
				Can:           func() bool { return true },
				Do:            func() { mylog.Info("Cart item clicked") },
				AppendDivider: false,
				Clickable:     widget.Clickable{},
			})
		})
		m.Set(k, v)
	}
	view.elems = m
	view.flow = NewFlow(5, view.elems)
	return view
}

func (v *IconView) Layout(gtx layout.Context) layout.Dimensions {
	// v.filter() //todo need layout filterInput
	return v.flow.Layout(gtx)
}

// func (v *IconView) filter() {
//	i := 0
//	for name := range images.IconMap.Range() {
//		i++
//		if i > len(v.buttons)-1 {
//			break
//		}
//		if v.keyWords == "" || strings.Contains(strings.ToLower(name), strings.ToLower(v.keyWords)) {
//			v.flow.buttons[i].Show = false // todo set
//			continue
//		}
//		v.flow.buttons[i].Show = true // todo set
//	}
// }
