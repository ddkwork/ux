package ux

//
//import (
//	"gioui.org/layout"
//	"gioui.org/op/clip"
//	"gioui.org/op/paint"
//	"gioui.org/unit"
//	"gioui.org/widget"
//
//	"image"
//)
//
//type Action int
//
//const (
//	BothAction Action = iota
//	OnlyCancelAction
//	OnlyConfirmAction
//)
//
//type Confirm struct {
//	title            string
//	message          string
//	height           int
//	dragWidth            int
//	visible          bool
//	cancelClickable  widget.Clickable
//	confirmClickable widget.Clickable
//	clickerWidget    *widget.Clickable
//	action           Action
//	cancelFunc       func()
//	confirmFunc      func()
//	customAction     []layout.FlexChild
//}
//
//func NewConfirm() *Confirm {
//	modal := &Confirm{
//		height:        150,
//		dragWidth:         300,
//		title:         "操作确认",
//		//clickerWidget: NewClickable(th),
//		action:        BothAction,
//	}
//	return modal
//}
//func (c *Confirm) SetAction(action Action) *Confirm {
//	c.action = action
//	return c
//}
//func (c *Confirm) Confirm(fun func()) *Confirm {
//	c.confirmFunc = fun
//	return c
//}
//func (c *Confirm) Cancel(fun func()) *Confirm {
//	c.cancelFunc = fun
//	return c
//}
//
//func (c *Confirm) SetWidth(dragWidth int) *Confirm {
//	c.dragWidth = dragWidth
//	return c
//}
//func (c *Confirm) Visible() bool {
//	return c.visible
//}
//func (c *Confirm) SetTitle(title string) *Confirm {
//	c.title = title
//	return c
//}
//
//func (c *Confirm) SetCustomAction(customAction []layout.FlexChild) *Confirm {
//	c.customAction = customAction
//	return c
//}
//
//func (c *Confirm) SetHeight(height int) *Confirm {
//	c.height = height
//	return c
//}
//
//func (c *Confirm) Message(message string) *Confirm {
//	c.message = message
//	c.visible = true
//	return c
//}
//func (c *Confirm) Close() {
//	c.visible = false
//}
//
//func (c *Confirm) Layout(gtx layout.Context) layout.Dimensions {
//	if !c.visible {
//		return layout.Dimensions{}
//	}
//	if c.visible {
//		// 绘制全屏半透明遮罩层
//		paint.Fill(gtx.Ops, th.Color.DefaultMaskBgColor)
//	}
//	for c.cancelClickable.OnClicked(gtx) {
//		c.visible = false
//		if c.cancelFunc != nil {
//			c.cancelFunc()
//		}
//	}
//	for c.confirmClickable.OnClicked(gtx) {
//		c.visible = false
//		if c.confirmFunc != nil {
//			c.confirmFunc()
//		}
//	}
//	dragWidth := gtx.Dp(unit.Dp(c.dragWidth))
//	height := gtx.Dp(unit.Dp(c.height))
//
//	var actions []layout.FlexChild
//
//	if c.customAction != nil {
//		actions = c.customAction
//	} else {
//		cancelAction := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
//			but := DefaultButton(th, &c.cancelClickable, "取消", unit.Dp(70), layout.Inset{
//				Top: 3, Bottom: 3,
//				Left: 5, Right: 5,
//			})
//			return but.Layout(gtx)
//		})
//
//		confirmAction := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
//			but := SuccessButton(th, &c.confirmClickable, "确认", unit.Dp(70), layout.Inset{
//				Top: 3, Bottom: 3,
//				Left: 5, Right: 5,
//			})
//			return but.Layout(gtx)
//		})
//
//		space := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
//			return layout.Dimensions{Size: image.Point{X: 20, Y: 0}}
//		})
//
//		if c.action == BothAction {
//			actions = []layout.FlexChild{
//				cancelAction,
//				space,
//				confirmAction,
//			}
//
//		}
//		if c.action == OnlyCancelAction {
//			actions = []layout.FlexChild{
//				cancelAction,
//			}
//		} else if c.action == OnlyConfirmAction {
//			actions = []layout.FlexChild{
//				confirmAction,
//			}
//		}
//	}
//
//	return c.clickerWidget.SetWidget(func(gtx layout.Context) layout.Dimensions {
//		return layout.Inset{
//			Top: unit.Dp(10),
//		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
//			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
//				// Set the size of the Confirm
//				gtx.Constraints = layout.Exact(image.Point{X: dragWidth, Y: height})
//				rc := clip.RRect{
//					Rect: image.Rectangle{Max: image.Point{
//						X: gtx.Constraints.Min.X,
//						Y: gtx.Constraints.Min.Y,
//					}},
//					NW: 10, NE: 10, SE: 10, SW: 10,
//				}
//				paint.FillShape(gtx.Ops, th.Color.DefaultContentBgGrayColor, rc.Op(gtx.Ops))
//				// Center the text inside the Confirm
//				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
//					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
//						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
//							return layout.Inset{Left: 0, Right: 10, Bottom: 10, Top: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
//								return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
//									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
//										gtx.Constraints.Max.X = gtx.Dp(unit.Dp(20))
//										return resource.ActionInfoOutlineIcon.Layout(gtx, th.Color.GreenColor)
//									}),
//									layout.Rigid(layout.Spacer{Width: unit.Dp(5)}.Layout),
//									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
//										return material.Body1(th, c.title).Layout(gtx)
//									}),
//								)
//
//							})
//						})
//					}),
//					utils.DrawLineFlex(th.Color.DefaultLineColor, unit.Dp(1), unit.Dp(c.dragWidth)),
//					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
//						gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(50))
//						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
//							return layout.Inset{Left: 5, Right: 5, Bottom: 2, Top: 2}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
//								return Label(th, c.message).Layout(gtx)
//							})
//						})
//					}),
//					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
//						return layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
//							return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
//								return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, actions...)
//							})
//						})
//					}),
//				)
//			})
//		})
//	}).Layout(gtx)
//}
