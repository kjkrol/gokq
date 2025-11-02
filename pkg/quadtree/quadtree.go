package quadtree

import (
	"github.com/kjkrol/gokg/pkg/geometry"
)

const (
	CAPACITY  int = 4
	MAX_DEPTH int = 10
)

type Item[T geometry.SupportedNumeric] interface {
	AABB() geometry.AABB[T]
}

type QuadTree[T geometry.SupportedNumeric] struct {
	root     *Node[T]
	appender QuadTreeAppender[T]
	remover  QuadTreeRemover[T]
	finder   QuadTreeFinder[T]
}

func NewQuadTree[T geometry.SupportedNumeric](
	plane geometry.Plane[T],
	opts ...QuadTreeOption[T],
) *QuadTree[T] {
	rootBounds := geometry.NewAABB(geometry.Vec[T]{X: 0, Y: 0}, plane.Size())
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

func (t *QuadTree[T]) Add(item Item[T]) bool {
	return t.appender.add(t.root, item, 0)
}

func (t *QuadTree[T]) Remove(item Item[T]) bool {
	return t.remover.remove(t.root, item)
}

func (t *QuadTree[T]) Close() {
	t.root.close()
}

func (t *QuadTree[T]) Count() int {
	return len(t.root.allItems())
}

func (t *QuadTree[T]) Depth() int {
	return t.root.depth()
}

func (t *QuadTree[T]) AllItems() []Item[T] {
	return t.root.allItems()
}

func (t *QuadTree[T]) LeafRectangles() []geometry.AABB[T] {
	return t.root.leafRectangles()
}

func (t *QuadTree[T]) FindNeighbors(target Item[T], margin T) []Item[T] {
	return t.finder.FindNeighbors(t.root, target, margin)
}
