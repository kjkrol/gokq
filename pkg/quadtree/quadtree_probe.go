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

func (t *QuadTree[T]) wrapAABBCyclic(probe geometry.AABB[T]) []geometry.AABB[T] {
	return geometry.GenerateBoundaryFragments(
		probe.TopLeft,
		t.plane.Size(),
		t.plane.VectorMath(),
		func(offset geometry.Vec[T]) (geometry.AABB[T], geometry.AABB[T], bool) {
			wrapped := geometry.AABB[T]{
				TopLeft:     probe.TopLeft.Add(offset),
				BottomRight: probe.BottomRight.Add(offset),
				Center:      probe.Center.Add(offset),
			}
			return wrapped, wrapped, true
		})
}
