package qtree

import (
	"github.com/kjkrol/gokg/pkg/geom"
	"github.com/kjkrol/gokg/pkg/plane"
)

type NodeIntersectionDetection[T geom.Numeric] func(Node[T]) bool
type ItemsInRangeDetection[T geom.Numeric] func(Node[T], func(Item[T]))

type QuadTreeFinderStrategy[T geom.Numeric] interface {
	NodeIntersectionDetectionFactory(target Item[T], margin T) NodeIntersectionDetection[T]
	ItemsInRangeDetectionFactory(target Item[T], maring T) ItemsInRangeDetection[T]
}

// ----------------- DefaultQuadTreeFinderStrategy -----------------

type DefaultQuadTreeFinderStrategy[T geom.Numeric] struct {
	plane.Space[T]
}

func NewDefaultQuadTreeFinderStrategy[T geom.Numeric](
	plane plane.Space[T],
) QuadTreeFinderStrategy[T] {
	return DefaultQuadTreeFinderStrategy[T]{plane}
}

func (s DefaultQuadTreeFinderStrategy[T]) NodeIntersectionDetectionFactory(
	target Item[T],
	margin T,
) NodeIntersectionDetection[T] {
	probe := s.Space.WrapAABB(target.Bound())
	s.Space.Expand(&probe, margin)
	return func(node Node[T]) bool {
		intersection := node.bounds.Intersects(probe.AABB)
		for _, frag := range probe.Fragments() {
			intersection = intersection || node.bounds.Intersects(frag)
		}
		return intersection
	}
}

func (s DefaultQuadTreeFinderStrategy[T]) ItemsInRangeDetectionFactory(
	target Item[T],
	margin T,
) ItemsInRangeDetection[T] {
	boundingBoxDistance := s.Space.AABBDistance()
	return func(node Node[T], inRangeApply func(Item[T])) {
		for _, item := range node.items {
			if item.SameID(target) {
				continue
			}
			if boundingBoxDistance(target.Bound(), item.Bound()) <= margin {
				inRangeApply(item)
			}
		}
	}
}
