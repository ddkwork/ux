package ux

import (
	"image"
	"image/color"
	"io"
	"slices"
	"strings"

	"gioui.org/gesture"
	"gioui.org/io/clipboard"
	"gioui.org/io/event"
	"gioui.org/io/input"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/transfer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/ddkwork/ux/widget/material"
	"github.com/ddkwork/ux/x/component"
)

type state uint8

const (
	inactive state = iota
	hovered
	activated
	focused
)

type (
	ActionFun func(gtx layout.Context)
	Input     struct {
		editor    widget.Editor
		height    unit.Dp
		before    layout.Widget
		after     layout.Widget
		icon      *widget.Icon
		iconClick widget.Clickable

		click       gesture.Click
		state       state
		borderColor color.NRGBA
		bgColor     color.NRGBA
		hint        string
		radius      unit.Dp
		size        ElementStyle
		width       unit.Dp
		hasBorder   bool

		showPassword bool

		onIconClick ActionFun
		onFocus     ActionFun
		onLostFocus ActionFun
		onChange    func(text string)

		contextMenu *ContextMenu
		contextArea *component.ContextArea
	}
)

func (i *Input) SetHasBorder(hasBorder bool) *Input {
	i.hasBorder = hasBorder
	return i
}

func NewInput(hint string, text ...string) *Input {
	t := &Input{
		editor: widget.Editor{},
		// maxIndentWidth:  th.Size.DefaultElementWidth,
	}
	t.size = th.Size.Medium
	t.hint = hint
	t.radius = th.Size.DefaultElementRadiusSize
	if len(text) > 0 {
		t.editor.SetText(text[0])
	}
	t.editor.SingleLine = true
	return t
}

func NewTextArea(hint string, text ...string) *Input {
	t := NewInput(hint, text...)
	t.height = unit.Dp(100)
	t.editor.SingleLine = false
	return t
}

func (i *Input) SetOnFocus(f ActionFun) *Input {
	i.onFocus = f
	return i
}

func (i *Input) SetOnLostFocus(f ActionFun) *Input {
	i.onLostFocus = f
	return i
}

func (i *Input) SetHeight(height unit.Dp) *Input {
	i.height = height
	return i
}

func (i *Input) SetWidth(width unit.Dp) *Input {
	i.width = width
	return i
}

func (i *Input) SetOnIconClick(f ActionFun) *Input {
	i.onIconClick = f
	return i
}

func (i *Input) SetOnChanged(f func(text string)) *Input {
	i.onChange = f
	return i
}

func (i *Input) Password() *Input {
	i.editor.Mask = '*'
	i.icon = ActionVisibilityOffIcon
	// t.IconPositionEnd = IconPositionEnd
	i.showPassword = false
	return i
}

func (i *Input) SetIcon(icon *widget.Icon) *Input {
	i.icon = icon
	return i
}

// SetRadius 设置radius
func (i *Input) SetRadius(radius unit.Dp) *Input {
	i.radius = radius
	return i
}

func (i *Input) ReadOnly() *Input {
	i.editor.ReadOnly = true
	return i
}

func (i *Input) SetBefore(before layout.Widget) *Input {
	i.before = before
	return i
}

func (i *Input) SetAfter(after layout.Widget) *Input {
	i.after = after
	return i
}

func (i *Input) SetSize(size ElementStyle) *Input {
	i.size = size
	return i
}

func (i *Input) SetText(text string) *Input {
	i.editor.SetText(text)
	return i
}

func (i *Input) GetText() string {
	return i.editor.Text()
}

func (i *Input) update(gtx layout.Context) {
	if gtx.Focused(&i.editor) {
		if i.onFocus != nil {
			i.onFocus(gtx)
		}
	} else {
		if i.onLostFocus != nil {
			i.onLostFocus(gtx)
		}
	}
	disabled := gtx.Source == (input.Source{})
	for {
		ev, ok := i.click.Update(gtx.Source)
		if !ok {
			break
		}
		switch ev.Kind {
		case gesture.KindPress:
			gtx.Execute(key.FocusCmd{Tag: &i.editor})
		case gesture.KindClick:

		default:

		}
	}
	i.state = inactive
	if i.click.Hovered() && !disabled {
		i.state = hovered
	}
	// if t.editor.Len() > 0 {
	// 	t.state = activated
	// }
	if gtx.Source.Focused(&i.editor) && !disabled {
		i.state = focused
	}

	i.bgColor = th.Color.DefaultBgGrayColor

	if i.editor.ReadOnly {
		return
	}

	switch i.state {
	case inactive:
		i.borderColor = th.Color.InputInactiveBorderColor
	case hovered:
		i.borderColor = th.Color.InputHoveredBorderColor
	case focused:
		i.bgColor = th.Color.InputFocusedBgColor
		i.borderColor = th.Color.InputFocusedBorderColor
	case activated:
		i.borderColor = th.Color.InputActivatedBorderColor
	}
	for {
		e, ok := i.editor.Update(gtx)
		if !ok {
			break
		}
		if _, ok := e.(widget.ChangeEvent); ok {
			if i.onChange != nil {
				i.onChange(i.GetText())
			}
		}
	}
}

