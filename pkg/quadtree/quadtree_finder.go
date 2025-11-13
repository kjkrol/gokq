package quadtree

import (
	"sort"

	"github.com/kjkrol/gokg/pkg/geometry"
	"github.com/kjkrol/gokq/pkg/dfs"
)

type QuadTreeFinder[T geometry.SupportedNumeric, K comparable] struct {
	strategy QuadTreeFinderStrategy[T, K]
}

func NewQuadTreeFinder[T geometry.SupportedNumeric, K comparable](strategy QuadTreeFinderStrategy[T, K]) QuadTreeFinder[T, K] {
	return QuadTreeFinder[T, K]{strategy: strategy}
}

func (qf QuadTreeFinder[T, K]) FindNeighbors(root *Node[T, K], target Item[T, K], margin T) []Item[T, K] {

	nodeIntersectionDetection := qf.strategy.NodeIntersectionDetectionFactory(target, margin)
	itemsInRangeDetection := qf.strategy.ItemsInRangeDetectionFactory(target, margin)
	neighbors := make([]Item[T, K], 0)

	dfs.DFS(root, struct{}{}, func(node *Node[T, K], _ struct{}) (dfs.DFSControl, struct{}) {
		if !nodeIntersectionDetection(*node) {
			return dfs.DFSControl{Skip: true}, struct{}{}
		}
		itemsInRangeDetection(*node, func(item Item[T, K]) { neighbors = append(neighbors, item) })
		return dfs.DFSControl{}, struct{}{}
	})

	sortNeighbors(neighbors)
	return neighbors
}

func sortNeighbors[T geometry.SupportedNumeric, K comparable](items []Item[T, K]) {
	sort.Slice(items, func(i, j int) bool {
		ai, aj := items[i].Bound(), items[j].Bound()
		first, _ := geometry.SortBoxesBy(
			ai, aj,
			func(box geometry.BoundingBox[T]) T { return box.TopLeft.Y },
			func(box geometry.BoundingBox[T]) T { return box.TopLeft.X },
			func(box geometry.BoundingBox[T]) T { return box.BottomRight.Y },
			func(box geometry.BoundingBox[T]) T { return box.BottomRight.X },
		)
		return first.Equals(ai)
	})
}
