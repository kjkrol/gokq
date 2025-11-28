package bucketgrid

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/kjkrol/gokq/pkg/pow2grid"
)

const benchMaxXY = uint32(4096)

func BenchmarkBucketGridBulkMove(b *testing.B) {

	makeCases := func(bucketResolution pow2grid.Resolution) []struct {
		name             string
		bucketResolution pow2grid.Resolution
		totalEntries     int
		movingEntries    int
	} {
		return []struct {
			name             string
			bucketResolution pow2grid.Resolution
			totalEntries     int
			movingEntries    int
		}{
			{"500-200", bucketResolution, 500, 200},
			{"5k-2k", bucketResolution, 5000, 2000},
			{"50k-20k", bucketResolution, 50000, 20000},
			{"100k-40k", bucketResolution, 100000, 40000},
			{"500k-200k", bucketResolution, 500000, 200000},
			{"5000k-2000k", bucketResolution, 5000000, 2000000},
		}
	}

	for _, tc := range makeCases(pow2grid.Size16x16) {
		b.Run(tc.name, func(b *testing.B) {
			benchmarkBucketGridBulkMove(b, tc.bucketResolution, tc.totalEntries, tc.movingEntries)
		})
	}
}

func benchmarkBucketGridBulkMove(b *testing.B, bucketResolution pow2grid.Resolution, totalEntries, movingEntries int) {
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

	grid := NewBucketGrid[string](bucketResolution, pow2grid.AABB{
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

	bucketResolution := pow2grid.Size64x64
	aabbSize := pow2grid.Size16x16
	cases := []struct {
		name         string
		totalEntries int
		querySize    uint32
	}{
		// nazwa opisuje: [liczba wpisów]-[rozmiar bucketow]-[rozmiar okna aabb]
		{"1k-" + bucketResolution.String() + "-" + aabbSize.String(), 1000, aabbSize.Side()},
		{"10k-" + bucketResolution.String() + "-" + aabbSize.String(), 10000, aabbSize.Side()},
		{"100k-" + bucketResolution.String() + "-" + aabbSize.String(), 100000, aabbSize.Side()},
		{"500k-" + bucketResolution.String() + "-" + aabbSize.String(), 500000, aabbSize.Side()},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			benchmarkBucketGridQueryRange(b, bucketResolution, tc.totalEntries, tc.querySize)
		})
	}
}

func benchmarkBucketGridQueryRange(b *testing.B, bucketResolution pow2grid.Resolution, totalEntries int, querySize uint32) {
	src := rand.New(rand.NewSource(2))

	entries := make([]pow2grid.Entry[string], totalEntries)
	for i := range totalEntries {
		entries[i] = pow2grid.Entry[string]{
			Pos:   pow2grid.Pos{X: randCoord(src, benchMaxXY), Y: randCoord(src, benchMaxXY)},
			Value: strPtr(fmt.Sprintf("v%d", i)),
		}
	}

	grid := NewBucketGrid[string](bucketResolution, pow2grid.AABB{
		Min: pow2grid.Pos{X: 0, Y: 0},
		Max: pow2grid.Pos{X: benchMaxXY, Y: benchMaxXY},
	})
	grid.BulkInsert(entries)

	query := pow2grid.AABB{
		Min: pow2grid.Pos{X: 100, Y: 100},
		Max: pow2grid.Pos{X: 100 + querySize, Y: 100 + querySize},
	}

	results := make([]*string, totalEntries)

	for i := 0; b.Loop(); i++ {
		_ = grid.QueryRange(query, results)
	}
}
