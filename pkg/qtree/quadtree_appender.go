package qtree

import "github.com/kjkrol/gokg/pkg/geometry"

type QuadTreeAppender[T geometry.SupportedNumeric] struct {
	maxDepth int
	capacity int
}

func (qa QuadTreeAppender[T]) add(node *Node[T], item Item[T], depth int) bool {

	if !node.bounds.Contains(item.Bound()) {
		return false
	}

	if node.isNode() && depth < qa.maxDepth {
		if child := node.findFittingChild(item.Bound()); child != nil {
			if qa.add(child, item, depth+1) {
				return true
			}
		}
	}
	node.items = append(node.items, item)

	if len(node.items) > qa.capacity && node.isLeaf() && depth < qa.maxDepth {
		qa.createChilds(node)
		qa.redistributeItems(node, depth)
	}

	return true
}

func (qa QuadTreeAppender[T]) redistributeItems(node *Node[T], depth int) {
	remaining := make([]Item[T], 0, len(node.items))
	moved := 0

	for _, item := range node.items {
		if child := node.findFittingChild(item.Bound()); child != nil && qa.add(child, item, depth+1) {
			moved++
		} else {
			remaining = append(remaining, item)
		}
	}
	node.items = remaining

	if moved == 0 {
		for _, ch := range node.childs {
			node.items = append(node.items, ch.items...)
		}
		node.childs = nil
	}
}

func (qa QuadTreeAppender[T]) createChilds(node *Node[T]) {
	childRectangles := node.bounds.Split()
	node.childs = make([]*Node[T], qa.capacity)
	for i, rect := range childRectangles {
		node.childs[i] = newNode(rect, node)
	}
}
