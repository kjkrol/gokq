package quadtree

import "github.com/kjkrol/gokg/pkg/geometry"

type QuadTreeRemover[T geometry.SupportedNumeric] struct {
	capacity int
}

func (qr QuadTreeRemover[T]) remove(node *Node[T], item Item[T]) bool {
	if node.isNode() {
		if child := node.findFittingChild(item.Bound()); child != nil {
			if qr.remove(child, item) {
				qr.tryCompress(node)
				return true
			}
		}
	}

	for i, it := range node.items {
		if it == item {
			node.items = append(node.items[:i], node.items[i+1:]...)
			qr.tryCompress(node)
			return true
		}
	}

	return false
}

func (qr QuadTreeRemover[T]) tryCompress(node *Node[T]) {
	if !node.isNode() {
		return
	}

	collected := qr.collectItems(node)
	if len(collected) <= qr.capacity {
		node.items = collected
		node.childs = nil
	}
}

func (qr QuadTreeRemover[T]) collectItems(n *Node[T]) []Item[T] {
	items := append([]Item[T]{}, n.items...)
	for _, ch := range n.childs {
		items = append(items, qr.collectItems(ch)...)
	}
	return items
}
