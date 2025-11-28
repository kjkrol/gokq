package bucketgrid

import (
	"testing"

	"github.com/kjkrol/gokq/pkg/pow2grid"
)

func TestBucketGridInsertAndGet(t *testing.T) {
	maxXY := pow2grid.Size32x32.Side()
	grid := NewBucketGrid[string](
		pow2grid.Size8x8,
		pow2grid.AABB{
			Min: pow2grid.Pos{X: 0, Y: 0},
			Max: pow2grid.Pos{X: maxXY, Y: maxXY},
		},
	)

	entries := []pow2grid.Entry[string]{
		{Pos: pow2grid.Pos{X: 0, Y: 0}, Value: strPtr("a")},
		{Pos: pow2grid.Pos{X: 1, Y: 2}, Value: strPtr("b")},
		{Pos: pow2grid.Pos{X: 3, Y: 3}, Value: strPtr("c")},
		{Pos: pow2grid.Pos{X: 7, Y: 1}, Value: strPtr("d")},
	}

	grid.BulkInsert(entries)

	for _, e := range entries {
		val, ok := grid.Get(e.X, e.Y)
		if !ok {
			t.Fatalf("expected value at %v", e.Pos)
		}
		if val == nil || *val != *e.Value {
			t.Fatalf("unexpected value at %v: got %v want %v", e.Pos, val, e.Value)
		}
	}
	if grid.Count() != uint64(len(entries)) {
		t.Fatalf("unexpected count: got %d want %d", grid.Count(), len(entries))
	}
}

func TestBucketGridRemove(t *testing.T) {
	maxXY := pow2grid.Size32x32.Side()
	grid := NewBucketGrid[string](
		pow2grid.Size8x8,
		pow2grid.AABB{
			Min: pow2grid.Pos{X: 0, Y: 0},
			Max: pow2grid.Pos{X: maxXY, Y: maxXY},
		},
	)

	entries := []pow2grid.Entry[string]{
		{Pos: pow2grid.Pos{X: 0, Y: 0}, Value: strPtr("a")},
		{Pos: pow2grid.Pos{X: 1, Y: 1}, Value: strPtr("b")},
		{Pos: pow2grid.Pos{X: 2, Y: 2}, Value: strPtr("c")},
		{Pos: pow2grid.Pos{X: 3, Y: 3}, Value: strPtr("d")},
	}
	grid.BulkInsert(entries)

	grid.BulkRemove([]pow2grid.Entry[string]{
		{Pos: pow2grid.Pos{X: 1, Y: 1}},
		{Pos: pow2grid.Pos{X: 3, Y: 3}},
	})

	checks := []struct {
		pos      pow2grid.Pos
		present  bool
		expected string
	}{
		{pos: pow2grid.Pos{X: 0, Y: 0}, present: true, expected: "a"},
		{pos: pow2grid.Pos{X: 1, Y: 1}, present: false},
		{pos: pow2grid.Pos{X: 2, Y: 2}, present: true, expected: "c"},
		{pos: pow2grid.Pos{X: 3, Y: 3}, present: false},
	}

	for _, c := range checks {
		val, ok := grid.Get(c.pos.X, c.pos.Y)
		if c.present != ok {
			t.Fatalf("presence mismatch at %v: got %v want %v", c.pos, ok, c.present)
		}
		if c.present && (val == nil || *val != c.expected) {
			t.Fatalf("unexpected value at %v: got %v want %v", c.pos, val, c.expected)
		}
	}
	if grid.Count() != 2 {
		t.Fatalf("unexpected count after removals: got %d want 2", grid.Count())
	}
}

func TestBucketGridMove(t *testing.T) {
	maxXY := pow2grid.Size32x32.Side()
	grid := NewBucketGrid[string](
		pow2grid.Size8x8,
		pow2grid.AABB{
			Min: pow2grid.Pos{X: 0, Y: 0},
			Max: pow2grid.Pos{X: maxXY, Y: maxXY},
		},
	)

	b := strPtr("b")
	d := strPtr("d")

	entries := []pow2grid.Entry[string]{
		{Pos: pow2grid.Pos{X: 0, Y: 0}, Value: strPtr("a")},
		{Pos: pow2grid.Pos{X: 1, Y: 1}, Value: b},
		{Pos: pow2grid.Pos{X: 2, Y: 2}, Value: strPtr("c")},
		{Pos: pow2grid.Pos{X: 3, Y: 3}, Value: d},
	}
	grid.BulkInsert(entries)

	updates := pow2grid.NewEntriesMove[string](2)
	updates.Append(b, pow2grid.Pos{X: 1, Y: 1}, pow2grid.Pos{X: 4, Y: 1})
	updates.Append(d, pow2grid.Pos{X: 3, Y: 3}, pow2grid.Pos{X: 5, Y: 5})
	grid.BulkMove(updates)

	checks := []struct {
		pos      pow2grid.Pos
		present  bool
		expected string
	}{
		{pos: pow2grid.Pos{X: 0, Y: 0}, present: true, expected: "a"},
		{pos: pow2grid.Pos{X: 1, Y: 1}, present: false},
		{pos: pow2grid.Pos{X: 2, Y: 2}, present: true, expected: "c"},
		{pos: pow2grid.Pos{X: 3, Y: 3}, present: false},
		{pos: pow2grid.Pos{X: 4, Y: 1}, present: true, expected: "b"},
		{pos: pow2grid.Pos{X: 5, Y: 5}, present: true, expected: "d"},
	}

	for _, c := range checks {
		val, ok := grid.Get(c.pos.X, c.pos.Y)
		if c.present != ok {
			t.Fatalf("presence mismatch at %v: got %v want %v", c.pos, ok, c.present)
		}
		if c.present && (val == nil || *val != c.expected) {
			t.Fatalf("unexpected value at %v: got %v want %v", c.pos, val, c.expected)
		}
	}
	if grid.Count() != 4 {
		t.Fatalf("unexpected count after moves: got %d want 4", grid.Count())
	}
}

