package template

import "github.com/gabrielluizsf/antui/template/css"

// originPoint is the pivot of the transform in device pixels: the box's
// top-left plus transform-origin, the initial value landing dead centre.
func (cs *cssStyle) originPoint(st css.Style, x, y, w, h int) (float64, float64) {
	if !st.Has("transform-origin") {
		return float64(x + w/2), float64(y + h/2)
	}
	return float64(x + cs.length(st, st.TransformOrigin[0], w)),
		float64(y + cs.length(st, st.TransformOrigin[1], h))
}
