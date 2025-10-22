package quadtree

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geometry"
)

func TestBox_newBox(t *testing.T) {
	box := NewBox(geometry.ZERO_INT_VEC, geometry.Vec[int]{X: 10, Y: 10})
	expected := geometry.Vec[int]{X: 5, Y: 5}
	if box.Center != expected {
		t.Errorf("center %v not equal to expected %v", box.Center, expected)
	}
}

func TestBox_buildBox(t *testing.T) {
	center := geometry.Vec[int]{X: 5, Y: 5}
	box := BuildBox(center, 2)
	expectedTopLeft := geometry.Vec[int]{X: 3, Y: 3}
	expectedBottomRight := geometry.Vec[int]{X: 7, Y: 7}
	if box.TopLeft != expectedTopLeft {
		t.Errorf("topLeft %v not equal to expected %v", box.TopLeft, expectedTopLeft)
	}
	if box.BottomRight != expectedBottomRight {
		t.Errorf("bottomRight %v not equal to expected %v", box.BottomRight, expectedBottomRight)
	}
}

func TestBox_Split(t *testing.T) {
	parentBox := NewBox(geometry.ZERO_INT_VEC, geometry.Vec[int]{X: 10, Y: 10})
	splitted := parentBox.Split()

	expected := [4]Box[int]{
		NewBox(geometry.ZERO_INT_VEC, geometry.Vec[int]{X: 5, Y: 5}),
		NewBox(geometry.Vec[int]{X: 5, Y: 0}, geometry.Vec[int]{X: 10, Y: 5}),
		NewBox(geometry.Vec[int]{X: 0, Y: 5}, geometry.Vec[int]{X: 5, Y: 10}),
		NewBox(geometry.Vec[int]{X: 5, Y: 5}, geometry.Vec[int]{X: 10, Y: 10}),
	}
	for i := 0; i < 4; i++ {
		if splitted[i] != expected[i] {
			t.Errorf("split %v not equal to expected %v", splitted[i], expected[i])
		}
	}
}

func TestBox_intersects(t *testing.T) {

	intersects := []struct{ box1, box2 Box[int] }{
		{
			box1: NewBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: NewBox(geometry.Vec[int]{X: 2, Y: 2}, geometry.Vec[int]{X: 4, Y: 4}),
		},
		{
			box1: NewBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: NewBox(geometry.Vec[int]{X: 4, Y: 2}, geometry.Vec[int]{X: 6, Y: 4}),
		},
		{
			box1: NewBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: NewBox(geometry.Vec[int]{X: 2, Y: 4}, geometry.Vec[int]{X: 4, Y: 6}),
		},
		{
			box1: NewBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: NewBox(geometry.Vec[int]{X: 4, Y: 4}, geometry.Vec[int]{X: 6, Y: 6}),
		},
		{
			box1: NewBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: NewBox(geometry.Vec[int]{X: 2, Y: 2}, geometry.Vec[int]{X: 6, Y: 6}),
		},
		{
			box1: NewBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: NewBox(geometry.Vec[int]{X: 2, Y: 2}, geometry.Vec[int]{X: 3, Y: 4}),
		},
		{
			box1: NewBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: NewBox(geometry.Vec[int]{X: 2, Y: 4}, geometry.Vec[int]{X: 3, Y: 6}),
		},
		{
			box1: NewBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: NewBox(geometry.Vec[int]{X: 2, Y: 2}, geometry.Vec[int]{X: 6, Y: 3}),
		},
		{
			box1: NewBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: NewBox(geometry.Vec[int]{X: 2, Y: 5}, geometry.Vec[int]{X: 6, Y: 6}),
		},
		{
			box1: NewBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: NewBox(geometry.Vec[int]{X: 5, Y: 2}, geometry.Vec[int]{X: 6, Y: 6}),
		},
		{
			box1: NewBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: NewBox(geometry.Vec[int]{X: 5, Y: 2}, geometry.Vec[int]{X: 6, Y: 6}),
		},
	}

	for _, intersection := range intersects {
		if !intersection.box1.Intersects(intersection.box2) {
			t.Errorf("Box1 %v should intersects with box2 %v", intersection.box1, intersection.box2)
		}
		if !intersection.box2.Intersects(intersection.box1) {
			t.Errorf("Box2 %v should intersects with box1 %v", intersection.box2, intersection.box1)
		}
	}

	notIntersects := []struct{ box1, box2 Box[int] }{
		{
			box1: NewBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: NewBox(geometry.Vec[int]{X: 1, Y: 1}, geometry.Vec[int]{X: 2, Y: 2}),
		},
		{
			box1: NewBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: NewBox(geometry.Vec[int]{X: 6, Y: 0}, geometry.Vec[int]{X: 9, Y: 9}),
		},
		{
			box1: NewBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: NewBox(geometry.Vec[int]{X: 0, Y: 6}, geometry.Vec[int]{X: 9, Y: 9}),
		},
	}
	for _, intersection := range notIntersects {
		if intersection.box1.Intersects(intersection.box2) {
			t.Errorf("Box1 %v should not intersects with box2 %v", intersection.box1, intersection.box2)
		}
		if intersection.box2.Intersects(intersection.box1) {
			t.Errorf("Box2 %v should not intersects with box1 %v", intersection.box2, intersection.box1)
		}
	}
}

func TestBox_intersectsCyclic(t *testing.T) {
	intersects := []struct{ box1, box2 Box[int] }{
		{
			box1: NewBox(geometry.Vec[int]{X: 5, Y: 5}, geometry.Vec[int]{X: 15, Y: 15}),
			box2: NewBox(geometry.Vec[int]{X: 95, Y: 95}, geometry.Vec[int]{X: 105, Y: 105}),
		},
	}
	size := geometry.Vec[int]{X: 100, Y: 100}
	plane := geometry.NewBoundedPlane(size.X, size.Y)

	for _, intersection := range intersects {
		wrappedBoxes := WrapBoxCyclic(intersection.box2, size, plane.Contains)
		if !intersection.box1.intersectsAny(wrappedBoxes) {
			t.Errorf("Box1 %v should intersects with Box2 %v", intersection.box1, intersection.box2)
		}
	}
}

func TestBox_IntersectsAny_ReturnsFalse(t *testing.T) {
	// bazowy box w lewym górnym rogu
	base := Box[int]{
		TopLeft:     geometry.Vec[int]{X: 0, Y: 0},
		BottomRight: geometry.Vec[int]{X: 10, Y: 10},
		Center:      geometry.Vec[int]{X: 5, Y: 5},
	}

	// inne boxy są daleko od base → brak przecięcia
	others := []Box[int]{
		{TopLeft: geometry.Vec[int]{X: 20, Y: 20}, BottomRight: geometry.Vec[int]{X: 30, Y: 30}},
		{TopLeft: geometry.Vec[int]{X: 40, Y: 0}, BottomRight: geometry.Vec[int]{X: 50, Y: 10}},
		{TopLeft: geometry.Vec[int]{X: 0, Y: 40}, BottomRight: geometry.Vec[int]{X: 10, Y: 50}},
	}

	if base.intersectsAny(others) {
		t.Errorf("expected intersectsAny to return false, but got true")
	}
}
