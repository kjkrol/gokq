package quadtree

import (
	"github.com/kjkrol/gokg/pkg/geometry"
)

func (t *QuadTree[T]) probe(aabb geometry.AABB[T], margin T) []geometry.AABB[T] {
	probe := aabb.Expand(margin)
	rectangles := []geometry.AABB[T]{probe}
	if t.plane.Name() == "cyclic" {
		rectangles = append(rectangles, t.wrapAABBCyclic(probe)...)
	}
	return rectangles
}

func (t *QuadTree[T]) wrapAABBCyclic(
	r geometry.AABB[T],
) []geometry.AABB[T] {
	var wrappedRectangles []geometry.AABB[T]
	var size = t.plane.Size()
	var contains = t.plane.Contains

	// Predefined offset values for wrapping
	offsets := []geometry.Vec[T]{
		{X: size.X, Y: 0},      // Shift right
		{X: 0, Y: size.Y},      // Shift down
		{X: size.X, Y: size.Y}, // Shift right-down
	}

	vecMath := geometry.VectorMathByType[T]()

	// Generate wrapped versions for each offset
	for _, offset := range offsets {
		wrapped := geometry.AABB[T]{
			TopLeft:     r.TopLeft,
			BottomRight: r.BottomRight,
			Center:      r.Center,
		}
		vecMath.Wrap(&wrapped.TopLeft, offset)
		vecMath.Wrap(&wrapped.BottomRight, offset)
		vecMath.Wrap(&wrapped.Center, offset)

		if contains(wrapped.TopLeft) || contains(wrapped.BottomRight) {
			wrappedRectangles = append(wrappedRectangles, wrapped)
		}
	}

	return wrappedRectangles
}
