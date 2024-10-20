package ux

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type ProgressBar struct {
	progressBar     material.ProgressBarStyle
	initProgress    float32
	currentProgress float32
}

func NewProgressBar(initProgress float32) *ProgressBar {
	progressBar := &ProgressBar{
		initProgress:    initProgress,
		currentProgress: initProgress,
		progressBar:     material.ProgressBar(th.Theme, initProgress),
	}
	progressBar.progressBar.Color = th.Color.ProgressBarColor
	return progressBar
}

func (p *ProgressBar) SetProgress(progress float32) {
	p.currentProgress = progress
	p.progressBar.Progress = p.currentProgress
}

func (p *ProgressBar) Layout(gtx layout.Context) layout.Dimensions {
	return p.progressBar.Layout(gtx)
}
