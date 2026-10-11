package svg

import "github.com/gabrielluizsf/antui/canvas"

// rgbaColor is four numbers between zero and one as a pixel's colour.
func rgbaColor(r, g, b, a float64) canvas.Color {
	return canvas.Color(uint32(a*255+0.5)<<24 |
		uint32(r*255+0.5)<<16 |
		uint32(g*255+0.5)<<8 |
		uint32(b*255+0.5))
}
