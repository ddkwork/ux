package ux

import (
	"unsafe"

	"github.com/ddkwork/golibrary/std/mylog"
	"golang.org/x/sys/windows"
)

func setCaptionColor(handle uintptr, color uint32) {
	const DWMWA_CAPTION_COLOR = 35
	mylog.Check(windows.DwmSetWindowAttribute(windows.HWND(handle), DWMWA_CAPTION_COLOR, unsafe.Pointer(&color), uint32(unsafe.Sizeof(color))))
}
