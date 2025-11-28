package ux

import (
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/ux/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/font"
	"gioui.org/io/pointer"
)

// Segment 表示分段按钮的单个选项
type Segment struct {
	Text      string
	Icon      *widget.Icon // 可选图标
	Selected  bool
	Clickable widget.Clickable
}

// SegmentedButton 表示 Material Design 3 风格的分段按钮
type SegmentedButton struct {
	segments    []Segment
	multiSelect bool

	// 视觉样式
	selectedColor       color.NRGBA
	unselectedColor     color.NRGBA
	selectedTextColor   color.NRGBA
	unselectedTextColor color.NRGBA
	borderColor         color.NRGBA
	pressedColor        color.NRGBA

	// 内部状态
	pressedIndex int // 当前被按下的按钮索引
}

// NewSegmentedButton 创建新的分段按钮
func NewSegmentedButton() *SegmentedButton {
	return &SegmentedButton{
		segments:     make([]Segment, 0),
		multiSelect:  false,
		pressedIndex: -1,

		// Material Design 3 颜色
		selectedColor:       color.NRGBA{R: 103, G: 80, B: 164, A: 255},  // Primary container
		unselectedColor:     color.NRGBA{R: 232, G: 222, B: 248, A: 255}, // Surface variant
		selectedTextColor:   color.NRGBA{R: 103, G: 80, B: 164, A: 255},  // On primary container
		unselectedTextColor: color.NRGBA{R: 73, G: 69, B: 79, A: 255},    // On surface variant
		borderColor:         color.NRGBA{R: 121, G: 116, B: 126, A: 255}, // Outline
		pressedColor:        color.NRGBA{R: 103, G: 80, B: 164, A: 102},  // Primary with opacity for pressed state
	}
}

// AddSegment 添加分段选项
func (sb *SegmentedButton) AddSegment(text string, icon *widget.Icon) {
	sb.segments = append(sb.segments, Segment{
		Text:     text,
		Icon:     icon,
		Selected: false,
	})
}

// SetSelected 设置选中的选项
func (sb *SegmentedButton) SetSelected(index int) {
	if index >= 0 && index < len(sb.segments) {
		if !sb.multiSelect {
			// 单选模式：取消其他选择
			for i := range sb.segments {
				sb.segments[i].Selected = false
			}
		}
		sb.segments[index].Selected = !sb.segments[index].Selected
	}
}

// GetSelected 获取选中的索引（单选模式）或所有选中的索引（多选模式）
func (sb *SegmentedButton) GetSelected() []int {
	var selected []int
	for i, segment := range sb.segments {
		if segment.Selected {
			selected = append(selected, i)
		}
	}
	return selected
}

// SetMultiSelect 设置多选模式
func (sb *SegmentedButton) SetMultiSelect(multi bool) {
	sb.multiSelect = multi
}

// Layout 渲染分段按钮
func (sb *SegmentedButton) Layout(gtx layout.Context) layout.Dimensions {
	// 处理点击事件
	for i := range sb.segments {
		for sb.segments[i].Clickable.Clicked(gtx) {
			sb.SetSelected(i)
		}

		// 检测按下状态
		if sb.segments[i].Clickable.Pressed() {
			sb.pressedIndex = i
		} else if sb.pressedIndex == i {
			sb.pressedIndex = -1
		}
	}

	// 创建弹性布局
	var children []layout.FlexChild

	for i := range sb.segments {
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sb.renderSegment(gtx, i)
		}))
	}

	return layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceSides,
	}.Layout(gtx, children...)
}

// renderSegment 渲染单个分段
func (sb *SegmentedButton) renderSegment(gtx layout.Context, index int) layout.Dimensions {
	segment := &sb.segments[index]

	// 根据选择状态确定颜色
	bgColor := sb.unselectedColor
	textColor := sb.unselectedTextColor

	if segment.Selected {
		bgColor = sb.selectedColor
		textColor = sb.selectedTextColor
	}

	// 最小尺寸约束
	minSize := gtx.Constraints.Max
	minSize.Y = gtx.Dp(40) // 最小高度

	return layout.Stack{}.Layout(gtx,
		// 背景层
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// 绘制背景 - 改为直角矩形
			defer clip.Rect{
				Max: gtx.Constraints.Min,
			}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, bgColor)
			paint.Fill(gtx.Ops, bgColor)

			// 绘制按下状态覆盖层
			if sb.pressedIndex == index {
				paint.Fill(gtx.Ops, sb.pressedColor)
			}

			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),

		// 内容层
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
					Spacing:   layout.SpaceEnd,
				}.Layout(gtx,
					// 图标（如果有）
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if segment.Icon != nil {
							//iconSize := gtx.Dp(18)

							return segment.Icon.Layout(gtx, textColor)
						}
						return layout.Dimensions{}
					}),

					// 文本
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if segment.Text != "" {
							label := material.Label(th, unit.Sp(14), segment.Text)
							label.Color = textColor
							label.Font.Weight = font.Medium
							if segment.Icon != nil {
								// 如果有图标，添加间距
								return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, label.Layout)
							}
							return label.Layout(gtx)
						}
						return layout.Dimensions{}
					}),
				)
			})
		}),

		// 点击区域层
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			// 扩大点击区域
			defer pointer.PassOp{}.Push(gtx.Ops).Pop()
			return segment.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: gtx.Constraints.Min}
			})
		}),
	)
}

// 使用示例
func testNewSegmentedButton() {
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("Material Design 3 Segmented Button"),
			app.Size(unit.Dp(500), unit.Dp(300)),
		)

		segButton := NewSegmentedButton()

		// 添加带图标的选项（需要先创建图标）
		// 注意：这里使用占位符，实际使用时需要加载真实图标
		homeIcon := mylog.Check2(widget.NewIcon(icons.ActionBackup))
		workIcon := mylog.Check2(widget.NewIcon(icons.Action3DRotation))
		otherIcon := mylog.Check2(widget.NewIcon(icons.ActionAccountCircle))

		segButton.AddSegment("Home", homeIcon)
		segButton.AddSegment("Work", workIcon)
		segButton.AddSegment("Other", otherIcon)

		// 设置默认选中
		segButton.SetSelected(0)

		var ops op.Ops
		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				os.Exit(0)
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)

				// 在中心布局分段按钮
				layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Top: unit.Dp(20), Bottom: unit.Dp(20)}.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return segButton.Layout(gtx)
						})
				})

				e.Frame(gtx.Ops)
			}
		}
	}()

	app.Main()
}
