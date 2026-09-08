package font

// Cell size of the built-in font, in pixels.
const (
	Width  = 8
	Height = 16
)

// TextWidth is the 8x16 font's own answer, which is arithmetic.
func TextWidth(text string) int {
	width, line := 0, 0
	for _, r := range text {
		switch r {
		case '\n':
			width = max(width, line)
			line = 0
		case '\t':
			line += Width * 4
		default:
			line += Width
		}
	}
	return max(width, line)
}

const (
	ASCIILo = 32
	ASCIIHi = 126
	LatinLo = 160
	LatinHi = 255
	Glyphs  = 191
)
