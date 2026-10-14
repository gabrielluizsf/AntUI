package svg

import "github.com/gabrielluizsf/antui/canvas"

// strokeWidth is how wide the stroke comes out on the canvas: the width the
// drawing wrote, taken through everything that scaled the shape — the viewBox
// onto the canvas and the transforms the shape went through, which are the same
// ones that moved the geometry, so that a shape scaled twice is a stroke twice
// as wide. That is the default SVG has.
//
// `vector-effect="non-scaling-stroke"` asks for the other one: the width it
// says, in the pixels the drawing comes out at, with nothing scaled at all —
// neither the viewBox nor any transform, which is what keeps a hairline a
// hairline when the drawing is zoomed. The canvas is painted supersample times
// the size asked for and brought back down at the end, so the width is in those
// bigger pixels and comes down with the rest of the picture.
func (st Style) strokeWidth(m canvas.Matrix) float64 {
	if st.VectorEffect.NonScalingStroke {
		return st.Width * supersample
	}
	// The transform only counts where it was put into the path, under the same
	// condition the build used to put it there: a style with no transform of
	// its own scales nothing, and a zero matrix would scale the width to none.
	if st.HasTransform && st.Transform != (canvas.Matrix{}) {
		m = m.Mul(st.Transform)
	}
	return st.Width * scaleOf(m)
}
