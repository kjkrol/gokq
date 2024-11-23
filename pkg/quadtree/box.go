package quadtree

import (
	"github.com/kjkrol/gokg/pkg/geometry"
)

// box represents a rectangular area in a 2D space defined by its top-left and bottom-right corners.
// It also includes the center point of the box for convenience.
//
// Fields:
// - topLeft: The top-left corner of the box as a vector of integers.
// - bottomRight: The bottom-right corner of the box as a vector of integers.
// - center: The center point of the box as a vector of integers.
type box[T geometry.SupportedNumeric] struct {
	topLeft     geometry.Vec[T]
	bottomRight geometry.Vec[T]
	center      geometry.Vec[T]
}

func newBox[T geometry.SupportedNumeric](topLeft geometry.Vec[T], bottomRight geometry.Vec[T]) box[T] {
	centerX := (topLeft.X + bottomRight.X) / 2
	centerY := (topLeft.Y + bottomRight.Y) / 2
	center := geometry.Vec[T]{X: centerX, Y: centerY}
	box := box[T]{topLeft, bottomRight, center}
	return box
}

func buildBox[T geometry.SupportedNumeric](center geometry.Vec[T], d T) box[T] {
	topLeft := geometry.Vec[T]{X: center.X - d, Y: center.Y - d}
	bottomRight := geometry.Vec[T]{X: center.X + d, Y: center.Y + d}
	return newBox[T](topLeft, bottomRight)
}

func (b *box[T]) split() [4]box[T] {
	return [4]box[T]{
		newBox(b.topLeft, b.center), // top left
		newBox(geometry.Vec[T]{X: b.center.X, Y: b.topLeft.Y}, geometry.Vec[T]{X: b.bottomRight.X, Y: b.center.Y}), // top right
		newBox(geometry.Vec[T]{X: b.topLeft.X, Y: b.center.Y}, geometry.Vec[T]{X: b.center.X, Y: b.bottomRight.Y}), // bottom left
		newBox(b.center, b.bottomRight), // bottom right
	}
}

//-------------------------------------------------------------------------

func (b box[T]) intersects(other box[T]) bool {
	// x axis check
	xIntersects := axisIntersection(b, other, func(v geometry.Vec[T]) T { return v.X })
	if !xIntersects {
		return false
	}
	// y axis check
	yIntersects := axisIntersection(b, other, func(v geometry.Vec[T]) T { return v.Y })
	return yIntersects
}

func axisIntersection[T geometry.SupportedNumeric](aa, bb box[T], axisValue func(geometry.Vec[T]) T) bool {
	aa, bb = sortBy(aa, bb, axisValue)
	noIntersection := axisValue(aa.topLeft) < axisValue(bb.bottomRight) && axisValue(aa.bottomRight) < axisValue(bb.topLeft)
	return !noIntersection
}

func sortBy[T geometry.SupportedNumeric](a, b box[T], axisValue func(geometry.Vec[T]) T) (aa, bb box[T]) {
	if axisValue(a.topLeft) < axisValue(b.topLeft) {
		aa = a
		bb = b
	} else {
		aa = b
		bb = a
	}
	return
}

//-------------------------------------------------------------------------

func (b box[T]) intersectsAny(others []box[T]) bool {
	for _, wrapped := range others {
		if b.intersects(wrapped) {
			return true
		}
	}
	return false
}

// wrapBoxCyclic wraps the given box in a cyclic space defined by the size vector.
// It returns a slice of boxes that includes the original box and its wrapped versions
// based on predefined offset values.
//
// Parameters:
// - b: The original box to be wrapped.
// - size: The size of the cyclic space as a vector of integers.
// - vecMath: A utility for vector math operations.
//
// Returns:
// A slice of boxes including the original and its wrapped versions.
func wrapBoxCyclic[T geometry.SupportedNumeric](b box[T], size geometry.Vec[T], vecMath geometry.VectorMath[T]) []box[T] {
	var wrappedBoxes []box[T]

	// Append original box directly
	wrappedBoxes = append(wrappedBoxes, b)

	// Predefined offset values for wrapping
	offsets := []geometry.Vec[T]{
		{X: size.X, Y: 0},      // Shift right
		{X: 0, Y: size.Y},      // Shift down
		{X: size.X, Y: size.Y}, // Shift right-down
	}

	// Generate wrapped versions for each offset
	for _, offset := range offsets {
		wrapped := box[T]{
			topLeft:     b.topLeft,
			bottomRight: b.bottomRight,
			center:      b.center,
		}
		vecMath.Wrap(&wrapped.topLeft, offset)
		vecMath.Wrap(&wrapped.bottomRight, offset)
		vecMath.Wrap(&wrapped.center, offset)
		wrappedBoxes = append(wrappedBoxes, wrapped)
	}

	return wrappedBoxes
}

//-------------------------------------------------------------------------