func (i *Input) Layout(gtx layout.Context) layout.Dimensions {
	if i.width > 0 {
		gtx.Constraints.Max.X = gtx.Dp(i.width)
	}
	i.update(gtx)
	// gtx.Constraints.Min.X = gtx.Constraints.Max.X
	// gtx.Constraints.Min.Y = 0
	macro := op.Record(gtx.Ops)
	dims := layout.Inset{Top: 4}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return i.layout(gtx)
	})
	call := macro.Stop()
	defer pointer.PassOp{}.Push(gtx.Ops).Pop()
	defer clip.Rect(image.Rectangle{Max: dims.Size}).Push(gtx.Ops).Pop()
	i.click.Add(gtx.Ops)
	event.Op(gtx.Ops, &i.editor)
	call.Add(gtx.Ops)
	return dims
}

func (i *Input) layout(gtx layout.Context) layout.Dimensions {
	border := widget.Border{
		Color:        i.borderColor,
		Width:        unit.Dp(1),
		CornerRadius: i.radius,
	}
	return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Background{}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				rr := gtx.Dp(i.radius)
				defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, rr).Push(gtx.Ops).Pop()
				paint.Fill(gtx.Ops, i.bgColor)
				return layout.Dimensions{Size: gtx.Constraints.Min}
			},
			func(gtx layout.Context) layout.Dimensions {
				return i.size.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					inputLayout := layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						if i.width > 0 {
							gtx.Constraints.Max.X = gtx.Dp(i.width)
						}
						editor := material.Editor(th.Theme, &i.editor, i.hint)
						editor.HintColor = th.Color.HintTextColor
						editor.SelectionColor = th.Color.TextSelectionColor

						gtx.Constraints.Min.Y = gtx.Dp(i.size.Height) // 设置最小高度为 100dp
						gtx.Constraints.Max.Y = gtx.Constraints.Min.Y // 限制最大高度与最小高度相同
						editor.TextSize = i.size.TextSize

						if i.height > 0 {
							gtx.Constraints.Min.Y = gtx.Dp(i.height)      // 设置最小高度为 100dp
							gtx.Constraints.Max.Y = gtx.Constraints.Min.Y // 限制最大高度与最小高度相同
						}
						if i.editor.ReadOnly {
							editor.Color = th.Color.HintTextColor
						} else {
							editor.Color = th.Color.DefaultTextWhiteColor
						}
						return layout.Stack{Alignment: layout.Center}.Layout(gtx,
							layout.Stacked(func(gtx layout.Context) layout.Dimensions {
								return editor.Layout(gtx)
							}),
							layout.Expanded(func(gtx layout.Context) layout.Dimensions { // 这也应该适用于代码编辑器控件
								if i.contextArea == nil {
									i.contextArea = &component.ContextArea{
										Activation:       pointer.ButtonSecondary,
										AbsolutePosition: true,
										PositionHint:     layout.W,
									}
								}
								if i.contextMenu == nil {
									i.contextMenu = NewContextMenu()
									items := []ContextMenuItem{
										{
											Title: "copy",
											Icon:  SvgIconCopy,
											Can:   func() bool { return true },
											Do: func() {
												// i.editor.SelectedText()
												// todo add selectAll api,或者在这里发送按键按下的事件?
												// e.text.SetCaret(0, e.text.Len())//ctrl +a,复制的同时应该显示选中全部，不然不知道复制的区域
												// 安卓上应该增加一个全选菜单，以及选中的左右拉伸拓展选中区域的功能
												gtx.Execute(clipboard.WriteCmd{Data: io.NopCloser(strings.NewReader(i.editor.Text()))})
											},
											AppendDivider: false,
											Clickable:     widget.Clickable{},
										},
										{
											Title: "paste",
											Icon:  SvgIconContentPasteTwotone,
											Can:   func() bool { return true },
											Do: func() {
												for {
													ke, ok := gtx.Event(transfer.TargetFilter{Target: &i.editor, Type: "application/text"})
													if !ok {
														break
													}
													switch ke := ke.(type) {
													case transfer.DataEvent:
														content, err := io.ReadAll(ke.Open())
														if err == nil {
															if i.editor.Insert(string(content)) != 0 {
																break
															}
														}
													}
												}
												gtx.Execute(clipboard.ReadCmd{Tag: &i.editor})
											},
											AppendDivider: false,
											Clickable:     widget.Clickable{},
										},
										{
											Title:         "clean",
											Icon:          SvgIconTrash,
											Can:           func() bool { return true },
											Do:            func() { i.editor.SetText("") },
											AppendDivider: false,
											Clickable:     widget.Clickable{},
										},
									}
									for _, item := range items {
										if item.Can() {
											i.contextMenu.AddItem(item)
										}
									}
								}
								i.contextMenu.OnClicked(gtx)
								return i.contextArea.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return drawContextArea(gtx, &i.contextMenu.MenuState)
								})
							}),
						)
					})

					var widgets []layout.FlexChild
					if i.before != nil {
						widgets = append(widgets, layout.Rigid(i.before))
					}
					widgets = append(widgets, inputLayout)
					if i.icon != nil {
						iconLayout := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if i.iconClick.Clicked(gtx) {
								if i.onIconClick != nil {
									i.onIconClick(gtx)
								}
								if !i.showPassword {
									i.editor.Mask = 0
									i.icon = ActionVisibilityIcon
									i.showPassword = true
								} else {
									i.editor.Mask = '*'
									i.icon = ActionVisibilityOffIcon
									i.showPassword = false
								}
							}
							return i.iconClick.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return i.icon.Layout(gtx, th.Color.DefaultIconColor)
							})
						})
						// widgets = append(widgets, iconLayout)
						widgets = slices.Insert(widgets, 0, iconLayout)
					} else {
						if i.after != nil {
							widgets = append(widgets, layout.Rigid(i.after))
						}
					}
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx, widgets...)
				})
			},
		)
	})
}
