package canvas

import "github.com/gabrielluizsf/antui/blendasm"

// BlendRow composites src over dst, pixel by pixel, honouring the source
// alpha. It blends min(len(dst), len(src)) pixels.
func BlendRow(dst, src []Color) {
	n := min(len(dst), len(src))
	if n == 0 {
		return
	}
	i := 0
	if hasSIMD && n >= 4 {
		whole := n &^ 3
		blendasm.Rows(pixels(dst), pixels(src), whole)
		i = whole
	}
	blendRowGo(dst[i:n], src[i:n])
}
