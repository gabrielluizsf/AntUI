package canvas

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// isqrt64 is an integer square root. The C original avoided it to stay clear
// of math.h and the -lm flag; here it survives because the rasterizers want
// an exact integer answer for coverage tests, and a float sqrt rounded back
// to an integer is not reliably the same number on every machine.
func isqrt64(value uint64) uint32 {
	var rem, root uint64
	for i := 0; i < 32; i++ {
		root <<= 1
		rem = (rem << 2) | (value >> 62)
		value <<= 2
		if root < rem {
			rem -= root + 1
			root += 2
		}
	}
	return uint32(root >> 1)
}
