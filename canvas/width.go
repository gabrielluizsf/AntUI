package canvas

import (
	"math"

	"github.com/gabrielluizsf/antui/font"
)

// Width is how wide a string is, in pixels, measuring the widest line.
func (f *Face) Width(text string) int {
	if f == nil || f.file == nil {
		return int(math.Round(float64(font.TextWidth(text)) * f.scaleOf()))
	}
	widest, line := 0, 0.0
	for _, r := range text {
		switch r {
		case '\n':
			widest = max(widest, int(line+0.5))
			line = 0
		case '\r':
		case '\t':
			line += f.advance(' ') * 4
		default:
			line += f.advance(r)
		}
	}
	return max(widest, int(line+0.5))
}
