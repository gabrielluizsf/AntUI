package svg

import "unicode/utf16"

// utf16ToUTF8 turns UTF-16 text — its byte-order mark, then one unit per two
// bytes — into the UTF-8 the parser reads. The mark says which way round the
// bytes come, and a unit in the surrogate range pairs with the one after it to
// stand for a single character beyond the basic plane.
func utf16ToUTF8(src string) string {
	le := src[0] == 0xFF
	src = src[2:]
	units := make([]uint16, 0, len(src)/2)
	for i := 0; i+1 < len(src); i += 2 {
		if le {
			units = append(units, uint16(src[i])|uint16(src[i+1])<<8)
		} else {
			units = append(units, uint16(src[i+1])|uint16(src[i])<<8)
		}
	}
	return string(utf16.Decode(units))
}
