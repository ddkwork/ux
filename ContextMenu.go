package ux

import (
	"image"
	"image/color"
	"strconv"
	"sync"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/ux/resources/icons"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/ddkwork/ux/widget/material"
	"github.com/ddkwork/ux/x/component"
)

type ContextMenuItem struct {
	Title            string      // 菜单项标题
	Icon             any         // 可选的图标
	Can              func() bool // 是否绘制取决于当前渲染的行，回调内需要传递当前渲染的节点给回调，说白了这里是绘制条件，下面的do是业务逻辑，回调内传入的形参节点不一样
	Do               func()      // 调用被选中节点来操作业务逻辑
	AppendDivider    bool        // 是否添加分割线
	widget.Clickable             // 可点击的控件
}

type ContextMenu struct {
	Items           []*ContextMenuItem
	area            component.ContextArea
	state           component.MenuState
	list            widget.List
	rowClicks       []widget.Clickable
	rootRows        []layout.Widget
	ClickedRowindex int
	sync.Once
}

func (m *ContextMenu) AppendRootRows(rootRow layout.Widget) {
	m.rootRows = append(m.rootRows, rootRow)
}

func NewContextMenuWithRootRows(rootRows ...layout.Widget) *ContextMenu {
	m := NewContextMenu()
	m.rootRows = rootRows
	return m
}

func NewContextMenu() *ContextMenu {
	return &ContextMenu{
		Items: nil,
		area:  component.ContextArea{},
		state: component.MenuState{},
		list: widget.List{
			Scrollbar: widget.Scrollbar{},
			List: layout.List{
				Axis:        layout.Vertical,
				ScrollToEnd: false,
				Alignment:   0,
				Position:    layout.Position{},
			},
		},
		rowClicks:       nil,
		rootRows:        make([]layout.Widget, 0),
		ClickedRowindex: 0,
		Once:            sync.Once{},
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

func (m *ContextMenu) InitMenuItems(items ...ContextMenuItem) {
	m.Once.Do(func() {
		for _, item := range items {
			if item.Can() {
				m.AddItem(item)
			}
		}
	})
}

func (m *ContextMenu) OnClicked(gtx layout.Context) {
	for _, item := range m.Items {
		if item.Clicked(gtx) {
			if item.Do != nil {
				item.Do()
			}
		}
	}
}

// Layout 线性的list，表格以及非线性的树形表格(核心:直接rootRow当做线性表格即可，顶层调用menu布局。转不过弯来一直去处理容器节点是否渲染menu布局的问题)均通过测试
func (m *ContextMenu) Layout(gtx layout.Context) layout.Dimensions {
	mylog.CheckNil(m.rootRows)
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return material.List(th, &m.list).Layout(gtx, len(m.rootRows), func(gtx layout.Context, index int) layout.Dimensions {
				if m.rowClicks == nil {
					m.rowClicks = make([]widget.Clickable, len(m.rootRows))
				}
				if len(m.rowClicks) != len(m.rootRows) {
					m.rowClicks = make([]widget.Clickable, len(m.rootRows))
					// mylog.Warning("remake row clicks")//todo 树形表格调用这个太频繁，得查一下原因
				}
				rowClick := &m.rowClicks[index]
				return material.Clickable(gtx, rowClick, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					if event, b := gtx.Event(pointer.Filter{Target: rowClick, Kinds: pointer.Press | pointer.Release}); b {
						if e, ok := event.(pointer.Event); ok {
							if e.Kind == pointer.Press {
								switch {
								case e.Buttons.Contain(pointer.ButtonPrimary):
									m.ClickedRowindex = index // todo 移除树形表格的这个意思的字段?
									println("Row selected (left click) " + strconv.Itoa(index))
								case e.Buttons.Contain(pointer.ButtonSecondary):
									m.ClickedRowindex = index
									println("Row selected (right click)" + strconv.Itoa(index))
								}
							}
						}
					}
					return m.rootRows[index](gtx)
				})
			})
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return m.area.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min = image.Point{}
				m.OnClicked(gtx)
				return m.drawContextArea(gtx)
				return component.Menu(th, &m.state).Layout(gtx) // 所有行的item共用一个popup菜单而不是每行popup一个
			})
		}),
	)
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

