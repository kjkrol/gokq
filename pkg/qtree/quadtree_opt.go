package qtree

import (
	"github.com/kjkrol/gokg/pkg/geom"
)

type QuadTreeOption[T geom.Numeric] func(*QuadTree[T])

func WithMaxDepth[T geom.Numeric](depth int) QuadTreeOption[T] {
	return func(qt *QuadTree[T]) {
		qt.appender.maxDepth = depth
	}
}

func WithBatchCompressThreshold[T geom.Numeric](threshold int) QuadTreeOption[T] {
	return func(qt *QuadTree[T]) {
		if threshold > 0 {
			qt.coordinator.threshold = threshold
		}
	}
}
