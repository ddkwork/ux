package ux

import (
	"fmt"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/ddkwork/ux/widget/material"
	"github.com/ddkwork/ux/x/component"
	"image"
	"image/color"
	"strconv"
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
	m.Options = append(m.Options, func(gtx C) D {
		return menuItem.Layout(gtx)
	})
	if item.AppendDivider {
		m.Options = append(m.Options, component.Divider(th).Layout)
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

func (m *ContextMenu) drawRowDefault(gtx layout.Context, rowClick *widget.Clickable, index int) layout.Dimensions {
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

func (m *ContextMenu) LayoutRow(gtx layout.Context, index int) layout.Dimensions { //官方表格控件已经有水平和垂直滚动条,所以不使用list布局，兼用一下
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			rowClick := &m.RowClicks[index]
			return material.Clickable(gtx, rowClick, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				if event, b := gtx.Event(pointer.Filter{Target: rowClick, Kinds: pointer.Press | pointer.Release}); b {
					if e, ok := event.(pointer.Event); ok {
						if e.Kind == pointer.Press {
							switch {
							case e.Buttons.Contain(pointer.ButtonPrimary):
								m.ClickedRowindex = index
								println("Row selected (left click) " + strconv.Itoa(index))
							case e.Buttons.Contain(pointer.ButtonSecondary):
								m.ClickedRowindex = index
								println("Row selected (right click)" + strconv.Itoa(index))
							}
						}
					}
				}
				return m.DrawRow(gtx, index)
			})
		}),
		layout.Expanded(func(gtx C) D {
			return m.ContextArea.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min = image.Point{}
				m.OnClicked(gtx)
				//return m.drawContextArea(gtx, th)
				return component.Menu(th, &m.MenuState).Layout(gtx) //所有行的item公用一个popup菜单而不是每行popup一个
			})
		}),
	)
}

func (m *ContextMenu) LayoutOld(gtx layout.Context) layout.Dimensions {
	return material.List(th, &m.List).Layout(gtx, len(m.RowClicks), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Stack{}.Layout(gtx,
			layout.Stacked(func(gtx C) D {
				rowClick := &m.RowClicks[index]
				return material.Clickable(gtx, rowClick, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					if event, b := gtx.Event(pointer.Filter{Target: rowClick, Kinds: pointer.Press | pointer.Release}); b {
						if e, ok := event.(pointer.Event); ok {
							if e.Kind == pointer.Press {
								switch {
								case e.Buttons.Contain(pointer.ButtonPrimary):
									m.ClickedRowindex = index
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
			}),
			layout.Expanded(func(gtx C) D {
				return m.ContextArea.Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min = image.Point{}
					m.OnClicked(gtx)
					return m.drawContextArea(gtx, th)
					return component.Menu(th, &m.MenuState).Layout(gtx) //所有行的item公用一个popup菜单而不是每行popup一个
				})
			}),
		)
	})
}

func (m *ContextMenu) Layout(gtx layout.Context) layout.Dimensions {
	//return m.LayoutOld(gtx)
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return material.List(th, &m.List).Layout(gtx, len(m.RowClicks), func(gtx layout.Context, index int) layout.Dimensions {
				rowClick := &m.RowClicks[index]
				return material.Clickable(gtx, rowClick, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					if event, b := gtx.Event(pointer.Filter{Target: rowClick, Kinds: pointer.Press | pointer.Release}); b {
						if e, ok := event.(pointer.Event); ok {
							if e.Kind == pointer.Press {
								switch {
								case e.Buttons.Contain(pointer.ButtonPrimary):
									m.ClickedRowindex = index
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
		layout.Expanded(func(gtx C) D {
			return m.ContextArea.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min = image.Point{}
				m.OnClicked(gtx)
				return m.drawContextArea(gtx, th)
				return component.Menu(th, &m.MenuState).Layout(gtx) //所有行的item公用一个popup菜单而不是每行popup一个
			})
		}),
	)
}

func (m *ContextMenu) drawContextArea(gtx C, th *material.Theme) D { //popup区域的背景色，位置，四角弧度
	menuStyle := component.Menu(th, &m.MenuState)
	menuStyle.SurfaceStyle = component.SurfaceStyle{
		Theme: th,
		ShadowStyle: component.ShadowStyle{
			CornerRadius: 18,
			//Elevation:     0,//todo test elevation
			//AmbientColor:  color.NRGBA{},//todo test ambient color
			//PenumbraColor: color.NRGBA{},
			//UmbraColor:    color.NRGBA{},
		},
		Fill: color.NRGBA{R: 50, G: 50, B: 50, A: 255},
	}
	return menuStyle.Layout(gtx)
}
