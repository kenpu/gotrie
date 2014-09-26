package gotrie

func Uint64_string(x uint64) string {
	var s string
	var i uint
	for i = 0; i < 64; i++ {
		if i > 0 && i%4 == 0 {
			s += " "
		}
		if x&(1<<i) == 0 {
			s += "0"
		} else {
			s += "1"
		}
	}
	return s
}

var index64 [64]uint8 = [64]uint8{
	0, 47, 1, 56, 48, 27, 2, 60,
	57, 49, 41, 37, 28, 16, 3, 61,
	54, 58, 35, 52, 50, 42, 21, 44,
	38, 32, 29, 23, 17, 11, 4, 62,
	46, 55, 26, 59, 40, 36, 15, 53,
	34, 51, 20, 43, 31, 22, 10, 45,
	25, 39, 14, 33, 19, 30, 9, 24,
	13, 18, 8, 12, 7, 6, 5, 63,
}

const debruijn64 uint64 = 0x03f79d71b4cb0a89

func LeadZeros(bb uint64) uint8 {
	return index64[((bb^(bb-1))*debruijn64)>>58]
}

func TestBit(x uint64, pos uint8) bool {
	return ((x>>pos)&1 != 0)
}
func TestBit_Int(x uint64, pos uint8) int {
	return int((x >> pos) & 1)
}

func CommonPrefixLength(x, y uint64) uint8 {
	return LeadZeros(x ^ y)
}

func PopCount(i uint64) uint8 {
	i = i - ((i >> 1) & 0x5555555555555555)
	i = (i & 0x3333333333333333) + ((i >> 2) & 0x3333333333333333)
	i = (i + (i >> 4)) & 0xF0F0F0F0F0F0F0F
	return uint8((i * 0x101010101010101) >> 56)
}

const One uint64 = ^uint64(0)

func PopCountPartial(x uint64, pos uint8) uint8 {
	return PopCount(x << (64 - pos))
}
