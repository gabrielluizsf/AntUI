package canvas

import (
	"unsafe"

	"github.com/gabrielluizsf/antui/blendasm"
)

// hasSIMD is whether the assembly exists for this build.
const hasSIMD = blendasm.Have

// pixels is a run of colours as the plain numbers the assembly takes.
func pixels(run []Color) *uint32 {
	return (*uint32)(unsafe.Pointer(&run[0]))
}

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

// blendRowGo is the reference implementation, and the one that runs where
// there is no assembly for the architecture.
func blendRowGo(dst, src []Color) {
	for i, c := range src {
		switch c.A() {
		case 0:
			// Nothing to add, and dst keeps its own alpha.
		case 255:
			dst[i] = c
		default:
			dst[i] = Blend(dst[i], c)
		}
	}
}

// BlendRowSolid composites one colour over a run of pixels. An opaque colour
// is not a blend at all — the run is written rather than read, multiplied and
// written back — which is the case a cleared background takes.
func BlendRowSolid(dst []Color, c Color) {
	switch c.A() {
	case 0:
		return
	case 255:
		solid := c | 0xFF000000
		for i := range dst {
			dst[i] = solid
		}
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

func blendRowSolidGo(dst []Color, c Color) {
	for i := range dst {
		dst[i] = Blend(dst[i], c)
	}
}
