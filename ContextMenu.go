package ux

import (
	"gioui.org/widget"
	"gioui.org/x/component"
)

type ContextMenuItem struct {
	Title         string
	Icon          *widget.Icon
	Can           func() bool
	Do            func()
	AppendDivider bool
	widget.Clickable
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

func (m *ContextMenu) Clicked(gtx C) {
	for _, item := range m.Items {
		if item.Clicked(gtx) {
			if item.Do != nil {
				item.Do()
			}
		}
	}
}
