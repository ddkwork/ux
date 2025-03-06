package ux

import "github.com/ddkwork/golibrary/mylog"

func LogView() Widget {
	logView := NewCodeEditor(mylog.Row(), CodeLanguageGolang)
	mylog.SetCallBack(func() {
		logView.AppendText(mylog.Row())
	})
	return logView.Layout
}
