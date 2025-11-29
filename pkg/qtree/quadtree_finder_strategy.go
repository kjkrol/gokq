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
	plane.Space2D[T]
}

func NewDefaultQuadTreeFinderStrategy[T geom.Numeric](
	plane plane.Space2D[T],
) QuadTreeFinderStrategy[T] {
	return DefaultQuadTreeFinderStrategy[T]{plane}
}

func (s DefaultQuadTreeFinderStrategy[T]) NodeIntersectionDetectionFactory(
	target Item[T],
	margin T,
) NodeIntersectionDetection[T] {
	probe := s.Space2D.WrapAABB(target.Bound())
	s.Space2D.Expand(&probe, margin)
	return func(node Node[T]) bool {
		intersection := node.bounds.Intersects(probe.AABB)
		probe.VisitFragments(func(pos plane.FragPosition, aabb geom.AABB[T]) bool {
			intersection = intersection || node.bounds.Intersects(aabb)
			return true
		})
		return intersection
	}
}

func (s DefaultQuadTreeFinderStrategy[T]) ItemsInRangeDetectionFactory(
	target Item[T],
	margin T,
) ItemsInRangeDetection[T] {
	boundingBoxDistance := s.Space2D.AABBDistance()
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
