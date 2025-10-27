package quadtree

import (
	"github.com/kjkrol/gokg/pkg/geometry"
	"github.com/kjkrol/gokg/pkg/geometry/spatial"
)

func (t *QuadTree[T]) probe(spatialItem spatial.Spatial[T], margin T) []spatial.Rectangle[T] {
	probe := spatialItem.Bounds().Expand(margin)
	rectangles := []spatial.Rectangle[T]{probe}
	if t.plane.Name() == "cyclic" {
		rectangles = append(rectangles, t.wrapRectangleCyclic(probe)...)
	}
	return rectangles
}

func (t *QuadTree[T]) wrapRectangleCyclic(
	r spatial.Rectangle[T],
) []spatial.Rectangle[T] {
	var wrappedRectangles []spatial.Rectangle[T]
	var size = t.plane.Size()
	var contains = t.plane.Contains

	// Predefined offset values for wrapping
	offsets := []spatial.Vec[T]{
		{X: size.X, Y: 0},      // Shift right
		{X: 0, Y: size.Y},      // Shift down
		{X: size.X, Y: size.Y}, // Shift right-down
	}

	vecMath := geometry.VectorMathByType[T]()

	// Generate wrapped versions for each offset
	for _, offset := range offsets {
		wrapped := spatial.Rectangle[T]{
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
