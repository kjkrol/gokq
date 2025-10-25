package quadtree

import (
	"fmt"

	"github.com/kjkrol/gokg/pkg/geometry"
)

type ExampleItem[T geometry.SupportedNumeric] struct {
	spatial geometry.Spatial[T]
}

func (ei *ExampleItem[T]) Value() geometry.Spatial[T] {
	return ei.spatial
}

func ExampleQuadTree_FindNeighbors_targetInTree() {
	// Create a bounded plane and a quadtree
	boundedPlane := geometry.NewBoundedPlane(64, 64)
	qtree := NewQuadTree(boundedPlane)
	defer qtree.Close()

	// Add items to the quadtree
	items := []*ExampleItem[int]{
		{spatial: &geometry.Vec[int]{X: 32, Y: 32}},
		{spatial: &geometry.Vec[int]{X: 32, Y: 31}},
		{spatial: &geometry.Vec[int]{X: 32, Y: 33}},
		{spatial: &geometry.Vec[int]{X: 31, Y: 32}},
		{spatial: &geometry.Vec[int]{X: 33, Y: 32}},
	}
	for _, item := range items {
		qtree.Add(item)
	}

	// Find neighbors of a target item
	target := items[0]
	neighbors := qtree.FindNeighbors(target, 1)

	// Print the neighbors
	for _, neighbor := range neighbors {
		fmt.Println(neighbor.Value())
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
	items := []*ExampleItem[int]{
		{spatial: &geometry.Vec[int]{X: 32, Y: 32}},
		{spatial: &geometry.Vec[int]{X: 32, Y: 31}},
		{spatial: &geometry.Vec[int]{X: 32, Y: 33}},
		{spatial: &geometry.Vec[int]{X: 31, Y: 32}},
		{spatial: &geometry.Vec[int]{X: 33, Y: 32}},
	}
	for _, item := range items {
		qtree.Add(item)
	}

	// Find neighbors of a target item
	target := &ExampleItem[int]{spatial: &geometry.Vec[int]{X: 32, Y: 32}}
	neighbors := qtree.FindNeighbors(target, 1)

	// Print the neighbors
	for _, neighbor := range neighbors {
		fmt.Println(neighbor.Value())
	}

	// Output:
	// (32,31)
	// (31,32)
	// (32,32)
	// (33,32)
	// (32,33)
}
