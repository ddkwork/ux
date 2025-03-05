package ux

import (
	"fmt"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/ux/dlgs"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type DirSelector struct {
	input       *Input
	dirName     string
	actionClick widget.Clickable
	windowTitle string
	onSelectDir func(dir string)
	changed     bool
	width       unit.Dp
}

func NewDirSelector(hint string, dirName ...string) *DirSelector {
	bf := &DirSelector{
		input:       NewInput(hint, dirName...),
		width:       th.Size.DefaultElementWidth,
		windowTitle: "Select Directory",
	}
	if len(dirName) > 0 {
		bf.dirName = dirName[0]
		bf.input.SetText(dirName[0])
	}
	bf.updateIcon()
	return bf
}

// SetWidth 设置width
func (b *DirSelector) SetWidth(width unit.Dp) *DirSelector {
	b.width = width
	return b
}

// SetWindowTitle 设置windowTitle
func (b *DirSelector) SetWindowTitle(title string) *DirSelector {
	b.windowTitle = title
	return b
}

func (b *DirSelector) action(gtx layout.Context) {
	if b.actionClick.Clicked(gtx) {
		if b.dirName != "" {
			b.RemoveDir()
			b.changed = true
			return
		} else {
			dir, _ := mylog.Check3(dlgs.File(b.windowTitle, "", true))
			fmt.Println("Selected Directory:", dir)
			if dir == "" {
				return
			}
			b.setDirName(dir)
			b.changed = true
			if b.onSelectDir != nil {
				b.onSelectDir(dir)
			}
		}
	}
}

func (b *DirSelector) SetOnSelectDir(f func(dir string)) *DirSelector {
	b.onSelectDir = f
	return b
}

func (b *DirSelector) setDirName(name string) *DirSelector {
	b.dirName = name
	b.input.SetText(name)
	b.updateIcon()
	b.changed = true
	return b
}

func (b *DirSelector) Changed() bool {
	out := b.changed
	b.changed = false
	return out
}

func (b *DirSelector) RemoveDir() {
	b.dirName = ""
	b.input.SetText("")
	b.updateIcon()
	b.changed = true
}

func (b *DirSelector) GetDirPath() string {
	return b.dirName
}

func (b *DirSelector) updateIcon() {
	if b.dirName != "" {
		b.input.SetAfter(func(gtx layout.Context) layout.Dimensions {
			return b.actionClick.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.X = gtx.Dp(th.Size.DefaultIconSize)
				return ActionDeleteIcon.Layout(gtx, th.Color.DefaultIconColor)
			})
		})
	} else {
		b.input.SetAfter(func(gtx layout.Context) layout.Dimensions {
			return b.actionClick.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.X = gtx.Dp(th.Size.DefaultIconSize)
				return FileFileUploadIcon.Layout(gtx, th.Color.DefaultIconColor)
			})
		})
	}
}

func (b *DirSelector) Layout(gtx layout.Context) layout.Dimensions {
	// gtx.Constraints.Max.Y = gtx.Dp(42)
	b.action(gtx)
	gtx.Constraints.Max.X = gtx.Dp(b.width)
	return b.input.Layout(gtx)
}
