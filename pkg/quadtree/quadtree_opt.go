package quadtree

import (
	"github.com/kjkrol/gokg/pkg/geometry"
	"github.com/kjkrol/gokg/pkg/geometry/spatial"
)

type QuadTreeOption[T spatial.SupportedNumeric] func(*QuadTree[T])

func WithDistance[T spatial.SupportedNumeric](d geometry.Distance[T]) QuadTreeOption[T] {
	return func(qt *QuadTree[T]) {
		qt.distance = d
	}
}

func WithMaxDepth[T spatial.SupportedNumeric](depth int) QuadTreeOption[T] {
	return func(qt *QuadTree[T]) {
		qt.maxDepth = depth
	}
}
