package quadtree

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geometry"
	"github.com/kjkrol/goku/pkg/sliceutils"
)

func newTestBoxItem[T geometry.SupportedNumeric](x1, y1, x2, y2 T) *ExampleItem[T] {
	rect := geometry.NewRectangle(
		geometry.Vec[T]{X: x1, Y: y1},
		geometry.Vec[T]{X: x2, Y: y2},
	)
	return &ExampleItem[T]{spatial: &rect}
}

func TestQuadTreeBox_FindNeighborsSimple(t *testing.T) {
	boundedPlane := geometry.NewBoundedPlane(64.0, 64.0)
	qtree := NewQuadTree(boundedPlane, geometry.BoundingBoxDistanceForPlane(boundedPlane))
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
	qtree := NewQuadTree(boundedPlane, geometry.BoundingBoxDistanceForPlane(boundedPlane))
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
	qtree := NewQuadTree(cyclicPlane, geometry.BoundingBoxDistanceForPlane(cyclicPlane))
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
	qtree := NewQuadTree(cyclicPlane, geometry.BoundingBoxDistanceForPlane(cyclicPlane))
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

func TestQuadTree_RemoveCascadeCompression_Box(t *testing.T) {
	plane := geometry.NewBoundedPlane(64.0, 64.0)
	qtree := NewQuadTree(plane, geometry.BoundingBoxDistanceForPlane(plane))
	defer qtree.Close()

	makeBox := func(x, y float64) *ExampleItem[float64] {
		// malutki box (1x1), żeby dało się wcisnąć w child
		rect := geometry.BuildRectangle(
			geometry.Vec[float64]{X: x, Y: y}, 0.5,
		)
		return &ExampleItem[float64]{spatial: &rect}
	}

	// Dodajemy 4 elementy -> root nadal liść
	items := []*ExampleItem[float64]{
		makeBox(1, 1),
		makeBox(2, 2),
		makeBox(3, 3),
		makeBox(4, 4),
	}
	for _, it := range items {
		qtree.Add(it)
	}
	if !qtree.root.isLeaf() {
		t.Fatalf("expected root to be leaf after 4 inserts")
	}

	// 5-ty element -> root się dzieli
	item5 := makeBox(5, 5)
	qtree.Add(item5)
	if qtree.root.isLeaf() {
		t.Fatalf("expected root to split after 5th insert")
	}

	// dodajemy kilka kolejnych boxów w tym samym rejonie
	item6 := makeBox(6, 6)
	item7 := makeBox(7, 7)
	item8 := makeBox(8, 8)
	qtree.Add(item6)
	qtree.Add(item7)
	qtree.Add(item8)

	child := qtree.root.childs[0]
	if child == nil {
		t.Fatalf("expected root.childs[0] to exist")
	}

	// 9-ty element -> powoduje split childa
	item9 := makeBox(9, 9)
	qtree.Add(item9)
	if child.isLeaf() {
		t.Fatalf("expected child to split after 9th insert")
	}

	// Usuwamy item5..9 -> root powinien się skompresować (wrócić do liścia z itemami)
	for _, it := range []*ExampleItem[float64]{item5, item6, item7, item8, item9} {
		removed := qtree.Remove(it)
		if !removed {
			t.Fatalf("expected %+v to be removed", it)
		}
	}

	// root powinien znów być liściem (choć z itemami, bo boxy mogły nie spłynąć do dzieci)
	if !qtree.root.isLeaf() {
		t.Errorf("expected root to compress back to leaf after removing 5..9")
	}
}

func TestQuadTree_BoxItems_LargeStayInParent_SmallGoToChildren(t *testing.T) {
	plane := geometry.NewBoundedPlane(64.0, 64.0)
	qtree := NewQuadTree(plane, geometry.BoundingBoxDistanceForPlane(plane))
	defer qtree.Close()

	// Dodajemy 6 dużych boxów, każdy obejmuje połowę przestrzeni
	for i := 0; i < 6; i++ {
		rect := geometry.NewRectangle(
			geometry.Vec[float64]{X: 0, Y: 0},
			geometry.Vec[float64]{X: 64, Y: 32}, // duży box, nie mieści się w jednym childzie
		)
		large := &ExampleItem[float64]{
			spatial: &rect,
		}
		qtree.Add(large)
	}

	// Root powinien być nadal liściem, bo żaden box nie mieścił się w childach
	if !qtree.root.isLeaf() {
		t.Fatalf("expected root to remain leaf for large boxes")
	}
	if len(qtree.root.items) != 6 {
		t.Fatalf("expected 6 items in root, got %d", len(qtree.root.items))
	}

	// Teraz dodajemy kilka małych boxów, które zmieszczą się w ćwiartkach
	rectSmall1 := geometry.NewRectangle(
		geometry.Vec[float64]{X: 1, Y: 1},
		geometry.Vec[float64]{X: 2, Y: 2},
	)
	small1 := &ExampleItem[float64]{
		spatial: &rectSmall1,
	}
	rectSmall2 := geometry.NewRectangle(
		geometry.Vec[float64]{X: 10, Y: 10},
		geometry.Vec[float64]{X: 11, Y: 11},
	)
	small2 := &ExampleItem[float64]{
		spatial: &rectSmall2,
	}
	qtree.Add(small1)
	qtree.Add(small2)

	// Teraz root powinien się podzielić
	if qtree.root.isLeaf() {
		t.Fatalf("expected root to split after adding small boxes")
	}

	// Duże boxy powinny nadal siedzieć w root.items
	if len(qtree.root.items) != 6 {
		t.Errorf("expected large boxes to remain in root, got %d", len(qtree.root.items))
	}

	// Małe boxy powinny trafić do dzieci
	childItems := 0
	for _, ch := range qtree.root.childs {
		childItems += len(ch.items)
	}
	if childItems < 2 {
		t.Errorf("expected small boxes to be distributed to children, got %d", childItems)
	}
}

func TestQuadTree_Box_CountDepthAllItemsLeafRectangles(t *testing.T) {
	plane := geometry.NewBoundedPlane(16.0, 16.0)
	qtree := NewQuadTree(plane, geometry.BoundingBoxDistanceForPlane(plane))
	defer qtree.Close()

	// początkowo puste
	if qtree.Count() != 0 {
		t.Errorf("expected empty tree count=0, got %d", qtree.Count())
	}
	if qtree.Depth() != 1 {
		t.Errorf("expected depth=1 for empty tree, got %d", qtree.Depth())
	}

	// dodajemy 2 duże boxy (obejmują większą część przestrzeni, nie zmieszczą się w childach)

	large1 := newTestBoxItem(0., 0, 16, 8) // górna połowa

	large2 := newTestBoxItem(0., 8, 16, 16) // dolna połowa

	qtree.Add(large1)
	qtree.Add(large2)

	if qtree.Count() != 2 {
		t.Errorf("expected count=2, got %d", qtree.Count())
	}
	// nadal depth=1, bo duże boxy nie mogą zejść do childów
	if qtree.Depth() != 1 {
		t.Errorf("expected depth=1 with only large boxes, got %d", qtree.Depth())
	}

	// dodajemy małe boxy, które zmieszczą się w childach
	smallBoxes := []*ExampleItem[float64]{
		newTestBoxItem(1., 1, 2, 2),
		newTestBoxItem(3., 3, 4, 4),
		newTestBoxItem(14., 14, 15, 15),
		newTestBoxItem(12., 1, 13, 2),
	}
	for _, sb := range smallBoxes {
		qtree.Add(sb)
	}

	if qtree.Count() != 6 {
		t.Errorf("expected count=6 (2 large + 4 small), got %d", qtree.Count())
	}

	// teraz powinien być split, bo małe boxy trafiły do children
	if qtree.Depth() <= 1 {
		t.Errorf("expected depth>1 after adding small boxes, got %d", qtree.Depth())
	}

	// AllItems powinno zawierać wszystkie 6 elementów
	all := qtree.AllItems()
	if len(all) != 6 {
		t.Errorf("expected 6 items in AllItems, got %d", len(all))
	}

	// LeafBoxes powinny być >1, bo root się podzielił
	leafs := qtree.LeafRectangles()
	if len(leafs) <= 1 {
		t.Errorf("expected more than 1 leaf box after split, got %d", len(leafs))
	}
	// wszystkie leaf boxy muszą być w obrębie płaszczyzny
	for _, lb := range leafs {
		if !plane.Contains(lb.TopLeft) && !plane.Contains(lb.BottomRight) {
			t.Errorf("leaf box %+v is outside plane", lb)
		}
	}
}

func TestSortNeighbors_BottomRightTieBreak(t *testing.T) {
	plane := geometry.NewBoundedPlane(16.0, 16.0)
	qtree := NewQuadTree(plane, geometry.BoundingBoxDistanceForPlane(plane))
	defer qtree.Close()

	// Box A i B mają identyczne TopLeft
	rectA := geometry.NewRectangle(
		geometry.Vec[float64]{X: 1, Y: 1},
		geometry.Vec[float64]{X: 3, Y: 3},
	)
	a := &ExampleItem[float64]{spatial: &rectA}
	rectB := geometry.NewRectangle(
		geometry.Vec[float64]{X: 1, Y: 1},
		geometry.Vec[float64]{X: 3, Y: 4}, // różni się tylko BottomRight.Y
	)
	b := &ExampleItem[float64]{spatial: &rectB}
	rectC := geometry.NewRectangle(
		geometry.Vec[float64]{X: 1, Y: 1},
		geometry.Vec[float64]{X: 4, Y: 3}, // różni się tylko BottomRight.X
	)
	c := &ExampleItem[float64]{spatial: &rectC}

	items := []Item[float64]{b, a, c}

	// sortujemy
	sortNeighbors(items)

	// oczekiwana kolejność:
	// najpierw a (BottomRight.Y=3),
	// potem b (BottomRight.Y=4),
	// na końcu c (BottomRight.Y=3 ale BottomRight.X=4 > 3).
	expected := []Item[float64]{a, c, b}

	for i, it := range items {
		if it != expected[i] {
			t.Errorf("unexpected order at %d: got %+v, expected %+v", i, it, expected[i])
		}
	}
}
