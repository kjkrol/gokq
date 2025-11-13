// Package dfs provides a generic depth-first traversal helper that operates on
// any node type exposing its children through the ChildCarrier interface.
package dfs

// ChildCarrier exposes a node's children so DFS can traverse without
// knowing the concrete tree type.
type ChildCarrier[N any] interface {
	Children() []N
}

// DFSControl tweaks traversal: set Skip to avoid visiting a node's children,
// or Break to stop the DFS immediately.
type DFSControl struct {
	Skip  bool
	Break bool
}

// DFSStepFunc processes one node and returns traversal controls plus
// the accumulator value to forward to that node's descendants.
type DFSStepFunc[N ChildCarrier[N], A any] func(node N, acc A) (DFSControl, A)

// DFS walks the structure rooted at node using an explicit stack. accInitial is
// copied into the first frame, and each returned accumulator is stored on the
// branch that produced it, so sibling subtrees never share state unless A is a
// reference type.
func DFS[N ChildCarrier[N], A any](
	root N,
	accInitial A,
	step DFSStepFunc[N, A],
) {
	type frame struct {
		node N
		acc  A
	}
	stack := []frame{{root, accInitial}}
	for len(stack) > 0 {
		entry := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		control, nextAcc := step(entry.node, entry.acc)

		if control.Skip {
			continue
		}

		if control.Break {
			break
		}

		for _, child := range entry.node.Children() {
			stack = append(stack, frame{child, nextAcc})
		}
	}
}
