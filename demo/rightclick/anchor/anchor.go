package anchor

import (
	"fmt"

	"gioui.org/f32"
)

// Anchor is an opaque reference to a global coordinate position.
// It can be provided to methods in this package as a reference
// to a global coordinate.
type Anchor struct {
	point f32.Point
}

// AnchorFrom wraps an f32.Point within an Anchor, preventing the
// coordinates within from being used in any way other than determining
// an offset using the OffsetWithin method.
func AnchorFrom(point f32.Point) Anchor {
	return Anchor{point}
}

// String is provided for debugging purposes.
func (a Anchor) String() string {
	return fmt.Sprintf("anchor (%f,%f)", a.point.X, a.point.Y)
}

// OffsetWithin returns an offset that will allow a widget of size contentSize
// to be rendered within the provided bounds. The offset is as close as possible
// to the coordinates wrapped within the anchor.
func (a Anchor) OffsetWithin(contentSize, bounds f32.Point) f32.Point {
	var offset f32.Point = a.point
	if contentSize.X+a.point.X > bounds.X {
		offset.X = bounds.X - contentSize.X
	}
	if contentSize.Y+a.point.Y > bounds.Y {
		offset.Y = bounds.Y - contentSize.Y
	}
	return offset
}
