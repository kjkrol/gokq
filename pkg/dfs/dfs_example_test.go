package dfs

import (
	"fmt"
)

type exampleNode struct {
	label    string
	children []*exampleNode
}

func (n *exampleNode) Children() []*exampleNode {
	return n.children
}

// ExampleDFS_controlFlow documents how Skip and Break influence traversal order
// and shows that the accumulator delivered to each branch is isolated.
func ExampleDFS_controlFlow() {
	afterBreak := &exampleNode{label: "after-break"}
	stop := &exampleNode{label: "stop"}
	skipChild := &exampleNode{label: "skip-child"}
	skip := &exampleNode{label: "skip", children: []*exampleNode{skipChild}}

	root := &exampleNode{
		label:    "root",
		children: []*exampleNode{afterBreak, stop, skip},
	}

	var visited []string

	DFS(root, nil, func(node *exampleNode, acc []string) (DFSControl, []string) {
		nextAcc := append(append([]string(nil), acc...), node.label)
		visited = append(visited, fmt.Sprintf("%s %v", node.label, nextAcc))

		switch node.label {
		case "skip":
			return DFSControl{Skip: true}, nextAcc
		case "stop":
			return DFSControl{Break: true}, nextAcc
		default:
			return DFSControl{}, nextAcc
		}
	})

	for _, entry := range visited {
		fmt.Println(entry)
	}

	// Output:
	// root [root]
	// skip [root skip]
	// stop [root stop]
}
