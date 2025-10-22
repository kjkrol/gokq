package quadtree

import "github.com/kjkrol/gokg/pkg/geometry"

type Node[V SpatialValue[T], T geometry.SupportedNumeric] struct {
	box    Box[T]
	items  []Item[V, T]
	parent *Node[V, T]
	childs []*Node[V, T]
}

func newNode[V SpatialValue[T], T geometry.SupportedNumeric](box Box[T], parent *Node[V, T]) *Node[V, T] {
	return &Node[V, T]{box: box, items: make([]Item[V, T], 0), parent: parent}
}

func (n *Node[V, T]) isLeaf() bool { return len(n.childs) == 0 }
func (n *Node[V, T]) isNode() bool { return len(n.childs) > 0 }

func (n *Node[V, T]) add(item Item[V, T], ops NeighborOps[V, T]) {
	if n.isNode() {
		if child := n.findFittingChild(ops.BoundsOf(item.Value())); child != nil {
			child.add(item, ops)
			return
		}
	}
	n.items = append(n.items, item)
	if len(n.items) > CAPACITY && n.isLeaf() {
		n.createChilds()
		n.redistribute(ops)
	}
}

func (n *Node[V, T]) remove(item Item[V, T], ops NeighborOps[V, T]) bool {
	// jeśli to węzeł wewnętrzny, spróbuj zejść do dziecka
	if n.isNode() {
		if child := n.findFittingChild(ops.BoundsOf(item.Value())); child != nil {
			if child.remove(item, ops) {
				// po usunięciu sprawdzamy czy dzieci nie są puste
				n.tryCompress(CAPACITY)
				return true
			}
		}
	}

	// sprawdzamy w bieżącym węźle
	for i, it := range n.items {
		if it == item { // porównanie referencji
			n.items = append(n.items[:i], n.items[i+1:]...)
			// jeśli nie ma dzieci i brak itemów → też można oczyścić
			n.tryCompress(CAPACITY)
			return true
		}
	}

	return false
}

func (n *Node[V, T]) tryCompress(capacity int) {
	if !n.isNode() {
		return
	}

	// rekurencyjnie zbierz wszystkie itemy w poddrzewie
	collected := collectItems(n)

	if len(collected) <= capacity {
		n.items = collected
		n.childs = nil
	}
}

func collectItems[V SpatialValue[T], T geometry.SupportedNumeric](n *Node[V, T]) []Item[V, T] {
	items := append([]Item[V, T]{}, n.items...)
	for _, ch := range n.childs {
		items = append(items, collectItems(ch)...)
	}
	return items
}

func (n *Node[V, T]) findFittingChild(b Box[T]) *Node[V, T] {
	for _, child := range n.childs {
		if child.box.ContainsBox(b) {
			return child
		}
	}
	return nil
}
func (n *Node[V, T]) redistribute(ops NeighborOps[V, T]) {
	remaining := make([]Item[V, T], 0, len(n.items))
	moved := 0

	for _, item := range n.items {
		if child := n.findFittingChild(ops.BoundsOf(item.Value())); child != nil {
			child.add(item, ops)
			moved++
		} else {
			remaining = append(remaining, item)
		}
	}
	n.items = remaining

	// rollback: jeżeli żaden item nie został przesunięty do dzieci,
	// to nie ma sensu utrzymywać dzieci
	if moved == 0 {
		// scal z powrotem: wszystkie itemy do rodzica
		for _, ch := range n.childs {
			n.items = append(n.items, ch.items...)
		}
		n.childs = nil
	}
}

func (n *Node[V, T]) createChilds() {
	childBoxes := n.box.Split()
	n.childs = make([]*Node[V, T], 4)
	for i, box := range childBoxes {
		n.childs[i] = newNode(box, n)
	}
}

func (n *Node[V, T]) close() {
	for _, child := range n.childs {
		child.close()
	}
	n.items = nil
	n.childs = nil
	n.parent = nil
}

func (n *Node[V, T]) findIntersectingNodesUnique(probe Box[T], set map[*Node[V, T]]struct{}) {
	if !n.box.Intersects(probe) {
		return
	}
	set[n] = struct{}{}
	for _, child := range n.childs {
		child.findIntersectingNodesUnique(probe, set)
	}
}

func (n *Node[V, T]) allItems() []Item[V, T] {
	items := append([]Item[V, T]{}, n.items...)
	for _, ch := range n.childs {
		items = append(items, ch.allItems()...)
	}
	return items
}

func (n *Node[V, T]) depth() int {
	if n.isLeaf() {
		return 1
	}
	maxChildDepth := 0
	for _, ch := range n.childs {
		if d := ch.depth(); d > maxChildDepth {
			maxChildDepth = d
		}
	}
	return 1 + maxChildDepth
}

func (n *Node[V, T]) leafBoxes() []Box[T] {
	if n.isLeaf() {
		return []Box[T]{n.box}
	}
	var boxes []Box[T]
	for _, ch := range n.childs {
		boxes = append(boxes, ch.leafBoxes()...)
	}
	return boxes
}
