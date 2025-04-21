package animation

import (
	gween2 "github.com/ddkwork/ux/internal/animation/gween"
	"github.com/ddkwork/ux/internal/animation/gween/ease"
)

func ButtonClick() *Animation {
	sequence := gween2.NewSequence(
		gween2.New(1, .98, .1, ease.Linear),
		gween2.New(.98, 1, .4, ease.OutBounce),
	)
	return NewAnimation(false, sequence)
}

func ButtonEnter() *Animation {
	sequence := gween2.NewSequence(
		gween2.New(1, .99, .1, ease.Linear),
	)
	return NewAnimation(false, sequence)
}

func ButtonLeave() *Animation {
	sequence := gween2.NewSequence(
		gween2.New(.99, 1, .1, ease.Linear),
	)
	return NewAnimation(false, sequence)
}
