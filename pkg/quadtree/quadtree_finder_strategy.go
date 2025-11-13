package quadtree

import "github.com/kjkrol/gokg/pkg/geometry"

type NodeIntersectionDetection[T geometry.SupportedNumeric, K comparable] func(Node[T, K]) bool
type ItemsInRangeDetection[T geometry.SupportedNumeric, K comparable] func(Node[T, K], func(Item[T, K]))

type QuadTreeFinderStrategy[T geometry.SupportedNumeric, K comparable] interface {
	NodeIntersectionDetectionFactory(target Item[T, K], margin T) NodeIntersectionDetection[T, K]
	ItemsInRangeDetectionFactory(target Item[T, K], maring T) ItemsInRangeDetection[T, K]
}

// ----------------- DefaultQuadTreeFinderStrategy -----------------

type DefaultQuadTreeFinderStrategy[T geometry.SupportedNumeric, K comparable] struct {
	geometry.Plane[T]
}

func NewDefaultQuadTreeFinderStrategy[T geometry.SupportedNumeric, K comparable](
	plane geometry.Plane[T],
) QuadTreeFinderStrategy[T, K] {
	return DefaultQuadTreeFinderStrategy[T, K]{plane}
}

func (s DefaultQuadTreeFinderStrategy[T, K]) NodeIntersectionDetectionFactory(
	target Item[T, K],
	margin T,
) NodeIntersectionDetection[T, K] {
	probe := s.Plane.WrapBoundingBox(target.Bound())
	s.Plane.Expand(&probe, margin)
	return func(node Node[T, K]) bool {
		intersection := node.bounds.Intersects(probe.BoundingBox)
		for _, frag := range probe.Fragments() {
			intersection = intersection || node.bounds.Intersects(frag)
		}
		return intersection
	}
}

func (s DefaultQuadTreeFinderStrategy[T, K]) ItemsInRangeDetectionFactory(
	target Item[T, K],
	margin T,
) ItemsInRangeDetection[T, K] {
	boundingBoxDistance := s.Plane.BoundingBoxDistance()
	return func(node Node[T, K], inRangeApply func(Item[T, K])) {
		for _, item := range node.items {
			if item.Id() == target.Id() {
				continue
			}
			if boundingBoxDistance(target.Bound(), item.Bound()) <= margin {
				inRangeApply(item)
			}
		}
	}
}
