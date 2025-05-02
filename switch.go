package ux

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/ddkwork/ux/widget/material"
)

type Switch struct {
	open material.SwitchStyle
	ok   widget.Bool
}

func NewSwitch(description string) *Switch {
	ok := widget.Bool{}
	openSwitch := &Switch{
		ok:   ok,
		open: material.Switch(th, &ok, description),
	}
	openSwitch.open.Color.Enabled = th.Color.GreenColor
	openSwitch.open.Color.Disabled = th.Color.InfoColor
	return openSwitch
}

func (s *Switch) Open() bool {
	return s.ok.Value
}

func (s *Switch) Layout(gtx layout.Context) layout.Dimensions {
	return s.open.Layout(gtx)
}
