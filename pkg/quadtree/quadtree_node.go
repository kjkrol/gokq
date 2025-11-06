package quadtree

import "github.com/kjkrol/gokg/pkg/geometry"

type Node[T geometry.SupportedNumeric, K comparable] struct {
	bounds geometry.BoundingBox[T]
	items  []Item[T, K]
	parent *Node[T, K]
	childs []*Node[T, K]
}

func newNode[T geometry.SupportedNumeric, K comparable](bounds geometry.BoundingBox[T], parent *Node[T, K]) *Node[T, K] {
	return &Node[T, K]{bounds: bounds, items: make([]Item[T, K], 0), parent: parent}
}

func (n *Node[T, K]) isLeaf() bool { return len(n.childs) == 0 }
func (n *Node[T, K]) isNode() bool { return len(n.childs) > 0 }

func (n *Node[T, K]) findFittingChild(r geometry.BoundingBox[T]) *Node[T, K] {
	for _, child := range n.childs {
		if child.bounds.Contains(r) {
			return child
		}
	}
	return nil
}

func (n *Node[T, K]) close() {
	for _, child := range n.childs {
		child.close()
	}
	n.items = nil
	n.childs = nil
	n.parent = nil
}

func (n *Node[T, K]) allItems() []Item[T, K] {
	items := append([]Item[T, K]{}, n.items...)
	for _, ch := range n.childs {
		items = append(items, ch.allItems()...)
	}
	return items
}

func (n *Node[T, K]) depth() int {
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

func (n *Node[T, K]) leafRectangles() []geometry.BoundingBox[T] {
	if n.isLeaf() {
		return []geometry.BoundingBox[T]{n.bounds}
	}
	var rectangles []geometry.BoundingBox[T]
	for _, ch := range n.childs {
		rectangles = append(rectangles, ch.leafRectangles()...)
	}
	return rectangles
}
