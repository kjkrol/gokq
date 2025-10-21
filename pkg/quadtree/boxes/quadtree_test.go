package boxes

import (
	"testing"

	quadcore "github.com/kjkrol/goka/pkg/quadtree/base"
	"github.com/kjkrol/gokg/pkg/geometry"
	"github.com/kjkrol/goku/pkg/sliceutils"
)

type TestBoxItem[T geometry.SupportedNumeric] struct {
	box quadcore.Box[T]
}

func (ti *TestBoxItem[T]) Bounds() quadcore.Box[T] {
	return ti.box
}

func newTestBoxItem[T geometry.SupportedNumeric](x1, y1, x2, y2 T) *TestBoxItem[T] {
	return &TestBoxItem[T]{box: quadcore.NewBox(
		geometry.Vec[T]{X: x1, Y: y1},
		geometry.Vec[T]{X: x2, Y: y2},
	)}
}

func TestQuadTreeBox_FindNeighborsSimple(t *testing.T) {
	boundedPlane := geometry.NewBoundedPlane(64.0, 64.0)
	qtree := NewQuadTree(boundedPlane)
	defer qtree.Close()

	target := newTestBoxItem(10.0, 10.0, 12.0, 12.0)
	item1 := newTestBoxItem(10.0, 8.0, 12.0, 10.0)  // nad targetem
	item2 := newTestBoxItem(8.0, 10.0, 10.0, 12.0)  // z lewej
	item3 := newTestBoxItem(12.0, 10.0, 14.0, 12.0) // z prawej
	item4 := newTestBoxItem(10.0, 12.0, 12.0, 14.0) // pod spodem
	itemFar := newTestBoxItem(30.0, 30.0, 32.0, 32.0)

	qtree.Add(item1)
	qtree.Add(item2)
	qtree.Add(item3)
	qtree.Add(item4)
	qtree.Add(itemFar)

	expected := []Item[float64]{item1, item2, item3, item4}
	neighbors := qtree.FindNeighbors(target, 1.0)

	if !sliceutils.SameElements(neighbors, expected) {
		t.Errorf("result %v not equal to expected %v", neighbors, expected)
	}
}

func TestQuadTreeBox_ForBoundedPlane(t *testing.T) {
	boundedPlane := geometry.NewBoundedPlane(4, 4)
	qtree := NewQuadTree(boundedPlane)
	defer qtree.Close()

	target := newTestBoxItem(0, 0, 1, 1)
	item1 := newTestBoxItem(0, 1, 1, 2)
	item2 := newTestBoxItem(1, 0, 2, 1)
	item3 := newTestBoxItem(3, 0, 4, 1)

	qtree.Add(item1)
	qtree.Add(item2)
	qtree.Add(item3)

	expected := []Item[int]{item1, item2}
	neighbors := qtree.FindNeighbors(target, 1)

	if !sliceutils.SameElements(neighbors, expected) {
		t.Errorf("result %v not equal to expected %v", neighbors, expected)
	}
}

func TestQuadTreeBox_ForCyclicBoundedPlane(t *testing.T) {
	cyclicPlane := geometry.NewCyclicBoundedPlane(4, 4)
	qtree := NewQuadTree(cyclicPlane)
	defer qtree.Close()

	target := newTestBoxItem(0, 0, 1, 1)
	item1 := newTestBoxItem(0, 3, 1, 4) // wrap od góry (ale przy margin=1 nie wejdzie)
	item2 := newTestBoxItem(3, 0, 4, 1) // wrap od lewej (ale przy margin=1 nie wejdzie)
	item3 := newTestBoxItem(1, 0, 2, 1) // normalny sąsiad
	item4 := newTestBoxItem(0, 1, 1, 2) // normalny sąsiad

	qtree.Add(item1)
	qtree.Add(item2)
	qtree.Add(item3)
	qtree.Add(item4)

	// Przy margin=1 znajdziemy tylko item3 i item4
	expected := []Item[int]{item3, item4}
	neighbors := qtree.FindNeighbors(target, 1)

	if !sliceutils.SameElements(neighbors, expected) {
		t.Errorf("result %v not equal to expected %v", neighbors, expected)
	}
}

func TestQuadTreeBox_ForCyclicBoundedPlane_WithWraps(t *testing.T) {
	cyclicPlane := geometry.NewCyclicBoundedPlane(4, 4)
	qtree := NewQuadTree(cyclicPlane)
	defer qtree.Close()

	target := newTestBoxItem(0, 0, 1, 1)
	item1 := newTestBoxItem(0, 3, 1, 4) // wrap od góry
	item2 := newTestBoxItem(3, 0, 4, 1) // wrap od lewej
	item3 := newTestBoxItem(1, 0, 2, 1) // normalny sąsiad
	item4 := newTestBoxItem(0, 1, 1, 2) // normalny sąsiad

	qtree.Add(item1)
	qtree.Add(item2)
	qtree.Add(item3)
	qtree.Add(item4)

	// Przy margin=2 znajdziemy także boxy wrapowane
	expected := []Item[int]{item1, item2, item3, item4}
	neighbors := qtree.FindNeighbors(target, 2)

	if !sliceutils.SameElements(neighbors, expected) {
		t.Errorf("result %v not equal to expected %v", neighbors, expected)
	}
}
