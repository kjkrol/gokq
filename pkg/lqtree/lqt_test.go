package lqtree

import "testing"

func strPtr(s string) *string { return &s }

func TestLinearQuadTreeInsertAndGet(t *testing.T) {
	qt := NewLinearQuadTree[string](Size64)

	entries := []Entry[string]{
		{Pos: Pos{X: 0, Y: 0}, Value: strPtr("a")},
		{Pos: Pos{X: 1, Y: 2}, Value: strPtr("b")},
		{Pos: Pos{X: 3, Y: 3}, Value: strPtr("c")},
		{Pos: Pos{X: 7, Y: 1}, Value: strPtr("d")},
	}

	qt.BulkInsert(entries)

	for _, e := range entries {
		val, ok := qt.Get(e.X, e.Y)
		if !ok {
			t.Fatalf("expected value at %v", e.Pos)
		}
		if val == nil || *val != *e.Value {
			t.Fatalf("unexpected value at %v: got %v want %v", e.Pos, val, e.Value)
		}
	}
	if qt.Count() != uint64(len(entries)) {
		t.Fatalf("unexpected count: got %d want %d", qt.Count(), len(entries))
	}
}

func TestLinearQuadTreeRemove(t *testing.T) {
	qt := NewLinearQuadTree[string](Size64)

	entries := []Entry[string]{
		{Pos: Pos{X: 0, Y: 0}, Value: strPtr("a")},
		{Pos: Pos{X: 1, Y: 1}, Value: strPtr("b")},
		{Pos: Pos{X: 2, Y: 2}, Value: strPtr("c")},
		{Pos: Pos{X: 3, Y: 3}, Value: strPtr("d")},
	}
	qt.BulkInsert(entries)

	qt.BulkRemove([]Entry[string]{
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
		val, ok := qt.Get(c.pos.X, c.pos.Y)
		if c.present != ok {
			t.Fatalf("presence mismatch at %v: got %v want %v", c.pos, ok, c.present)
		}
		if c.present && (val == nil || *val != c.expected) {
			t.Fatalf("unexpected value at %v: got %v want %v", c.pos, val, c.expected)
		}
	}
	if qt.Count() != 2 {
		t.Fatalf("unexpected count after removals: got %d want 2", qt.Count())
	}
}

func TestLinearQuadTreeUpdateWithMove(t *testing.T) {
	qt := NewLinearQuadTree[string](Size64)

	entries := []Entry[string]{
		{Pos: Pos{X: 0, Y: 0}, Value: strPtr("a")},
		{Pos: Pos{X: 1, Y: 1}, Value: strPtr("b")},
		{Pos: Pos{X: 2, Y: 2}, Value: strPtr("c")},
		{Pos: Pos{X: 3, Y: 3}, Value: strPtr("d")},
	}
	qt.BulkInsert(entries)

	updates := NewEntriesMove[string](2)
	updates.Append(strPtr("b2"), Pos{X: 1, Y: 1}, Pos{X: 4, Y: 1})
	updates.Append(strPtr("d2"), Pos{X: 3, Y: 3}, Pos{X: 4, Y: 4})
	qt.BulkMove(updates)

	checks := []struct {
		pos      Pos
		present  bool
		expected string
	}{
		{pos: Pos{X: 0, Y: 0}, present: true, expected: "a"},
		{pos: Pos{X: 1, Y: 1}, present: false},
		{pos: Pos{X: 2, Y: 2}, present: true, expected: "c"},
		{pos: Pos{X: 3, Y: 3}, present: false},
		{pos: Pos{X: 4, Y: 1}, present: true, expected: "b2"},
		{pos: Pos{X: 4, Y: 4}, present: true, expected: "d2"},
	}

	for _, c := range checks {
		val, ok := qt.Get(c.pos.X, c.pos.Y)
		if c.present != ok {
			t.Fatalf("presence mismatch at %v: got %v want %v", c.pos, ok, c.present)
		}
		if c.present && (val == nil || *val != c.expected) {
			t.Fatalf("unexpected value at %v: got %v want %v", c.pos, val, c.expected)
		}
	}
	if qt.Count() != 4 {
		t.Fatalf("unexpected count after moves: got %d want 4", qt.Count())
	}
}

func TestLinearQuadTreeQueryRange(t *testing.T) {
	qt := NewLinearQuadTree[string](Size64)

	// Clustered points: one center with 4 neighbors
	cluster := []Entry[string]{
		{Pos: Pos{X: 3, Y: 3}, Value: strPtr("center")},
		{Pos: Pos{X: 2, Y: 3}, Value: strPtr("west")},
		{Pos: Pos{X: 4, Y: 3}, Value: strPtr("east")},
		{Pos: Pos{X: 3, Y: 2}, Value: strPtr("north")},
		{Pos: Pos{X: 3, Y: 4}, Value: strPtr("south")},
	}

	// Far points that should not be returned
	far := []Entry[string]{
		{Pos: Pos{X: 0, Y: 0}, Value: strPtr("far1")},
		{Pos: Pos{X: 7, Y: 7}, Value: strPtr("far2")},
		{Pos: Pos{X: 6, Y: 1}, Value: strPtr("far3")},
		{Pos: Pos{X: 1, Y: 6}, Value: strPtr("far4")},
	}

	qt.BulkInsert(append(cluster, far...))

	buf := make([]*string, 0)
	results := qt.QueryRange(AABB{
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
