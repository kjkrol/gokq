package quadtree

import (
	"github.com/kjkrol/gokg/pkg/geometry"
)

const (
	CAPACITY  int = 4
	MAX_DEPTH int = 10
)

// Item represents an object with an axis-aligned bounding box.
type Item[T geometry.SupportedNumeric] interface {
	Bound() geometry.AABB[T]
}

// QuadTree stores spatial items in a hierarchical grid for fast range queries.
type QuadTree[T geometry.SupportedNumeric] struct {
	root     *Node[T]
	appender QuadTreeAppender[T]
	remover  QuadTreeRemover[T]
	finder   QuadTreeFinder[T]
}

// NewQuadTree builds a QuadTree covering the supplied plane viewport.
func NewQuadTree[T geometry.SupportedNumeric](
	plane geometry.Plane[T],
	opts ...QuadTreeOption[T],
) *QuadTree[T] {
	rootBounds := plane.Viewport()
	root := newNode(rootBounds, nil)
	qt := &QuadTree[T]{
		root:     root,
		appender: QuadTreeAppender[T]{maxDepth: MAX_DEPTH},
		remover:  QuadTreeRemover[T]{capacity: CAPACITY},
		finder:   NewQuadTreeFinder(plane),
	}
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
}

// Count returns the number of items stored in the tree.
func (t *QuadTree[T]) Count() int {
	return len(t.root.allItems())
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
func (t *QuadTree[T]) LeafBounds() []geometry.AABB[T] {
	return t.root.leafRectangles()
}

// FindNeighbors retrieves items within margin of the target's bounds.
func (t *QuadTree[T]) FindNeighbors(target Item[T], margin T) []Item[T] {
	return t.finder.FindNeighbors(t.root, target, margin)
}
