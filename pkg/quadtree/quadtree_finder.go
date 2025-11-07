package quadtree

import (
	"sort"

	"github.com/kjkrol/gokg/pkg/geometry"
)

type QuadTreeFinder[T geometry.SupportedNumeric, K comparable] struct {
	plane    geometry.Plane[T]
	distance geometry.Distance[T]
}

func NewQuadTreeFinder[T geometry.SupportedNumeric, K comparable](plane geometry.Plane[T]) QuadTreeFinder[T, K] {
	return QuadTreeFinder[T, K]{
		plane:    plane,
		distance: geometry.BoundingBoxDistance(plane),
	}
}

func (qf QuadTreeFinder[T, K]) FindNeighbors(root *Node[T, K], target Item[T, K], margin T) []Item[T, K] {
	targetBounds := target.Bound()
	predicate := func(it Item[T, K]) bool {
		if it.Id() == target.Id() {
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

func (qf QuadTreeFinder[T, K]) probe(box geometry.BoundingBox[T], margin T) []geometry.BoundingBox[T] {
	probe := qf.plane.WrapBoundingBox(box)
	qf.plane.Expand(&probe, margin)
	rectangles := []geometry.BoundingBox[T]{probe.BoundingBox}
	for _, frag := range probe.Fragments() {
		rectangles = append(rectangles, frag)
	}
	return rectangles
}

func singleProbeFind[T geometry.SupportedNumeric, K comparable](
	root *Node[T, K],
	probe geometry.BoundingBox[T],
	predicate func(it Item[T, K]) bool,
) []Item[T, K] {
	neighbors := make([]Item[T, K], 0)
	forEachIntersectingItem(root, probe, predicate, nil, func(item Item[T, K]) {
		neighbors = append(neighbors, item)
	})
	sortNeighbors(neighbors)
	return neighbors
}

func multiProbeFind[T geometry.SupportedNumeric, K comparable](
	root *Node[T, K],
	probes []geometry.BoundingBox[T],
	predicate func(it Item[T, K]) bool,
) []Item[T, K] {
	candidateSet := make(map[Item[T, K]]struct{})
	for _, probe := range probes {
		visited := make(map[*Node[T, K]]struct{})
		forEachIntersectingItem(root, probe, predicate, visited, func(item Item[T, K]) {
			candidateSet[item] = struct{}{}
		})
	}
	neighbors := make([]Item[T, K], 0, len(candidateSet))
	for item := range candidateSet {
		neighbors = append(neighbors, item)
	}
	sortNeighbors(neighbors)
	return neighbors
}

func forEachIntersectingItem[T geometry.SupportedNumeric, K comparable](
	root *Node[T, K],
	probe geometry.BoundingBox[T],
	predicate func(Item[T, K]) bool,
	visited map[*Node[T, K]]struct{},
	visit func(Item[T, K]),
) {
	if root == nil {
		return
	}

	stack := []*Node[T, K]{root}
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

func sortNeighbors[T geometry.SupportedNumeric, K comparable](items []Item[T, K]) {
	sort.Slice(items, func(i, j int) bool {
		ai, aj := items[i].Bound(), items[j].Bound()
		first, _ := geometry.SortBoxesBy(
			ai, aj,
			func(box geometry.BoundingBox[T]) T { return box.TopLeft.Y },
			func(box geometry.BoundingBox[T]) T { return box.TopLeft.X },
			func(box geometry.BoundingBox[T]) T { return box.BottomRight.Y },
			func(box geometry.BoundingBox[T]) T { return box.BottomRight.X },
		)
		return first.Equals(ai)
	})
}
