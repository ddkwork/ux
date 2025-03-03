package ux

//
//import (
//	"time"
//)
//
//var NotificationController = NewNotification()
////var SystemNoticeController = &SystemNotice{}
//
//func SendAppInfoNotice(text string, duration ...time.Duration) {
//	dru := time.Second * 4
//	if len(duration) > 0 {
//		dru = duration[0]
//	}
//	notice := NewNoticeItem()
//	notice.Text = text
//	notice.msgType = InfoMsg
//	notice.EndAt = time.Now().Add(dru)
//	NotificationController.AppendNotice(notice)
//}
//func SendAppSuccessNotice(text string, duration ...time.Duration) {
//	dru := time.Second * 4
//	if len(duration) > 0 {
//		dru = duration[0]
//	}
//	notice := NewNoticeItem()
//	notice.Text = text
//	notice.msgType = SuccessMsg
//	notice.EndAt = time.Now().Add(dru)
//	NotificationController.AppendNotice(notice)
//}
//func SendAppWaringNotice(text string, duration ...time.Duration) {
//	dru := time.Second * 4
//	if len(duration) > 0 {
//		dru = duration[0]
//	}
//	notice := NewNoticeItem()
//	notice.Text = text
//	notice.msgType = WaringMsg
//	notice.EndAt = time.Now().Add(dru)
//	NotificationController.AppendNotice(notice)
//}
//func SendAppErrorNotice(text string, duration ...time.Duration) {
//	dru := time.Second * 4
//	if len(duration) > 0 {
//		dru = duration[0]
//	}
//	notice := NewNoticeItem()
//	notice.Text = text
//	notice.msgType = ErrorMsg
//	notice.EndAt = time.Now().Add(dru)
//	NotificationController.AppendNotice(notice)
//}
//
//func SendSystemNotice(message string) {
//	//_ = SystemNoticeController.Notice(message)
//}
