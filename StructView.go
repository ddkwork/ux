package ux

import (
	"image/color"
	"reflect"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/x/outlay"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/safemap"
	"github.com/ddkwork/ux/widget/material"
	"github.com/ddkwork/ux/x/component"
)

// 相当于StructView，formView，但是他是绑定结构体的
// 使用有序map来存储布局信息
//type FieldRow struct {
//	Key   string
//	Value *Input
//}

type StructView[T any] struct { // 其实就是一个标题row+list滚动多个row

	// 布局预期视觉样式：1个元素的一行
	Title string // 标题，label渲染，不进入list滚动，需要固定位置，字体大小是h4

	// select keygen需要一个下拉框，
	// 布局预期视觉样式：2个元素的一行，key右对齐，value左对齐

	// 布局预期视觉样式：2个元素的一行，key右对齐，value左对齐
	rows *safemap.M[string, *Input] //

	// 底部布局成一行，对于密码箱的rsa，这里需要额外的按钮，所以应该增加个outlay布局，
	// 填充到这一行的左侧，关闭和应用始终在右下角 layout.N
	// 这样如果不在list内滚动就不会渲染，不知道是什么原因，总感觉一个整体的widget只能有一个list，先滚动吧
	// 布局预期视觉样式：多个元素的一行
	applyBtn widget.Clickable
	closeBtn widget.Clickable
	onApply  func()

	unmarshal unmarshalFun // 反序列化函数，用于将输入框的值反序列化成结构体

	widget.List // 滚动所有行
	*component.ModalState
	Visible bool
}

type (
	marshalFun   func(any) []string
	unmarshalFun func([]string) any
)

func ReflectStruct2Map(object any, marshal marshalFun, unmarshal unmarshalFun) *safemap.M[string, string] {
	keys := reflect.VisibleFields(reflect.TypeOf(object))
	values := marshal(object)
	m := new(safemap.M[string, string])
	for i, key := range keys {
		if key.Tag.Get("table") == "-" || key.Tag.Get("json") == "-" {
			//	continue  //todo
		}
		if !key.IsExported() {
			mylog.Trace("field name is not exported: ", key.Name) // 用于树形表格序列化json保存到文件，没有导出则json会失败
			continue
		}
		m.Set(key.Name, values[i])
	}
	return m
}

func NewStructView[T any](title string, object T, marshal marshalFun, unmarshal unmarshalFun) *StructView[T] {
	visibleFields := ReflectStruct2Map(object, marshal, unmarshal)
	FieldRows := new(safemap.M[string, *Input])
	for k, v := range visibleFields.Range() {
		FieldRows.Set(k+"：", NewInput(v))
	}
	// modalLayer := component.NewModal()
	const defaultModalAnimationDuration = time.Millisecond * 250
	return &StructView[T]{
		Title:     title,
		rows:      FieldRows,
		applyBtn:  widget.Clickable{},
		closeBtn:  widget.Clickable{},
		onApply:   nil, // 用于刷新编辑够的节点元数据
		unmarshal: unmarshal,
		List: widget.List{
			Scrollbar: widget.Scrollbar{},
			List: layout.List{
				Axis:        layout.Vertical,
				ScrollToEnd: false,
				Position:    layout.Position{},
			},
		},
		ModalState: &component.ModalState{
			ScrimState: component.ScrimState{
				Clickable: widget.Clickable{},
				VisibilityAnimation: component.VisibilityAnimation{
					Duration: defaultModalAnimationDuration,
					// State:    component.Invisible,
					// Started:  gtx.Now,
				},
			},
		},
		Visible: true,
	}
}

func (s *StructView[T]) SetOnApply(f func()) { s.onApply = f }

// RightAlignLabel keygen需要一个下拉框select,form表单, structView, input都可以用这个方法右对齐label
func RightAlignLabel(gtx layout.Context, width unit.Dp, text string) layout.Dimensions { // 布局label右对齐文本
	gtx.Constraints.Min.X = int(width)
	gtx.Constraints.Max.X = gtx.Constraints.Min.X
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		outlay.EmptyFlexed(), // 占据label宽度右对齐,这个的取值方面理解，解决label Alignment:       text.End, 不生效的问题
		//layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		//	return layout.Spacer{}.Layout(gtx)
		//}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Body1(th.Theme, text).Layout(gtx)
			})
		}),
	)
}

func (s *StructView[T]) Layout(gtx layout.Context) layout.Dimensions {
	if s.closeBtn.Clicked(gtx) {
		s.Visible = false
	}

	if s.onApply != nil && s.applyBtn.Clicked(gtx) {
		s.onApply()
	}

	border := widget.Border{
		Color:        color.NRGBA{R: 0xFD, G: 0xB5, B: 0x0E, A: 0xFF},
		CornerRadius: unit.Dp(18),
		Width:        unit.Dp(1),
	}

	rows := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.H4(th.Theme, s.Title).Layout(gtx)
			})
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(20)}.Layout(gtx)
		},
		//func(gtx layout.Context) layout.Dimensions {
		//	//select keygen需要一个下拉框，
		//	return layout.Dimensions{}
		//},
		//func(gtx layout.Context) layout.Dimensions {
		//	//rows
		//	return layout.Dimensions{}
		//},
		//func(gtx layout.Context) layout.Dimensions {
		//	//apply close
		//	return layout.Dimensions{}
		//},
	}

	for k, v := range s.rows.Range() {
		rows = append(rows,
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:    layout.Horizontal,
					Spacing: layout.SpaceSides,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions { // 布局label右对齐文本
						return RightAlignLabel(gtx, MaxLabelWidth(gtx, s.rows.Keys()), k)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = 500
						gtx.Constraints.Max.X = 500
						return layout.Flex{
							Axis: layout.Horizontal,
							// Spacing: 0,
							// Alignment: layout.Middle,
							// WeightSum: 0,
						}.Layout(gtx,
							outlay.EmptyRigidVertical(20),
							// layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return v.Layout(gtx) // todo fake make input
							}))
					}))
			})
	}

	rows = append(rows, func(gtx layout.Context) layout.Dimensions {
		return layout.Spacer{Height: unit.Dp(20)}.Layout(gtx)
	})
	rows = append(rows, func(gtx layout.Context) layout.Dimensions {
		return layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return Button(&s.closeBtn, NavigationCloseIcon, "Close").Layout(gtx)
				}),
				outlay.EmptyRigidHorizontal(10),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return Button(&s.applyBtn, ActionAssignmentTurnedInIcon, "Apply").Layout(gtx)
				}),
				outlay.EmptyRigidHorizontal(10),
			)
		})
	})

	return layout.Inset{Top: unit.Dp(80)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			s.ModalState.Show(gtx.Now, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return material.List(th.Theme, &s.List).Layout(gtx, len(rows), func(gtx layout.Context, index int) layout.Dimensions {
						BackgroundDark(gtx)
						//exact := layout.Exact(image.Point{
						//	X: gtx.Constraints.Max.X / 2,
						//	Y: gtx.Constraints.Max.Y / 2,
						//})
						//rect := clip.Rect{Max: exact.Max}
						//paint.FillShape(gtx.Ops, ColorHeaderFg, rect.Op())
						return rows[index](gtx)
					})
				})
			})
			return component.ModalStyle{
				ModalState: s.ModalState,
				Scrim:      component.NewScrim(th.Theme, &s.ModalState.ScrimState, 0),
			}.Layout(gtx)
		})
		//})
	})
}
