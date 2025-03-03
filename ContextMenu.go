package ux

import (
	"gioui.org/widget"
	"github.com/ddkwork/ux/x/component"

	"github.com/ddkwork/ux/giosvg"
	//
)

type ContextMenuItem struct {
	Title            string       // 菜单项标题
	Icon             *giosvg.Icon // 可选的图标
	Can              func() bool  // 是否绘制取决于当前渲染的行，回调内需要传递当前渲染的节点给回调，说白了这里是绘制条件，下面的do是业务逻辑，回调内传入的形参节点不一样
	Do               func()       // 调用被选中节点来操作业务逻辑
	AppendDivider    bool         // 是否添加分割线
	widget.Clickable              // 可点击的控件
}

type ContextMenu struct {
	Items []*ContextMenuItem
	component.MenuState
}

func NewContextMenu() *ContextMenu {
	return &ContextMenu{
		Items:     nil,
		MenuState: component.MenuState{},
	}
}

func (m *ContextMenu) AddItem(item ContextMenuItem) {
	menuItem := component.MenuItem(th.Theme, &item.Clickable, item.Title)
	menuItem.Icon = item.Icon
	m.Options = append(m.Options, func(gtx C) D {
		return menuItem.Layout(gtx)
	})
	if item.AppendDivider {
		m.Options = append(m.Options, component.Divider(th.Theme).Layout)
	}
	m.Items = append(m.Items, &item)
}

func (m *ContextMenu) OnClicked(gtx C) {
	for _, item := range m.Items {
		if item.Clicked(gtx) {
			if item.Do != nil {
				item.Do()
			}
		}
	}
}
