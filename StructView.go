package ux

import (
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/x/outlay"
	"github.com/ddkwork/golibrary/safemap"
	"github.com/ddkwork/golibrary/stream"
	"github.com/ddkwork/ux/resources/colors"
	"github.com/ddkwork/ux/resources/images"
	"github.com/ddkwork/ux/widget/material"
	"github.com/ddkwork/ux/x/component"
)

// 相当于StructView，formView，但是他是绑定结构体的

type StructView[T any] struct { // 其实就是一个标题row+list滚动多个row
	Data T          // 结构体元数据
	Rows []CellData // 在onApply回调中反序列化成结构体

	// 布局预期视觉样式：1个元素的一行
	Title string // 标题，label渲染，不进入list滚动，需要固定位置，字体大小是h4

	// select keygen需要一个下拉框，
	// 布局预期视觉样式：2个元素的一行，key右对齐，value左对齐

	// 布局预期视觉样式：2个元素的一行，key右对齐，value左对齐
	fields *safemap.M[string, *Input] // 垂直flex，树形表格布局却是水平flex为row,unmarshal应该调用这个

	// 底部布局成一行，对于密码箱的rsa，这里需要额外的按钮，所以应该增加个outlay布局，
	// 填充到这一行的左侧，关闭和应用始终在右下角 layout.N
	// 这样如果不在list内滚动就不会渲染，不知道是什么原因，总感觉一个整体的widget只能有一个list，先滚动吧
	// 布局预期视觉样式：多个元素的一行
	applyBtn *widget.Clickable
	closeBtn *widget.Clickable
	onApply  func()

	widget.List // 滚动所有行
	*component.ModalState
	Visible bool
	Modal   bool
	layout.Inset
}

func (s *StructView[T]) Unmarshal(callback UnmarshalRowCallback) T {
	s.Data = UnmarshalRow[T](s.Rows, callback)
	return s.Data
}

func NewStructView[T any](title string, object T, callback MarshalRowCallback) *StructView[T] {
	FieldRows := new(safemap.M[string, *Input])
	rows := MarshalRow(object, callback)
	for _, row := range rows {
		input := NewInput(row.Key)
		input.SetText(row.Value)
		input.editor.Alignment = text.Start // 左对齐 todo bug, not work
		FieldRows.Update(row.Key+"：", input)
	}
	const defaultModalAnimationDuration = time.Millisecond * 250
	return &StructView[T]{
		Data:     object,
		Rows:     rows,
		Title:    title,
		fields:   FieldRows,
		applyBtn: &widget.Clickable{},
		closeBtn: &widget.Clickable{},
		onApply:  nil, // 用于刷新编辑够的节点元数据
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
		Modal:   false,
	}
}

func (s *StructView[T]) SetOnApply(f func()) { s.onApply = f }

// RightAlignLabel keygen需要一个下拉框select,form表单, structView, input都可以用这个方法右对齐label
func RightAlignLabel(gtx layout.Context, maxWidth unit.Dp, text string) layout.Dimensions { // 布局label右对齐文本
	gtx.Constraints.Min.X = int(maxWidth)
	gtx.Constraints.Max.X = gtx.Constraints.Min.X
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		outlay.EmptyFlexed(), // 占据label宽度右对齐,这个的取值方面理解，解决label Alignment:       text.End, 不生效的问题
		// layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		//	return layout.Spacer{}.Layout(gtx)
		// }),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Body1(th, text).Layout(gtx)
			})
		}),
	)
}

func (s *StructView[T]) Layout(gtx layout.Context) layout.Dimensions {
	for i, row := range s.Rows {
		s.Rows[i].Value = s.fields.GetMust(row.Key + "：").GetText() // 刷新输入框的值
	}

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
				return material.H4(th, s.Title).Layout(gtx)
			})
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(20)}.Layout(gtx)
		},
		// func(gtx layout.Context) layout.Dimensions {
		//	//select keygen需要一个下拉框，
		//	return layout.Dimensions{}
		// },
		// func(gtx layout.Context) layout.Dimensions {
		//	//fields
		//	return layout.Dimensions{}
		// },
		// func(gtx layout.Context) layout.Dimensions {
		//	//apply close
		//	return layout.Dimensions{}
		// },
	}

	for k, v := range s.fields.Range() {
		rows = append(rows,
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceSides}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions { // 布局label右对齐文本
						return RightAlignLabel(gtx, MaxLabelWidth(gtx, s.Rows), k)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = 700
						gtx.Constraints.Max.X = 700
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
		return layout.Spacer{Height: unit.Dp(20)}.Layout(gtx) // 垂直间距
	})
	rows = append(rows, func(gtx layout.Context) layout.Dimensions {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				outlay.EmptyRigidHorizontal(300), // 标签站位
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return Button(s.closeBtn, images.NavigationCloseIcon, "Close").Layout(gtx)
				}),
				outlay.EmptyRigidHorizontal(20), // 按钮间距
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return Button(s.applyBtn, images.ActionAssignmentTurnedInIcon, "Apply").Layout(gtx)
				}),
				outlay.EmptyRigidHorizontal(10),
			)
		})
	})
	if stream.IsAndroid() {
		s.Inset = layout.Inset{}
	}
	return s.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			s.ModalState.Show(gtx.Now, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return material.List(th, &s.List).Layout(gtx, len(rows), func(gtx layout.Context, index int) layout.Dimensions {
						return Background{Color: colors.BackgroundColor}.Layout(gtx, func(gtx layout.Context) layout.Dimensions { // todo 这样把边框整没了
							return rows[index](gtx)
						})
					})
				})
			})
			alpha := byte(0)
			if s.Modal {
				alpha = 20
			}
			return component.ModalStyle{ModalState: s.ModalState, Scrim: component.NewScrim(th, &s.ModalState.ScrimState, alpha)}.Layout(gtx)
		})
	})
}
