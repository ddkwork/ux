package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// ContextMenu 封装右键菜单的所有功能
type ContextMenu struct {
	Visible bool
	Pos     image.Point
	Hovered int // 当前悬停的菜单项索引

	Items []MenuItem // 菜单项列表
}

// MenuItem 定义菜单项
type MenuItem struct {
	Text   string
	Action func()
}

// Layout 封装右键事件的监听和菜单显示
func (c *ContextMenu) Layout(gtx layout.Context, th *material.Theme, content layout.Widget) layout.Dimensions {
	// 1. 首先布置内容区域
	dims := content(gtx)

	// 2. 监听右键事件
	c.listenRightClick(gtx)

	// 3. 如果菜单可见，则显示菜单
	if c.Visible {
		c.drawMenu(gtx, th)
	}

	return dims
}

// 监听右键点击事件
func (c *ContextMenu) listenRightClick(gtx layout.Context) {
	// 注册右键事件
	defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()

	// 注册事件过滤器
	event.Op(gtx.Ops, pointer.Filter{
		Target: c,             // 目标为ContextMenu
		Kinds:  pointer.Press, // 只关心按下事件
	})

	// 处理事件
	for {
		ev, ok := gtx.Event(pointer.Filter{Target: c})
		if !ok {
			break
		}

		if e, ok := ev.(pointer.Event); ok {
			if e.Buttons == pointer.ButtonSecondary && e.Kind == pointer.Press {
				c.Visible = true
				c.Pos = image.Point{int(e.Position.X), int(e.Position.Y)}
				c.Hovered = -1
			}
		}
	}
}

// 监听全局点击关闭菜单事件
func (c *ContextMenu) listenGlobalClick(gtx layout.Context) {
	// 设置全局点击区域
	defer clip.Rect{
		Min: image.Point{-10000, -10000},
		Max: image.Point{10000, 10000},
	}.Push(gtx.Ops).Pop()

	// 注册事件过滤器
	event.Op(gtx.Ops, pointer.Filter{
		Target: c, // 目标为ContextMenu
		Kinds:  pointer.Press,
	})

	// 处理事件
	for {
		ev, ok := gtx.Event(pointer.Filter{Target: c})
		if !ok {
			break
		}

		if e, ok := ev.(pointer.Event); ok {
			if e.Kind == pointer.Press {
				// 计算菜单区域
				menuRect := image.Rectangle{
					Min: c.Pos,
					Max: c.Pos.Add(image.Point{gtx.Dp(150), len(c.Items) * gtx.Dp(32)}),
				}

				// 如果点击在菜单外，则关闭菜单
				if e.Position.X > float32(menuRect.Max.X) {
					c.Visible = false
				}
			}
		}
	}
}

// 绘制菜单
func (c *ContextMenu) drawMenu(gtx layout.Context, th *material.Theme) {
	const (
		menuWidth  = unit.Dp(150)
		itemHeight = unit.Dp(32)
	)

	// 设置菜单位置
	defer op.Offset(c.Pos).Push(gtx.Ops).Pop()

	// 创建菜单约束
	gtx.Constraints.Min.X = gtx.Dp(menuWidth)
	gtx.Constraints.Min.Y = len(c.Items) * gtx.Dp(itemHeight)
	gtx.Constraints.Max = gtx.Constraints.Min

	// 绘制阴影
	shadowRect := image.Rectangle{Max: gtx.Constraints.Min.Add(image.Pt(4, 4))}
	paint.FillShape(gtx.Ops, color.NRGBA{A: 100}, clip.Rect(shadowRect).Op())

	// 绘制菜单背景
	menuArea := clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops)
	paint.Fill(gtx.Ops, color.NRGBA{R: 250, G: 250, B: 250, A: 255})
	menuArea.Pop()

	// 绘制边框
	border := clip.Stroke{
		Path:  clip.Rect{Max: gtx.Constraints.Min}.Path(),
		Width: 1,
	}.Op().Push(gtx.Ops)
	paint.Fill(gtx.Ops, color.NRGBA{A: 50})
	border.Pop()

	// 监听全局点击关闭菜单
	c.listenGlobalClick(gtx)

	// 绘制菜单项
	for i := range c.Items {
		c.drawMenuItem(gtx, th, i)
	}
}

