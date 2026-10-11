package font

// Glyph returns the 16 rows for a codepoint, falling back to '?' for
// anything the font does not carry.
func Glyph(r rune) []byte {
	var index int
	switch {
	case r >= ASCIILo && r <= ASCIIHi:
		index = int(r) - ASCIILo
	case r >= LatinLo && r <= LatinHi:
		index = (ASCIIHi - ASCIILo + 1) + int(r) - LatinLo
	default:
		index = '?' - ASCIILo
	}
	return font8x16[index*Height : index*Height+Height]
}
