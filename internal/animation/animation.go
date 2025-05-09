package animation

import (
	"time"

	"github.com/ddkwork/ux/internal/animation/gween"

	"gioui.org/layout"
	"gioui.org/op"
)

type AnimState struct {
	Value    float32
	Active   bool
	Finished bool
}

type Animation struct {
	Sequence *gween.Sequence

	active        bool
	stop          bool
	delay         time.Duration
	lastFrameTime time.Time
	startTime     time.Time
}

func NewAnimation(startImmediately bool, sequence *gween.Sequence) *Animation {
	return &Animation{
		Sequence: sequence,
		stop:     !startImmediately,
		active:   startImmediately,
	}
}

func (animation *Animation) Update(gtx layout.Context) AnimState {
	var dt time.Duration

	if animation.startTime.IsZero() {
		animation.startTime = gtx.Now
	}

	if !animation.lastFrameTime.IsZero() {
		dt = gtx.Now.Sub(animation.lastFrameTime)
	}

	if gtx.Now.Sub(animation.startTime) > animation.delay && !animation.stop {
		animation.lastFrameTime = gtx.Now
	}

	seconds := float32(dt.Seconds())
	value, _, finished := animation.Sequence.Update(seconds)

	if !animation.stop {
		gtx.Execute(op.InvalidateCmd{})
	}

	if finished {
		animation.stop = true
	}

	return AnimState{
		Value:    value,
		Active:   animation.active,
		Finished: finished,
	}
}

func (animation *Animation) Start() *Animation {
	if animation.stop {
		animation.Reset()
		animation.stop = false
		animation.active = true
	}

	return animation
}

func (animation *Animation) StartWithDelay(delay time.Duration) *Animation {
	if animation.stop {
		animation.Reset()
		animation.delay = delay
		animation.stop = false
		animation.active = true
	}

	return animation
}

func (animation *Animation) Resume() *Animation {
	animation.lastFrameTime = time.Time{}
	animation.stop = false
	return animation
}

func (animation *Animation) Pause() *Animation {
	animation.lastFrameTime = time.Time{}
	animation.stop = true
	return animation
}

func (animation *Animation) Reset() *Animation {
	animation.active = false
	animation.delay = 0
	animation.lastFrameTime = time.Time{}
	animation.startTime = time.Time{}
	animation.stop = true
	animation.Sequence.Reset()
	return animation
}
