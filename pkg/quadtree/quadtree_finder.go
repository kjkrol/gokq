package quadtree

import (
	"sort"

	"github.com/kjkrol/gokg/pkg/geometry"
)

type QuadTreeFinder[T geometry.SupportedNumeric] struct {
	plane    geometry.Plane[T]
	distance geometry.Distance[T]
}

func NewQuadTreeFinder[T geometry.SupportedNumeric](plane geometry.Plane[T]) QuadTreeFinder[T] {
	return QuadTreeFinder[T]{
		plane:    plane,
		distance: geometry.BoundingBoxDistanceForPlane(plane),
	}
}

func (qf QuadTreeFinder[T]) FindNeighbors(root *Node[T], target Item[T], margin T) []Item[T] {
	targetBounds := target.Bound()
	predicate := func(it Item[T]) bool {
		if it == target {
			return false
		}
		return qf.distance(targetBounds, it.Bound()) <= margin
	}

	probes := qf.probe(targetBounds, margin)
	if len(probes) == 1 {
		return singleProbeFind(root, probes[0], predicate)
	} else {
		return multiProbeFind(root, probes, predicate)
	}
}

func (qf QuadTreeFinder[T]) probe(aabb geometry.AABB[T], margin T) []geometry.AABB[T] {
	probe := aabb
	qf.plane.Expand(&probe, margin)
	rectangles := []geometry.AABB[T]{probe}
	if qf.plane.Name() == geometry.CYCLIC {
		for _, frag := range probe.Fragments() {
			rectangles = append(rectangles, frag)
		}
	}
	return rectangles
}

func singleProbeFind[T geometry.SupportedNumeric](
	root *Node[T],
	probe geometry.AABB[T],
	predicate func(it Item[T]) bool,
) []Item[T] {
	neighbors := make([]Item[T], 0)
	forEachIntersectingItem(root, probe, predicate, nil, func(item Item[T]) {
		neighbors = append(neighbors, item)
	})
	sortNeighbors(neighbors)
	return neighbors
}

func multiProbeFind[T geometry.SupportedNumeric](
	root *Node[T],
	probes []geometry.AABB[T],
	predicate func(it Item[T]) bool,
) []Item[T] {
	candidateSet := make(map[Item[T]]struct{})
	for _, probe := range probes {
		visited := make(map[*Node[T]]struct{})
		forEachIntersectingItem(root, probe, predicate, visited, func(item Item[T]) {
			candidateSet[item] = struct{}{}
		})
	}
	neighbors := make([]Item[T], 0, len(candidateSet))
	for item := range candidateSet {
		neighbors = append(neighbors, item)
	}
	sortNeighbors(neighbors)
	return neighbors
}

func forEachIntersectingItem[T geometry.SupportedNumeric](
	root *Node[T],
	probe geometry.AABB[T],
	predicate func(Item[T]) bool,
	visited map[*Node[T]]struct{},
	visit func(Item[T]),
) {
	if root == nil {
		return
	}

	stack := []*Node[T]{root}
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited != nil {
			if _, seen := visited[node]; seen {
				continue
			}
		}

		if !node.bounds.Intersects(probe) {
			continue
		}

		if visited != nil {
			visited[node] = struct{}{}
		}

		for _, item := range node.items {
			if predicate(item) {
				visit(item)
			}
		}

		stack = append(stack, node.childs...)
	}
}

func sortNeighbors[T geometry.SupportedNumeric](items []Item[T]) {
	sort.Slice(items, func(i, j int) bool {
		ai, aj := items[i].Bound(), items[j].Bound()
		first, _ := geometry.SortRectanglesBy(
			ai, aj,
			func(box geometry.AABB[T]) T { return box.TopLeft.Y },
			func(box geometry.AABB[T]) T { return box.TopLeft.X },
			func(box geometry.AABB[T]) T { return box.BottomRight.Y },
			func(box geometry.AABB[T]) T { return box.BottomRight.X },
		)
		return first.Equals(ai)
	})
}
