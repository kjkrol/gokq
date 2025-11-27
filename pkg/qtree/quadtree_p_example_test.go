package qtree

import (
	"fmt"
	"sync/atomic"

	"github.com/kjkrol/gokg/pkg/geometry"
)

var globalID uint64

type TestItem[T geometry.SupportedNumeric] struct {
	geometry.BoundingBox[T]
	id uint64
}

func (t *TestItem[T]) Bound() geometry.BoundingBox[T] {
	return t.BoundingBox
}

func (t *TestItem[T]) SameID(other Item[T]) bool {
	o, ok := other.(*TestItem[T])
	if !ok {
		return false
	}
	return t.id == o.id
}

func newTestItemPointAtPos[T geometry.SupportedNumeric](x, y T) *TestItem[T] {
	vec := geometry.NewVec(x, y)
	box := geometry.NewBoundingBoxAround(vec, 0)
	id := atomic.AddUint64(&globalID, 1)
	return &TestItem[T]{BoundingBox: box, id: id}
}

func newTestItemFromPos[T geometry.SupportedNumeric](x1, y1, x2, y2 T) *TestItem[T] {
	box := geometry.NewBoundingBoxAt(geometry.NewVec(x1, y1), x2, y2)
	id := atomic.AddUint64(&globalID, 1)
	return &TestItem[T]{BoundingBox: box, id: id}
}

func newTestItemFromVec[T geometry.SupportedNumeric](vec geometry.Vec[T]) *TestItem[T] {
	box := geometry.NewBoundingBoxAround(vec, 0)
	id := atomic.AddUint64(&globalID, 1)
	return &TestItem[T]{BoundingBox: box, id: id}
}

func newTestItemFromBox[T geometry.SupportedNumeric](box geometry.BoundingBox[T]) *TestItem[T] {
	id := atomic.AddUint64(&globalID, 1)
	return &TestItem[T]{BoundingBox: box, id: id}
}

func ExampleQuadTree_FindNeighbors_targetInTree() {
	// Create a bounded plane and a quadtree
	boundedPlane := geometry.NewBoundedPlane(64, 64)
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
	boundedPlane := geometry.NewBoundedPlane(64, 64)
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
