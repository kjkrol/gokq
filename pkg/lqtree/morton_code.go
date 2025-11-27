package lqtree

type MortonCode uint64

func NewMortonCode(x, y uint32) MortonCode {
	return MortonCode((splitBy1(x)) | splitBy1(y)<<1)
}

func (c MortonCode) Decode() (x, y uint32) {
	code := uint64(c)
	x = compact1By1(code)
	y = compact1By1(code >> 1)
	return
}

// splitBy1 spreads the lower 32 bits of input so that there is 1 zero bit between each bit.
func splitBy1(a uint32) uint64 {
	x := uint64(a)
	x = (x | (x << 16)) & 0x0000FFFF0000FFFF
	x = (x | (x << 8)) & 0x00FF00FF00FF00FF
	x = (x | (x << 4)) & 0x0F0F0F0F0F0F0F0F
	x = (x | (x << 2)) & 0x3333333333333333
	x = (x | (x << 1)) & 0x5555555555555555
	return x
}

func compact1By1(a uint64) uint32 {
	x := a & 0x5555555555555555
	x = (x | (x >> 1)) & 0x3333333333333333
	x = (x | (x >> 2)) & 0x0F0F0F0F0F0F0F0F
	x = (x | (x >> 4)) & 0x00FF00FF00FF00FF
	x = (x | (x >> 8)) & 0x0000FFFF0000FFFF
	x = (x | (x >> 16)) & 0x00000000FFFFFFFF
	return uint32(x)
}
