package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// transformMatrix is the whole transform list pinned to the box and its
// origin: the user's functions run around the origin point, which itself sits
// at the box's top-left plus transform-origin. An undeclared origin is the
// CSS initial "50% 50%", the centre.
func (cs *cssStyle) transformMatrix(st css.Style, x, y, w, h int) canvas.Matrix {
	ox, oy := cs.originPoint(st, x, y, w, h)

	m := canvas.Translate(ox, oy)
	for i := range st.Transform {
		f := st.Transform[i]
		switch f.Kind {
		case css.TransformTranslate:
			m = m.Mul(canvas.Translate(float64(cs.length(st, f.Dx, w)), float64(cs.length(st, f.Dy, h))))
		case css.TransformScale:
			m = m.Mul(canvas.Scale(f.Sx, f.Sy))
		case css.TransformRotate:
			m = m.Mul(canvas.Rotate(f.Ax.Rad()))
		case css.TransformSkew:
			m = m.Mul(canvas.Skew(f.Ax.Rad(), f.Ay.Rad()))
		case css.TransformMatrix:
			m = m.Mul(canvas.Matrix{A: f.M[0], B: f.M[1], C: f.M[2], D: f.M[3], E: f.M[4], F: f.M[5]})
		}
	}
	return m.Mul(canvas.Translate(-ox, -oy))
}
