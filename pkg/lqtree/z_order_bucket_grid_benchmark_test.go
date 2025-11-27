package lqtree

import (
	"fmt"
	"math/rand"
	"testing"
)

const benchMaxXY = uint32(4096)

func BenchmarkZOrderBucketGridBulkMove(b *testing.B) {
	cases := []struct {
		name          string
		totalEntries  int
		movingEntries int
	}{
		{"500-200", 500, 200},
		{"5k-2k", 5000, 2000},
		{"50k-20k", 50000, 20000},
		{"100k-40k", 100000, 40000},
		{"500k-200k", 500000, 200000},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			benchmarkZOrderBucketGridBulkMove(b, tc.totalEntries, tc.movingEntries)
		})
	}
}

func benchmarkZOrderBucketGridBulkMove(b *testing.B, totalEntries, movingEntries int) {
	src := rand.New(rand.NewSource(1))

	entries := make([]Entry[string], totalEntries)
	for i := range totalEntries {
		entries[i] = Entry[string]{
			Pos:   Pos{X: randCoord(src, benchMaxXY), Y: randCoord(src, benchMaxXY)},
			Value: strPtr(fmt.Sprintf("v%d", i)),
		}
	}

	forwardMoves := make([]EntryMove[string], movingEntries)
	backwardMoves := make([]EntryMove[string], movingEntries)
	for i := range movingEntries {
		oldPos := entries[i].Pos
		newPos := Pos{X: randCoord(src, benchMaxXY), Y: randCoord(src, benchMaxXY)}

		forwardMoves[i] = EntryMove[string]{Old: oldPos, New: newPos}
		backwardMoves[i] = EntryMove[string]{Old: newPos, New: oldPos}
	}

	grid := NewZOrderBucketGrid[string](Size65536, AABB{
		Min: Pos{0, 0},
		Max: Pos{benchMaxXY, benchMaxXY},
	})
	grid.BulkInsert(entries)

	for i := 0; b.Loop(); i++ {
		if i%2 == 0 {
			grid.BulkMove(forwardMoves)
		} else {
			grid.BulkMove(backwardMoves)
		}
	}
}

func randCoord(r *rand.Rand, max uint32) uint32 {
	return uint32(r.Intn(int(max + 1)))
}
