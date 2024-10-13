package ux

//
//import (
//	"gioui.org/op/clip"
//	"gioui.org/op/paint"
//	"gioui.org/widget"
//	"gioui.org/widget/material"
//	"github.com/x-module/gioui-plugins/resource"
//	"github.com/x-module/gioui-plugins/theme"
//	"image"
//	"time"
//
//	"gioui.org/layout"
//	"gioui.org/unit"
//)
//
//type MsgType int
//
//const (
//	InfoMsg MsgType = iota
//	SuccessMsg
//	WaringMsg
//	ErrorMsg
//)
//
//type NoticeItem struct {
//	Text      string
//	EndAt     time.Time
//	closeIcon *Button
//	msgType   MsgType // 通知类型
//}
//
//func NewNoticeItem() *NoticeItem {
//	notice := &NoticeItem{
//		closeIcon: NewIconButton(th, resource.CloseIcon),
//	}
//	notice.closeIcon.SetOnClick(func(gtx layout.Context) {
//		notice.EndAt = time.Now().Add(-time.Hour)
//	})
//	return notice
//}
//
//type Notice struct {
//	notice []*NoticeItem
//	theme  *theme.Theme
//	list   *widget.List
//}
//
//func NewNotification() *Notice {
//	return &Notice{
//		list: &widget.List{
//			List: layout.List{
//				Axis: layout.Vertical,
//				// ScrollToEnd: true,
//			},
//		},
//	}
//}
//
//// 过滤过期的通知
//func (n *Notice) update() {
//	for i := 0; i < len(n.notice); i++ {
//		if time.Now().After(n.notice[i].EndAt) {
//			n.notice = append(n.notice[:i], n.notice[i+1:]...)
//			i--
//		}
//	}
//}
//
//func (n *Notice) AppendNotice(notice *NoticeItem) {
//	n.notice = append(n.notice, notice)
//}
//
//func (n *NoticeItem) layout(gtx layout.Context, theme *theme.Theme) layout.Dimensions {
//	// set max width for the notification
//	gtx.Constraints.Max.X = gtx.Dp(350)
//	// // set max height for the notification
//	gtx.Constraints.Max.Y = gtx.Dp(50)
//
//	// utils.ColorBackground(gtx, gtx.Constraints.Max, resource.GreenColor)
//	return layout.Background{}.Layout(gtx,
//		func(gtx layout.Context) layout.Dimensions {
//			defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, 8).Push(gtx.Ops).Pop()
//			paint.Fill(gtx.Ops, n.theme.Color.NotificationBgColor)
//			return layout.Dimensions{Size: gtx.Constraints.Min}
//		},
//		func(gtx layout.Context) layout.Dimensions {
//			return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
//				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
//					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
//						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
//							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
//								if n.msgType == InfoMsg {
//									return resource.AlertErrorIcon.Layout(gtx, n.theme.Color.NoticeInfoColor)
//								} else if n.msgType == SuccessMsg {
//									return resource.ActionCheckCircleIcon.Layout(gtx, n.theme.Color.NoticeSuccessColor)
//								} else if n.msgType == WaringMsg {
//									return resource.ActionCheckCircleIcon.Layout(gtx, n.theme.Color.NoticeWaringColor)
//								} else if n.msgType == ErrorMsg {
//									return resource.NavigationCancelIcon.Layout(gtx, n.theme.Color.NoticeErrorColor)
//								}
//								return resource.AlertErrorIcon.Layout(gtx, n.theme.Color.NoticeInfoColor)
//							}),
//							layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
//							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
//								return layout.Inset{Top: unit.Dp(3)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
//									bd := material.Body1(theme.Material(), n.Text)
//									bd.Color = n.theme.Color.NotificationTextWhiteColor
//									bd.TextSize = n.theme.Size.DefaultTextSize
//									return bd.Layout(gtx)
//								})
//							}),
//							layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
//						)
//					}),
//					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
//						return n.closeIcon.Layout(gtx)
//					}),
//				)
//			})
//		},
//	)
//}
//func (n *Notice) Layout(gtx layout.Context, theme *theme.Theme) layout.Dimensions {
//	n.update()
//	if len(n.notice) == 0 {
//		return layout.Dimensions{}
//	}
//	return layout.NE.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
//		return layout.Inset{Top: unit.Dp(20), Right: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
//			return material.List(theme.Material(), n.list).Layout(gtx, len(n.notice), func(gtx layout.Context, index int) layout.Dimensions {
//				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
//					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
//						return n.notice[index].layout(gtx, theme)
//					}),
//					layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
//				)
//			})
//		})
//	})
//}
