package ux

import (
	"image"
	"image/color"
	"sync"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/ux/resources/images"
	"github.com/ddkwork/ux/widget/material"
	"github.com/ddkwork/ux/x/component"
)

type ContextMenuItem struct {
	Title            string      // 菜单项标题
	Icon             []byte      // 可选的图标
	Can              func() bool // 是否绘制取决于当前渲染的行，回调内需要传递当前渲染的节点给回调，说白了这里是绘制条件，下面的do是业务逻辑，回调内传入的形参节点不一样
	Do               func()      // 调用被选中节点来操作业务逻辑
	AppendDivider    bool        // 是否添加分割线
	widget.Clickable             // 可点击的控件
}

type ContextMenu struct {
	Items []*ContextMenuItem
	area  *component.ContextArea
	state component.MenuState
	list  *widget.List
	sync.Once
}

func NewContextMenu() *ContextMenu {
	return &ContextMenu{
		Items: nil,
		area: &component.ContextArea{
			Activation:       pointer.ButtonSecondary,
			AbsolutePosition: false,
			PositionHint:     0,
		},
		state: component.MenuState{},
		list: &widget.List{
			Scrollbar: widget.Scrollbar{},
			List: layout.List{
				Axis:        layout.Vertical,
				ScrollToEnd: false,
				Alignment:   layout.Middle,
				// ScrollAnyAxis: true,
				Position: layout.Position{
					BeforeEnd:  false,
					First:      0,
					Offset:     0,
					OffsetLast: 0,
					Count:      0,
					Length:     0,
				},
			},
		},
		Once: sync.Once{},
	}
}

func (m *ContextMenu) AddItem(item ContextMenuItem) {
	menuItem := component.MenuItem(th, &item.Clickable, item.Title)
	menuItem.Icon = item.Icon
	m.state.Options = append(m.state.Options, func(gtx layout.Context) layout.Dimensions {
		return menuItem.Layout(gtx)
	})
	if item.AppendDivider {
		m.state.Options = append(m.state.Options, component.Divider(th).Layout)
	}
	m.Items = append(m.Items, &item)
}

func (m *ContextMenu) onClicked(gtx layout.Context) {
	for _, item := range m.Items {
		if item.Clicked(gtx) {
			if item.Do != nil {
				item.Do()
			}
		}
	}
}

func (m *ContextMenu) Layout(gtx layout.Context, rootRows []layout.Widget) layout.Dimensions {
	if len(rootRows) == 0 { // mitmproxy start
		return layout.Dimensions{}
	}
	gtx.Values[""] = layout.Widget(func(gtx layout.Context) layout.Dimensions {
		return m.area.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = image.Point{}
			m.onClicked(gtx)
			return m.drawContextArea(gtx)
			return component.Menu(th, &m.state).Layout(gtx) // 所有行的item共用一个popup菜单而不是每行popup一个
		})
	})
	return material.List(th, m.list).Layout(gtx, len(rootRows), func(gtx layout.Context, index int) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return rootRows[index](gtx)
	})

	//return layout.Stack{}.Layout(gtx, layout.Stacked(func(gtx layout.Context) layout.Dimensions {
	//
	//}),
	//	layout.Expanded(func(gtx layout.Context) layout.Dimensions {
	//		// skip := false
	//		// for {
	//		// 	ev, ok := gtx.Event(pointer.Filter{
	//		// 		Target: &m.area,
	//		// 		Kinds:  pointer.Press,
	//		// 	})
	//		// 	if !ok {
	//		// 		break
	//		// 	}
	//		// 	e, ok := ev.(pointer.Event)
	//		// 	if !ok {
	//		// 		continue
	//		// 	}
	//		// 	if e.Buttons.Contain(pointer.ButtonPrimary) && e.Kind == pointer.Press {
	//		// 		skip = true
	//		// 	}
	//		// }
	//		// if skip {
	//		// 	return layout.Dimensions{}
	//		// }
	//		return m.area.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
	//			gtx.Constraints.Min = image.Point{}
	//			m.onClicked(gtx)
	//			return m.drawContextArea(gtx)
	//			return component.Menu(th, &m.state).Layout(gtx) // 所有行的item共用一个popup菜单而不是每行popup一个
	//		})
	//	})
	//)
}

func (m *ContextMenu) drawContextArea(gtx layout.Context) layout.Dimensions { // popup区域的背景色，位置，四角弧度
	menuStyle := component.Menu(th, &m.state)
	menuStyle.SurfaceStyle = component.SurfaceStyle{
		Theme: th,
		ShadowStyle: component.ShadowStyle{
			CornerRadius: 18,
			// Elevation:     0,//todo test elevation
			// AmbientColor:  color.NRGBA{},//todo test ambient color
			// PenumbraColor: color.NRGBA{},
			// UmbraColor:    color.NRGBA{},
		},
		Fill: color.NRGBA{R: 50, G: 50, B: 50, A: 255},
	}
	return menuStyle.Layout(gtx)
}

func (m *ContextMenu) LayoutTest(gtx layout.Context, rootRows []layout.Widget) layout.Dimensions {
	m.Once.Do(func() {
		m.AddItem(ContextMenuItem{
			Title:         "Red",
			Icon:          nil,
			Can:           func() bool { return true },
			Do:            func() { mylog.Info("red item clicked") },
			AppendDivider: false,
			Clickable:     widget.Clickable{},
		})
		m.AddItem(ContextMenuItem{
			Title:         "Green",
			Icon:          nil,
			Can:           func() bool { return true },
			Do:            func() { mylog.Info("Green item clicked") },
			AppendDivider: false,
			Clickable:     widget.Clickable{},
		})
		m.AddItem(ContextMenuItem{
			Title:         "Blue",
			Icon:          nil,
			Can:           func() bool { return true },
			Do:            func() { mylog.Info("Blue item clicked") },
			AppendDivider: false,
			Clickable:     widget.Clickable{},
		})
		m.AddItem(ContextMenuItem{
			Title:         "Balance",
			Icon:          images.ActionAccountBalanceIcon,
			Can:           func() bool { return true },
			Do:            func() { mylog.Info("Balance item clicked") },
			AppendDivider: false,
			Clickable:     widget.Clickable{},
		})
		m.AddItem(ContextMenuItem{
			Title:         "Account",
			Icon:          images.ActionAccountBoxIcon,
			Can:           func() bool { return true },
			Do:            func() { mylog.Info("Account item clicked") },
			AppendDivider: false,
			Clickable:     widget.Clickable{},
		})
		m.AddItem(ContextMenuItem{
			Title:         "Cart",
			Icon:          images.ActionAddShoppingCartIcon,
			Can:           func() bool { return true },
			Do:            func() { mylog.Info("Cart item clicked") },
			AppendDivider: false,
			Clickable:     widget.Clickable{},
		})
	})
	return m.Layout(gtx, rootRows)
}
