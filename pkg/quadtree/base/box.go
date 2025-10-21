package quadcore

import (
	"sort"

	"github.com/kjkrol/gokg/pkg/geometry"
)

// Box represents a rectangular area in a 2D space defined by its top-left and bottom-right corners.
// It also includes the center point of the Box for convenience.
//
// Fields:
// - topLeft: The top-left corner of the Box as a vector of integers.
// - bottomRight: The bottom-right corner of the Box as a vector of integers.
// - center: The center point of the Box as a vector of integers.
type Box[T geometry.SupportedNumeric] struct {
	TopLeft     geometry.Vec[T]
	BottomRight geometry.Vec[T]
	Center      geometry.Vec[T]
}

func NewBox[T geometry.SupportedNumeric](topLeft geometry.Vec[T], bottomRight geometry.Vec[T]) Box[T] {
	centerX := (topLeft.X + bottomRight.X) / 2
	centerY := (topLeft.Y + bottomRight.Y) / 2
	center := geometry.Vec[T]{X: centerX, Y: centerY}
	Box := Box[T]{topLeft, bottomRight, center}
	return Box
}

func BuildBox[T geometry.SupportedNumeric](center geometry.Vec[T], d T) Box[T] {
	topLeft := geometry.Vec[T]{X: center.X - d, Y: center.Y - d}
	bottomRight := geometry.Vec[T]{X: center.X + d, Y: center.Y + d}
	return NewBox[T](topLeft, bottomRight)
}

func (b *Box[T]) Split() [4]Box[T] {
	return [4]Box[T]{
		NewBox(b.TopLeft, b.Center), // top left
		NewBox(geometry.Vec[T]{X: b.Center.X, Y: b.TopLeft.Y}, geometry.Vec[T]{X: b.BottomRight.X, Y: b.Center.Y}), // top right
		NewBox(geometry.Vec[T]{X: b.TopLeft.X, Y: b.Center.Y}, geometry.Vec[T]{X: b.Center.X, Y: b.BottomRight.Y}), // bottom left
		NewBox(b.Center, b.BottomRight), // bottom right
	}
}

func (b Box[T]) ContainsBox(other Box[T]) bool {
	return other.TopLeft.X >= b.TopLeft.X &&
		other.TopLeft.Y >= b.TopLeft.Y &&
		other.BottomRight.X <= b.BottomRight.X &&
		other.BottomRight.Y <= b.BottomRight.Y
}

func (b Box[T]) Expand(margin T) Box[T] {
	return NewBox(
		geometry.Vec[T]{X: b.TopLeft.X - margin, Y: b.TopLeft.Y - margin},
		geometry.Vec[T]{X: b.BottomRight.X + margin, Y: b.BottomRight.Y + margin},
	)
}

//-------------------------------------------------------------------------

func (b Box[T]) Intersects(other Box[T]) bool {
	// x axis check
	xIntersects := axisIntersection(b, other, func(v geometry.Vec[T]) T { return v.X })
	if !xIntersects {
		return false
	}
	// y axis check
	yIntersects := axisIntersection(b, other, func(v geometry.Vec[T]) T { return v.Y })
	return yIntersects
}

func axisIntersection[T geometry.SupportedNumeric](aa, bb Box[T], axisValue func(geometry.Vec[T]) T) bool {
	aa, bb = SortBy(aa, bb, axisValue)
	noIntersection := axisValue(aa.TopLeft) < axisValue(bb.BottomRight) && axisValue(aa.BottomRight) < axisValue(bb.TopLeft)
	return !noIntersection
}

func SortBy[T geometry.SupportedNumeric](a, b Box[T], axisValue func(geometry.Vec[T]) T) (aa, bb Box[T]) {
	if axisValue(a.TopLeft) < axisValue(b.TopLeft) {
		aa = a
		bb = b
	} else {
		aa = b
		bb = a
	}
	return
}

//-------------------------------------------------------------------------

func (b Box[T]) intersectsAny(others []Box[T]) bool {
	for _, wrapped := range others {
		if b.Intersects(wrapped) {
			return true
		}
	}
	return false
}

func WrapBoxCyclic[T geometry.SupportedNumeric](
	b Box[T],
	size geometry.Vec[T],
	contains func(geometry.Vec[T]) bool,
) []Box[T] {
	var wrappedBoxes []Box[T]

	// Predefined offset values for wrapping
	offsets := []geometry.Vec[T]{
		{X: size.X, Y: 0},      // Shift right
		{X: 0, Y: size.Y},      // Shift down
		{X: size.X, Y: size.Y}, // Shift right-down
	}

	vecMath := geometry.VectorMathByType[T]()

	// Generate wrapped versions for each offset
	for _, offset := range offsets {
		wrapped := Box[T]{
			TopLeft:     b.TopLeft,
			BottomRight: b.BottomRight,
			Center:      b.Center,
		}
		vecMath.Wrap(&wrapped.TopLeft, offset)
		vecMath.Wrap(&wrapped.BottomRight, offset)
		vecMath.Wrap(&wrapped.Center, offset)

		if contains(wrapped.TopLeft) || contains(wrapped.BottomRight) {
			wrappedBoxes = append(wrappedBoxes, wrapped)
		}
	}

	return wrappedBoxes
}

//-------------------------------------------------------------------------

func SortBoxes[T geometry.SupportedNumeric](boxes []Box[T], axis func(geometry.Vec[T]) T) {
	sort.Slice(boxes, func(i, j int) bool {
		return axis(boxes[i].TopLeft) < axis(boxes[j].TopLeft)
	})
}
