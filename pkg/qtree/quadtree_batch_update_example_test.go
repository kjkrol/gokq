package qtree

import (
	"fmt"

	"github.com/kjkrol/gokg/pkg/geometry"
)

func ExampleBatchUpdateCoordinator() {
	plane := geometry.NewBoundedPlane(16.0, 16.0)
	tree := NewQuadTree(plane)
	defer tree.Close()

	// Stationary items occupy two opposite quadrants.
	stationary := []*TestItem[float64]{
		newTestItemPointAtPos[float64](2, 2),   // top-left
		newTestItemPointAtPos[float64](14, 14), // bottom-right
	}
	for _, item := range stationary {
		tree.Add(item)
	}

	// Moving items fill the remaining quadrants so every child will be populated
	// once the split happens.
	moving := []*TestItem[float64]{
		newTestItemPointAtPos[float64](14, 2), // top-right
		newTestItemPointAtPos[float64](2, 14), // bottom-left
	}
	for _, item := range moving {
		tree.Add(item)
	}

	// 5th point triggers the split; it is placed inside the top-left quadrant.
	trigger := newTestItemPointAtPos(7.5, 7.5)
	tree.Add(trigger)

	fmt.Println("children before:", formatChildCounts(tree.root.childs))

	toRemove := []Item[float64]{moving[0], moving[1]}
	toAdd := []Item[float64]{
		newTestItemPointAtPos[float64](3, 3),   // follow top-left stationary
		newTestItemPointAtPos[float64](13, 13), // follow bottom-right stationary
	}

	tree.BatchUpdate(toRemove, toAdd, true)
	fmt.Println("children after :", formatChildCounts(tree.root.childs))

	fmt.Println("count:", tree.Count())
	all := tree.AllItems()
	positions := make([]string, len(all))
	for i, item := range all {
		positions[i] = item.Bound().String()
	}

	for _, pos := range positions {
		fmt.Println(pos)
	}

	// Output:
	// children before: [2 1 1 1]
	// children after : [3 0 0 2]
	// count: 5
	// {(2,2) (2,2)}
	// {(3,3) (3,3)}
	// {(7.5,7.5) (7.5,7.5)}
	// {(13,13) (13,13)}
	// {(14,14) (14,14)}
}
func formatChildCounts(children []*Node[float64]) []int {
	counts := make([]int, len(children))
	for i, child := range children {
		if child != nil {
			counts[i] = len(child.items)
		}
	}
	return counts
}
