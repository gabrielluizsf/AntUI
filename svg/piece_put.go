package svg

import "github.com/gabrielluizsf/antui/canvas"

// put draws this piece of writing in one colour, moved by dx and dy. A run with
// an outline is written along it; one with a place for each letter, or a gap
// left after every one, is written a letter at a time; anything else is written
// in one go, a heavy weight being the same writing a second time one pixel on.
func (p *piece) put(target *canvas.Canvas, c canvas.Color, dx, dy float64) {
	st := p.run.Style
	switch {
	case p.run.Path != nil:
		writeAlong(target, p.face, p.run, p.width, p.scale, c,
			p.at.Mul(canvas.Translate(dx, dy)))
	case len(p.run.XList) > 0 || len(p.run.YList) > 0 || st.LetterSpacing != 0:
		// The pen is not in the matrix here — the places a list names are
		// places of the drawing itself, not distances from where the pen
		// already was.
		writeLetters(target, p.face, p.run, p.at, p.scale, c,
			p.penX+dx, p.penY+dy, p.width)
	default:
		at := p.at.Mul(canvas.Translate(p.penX+dx, p.penY+dy))
		drawRun(target, p.face, p.run.Text, c, at)
		if heavy(st.FontWeight) {
			// The weight the face does not carry is the same writing a second
			// time, one pixel further on.
			drawRun(target, p.face, p.run.Text, c,
				at.Mul(canvas.Translate(1/p.scale, 0)))
		}
	}
}
