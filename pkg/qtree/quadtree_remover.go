package qtree

import (
	"github.com/kjkrol/gokg/pkg/geom"
	"github.com/kjkrol/gokq/pkg/dfs"
)

type QuadTreeRemover[T geom.Numeric] struct {
	capacity int
}

func (qr QuadTreeRemover[T]) remove(node *Node[T], item Item[T]) bool {
	if removedNode, ok := qr.removeInternal(node, item); ok {
		qr.compressPath(removedNode)
		return true
	}
	return false
}

func (qr QuadTreeRemover[T]) RemoveWithoutCompression(node *Node[T], item Item[T]) (*Node[T], bool) {
	return qr.removeInternal(node, item)
}

func (qr QuadTreeRemover[T]) removeInternal(node *Node[T], item Item[T]) (*Node[T], bool) {
	if node.isNode() {
		if child := node.findFittingChild(item.Bound()); child != nil {
			if removedNode, ok := qr.removeInternal(child, item); ok {
				return removedNode, true
			}
		}
	}

	for i, it := range node.items {
		if it == item {
			node.items = append(node.items[:i], node.items[i+1:]...)
			return node, true
		}
	}

	return nil, false
}

func (qr QuadTreeRemover[T]) compressPath(node *Node[T]) {
	for n := node; n != nil; n = n.parent {
		qr.tryCompress(n)
	}
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
	items := make([]Item[T], 0, len(n.items))
	dfs.DFS(n, struct{}{}, func(node *Node[T], _ struct{}) (dfs.DFSControl, struct{}) {
		items = append(items, node.items...)
		return dfs.DFSControl{}, struct{}{}
	})
	return items
}
