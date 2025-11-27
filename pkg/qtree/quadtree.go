package qtree

import (
	"github.com/kjkrol/gokg/pkg/geometry"
	"github.com/kjkrol/gokq/pkg/dfs"
)

const (
	CAPACITY  int = 4
	MAX_DEPTH int = 10
)

// Item represents an object with an axis-aligned bounding box and a stable identifier.
// SameID must return true when both operands refer to the same logical entity.
type Item[T geometry.SupportedNumeric] interface {
	Bound() geometry.BoundingBox[T]
	SameID(other Item[T]) bool
}

// QuadTree stores spatial items in a hierarchical grid for fast range queries.
type QuadTree[T geometry.SupportedNumeric] struct {
	root        *Node[T]
	appender    QuadTreeAppender[T]
	remover     QuadTreeRemover[T]
	finder      QuadTreeFinder[T]
	coordinator BatchUpdateCoordinator[T]
}

// NewQuadTree builds a QuadTree covering the supplied plane viewport.
func NewQuadTree[T geometry.SupportedNumeric](
	plane geometry.Plane[T],
	opts ...QuadTreeOption[T],
) *QuadTree[T] {
	rootBounds := plane.Viewport()
	root := newNode(rootBounds, nil)
	finderStrategy := NewDefaultQuadTreeFinderStrategy(plane)
	qt := &QuadTree[T]{
		root:     root,
		appender: QuadTreeAppender[T]{maxDepth: MAX_DEPTH, capacity: CAPACITY},
		remover:  QuadTreeRemover[T]{capacity: CAPACITY},
		finder:   NewQuadTreeFinder(finderStrategy),
	}
	qt.coordinator = NewBatchUpdateCoordinator(qt.appender, qt.remover)
	for _, opt := range opts {
		opt(qt)
	}
	return qt
}

// Add inserts item into the tree; returns false if it cannot be placed.
func (t *QuadTree[T]) Add(item Item[T]) bool {
	return t.appender.add(t.root, item, 0)
}

// Remove deletes item from the tree; returns false when nothing was removed.
func (t *QuadTree[T]) Remove(item Item[T]) bool {
	return t.remover.remove(t.root, item)
}

// Close releases internal resources held by the tree.
func (t *QuadTree[T]) Close() {
	t.root.close()
	t.coordinator.Close()
}

// Count returns the number of items stored in the tree.
func (t *QuadTree[T]) Count() int {
	total := 0

	dfs.DFS(t.root, struct{}{}, func(node *Node[T], _ struct{}) (dfs.DFSControl, struct{}) {
		total += len(node.items)
		return dfs.DFSControl{}, struct{}{}
	})

	return total
}

// Depth reports the maximum depth for active nodes.
func (t *QuadTree[T]) Depth() int {
	return t.root.depth()
}

// AllItems returns a snapshot of every stored item.
func (t *QuadTree[T]) AllItems() []Item[T] {
	return t.root.allItems()
}

// LeafBounds returns the bounding boxes of all current leaf nodes.
func (t *QuadTree[T]) LeafBounds() []geometry.BoundingBox[T] {
	return t.root.leafBounds()
}

// FindNeighbors retrieves items within margin of the target's bounds.
func (t *QuadTree[T]) FindNeighbors(target Item[T], margin T) []Item[T] {
	return t.finder.FindNeighbors(t.root, target, margin)
}

// BatchUpdate removes a batch of items, re-inserts the supplied replacements and
// optionally compresses the affected nodes. Compression is triggered when
// triggerCompression is true or when the number of touched nodes exceeds the
// configured threshold.
func (t *QuadTree[T]) BatchUpdate(toRemove []Item[T], toAdd []Item[T], triggerCompression bool) {
	if t == nil || t.root == nil {
		return
	}

	t.coordinator.BatchUpdate(t.root, toRemove, toAdd, triggerCompression)
}
