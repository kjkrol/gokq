package quadtree

import (
	"github.com/kjkrol/gokg/pkg/geometry"
)

type QuadTreeOption[T geometry.SupportedNumeric, K comparable] func(*QuadTree[T, K])

func WithMaxDepth[T geometry.SupportedNumeric, K comparable](depth int) QuadTreeOption[T, K] {
	return func(qt *QuadTree[T, K]) {
		qt.appender.maxDepth = depth
	}
}
