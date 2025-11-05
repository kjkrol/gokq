package quadtree

import (
	"sort"
	"testing"

	"github.com/kjkrol/gokg/pkg/geometry"
	"github.com/kjkrol/goku/pkg/sliceutils"
)

func TestQuadTree_Probe_For_BoundedPlane(t *testing.T) {
	boundedPlane := geometry.NewBoundedPlane(100, 100)
	qtree := NewQuadTree(boundedPlane)
	aabb := geometry.NewAABB(geometry.NewVec(0, 0), 2, 2)

	probes := qtree.finder.probe(aabb, 1)
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
	qtree := NewQuadTree(boundedPlane)
	aabb := geometry.NewAABB(geometry.NewVec(8, 8), 2, 2)
	probes := qtree.finder.probe(aabb, 0)
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
	qtree := NewQuadTree(boundedPlane)
	aabb := geometry.NewAABB(geometry.NewVec(0, 0), 2, 2)
	probes := qtree.finder.probe(aabb, 2)
	if len(probes) != 4 {
		t.Fatalf("expected 4 probes, got %d", len(probes))
	}

	want := []geometry.AABB[int]{
		geometry.NewAABB(geometry.NewVec(0, 0), 4, 4),
		geometry.NewAABB(geometry.NewVec(8, 0), 2, 4),
		geometry.NewAABB(geometry.NewVec(0, 8), 4, 2),
		geometry.NewAABB(geometry.NewVec(8, 8), 2, 2),
	}

	sort.Slice(probes, func(i, j int) bool {
		ai, aj := probes[i], probes[j]
		first, _ := geometry.SortRectanglesBy(
			ai, aj,
			func(box geometry.AABB[int]) int { return box.TopLeft.Y },
			func(box geometry.AABB[int]) int { return box.TopLeft.X },
			func(box geometry.AABB[int]) int { return box.BottomRight.Y },
			func(box geometry.AABB[int]) int { return box.BottomRight.X },
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
	qtree := NewQuadTree(cyclicPlane)
	defer qtree.Close()

	target := newAABBItem(0, 0, 1, 1)
	item1 := newAABBItem(0, 2, 2, 2) // wrap od góry
	cyclicPlane.Translate(&item1.AABB, geometry.NewVec(0, 1))
	item1Frags := make([]*ShapeItem[int], len(item1.Bound().Fragments()))
	i := 0
	for _, frag := range item1.Bound().Fragments() {
		item1Frags[i] = &ShapeItem[int]{frag}
		t.Log(item1Frags[i].Bound().String())
		qtree.Add(item1Frags[i])
		i++
	}

	item2 := newAABBItem(0, 0, 4, 1) // wrap od lewej
	cyclicPlane.Translate(&item2.AABB, geometry.NewVec(1, 0))
	item2Frags := make([]*ShapeItem[int], len(item2.Bound().Fragments()))
	qtree.Add(item1)
	i = 0
	for _, frag := range item1.Bound().Fragments() {
		item2Frags[i] = &ShapeItem[int]{frag}
		t.Log(item2Frags[i].Bound().String())
		qtree.Add(item2Frags[i])
		i++
	}
	qtree.Add(item2)

	item3 := newAABBItem(1, 0, 2, 1) // normalny sąsiad
	qtree.Add(item3)

	item4 := newAABBItem(0, 1, 1, 2) // normalny sąsiad
	qtree.Add(item4)

	// Przy margin=2 znajdziemy także boxy wrapowane
	expected := []Item[int]{item1Frags[0], item2Frags[0], item3, item1, item4, item2}
	neighbors := qtree.FindNeighbors(target, 2)

	if !sliceutils.SameElements(neighbors, expected) {
		t.Errorf("result %v not equal to expected %v", neighbors, expected)
	}
}
