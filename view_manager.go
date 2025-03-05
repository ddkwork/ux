package ux

import (
	"net/url"

	"gioui.org/layout"
	"gioui.org/widget"
)

type ViewID struct {
	name string
	path string
}

type ViewAction struct {
	Name      string
	Icon      *widget.Icon
	OnClicked func(gtx C)
}

type View interface {
	Actions() []ViewAction
	Layout(gtx layout.Context) layout.Dimensions
	ID() ViewID
	Location() url.URL
	Title() string
	// set the view to finished state and do some cleanup ops.
	OnFinish()
	Finished() bool
}
