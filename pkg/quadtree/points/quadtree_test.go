package points

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geometry"
	"github.com/kjkrol/goku/pkg/sliceutils"
)

type TestItem[T geometry.SupportedNumeric] struct {
	pos geometry.Vec[T]
}

func (ts *TestItem[T]) Vector() geometry.Vec[T] {
	return ts.pos
}

func (ts *TestItem[T]) String() string {
	return ts.pos.String()
}

func newTestItem[T geometry.SupportedNumeric](x, y T) *TestItem[T] {
	return &TestItem[T]{pos: geometry.Vec[T]{X: x, Y: y}}
}

func TestQuadTreeFindNeighbors(t *testing.T) {
	boundedPlane := geometry.NewBoundedPlane(64.0, 64.0)
	qtree := NewQuadTree(boundedPlane)

	defer qtree.Close()

	target := newTestItem[float64](32, 32)
	item1 := newTestItem[float64](32, 31)
	item2 := newTestItem[float64](32, 33)
	item3 := newTestItem[float64](31, 32)
	item4 := newTestItem[float64](33, 32)
	qtree.Add(item1)
	qtree.Add(item2)
	qtree.Add(item3)
	qtree.Add(item4)

	item5 := newTestItem[float64](12, 11)
	item6 := newTestItem[float64](12, 13)
	item7 := newTestItem[float64](11, 12)
	item8 := newTestItem[float64](13, 12)
	qtree.Add(item5)
	qtree.Add(item6)
	qtree.Add(item7)
	qtree.Add(item8)

	expected := [4]Item[float64]{item1, item2, item3, item4}
	for i := 0; i < 100; i++ {
		foundNeightbors := qtree.FindNeighbors(target, 1)
		if !sliceutils.SameElements(foundNeightbors[:], expected[:]) {
			t.Errorf("result %v not equal to expected %v", foundNeightbors, expected)
		}
	}
}

func TestQuadTreeForBoundedPlane(t *testing.T) {
	boundedPlane := geometry.NewBoundedPlane(4, 4)
	qtree := NewQuadTree(boundedPlane)

	defer qtree.Close()

	target := newTestItem[int](0, 0)
	item1 := newTestItem[int](0, 1)
	item2 := newTestItem[int](1, 0)
	item3 := newTestItem[int](2, 0)
	item4 := newTestItem[int](0, 2)

	qtree.Add(item1)
	qtree.Add(item2)
	qtree.Add(item3)
	qtree.Add(item4)

	item5 := newTestItem[int](2, 2)
	item6 := newTestItem[int](1, 1)
	qtree.Add(item5)
	qtree.Add(item6)

	expected := [2]Item[int]{item1, item2}
	for i := 0; i < 100000; i++ {
		foundNeightbors := qtree.FindNeighbors(target, 1)
		if !sliceutils.SameElements(foundNeightbors[:], expected[:]) {
			t.Errorf("result %v not equal to expected %v", foundNeightbors, expected)
		}
	}
}

func TestQuadTreeForCyclicBoundedPlane(t *testing.T) {
	cyclicBoundedPlane := geometry.NewCyclicBoundedPlane(4, 4)
	qtree := NewQuadTree(cyclicBoundedPlane)

	defer qtree.Close()

	target := newTestItem[int](0, 0)
	item1 := newTestItem[int](0, 1)
	item2 := newTestItem[int](1, 0)
	item3 := newTestItem[int](2, 0)
	item4 := newTestItem[int](0, 2)

	qtree.Add(item1)
	qtree.Add(item2)
	qtree.Add(item3)
	qtree.Add(item4)

	item5 := newTestItem[int](2, 2)
	item6 := newTestItem[int](1, 1)
	item7 := newTestItem[int](0, 3)
	item8 := newTestItem[int](3, 0)
	qtree.Add(item5)
	qtree.Add(item6)
	qtree.Add(item7)
	qtree.Add(item8)

	expected := [4]Item[int]{item1, item2, item7, item8}
	for i := 0; i < 100; i++ {
		foundNeighbors := qtree.FindNeighbors(target, 1)
		if !sliceutils.SameElements(foundNeighbors[:], expected[:]) {
			t.Errorf("result %v not equal to expected %v", foundNeighbors, expected)
		}
	}
}

func TestQuadTreeForCyclicBoundedPlaneWithLeavesIn2ndGeneration(t *testing.T) {
	cyclicBoundedPlane := geometry.NewCyclicBoundedPlane(100, 100)
	qtree := NewQuadTree(cyclicBoundedPlane)

	defer qtree.Close()

	item1 := &TestItem[int]{pos: geometry.Vec[int]{X: 0, Y: 0}}
	item2 := &TestItem[int]{pos: geometry.Vec[int]{X: 0, Y: 99}}
	item3 := &TestItem[int]{pos: geometry.Vec[int]{X: 99, Y: 0}}

	qtree.Add(item1)
	qtree.Add(item2)
	qtree.Add(item3)

	// Insert points such that the first leaves appear only in the 2th generation/layer of nodes
	for _, point := range []geometry.Vec[int]{
		{X: 10, Y: 10}, {X: 90, Y: 10}, {X: 10, Y: 90}, {X: 90, Y: 90},
		{X: 30, Y: 30}, {X: 70, Y: 30}, {X: 30, Y: 70}, {X: 70, Y: 70},
		{X: 20, Y: 20}, {X: 80, Y: 20}, {X: 20, Y: 80}, {X: 80, Y: 80},
		{X: 40, Y: 40}, {X: 60, Y: 40}, {X: 40, Y: 60}, {X: 60, Y: 60},
	} {
		item := &TestItem[int]{pos: point}
		qtree.Add(item)
	}

	// Verify that the first leaves appear only in the 2nd generation/layer
	verifyNodeDepth(t, qtree.root, 0, 2)

	expected := [2]Item[int]{item2, item3}
	for i := 0; i < 100; i++ {
		foundNeighbors := qtree.FindNeighbors(item1, 1)
		if !sliceutils.SameElements(foundNeighbors[:], expected[:]) {
			t.Errorf("result %v not equal to expected %v", foundNeighbors, expected)
		}
	}

}

