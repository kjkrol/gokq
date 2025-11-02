package quadtree

import "github.com/kjkrol/gokg/pkg/geometry"

type Node[T geometry.SupportedNumeric] struct {
	bounds geometry.AABB[T]
	items  []Item[T]
	parent *Node[T]
	childs []*Node[T]
}

func newNode[T geometry.SupportedNumeric](bounds geometry.AABB[T], parent *Node[T]) *Node[T] {
	return &Node[T]{bounds: bounds, items: make([]Item[T], 0), parent: parent}
}

func (n *Node[T]) isLeaf() bool { return len(n.childs) == 0 }
func (n *Node[T]) isNode() bool { return len(n.childs) > 0 }

func (n *Node[T]) findFittingChild(r geometry.AABB[T]) *Node[T] {
	for _, child := range n.childs {
		if child.bounds.Contains(r) {
			return child
		}
	}
	return nil
}

func (n *Node[T]) close() {
	for _, child := range n.childs {
		child.close()
	}
	n.items = nil
	n.childs = nil
	n.parent = nil
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

func (n *Node[T]) leafRectangles() []geometry.AABB[T] {
	if n.isLeaf() {
		return []geometry.AABB[T]{n.bounds}
	}
	var rectangles []geometry.AABB[T]
	for _, ch := range n.childs {
		rectangles = append(rectangles, ch.leafRectangles()...)
	}
	return rectangles
}
