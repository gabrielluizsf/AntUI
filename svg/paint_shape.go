package svg

import "github.com/gabrielluizsf/antui/canvas"

// paintShape draws one node's shape: its fill first and its stroke over it,
// which is the order SVG paints, so a stroke half its width over a fill covers
// the fill's edge as it should, and the marker over the pair of them last of
// all — an arrowhead on the line rather than under it.
func paintShape(cv *canvas.Canvas, n *Node, m canvas.Matrix, current canvas.Color, masking []*maskDef, patterning []*pattern) {
	st := n.Style
	if st.HasFill {
		switch {
		case st.fillGradient != nil:
			// A gradient is a colour at every point rather than one colour, so
			// the fill is laid one pixel at a time. The path is still transformed
			// first: the gradient moves with the shape it paints, and what is
			// asked about each pixel is where that pixel falls in the shape's own
			// coordinates rather than where it falls on the canvas.
			fill := clonePath(n.Path)
			fill.Transform(m)
			cv.FillPathFunc(fill, st.FillRule, gradientShade(st.fillGradient, n, m, current, st.Opacity*st.FillOpacity))
		case st.fillPattern != nil:
			// A pattern is a picture at every point rather than one colour, so
			// the fill is laid one pixel at a time the same way. Where the
			// pattern has nothing to paint with — it names itself, the tile has
			// no size to it, the transform has no way back — the colour after
			// the `url(...)` is what the shape falls back on, and nothing at
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
				// The path is drawn in its own coordinates and the transform is
				// put on a copy of it, so the node keeps the shape the drawing
				// wrote while the picture comes out where the transform put it.
				// The stroke takes its width from the same transform rather
				// than from the geometry — see [Style.strokeWidth].
				fill := clonePath(n.Path)
				fill.Transform(m)
				cv.FillPath(fill, c, st.FillRule)
			}
		}
	}
	if st.HasStroke && st.Width > 0 {
		paintStroke(cv, n, m, current, patterning)
	}
	// The marker comes after the fill and the stroke, which is the order the
	// spec puts them in. It is drawn where the shape has no stroke at all, and
	// where it encloses no fill either: a marker is a picture put on a vertex,
	// and the line beside it is not what puts it there — the vertex is.
	if len(n.markers) > 0 {
		paintMarkers(cv, n, m, current, masking, patterning)
	}
}
