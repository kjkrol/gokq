package quadtree

import (
	"github.com/kjkrol/gokg/pkg/geometry"
	"github.com/kjkrol/gokq/pkg/dfs"
)

const defaultBatchCompressThreshold = 128

// BatchUpdateCoordinator tracks nodes touched during bulk updates so compression
// can be deferred and executed once per batch instead of after every removal.
type BatchUpdateCoordinator[T geometry.SupportedNumeric] struct {
	touched   map[*Node[T]]struct{}
	threshold int
	QuadTreeAppender[T]
	QuadTreeRemover[T]
}

// NewBatchUpdateCoordinator creates a coordinator bound to the supplied tree.
func NewBatchUpdateCoordinator[T geometry.SupportedNumeric](
	appender QuadTreeAppender[T],
	remover QuadTreeRemover[T],
) BatchUpdateCoordinator[T] {
	return BatchUpdateCoordinator[T]{
		touched:          make(map[*Node[T]]struct{}),
		threshold:        defaultBatchCompressThreshold,
		QuadTreeAppender: appender,
		QuadTreeRemover:  remover,
	}
}

// BatchUpdate removes items, adds replacements via provided callback, and returns
// the number of items removed. Compression is deferred until Compress() is called.
func (c *BatchUpdateCoordinator[T]) BatchUpdate(
	root *Node[T],
	toRemove []Item[T],
	toAdd []Item[T],
	triggerCompression bool,
) int {

	removed := c.removeBatch(root, toRemove)

	for _, item := range toAdd {
		c.QuadTreeAppender.add(root, item, 0)
	}

	if triggerCompression || c.shouldCompress() {
		c.compress()
	}

	return removed
}

func (c *BatchUpdateCoordinator[T]) removeBatch(root *Node[T], items []Item[T]) int {
	if c == nil || root == nil || len(items) == 0 {
		return 0
	}

	set := newBatchRemovalSet(items)
	removed := 0

	dfs.DFS(root, struct{}{}, func(node *Node[T], _ struct{}) (dfs.DFSControl, struct{}) {
		if len(set.items) == 0 {
			return dfs.DFSControl{Break: true}, struct{}{}
		}
		if count := c.removeFromNode(node, set); count > 0 {
			removed += count
		}
		return dfs.DFSControl{}, struct{}{}
	})

	return removed
}

func (c *BatchUpdateCoordinator[T]) removeFromNode(node *Node[T], set *batchRemovalSet[T]) int {
	keep := node.items[:0]
	removed := 0

	for _, item := range node.items {
		if set.consume(item) {
			removed++
			continue
		}
		keep = append(keep, item)
	}

	if removed > 0 {
		node.items = keep
		c.track(node)
	}

	return removed
}

// compress runs compression for all touched nodes (and their ancestors) and
// clears the pending set so the coordinator can be reused for another batch.
func (c *BatchUpdateCoordinator[T]) compress() {

	for node := range c.touched {
		c.QuadTreeRemover.compressPath(node)
	}
	c.reset()
}

// reset clears the current batch state without compressing anything.
func (c *BatchUpdateCoordinator[T]) reset() {
	for node := range c.touched {
		delete(c.touched, node)
	}
}

// ShouldCompress reports whether the pending set crossed the threshold.
func (c *BatchUpdateCoordinator[T]) shouldCompress() bool {
	if c.threshold <= 0 {
		return false
	}
	return len(c.touched) >= c.threshold
}

// Close releases references held by the coordinator so it can be garbage
// collected promptly after use.
func (c *BatchUpdateCoordinator[T]) Close() {
	c.touched = nil
}

func (c *BatchUpdateCoordinator[T]) track(node *Node[T]) {
	if node == nil {
		return
	}
	if c.touched == nil {
		c.touched = make(map[*Node[T]]struct{})
	}
	c.touched[node] = struct{}{}
}

type batchRemovalSet[T geometry.SupportedNumeric] struct {
	items []Item[T]
}

func newBatchRemovalSet[T geometry.SupportedNumeric](items []Item[T]) *batchRemovalSet[T] {
	copied := make([]Item[T], len(items))
	copy(copied, items)
	return &batchRemovalSet[T]{items: copied}
}

func (s *batchRemovalSet[T]) consume(target Item[T]) bool {
	for i, item := range s.items {
		if sameItem(item, target) {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return true
		}
	}
	return false
}

func sameItem[T geometry.SupportedNumeric](a, b Item[T]) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a == b {
		return true
	}
	return a.SameID(b)
}
