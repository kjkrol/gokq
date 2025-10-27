package quadtree

import (
	"fmt"

	"github.com/kjkrol/gokg/pkg/geometry"
	"github.com/kjkrol/gokg/pkg/geometry/spatial"
)

func newRectangleSpatial[T spatial.SupportedNumeric](topLeft, bottomRight spatial.Vec[T]) spatial.Spatial[T] {
	rect := spatial.NewRectangle(topLeft, bottomRight)
	return &rect
}

func ExampleQuadTree() {
	// Tworzymy płaszczyznę i QuadTree
	boundedPlane := geometry.NewBoundedPlane(64, 64)
	qtree := NewQuadTree(boundedPlane)
	defer qtree.Close()

	// Dodajemy kilka obiektów jako boxy 1x1
	items := []*ExampleItem[int]{
		{spatial: newRectangleSpatial(spatial.Vec[int]{X: 32, Y: 33}, spatial.Vec[int]{X: 10, Y: 5})},
		{spatial: newRectangleSpatial(spatial.Vec[int]{X: 31, Y: 32}, spatial.Vec[int]{X: 11, Y: 12})},
		{spatial: newRectangleSpatial(spatial.Vec[int]{X: 33, Y: 32}, spatial.Vec[int]{X: 40, Y: 13})},
		{spatial: newRectangleSpatial(spatial.Vec[int]{X: 32, Y: 31}, spatial.Vec[int]{X: 33, Y: 32})}, // sąsiad (góra)
		{spatial: newRectangleSpatial(spatial.Vec[int]{X: 32, Y: 33}, spatial.Vec[int]{X: 33, Y: 34})}, // sąsiad (dół)
		{spatial: newRectangleSpatial(spatial.Vec[int]{X: 31, Y: 32}, spatial.Vec[int]{X: 32, Y: 33})}, // sąsiad (lewo)
		{spatial: newRectangleSpatial(spatial.Vec[int]{X: 33, Y: 32}, spatial.Vec[int]{X: 34, Y: 33})}, // sąsiad (prawo)
	}
	for _, item := range items {
		qtree.Add(item)
	}

	// Wybieramy target
	target := &ExampleItem[int]{spatial: newRectangleSpatial(spatial.Vec[int]{X: 32, Y: 32}, spatial.Vec[int]{X: 33, Y: 33})}

	// Szukamy sąsiadów targeta z marginesem 0 (czyli boxy przecinające się dokładnie z nim)
	neighbors := qtree.FindNeighbors(target, 0)

	// Wypisujemy sąsiadów (ich granice)
	for _, neighbor := range neighbors {
		fmt.Println(neighbor.Value())
	}

	// Output:
	// {(32,31) (33,32) (32,31)}
	// {(31,32) (32,33) (31,32)}
	// {(33,32) (34,33) (33,32)}
	// {(32,33) (33,34) (32,33)}

}

// ExampleQuadTree_largeBoxes pokazuje działanie szukania sąsiadów.
// Wizualizacja tego przykładu znajduje się w docs/ExampleQuadTree_largeBoxes.svg
func ExampleQuadTree_largeBoxes() {
	// Tworzymy płaszczyznę i QuadTree
	boundedPlane := geometry.NewBoundedPlane(64.0, 64.0)
	qtree := NewQuadTree(boundedPlane)
	defer qtree.Close()

	// Dodajemy boxy 2x2
	items := []*ExampleItem[float64]{
		{spatial: newRectangleSpatial(spatial.Vec[float64]{X: 9, Y: 7}, spatial.Vec[float64]{X: 11, Y: 9})},   // powyżej w odległości 0
		{spatial: newRectangleSpatial(spatial.Vec[float64]{X: 9, Y: 11}, spatial.Vec[float64]{X: 11, Y: 13})}, // poniżej w odległości 0
		{spatial: newRectangleSpatial(spatial.Vec[float64]{X: 7, Y: 9}, spatial.Vec[float64]{X: 9, Y: 11})},   // z lewej w odległości 0
		{spatial: newRectangleSpatial(spatial.Vec[float64]{X: 11, Y: 9}, spatial.Vec[float64]{X: 13, Y: 11})}, // z prawej w odległości 0
		{spatial: newRectangleSpatial(spatial.Vec[float64]{X: 9, Y: 4}, spatial.Vec[float64]{X: 11, Y: 6})},   // powyżej w odległości 3
		{spatial: newRectangleSpatial(spatial.Vec[float64]{X: 9, Y: 14}, spatial.Vec[float64]{X: 11, Y: 16})}, // poniżej w odległości 3
		{spatial: newRectangleSpatial(spatial.Vec[float64]{X: 4, Y: 9}, spatial.Vec[float64]{X: 6, Y: 11})},   // z lewej w odległości 3
		{spatial: newRectangleSpatial(spatial.Vec[float64]{X: 14, Y: 9}, spatial.Vec[float64]{X: 16, Y: 11})}, // z prawej w odległości 3
		{spatial: newRectangleSpatial(spatial.Vec[float64]{X: 9, Y: 5}, spatial.Vec[float64]{X: 11, Y: 7})},   // powyzej w odleglosci 2
		{spatial: newRectangleSpatial(spatial.Vec[float64]{X: 6, Y: 6}, spatial.Vec[float64]{X: 8, Y: 8})},    // powyzej w odleglosci 2
	}
	for _, item := range items {
		qtree.Add(item)
	}

	// Target to pierwszy box
	target := &ExampleItem[float64]{spatial: newRectangleSpatial(spatial.Vec[float64]{X: 9, Y: 9}, spatial.Vec[float64]{X: 11, Y: 11})}

	// Szukamy sąsiadów targeta z marginesem 1.5
	neighbors := qtree.FindNeighbors(target, 1.5)

	for _, neighbor := range neighbors {
		fmt.Println(neighbor.Value())
	}

	// Output:
	// {(6,6) (8,8) (7,7)}
	// {(9,7) (11,9) (10,8)}
	// {(7,9) (9,11) (8,10)}
	// {(11,9) (13,11) (12,10)}
	// {(9,11) (11,13) (10,12)}

}
