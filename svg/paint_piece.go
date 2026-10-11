package svg

import "github.com/gabrielluizsf/antui/canvas"

// paintPiece puts one run of writing down: its stroke first and its fill over
// the stroke, which is the order SVG paints and the order that leaves the half
// of a stroke inside the letters covered by the fill and the half outside
// showing as the outline it is.
func paintPiece(cv *canvas.Canvas, n *Node, p *piece, m canvas.Matrix,
	current canvas.Color, patterning []*pattern) {
	st := p.run.Style
	put := func(target *canvas.Canvas, c canvas.Color, dx, dy float64) {
		switch {
		case p.run.Path != nil:
			writeAlong(target, p.face, p.run, p.width, p.scale, c,
				p.at.Mul(canvas.Translate(dx, dy)))
		case len(p.run.XList) > 0 || len(p.run.YList) > 0 || st.LetterSpacing != 0:
			// A place for each letter, or a gap left after every one of them,
			// or both: either puts the letters down one at a time, because the
			// pen has to be moved between them rather than handed the whole
			// piece and told to write it in one go. The pen is not in the
			// matrix here — the places a list names are places of the drawing
			// itself, not distances from where the pen already was.
			writeLetters(target, p.face, p.run, p.at, p.scale, c,
				p.penX+dx, p.penY+dy, p.width)
		default:
			at := p.at.Mul(canvas.Translate(p.penX+dx, p.penY+dy))
			drawRun(target, p.face, p.run.Text, c, at)
			if heavy(st.FontWeight) {
				// The weight the face does not carry is the same writing a
				// second time, one pixel further on.
				drawRun(target, p.face, p.run.Text, c,
					at.Mul(canvas.Translate(1/p.scale, 0)))
			}
		}
	}
	if st.HasStroke && st.Width > 0 {
		c := st.Stroke
		if st.StrokeCurrent {
			c = current
		}
		if c := fade(c, st.Opacity*st.StrokeOpacity); c.A() > 0 {
			strokeWriting(cv, put, st.strokeWidth(m), c, p.scale)
		}
	}
	var shade func(x, y int) canvas.Color
	var c canvas.Color
	switch {
	case !st.HasFill:
		return
	case st.fillGradient != nil:
		shade = gradientShade(st.fillGradient, n, m, current, st.Opacity*st.FillOpacity)
	case st.fillPattern != nil:
		// Where the pattern has nothing to paint with — it names itself, the
		// tile has no size to it, the transform has no way back — the colour
		// after the `url(...)` is what the writing falls back on, and nothing
		// at all where it brought none.
		if s, ok := patternShade(st.fillPattern, n, m, cv.Width, cv.Height, current,
			st.Opacity*st.FillOpacity, patterning); ok {
			shade = s
		} else if fb, ok := st.fillFallbackColour(current); ok {
			c = fade(fb, st.Opacity*st.FillOpacity)
		}
	default:
		c = st.Fill
		if st.FillCurrent {
			c = current
		}
		c = fade(c, st.Opacity*st.FillOpacity)
	}
	if shade == nil {
		if c.A() > 0 {
			put(cv, c, 0, 0)
		}
		return
	}
	layShaded(cv, put, shade)
}
