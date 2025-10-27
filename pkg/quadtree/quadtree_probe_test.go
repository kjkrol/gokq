package quadtree

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geometry"
	"github.com/kjkrol/gokg/pkg/geometry/spatial"
)

func TestQuadTree_ProbeLine(t *testing.T) {
	boundedPlane := geometry.NewBoundedPlane(10, 10)
	qtree := NewQuadTree(boundedPlane)
	line := spatial.NewLine(spatial.Vec[int]{X: 1, Y: 1}, spatial.Vec[int]{X: 4, Y: 3})

	probes := qtree.probe(&line, 1)
	if len(probes) != 1 {
		t.Fatalf("expected a single probe rectangle, got %d", len(probes))
	}
	expectedProbe := spatial.NewRectangle(spatial.Vec[int]{X: 0, Y: 0}, spatial.Vec[int]{X: 5, Y: 4})
	if !probes[0].Equals(expectedProbe) {
		t.Errorf("expected expanded rectangle %v, got %v", expectedProbe, probes[0])
	}
}

func TestQuadTree_ProbePolygon(t *testing.T) {
	poly := spatial.NewPolygon(
		spatial.Vec[int]{X: 0, Y: 0},
		spatial.Vec[int]{X: 2, Y: 0},
		spatial.Vec[int]{X: 1, Y: 2},
	)
	boundedPlane := geometry.NewBoundedPlane(100, 100)
	qtree := NewQuadTree(boundedPlane)

	probes := qtree.probe(&poly, 1)
	if len(probes) != 1 {
		t.Fatalf("expected single probe rectangle, got %d", len(probes))
	}

	expectedTopLeft := spatial.Vec[int]{X: -1, Y: -1}
	expectedBottomRight := spatial.Vec[int]{X: 3, Y: 3}
	if probes[0].TopLeft != expectedTopLeft {
		t.Errorf("expected probe top-left %v, got %v", expectedTopLeft, probes[0].TopLeft)
	}
	if probes[0].BottomRight != expectedBottomRight {
		t.Errorf("expected probe bottom-right %v, got %v", expectedBottomRight, probes[0].BottomRight)
	}
}

func TestQuadTree_ForCyclicBoundedPlane_ProbeRectangle(t *testing.T) {
	rect := spatial.NewRectangle(spatial.Vec[int]{X: 8, Y: 8}, spatial.Vec[int]{X: 10, Y: 10})
	boundedPlane := geometry.NewCyclicBoundedPlane(10, 10)
	qtree := NewQuadTree(boundedPlane)
	probes := qtree.probe(&rect, 0)
	if len(probes) < 2 {
		t.Fatalf("expected wrapped probes, got %d", len(probes))
	}

	wrappedFound := false
	for _, p := range probes {
		if !p.Equals(rect) {
			wrappedFound = true
		}
	}
	if !wrappedFound {
		t.Errorf("expected wrapped rectangle in probes %v", probes)
	}
}

func TestRectangle_IntersectsCyclic(t *testing.T) {
	intersects := []struct{ rect1, rect2 spatial.Rectangle[int] }{
		{
			rect1: spatial.NewRectangle(spatial.Vec[int]{X: 5, Y: 5}, spatial.Vec[int]{X: 15, Y: 15}),
			rect2: spatial.NewRectangle(spatial.Vec[int]{X: 95, Y: 95}, spatial.Vec[int]{X: 105, Y: 105}),
		},
	}
	size := spatial.Vec[int]{X: 100, Y: 100}
	boundedPlane := geometry.NewBoundedPlane(size.X, size.Y)
	qtree := NewQuadTree(boundedPlane)

	for _, intersection := range intersects {
		wrapped := qtree.wrapRectangleCyclic(intersection.rect2)
		if !intersection.rect1.IntersectsAny(wrapped) {
			t.Errorf("rect1 %v should intersect with rect2 %v", intersection.rect1, intersection.rect2)
		}
	}
}
