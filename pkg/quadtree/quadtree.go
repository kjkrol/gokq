// Package quadtree provides an implementation of a quadtree data structure,
// which is used to partition a two-dimensional space by recursively subdividing
// it into four quadrants or regions. This structure is useful for various
// spatial indexing applications, such as range searching, nearest neighbor
// searching, and collision detection in 2D space.
package quadtree

import (
	"github.com/kjkrol/gokg/pkg/geometry"
)

// QuadTree represents a quadtree data structure for spatial indexing.
type QuadTree[T geometry.SupportedNumeric] struct {
	root  *Node[T]
	plane geometry.Plane[T]
}

// NewQuadTree creates a new quadtree with the specified plane.
func NewQuadTree[T geometry.SupportedNumeric](plane geometry.Plane[T]) QuadTree[T] {
	box := newBox(geometry.Vec[T]{X: 0, Y: 0}, plane.Size())
	root := newNode[T](box, nil)
	return QuadTree[T]{root, plane}
}

// Add inserts an item into the quadtree.
func (t *QuadTree[T]) Add(item Item[T]) {
	node := t.root
	for node.isNode() {
		node = node.traverseToChild(item.Vector())
	}
	node.add(item)
}

// Close releases resources associated with the QuadTree.
func (t *QuadTree[T]) Close() {
	t.root.close()
}

// FindNeighbors returns a list of items that are within the specified distance of the target item.
func (t *QuadTree[T]) FindNeighbors(target Item[T], distance T) []Item[T] {
	probeBox := buildBox(target.Vector(), distance)
	neighborNodes := selectingNeighbors(probeBox, t)
	predicate := predicate(target, distance, t.plane.Metric)
	return scan(neighborNodes, predicate)
}

func selectingNeighbors[T geometry.SupportedNumeric](probeBox box[T], t *QuadTree[T]) []*Node[T] {
	neighborNodes := make([]*Node[T], 0)
	if t.plane.Name() == "cyclic" {
		wrappedBoxes := wrapBoxCyclic(probeBox, t.plane.Size(), t.plane.Contains)
		return findIntersectingNodesMultiple(t.root, wrappedBoxes, neighborNodes)
	} else {
		return findIntersectingNodes(t.root, probeBox, neighborNodes)
	}
}

// ----------------------------------------------------------------------------

func findIntersectingNodes[T geometry.SupportedNumeric](node *Node[T], box box[T], neighborNodes []*Node[T]) []*Node[T] {
	if node.isLeaf() {
		return append(neighborNodes, node)
	}
	for _, childNode := range node.childs {
		if childNode.box.intersects(box) {
			neighborNodes = findIntersectingNodes(childNode, box, neighborNodes)
		}
	}
	return neighborNodes
}

func findIntersectingNodesMultiple[T geometry.SupportedNumeric](node *Node[T], boxes []box[T], neighborNodes []*Node[T]) []*Node[T] {
	if node.isLeaf() {
		return append(neighborNodes, node)
	}
	for _, childNode := range node.childs {
		if childNode.box.intersectsAny(boxes) {
			neighborNodes = findIntersectingNodesMultiple(childNode, boxes, neighborNodes)
		}
	}
	return neighborNodes
}

func predicate[T geometry.SupportedNumeric](
	target Item[T],
	distance T,
	metric func(geometry.Vec[T], geometry.Vec[T]) T,
) func(item *Item[T]) bool {
	if metric == nil {
		panic("metric function is required")
	}
	return func(item *Item[T]) bool {
		if (*item).Vector() == target.Vector() {
			return false
		}
		a := target.Vector()
		b := (*item).Vector()
		return metric(a, b) <= distance
	}
}

func scan[T geometry.SupportedNumeric](neighborNodes []*Node[T], predicate func(*Item[T]) bool) []Item[T] {
	neighborItems := make([]Item[T], 0)
	for _, node := range neighborNodes {
		for _, neighborItem := range node.items {
			if predicate(&neighborItem) {
				neighborItems = append(neighborItems, neighborItem)
			}
		}
	}
	return neighborItems
}

//-----------------------------------------------------------------------------

// Item:
// - Represents an item that can be stored in the quadtree.
// - Methods:
//   - Vector() geometry.Vec[int]: Returns the vector representation of the item.
type Item[T geometry.SupportedNumeric] interface {
	Vector() geometry.Vec[T]
}

// Node:
// - Represents a node in the quadtree.
// - Fields:
//   - box: The bounding box of the node.
//   - items: The items contained in the node.
//   - parent: A pointer to the parent node.
//   - childs: The child nodes of the current node.
type Node[T geometry.SupportedNumeric] struct {
	box    box[T]
	items  []Item[T]
	parent *Node[T]
	childs []*Node[T]
}

func newNode[T geometry.SupportedNumeric](box box[T], parent *Node[T]) *Node[T] {
	items := make([]Item[T], 0)
	return &Node[T]{box: box, items: items, parent: parent}
}

func (n *Node[T]) isLeaf() bool { return len(n.childs) == 0 }

func (n *Node[T]) isNode() bool { return len(n.childs) > 0 }

func (n *Node[T]) add(item Item[T]) {
	n.items = append(n.items, item)
	if len(n.items) > 3 {
		n.createChilds()
		n.arrange()
		clear(n.items)
	}
}

func (n *Node[T]) createChilds() {
	var childBoxes [4]box[T] = n.box.split()
	n.childs = make([]*Node[T], 4)
	for i, box := range childBoxes {
		n.childs[i] = newNode(box, n)
	}
}

func (n *Node[T]) arrange() {
	for _, item := range n.items {
		n.traverseToChild(item.Vector()).add(item)
	}
}

// traverseToChild determines the appropriate child node to traverse to based on the given vector's coordinates.
// It calculates the child index by comparing the vector's X and Y coordinates with the center of the current node's bounding box.
// The child nodes are indexed as follows:
// |0|1|
// |2|3|
//
// Parameters:
//   - vector: A geometry.Vec[int] representing the coordinates to compare.
//
// Returns:
//   - *Node: The child node corresponding to the calculated index.
func (n *Node[T]) traverseToChild(vector geometry.Vec[T]) *Node[T] {
	childIdx := 0
	if vector.X > n.box.center.X {
		childIdx += 1
	}
	if vector.Y > n.box.center.Y {
		childIdx += 2
	}
	return n.childs[childIdx]
}

func (n *Node[T]) close() {
	for _, child := range n.childs {
		child.close()
	}
	n.items = nil
	n.childs = nil
	n.parent = nil
}
