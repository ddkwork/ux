package ux

import (
	"image"
	"unicode"

	"github.com/ddkwork/ux/resources/colors"
	"github.com/ddkwork/ux/resources/icons"

	"github.com/ddkwork/ux/widget/material"

	"gioui.org/font"
	"gioui.org/gesture"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/safemap"
)

var (
	horizontalInset = layout.Inset{Left: unit.Dp(2)}
	verticalInset   = layout.Inset{Top: unit.Dp(2)}
	horizontalFlex  = layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}
	verticalFlex    = layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}
)

type TabView struct {
	Axis         layout.Axis
	list         layout.List
	tabItems     []*TabItem
	index        int
	headerLength int
}

type TabItem struct {
	// Title of the tab.
	title string
	// Main part of the tab content.
	Content layout.Widget
	// Title padding of the tab item.
	Inset     layout.Inset
	direction layout.Direction
	click     gesture.Click
	hovering  bool
	selected  bool
	btn       widget.Clickable
	index     int
	// number of characters to show in the tab title
	maxTitleWidth    int
	onSelectedChange func(int)

	Identifier string

	closable       bool
	CloseClickable widget.Clickable

	isDataChanged bool
	onClose       func(t *TabItem)
	isClosed      bool

	Meta *safemap.M[string, string]
}

func (t *TabItem) SetOnSelectedChange(onSelectedChange func(int)) *TabItem {
	t.onSelectedChange = onSelectedChange
	return t
}

func (t *TabItem) SetClosable(closable bool) *TabItem {
	t.closable = closable
	return t
}

func (t *TabItem) Update(gtx layout.Context) bool {
	for {
		event, ok := gtx.Event(
			pointer.Filter{Target: t, Kinds: pointer.Enter | pointer.Leave},
		)
		if !ok {
			break
		}

		switch e := event.(type) {
		case pointer.Event:
			switch e.Kind {
			case pointer.Enter:
				t.hovering = true
			case pointer.Leave:
				t.hovering = false
			case pointer.Cancel:
				t.hovering = false
			}
		}
	}

	var clicked bool
	for {
		e, ok := t.click.Update(gtx.Source)
		if !ok {
			break
		}
		if e.Kind == gesture.KindClick {
			clicked = true
			t.selected = true
		}
	}

	return clicked
}

func (t *TabItem) LayoutTitle(gtx layout.Context) layout.Dimensions {
	t.Update(gtx)

	macro := op.Record(gtx.Ops)
	dims := t.layoutTitle(gtx)
	call := macro.Stop()

	rect := clip.Rect(image.Rectangle{Max: dims.Size})
	defer rect.Push(gtx.Ops).Pop()

	t.click.Add(gtx.Ops)
	// register tag
	event.Op(gtx.Ops, t)
	call.Add(gtx.Ops)

	return dims
}

func (t *TabItem) layoutTitle(gtx layout.Context) layout.Dimensions {
	return t.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return t.direction.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			if t.btn.Clicked(gtx) {
				if t.onSelectedChange != nil {
					go t.onSelectedChange(t.index)
					gtx.Execute(op.InvalidateCmd{})
				}
			}

			if t.closable && t.onClose != nil && t.CloseClickable.Clicked(gtx) {
				t.onClose(t)
				gtx.Execute(op.InvalidateCmd{})
			}

			if t.btn.Hovered() {
				paint.FillShape(gtx.Ops, th.Palette.ContrastBg, clip.Rect{Max: gtx.Constraints.Min}.Op())
			}

			var tabWidth int
			return layout.Stack{Alignment: layout.S}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					var dims layout.Dimensions
					if t.closable {
						dims = material.Clickable(gtx, &t.btn, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									maxTitleWidth := 20
									label := material.Label(th, unit.Sp(13), ellipticalTruncate(t.title, maxTitleWidth))
									label.Color = th.Color.DefaultTextWhiteColor
									if t.btn.Hovered() {
										// label.Font.Weight = font.Bold
										label.TextSize++
										label.Color = colors.Red100
									}
									return layout.UniformInset(unit.Dp(7)).Layout(gtx, label.Layout)
								}),

								layout.Rigid(layout.Spacer{Width: unit.Dp(2)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									// bkColor := color.NRGBA{}
									// hoveredColor := Hovered(bkColor)
									if t.btn.Hovered() {
										// bkColor = hoveredColor
									}
									// iconColor := th.ContrastFg
									// closeIcon := NavigationCloseIcon
									// iconSize := unit.Dp(16)
									padding := unit.Dp(4)
									if t.isDataChanged {
										// yellow
										// iconColor = color.NRGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff}
										icons.NavigationCloseIcon = icons.ImageLensIcon
										// iconSize = unit.Dp(10)
										padding = unit.Dp(8)
									}

									//ib := &IconButton2{//todo
									//	Icon:                 closeIcon,
									//	Color:                iconColor,
									//	BackgroundColor:      bkColor,
									//	BackgroundColorHover: hoveredColor,
									//	Size:                 iconSize,
									//	Clickable:            &t.CloseClickable,
									//}
									return layout.UniformInset(padding).Layout(gtx,
										func(gtx layout.Context) layout.Dimensions {
											return Button(&t.CloseClickable, icons.NavigationCloseIcon, "").Layout(gtx)
										},
									)
								}),
							)
						})
					} else {
						dims = material.Clickable(gtx, &t.btn, func(gtx layout.Context) layout.Dimensions {
							label := material.Label(th, unit.Sp(13), t.title)
							label.Color = th.Color.DefaultTextWhiteColor
							if t.btn.Hovered() {
								label.Font.Weight = font.Bold
								label.Color = colors.Red100
							}
							return layout.UniformInset(unit.Dp(7)).Layout(gtx,
								label.Layout,
							)
						})
					}
					tabWidth = dims.Size.X
					return dims
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					if !t.selected {
						return layout.Dimensions{}
					}
					tabHeight := gtx.Dp(unit.Dp(4))
					tabRect := image.Rect(0, 0, tabWidth, tabHeight)
					paint.FillShape(gtx.Ops, colors.ColorPink, clip.Rect(tabRect).Op())
					return layout.Dimensions{
						Size: image.Point{X: tabWidth, Y: tabHeight},
					}
				}),
			)
		})
	})
}

