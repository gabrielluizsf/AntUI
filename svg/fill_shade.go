package svg

import "github.com/gabrielluizsf/antui/canvas"

// fillShade works out what the writing is filled with: a gradient's shade, a
// pattern's, or one colour. shade is nil where the fill is one flat colour and
// ok is false where there is no fill at all.
func fillShade(n *Node, st Style, m canvas.Matrix, paints palette, w, h int,
	patterning []*pattern) (shade func(x, y int) canvas.Color, c canvas.Color, ok bool) {
	switch {
	case !st.HasFill:
		return nil, 0, false
	case st.fillGradient != nil:
		return gradientShade(st.fillGradient, n, m, paints.current, st.Opacity*st.FillOpacity), 0, true
	case st.fillPattern != nil:
		// Where the pattern has nothing to paint with — it names itself, the
		// tile has no size to it, the transform has no way back — the colour
		// after the `url(...)` is what the writing falls back on, and nothing
		// at all where it brought none.
		if s, ok := patternShade(st.fillPattern, n, m, w, h, paints.current,
			st.Opacity*st.FillOpacity, patterning); ok {
			return s, 0, true
		}
		if fb, ok := st.fillFallbackColour(paints.current); ok {
			return nil, fade(fb, st.Opacity*st.FillOpacity), true
		}
		return nil, 0, true
	default:
		c = st.Fill
		if st.FillCurrent {
			c = paints.current
		}
		if st.FillContext {
			c = paints.fill
		}
		return nil, fade(c, st.Opacity*st.FillOpacity), true
	}
}
