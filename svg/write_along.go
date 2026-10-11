package svg

import "github.com/gabrielluizsf/antui/canvas"

// writeAlong lays one piece of writing along an outline: every letter is put
// down at the place the outline has reached by the time the letters before it
// have been written, turned to the way the outline is going there, so the
// writing curves with the line rather than lying flat across it.
func writeAlong(cv *canvas.Canvas, face *canvas.Face, run TextRun, width, scale float64, c canvas.Color, m canvas.Matrix) {
	start := run.Offset
	switch run.Style.Anchor {
	case AnchorMiddle:
		start -= width / 2
	case AnchorEnd:
		start -= width
	}
	d := start
	spacing := run.Style.LetterSpacing
	bold := heavy(run.Style.FontWeight)
	for _, r := range run.Text {
		if r == '\n' || r == '\r' {
			// An outline is one line: there is nowhere below it to go.
			continue
		}
		s := string(r)
		if r == '\t' {
			s = "    "
		}
		adv := float64(face.Width(s)) / scale
		if r != ' ' && r != '\t' {
			x, y, ang := pointAt(run.Path, d)
			at := m.Mul(canvas.Translate(x, y)).Mul(canvas.Rotate(ang))
			drawRun(cv, face, s, c, at)
			if bold {
				drawRun(cv, face, s, c, at.Mul(canvas.Translate(1/scale, 0)))
			}
		}
		d += adv + spacing
	}
}
