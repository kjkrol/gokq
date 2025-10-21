package boxes

import (
	quadcore "github.com/kjkrol/goka/pkg/quadtree/base"
	"github.com/kjkrol/gokg/pkg/geometry"
)

// QuadTree reprezentuje drzewo czwórkowe dla indeksowania przestrzennego.
type QuadTree[T geometry.SupportedNumeric] struct {
	root  *Node[T]
	plane geometry.Plane[T]
}

// NewQuadTree tworzy nowe drzewo quadtree dla podanej płaszczyzny.
func NewQuadTree[T geometry.SupportedNumeric](plane geometry.Plane[T]) QuadTree[T] {
	box := quadcore.NewBox(geometry.Vec[T]{X: 0, Y: 0}, plane.Size())
	root := newNode(box, nil)
	return QuadTree[T]{root, plane}
}

// Add dodaje element (o zadanym prostokącie) do quadtree.
func (t *QuadTree[T]) Add(item Item[T]) {
	t.root.add(item)
}

// Close zwalnia zasoby powiązane z quadtree.
func (t *QuadTree[T]) Close() {
	t.root.close()
}

// FindNeighbors zwraca elementy, których boxy przecinają się z powiększonym boxem targeta.
func (t *QuadTree[T]) FindNeighbors(target Item[T], margin T) []Item[T] {
	// 1. bierzemy granice targeta i je powiększamy o margin
	expanded := target.Bounds().Expand(margin)

	// 2. znajdujemy kandydatów w węzłach, których box przecina expanded
	neighborNodes := selectingNeighbors(expanded, t)

	// 3. filtrujemy — zostawiamy tylko te, które naprawdę się przecinają
	return scan(neighborNodes, func(item *Item[T]) bool {
		dist := boxDistance(target.Bounds(), (*item).Bounds(), t.plane.Metric)
		return dist <= margin
	})
}

// func boxDistance[T geometry.SupportedNumeric](
// 	a, b quadcore.Box[T],
// 	metric func(geometry.Vec[T], geometry.Vec[T]) T,
// ) T {
// 	// jeśli się przecinają, dystans = 0
// 	if a.Intersects(b) {
// 		return 0
// 	}
// 	// inaczej bierz najbliższe punkty z obu boxów
// 	// tu uproszczona wersja: licz dystans między środkami
// 	// ALE można zrobić dokładniej: min odległość między krawędziami
// 	return metric(a.Center, b.Center)
// }

func boxDistance[T geometry.SupportedNumeric](
	a, b quadcore.Box[T],
	metric func(geometry.Vec[T], geometry.Vec[T]) T,
) T {
	// jeśli się przecinają, dystans = 0
	if a.Intersects(b) {
		return 0
	}

	// minimalny dystans na osi X
	dx := axisDistance(a, b, func(v geometry.Vec[T]) T { return v.X })

	// minimalny dystans na osi Y
	dy := axisDistance(a, b, func(v geometry.Vec[T]) T { return v.Y })

	// wektor dystansu (dx, dy) porównujemy do (0,0) w wybranej metryce
	return metric(
		geometry.Vec[T]{X: dx, Y: dy},
		geometry.Vec[T]{X: 0, Y: 0},
	)
}

// axisDistance zwraca minimalny odstęp między przedziałami dwóch boxów na danej osi.
// Działa analogicznie do axisIntersection, ale zwraca wartość liczbową zamiast bool.
func axisDistance[T geometry.SupportedNumeric](
	aa, bb quadcore.Box[T],
	axisValue func(geometry.Vec[T]) T,
) T {
	aa, bb = quadcore.SortBy(aa, bb, axisValue)

	// jeśli zachodzi nakładanie → odstęp = 0
	if axisValue(aa.BottomRight) >= axisValue(bb.TopLeft) {
		return 0
	}

	// w przeciwnym razie odstęp to różnica między "początkiem bb" a "końcem aa"
	return axisValue(bb.TopLeft) - axisValue(aa.BottomRight)
}

//-----------------------------------------------------------------------------

// Item reprezentuje obiekt przechowywany w quadtree.
type Item[T geometry.SupportedNumeric] interface {
	Bounds() quadcore.Box[T]
}

