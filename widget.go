package ux

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/stream"
	"github.com/ddkwork/ux/languages"
	"github.com/ddkwork/ux/resources/images"
	"github.com/ddkwork/ux/widget/material"
)

type Widget interface {
	Layout(gtx layout.Context) layout.Dimensions
}

type Widgets interface { // map[string]Widget
	Button() Widget
	Card() Widget
	Checkbox() Widget
	CodeEditor() Widget
	ColorPicker() Widget
	DateInput() Widget
	DateTimeInput() Widget
	Dialog() Widget
	Divider() Widget
	FileInput() Widget
	Flex() Widget
	FormItem() Widget
	Grid() Widget
	Icon() Widget
	Image() Widget
	Input() Widget
	Label() Widget
	Link() Widget
	Menu() Widget
	Radio() Widget
	Slider() Widget
	Space() Widget
	Spin() Widget
	Switch() Widget
	Tabs() Widget
	Text() Widget
	TextField() Widget
	TimeInput() Widget
	Tooltip() Widget
	LogView() Widget
	Modal() Widget
	Select() Widget
	SelectEntry() Widget
	Form() Widget
	StructView() Widget
	ProgressBar() Widget
	Toolbar() Widget
	List() Widget
	TreeTable() Widget
	terminal() Widget
	calculator() Widget
	Calendar() Widget
	Markdown() Widget
}

type LogView struct {
	*ContextMenu
	*CodeEditor
}

func NewLogView() *LogView {
	l := &LogView{
		ContextMenu: NewContextMenu(),
		CodeEditor:  NewCodeEditor("", languages.GoKind),
	}
	mylog.SetCallBack(func(row string) { l.CodeEditor.AppendText(row) })
	l.ContextMenu.AddItem(ContextMenuItem{
		Title: "save",
		Icon:  images.ContentSaveIcon,
		Can:   func() bool { return true },
		Do: func() {
			if !stream.IsAndroid() {
				return
			}
			go func() {
				f := mylog.Check2(explore.CreateFile("sata.json"))
				mylog.Check2(f.Write([]byte(l.CodeEditor.editor.Text())))
				mylog.Check(f.Close())
			}()
		},
		AppendDivider: false,
		Clickable:     widget.Clickable{},
	})
	return l
}

func (l *LogView) Layout(gtx layout.Context) layout.Dimensions {
	return l.ContextMenu.Layout(gtx, []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return l.CodeEditor.Layout(gtx)
		},
	})
}

func LayoutErrorLabel(gtx layout.Context, e error) layout.Dimensions {
	if e != nil {
		return layout.Inset{
			Top:    10,
			Bottom: 10,
			Left:   15,
			Right:  15,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			label := material.Label(th, th.TextSize*0.8, e.Error())
			label.Color = color.NRGBA{R: 255, A: 255}
			label.Alignment = text.Middle
			return label.Layout(gtx)
		})
	} else {
		return layout.Dimensions{}
	}
}

// func layoutEmoji(gtx layout.Context) layout.Dimensions {
//	var sel widget.Selectable
//	message := "🥳🧁🍰🎁🎂🎈🎺🎉🎊\n📧〽️🧿🌶️🔋\n😂❤️😍🤣😊\n🥺🙏💕😭😘\n👍😅👏"
//	var customTruncator widget.Bool
//	var maxLines widget.Float
//	maxLines.Value = 0
//
//	const (
//		minLinesRange = 1
//		maxLinesRange = 5
//	)
//
//	inset := layout.UniformInset(5)
//
//	l := material.H4(th, message)
//	if customTruncator.Value {
//		l.Truncator = "cont..."
//	} else {
//		l.Truncator = ""
//	}
//	l.MaxLines = minLinesRange + int(math.Round(float64(maxLines.Value)*(maxLinesRange-minLinesRange)))
//	l.State = &sel
//	return inset.Layout(gtx, l.Layout)
// }
