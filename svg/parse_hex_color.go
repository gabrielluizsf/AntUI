package svg

import "github.com/gabrielluizsf/antui/canvas"

// parseHexColor reads `#rgb`, `#rgba`, `#rrggbb` and `#rrggbbaa`. The three and
// four digit forms repeat each digit, which is how `#f0a` is the same colour as
// `#ff00aa` and not a shorter way of writing something else.
func parseHexColor(body string) (canvas.Color, bool) {
	digit := func(c byte) (uint8, bool) {
		switch {
		case c >= '0' && c <= '9':
			return c - '0', true
		case c >= 'a' && c <= 'f':
			return c - 'a' + 10, true
		case c >= 'A' && c <= 'F':
			return c - 'A' + 10, true
		}
		return 0, false
	}
	// read takes a run of digits and widens it to the eight of a colour, each
	// digit standing for itself and each byte for both of its halves.
	read := func(s string) (r, g, b, a uint8, ok bool) {
		switch len(s) {
		case 3, 4:
			v := make([]uint8, 0, 4)
			for i := range len(s) {
				d, good := digit(s[i])
				if !good {
					return 0, 0, 0, 0, false
				}
				v = append(v, d*17) // #f is #ff
			}
			r, g, b = v[0], v[1], v[2]
			a = 0xFF
			if len(v) == 4 {
				a = v[3]
			}
			return r, g, b, a, true
		case 6, 8:
			v := make([]uint8, 0, 4)
			for i := 0; i < len(s); i += 2 {
				hi, ok1 := digit(s[i])
				lo, ok2 := digit(s[i+1])
				if !ok1 || !ok2 {
					return 0, 0, 0, 0, false
				}
				v = append(v, hi<<4|lo)
			}
			r, g, b = v[0], v[1], v[2]
			a = 0xFF
			if len(v) == 4 {
				a = v[3]
			}
			return r, g, b, a, true
		}
		return 0, 0, 0, 0, false
	}
	r, g, b, a, ok := read(body)
	if !ok {
		return 0, false
	}
	return canvas.RGBA(r, g, b, a), true
}
