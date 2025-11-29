package qtree

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
	"github.com/kjkrol/gokg/pkg/plane"
	"github.com/kjkrol/goku/pkg/sliceutils"
)

func TestQuadTree_Probe_For_BoundedPlane(t *testing.T) {
	boundedPlane := plane.NewEuclidean2D(100, 100)
	qtree := NewQuadTree(boundedPlane)
	box := geom.NewAABBAt(geom.NewVec(0, 0), 2, 2)
	item := newTestItemFromBox(box)

	nodeIntersectionDetection := qtree.finder.strategy.NodeIntersectionDetectionFactory(item, 1)
	nodeIntersectionDetection(*qtree.root)
	if !nodeIntersectionDetection(*qtree.root) {
		t.Fatalf("expected intersection with root node")
	}
}

func TestQuadTree_FindNeighbors_ForCyclicBoundedPlane_WithFrags(t *testing.T) {
	torus := plane.NewToroidal2D(4, 4)
	qtree := NewQuadTree(torus)
	defer qtree.Close()

	target := newTestItemFromPos(0, 0, 1, 1)

	planeBox1 := torus.WrapAABB(geom.NewAABBAt(geom.NewVec(0, 0), 2, 2))

	torus.Translate(&planeBox1, geom.NewVec(0, 1))

	item1 := newTestItemFromBox(planeBox1.AABB) // wrap od góry
	qtree.Add(item1)

	planeBox1.VisitFragments(func(pos plane.FragPosition, aabb geom.AABB[int]) bool {
		fragItem := newTestItemFromBox(aabb)
		t.Log(fragItem.Bound().String())
		qtree.Add(fragItem)
		return true
	})

	planeBox2 := torus.WrapAABB(geom.NewAABBAt(geom.NewVec(0, 0), 4, 1))

	torus.Translate(&planeBox2, geom.NewVec(1, 0))

	item2 := newTestItemFromBox(planeBox2.AABB) // wrap od lewej

	frag2Items := make(map[plane.FragPosition]*TestItem[int], 4)
	planeBox2.VisitFragments(func(pos plane.FragPosition, aabb geom.AABB[int]) bool {
		frag2Items[pos] = newTestItemFromBox(aabb)
		t.Log(frag2Items[pos].Bound().String())
		qtree.Add(frag2Items[pos])
		return true
	})

	qtree.Add(item2)

	item3 := newTestItemFromPos(1, 0, 2, 1) // normalny sąsiad
	qtree.Add(item3)

	item4 := newTestItemFromPos(0, 1, 1, 2) // normalny sąsiad
	qtree.Add(item4)

	// Przy margin=2 znajdziemy także boxy wrapowane
	expected := []Item[int]{frag2Items[0], item3, item1, item4, item2}
	neighbors := qtree.FindNeighbors(target, 2)

	if !sliceutils.SameElements(neighbors, expected) {
		t.Errorf("result %v not equal to expected %v", neighbors, expected)
	}
}
