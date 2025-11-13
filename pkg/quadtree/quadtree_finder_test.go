package quadtree

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geometry"
	"github.com/kjkrol/goku/pkg/sliceutils"
)

func TestQuadTree_Probe_For_BoundedPlane(t *testing.T) {
	boundedPlane := geometry.NewBoundedPlane(100, 100)
	qtree := NewQuadTree[int, uint64](boundedPlane)
	box := geometry.NewBoundingBoxAt(geometry.NewVec(0, 0), 2, 2)
	item := newTestItemFromBox(box)

	nodeIntersectionDetection := qtree.finder.strategy.NodeIntersectionDetectionFactory(item, 1)
	nodeIntersectionDetection(*qtree.root)
	if !nodeIntersectionDetection(*qtree.root) {
		t.Fatalf("expected intersection with root node")
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
