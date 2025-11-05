package quadtree

import (
	"fmt"

	"github.com/kjkrol/gokg/pkg/geometry"
)

func ExampleQuadTree() {
	// Tworzymy płaszczyznę i QuadTree
	boundedPlane := geometry.NewBoundedPlane(64, 64)
	qtree := NewQuadTree(boundedPlane)
	defer qtree.Close()

	// Dodajemy kilka obiektów jako boxy 1x1
	items := []*ShapeItem[int]{
		{geometry.NewAABB(geometry.NewVec(32, 30), 1, 1)}, // wyzej
		{geometry.NewAABB(geometry.NewVec(31, 31), 1, 1)}, // lewy górny
		{geometry.NewAABB(geometry.NewVec(32, 31), 1, 1)}, // góra
		{geometry.NewAABB(geometry.NewVec(33, 31), 1, 1)}, // prawy górny
		{geometry.NewAABB(geometry.NewVec(31, 32), 1, 1)}, // lewo
		{geometry.NewAABB(geometry.NewVec(33, 32), 1, 1)}, // prawo
		{geometry.NewAABB(geometry.NewVec(31, 33), 1, 1)}, // lewy dolny
		{geometry.NewAABB(geometry.NewVec(32, 33), 1, 1)}, // dół
		{geometry.NewAABB(geometry.NewVec(33, 33), 1, 1)}, // prawy dolny
		{geometry.NewAABB(geometry.NewVec(32, 34), 1, 1)}, // nizej
	}
	for _, item := range items {
		qtree.Add(item)
	}

	// Wybieramy target
	target := &ShapeItem[int]{geometry.NewAABB(geometry.NewVec(32, 32), 1, 1)}

	// Szukamy sąsiadów targeta z marginesem 0 (czyli boxy przecinające się dokładnie z nim)
	neighbors := qtree.FindNeighbors(target, 0)

	// Wypisujemy sąsiadów (ich granice)
	for _, neighbor := range neighbors {
		fmt.Println(neighbor.Bound())
	}

	// Output:
	// {(31,31) (32,32)}
	// {(32,31) (33,32)}
	// {(33,31) (34,32)}
	// {(31,32) (32,33)}
	// {(33,32) (34,33)}
	// {(31,33) (32,34)}
	// {(32,33) (33,34)}
	// {(33,33) (34,34)}

}

// ExampleQuadTree_largeBoxes pokazuje działanie szukania sąsiadów.
// Wizualizacja tego przykładu znajduje się w docs/ExampleQuadTree_largeBoxes.svg
func ExampleQuadTree_largeBoxes() {
	// Tworzymy płaszczyznę i QuadTree
	boundedPlane := geometry.NewBoundedPlane(64.0, 64.0)
	qtree := NewQuadTree(boundedPlane)
	defer qtree.Close()

	// Dodajemy boxy 2x2
	items := []*ShapeItem[float64]{
		{geometry.NewAABB(geometry.NewVec(9.0, 7), 2.0, 2.0)},  // powyżej w odległości 0
		{geometry.NewAABB(geometry.NewVec(9.0, 11), 2.0, 2.0)}, // poniżej w odległości 0
		{geometry.NewAABB(geometry.NewVec(7.0, 9), 2.0, 2.0)},  // z lewej w odległości 0
		{geometry.NewAABB(geometry.NewVec(11.0, 9), 2.0, 2.0)}, // z prawej w odległości 0
		{geometry.NewAABB(geometry.NewVec(9.0, 4), 2.0, 2.0)},  // powyżej w odległości 3
		{geometry.NewAABB(geometry.NewVec(9.0, 14), 2.0, 2.0)}, // poniżej w odległości 3
		{geometry.NewAABB(geometry.NewVec(4.0, 9), 2.0, 2.0)},  // z lewej w odległości 3
		{geometry.NewAABB(geometry.NewVec(14.0, 9), 2.0, 2.0)}, // z prawej w odległości 3
		{geometry.NewAABB(geometry.NewVec(9.0, 5), 2.0, 2.0)},  // powyzej w odleglosci 2
		{geometry.NewAABB(geometry.NewVec(6.0, 6), 2.0, 2.0)},  // powyzej w odleglosci 2
	}
	for _, item := range items {
		qtree.Add(item)
	}

	// Target to pierwszy box
	target := &ShapeItem[float64]{geometry.NewAABB(geometry.NewVec(9.0, 9), 2.0, 2.0)}

	// Szukamy sąsiadów targeta z marginesem 1.5
	neighbors := qtree.FindNeighbors(target, 1.5)

	for _, neighbor := range neighbors {
		fmt.Println(neighbor.Bound())
	}

	// Output:
	// {(6,6) (8,8)}
	// {(9,7) (11,9)}
	// {(7,9) (9,11)}
	// {(11,9) (13,11)}
	// {(9,11) (11,13)}

}
