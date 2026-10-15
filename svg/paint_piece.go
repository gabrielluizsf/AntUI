package svg

import "github.com/gabrielluizsf/antui/canvas"

// paintPiece puts one run of writing down: its stroke first and its fill over
// the stroke, which is the order SVG paints and the order that leaves the half
// of a stroke inside the letters covered by the fill and the half outside
// showing as the outline it is.
func paintPiece(cv *canvas.Canvas, n *Node, p *piece, m canvas.Matrix,
	paints palette, patterning []*pattern) {
	st := p.run.Style
	if st.HasStroke && st.Width > 0 {
		c := st.Stroke
		if st.StrokeCurrent {
			c = paints.current
		}
		if c := fade(c, st.Opacity*st.StrokeOpacity); c.A() > 0 {
			strokeWriting(cv, p.put, st.strokeWidth(m), c, p.scale)
		}
	}
	shade, c, ok := fillShade(n, st, m, paints, cv.Width, cv.Height, patterning)
	switch {
	case !ok:
		return
	case shade == nil:
		if c.A() > 0 {
			p.put(cv, c, 0, 0)
		}
	default:
		layShaded(cv, p.put, shade)
	}
}