func (m *ContextMenu) LayoutTest(gtx layout.Context) layout.Dimensions {
	m.InitMenuItems(
		ContextMenuItem{
			Title:         "Red",
			Icon:          nil,
			Can:           func() bool { return true },
			Do:            func() { mylog.Info(m.ClickedRowindex, "red item clicked") },
			AppendDivider: false,
			Clickable:     widget.Clickable{},
		},
		ContextMenuItem{
			Title:         "Green",
			Icon:          nil,
			Can:           func() bool { return true },
			Do:            func() { mylog.Info(m.ClickedRowindex, "Green item clicked") },
			AppendDivider: false,
			Clickable:     widget.Clickable{},
		},
		ContextMenuItem{
			Title:         "Blue",
			Icon:          nil,
			Can:           func() bool { return true },
			Do:            func() { mylog.Info(m.ClickedRowindex, "Blue item clicked") },
			AppendDivider: false,
			Clickable:     widget.Clickable{},
		},
		ContextMenuItem{
			Title:         "Balance",
			Icon:          icons.ActionAccountBalanceIcon,
			Can:           func() bool { return true },
			Do:            func() { mylog.Info(m.ClickedRowindex, "Balance item clicked") },
			AppendDivider: false,
			Clickable:     widget.Clickable{},
		},
		ContextMenuItem{
			Title:         "Account",
			Icon:          icons.ActionAccountBoxIcon,
			Can:           func() bool { return true },
			Do:            func() { mylog.Info(m.ClickedRowindex, "Account item clicked") },
			AppendDivider: false,
			Clickable:     widget.Clickable{},
		},
		ContextMenuItem{
			Title:         "Cart",
			Icon:          icons.ActionAddShoppingCartIcon,
			Can:           func() bool { return true },
			Do:            func() { mylog.Info(m.ClickedRowindex, "Cart item clicked") },
			AppendDivider: false,
			Clickable:     widget.Clickable{},
		},
	)
	return m.Layout(gtx)
}

// LayoutRow 对先行表格不感兴趣，不过在demo/other/hashicon/example/main.go下已经通过测试，不会去维护官方那个表格，不好用
// 官方表格控件已经有水平和垂直滚动条,所以不使用list布局，兼容一下
//func (m *ContextMenu) LayoutRow(gtx layout.Context, index int) layout.Dimensions {
//	return layout.Stack{}.Layout(gtx,
//		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
//			rowClick := &m.rowClicks[index]
//			return material.Clickable(gtx, rowClick, func(gtx layout.Context) layout.Dimensions {
//				gtx.Constraints.Min.X = gtx.Constraints.Max.X
//				if event, b := gtx.Event(pointer.Filter{Target: rowClick, Kinds: pointer.Press | pointer.Release}); b {
//					if e, ok := event.(pointer.Event); ok {
//						if e.Kind == pointer.Press {
//							switch {
//							case e.Buttons.Contain(pointer.ButtonPrimary):
//								m.ClickedRowindex = index
//								println("Row selected (left click) " + strconv.Itoa(index))
//							case e.Buttons.Contain(pointer.ButtonSecondary):
//								m.ClickedRowindex = index
//								println("Row selected (right click)" + strconv.Itoa(index))
//							}
//						}
//					}
//				}
//				return m.DrawRowCallback(gtx, index)
//			})
//		}),
//		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
//			return m.ContextArea.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
//				gtx.Constraints.Min = image.Point{}
//				m.OnClicked(gtx)
//				// return m.drawContextArea(gtx, th)
//				return component.Menu(th, &m.MenuState).Layout(gtx) // 所有行的item共用一个popup菜单而不是每行popup一个
//			})
//		}),
//	)
//}
