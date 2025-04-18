package main

import (
	"gioui.org/layout"
	"github.com/ddkwork/ux"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Page struct {
	*ux.PopupTest
}

// New constructs a Page with the provided router.
func New() *Page {
	return &Page{
		PopupTest: ux.NewPopupTest(100),
	}
}

var th = ux.NewTheme()

func (p *Page) Layout(gtx C) D {
	return p.PopupTest.Layout(gtx, nil)
}
