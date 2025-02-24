package ux

import (
	"slices"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/ddkwork/golibrary/stream"
)

type StructView struct {
	keys    []string
	element []layout.Widget
	widget.List
}

func NewStructView(data any, marshal func() (elems []CellData)) *StructView {
	visibleFields := stream.ReflectVisibleFields(data)
	keys := make([]string, len(visibleFields))
	values := make([]string, len(visibleFields))
	elems := marshal()
	element := make([]layout.Widget, len(visibleFields))
	for i, field := range visibleFields {
		keys[i] = field.Name
		values[i] = elems[i].Text
		element[i] = NewInput(field.Name, elems[i].Text).Layout
	}
	return &StructView{
		keys:    keys,
		element: element,
	}
}

func (s *StructView) InsertAt(index int, label string, widget layout.Widget) *StructView {
	s.keys = slices.Insert(s.keys, index, label)
	s.element = slices.Insert(s.element, index, widget)
	return s
}

func (s *StructView) Add(label string, widget layout.Widget) *StructView {
	s.keys = append(s.keys, label)
	s.element = append(s.element, widget)
	return s
}

func (s *StructView) maxLabelWidth(gtx layout.Context) unit.Dp {
	originalConstraints := gtx.Constraints
	maxWidth := unit.Dp(0)
	for _, data := range s.keys {
		currentWidth := LabelWidth(gtx, data)
		if currentWidth > maxWidth {
			maxWidth = currentWidth
		}
	}
	gtx.Constraints = originalConstraints
	return maxWidth
}

func (s *StructView) Layout(gtx layout.Context) layout.Dimensions {
	rowHeight := unit.Dp(4)
	maxLabelWidth := s.maxLabelWidth(gtx)
	var elements []layout.Widget
	for i := range s.keys {
		elements = append(elements, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.Y = gtx.Dp(rowHeight)
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Dp(maxLabelWidth)
					gtx.Constraints.Max.X = gtx.Dp(maxLabelWidth)
					return layout.Inset{Top: rowHeight}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { // 占据label宽度右对齐
								return layout.Spacer{}.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Top: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									labelText := s.keys[i]
									if labelText != "" {
										labelText = labelText + "："
									}
									return material.Label(th.Theme, th.TextSize, labelText).Layout(gtx)
								})
							}),
						)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Top: rowHeight}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return s.element[i](gtx)
					})
				}),
			)
		})
	}
	s.List.Axis = layout.Vertical
	return material.List(th.Theme, &s.List).Layout(gtx, len(elements), func(gtx layout.Context, index int) layout.Dimensions {
		return elements[index](gtx)
	})
}