// Node reprezentuje węzeł quadtree.
type Node[T geometry.SupportedNumeric] struct {
	box    quadcore.Box[T]
	items  []Item[T]
	parent *Node[T]
	childs []*Node[T]
}

func newNode[T geometry.SupportedNumeric](box quadcore.Box[T], parent *Node[T]) *Node[T] {
	items := make([]Item[T], 0)
	return &Node[T]{box: box, items: items, parent: parent}
}

func (n *Node[T]) isLeaf() bool { return len(n.childs) == 0 }
func (n *Node[T]) isNode() bool { return len(n.childs) > 0 }

// add dodaje obiekt do węzła (lub jego potomków).
func (n *Node[T]) add(item Item[T]) {
	if n.isNode() {
		if child := n.findFittingChild(item.Bounds()); child != nil {
			child.add(item)
			return
		}
	}

	n.items = append(n.items, item)

	// Rozdzielamy, jeśli node jest przepełniony i jeszcze nie ma dzieci
	if len(n.items) > 3 && n.isLeaf() {
		n.createChilds()
		n.redistribute()
	}
}

// findFittingChild zwraca dziecko, w którym box mieści się w całości.
func (n *Node[T]) findFittingChild(b quadcore.Box[T]) *Node[T] {
	for _, child := range n.childs {
		if child.box.ContainsBox(b) {
			return child
		}
	}
	return nil
}

// redistribute próbuje przepchnąć istniejące itemy do dzieci, jeśli się mieszczą.
func (n *Node[T]) redistribute() {
	remaining := make([]Item[T], 0, len(n.items))
	for _, item := range n.items {
		if child := n.findFittingChild(item.Bounds()); child != nil {
			child.add(item)
		} else {
			remaining = append(remaining, item)
		}
	}
	n.items = remaining
}

func (n *Node[T]) createChilds() {
	var childBoxes [4]quadcore.Box[T] = n.box.Split()
	n.childs = make([]*Node[T], 4)
	for i, box := range childBoxes {
		n.childs[i] = newNode(box, n)
	}
}

func (n *Node[T]) close() {
	for _, child := range n.childs {
		child.close()
	}
	n.items = nil
	n.childs = nil
	n.parent = nil
}

//-----------------------------------------------------------------------------

// selectingNeighbors działa jak wcześniej, ale operuje na boxach.
func selectingNeighbors[T geometry.SupportedNumeric](probeBox quadcore.Box[T], t *QuadTree[T]) []*Node[T] {
	neighborNodes := make([]*Node[T], 0)
	probeBoxes := []quadcore.Box[T]{probeBox}
	if t.plane.Name() == "cyclic" {
		wrappedBoxes := quadcore.WrapBoxCyclic(probeBox, t.plane.Size(), t.plane.Contains)
		probeBoxes = append(probeBoxes, wrappedBoxes...)
	}
	for _, pBox := range probeBoxes {
		findIntersectingNodes(t.root, pBox, &neighborNodes)
	}

	return neighborNodes
}

func findIntersectingNodes[T geometry.SupportedNumeric](node *Node[T], probeBox quadcore.Box[T], neighborNodes *[]*Node[T]) {
	// Jeżeli węzeł w ogóle nie przecina obszaru sondy, to nie schodzimy w dół.
	if !node.box.Intersects(probeBox) {
		return
	}
	// Dodajemy węzeł bez względu na to, czy jest liściem — bo itemy mogą być też w węźle wewnętrznym.
	*neighborNodes = append(*neighborNodes, node)

	// Schodzimy do dzieci (jeśli są).
	for _, child := range node.childs {
		findIntersectingNodes(child, probeBox, neighborNodes)
	}
}

func scan[T geometry.SupportedNumeric](neighborNodes []*Node[T], predicate func(*Item[T]) bool) []Item[T] {
	neighborItems := make([]Item[T], 0)
	for _, node := range neighborNodes {
		for _, neighborItem := range node.items {
			if predicate(&neighborItem) {
				neighborItems = append(neighborItems, neighborItem)
			}
		}
	}
	return neighborItems
}
