package lqtree

// -----------------------------------------
// Enum where the name reflects total cell count (side length squared).
// -----------------------------------------
type LQTSize uint8

const (
	Size1       LQTSize = iota // 1x1
	Size4                      // 2x2
	Size16                     // 4x4
	Size64                     // 8x8
	Size256                    // 16x16
	Size1024                   // 32x32
	Size4096                   // 64x64
	Size16384                  // 128x128
	Size65536                  // 256x256
	Size262144                 // 512x512
	Size1048576                // 1024x1024
)

// Resolution returns the max coordinate on each axis (side length - 1), e.g. Size1024 (32x32) → 31.
func (s LQTSize) Resolution() uint32 {
	return (1 << s) - 1
}

func (s LQTSize) ArraySize() uint64 {
	return 1 << (2 * s)
}
