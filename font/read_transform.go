package font

// readTransform reads the 2x2 transform a component carries, which is the
// identity unless its flags name one of the three scale forms.
func readTransform(data []byte, p int, flags uint16) (a, b, c, d float64, next int, ok bool) {
	a, b, c, d = 1.0, 0.0, 0.0, 1.0
	switch {
	case flags&0x0008 != 0: // WE_HAVE_A_SCALE
		if p+2 > len(data) {
			return 0, 0, 0, 0, p, false
		}
		a = f2dot14(be16(data, p))
		d = a
		p += 2
	case flags&0x0040 != 0: // X_AND_Y_SCALE
		if p+4 > len(data) {
			return 0, 0, 0, 0, p, false
		}
		a, d = f2dot14(be16(data, p)), f2dot14(be16(data, p+2))
		p += 4
	case flags&0x0080 != 0: // TWO_BY_TWO
		if p+8 > len(data) {
			return 0, 0, 0, 0, p, false
		}
		a, b = f2dot14(be16(data, p)), f2dot14(be16(data, p+2))
		c, d = f2dot14(be16(data, p+4)), f2dot14(be16(data, p+6))
		p += 8
	}
	return a, b, c, d, p, true
}
