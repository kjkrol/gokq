package quadtree

import (
	"github.com/kjkrol/gokg/pkg/geometry"
)

// Box represents a rectangular area in a 2D space defined by its top-left and bottom-right corners.
// It also includes the center point of the box for convenience.
//
// Fields:
// - topLeft: The top-left corner of the box as a vector of integers.
// - bottomRight: The bottom-right corner of the box as a vector of integers.
// - center: The center point of the box as a vector of integers.
type Box struct {
	topLeft     geometry.Vec[int]
	bottomRight geometry.Vec[int]
	center      geometry.Vec[int]
}

func buildBox(center geometry.Vec[int], d int) Box {
	return Box{
		geometry.Vec[int]{X: center.X - d, Y: center.Y - d},
		geometry.Vec[int]{X: center.X + d, Y: center.Y + d},
		center,
	}
}

// NewBox creates a new Box given the top-left and bottom-right coordinates.
// It calculates the center point of the box and returns a Box struct.
//
// Parameters:
//   - topLeft: a Vec[int] representing the top-left corner of the box.
//   - bottomRight: a Vec[int] representing the bottom-right corner of the box.
//
// Returns:
//   - A Box struct with the specified top-left and bottom-right corners, and the calculated center point.
func NewBox(topLeft geometry.Vec[int], bottomRight geometry.Vec[int]) Box {
	centerX := (topLeft.X + bottomRight.X) / 2
	centerY := (topLeft.Y + bottomRight.Y) / 2
	center := geometry.Vec[int]{X: centerX, Y: centerY}
	return Box{topLeft, bottomRight, center}
}

func (b *Box) split() [4]Box {
	return [4]Box{
		NewBox(b.topLeft, b.center), // top left
		NewBox(geometry.Vec[int]{X: b.center.X, Y: b.topLeft.Y}, geometry.Vec[int]{X: b.bottomRight.X, Y: b.center.Y}), // top right
		NewBox(geometry.Vec[int]{X: b.topLeft.X, Y: b.center.Y}, geometry.Vec[int]{X: b.center.X, Y: b.bottomRight.Y}), // bottom left
		NewBox(b.center, b.bottomRight), // bottom right
	}
}

func (b *Box) intersects(other *Box) bool {
	// x axis check
	xIntersects := axisIntersection(b, other, func(v geometry.Vec[int]) int { return v.X })
	if !xIntersects {
		return false
	}
	// y axis check
	yIntersects := axisIntersection(b, other, func(v geometry.Vec[int]) int { return v.Y })
	return yIntersects
}

func axisIntersection(aa, bb *Box, axisValue func(geometry.Vec[int]) int) bool {
	aa, bb = sortBy(aa, bb, axisValue)
	noIntersection := axisValue(aa.topLeft) < axisValue(bb.bottomRight) && axisValue(aa.bottomRight) < axisValue(bb.topLeft)
	return !noIntersection
}

func sortBy(a, b *Box, axisValue func(geometry.Vec[int]) int) (aa, bb *Box) {
	if axisValue(a.topLeft) < axisValue(b.topLeft) {
		aa = a
		bb = b
	} else {
		aa = b
		bb = a
	}
	return
}
