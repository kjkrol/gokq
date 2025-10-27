package quadtree

import (
	"github.com/kjkrol/gokg/pkg/geometry/spatial"
)

type Node[T spatial.SupportedNumeric] struct {
	bounds spatial.Rectangle[T]
	items  []Item[T]
	parent *Node[T]
	childs []*Node[T]
}

func newNode[T spatial.SupportedNumeric](bounds spatial.Rectangle[T], parent *Node[T]) *Node[T] {
	return &Node[T]{bounds: bounds, items: make([]Item[T], 0), parent: parent}
}

func (n *Node[T]) isLeaf() bool { return len(n.childs) == 0 }
func (n *Node[T]) isNode() bool { return len(n.childs) > 0 }

func (n *Node[T]) add(item Item[T]) {
	if n.isNode() {
		if child := n.findFittingChild(item.Value().Bounds()); child != nil {
			child.add(item)
			return
		}
	}
	n.items = append(n.items, item)
	if len(n.items) > CAPACITY && n.isLeaf() {
		n.createChilds()
		n.redistribute()
	}
}

func (n *Node[T]) remove(item Item[T]) bool {
	if n.isNode() {
		if child := n.findFittingChild(item.Value().Bounds()); child != nil {
			if child.remove(item) {
				n.tryCompress(CAPACITY)
				return true
			}
		}
	}

	for i, it := range n.items {
		if it == item {
			n.items = append(n.items[:i], n.items[i+1:]...)
			n.tryCompress(CAPACITY)
			return true
		}
	}

	return false
}

func (n *Node[T]) tryCompress(capacity int) {
	if !n.isNode() {
		return
	}

	collected := collectItems(n)
	if len(collected) <= capacity {
		n.items = collected
		n.childs = nil
	}
}

func collectItems[T spatial.SupportedNumeric](n *Node[T]) []Item[T] {
	items := append([]Item[T]{}, n.items...)
	for _, ch := range n.childs {
		items = append(items, collectItems(ch)...)
	}
	return items
}

func (n *Node[T]) findFittingChild(r spatial.Rectangle[T]) *Node[T] {
	for _, child := range n.childs {
		if child.bounds.Contains(r) {
			return child
		}
	}
	return nil
}

func (n *Node[T]) redistribute() {
	remaining := make([]Item[T], 0, len(n.items))
	moved := 0

	for _, item := range n.items {
		if child := n.findFittingChild(item.Value().Bounds()); child != nil {
			child.add(item)
			moved++
		} else {
			remaining = append(remaining, item)
		}
	}
	n.items = remaining

	if moved == 0 {
		for _, ch := range n.childs {
			n.items = append(n.items, ch.items...)
		}
		n.childs = nil
	}
}

func (n *Node[T]) createChilds() {
	childRectangles := n.bounds.Split()
	n.childs = make([]*Node[T], 4)
	for i, rect := range childRectangles {
		n.childs[i] = newNode[T](rect, n)
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

// TODO: wydaje mi sie, ze ta metoda jest nie optymalna; dobra dla cyklicznej przestrzeni; ale dla bounded
// absolutne nie efektywan - trzeba to przemyslec i napisac lepiej (chyba szybka poprawka; ale trzeba sie skupic)
func (n *Node[T]) findIntersectingNodesUnique(probe spatial.Rectangle[T], set map[*Node[T]]struct{}) {
	if !n.bounds.Intersects(probe) {
		return
	}
	set[n] = struct{}{}
	for _, child := range n.childs {
		child.findIntersectingNodesUnique(probe, set)
	}
}

func (n *Node[T]) allItems() []Item[T] {
	items := append([]Item[T]{}, n.items...)
	for _, ch := range n.childs {
		items = append(items, ch.allItems()...)
	}
	return items
}

func (n *Node[T]) depth() int {
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

func (n *Node[T]) leafRectangles() []spatial.Rectangle[T] {
	if n.isLeaf() {
		return []spatial.Rectangle[T]{n.bounds}
	}
	var rectangles []spatial.Rectangle[T]
	for _, ch := range n.childs {
		rectangles = append(rectangles, ch.leafRectangles()...)
	}
	return rectangles
}
