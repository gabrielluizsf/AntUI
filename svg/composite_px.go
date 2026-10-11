package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// compositePx is one pixel of feComposite: in is the source — `in` — and out
// the destination — `in2` — with the operations the filter effects spec
// names, and the arithmetic the same on every channel including the alpha.
// The two are held against their own alpha while the operation runs and come
// out of it against theirs, which is what the spec's formulas write.
func compositePx(in, out canvas.Color, op compositeOp, k [4]float64) canvas.Color {
	if op == compOver {
		return canvas.BlendOver(out, in)
	}
	qa, qb := frac(in.A()), frac(out.A())
	ca := [3]float64{frac(in.R()) * qa, frac(in.G()) * qa, frac(in.B()) * qa}
	cb := [3]float64{frac(out.R()) * qb, frac(out.G()) * qb, frac(out.B()) * qb}
	var cr [3]float64
	var qr float64
	switch op {
	case compIn:
		qr = qa * qb
		for i := range cr {
			cr[i] = ca[i] * qb
		}
	case compOut:
		qr = qa * (1 - qb)
		for i := range cr {
			cr[i] = ca[i] * (1 - qb)
		}
	case compAtop:
		qr = qb
		for i := range cr {
			cr[i] = ca[i]*qb + cb[i]*(1-qa)
		}
	case compXor:
		qr = qa*(1-qb) + qb*(1-qa)
		for i := range cr {
			cr[i] = ca[i]*(1-qb) + cb[i]*(1-qa)
		}
	case compLighter:
		qr = math.Min(1, qa+qb)
		for i := range cr {
			cr[i] = math.Min(1, ca[i]+cb[i])
		}
	case compArithmetic:
		qr = clamp01(k[0]*qa*qb + k[1]*qa + k[2]*qb + k[3])
		for i := range cr {
			cr[i] = clamp01(k[0]*ca[i]*cb[i] + k[1]*ca[i] + k[2]*cb[i] + k[3])
		}
	}
	return unPremultiply(cr, qr)
}