func TestQuadTreeForCyclicBoundedPlaneWithLeavesIn5thGeneration(t *testing.T) {
	cyclicBoundedPlane := geometry.NewCyclicBoundedPlane(100.0, 100.0)
	qtree := NewQuadTree(cyclicBoundedPlane)

	defer qtree.Close()

	item1 := &TestItem[float64]{pos: geometry.Vec[float64]{X: 0, Y: 0}}
	item2 := &TestItem[float64]{pos: geometry.Vec[float64]{X: 0, Y: 99}}
	item3 := &TestItem[float64]{pos: geometry.Vec[float64]{X: 99, Y: 0}}

	qtree.Add(item1)
	qtree.Add(item2)
	qtree.Add(item3)

	// Insert points such that the first leaves appear only in the 5th generation/layer of nodes
	for _, point := range []geometry.Vec[float64]{
		{X: 1, Y: 1}, {X: 2, Y: 2}, {X: 3, Y: 3}, {X: 4, Y: 4},
		{X: 5, Y: 5}, {X: 6, Y: 6}, {X: 7, Y: 7}, {X: 8, Y: 8},
		{X: 9, Y: 9}, {X: 10, Y: 10}, {X: 11, Y: 11}, {X: 12, Y: 12},
		{X: 13, Y: 13}, {X: 14, Y: 14}, {X: 15, Y: 15}, {X: 16, Y: 16},
		{X: 17, Y: 17}, {X: 18, Y: 18}, {X: 19, Y: 19}, {X: 20, Y: 20},
		{X: 21, Y: 21}, {X: 22, Y: 22}, {X: 23, Y: 23}, {X: 24, Y: 24},
		{X: 25, Y: 25}, {X: 26, Y: 26}, {X: 27, Y: 27}, {X: 28, Y: 28},
		{X: 29, Y: 29}, {X: 30, Y: 30}, {X: 31, Y: 31}, {X: 32, Y: 32},
		{X: 33, Y: 33}, {X: 34, Y: 34}, {X: 35, Y: 35}, {X: 36, Y: 36},
		{X: 37, Y: 37}, {X: 38, Y: 38}, {X: 39, Y: 39}, {X: 40, Y: 40},
		{X: 41, Y: 41}, {X: 42, Y: 42}, {X: 43, Y: 43}, {X: 44, Y: 44},
		{X: 45, Y: 45}, {X: 46, Y: 46}, {X: 47, Y: 47}, {X: 48, Y: 48},
		{X: 49, Y: 49}, {X: 50, Y: 50}, {X: 51, Y: 51}, {X: 52, Y: 52},
		{X: 53, Y: 53}, {X: 54, Y: 54}, {X: 55, Y: 55}, {X: 56, Y: 56},
		{X: 57, Y: 57}, {X: 58, Y: 58}, {X: 59, Y: 59}, {X: 60, Y: 60},
		{X: 61, Y: 61}, {X: 62, Y: 62}, {X: 63, Y: 63}, {X: 64, Y: 64},
		{X: 65, Y: 65}, {X: 66, Y: 66}, {X: 67, Y: 67}, {X: 68, Y: 68},
		{X: 69, Y: 69}, {X: 70, Y: 70}, {X: 71, Y: 71}, {X: 72, Y: 72},
		{X: 73, Y: 73}, {X: 74, Y: 74}, {X: 75, Y: 75}, {X: 76, Y: 76},
		{X: 77, Y: 77}, {X: 78, Y: 78}, {X: 79, Y: 79}, {X: 80, Y: 80},
		{X: 81, Y: 81}, {X: 82, Y: 82}, {X: 83, Y: 83}, {X: 84, Y: 84},
		{X: 85, Y: 85}, {X: 86, Y: 86}, {X: 87, Y: 87}, {X: 88, Y: 88},
		{X: 89, Y: 89}, {X: 90, Y: 90}, {X: 91, Y: 91}, {X: 92, Y: 92},
		{X: 93, Y: 93}, {X: 94, Y: 94}, {X: 95, Y: 95}, {X: 96, Y: 96},
		{X: 97, Y: 97}, {X: 98, Y: 98}, {X: 99, Y: 99},
	} {
		item := &TestItem[float64]{pos: point}
		qtree.Add(item)
	}

	expected := [2]Item[float64]{item2, item3}
	for i := 0; i < 100; i++ {
		foundNeighbors := qtree.FindNeighbors(item1, 1)
		if !sliceutils.SameElements(foundNeighbors[:], expected[:]) {
			t.Errorf("result %v not equal to expected %v", foundNeighbors, expected)
		}
	}
}

func verifyNodeDepth(t *testing.T, node *Node[int], currentDepth, targetDepth int) {
	if node == nil {
		return
	}

	if currentDepth == targetDepth {
		if !node.isLeaf() {
			t.Errorf("Node at depth %d is not a leaf", currentDepth)
		}
		return
	}

	if node.isLeaf() {
		t.Errorf("Node at depth %d is a leaf, but it should not be", currentDepth)
		t.Logf("Node at depth %d has items: %v", currentDepth, node.items)
		return
	}

	for _, child := range node.childs {
		verifyNodeDepth(t, child, currentDepth+1, targetDepth)
	}
}
