package qtree

import "github.com/kjkrol/gokg/pkg/geometry"

type NodeIntersectionDetection[T geometry.SupportedNumeric] func(Node[T]) bool
type ItemsInRangeDetection[T geometry.SupportedNumeric] func(Node[T], func(Item[T]))

type QuadTreeFinderStrategy[T geometry.SupportedNumeric] interface {
	NodeIntersectionDetectionFactory(target Item[T], margin T) NodeIntersectionDetection[T]
	ItemsInRangeDetectionFactory(target Item[T], maring T) ItemsInRangeDetection[T]
}

// ----------------- DefaultQuadTreeFinderStrategy -----------------

type DefaultQuadTreeFinderStrategy[T geometry.SupportedNumeric] struct {
	geometry.Plane[T]
}

func NewDefaultQuadTreeFinderStrategy[T geometry.SupportedNumeric](
	plane geometry.Plane[T],
) QuadTreeFinderStrategy[T] {
	return DefaultQuadTreeFinderStrategy[T]{plane}
}

func (s DefaultQuadTreeFinderStrategy[T]) NodeIntersectionDetectionFactory(
	target Item[T],
	margin T,
) NodeIntersectionDetection[T] {
	probe := s.Plane.WrapBoundingBox(target.Bound())
	s.Plane.Expand(&probe, margin)
	return func(node Node[T]) bool {
		intersection := node.bounds.Intersects(probe.BoundingBox)
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
	boundingBoxDistance := s.Plane.BoundingBoxDistance()
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
