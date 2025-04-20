package ux

import (
	"fmt"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/ux/resources/icons"
	"image"
	"image/color"
	"strconv"
	"sync"

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
	Items []*ContextMenuItem
	component.ContextArea
	component.MenuState
	widget.List
	RowClicks       []widget.Clickable
	DrawRow         func(gtx layout.Context, index int) layout.Dimensions
	ClickedRowindex int
	sync.Once
}

func NewContextMenu(length int, drawRow func(gtx layout.Context, index int) layout.Dimensions) *ContextMenu {
	return &ContextMenu{
		Items:       nil,
		ContextArea: component.ContextArea{},
		MenuState:   component.MenuState{},
		List: widget.List{
			Scrollbar: widget.Scrollbar{},
			List: layout.List{
				Axis:        layout.Vertical,
				ScrollToEnd: false,
				Alignment:   0,
				Position:    layout.Position{},
			},
		},
		RowClicks: make([]widget.Clickable, length),
		DrawRow:   drawRow,
	}
}

func (m *ContextMenu) AddItem(item ContextMenuItem) {
	menuItem := component.MenuItem(th, &item.Clickable, item.Title)
	menuItem.Icon = item.Icon
	m.Options = append(m.Options, func(gtx layout.Context) layout.Dimensions {
		return menuItem.Layout(gtx)
	})
	if item.AppendDivider {
		m.Options = append(m.Options, component.Divider(th).Layout)
	}
	m.Items = append(m.Items, &item)
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

// 测试用例，现在不需要了
func (m *ContextMenu) drawRowDefault(gtx layout.Context, rowClick *widget.Clickable, index int) layout.Dimensions {
	m.Once.Do(func() {
		m.AddItem(ContextMenuItem{
			Title:         "Red",
			Icon:          nil,
			Can:           func() bool { return false },
			Do:            func() { mylog.Info(m.ClickedRowindex, "red item clicked") },
			AppendDivider: false,
			Clickable:     widget.Clickable{},
		})
		m.AddItem(ContextMenuItem{
			Title:         "Green",
			Icon:          nil,
			Can:           func() bool { return false },
			Do:            func() { mylog.Info(m.ClickedRowindex, "Green item clicked") },
			AppendDivider: false,
			Clickable:     widget.Clickable{},
		})
		m.AddItem(ContextMenuItem{
			Title:         "Blue",
			Icon:          nil,
			Can:           func() bool { return false },
			Do:            func() { mylog.Info(m.ClickedRowindex, "Blue item clicked") },
			AppendDivider: false,
			Clickable:     widget.Clickable{},
		})
		m.AddItem(ContextMenuItem{
			Title:         "Balance",
			Icon:          icons.ActionAccountBalanceIcon,
			Can:           func() bool { return false },
			Do:            func() { mylog.Info(m.ClickedRowindex, "Balance item clicked") },
			AppendDivider: false,
			Clickable:     widget.Clickable{},
		})
		m.AddItem(ContextMenuItem{
			Title:         "Account",
			Icon:          icons.ActionAccountBoxIcon,
			Can:           func() bool { return false },
			Do:            func() { mylog.Info(m.ClickedRowindex, "Account item clicked") },
			AppendDivider: false,
			Clickable:     widget.Clickable{},
		})
		m.AddItem(ContextMenuItem{
			Title:         "Cart",
			Icon:          icons.ActionAddShoppingCartIcon,
			Can:           func() bool { return false },
			Do:            func() { mylog.Info(m.ClickedRowindex, "Cart item clicked") },
			AppendDivider: false,
			Clickable:     widget.Clickable{},
		})
	})
	if event, b := gtx.Event(pointer.Filter{Target: rowClick, Kinds: pointer.Press | pointer.Release}); b {
		if e, ok := event.(pointer.Event); ok {
			if e.Kind == pointer.Press {
				switch {
				case e.Buttons.Contain(pointer.ButtonPrimary):
					println("Row selected (left click) " + strconv.Itoa(index))
				case e.Buttons.Contain(pointer.ButtonSecondary):
					println("Row selected (right click)" + strconv.Itoa(index))
				}
			}
		}
	}
	buttonStyle := material.Button(th, rowClick, "item"+fmt.Sprintf("%d", index))
	buttonStyle.Color = RowColor(index)
	return buttonStyle.Layout(gtx)
}

// LayoutRow 对先行表格不感兴趣，不过在demo/other/hashicon/example/main.go下已经通过测试，不会去维护官方那个表格，不好用
// 官方表格控件已经有水平和垂直滚动条,所以不使用list布局，兼容一下
//func (m *ContextMenu) LayoutRow(gtx layout.Context, index int) layout.Dimensions {
//	return layout.Stack{}.Layout(gtx,
//		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
//			rowClick := &m.RowClicks[index]
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
//				return m.DrawRow(gtx, index)
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

// Layout 线性的list，表格以及非线性的树形表格(核心:直接rootRow当做线性表格即可，顶层调用menu布局。转不过弯来一直去处理容器节点是否渲染menu布局的问题)均通过测试
func (m *ContextMenu) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return material.List(th, &m.List).Layout(gtx, len(m.RowClicks), func(gtx layout.Context, index int) layout.Dimensions {
				rowClick := &m.RowClicks[index]
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
					if m.DrawRow == nil {
						return m.drawRowDefault(gtx, rowClick, index)
					}
					return m.DrawRow(gtx, index)
				})
			})
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return m.ContextArea.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min = image.Point{}
				m.OnClicked(gtx)
				return m.drawContextArea(gtx)
				return component.Menu(th, &m.MenuState).Layout(gtx) // 所有行的item共用一个popup菜单而不是每行popup一个
			})
		}),
	)
}

func (m *ContextMenu) drawContextArea(gtx layout.Context) layout.Dimensions { // popup区域的背景色，位置，四角弧度
	menuStyle := component.Menu(th, &m.MenuState)
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
