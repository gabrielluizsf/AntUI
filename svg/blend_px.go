package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// blendPx is one pixel of feBlend, with the formulas the filter effects spec
// writes: A is `in` — the top of the two — and B is `in2`, the bottom. The
// colour is held against each picture's own alpha while the mode runs and
// comes out against the alpha of what the mode gave back.
func blendPx(a, b canvas.Color, mode blendMode) canvas.Color {
	if mode == blendNormal {
		return canvas.BlendOver(b, a)
	}
	qa, qb := frac(a.A()), frac(b.A())
	ca := [3]float64{frac(a.R()) * qa, frac(a.G()) * qa, frac(a.B()) * qa}
	cb := [3]float64{frac(b.R()) * qb, frac(b.G()) * qb, frac(b.B()) * qb}
	var cr [3]float64
	switch mode {
	case blendMultiply:
		for i := range cr {
			cr[i] = (1-qa)*cb[i] + (1-qb)*ca[i] + ca[i]*cb[i]
		}
	case blendScreen:
		for i := range cr {
			cr[i] = cb[i] + ca[i] - ca[i]*cb[i]
		}
	case blendDarken:
		for i := range cr {
			cr[i] = math.Min((1-qa)*cb[i]+ca[i], (1-qb)*ca[i]+cb[i])
		}
	case blendLighten:
		for i := range cr {
			cr[i] = math.Max((1-qa)*cb[i]+ca[i], (1-qb)*ca[i]+cb[i])
		}
	}
	return unPremultiply(cr, 1-(1-qa)*(1-qb))
}
