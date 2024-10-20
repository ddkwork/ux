package ux

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/outlay"
)

type (
	FlowWrap struct {
		Cards []Card
		widget.List
		Wrap       outlay.FlowWrap
		Contextual interface{}
		Loaded     bool
	}
)

func (gr *FlowWrap) Layout(gtx C) D {
	if !gr.Loaded {
		gr.Wrap.Alignment = layout.Middle
		gr.List.Axis = layout.Vertical
		gr.List.Alignment = layout.Middle

		gr.Cards = append(gr.Cards,
			Card{
				Name: "Card 1",
			},
			Card{
				Name: "Card 2",
			},
			Card{
				Name: "Card 3",
			},
			Card{
				Name: "Card 4",
			},
			Card{
				Name: "Card 5",
			},
			Card{
				Name: "Card 6",
			},
		)
		gr.Loaded = true
	}

	return material.List(th.Theme, &gr.List).Layout(gtx, 1, func(gtx C, _ int) D {
		return layout.Flex{Spacing: layout.SpaceSides}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return gr.Wrap.Layout(gtx, len(gr.Cards), func(gtx C, i int) D {
					var content D

					// copy only this specific card
					if gr.Cards[i].copyToClipBtn.Clicked(gtx) {
						//res, _ := json.MarshalIndent(gr.cards[i], "", "\t")
						//clipboard.WriteOp{
						//	Content: string(res),
						//}.Add(gtx.Ops)
						//globals.ClipBoardVal = string(res)
					}

					if gr.Cards[i].btn.Clicked(gtx) {
						gr.Contextual = gr.Cards[i] // interface to assert type when enabling ContextualAppBar
						// gr.cards[i].IsCtxtActive = true
						// op.InvalidateOp{}.Add(gtx.Ops)
					}

					if gr.Cards[i].selectBtn.Clicked(gtx) {
						// data.Cached[i].Selected = true
						// op.InvalidateOp{}.Add(gtx.Ops)
					} else if gr.Cards[i].deselectBtn.Clicked(gtx) {
						// data.Cached[i].Selected = false
						// op.InvalidateOp{}.Add(gtx.Ops)
					}

					// if gr.cards[i].IsSearchedFor && gr.cards[i].IsActiveContinent {
					content = layout.Inset{
						Top:    unit.Dp(15),
						Bottom: unit.Dp(15),
						Left:   unit.Dp(25),
						Right:  unit.Dp(25),
					}.Layout(gtx, func(gtx C) D {
						return gr.Cards[i].LayCard(gtx)
					})
					//}
					return content
				})
			}))
	})
}
