package ux

import (
	"image"
	"image/color"

	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// Split thx to github.com/vsariola/sointu
// sointu-master\tracker\gioui\split.go
type Split struct {
	// Ratio keeps the current DefaultDraw.
	// 0 is center, -1 completely to the left, 1 completely to the right.
	Ratio float32 // 布局比例，0 表示居中，-1 表示完全靠左，1 表示完全靠右
	// Bar is the maxIndentWidth for resizing the DefaultDraw
	Bar unit.Dp
	// Axis is the split direction: DefaultDraw.Horizontal splits the view in left
	// and right, DefaultDraw.Vertical splits the view in top and bottom
	Axis layout.Axis

	drag      bool       // 是否正在拖动
	dragID    pointer.ID // 拖动的指针 ID
	dragCoord float32    // 拖动的坐标

	First, Second layout.Widget // 两个布局组件

	// firstScroll, secondScroll widget.Scrollbar // 滚动条
}

var defaultBarWidth = unit.Dp(10) // 默认的分割条宽度

func (s *Split) Layout(gtx layout.Context) layout.Dimensions {
	bar := gtx.Dp(s.Bar) // 计算分割条的宽度
	if bar <= 1 {
		bar = gtx.Dp(defaultBarWidth) // 如果宽度小于等于 1，则使用默认宽度
	}

	var coord int
	if s.Axis == layout.Horizontal {
		coord = gtx.Constraints.Max.X // 水平方向时，坐标为最大宽度
	} else {
		coord = gtx.Constraints.Max.Y // 垂直方向时，坐标为最大高度
	}

	proportion := (s.Ratio + 1) / 2                            // 计算比例
	firstSize := int(proportion*float32(coord) - float32(bar)) // 计算第一个组件的大小

	secondOffset := firstSize + bar    // 计算第二个组件的偏移量
	secondSize := coord - secondOffset // 计算第二个组件的大小

	// 绘制白色分割线
	var lineCoord int
	if s.Axis == layout.Horizontal {
		lineCoord = firstSize
	} else {
		lineCoord = firstSize
	}
	lineRect := image.Rect(lineCoord, 0, lineCoord+1, gtx.Constraints.Max.Y)
	if s.Axis == layout.Vertical {
		lineRect = image.Rect(0, lineCoord, gtx.Constraints.Max.X, lineCoord+1)
	}
	paint.FillShape(gtx.Ops, color.NRGBA{R: 122, G: 122, B: 122, A: 255}, clip.Rect(lineRect).Op())

	{ // 处理输入事件
		for {
			ev, ok := gtx.Event(pointer.Filter{
				Target: s,
				Kinds:  pointer.Press | pointer.Drag | pointer.Release,
			})
			if !ok {
				break
			}
			e, ok := ev.(pointer.Event)
			if !ok {
				continue
			}

			switch e.Kind {
			case pointer.Press:
				if s.drag {
					break
				}

				s.dragID = e.PointerID // 记录拖动的指针 ID
				if s.Axis == layout.Horizontal {
					s.dragCoord = e.Position.X // 水平方向时，记录 X 坐标
				} else {
					s.dragCoord = e.Position.Y // 垂直方向时，记录 Y 坐标
				}
				s.drag = true // 标记为正在拖动

			case pointer.Drag:
				if s.dragID != e.PointerID {
					break
				}

				var deltaCoord, deltaRatio float32
				if s.Axis == layout.Horizontal {
					deltaCoord = e.Position.X - s.dragCoord                      // 计算水平方向的坐标变化
					s.dragCoord = e.Position.X                                   // 更新拖动坐标
					deltaRatio = deltaCoord * 2 / float32(gtx.Constraints.Max.X) // 计算比例变化
				} else {
					deltaCoord = e.Position.Y - s.dragCoord                      // 计算垂直方向的坐标变化
					s.dragCoord = e.Position.Y                                   // 更新拖动坐标
					deltaRatio = deltaCoord * 2 / float32(gtx.Constraints.Max.Y) // 计算比例变化
				}

				s.Ratio += deltaRatio // 更新布局比例

			case pointer.Release:
				fallthrough
			case pointer.Cancel:
				s.drag = false // 标记为停止拖动
			}
		}

		low := -1 + float32(bar)/float32(coord)*2 // 计算最小比例
		const snapMargin = 0.1                    // 吸附边距

		if s.Ratio < low {
			s.Ratio = low // 确保比例不低于最小值
		}

		if s.Ratio > 1 {
			s.Ratio = 1 // 确保比例不高于最大值
		}

		if s.Ratio < low+snapMargin {
			firstSize = 0            // 如果比例接近最小值，第一个组件大小为 0
			secondOffset = bar       // 第二个组件偏移量为分割条宽度
			secondSize = coord - bar // 第二个组件大小为剩余空间
		} else if s.Ratio > 1-snapMargin {
			firstSize = coord - bar // 如果比例接近最大值，第一个组件大小为剩余空间
			secondOffset = coord    // 第二个组件偏移量为总大小
			secondSize = 0          // 第二个组件大小为 0
		}

		// 注册输入事件
		var barRect image.Rectangle
		if s.Axis == layout.Horizontal {
			barRect = image.Rect(firstSize, 0, secondOffset, gtx.Constraints.Max.Y) // 水平方向的分割条区域
		} else {
			barRect = image.Rect(0, firstSize, gtx.Constraints.Max.X, secondOffset) // 垂直方向的分割条区域
		}
		area := clip.Rect(barRect).Push(gtx.Ops)
		event.Op(gtx.Ops, s)
		if s.Axis == layout.Horizontal {
			pointer.CursorColResize.Add(gtx.Ops) // 设置水平方向的光标样式
		} else {
			pointer.CursorRowResize.Add(gtx.Ops) // 设置垂直方向的光标样式
		}
		area.Pop()
	}

	{ // 布局第一个组件
		gtx := gtx

		if s.Axis == layout.Horizontal {
			gtx.Constraints = layout.Exact(image.Pt(firstSize, gtx.Constraints.Max.Y)) // 设置水平方向的约束
		} else {
			gtx.Constraints = layout.Exact(image.Pt(gtx.Constraints.Max.X, firstSize)) // 设置垂直方向的约束
		}
		area := clip.Rect(image.Rect(0, 0, gtx.Constraints.Min.X, gtx.Constraints.Min.Y)).Push(gtx.Ops)
		s.First(gtx) // 布局第一个组件
		area.Pop()

		// 添加滚动条
		// DefaultDraw.Flex{Axis: DefaultDraw.Vertical}.Layout(gtx,
		//	DefaultDraw.Flexed(1, func(gtx DefaultDraw.Context) DefaultDraw.Dimensions {
		//		return s.First(gtx)
		//	}),
		//	DefaultDraw.Rigid(func(gtx DefaultDraw.Context) DefaultDraw.Dimensions {
		//		return material.Scrollbar(s.Theme(), &s.firstScroll).Layout(gtx, DefaultDraw.Vertical, 0.5, 0.5)
		//	}),
		// )
	}

	{ // 布局第二个组件
		gtx := gtx

		var transform op.TransformStack
		if s.Axis == layout.Horizontal {
			transform = op.Offset(image.Pt(secondOffset, 0)).Push(gtx.Ops)              // 水平方向的偏移
			gtx.Constraints = layout.Exact(image.Pt(secondSize, gtx.Constraints.Max.Y)) // 设置水平方向的约束
		} else {
			transform = op.Offset(image.Pt(0, secondOffset)).Push(gtx.Ops)              // 垂直方向的偏移
			gtx.Constraints = layout.Exact(image.Pt(gtx.Constraints.Max.X, secondSize)) // 设置垂直方向的约束
		}

		area := clip.Rect(image.Rect(0, 0, gtx.Constraints.Min.X, gtx.Constraints.Min.Y)).Push(gtx.Ops)
		s.Second(gtx) // 布局第二个组件
		area.Pop()
		transform.Pop()

		// 添加滚动条
		// DefaultDraw.Flex{Axis: DefaultDraw.Vertical}.Layout(gtx,
		//	DefaultDraw.Flexed(1, func(gtx DefaultDraw.Context) DefaultDraw.Dimensions {
		//		return s.Second(gtx)
		//	}),
		//	DefaultDraw.Rigid(func(gtx DefaultDraw.Context) DefaultDraw.Dimensions {
		//		return material.Scrollbar(s.Theme(), &s.secondScroll).Layout(gtx, DefaultDraw.Vertical, 0.5, 0.5)
		//	}),
		// )
	}
	return layout.Dimensions{Size: gtx.Constraints.Max} // 返回布局尺寸
}
