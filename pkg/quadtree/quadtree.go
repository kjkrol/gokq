package quadtree

import (
	"sort"

	"github.com/kjkrol/gokg/pkg/geometry"
)

const CAPACITY int = 4

type Item[T geometry.SupportedNumeric] interface {
	Value() geometry.AABB[T]
}

type QuadTree[T geometry.SupportedNumeric] struct {
	root     *Node[T]
	plane    geometry.Plane[T]
	distance geometry.Distance[T]
	maxDepth int
}

func NewQuadTree[T geometry.SupportedNumeric](
	plane geometry.Plane[T],
	opts ...QuadTreeOption[T],
) *QuadTree[T] {
	rootBounds := geometry.NewAABB(geometry.Vec[T]{X: 0, Y: 0}, plane.Size())
	root := newNode(rootBounds, nil)
	qt := &QuadTree[T]{
		root:     root,
		plane:    plane,
		distance: geometry.BoundingBoxDistanceForPlane(plane),
		maxDepth: 10,
	}
	for _, opt := range opts {
		opt(qt)
	}
	return qt
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

func (t *QuadTree[T]) LeafRectangles() []geometry.AABB[T] {
	return t.root.leafRectangles()
}

func (t *QuadTree[T]) FindNeighbors(target Item[T], margin T) []Item[T] {
	probes := t.probe(target.Value(), margin)

	neighborSet := make(map[*Node[T]]struct{})
	for _, p := range probes {
		t.root.findIntersectingNodesUnique(p, neighborSet)
	}
	neighborNodes := make([]*Node[T], 0, len(neighborSet))
	for n := range neighborSet {
		neighborNodes = append(neighborNodes, n)
	}

	targetValue := target.Value()
	neighbors := scan(neighborNodes, func(it Item[T]) bool {
		if it == target {
			return false
		}
		return t.distance(targetValue, it.Value()) <= margin
	})

	sortNeighbors(neighbors)
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

// TODO: to wyglada na duplikat SortRectanglesBy z rectangle.go
func sortNeighbors[T geometry.SupportedNumeric](items []Item[T]) {
	sort.Slice(items, func(i, j int) bool {
		ai, aj := items[i].Value(), items[j].Value()
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
