package quadtree

import "github.com/kjkrol/gokg/pkg/geometry"

type QuadTreeOption[T geometry.SupportedNumeric] func(*QuadTree[T])

func WithDistance[T geometry.SupportedNumeric](d geometry.Distance[T]) QuadTreeOption[T] {
	return func(qt *QuadTree[T]) {
		qt.distance = d
	}
}

func WithMaxDepth[T geometry.SupportedNumeric](depth int) QuadTreeOption[T] {
	return func(qt *QuadTree[T]) {
		qt.maxDepth = depth
	}
}
