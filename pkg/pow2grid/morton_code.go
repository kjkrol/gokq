package pow2grid

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

func (c MortonCode) Offset(dx, dy int32) MortonCode {
	x, y := c.Decode()

	// tu możesz dodać obsługę zakresów / clamp / wrapowanie
	nx := uint32(int64(x) + int64(dx))
	ny := uint32(int64(y) + int64(dy))

	return NewMortonCode(nx, ny)
}

func MortonCodeArea(aabb AABB) []MortonCode {
	// Pusta / niepoprawna AABB
	if aabb.Max.X < aabb.Min.X || aabb.Max.Y < aabb.Min.Y {
		return nil
	}

	width := aabb.Max.X - aabb.Min.X + 1
	height := aabb.Max.Y - aabb.Min.Y + 1

	// Ile elementów w sumie
	count := uint64(width) * uint64(height)
	if count == 0 {
		return nil
	}

	res := make([]MortonCode, int(count))

	// Kod lewego-górnego rogu (Min.X, Min.Y)
	rowStart := NewMortonCode(aabb.Min.X, aabb.Min.Y)

	idx := 0
	for range height {
		code := rowStart

		for range width {
			res[idx] = code
			idx++

			code = code.IncX() // (x+1, y)
		}

		// Następny wiersz: (Min.X, y+1)
		rowStart = rowStart.IncY()
	}

	return res
}

func MortonCodeAreaConsume(aabb AABB, fn func(int, MortonCode)) {
	// Pusta / niepoprawna AABB
	if aabb.Max.X < aabb.Min.X || aabb.Max.Y < aabb.Min.Y {
		return
	}

	width := aabb.Max.X - aabb.Min.X + 1
	height := aabb.Max.Y - aabb.Min.Y + 1

	// Ile elementów w sumie
	count := uint64(width) * uint64(height)
	if count == 0 {
		return
	}

	// Kod lewego-górnego rogu (Min.X, Min.Y)
	rowStart := NewMortonCode(aabb.Min.X, aabb.Min.Y)

	idx := 0
	for range height {
		code := rowStart

		for range width {
			fn(idx, code)
			idx++
			code = code.IncX()
		}

		rowStart = rowStart.IncY()
	}
}

const (
	xMask uint64 = 0x5555555555555555 // bity X na pozycjach 0,2,4,...
	yMask uint64 = 0xAAAAAAAAAAAAAAAA // bity Y na pozycjach 1,3,5,...
)

func (c MortonCode) IncX() MortonCode {
	code := uint64(c)
	x := code & xMask
	y := code & yMask

	// splitBy1(x+1) == (splitBy1(x) - xMask) & xMask
	x = (x - xMask) & xMask

	return MortonCode(x | y)
}

func (c MortonCode) IncY() MortonCode {
	code := uint64(c)
	x := code & xMask
	y := (code & yMask) >> 1 // teraz Y ma bity jak splitBy1

	y = (y - xMask) & xMask // y+1 w splitBy1-space
	y = y << 1              // wracamy na nieparzyste pozycje

	return MortonCode(x | y)
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
