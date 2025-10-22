package quadtree

import (
	"sort"

	"github.com/kjkrol/gokg/pkg/geometry"
)

const CAPACITY int = 4

type Item[T geometry.SupportedNumeric] interface {
	Value() geometry.Spatial[T]
}

type QuadTree[T geometry.SupportedNumeric] struct {
	root  *Node[T]
	plane geometry.Plane[T]
}

func NewQuadTree[T geometry.SupportedNumeric](plane geometry.Plane[T]) *QuadTree[T] {
	rootBounds := geometry.NewRectangle(geometry.Vec[T]{X: 0, Y: 0}, plane.Size())
	root := newNode[T](rootBounds, nil)
	return &QuadTree[T]{root: root, plane: plane}
}

func (t *QuadTree[T]) Add(item Item[T]) {
	t.root.add(item)
}

func (t *QuadTree[T]) Remove(item Item[T]) bool {
	return t.root.remove(item)
}

func (t *QuadTree[T]) Close() {
	t.root.close()
}

func (t *QuadTree[T]) Count() int {
	return len(t.root.allItems())
}

func (t *QuadTree[T]) Depth() int {
	return t.root.depth()
}

func (t *QuadTree[T]) AllItems() []Item[T] {
	return t.root.allItems()
}

func (t *QuadTree[T]) LeafRectangles() []geometry.Rectangle[T] {
	return t.root.leafRectangles()
}

func (t *QuadTree[T]) FindNeighbors(target Item[T], margin T) []Item[T] {
	probes := target.Value().Probe(margin, t.plane)

	neighborSet := make(map[*Node[T]]struct{})
	for _, p := range probes {
		t.root.findIntersectingNodesUnique(p, neighborSet)
	}
	neighborNodes := make([]*Node[T], 0, len(neighborSet))
	for n := range neighborSet {
		neighborNodes = append(neighborNodes, n)
	}

	targetValue := target.Value()
	neighbors := scan[T](neighborNodes, func(it Item[T]) bool {
		if it == target {
			return false
		}
		return targetValue.DistanceTo(it.Value(), t.plane.Metric) <= margin
	})

	sortNeighbors[T](neighbors)
	return neighbors
}

// ---------------------- helpers ----------------------

func scan[T geometry.SupportedNumeric](
	neighborNodes []*Node[T],
	predicate func(Item[T]) bool,
) []Item[T] {
	neighborItems := make([]Item[T], 0)
	for _, node := range neighborNodes {
		for _, neighborItem := range node.items {
			if predicate(neighborItem) {
				neighborItems = append(neighborItems, neighborItem)
			}
		}
	}
	return neighborItems
}

func sortNeighbors[T geometry.SupportedNumeric](items []Item[T]) {
	sort.Slice(items, func(i, j int) bool {
		ai, aj := items[i].Value().Bounds(), items[j].Value().Bounds()
		if ai.TopLeft.Y != aj.TopLeft.Y {
			return ai.TopLeft.Y < aj.TopLeft.Y
		}
		if ai.TopLeft.X != aj.TopLeft.X {
			return ai.TopLeft.X < aj.TopLeft.X
		}
		if ai.BottomRight.Y != aj.BottomRight.Y {
			return ai.BottomRight.Y < aj.BottomRight.Y
		}
		return ai.BottomRight.X < aj.BottomRight.X
	})
}
