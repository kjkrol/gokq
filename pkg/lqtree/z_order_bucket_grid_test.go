package lqtree

import "testing"

func TestZOrderBucketGridInsertAndGet(t *testing.T) {
	maxXY := Size32x32.Resolution()
	grid := NewZOrderBucketGrid[string](Size8x8, AABB{Min: Pos{0, 0}, Max: Pos{maxXY, maxXY}})

	entries := []Entry[string]{
		{Pos: Pos{X: 0, Y: 0}, Value: strPtr("a")},
		{Pos: Pos{X: 1, Y: 2}, Value: strPtr("b")},
		{Pos: Pos{X: 3, Y: 3}, Value: strPtr("c")},
		{Pos: Pos{X: 7, Y: 1}, Value: strPtr("d")},
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

func TestZOrderBucketGridRemove(t *testing.T) {
	maxXY := Size32x32.Resolution()
	grid := NewZOrderBucketGrid[string](Size8x8, AABB{Min: Pos{0, 0}, Max: Pos{maxXY, maxXY}})

	entries := []Entry[string]{
		{Pos: Pos{X: 0, Y: 0}, Value: strPtr("a")},
		{Pos: Pos{X: 1, Y: 1}, Value: strPtr("b")},
		{Pos: Pos{X: 2, Y: 2}, Value: strPtr("c")},
		{Pos: Pos{X: 3, Y: 3}, Value: strPtr("d")},
	}
	grid.BulkInsert(entries)

	grid.BulkRemove([]Entry[string]{
		{Pos: Pos{X: 1, Y: 1}},
		{Pos: Pos{X: 3, Y: 3}},
	})

	checks := []struct {
		pos      Pos
		present  bool
		expected string
	}{
		{pos: Pos{X: 0, Y: 0}, present: true, expected: "a"},
		{pos: Pos{X: 1, Y: 1}, present: false},
		{pos: Pos{X: 2, Y: 2}, present: true, expected: "c"},
		{pos: Pos{X: 3, Y: 3}, present: false},
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

func TestZOrderBucketGridMove(t *testing.T) {
	maxXY := Size32x32.Resolution()
	grid := NewZOrderBucketGrid[string](Size8x8, AABB{Min: Pos{0, 0}, Max: Pos{maxXY, maxXY}})

	b := strPtr("b")
	d := strPtr("d")

	entries := []Entry[string]{
		{Pos: Pos{X: 0, Y: 0}, Value: strPtr("a")},
		{Pos: Pos{X: 1, Y: 1}, Value: b},
		{Pos: Pos{X: 2, Y: 2}, Value: strPtr("c")},
		{Pos: Pos{X: 3, Y: 3}, Value: d},
	}
	grid.BulkInsert(entries)

	updates := NewEntriesMove[string](2)
	updates.Append(b, Pos{X: 1, Y: 1}, Pos{X: 4, Y: 1})
	updates.Append(d, Pos{X: 3, Y: 3}, Pos{X: 5, Y: 5})
	grid.BulkMove(updates)

	checks := []struct {
		pos      Pos
		present  bool
		expected string
	}{
		{pos: Pos{X: 0, Y: 0}, present: true, expected: "a"},
		{pos: Pos{X: 1, Y: 1}, present: false},
		{pos: Pos{X: 2, Y: 2}, present: true, expected: "c"},
		{pos: Pos{X: 3, Y: 3}, present: false},
		{pos: Pos{X: 4, Y: 1}, present: true, expected: "b"},
		{pos: Pos{X: 5, Y: 5}, present: true, expected: "d"},
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

func TestZOrderBucketGridQueryRange(t *testing.T) {
	maxXY := Size32x32.Resolution()
	grid := NewZOrderBucketGrid[string](Size8x8, AABB{Min: Pos{0, 0}, Max: Pos{maxXY, maxXY}})

	cluster := []Entry[string]{
		{Pos: Pos{X: 3, Y: 3}, Value: strPtr("center")},
		{Pos: Pos{X: 2, Y: 3}, Value: strPtr("west")},
		{Pos: Pos{X: 4, Y: 3}, Value: strPtr("east")},
		{Pos: Pos{X: 3, Y: 2}, Value: strPtr("north")},
		{Pos: Pos{X: 3, Y: 4}, Value: strPtr("south")},
	}

	far := []Entry[string]{
		{Pos: Pos{X: 0, Y: 0}, Value: strPtr("far1")},
		{Pos: Pos{X: 7, Y: 7}, Value: strPtr("far2")},
		{Pos: Pos{X: 6, Y: 1}, Value: strPtr("far3")},
		{Pos: Pos{X: 1, Y: 6}, Value: strPtr("far4")},
	}

	grid.BulkInsert(append(cluster, far...))

	buf := make([]*string, 0)
	results := grid.QueryRange(AABB{
		Min: Pos{X: 2, Y: 2},
		Max: Pos{X: 4, Y: 4},
	}, buf)

	found := make(map[string]bool)
	for _, v := range results {
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

func TestZOrderBucketGridQueryRangeCrossChunk(t *testing.T) {
	maxXY := Size16x16.Resolution()
	grid := NewZOrderBucketGrid[string](Size4x4, AABB{Min: Pos{0, 0}, Max: Pos{maxXY, maxXY}})

	cluster := []Entry[string]{
		{Pos: Pos{X: 4, Y: 4}, Value: strPtr("center")}, // chunk origin (size1 -> 1x1)
		{Pos: Pos{X: 3, Y: 4}, Value: strPtr("west")},   // neighboring chunk to the west
		{Pos: Pos{X: 5, Y: 4}, Value: strPtr("north")},  // neighboring chunk to the north
		{Pos: Pos{X: 4, Y: 3}, Value: strPtr("east")},   // neighboring chunk to the east
		{Pos: Pos{X: 4, Y: 5}, Value: strPtr("south")},  // neighboring chunk to the south
	}

	far := []Entry[string]{
		{Pos: Pos{X: 0, Y: 0}, Value: strPtr("far1")},
		{Pos: Pos{X: 7, Y: 7}, Value: strPtr("far2")},
		{Pos: Pos{X: 6, Y: 1}, Value: strPtr("far3")},
		{Pos: Pos{X: 10, Y: 10}, Value: strPtr("far4")},
	}

	grid.BulkInsert(append(cluster, far...))

	buf := make([]*string, 0)
	results := grid.QueryRange(AABB{
		Min: Pos{X: 3, Y: 3},
		Max: Pos{X: 5, Y: 5},
	}, buf)

	found := make(map[string]bool)
	for _, v := range results {
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
