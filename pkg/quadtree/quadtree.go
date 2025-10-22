package quadtree

import (
	"sort"

	"github.com/kjkrol/gokg/pkg/geometry"
)

// ---------------------- Public API ----------------------

const CAPACITY int = 4

type SpatialValue[T geometry.SupportedNumeric] interface {
	geometry.Vec[T] | Box[T]
}

type Item[V SpatialValue[T], T geometry.SupportedNumeric] interface {
	Value() V
}

type QuadTree[V SpatialValue[T], T geometry.SupportedNumeric] struct {
	root  *Node[V, T]
	plane geometry.Plane[T]
	ops   NeighborOps[V, T]
}

func NewQuadTree[V SpatialValue[T], T geometry.SupportedNumeric](plane geometry.Plane[T], ops NeighborOps[V, T]) *QuadTree[V, T] {
	box := NewBox(geometry.Vec[T]{X: 0, Y: 0}, plane.Size())
	root := newNode[V](box, nil)
	return &QuadTree[V, T]{root: root, plane: plane, ops: ops}
}

func NewBoxQuadTree[T geometry.SupportedNumeric](plane geometry.Plane[T]) *QuadTree[Box[T], T] {
	return NewQuadTree(plane, BoxOps[T]{})
}

func NewPointQuadTree[T geometry.SupportedNumeric](plane geometry.Plane[T]) *QuadTree[geometry.Vec[T], T] {
	return NewQuadTree(plane, PointOps[T]{})
}

func (t *QuadTree[V, T]) Add(item Item[V, T]) {
	t.root.add(item, t.ops)
}

func (t *QuadTree[V, T]) Remove(item Item[V, T]) bool {
	return t.root.remove(item, t.ops)
}

func (t *QuadTree[V, T]) Close() {
	t.root.close()
}

func (t *QuadTree[V, T]) Count() int {
	return len(t.root.allItems())
}

func (t *QuadTree[V, T]) Depth() int {
	return t.root.depth()
}

func (t *QuadTree[V, T]) AllItems() []Item[V, T] {
	return t.root.allItems()
}

func (t *QuadTree[V, T]) LeafBoxes() []Box[T] {
	return t.root.leafBoxes()
}

func (t *QuadTree[V, T]) FindNeighbors(target Item[V, T], margin T) []Item[V, T] {
	probes := t.ops.Probe(target.Value(), margin, t.plane)

	neighborSet := make(map[*Node[V, T]]struct{})
	for _, p := range probes {
		t.root.findIntersectingNodesUnique(p, neighborSet)
	}
	neighborNodes := make([]*Node[V, T], 0, len(neighborSet))
	for n := range neighborSet {
		neighborNodes = append(neighborNodes, n)
	}

	neighbors := scan(neighborNodes, func(it *Item[V, T]) bool {
		return *it != target &&
			t.ops.Distance(target.Value(), (*it).Value(), t.plane.Metric) <= margin
	})

	sortNeighbors(neighbors, t.ops)
	return neighbors
}

// ---------------------- helpers ----------------------

func scan[V SpatialValue[T], T geometry.SupportedNumeric](neighborNodes []*Node[V, T], predicate func(*Item[V, T]) bool) []Item[V, T] {
	neighborItems := make([]Item[V, T], 0)
	for _, node := range neighborNodes {
		for _, neighborItem := range node.items {
			if predicate(&neighborItem) {
				neighborItems = append(neighborItems, neighborItem)
			}
		}
	}
	return neighborItems
}

func sortNeighbors[V SpatialValue[T], T geometry.SupportedNumeric](items []Item[V, T], ops NeighborOps[V, T]) {
	sort.Slice(items, func(i, j int) bool {
		ai, aj := ops.BoundsOf(items[i].Value()), ops.BoundsOf(items[j].Value())
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

//
// ---------------------- Strategy (NeighborOps) ----------------------
//

// NeighborOps definiuje operacje przestrzenne zależne od rodzaju elementu.
type NeighborOps[V any, T geometry.SupportedNumeric] interface {
	BoundsOf(v V) Box[T]
	Probe(v V, margin T, plane geometry.Plane[T]) []Box[T]
	Distance(a, b V, metric func(geometry.Vec[T], geometry.Vec[T]) T) T
}

// ------------------------------------------------------------
// PointOps – implementacja NeighborOps dla punktów.
// ------------------------------------------------------------
type PointOps[T geometry.SupportedNumeric] struct{}

func (PointOps[T]) BoundsOf(v geometry.Vec[T]) Box[T] {
	return BuildBox(v, 0)
}
func (PointOps[T]) Probe(v geometry.Vec[T], margin T, plane geometry.Plane[T]) []Box[T] {
	probe := BuildBox(v, margin)
	boxes := []Box[T]{probe}
	if plane.Name() == "cyclic" {
		boxes = append(boxes, WrapBoxCyclic(probe, plane.Size(), plane.Contains)...)
	}
	return boxes
}
func (PointOps[T]) Distance(a, b geometry.Vec[T], metric func(geometry.Vec[T], geometry.Vec[T]) T) T {
	return metric(a, b)
}

// ------------------------------------------------------------
// BoxOps – implementacja NeighborOps dla boxów.
// ------------------------------------------------------------
type BoxOps[T geometry.SupportedNumeric] struct{}

func (BoxOps[T]) BoundsOf(b Box[T]) Box[T] {
	return b
}
func (BoxOps[T]) Probe(b Box[T], margin T, plane geometry.Plane[T]) []Box[T] {
	probe := b.Expand(margin)
	boxes := []Box[T]{probe}
	if plane.Name() == "cyclic" {
		boxes = append(boxes, WrapBoxCyclic(probe, plane.Size(), plane.Contains)...)
	}
	return boxes
}
func (BoxOps[T]) Distance(a, b Box[T], metric func(geometry.Vec[T], geometry.Vec[T]) T) T {
	if a.Intersects(b) {
		return 0
	}
	dx := axisDistance(a, b, func(v geometry.Vec[T]) T { return v.X })
	dy := axisDistance(a, b, func(v geometry.Vec[T]) T { return v.Y })
	return metric(geometry.Vec[T]{X: dx, Y: dy}, geometry.Vec[T]{X: 0, Y: 0})
}

func axisDistance[T geometry.SupportedNumeric](
	aa, bb Box[T],
	axisValue func(geometry.Vec[T]) T,
) T {
	aa, bb = SortBy(aa, bb, axisValue)

	if axisValue(aa.BottomRight) >= axisValue(bb.TopLeft) {
		return 0
	}
	return axisValue(bb.TopLeft) - axisValue(aa.BottomRight)
}

// ------------------------------------------------------------
