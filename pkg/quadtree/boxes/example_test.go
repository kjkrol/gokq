package boxes

import (
	"fmt"

	quadcore "github.com/kjkrol/goka/pkg/quadtree/base"
	"github.com/kjkrol/gokg/pkg/geometry"
)

// ExampleBoxItem implementuje Item dla QuadTree boxowego.
type ExampleBoxItem struct {
	box quadcore.Box[int]
}

func (ei *ExampleBoxItem) Bounds() quadcore.Box[int] {
	return ei.box
}

func ExampleQuadTree() {
	// Tworzymy płaszczyznę i QuadTree
	boundedPlane := geometry.NewBoundedPlane(64, 64)
	qtree := NewQuadTree(boundedPlane)
	defer qtree.Close()

	// Dodajemy kilka obiektów jako boxy 1x1
	items := []*ExampleBoxItem{
		{box: quadcore.NewBox(geometry.Vec[int]{X: 32, Y: 33}, geometry.Vec[int]{X: 10, Y: 5})},
		{box: quadcore.NewBox(geometry.Vec[int]{X: 31, Y: 32}, geometry.Vec[int]{X: 11, Y: 12})},
		{box: quadcore.NewBox(geometry.Vec[int]{X: 33, Y: 32}, geometry.Vec[int]{X: 40, Y: 13})},
		{box: quadcore.NewBox(geometry.Vec[int]{X: 32, Y: 31}, geometry.Vec[int]{X: 33, Y: 32})}, // sąsiad (góra)
		{box: quadcore.NewBox(geometry.Vec[int]{X: 32, Y: 33}, geometry.Vec[int]{X: 33, Y: 34})}, // sąsiad (dół)
		{box: quadcore.NewBox(geometry.Vec[int]{X: 31, Y: 32}, geometry.Vec[int]{X: 32, Y: 33})}, // sąsiad (lewo)
		{box: quadcore.NewBox(geometry.Vec[int]{X: 33, Y: 32}, geometry.Vec[int]{X: 34, Y: 33})}, // sąsiad (prawo)
	}
	for _, item := range items {
		qtree.Add(item)
	}

	// Wybieramy target
	target := &ExampleBoxItem{box: quadcore.NewBox(geometry.Vec[int]{X: 32, Y: 32}, geometry.Vec[int]{X: 33, Y: 33})}

	// Szukamy sąsiadów targeta z marginesem 0 (czyli boxy przecinające się dokładnie z nim)
	neighbors := qtree.FindNeighbors(target, 0)

	// Wypisujemy sąsiadów (ich granice)
	for _, neighbor := range neighbors {
		fmt.Println(neighbor.Bounds())
	}

	// Output:
	// {(32,31) (33,32) (32,31)}
	// {(31,32) (32,33) (31,32)}
	// {(32,33) (33,34) (32,33)}
	// {(33,32) (34,33) (33,32)}

}

// ExampleQuadTree_largeBoxes pokazuje działanie szukania sąsiadów.
// Wizualizacja tego przykładu znajduje się w docs/ExampleQuadTree_largeBoxes.png
func ExampleQuadTree_largeBoxes() {
	// Tworzymy płaszczyznę i QuadTree
	boundedPlane := geometry.NewBoundedPlane(64, 64)
	qtree := NewQuadTree(boundedPlane)
	defer qtree.Close()

	// Dodajemy boxy 2x2
	items := []*ExampleBoxItem{
		{box: quadcore.NewBox(geometry.Vec[int]{X: 9, Y: 7}, geometry.Vec[int]{X: 11, Y: 9})},   // powyżej w odległości 0
		{box: quadcore.NewBox(geometry.Vec[int]{X: 9, Y: 11}, geometry.Vec[int]{X: 11, Y: 13})}, // poniżej w odległości 0
		{box: quadcore.NewBox(geometry.Vec[int]{X: 7, Y: 9}, geometry.Vec[int]{X: 9, Y: 11})},   // z lewej w odległości 0
		{box: quadcore.NewBox(geometry.Vec[int]{X: 11, Y: 9}, geometry.Vec[int]{X: 13, Y: 11})}, // z prawej w odległości 0
		{box: quadcore.NewBox(geometry.Vec[int]{X: 9, Y: 4}, geometry.Vec[int]{X: 11, Y: 6})},   // powyżej w odległości 3
		{box: quadcore.NewBox(geometry.Vec[int]{X: 9, Y: 14}, geometry.Vec[int]{X: 11, Y: 16})}, // poniżej w odległości 3
		{box: quadcore.NewBox(geometry.Vec[int]{X: 4, Y: 9}, geometry.Vec[int]{X: 6, Y: 11})},   // z lewej w odległości 3
		{box: quadcore.NewBox(geometry.Vec[int]{X: 14, Y: 9}, geometry.Vec[int]{X: 16, Y: 11})}, // z prawej w odległości 3
		{box: quadcore.NewBox(geometry.Vec[int]{X: 9, Y: 5}, geometry.Vec[int]{X: 11, Y: 7})},   // powyzej w odleglosci 2
	}
	for _, item := range items {
		qtree.Add(item)
	}

	// Target to pierwszy box
	target := &ExampleBoxItem{box: quadcore.NewBox(geometry.Vec[int]{X: 9, Y: 9}, geometry.Vec[int]{X: 11, Y: 11})}

	// Szukamy sąsiadów targeta z marginesem 2
	neighbors := qtree.FindNeighbors(target, 2)

	for _, neighbor := range neighbors {
		fmt.Println(neighbor.Bounds())
	}

	// Output:
	// {(9,7) (11,9) (10,8)}
	// {(7,9) (9,11) (8,10)}
	// {(9,5) (11,7) (10,6)}
	// {(9,11) (11,13) (10,12)}
	// {(11,9) (13,11) (12,10)}
}
