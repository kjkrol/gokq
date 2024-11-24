package quadtree

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geometry"
)

func TestBox_newBox(t *testing.T) {
	box := newBox(geometry.ZERO_INT_VEC, geometry.Vec[int]{X: 10, Y: 10})
	expected := geometry.Vec[int]{X: 5, Y: 5}
	if box.center != expected {
		t.Errorf("center %v not equal to expected %v", box.center, expected)
	}
}

func TestBox_buildBox(t *testing.T) {
	center := geometry.Vec[int]{X: 5, Y: 5}
	box := buildBox(center, 2)
	expectedTopLeft := geometry.Vec[int]{X: 3, Y: 3}
	expectedBottomRight := geometry.Vec[int]{X: 7, Y: 7}
	if box.topLeft != expectedTopLeft {
		t.Errorf("topLeft %v not equal to expected %v", box.topLeft, expectedTopLeft)
	}
	if box.bottomRight != expectedBottomRight {
		t.Errorf("bottomRight %v not equal to expected %v", box.bottomRight, expectedBottomRight)
	}
}

func TestBox_split(t *testing.T) {
	parentBox := newBox(geometry.ZERO_INT_VEC, geometry.Vec[int]{X: 10, Y: 10})
	splitted := parentBox.split()

	expected := [4]box[int]{
		newBox(geometry.ZERO_INT_VEC, geometry.Vec[int]{X: 5, Y: 5}),
		newBox(geometry.Vec[int]{X: 5, Y: 0}, geometry.Vec[int]{X: 10, Y: 5}),
		newBox(geometry.Vec[int]{X: 0, Y: 5}, geometry.Vec[int]{X: 5, Y: 10}),
		newBox(geometry.Vec[int]{X: 5, Y: 5}, geometry.Vec[int]{X: 10, Y: 10}),
	}
	for i := 0; i < 4; i++ {
		if splitted[i] != expected[i] {
			t.Errorf("split %v not equal to expected %v", splitted[i], expected[i])
		}
	}
}

func TestBox_intersects(t *testing.T) {

	intersects := []struct{ box1, box2 box[int] }{
		{
			box1: newBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: newBox(geometry.Vec[int]{X: 2, Y: 2}, geometry.Vec[int]{X: 4, Y: 4}),
		},
		{
			box1: newBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: newBox(geometry.Vec[int]{X: 4, Y: 2}, geometry.Vec[int]{X: 6, Y: 4}),
		},
		{
			box1: newBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: newBox(geometry.Vec[int]{X: 2, Y: 4}, geometry.Vec[int]{X: 4, Y: 6}),
		},
		{
			box1: newBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: newBox(geometry.Vec[int]{X: 4, Y: 4}, geometry.Vec[int]{X: 6, Y: 6}),
		},
		{
			box1: newBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: newBox(geometry.Vec[int]{X: 2, Y: 2}, geometry.Vec[int]{X: 6, Y: 6}),
		},
		{
			box1: newBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: newBox(geometry.Vec[int]{X: 2, Y: 2}, geometry.Vec[int]{X: 3, Y: 4}),
		},
		{
			box1: newBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: newBox(geometry.Vec[int]{X: 2, Y: 4}, geometry.Vec[int]{X: 3, Y: 6}),
		},
		{
			box1: newBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: newBox(geometry.Vec[int]{X: 2, Y: 2}, geometry.Vec[int]{X: 6, Y: 3}),
		},
		{
			box1: newBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: newBox(geometry.Vec[int]{X: 2, Y: 5}, geometry.Vec[int]{X: 6, Y: 6}),
		},
		{
			box1: newBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: newBox(geometry.Vec[int]{X: 5, Y: 2}, geometry.Vec[int]{X: 6, Y: 6}),
		},
		{
			box1: newBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: newBox(geometry.Vec[int]{X: 5, Y: 2}, geometry.Vec[int]{X: 6, Y: 6}),
		},
	}

	for _, intersection := range intersects {
		if !intersection.box1.intersects(intersection.box2) {
			t.Errorf("Box1 %v should intersects with box2 %v", intersection.box1, intersection.box2)
		}
		if !intersection.box2.intersects(intersection.box1) {
			t.Errorf("Box2 %v should intersects with box1 %v", intersection.box2, intersection.box1)
		}
	}

	notIntersects := []struct{ box1, box2 box[int] }{
		{
			box1: newBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: newBox(geometry.Vec[int]{X: 1, Y: 1}, geometry.Vec[int]{X: 2, Y: 2}),
		},
		{
			box1: newBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: newBox(geometry.Vec[int]{X: 6, Y: 0}, geometry.Vec[int]{X: 9, Y: 9}),
		},
		{
			box1: newBox(geometry.Vec[int]{X: 3, Y: 3}, geometry.Vec[int]{X: 5, Y: 5}),
			box2: newBox(geometry.Vec[int]{X: 0, Y: 6}, geometry.Vec[int]{X: 9, Y: 9}),
		},
	}
	for _, intersection := range notIntersects {
		if intersection.box1.intersects(intersection.box2) {
			t.Errorf("Box1 %v should not intersects with box2 %v", intersection.box1, intersection.box2)
		}
		if intersection.box2.intersects(intersection.box1) {
			t.Errorf("Box2 %v should not intersects with box1 %v", intersection.box2, intersection.box1)
		}
	}
}

func TestBox_intersectsCyclic(t *testing.T) {
	intersects := []struct{ box1, box2 box[int] }{
		{
			box1: newBox(geometry.Vec[int]{X: 5, Y: 5}, geometry.Vec[int]{X: 15, Y: 15}),
			box2: newBox(geometry.Vec[int]{X: 95, Y: 95}, geometry.Vec[int]{X: 105, Y: 105}),
		},
	}
	size := geometry.Vec[int]{X: 100, Y: 100}
	plane := geometry.NewBoundedPlane(size.X, size.Y)

	for _, intersection := range intersects {
		wrappedBoxes := wrapBoxCyclic(intersection.box2, size, plane.Contains)
		if !intersection.box1.intersectsAny(wrappedBoxes) {
			t.Errorf("Box1 %v should intersects with Box2 %v", intersection.box1, intersection.box2)
		}
	}
}
