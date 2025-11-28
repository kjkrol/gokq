package grid

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/kjkrol/gokq/pkg/pow2grid"
)

const (
	gridBenchResolution = pow2grid.Size1024x1024
	gridMaxCoord        = (uint32(1) << gridBenchResolution) - 1
)

func BenchmarkGridQueryRange(b *testing.B) {
	gridResolution := gridBenchResolution
	aabbSize := pow2grid.Size16x16

	cases := []struct {
		name         string
		totalEntries int
		querySize    uint32
	}{
		// name: [entries]-[tree size]-[aabb size]
		{"1k-" + gridResolution.String() + "-" + aabbSize.String(), 1000, aabbSize.Side()},
		{"10k-" + gridResolution.String() + "-" + aabbSize.String(), 10000, aabbSize.Side()},
		{"100k-" + gridResolution.String() + "-" + aabbSize.String(), 100000, aabbSize.Side()},
		{"500k-" + gridResolution.String() + "-" + aabbSize.String(), 500000, aabbSize.Side()},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			benchmarkGridQueryRange(b, tc.totalEntries, tc.querySize)
		})
	}
}

func benchmarkGridQueryRange(b *testing.B, totalEntries int, querySize uint32) {
	src := rand.New(rand.NewSource(3))

	entries := make([]pow2grid.Entry[string], totalEntries)
	for i := range totalEntries {
		entries[i] = pow2grid.Entry[string]{
			Pos:   pow2grid.Pos{X: randCoord(src, gridMaxCoord), Y: randCoord(src, gridMaxCoord)},
			Value: strPtr(fmt.Sprintf("v%d", i)),
		}
	}

	grid := NewGrid[string](gridBenchResolution)
	grid.BulkInsert(entries)

	query := pow2grid.AABB{
		Min: pow2grid.Pos{X: 100, Y: 100},
		Max: pow2grid.Pos{X: 100 + querySize, Y: 100 + querySize},
	}

	results := make([]*string, querySize*querySize)

	for i := 0; b.Loop(); i++ {
		_ = grid.QueryRange(query, results)
	}
}

func randCoord(r *rand.Rand, max uint32) uint32 {
	return uint32(r.Intn(int(max + 1)))
}
