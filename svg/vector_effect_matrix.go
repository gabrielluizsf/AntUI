package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// vectorEffectMatrix is the matrix a node and everything inside it are painted
// through when the element wrote a `vector-effect` that holds its own user
// coordinate system still. It is the matrix [paintNodes] carries in place of
// the drawing's, so a painter that puts its geometry through m·Style.Transform
// lands where the effect says and a shape inside the element is left in the
// coordinate system the effect made of it.
//
// The transform from the element's coordinates to the host — the viewBox fit
// composed with everything the element inherited and wrote — is the spec's CTM,
// and the spec's formulas rewrite it: non-scaling-size divides its linear part
// by the square root of its determinant, which leaves the rotation and the skew
// and makes the scale one; non-rotation takes that square root alone, keeping
// the scale and dropping rotation and skew; the two together are the identity;
// fixed-position spends the shift the element's own transform writes, so the
// element keeps the host's rotation and scale but not its move.
//
// The determinant and the shift are in the pixels the canvas is painted at,
// supersample times the pixels the drawing comes out at, so the scale that is
// dropped has the supersampling put back: a shape held at scale one is held at
// one final pixel per unit and not at one of the bigger ones. A host transform
// with no determinant, or one with no way back, leaves the element as it is —
// the spec's own reading for a determinant of zero — rather than making an
// infinity of it.
func vectorEffectMatrix(m canvas.Matrix, st Style) canvas.Matrix {
	at := m.Mul(st.Transform)
	a, b, c, d, e, f := at.A, at.B, at.C, at.D, at.E, at.F
	if st.VectorEffect.FixedPosition {
		e, f = st.VectorEffect.OriginX*supersample, st.VectorEffect.OriginY*supersample
	}
	det := a*d - b*c
	if det == 0 {
		return m
	}
	back, ok := st.Transform.Inverse()
	if !ok {
		return m
	}
	k := math.Sqrt(math.Abs(det))
	switch {
	case st.VectorEffect.NonScalingSize && st.VectorEffect.NonRotation:
		a, b, c, d = supersample, 0, 0, supersample
	case st.VectorEffect.NonScalingSize:
		a, b, c, d = a*supersample/k, b*supersample/k, c*supersample/k, d*supersample/k
	case st.VectorEffect.NonRotation:
		a, b, c, d = k, 0, 0, k
	}
	return canvas.Matrix{A: a, B: b, C: c, D: d, E: e, F: f}.Mul(back)
}
