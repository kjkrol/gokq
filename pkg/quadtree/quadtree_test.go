package quadtree_test

import (
	"kjkrol/pkg/quadtree"
	"slices"
	"testing"

	"github.com/kjkrol/gokg/pkg/geometry"
)

type TestItem struct {
	pos geometry.Vec[int]
}

func (ts *TestItem) Vector() geometry.Vec[int] {
	return ts.pos
}

func (ts *TestItem) String() string {
	return ts.pos.String()
}

func TestQuadTreeAdd(t *testing.T) {
	topLeft := geometry.Vec[int]{X: 0, Y: 0}
	bottomRight := geometry.Vec[int]{X: 64, Y: 64}
	box := quadtree.NewBox(topLeft, bottomRight)
	qtree := quadtree.NewQuadTree(box)

	target := &TestItem{pos: geometry.Vec[int]{X: 32, Y: 32}}
	item1 := &TestItem{pos: geometry.Vec[int]{X: 32, Y: 31}}
	item2 := &TestItem{pos: geometry.Vec[int]{X: 32, Y: 33}}
	item3 := &TestItem{pos: geometry.Vec[int]{X: 31, Y: 32}}
	item4 := &TestItem{pos: geometry.Vec[int]{X: 33, Y: 32}}
	qtree.Add(item1)
	qtree.Add(item2)
	qtree.Add(item3)
	qtree.Add(item4)

	item5 := &TestItem{pos: geometry.Vec[int]{X: 12, Y: 11}}
	item6 := &TestItem{pos: geometry.Vec[int]{X: 12, Y: 13}}
	item7 := &TestItem{pos: geometry.Vec[int]{X: 11, Y: 12}}
	item8 := &TestItem{pos: geometry.Vec[int]{X: 13, Y: 12}}
	qtree.Add(item5)
	qtree.Add(item6)
	qtree.Add(item7)
	qtree.Add(item8)

	expected := [4]quadtree.Item{item1, item2, item3, item4}
	for range 100000 {
		foundNeightbors := qtree.FindNeighbors(target, 1, nil)
		if !allElementsPresent(foundNeightbors[:], expected[:]) {
			t.Errorf("result %v not equal to expected %v", foundNeightbors, expected)
		}
	}
}

func TestQuadTreeForBoundedPlane(t *testing.T) {
	rectBoardGemoetry := geometry.NewDiscreteBoundedPlane(4, 4)
	topLeft := geometry.Vec[int]{X: 0, Y: 0}
	bottomRight := geometry.Vec[int]{X: 3, Y: 3}
	box := quadtree.NewBox(topLeft, bottomRight)
	qtree := quadtree.NewQuadTree(box)

	target := &TestItem{pos: geometry.Vec[int]{X: 0, Y: 0}}
	item1 := &TestItem{pos: geometry.Vec[int]{X: 0, Y: 1}}
	item2 := &TestItem{pos: geometry.Vec[int]{X: 1, Y: 0}}
	item3 := &TestItem{pos: geometry.Vec[int]{X: 2, Y: 0}}
	item4 := &TestItem{pos: geometry.Vec[int]{X: 0, Y: 2}}

	qtree.Add(item1)
	qtree.Add(item2)
	qtree.Add(item3)
	qtree.Add(item4)

	item5 := &TestItem{pos: geometry.Vec[int]{X: 2, Y: 2}}
	item6 := &TestItem{pos: geometry.Vec[int]{X: 1, Y: 1}}
	qtree.Add(item5)
	qtree.Add(item6)

	expected := [2]quadtree.Item{item1, item2}
	for range 100000 {
		foundNeightbors := qtree.FindNeighbors(target, 1, rectBoardGemoetry.Metric)
		if !allElementsPresent(foundNeightbors[:], expected[:]) {
			t.Errorf("result %v not equal to expected %v", foundNeightbors, expected)
		}
	}
}

func TestQuadTreeForCyclicBoundedPlane(t *testing.T) {
	rectBoardGemoetry := geometry.NewDiscreteCyclicBoundedPlane(4, 4)
	topLeft := geometry.Vec[int]{X: 0, Y: 0}
	bottomRight := geometry.Vec[int]{X: 3, Y: 3}
	box := quadtree.NewBox(topLeft, bottomRight)
	qtree := quadtree.NewQuadTree(box)

	target := &TestItem{pos: geometry.Vec[int]{X: 0, Y: 0}}
	item1 := &TestItem{pos: geometry.Vec[int]{X: 0, Y: 1}}
	item2 := &TestItem{pos: geometry.Vec[int]{X: 1, Y: 0}}
	item3 := &TestItem{pos: geometry.Vec[int]{X: 2, Y: 0}}
	item4 := &TestItem{pos: geometry.Vec[int]{X: 0, Y: 2}}

	qtree.Add(item1)
	qtree.Add(item2)
	qtree.Add(item3)
	qtree.Add(item4)

	item5 := &TestItem{pos: geometry.Vec[int]{X: 2, Y: 2}}
	item6 := &TestItem{pos: geometry.Vec[int]{X: 1, Y: 1}}
	item7 := &TestItem{pos: geometry.Vec[int]{X: 0, Y: 3}}
	item8 := &TestItem{pos: geometry.Vec[int]{X: 3, Y: 0}}
	qtree.Add(item5)
	qtree.Add(item6)
	qtree.Add(item7)
	qtree.Add(item8)

	expected := [4]quadtree.Item{item1, item2, item7, item8}
	for range 100000 {
		foundNeightbors := qtree.FindNeighbors(target, 1, rectBoardGemoetry.Metric)
		if !allElementsPresent(foundNeightbors[:], expected[:]) {
			t.Errorf("result %v not equal to expected %v", foundNeightbors, expected)
		}
	}
}

func allElementsPresent[T comparable](slice1 []T, slice2 []T) bool {
	for _, val := range slice2 {
		if !slices.Contains(slice1, val) {
			return false
		}
	}
	return true
}
