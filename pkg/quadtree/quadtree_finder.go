package quadtree

import (
	"github.com/kjkrol/gokg/pkg/geometry"
	"github.com/kjkrol/gokq/pkg/dfs"
)

type QuadTreeFinder[T geometry.SupportedNumeric] struct {
	strategy QuadTreeFinderStrategy[T]
}

func NewQuadTreeFinder[T geometry.SupportedNumeric](strategy QuadTreeFinderStrategy[T]) QuadTreeFinder[T] {
	return QuadTreeFinder[T]{strategy: strategy}
}

func (qf QuadTreeFinder[T]) FindNeighbors(root *Node[T], target Item[T], margin T) []Item[T] {

	nodeIntersectionDetection := qf.strategy.NodeIntersectionDetectionFactory(target, margin)
	itemsInRangeDetection := qf.strategy.ItemsInRangeDetectionFactory(target, margin)
	neighbors := make([]Item[T], 0)

	dfs.DFS(root, struct{}{}, func(node *Node[T], _ struct{}) (dfs.DFSControl, struct{}) {
		if !nodeIntersectionDetection(*node) {
			return dfs.DFSControl{Skip: true}, struct{}{}
		}
		itemsInRangeDetection(*node, func(item Item[T]) { neighbors = append(neighbors, item) })
		return dfs.DFSControl{}, struct{}{}
	})

	sortItems(neighbors)
	return neighbors
}
