package qtree

import (
	"github.com/kjkrol/gokg/pkg/geom"
	"github.com/kjkrol/gokq/pkg/dfs"
)

type Node[T geom.Numeric] struct {
	bounds geom.AABB[T]
	items  []Item[T]
	parent *Node[T]
	childs []*Node[T]
}

func newNode[T geom.Numeric](bounds geom.AABB[T], parent *Node[T]) *Node[T] {
	return &Node[T]{bounds: bounds, items: make([]Item[T], 0), parent: parent}
}

func (n *Node[T]) isLeaf() bool { return len(n.childs) == 0 }
func (n *Node[T]) isNode() bool { return len(n.childs) > 0 }

func (n *Node[T]) findFittingChild(r geom.AABB[T]) *Node[T] {
	for _, child := range n.childs {
		if child.bounds.Contains(r) {
			return child
		}
	}
	return nil
}

func (n *Node[T]) Children() []*Node[T] {
	return n.childs
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
	items := []Item[T]{}

	dfs.DFS(n, struct{}{}, func(node *Node[T], _ struct{}) (dfs.DFSControl, struct{}) {
		items = append(items, node.items...)
		return dfs.DFSControl{}, struct{}{}
	})

	sortItems(items)
	return items
}

func (n *Node[T]) depth() int {
	maxDepth := 0

	dfs.DFS(n, 0, func(node *Node[T], acc int) (dfs.DFSControl, int) {
		acc = acc + 1
		if acc > maxDepth {
			maxDepth = acc
		}
		return dfs.DFSControl{}, acc
	})

	return maxDepth
}

func (n *Node[T]) leafBounds() []geom.AABB[T] {
	rectangles := []geom.AABB[T]{}

	dfs.DFS(n, struct{}{}, func(node *Node[T], _ struct{}) (dfs.DFSControl, struct{}) {
		if node.isLeaf() {
			rectangles = append(rectangles, node.bounds)
			return dfs.DFSControl{Skip: true}, struct{}{}
		}
		return dfs.DFSControl{}, struct{}{}
	})

	return rectangles
}
