package svg

// hexBytes turns the body of a hex colour into channels, widening the three and
// four digit forms to the eight of a colour: each digit stands for itself and
// each byte for both of its halves, which is why `#f` is `#ff`.
func hexBytes(s string) (r, g, b, a uint8, ok bool) {
	switch len(s) {
	case 3, 4:
		v := make([]uint8, 0, 4)
		for i := range len(s) {
			d, good := hexDigit(s[i])
			if !good {
				return 0, 0, 0, 0, false
			}
			v = append(v, d*17) // #f is #ff
		}
		r, g, b, a = v[0], v[1], v[2], 0xFF
		if len(v) == 4 {
			a = v[3]
		}
		return r, g, b, a, true
	case 6, 8:
		v := make([]uint8, 0, 4)
		for i := 0; i < len(s); i += 2 {
			hi, ok1 := hexDigit(s[i])
			lo, ok2 := hexDigit(s[i+1])
			if !ok1 || !ok2 {
				return 0, 0, 0, 0, false
			}
			v = append(v, hi<<4|lo)
		}
		r, g, b, a = v[0], v[1], v[2], 0xFF
		if len(v) == 4 {
			a = v[3]
		}
		return r, g, b, a, true
	}
	return 0, 0, 0, 0, false
}
