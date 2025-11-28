package zorder

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/kjkrol/gokq/pkg/pow2grid"
)

const benchMaxXY = uint32(4096)
const bucketSize = pow2grid.Size256x256

func BenchmarkBucketGridBulkMove(b *testing.B) {
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
			benchmarkBucketGridBulkMove(b, tc.totalEntries, tc.movingEntries)
		})
	}
}

func benchmarkBucketGridBulkMove(b *testing.B, totalEntries, movingEntries int) {
	src := rand.New(rand.NewSource(1))

	entries := make([]pow2grid.Entry[string], totalEntries)
	for i := range totalEntries {
		entries[i] = pow2grid.Entry[string]{
			Pos:   pow2grid.Pos{X: randCoord(src, benchMaxXY), Y: randCoord(src, benchMaxXY)},
			Value: strPtr(fmt.Sprintf("v%d", i)),
		}
	}

	forwardMoves := pow2grid.NewEntriesMove[string](movingEntries)
	backwardMoves := pow2grid.NewEntriesMove[string](movingEntries)
	for i := range movingEntries {
		oldPos := entries[i].Pos
		newPos := pow2grid.Pos{X: randCoord(src, benchMaxXY), Y: randCoord(src, benchMaxXY)}

		forwardMoves.Append(entries[i].Value, oldPos, newPos)
		backwardMoves.Append(entries[i].Value, newPos, oldPos)
	}

	grid := NewBucketGrid[string](bucketSize, pow2grid.AABB{
		Min: pow2grid.Pos{X: 0, Y: 0},
		Max: pow2grid.Pos{X: benchMaxXY, Y: benchMaxXY},
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

func BenchmarkBucketGridQueryRange(b *testing.B) {
	cases := []struct {
		name         string
		totalEntries int
		queryCount   int
		querySize    uint32
	}{
		// nazwa opisuje: [liczba wpisów]-[rozmiar okna]
		{"1k-64sz", 1000, 128, 64},
		{"100k-64sz", 100000, 128, 64},
		{"500k-64sz", 500000, 128, 64},
		{"100k-256sz", 100000, 64, 256},
		{"500k-256sz", 500000, 64, 256},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			benchmarkBucketGridQueryRange(b, tc.totalEntries, tc.queryCount, tc.querySize)
		})
	}
}

func benchmarkBucketGridQueryRange(b *testing.B, totalEntries, queryCount int, querySize uint32) {
	src := rand.New(rand.NewSource(2))

	entries := make([]pow2grid.Entry[string], totalEntries)
	for i := range totalEntries {
		entries[i] = pow2grid.Entry[string]{
			Pos:   pow2grid.Pos{X: randCoord(src, benchMaxXY), Y: randCoord(src, benchMaxXY)},
			Value: strPtr(fmt.Sprintf("v%d", i)),
		}
	}

	grid := NewBucketGrid[string](bucketSize, pow2grid.AABB{
		Min: pow2grid.Pos{X: 0, Y: 0},
		Max: pow2grid.Pos{X: benchMaxXY, Y: benchMaxXY},
	})
	grid.BulkInsert(entries)

	// Pre-generujemy queryCount różnych okien zapytań.
	// W benchmarku każda iteracja b.Loop() wykonuje dokładnie JEDNO QueryRange,
	// za każdym razem biorąc inne okno (rotacja i%queryCount).
	queries := make([]pow2grid.AABB, queryCount)
	maxStart := benchMaxXY - querySize
	for i := range queryCount {
		minX := randCoord(src, maxStart)
		minY := randCoord(src, maxStart)
		queries[i] = pow2grid.AABB{
			Min: pow2grid.Pos{X: minX, Y: minY},
			Max: pow2grid.Pos{X: minX + querySize, Y: minY + querySize},
		}
	}

	results := make([]*string, 0, 1024)

	for i := 0; b.Loop(); i++ {
		q := queries[i%queryCount]
		results = grid.QueryRange(q, results[:0])
	}
}