// 绘制单个菜单项
func (c *ContextMenu) drawMenuItem(gtx layout.Context, th *material.Theme, index int) {
	item := &c.Items[index]

	// 设置菜单项位置
	itemTop := index * gtx.Dp(unit.Dp(32))
	defer op.Offset(image.Pt(0, itemTop)).Push(gtx.Ops).Pop()

	// 设置菜单项区域
	itemArea := clip.Rect{Max: image.Pt(gtx.Constraints.Min.X, gtx.Dp(32))}.Push(gtx.Ops)

	// 注册菜单项事件
	event.Op(gtx.Ops, pointer.Filter{
		Target: item, // 使用菜单项作为事件目标
		Kinds:  pointer.Enter | pointer.Leave | pointer.Press,
	})

	// 处理菜单项事件
	for {
		ev, ok := gtx.Event(pointer.Filter{Target: item})
		if !ok {
			break
		}

		if e, ok := ev.(pointer.Event); ok {
			switch e.Kind {
			case pointer.Enter:
				c.Hovered = index
			case pointer.Leave:
				if c.Hovered == index {
					c.Hovered = -1
				}
			case pointer.Press:
				if e.Buttons == pointer.ButtonPrimary {
					item.Action()
					c.Visible = false
				}
			}
		}
	}

	// 绘制悬停背景
	if index == c.Hovered {
		paint.Fill(gtx.Ops, color.NRGBA{R: 240, G: 240, B: 255, A: 255})
	}

	// 绘制菜单项文本
	label := material.Label(th, unit.Sp(14), item.Text)
	label.Color = color.NRGBA{R: 20, G: 20, B: 20}
	layout.UniformInset(unit.Dp(10)).Layout(gtx, label.Layout)

	itemArea.Pop()
}

func main() {
	go func() {
		// 创建窗口
		w := new(app.Window)
		w.Option(
			app.Title("右键菜单演示"),
			app.Size(unit.Dp(400), unit.Dp(300)),
		)

		// 初始化主题
		th := material.NewTheme()
		//th.Shaper = font.NewShaper(gofont.Collection())

		// 创建应用
		appInstance := &MyApp{
			Theme: th,
			Menu: &ContextMenu{
				Items: []MenuItem{
					{Text: "复制", Action: func() { log.Println("复制操作") }},
					{Text: "粘贴", Action: func() { log.Println("粘贴操作") }},
					{Text: "删除", Action: func() { log.Println("删除操作") }},
				},
			},
		}

		// 主事件循环
		var ops op.Ops
		for {
			e := w.Event()
			switch e := e.(type) {
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				appInstance.Layout(gtx)
				e.Frame(gtx.Ops)
				ops.Reset()
			case app.DestroyEvent:
				log.Println("窗口关闭")
				os.Exit(0)
				return
			}
		}
	}()

	// 启动应用
	app.Main()
}

// 重命名为 MyApp 以避免与 app 包冲突
type MyApp struct {
	Theme *material.Theme
	Menu  *ContextMenu
}

func (a *MyApp) Layout(gtx layout.Context) layout.Dimensions {
	// 使用封装好的右键菜单组件包裹内容
	return a.Menu.Layout(gtx, a.Theme, func(gtx layout.Context) layout.Dimensions {
		// 绘制背景
		paint.Fill(gtx.Ops, color.NRGBA{R: 240, G: 240, B: 240, A: 255})

		// 注册右键事件
		defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
		event.Op(gtx.Ops, pointer.Filter{
			Target: a.Menu,        // 目标为ContextMenu
			Kinds:  pointer.Press, // 只关心按下事件
		})

		// 添加说明文本
		label := material.Label(a.Theme, unit.Sp(18), "在任意位置点击右键显示菜单")
		layout.Center.Layout(gtx, label.Layout)

		return layout.Dimensions{Size: gtx.Constraints.Max}
	})
}
