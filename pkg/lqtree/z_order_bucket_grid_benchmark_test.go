package lqtree

import (
	"fmt"
	"math/rand"
	"testing"
)

const benchMaxXY = uint32(4096)
const bucketSize = Size65536

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
		{"5000k-2000k", 5000000, 2000000},
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
	for i := 0; i < totalEntries; i++ {
		entries[i] = Entry[string]{
			Pos:   Pos{X: randCoord(src, benchMaxXY), Y: randCoord(src, benchMaxXY)},
			Value: strPtr(fmt.Sprintf("v%d", i)),
		}
	}

	forwardMoves := NewEntriesUpdate[string](movingEntries)
	backwardMoves := NewEntriesUpdate[string](movingEntries)
	for i := 0; i < movingEntries; i++ {
		oldPos := entries[i].Pos
		newPos := Pos{X: randCoord(src, benchMaxXY), Y: randCoord(src, benchMaxXY)}

		forwardMoves.Append(entries[i].Value, oldPos, newPos)
		backwardMoves.Append(entries[i].Value, newPos, oldPos)
	}

	grid := NewZOrderBucketGrid[string](bucketSize, AABB{
		Min: Pos{0, 0},
		Max: Pos{benchMaxXY, benchMaxXY},
	})
	grid.BulkInsert(entries)

	for i := 0; b.Loop(); i++ {
		if i%2 == 0 {
			grid.BulkUpdate(forwardMoves)
		} else {
			grid.BulkUpdate(backwardMoves)
		}
	}
}

func randCoord(r *rand.Rand, max uint32) uint32 {
	return uint32(r.Intn(int(max + 1)))
}
