package quadtree

import (
	"github.com/kjkrol/gokg/pkg/geometry"
)

const (
	CAPACITY  int = 4
	MAX_DEPTH int = 10
)

// Item represents an object with an axis-aligned bounding box and a stable identifier.
// The ID is consumed by the Finder to compare and deduplicate quadtree nodes.
type Item[T geometry.SupportedNumeric, K comparable] interface {
	Bound() geometry.BoundingBox[T]
	Id() K
}

// QuadTree stores spatial items in a hierarchical grid for fast range queries.
type QuadTree[T geometry.SupportedNumeric, K comparable] struct {
	root     *Node[T, K]
	appender QuadTreeAppender[T, K]
	remover  QuadTreeRemover[T, K]
	finder   QuadTreeFinder[T, K]
}

// NewQuadTree builds a QuadTree covering the supplied plane viewport.
func NewQuadTree[T geometry.SupportedNumeric, K comparable](
	plane geometry.Plane[T],
	opts ...QuadTreeOption[T, K],
) *QuadTree[T, K] {
	rootBounds := plane.Viewport().BoundingBox
	root := newNode[T, K](rootBounds, nil)
	qt := &QuadTree[T, K]{
		root:     root,
		appender: QuadTreeAppender[T, K]{maxDepth: MAX_DEPTH},
		remover:  QuadTreeRemover[T, K]{capacity: CAPACITY},
		finder:   NewQuadTreeFinder[T, K](plane),
	}
	for _, opt := range opts {
		opt(qt)
	}
	return qt
}

// Add inserts item into the tree; returns false if it cannot be placed.
func (t *QuadTree[T, K]) Add(item Item[T, K]) bool {
	return t.appender.add(t.root, item, 0)
}

// Remove deletes item from the tree; returns false when nothing was removed.
func (t *QuadTree[T, K]) Remove(item Item[T, K]) bool {
	return t.remover.remove(t.root, item)
}

// Close releases internal resources held by the tree.
func (t *QuadTree[T, K]) Close() {
	t.root.close()
}

// Count returns the number of items stored in the tree.
func (t *QuadTree[T, K]) Count() int {
	return len(t.root.allItems())
}

// Depth reports the maximum depth for active nodes.
func (t *QuadTree[T, K]) Depth() int {
	return t.root.depth()
}

// AllItems returns a snapshot of every stored item.
func (t *QuadTree[T, K]) AllItems() []Item[T, K] {
	return t.root.allItems()
}

// LeafBounds returns the bounding boxes of all current leaf nodes.
func (t *QuadTree[T, K]) LeafBounds() []geometry.BoundingBox[T] {
	return t.root.leafRectangles()
}

// FindNeighbors retrieves items within margin of the target's bounds.
func (t *QuadTree[T, K]) FindNeighbors(target Item[T, K], margin T) []Item[T, K] {
	return t.finder.FindNeighbors(t.root, target, margin)
}
