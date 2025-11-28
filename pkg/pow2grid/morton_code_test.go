package pow2grid

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestMortonCodeOffset(t *testing.T) {
	// TODO: cover boundary/wrap behaviour once Offset defines how to handle edges.
	cases := []struct {
		name     string
		x, y     uint32
		dx, dy   int32
		expected Pos
	}{
		{
			name:     "shift_3_4",
			x:        10,
			y:        20,
			dx:       3,
			dy:       4,
			expected: Pos{X: 13, Y: 24},
		},
	}

	src := rand.New(rand.NewSource(1))
	for i := range 10 {
		x := uint32(src.Intn(1 << 20)) // stay far from edges
		y := uint32(src.Intn(1 << 20)) // stay far from edges
		dx := int32(src.Intn(6) + 1)   // positive shifts to avoid underflow
		dy := int32(src.Intn(6) + 1)   // positive shifts to avoid underflow
		cases = append(cases, struct {
			name     string
			x, y     uint32
			dx, dy   int32
			expected Pos
		}{
			name:     fmt.Sprintf("rand_%d", i),
			x:        x,
			y:        y,
			dx:       dx,
			dy:       dy,
			expected: Pos{X: x + uint32(dx), Y: y + uint32(dy)},
		})
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			code := NewMortonCode(tc.x, tc.y)
			shifted := code.Offset(tc.dx, tc.dy)
			gotX, gotY := shifted.Decode()

			if gotX != tc.expected.X || gotY != tc.expected.Y {
				t.Fatalf("Offset(%d,%d) from (%d,%d): got (%d,%d), want (%d,%d)",
					tc.dx, tc.dy, tc.x, tc.y, gotX, gotY, tc.expected.X, tc.expected.Y)
			}
		})
	}
}

func TestMortonCodeArea(t *testing.T) {
	// TODO: add clamping/wrapping coverage when MortonCodeArea behaviour at edges is defined.
	aabb := AABB{
		Min: Pos{X: 5, Y: 7},
		Max: Pos{X: 7, Y: 10}, // 3x4 rectangle (inclusive)
	}

	codes := MortonCodeArea(aabb)

	wantLen := 12
	if len(codes) != wantLen {
		t.Fatalf("len(MortonCodeArea) = %d, want %d", len(codes), wantLen)
	}

	idx := 0
	for y := aabb.Min.Y; y <= aabb.Max.Y; y++ {
		for x := aabb.Min.X; x <= aabb.Max.X; x++ {
			want := NewMortonCode(x, y)
			if codes[idx] != want {
				t.Fatalf("codes[%d] = %v for (%d,%d), want %v", idx, codes[idx], x, y, want)
			}
			idx++
		}
	}
}
