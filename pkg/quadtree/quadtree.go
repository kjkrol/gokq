// Package quadtree provides an implementation of a quadtree data structure,
// which is used to partition a two-dimensional space by recursively subdividing
// it into four quadrants or regions. This structure is useful for various
// spatial indexing applications, such as range searching, nearest neighbor
// searching, and collision detection in 2D space.
package quadtree

import (
	"github.com/kjkrol/gokg/pkg/geometry"
)

// QuadTree:
// - Represents the quadtree data structure.
// - Fields:
//   - root: A pointer to the root node of the quadtree.
//
// - Methods:
//   - NewQuadTree(box Box) QuadTree: Creates a new quadtree with the given bounding box.
//   - Add(item Item): Adds an item to the quadtree.
//   - FindNeighbors(target Item, distance int, metric func(geometry.Vec[int], geometry.Vec[int]) int) []Item:
//     Finds and returns the neighbors of the target item within the specified distance using the given metric.
type QuadTree struct {
	root *Node
}

func NewQuadTree(box Box) QuadTree {
	root := newNode(box, nil)
	return QuadTree{root}
}

// Add inserts an item into the QuadTree. It starts at the root node and
// traverses deeper into the tree until it finds a leaf node where the item
// can be added. The item's position is determined by its vector.
func (t *QuadTree) Add(item Item) {
	node := t.root
	for node.isNode() {
		node = node.goDeeper(item.Vector())
	}
	node.add(item)
}

// FindNeighbors searches for neighboring items within a specified distance from the target item.
// It uses a custom metric function to determine the distance between items.
//
// Parameters:
//   - target: The item for which neighbors are to be found.
//   - distance: The maximum distance within which to search for neighbors.
//   - metric: A function that calculates the distance between two vectors.
//
// Returns:
//
//	A slice of items that are neighbors of the target item within the specified distance.
func (t *QuadTree) FindNeighbors(
	target Item,
	distance int,
	metric func(geometry.Vec[int], geometry.Vec[int]) int,
) []Item {
	neighborNodes := make([]*Node, 0)
	probeBox := buildBox(target.Vector(), distance)
	neighborNodes = t.root.inspect(probeBox, neighborNodes)
	predicate := predicate(target, distance, metric)
	return scan(neighborNodes, predicate)
}

func (n *Node) inspect(box Box, neighborNodes []*Node) []*Node {
	if n.isLeaf() {
		return append(neighborNodes, n)
	}
	for _, child := range n.childs {
		if child.intersects(&box) {
			neighborNodes = child.inspect(box, neighborNodes)
		}
	}
	return neighborNodes
}

func predicate(
	target Item,
	distance int,
	metric func(geometry.Vec[int], geometry.Vec[int]) int,
) func(item *Item) bool {
	if metric == nil {
		metric = defaultMetric
	}
	return func(item *Item) bool {
		a := target.Vector()
		b := (*item).Vector()
		return metric(a, b) <= distance
	}
}

func defaultMetric(a, b geometry.Vec[int]) int {
	dx := a.X - b.X
	if dx < 0 {
		dx *= -1
	}
	dy := a.Y - b.Y
	if dy < 0 {
		dy *= -1
	}
	return dx + dy
}

func scan(neighborNodes []*Node, predicate func(*Item) bool) []Item {
	neighborItems := make([]Item, 0)
	for _, node := range neighborNodes {
		for _, neighborItem := range node.items {
			if predicate(&neighborItem) {
				neighborItems = append(neighborItems, neighborItem)
			}
		}
	}
	return neighborItems
}

// Item:
// - Represents an item that can be stored in the quadtree.
// - Methods:
//   - Vector() geometry.Vec[int]: Returns the vector representation of the item.
type Item interface {
	Vector() geometry.Vec[int]
}

// Node:
// - Represents a node in the quadtree.
// - Fields:
//   - box: The bounding box of the node.
//   - items: The items contained in the node.
//   - parent: A pointer to the parent node.
//   - childs: The child nodes of the current node.
type Node struct {
	box    Box
	items  []Item
	parent *Node
	childs []*Node
}

func newNode(box Box, parent *Node) *Node {
	items := make([]Item, 0)
	return &Node{box: box, items: items, parent: parent}
}

func (n *Node) isLeaf() bool { return len(n.childs) == 0 && len(n.items) > 0 }

func (n *Node) isNode() bool { return len(n.childs) > 0 }

func (n *Node) add(item Item) {
	n.items = append(n.items, item)
	if len(n.items) > 3 {
		n.createChilds()
		n.arrange()
		clear(n.items)
	}
}

func (n *Node) createChilds() {
	var childBoxes [4]Box = n.box.split()
	n.childs = make([]*Node, 4)
	for i, box := range childBoxes {
		n.childs[i] = newNode(box, n)
	}
}

func (n *Node) arrange() {
	for _, item := range n.items {
		n.goDeeper(item.Vector()).add(item)
	}
}

// goDeeper determines the appropriate child node to traverse to based on the given vector's coordinates.
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
func (n *Node) goDeeper(vector geometry.Vec[int]) *Node {
	childIdx := 0
	if vector.Y > n.box.center.Y {
		childIdx += 2
	}
	if vector.X > n.box.center.X {
		childIdx += 1
	}
	return n.childs[childIdx]
}

func (n *Node) intersects(box *Box) bool {
	return n.box.intersects(box)
}
