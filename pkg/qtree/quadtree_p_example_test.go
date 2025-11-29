package qtree

import (
	"fmt"
	"sync/atomic"

	"github.com/kjkrol/gokg/pkg/geom"
	"github.com/kjkrol/gokg/pkg/plane"
)

var globalID uint64

type TestItem[T geom.Numeric] struct {
	geom.AABB[T]
	id uint64
}

func (t *TestItem[T]) Bound() geom.AABB[T] {
	return t.AABB
}

func (t *TestItem[T]) SameID(other Item[T]) bool {
	o, ok := other.(*TestItem[T])
	if !ok {
		return false
	}
	return t.id == o.id
}

func newTestItemPointAtPos[T geom.Numeric](x, y T) *TestItem[T] {
	vec := geom.NewVec(x, y)
	box := geom.NewAABBAround(vec, 0)
	id := atomic.AddUint64(&globalID, 1)
	return &TestItem[T]{AABB: box, id: id}
}

func newTestItemFromPos[T geom.Numeric](x1, y1, x2, y2 T) *TestItem[T] {
	box := geom.NewAABBAt(geom.NewVec(x1, y1), x2, y2)
	id := atomic.AddUint64(&globalID, 1)
	return &TestItem[T]{AABB: box, id: id}
}

func newTestItemFromVec[T geom.Numeric](vec geom.Vec[T]) *TestItem[T] {
	box := geom.NewAABBAround(vec, 0)
	id := atomic.AddUint64(&globalID, 1)
	return &TestItem[T]{AABB: box, id: id}
}

func newTestItemFromBox[T geom.Numeric](box geom.AABB[T]) *TestItem[T] {
	id := atomic.AddUint64(&globalID, 1)
	return &TestItem[T]{AABB: box, id: id}
}

func ExampleQuadTree_FindNeighbors_targetInTree() {
	// Create a bounded plane and a quadtree
	boundedPlane := plane.NewEuclidean2D(64, 64)
	qtree := NewQuadTree(boundedPlane)
	defer qtree.Close()

	// Add items to the quadtree
	items := []*TestItem[int]{
		newTestItemPointAtPos(32, 32),
		newTestItemPointAtPos(32, 31),
		newTestItemPointAtPos(32, 33),
		newTestItemPointAtPos(31, 32),
		newTestItemPointAtPos(33, 32),
	}
	for _, item := range items {
		qtree.Add(item)
	}

	// Find neighbors of a target item
	target := items[0]
	neighbors := qtree.FindNeighbors(target, 1)

	// Print the neighbors
	for _, neighbor := range neighbors {
		fmt.Println(neighbor)
	}

	// Output:
	// {(32,31) (32,31)}
	// {(31,32) (31,32)}
	// {(33,32) (33,32)}
	// {(32,33) (32,33)}
}

func ExampleQuadTree_FindNeighbors_targetNotInTree() {
	// Create a bounded plane and a quadtree
	boundedPlane := plane.NewEuclidean2D(64, 64)
	qtree := NewQuadTree(boundedPlane)
	defer qtree.Close()

	// Add items to the quadtree
	items := []*TestItem[int]{
		newTestItemPointAtPos(32, 32),
		newTestItemPointAtPos(32, 31),
		newTestItemPointAtPos(32, 33),
		newTestItemPointAtPos(31, 32),
		newTestItemPointAtPos(33, 32),
	}
	for _, item := range items {
		qtree.Add(item)
	}

	// Find neighbors of a target item
	target := newTestItemPointAtPos(32, 32)
	neighbors := qtree.FindNeighbors(target, 1)

	// Print the neighbors
	for _, neighbor := range neighbors {
		fmt.Println(neighbor)
	}

	// Output:
	// {(32,31) (32,31)}
	// {(31,32) (31,32)}
	// {(32,32) (32,32)}
	// {(33,32) (33,32)}
	// {(32,33) (32,33)}
}
