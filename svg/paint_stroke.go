package svg

import "github.com/gabrielluizsf/antui/canvas"

// paintStroke widens the outline of one node's shape and paints what comes of
// it. It is the whole of the stroke rather than a part of the shape, so that
// the marker can be drawn after all of it rather than after only some.
func paintStroke(cv *canvas.Canvas, n *Node, m canvas.Matrix, current canvas.Color, patterning []*pattern) {
	st := n.Style
	// A dashed stroke is cut into pieces along the line before it is widened,
	// so that each dash is capped at both ends and a dash never starts halfway
	// through a corner. The cut was made where the shape was written and the
	// transform of the element has carried it along with the shape ever since,
	// so the outline taken here is already in the runs the pattern asked for,
	// and taking it to the canvas carries the cut the rest of the way — the
	// viewBox and the size asked for included.
	outline := n.Path
	if n.dashed != nil {
		outline = n.dashed
	}
	band := strokeBand(outline, st, m)
	if st.strokeGradient != nil {
		// The stroke is widened into the shape it covers, and that shape is
		// filled with the gradient, so the colour runs across the stroke the same
		// way it runs across the fill rather than being one colour of it.
		cv.FillPathFunc(band, canvas.FillNonZero, gradientShade(st.strokeGradient, n, m, current, st.Opacity*st.StrokeOpacity))
		return
	}
	if st.strokePattern != nil {
		// The stroke is widened into the shape it covers, and that shape is
		// painted with the pattern, so the picture runs across the stroke the
		// same way it runs across the fill — with the colour after the
		// `url(...)` where the pattern has nothing to paint with.
		if shade, ok := patternShade(st.strokePattern, n, m, cv.Width, cv.Height, current, st.Opacity*st.StrokeOpacity, patterning); ok {
			cv.FillPathFunc(band, canvas.FillNonZero, shade)
		} else if c, ok := st.strokeFallbackColour(current); ok {
			if c := fade(c, st.Opacity*st.StrokeOpacity); c.A() > 0 {
				cv.FillPath(band, c, canvas.FillNonZero)
			}
		}
		return
	}
	c := st.Stroke
	if st.StrokeCurrent {
		c = current
	}
	c = fade(c, st.Opacity*st.StrokeOpacity)
	if c.A() == 0 {
		return
	}
	cv.FillPath(band, c, canvas.FillNonZero)
}
