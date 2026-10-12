package svg

import "github.com/gabrielluizsf/antui/canvas"

// paintFill lays a node's fill: one colour, a gradient or a pattern. The path
// is drawn in its own coordinates and the transform is put on a copy of it, so
// the node keeps the shape the drawing wrote while the picture comes out where
// the transform put it. The stroke takes its width from the same transform
// rather than from the geometry — see [Style.strokeWidth].
func paintFill(cv *canvas.Canvas, n *Node, m canvas.Matrix, current canvas.Color, st Style, patterning []*pattern) {
	switch {
	case st.fillGradient != nil:
		// A gradient is a colour at every point rather than one colour, so the
		// fill is laid one pixel at a time from the gradient's own shade.
		fill := clonePath(n.Path)
		fill.Transform(m)
		cv.FillPathFunc(fill, st.FillRule, gradientShade(st.fillGradient, n, m, current, st.Opacity*st.FillOpacity))
	case st.fillPattern != nil:
		// A pattern is a picture at every point rather than one colour, laid
		// the same way. Where the pattern has nothing to paint with, the colour
		// after the `url(...)` is what the shape falls back on, and nothing at
		// all where it brought none.
		fill := clonePath(n.Path)
		fill.Transform(m)
		if shade, ok := patternShade(st.fillPattern, n, m, cv.Width, cv.Height, current, st.Opacity*st.FillOpacity, patterning); ok {
			cv.FillPathFunc(fill, st.FillRule, shade)
		} else if c, ok := st.fillFallbackColour(current); ok {
			if c := fade(c, st.Opacity*st.FillOpacity); c.A() > 0 {
				cv.FillPath(fill, c, st.FillRule)
			}
		}
	default:
		c := st.Fill
		if st.FillCurrent {
			c = current
		}
		if c := fade(c, st.Opacity*st.FillOpacity); c.A() > 0 {
			fill := clonePath(n.Path)
			fill.Transform(m)
			cv.FillPath(fill, c, st.FillRule)
		}
	}
}
