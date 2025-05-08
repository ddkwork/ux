package ux

import (
	"time"
)

const (
	alertSpeed           = 10 // units: fadeLevels per second
	defaultAlertDuration = time.Second * 3
)

type (
	Alert struct {
		Name      string
		Priority  AlertPriority
		Message   string
		Duration  time.Duration
		FadeLevel float64
	}

	AlertPriority  int
	AlertYieldFunc func(index int, alert Alert) bool
	// Alerts         Model
)

const (
	None AlertPriority = iota
	Info
	Warning
	Error
)
