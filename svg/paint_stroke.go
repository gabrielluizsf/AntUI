package svg

import "github.com/gabrielluizsf/antui/canvas"

// paintStroke widens the outline of one node's shape and paints what comes of
// it. It is the whole of the stroke rather than a part of the shape, so that
// the marker can be drawn after all of it rather than after only some.
func paintStroke(cv *canvas.Canvas, n *Node, m canvas.Matrix, paints palette, patterning []*pattern) {
	st := n.Style
	// A regular dashed stroke is cut where the shape was written, so the
	// element transform carries its runs along. A non-scaling stroke is cut by
	// nonScalingDashed after that transform, because its pattern is in pixels.
	outline := n.Path
	outlineMatrix := m
	if st.VectorEffect.NonScalingStroke {
		if fixed := nonScalingDashed(n, m); fixed != nil {
			outline = fixed
			outlineMatrix = canvas.Identity()
		}
	}
	if n.dashed != nil && !st.VectorEffect.NonScalingStroke {
		outline = n.dashed
	}
	band := strokeBand(outline, st, outlineMatrix)
	if st.strokeGradient != nil {
		// The stroke is widened into the shape it covers, and that shape is
		// filled with the gradient, so the colour runs across the stroke the same
		// way it runs across the fill rather than being one colour of it.
		cv.FillPathFunc(band, canvas.FillNonZero, gradientShade(st.strokeGradient, n, m, paints.current, st.Opacity*st.StrokeOpacity))
		return
	}
	if st.strokePattern != nil {
		// The stroke is widened into the shape it covers, and that shape is
		// painted with the pattern, so the picture runs across the stroke the
		// same way it runs across the fill — with the colour after the
		// `url(...)` where the pattern has nothing to paint with.
		if shade, ok := patternShade(st.strokePattern, n, m, cv.Width, cv.Height, paints.current, st.Opacity*st.StrokeOpacity, patterning); ok {
			cv.FillPathFunc(band, canvas.FillNonZero, shade)
		} else if c, ok := st.strokeFallbackColour(paints.current); ok {
			if c := fade(c, st.Opacity*st.StrokeOpacity); c.A() > 0 {
				cv.FillPath(band, c, canvas.FillNonZero)
			}
		}
		return
	}
	c := st.Stroke
	if st.StrokeCurrent {
		c = paints.current
	}
	c = fade(c, st.Opacity*st.StrokeOpacity)
	if c.A() == 0 {
		return
	}
	cv.FillPath(band, c, canvas.FillNonZero)
}
