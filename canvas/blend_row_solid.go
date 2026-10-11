package canvas

import "github.com/gabrielluizsf/antui/blendasm"

// BlendRowSolid composites one colour over a run of pixels. An opaque colour
// is not a blend at all — the run is written rather than read, multiplied and
// written back — which is the case a cleared background takes.
func BlendRowSolid(dst []Color, c Color) {
	if c.A() == 255 {
		solidRunGo(dst, c)
		return
	}
	if c.A() == 0 {
		return
	}

	i := 0
	if hasSIMD && len(dst) >= 4 {
		whole := len(dst) &^ 3
		blendasm.Solid(pixels(dst), uint32(c), whole)
		i = whole
	}
	blendRowSolidGo(dst[i:], c)
}
