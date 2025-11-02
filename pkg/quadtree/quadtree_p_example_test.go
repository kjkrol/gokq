package quadtree

import (
	"fmt"

	"github.com/kjkrol/gokg/pkg/geometry"
)

type ShapeItem[T geometry.SupportedNumeric] struct {
	shape geometry.Shape[T]
}

func (si *ShapeItem[T]) AABB() geometry.AABB[T] {
	return si.shape.Bounds()
}

func (si *ShapeItem[T]) String() string {
	return si.shape.String()
}

func ExampleQuadTree_FindNeighbors_targetInTree() {
	// Create a bounded plane and a quadtree
	boundedPlane := geometry.NewBoundedPlane(64, 64)
	qtree := NewQuadTree(boundedPlane)
	defer qtree.Close()

	// Add items to the quadtree
	items := []*ShapeItem[int]{
		{shape: &geometry.Vec[int]{X: 32, Y: 32}},
		{shape: &geometry.Vec[int]{X: 32, Y: 31}},
		{shape: &geometry.Vec[int]{X: 32, Y: 33}},
		{shape: &geometry.Vec[int]{X: 31, Y: 32}},
		{shape: &geometry.Vec[int]{X: 33, Y: 32}},
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
	// (32,31)
	// (31,32)
	// (33,32)
	// (32,33)
}

func ExampleQuadTree_FindNeighbors_targetNotInTree() {
	// Create a bounded plane and a quadtree
	boundedPlane := geometry.NewBoundedPlane(64, 64)
	qtree := NewQuadTree(boundedPlane)
	defer qtree.Close()

	// Add items to the quadtree
	items := []*ShapeItem[int]{
		{shape: &geometry.Vec[int]{X: 32, Y: 32}},
		{shape: &geometry.Vec[int]{X: 32, Y: 31}},
		{shape: &geometry.Vec[int]{X: 32, Y: 33}},
		{shape: &geometry.Vec[int]{X: 31, Y: 32}},
		{shape: &geometry.Vec[int]{X: 33, Y: 32}},
	}
	for _, item := range items {
		qtree.Add(item)
	}

	// Find neighbors of a target item
	target := &ShapeItem[int]{shape: &geometry.Vec[int]{X: 32, Y: 32}}
	neighbors := qtree.FindNeighbors(target, 1)

	// Print the neighbors
	for _, neighbor := range neighbors {
		fmt.Println(neighbor)
	}

	// Output:
	// (32,31)
	// (31,32)
	// (32,32)
	// (33,32)
	// (32,33)
}