func TestBucketGridQueryRange(t *testing.T) {
	maxXY := pow2grid.Size32x32.Side()
	grid := NewBucketGrid[string](
		pow2grid.Size8x8,
		pow2grid.AABB{
			Min: pow2grid.Pos{X: 0, Y: 0},
			Max: pow2grid.Pos{X: maxXY, Y: maxXY},
		},
	)

	cluster := []pow2grid.Entry[string]{
		{Pos: pow2grid.Pos{X: 3, Y: 3}, Value: strPtr("center")},
		{Pos: pow2grid.Pos{X: 2, Y: 3}, Value: strPtr("west")},
		{Pos: pow2grid.Pos{X: 4, Y: 3}, Value: strPtr("east")},
		{Pos: pow2grid.Pos{X: 3, Y: 2}, Value: strPtr("north")},
		{Pos: pow2grid.Pos{X: 3, Y: 4}, Value: strPtr("south")},
	}

	far := []pow2grid.Entry[string]{
		{Pos: pow2grid.Pos{X: 0, Y: 0}, Value: strPtr("far1")},
		{Pos: pow2grid.Pos{X: 7, Y: 7}, Value: strPtr("far2")},
		{Pos: pow2grid.Pos{X: 6, Y: 1}, Value: strPtr("far3")},
		{Pos: pow2grid.Pos{X: 1, Y: 6}, Value: strPtr("far4")},
	}

	grid.BulkInsert(append(cluster, far...))

	buf := make([]*string, 16)
	n := grid.QueryRange(pow2grid.AABB{
		Min: pow2grid.Pos{X: 2, Y: 2},
		Max: pow2grid.Pos{X: 4, Y: 4},
	}, buf)

	found := make(map[string]bool)
	for i := 0; i < n; i++ {
		v := buf[i]
		if v != nil {
			found[*v] = true
		}
	}

	expected := []string{"center", "west", "east", "north", "south"}
	for _, want := range expected {
		if !found[want] {
			t.Fatalf("expected to find %q in query results", want)
		}
	}
	if len(found) != len(expected) {
		t.Fatalf("unexpected extra results: got %v", found)
	}
}

func TestBucketGridQueryRangeCrossChunk(t *testing.T) {
	maxXY := pow2grid.Size16x16.Side()
	grid := NewBucketGrid[string](
		pow2grid.Size4x4,
		pow2grid.AABB{
			Min: pow2grid.Pos{X: 0, Y: 0},
			Max: pow2grid.Pos{X: maxXY, Y: maxXY},
		},
	)

	cluster := []pow2grid.Entry[string]{
		{Pos: pow2grid.Pos{X: 4, Y: 4}, Value: strPtr("center")}, // chunk origin (size1 -> 1x1)
		{Pos: pow2grid.Pos{X: 3, Y: 4}, Value: strPtr("west")},   // neighboring chunk to the west
		{Pos: pow2grid.Pos{X: 5, Y: 4}, Value: strPtr("north")},  // neighboring chunk to the north
		{Pos: pow2grid.Pos{X: 4, Y: 3}, Value: strPtr("east")},   // neighboring chunk to the east
		{Pos: pow2grid.Pos{X: 4, Y: 5}, Value: strPtr("south")},  // neighboring chunk to the south
	}

	far := []pow2grid.Entry[string]{
		{Pos: pow2grid.Pos{X: 0, Y: 0}, Value: strPtr("far1")},
		{Pos: pow2grid.Pos{X: 7, Y: 7}, Value: strPtr("far2")},
		{Pos: pow2grid.Pos{X: 6, Y: 1}, Value: strPtr("far3")},
		{Pos: pow2grid.Pos{X: 10, Y: 10}, Value: strPtr("far4")},
	}

	grid.BulkInsert(append(cluster, far...))

	buf := make([]*string, 16)
	n := grid.QueryRange(pow2grid.AABB{
		Min: pow2grid.Pos{X: 3, Y: 3},
		Max: pow2grid.Pos{X: 5, Y: 5},
	}, buf)

	found := make(map[string]bool)
	for i := 0; i < n; i++ {
		v := buf[i]
		if v != nil {
			found[*v] = true
		}
	}

	expected := []string{"center", "west", "east", "north", "south"}
	for _, want := range expected {
		if !found[want] {
			t.Fatalf("expected to find %q in query results", want)
		}
	}
	if len(found) != len(expected) {
		t.Fatalf("unexpected extra results: got %v", found)
	}
}