func (t *TabItem) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Inset{Top: 4}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return t.Content(gtx)
	})
}

func (v *TabView) Layout(gtx layout.Context) layout.Dimensions {
	v.Update(gtx)

	if len(v.tabItems) <= 0 {
		return layout.Dimensions{}
	}

	v.tabItems[v.index].index = v.index

	var direction layout.Direction
	var flex layout.Flex
	var tabAlign layout.Direction
	if v.Axis == layout.Horizontal {
		direction = layout.N
		flex = horizontalFlex
		tabAlign = layout.N
	} else {
		direction = layout.N
		flex = verticalFlex
		tabAlign = layout.W
	}

	return flex.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return direction.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				v.list.Axis = v.Axis
				v.list.Alignment = layout.Start
				listDims := v.list.Layout(gtx, len(v.tabItems), func(gtx layout.Context, index int) layout.Dimensions {
					item := v.tabItems[index]
					item.direction = tabAlign

					if index == 0 {
						return item.LayoutTitle(gtx)
					}

					if v.Axis == layout.Horizontal {
						return horizontalInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return item.LayoutTitle(gtx)
						})
					} else {
						return verticalInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return item.LayoutTitle(gtx)
						})
					}
				})

				if v.Axis == layout.Horizontal {
					v.headerLength = listDims.Size.X
				} else {
					v.headerLength = listDims.Size.Y
				}
				return listDims
			})
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if v.Axis == layout.Horizontal {
				gtx.Constraints.Min.X = v.headerLength
			} else {
				gtx.Constraints.Min.Y = v.headerLength
			}
			return Divider(v.Axis, unit.Dp(0.5)).Layout(gtx)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return v.tabItems[v.index].Layout(gtx)
		}),
	)
}

func (v *TabView) Update(gtx layout.Context) {
	for idx, item := range v.tabItems {
		if item.Update(gtx) {
			// unselect last item
			lastItem := v.tabItems[v.index]
			if lastItem != nil && idx != v.index {
				lastItem.selected = false
			}
			v.index = idx
		}
		if v.index == idx && !item.selected {
			item.selected = true
		}
	}
}

func NewTabView(axis layout.Axis, item ...*TabItem) *TabView {
	return &TabView{
		Axis:     axis,
		tabItems: item,
	}
}

func NewTabItem(title string, content layout.Widget) *TabItem {
	return &TabItem{
		title:     title,
		Content:   content,
		Inset:     layout.UniformInset(0),
		direction: 0,
		click:     gesture.Click{},
		hovering:  false,
		selected:  false,
	}
}

func ellipticalTruncate(text string, maxLen int) string {
	if maxLen == 0 {
		return text
	}

	lastSpaceIx := maxLen
	l := 0
	for i, r := range text {
		if unicode.IsSpace(r) {
			lastSpaceIx = i
		}
		l++
		if l > maxLen {
			return text[:lastSpaceIx] + "..."
		}
	}
	return text
}

func (v *TabView) SelectedTab() *TabItem {
	if len(v.tabItems) == 0 {
		return nil
	}
	return v.tabItems[v.index]
}

func (v *TabView) AddTab(tab *TabItem) int {
	v.tabItems = append(v.tabItems, tab)
	return len(v.tabItems) - 1
}

func (v *TabView) RemoveTabByID(id string) {
	tab := v.findTabByID(id)
	if tab == nil {
		return
	}

	tab.isClosed = true
}

func (v *TabView) findTabByID(id string) *TabItem {
	for _, t := range v.tabItems {
		if t.Identifier == id {
			return t
		}
	}
	return nil
}

func (v *TabView) SetSelectedByID(id string) {
	for i, t := range v.tabItems {
		if t.Identifier == id {
			v.index = i
			return
		}
	}
}

func (v *TabView) SetTabs(items []*TabItem) {
	v.tabItems = items
}

func (t *TabItem) SetMaxTitleWidth(maxWidth int) {
	t.maxTitleWidth = maxWidth
}

func (t *TabItem) Selected() bool {
	return t.selected
}

func (t *TabItem) SetDataChanged(changed bool) *TabItem {
	t.isDataChanged = changed
	return t
}

func (t *TabItem) IsDataChanged() bool {
	return t.isDataChanged
}

func (t *TabItem) SetIdentifier(id string) {
	t.Identifier = id
}

func (t *TabItem) SetSelected(index int) {
	t.index = index
}

func (t *TabItem) SetOnClose(f func(t *TabItem)) {
	t.onClose = f
}
