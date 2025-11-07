package quadtree

import (
	"sort"
	"testing"

	"github.com/kjkrol/gokg/pkg/geometry"
	"github.com/kjkrol/goku/pkg/sliceutils"
)

func TestQuadTree_Probe_For_BoundedPlane(t *testing.T) {
	boundedPlane := geometry.NewBoundedPlane(100, 100)
	qtree := NewQuadTree[int, uint64](boundedPlane)
	box := geometry.NewBoundingBoxAt(geometry.NewVec(0, 0), 2, 2)

	probes := qtree.finder.probe(box, 1)
	if len(probes) != 1 {
		t.Fatalf("expected single probe, got %d", len(probes))
	}

	expectedTopLeft := geometry.NewVec(0, 0)
	expectedBottomRight := geometry.NewVec(3, 3)
	if probes[0].TopLeft != expectedTopLeft {
		t.Errorf("expected probe top-left %v, got %v", expectedTopLeft, probes[0].TopLeft)
	}
	if probes[0].BottomRight != expectedBottomRight {
		t.Errorf("expected probe bottom-right %v, got %v", expectedBottomRight, probes[0].BottomRight)
	}
}

func TestQuadTree_Probe_For_CyclicPlane(t *testing.T) {
	boundedPlane := geometry.NewCyclicBoundedPlane(10, 10)
	qtree := NewQuadTree[int, uint64](boundedPlane)
	box := geometry.NewBoundingBoxAt(geometry.NewVec(8, 8), 2, 2)
	probes := qtree.finder.probe(box, 0)
	if len(probes) > 1 {
		t.Fatalf("expected single probe, got %d", len(probes))
	}

	expectedTopLeft := geometry.NewVec(8, 8)
	expectedBottomRight := geometry.NewVec(10, 10)
	if probes[0].TopLeft != expectedTopLeft {
		t.Errorf("expected probe top-left %v, got %v", expectedTopLeft, probes[0].TopLeft)
	}
	if probes[0].BottomRight != expectedBottomRight {
		t.Errorf("expected probe bottom-right %v, got %v", expectedBottomRight, probes[0].BottomRight)
	}
}

func TestQuadTree_Probe_For_CyclicPlane_Edge_Case(t *testing.T) {
	boundedPlane := geometry.NewCyclicBoundedPlane(10, 10)
	qtree := NewQuadTree[int, uint64](boundedPlane)
	box := geometry.NewBoundingBoxAt(geometry.NewVec(0, 0), 2, 2)
	probes := qtree.finder.probe(box, 2)
	if len(probes) != 4 {
		t.Fatalf("expected 4 probes, got %d", len(probes))
	}

	want := []geometry.BoundingBox[int]{
		geometry.NewBoundingBox(geometry.NewVec(0, 0), geometry.NewVec(4, 4)),
		geometry.NewBoundingBox(geometry.NewVec(8, 0), geometry.NewVec(10, 4)),
		geometry.NewBoundingBox(geometry.NewVec(0, 8), geometry.NewVec(4, 10)),
		geometry.NewBoundingBox(geometry.NewVec(8, 8), geometry.NewVec(10, 10)),
	}

	sort.Slice(probes, func(i, j int) bool {
		ai, aj := probes[i], probes[j]
		first, _ := geometry.SortBoxesBy(
			ai, aj,
			func(box geometry.BoundingBox[int]) int { return box.TopLeft.Y },
			func(box geometry.BoundingBox[int]) int { return box.TopLeft.X },
			func(box geometry.BoundingBox[int]) int { return box.BottomRight.Y },
			func(box geometry.BoundingBox[int]) int { return box.BottomRight.X },
		)
		return first.Equals(ai)
	})

	for i := range want {
		if !probes[i].Equals(want[i]) {
			t.Fatalf("probe[%d] = %+v, want %+v", i, probes[i], want[i])
		}
	}
}

func TestQuadTree_FindNeighbors_ForCyclicBoundedPlane_WithFrags(t *testing.T) {
	cyclicPlane := geometry.NewCyclicBoundedPlane(4, 4)
	qtree := NewQuadTree[int, uint64](cyclicPlane)
	defer qtree.Close()

	target := newTestItemFromPos(0, 0, 1, 1)

	planeBox1 := cyclicPlane.WrapBoundingBox(geometry.NewBoundingBoxAt(geometry.NewVec(0, 0), 2, 2))

	cyclicPlane.Translate(&planeBox1, geometry.NewVec(0, 1))

	item1 := newTestItemFromBox(planeBox1.BoundingBox) // wrap od góry
	qtree.Add(item1)

	item1Frags := make([]*TestItem[int], len(planeBox1.Fragments()))
	i := 0
	for _, frag := range planeBox1.Fragments() {
		item1Frags[i] = newTestItemFromBox(frag)
		t.Log(item1Frags[i].Bound().String())
		qtree.Add(item1Frags[i])
		i++
	}

	planeBox2 := cyclicPlane.WrapBoundingBox(geometry.NewBoundingBoxAt(geometry.NewVec(0, 0), 4, 1))

	cyclicPlane.Translate(&planeBox2, geometry.NewVec(1, 0))

	item2 := newTestItemFromBox(planeBox2.BoundingBox) // wrap od lewej
	item2Frags := make([]*TestItem[int], len(planeBox2.Fragments()))

	i = 0
	for _, frag := range planeBox2.Fragments() {
		item2Frags[i] = newTestItemFromBox(frag)
		t.Log(item2Frags[i].Bound().String())
		qtree.Add(item2Frags[i])
		i++
	}
	qtree.Add(item2)

	item3 := newTestItemFromPos(1, 0, 2, 1) // normalny sąsiad
	qtree.Add(item3)

	item4 := newTestItemFromPos(0, 1, 1, 2) // normalny sąsiad
	qtree.Add(item4)

	// Przy margin=2 znajdziemy także boxy wrapowane
	expected := []Item[int, uint64]{item2Frags[0], item3, item1, item4, item2}
	neighbors := qtree.FindNeighbors(target, 2)

	if !sliceutils.SameElements(neighbors, expected) {
		t.Errorf("result %v not equal to expected %v", neighbors, expected)
	}
}
