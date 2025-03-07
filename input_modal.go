package ux

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/x/outlay"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/safemap"
	"github.com/ddkwork/ux/widget/material"
	"github.com/ddkwork/ux/x/component"
	"reflect"
)

// 相当于StructView，formView，但是他是绑定结构体的
// 使用有序map来存储布局信息
//type FieldRow struct {
//	Key   string
//	Value *Input
//}

type InputModal[T any] struct { //其实就是一个标题row+list滚动多个row

	//布局预期视觉样式：1个元素的一行
	Title string //标题，label渲染，不进入list滚动，需要固定位置，字体大小是h4

	//select keygen需要一个下拉框，
	//布局预期视觉样式：2个元素的一行，key右对齐，value左对齐

	//布局预期视觉样式：2个元素的一行，key右对齐，value左对齐
	rows *safemap.M[string, *Input] //

	//底部布局成一行，对于密码箱的rsa，这里需要额外的按钮，所以应该增加个outlay布局，
	//填充到这一行的左侧，关闭和应用始终在右下角 layout.N
	//这样如果不在list内滚动就不会渲染，不知道是什么原因，总感觉一个整体的widget只能有一个list，先滚动吧
	//布局预期视觉样式：多个元素的一行
	applyBtn widget.Clickable
	closeBtn widget.Clickable
	onApply  func()
	onClose  func()

	unmarshal unmarshalFun //反序列化函数，用于将输入框的值反序列化成结构体

	widget.List //滚动所有行

	Visit bool //是否显示模态窗口
}

type marshalFun func(any) []string
type unmarshalFun func([]string) any

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

func NewInputModal[T any](title string, object T, marshal marshalFun, unmarshal unmarshalFun) *InputModal[T] {
	visibleFields := ReflectStruct2Map(object, marshal, unmarshal)
	FieldRows := new(safemap.M[string, *Input])
	for k, v := range visibleFields.Range() {
		FieldRows.Set(k+"：", NewInput(v))
	}
	return &InputModal[T]{
		Title:     title,
		rows:      FieldRows,
		applyBtn:  widget.Clickable{},
		closeBtn:  widget.Clickable{},
		onApply:   nil, //用于刷新编辑够的节点元数据
		onClose:   nil,
		unmarshal: unmarshal,
		List: widget.List{
			Scrollbar: widget.Scrollbar{},
			List: layout.List{
				Axis:        layout.Vertical,
				ScrollToEnd: false,
				Position:    layout.Position{},
			},
		},
		Visit: true, //默认显示，点击关闭按钮后变为false隐藏渲染
	}
}

func (m *InputModal[T]) SetOnClose(f func()) { m.onClose = f }
func (m *InputModal[T]) SetOnApply(f func()) { m.onApply = f }
func (m *InputModal[T]) Layout(gtx layout.Context) layout.Dimensions {
	ops := op.Record(gtx.Ops)
	dims := m.layout(gtx)
	defer op.Defer(gtx.Ops, ops.Stop())
	return dims
}

// RightAlignLabel keygen需要一个下拉框select,form表单, structView, input都可以用这个方法右对齐label
func RightAlignLabel(gtx layout.Context, width unit.Dp, text string) layout.Dimensions { //布局label右对齐文本
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

func (m *InputModal[T]) layout(gtx layout.Context) layout.Dimensions {
	if m.closeBtn.Clicked(gtx) {
		m.onClose()
		m.Visit = false
	}

	if m.onApply != nil && m.applyBtn.Clicked(gtx) {
		m.onApply()
	}

	border := widget.Border{
		Color:        th.BorderBlueColor,
		CornerRadius: unit.Dp(4),
		Width:        unit.Dp(1),
	}

	rows := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.H4(th.Theme, m.Title).Layout(gtx)
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

	for k, v := range m.rows.Range() {
		rows = append(rows,
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:    layout.Horizontal,
					Spacing: layout.SpaceSides,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions { //布局label右对齐文本
						return RightAlignLabel(gtx, MaxLabelWidth(gtx, m.rows.Keys()), k)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = 500
						gtx.Constraints.Max.X = 500
						return layout.Flex{
							Axis:      layout.Horizontal,
							Spacing:   0,
							Alignment: 0,
							WeightSum: 0,
						}.Layout(gtx,
							outlay.EmptyRigidVertical(20),
							//layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return v.Layout(gtx) //todo fake make input
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
					return Button(&m.closeBtn, NavigationCloseIcon, "Close").Layout(gtx)
				}),
				outlay.EmptyRigidHorizontal(10),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return Button(&m.applyBtn, ActionAssignmentTurnedInIcon, "Apply").Layout(gtx)
				}),
				outlay.EmptyRigidHorizontal(10),
			)
		})
	})

	return layout.Inset{Top: unit.Dp(80)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return component.NewModalSheet(component.NewModal()).Layout(gtx, th.Theme, &component.VisibilityAnimation{}, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return material.List(th.Theme, &m.List).Layout(gtx, len(rows), func(gtx layout.Context, index int) layout.Dimensions {
						return rows[index](gtx)
					})
				})
			})
		})
	})
}
