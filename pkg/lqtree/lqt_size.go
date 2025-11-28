package lqtree

// -----------------------------------------
// Enum where the name reflects total cell count (side length squared).
// -----------------------------------------
type LQTSize uint8

const (
	Size1 LQTSize = iota
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

// Resolution returns the max coordinate on each axis (side length - 1), e.g. Size1024 (32x32) → 31.
func (s LQTSize) Resolution() uint32 {
	return (1 << s) - 1
}

func (s LQTSize) ArraySize() uint64 {
	return 1 << (2 * s)
}
