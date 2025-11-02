package quadtree

import (
	"github.com/kjkrol/gokg/pkg/geometry"
)

type QuadTreeOption[T geometry.SupportedNumeric] func(*QuadTree[T])

func WithMaxDepth[T geometry.SupportedNumeric](depth int) QuadTreeOption[T] {
	return func(qt *QuadTree[T]) {
		qt.appender.maxDepth = depth
	}
}
