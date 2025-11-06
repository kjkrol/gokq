package quadtree

import "github.com/kjkrol/gokg/pkg/geometry"

type QuadTreeRemover[T geometry.SupportedNumeric, K comparable] struct {
	capacity int
}

func (qr QuadTreeRemover[T, K]) remove(node *Node[T, K], item Item[T, K]) bool {
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

func (qr QuadTreeRemover[T, K]) tryCompress(node *Node[T, K]) {
	if !node.isNode() {
		return
	}

	collected := qr.collectItems(node)
	if len(collected) <= qr.capacity {
		node.items = collected
		node.childs = nil
	}
}

func (qr QuadTreeRemover[T, K]) collectItems(n *Node[T, K]) []Item[T, K] {
	items := append([]Item[T, K]{}, n.items...)
	for _, ch := range n.childs {
		items = append(items, qr.collectItems(ch)...)
	}
	return items
}
