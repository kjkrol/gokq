# GOKQ

*aka "Golang kjkrol QuadTree"*

The library provides a [Quadtree](https://en.wikipedia.org/wiki/Quadtree) data structure  
for efficient spatial indexing and 2D range queries.

**Powered by `gokg` planes:** the tree understands both bounded and cyclic planes
defined in [`gokg`](../gokg). Queries grazing a cyclic edge automatically wrap and
fragment their search regions, so neighbour lookups stay accurate even when the area
straddles the boundary.

## Usage example

```go
package main

import (
	"fmt"

	"github.com/kjkrol/gokg/pkg/geometry"
	"github.com/kjkrol/gokq/pkg/quadtree"
)

type point struct {
	box geometry.AABB[float64]
	id  string
}

func newPoint(id string, x, y float64) *point {
	pos := geometry.NewVec(x, y)
	return &point{
		box: pos.Bounds(),
		id:  id,
	}
}

func (p *point) Bound() geometry.AABB[float64] { return p.box }

func main() {
	plane := geometry.NewCyclicBoundedPlane[float64](64, 64)
	tree := quadtree.NewQuadTree(plane)
	defer tree.Close()

	for _, pt := range []*point{
		newPoint("edge", 63.0, 63.0),
		newPoint("wrap-x", 1.0, 63.0),
		newPoint("wrap-y", 63.0, 1.0),
	} {
		tree.Add(pt)
	}

	target := newPoint("target", 63.5, 63.5)
	for _, neighbor := range tree.FindNeighbors(target, 2.0) {
		fmt.Println(neighbor.Bound())
	}
}
```

For more scenarios, explore the example-based tests in `pkg/quadtree`, which double as runnable docs.

----
*[Contributor Recommendations](docs/Contributor_Recommendations.md)
