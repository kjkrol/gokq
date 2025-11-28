package pow2grid

// -----------------------------------------
// Enum where the name reflects total cell count (side length squared).
// -----------------------------------------
type Resolution uint8

const (
	Size1x1 Resolution = iota
	Size2x2
	Size4x4
	Size8x8
	Size16x16
	Size32x32
	Size64x64
	Size128x128
	Size256x256
	Size512x512
	Size1024x1024
)

// Side returns the side length of the square grid, e.g. Size8x8 → 8.
func (s Resolution) Side() uint32 {
	return uint32(1) << s
}

// MaxCoord returns the maximum coordinate on each axis, i.e. side - 1.
func (s Resolution) MaxCoord() uint32 {
	return (uint32(1) << s) - 1
}

// Cells returns the total number of cells in the grid (side * side).
// To jest 1 << (2*s).
func (s Resolution) Cells() uint64 {
	return uint64(1) << (s << 1)
}
