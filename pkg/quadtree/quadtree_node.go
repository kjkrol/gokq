package quadtree

import (
	"github.com/kjkrol/gokg/pkg/geometry"
	"github.com/kjkrol/gokq/pkg/dfs"
)

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

func (n *Node[T, K]) Children() []*Node[T, K] {
	return n.childs
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
	items := []Item[T, K]{}

	dfs.DFS(n, struct{}{}, func(node *Node[T, K], _ struct{}) (dfs.DFSControl, struct{}) {
		items = append(items, node.items...)
		return dfs.DFSControl{}, struct{}{}
	})

	return items
}

func (n *Node[T, K]) depth() int {
	maxDepth := 0

	dfs.DFS(n, 0, func(node *Node[T, K], acc int) (dfs.DFSControl, int) {
		acc = acc + 1
		if acc > maxDepth {
			maxDepth = acc
		}
		return dfs.DFSControl{}, acc
	})

	return maxDepth
}

func (n *Node[T, K]) leafBounds() []geometry.BoundingBox[T] {
	rectangles := []geometry.BoundingBox[T]{}

	dfs.DFS(n, struct{}{}, func(node *Node[T, K], _ struct{}) (dfs.DFSControl, struct{}) {
		if node.isLeaf() {
			rectangles = append(rectangles, node.bounds)
			return dfs.DFSControl{Skip: true}, struct{}{}
		}
		return dfs.DFSControl{}, struct{}{}
	})

	return rectangles
}
