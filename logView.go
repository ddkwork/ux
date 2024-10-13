package ux

import "github.com/ddkwork/golibrary/mylog"

func LogView() Widget {
	logView := NewCodeEditor(mylog.Body(), CodeLanguageGO)
	mylog.SetCallBack(func() {
		logView.AppendText(mylog.Body())
	})
	return logView.Layout
}
