package svg

import "github.com/gabrielluizsf/antui/canvas"

// writeLetters puts a piece of writing down one letter at a time: each letter
// at the place the drawing named for it when it named a place for every letter
// of the piece, and otherwise at the pen the letter before it left, with the
// letter-spacing the drawing asked for left after every one of them. Bold is
// the letter drawn a second time a pixel further on, the same as it is when
// the whole piece goes down at once.
//
// m is the transform of the writing without the pen's place in it, because a
// place a list names is a place of the drawing rather than a distance from
// wherever the pen was; penX and penY are the pen already moved by the anchor,
// for the letters the list did not reach, and width is what the anchor worked
// that movement out from.
func writeLetters(cv *canvas.Canvas, face *canvas.Face, run TextRun, m canvas.Matrix,
	scale float64, c canvas.Color, penX, penY, width float64) {
	st := run.Style
	spacing := st.LetterSpacing
	bold := heavy(st.FontWeight)
	shift := 0.0
	switch st.Anchor {
	case AnchorMiddle:
		shift = -width / 2
	case AnchorEnd:
		shift = -width
	}
	lineStart := penX
	for i, r := range run.Text {
		switch r {
		case '\r':
			continue
		case '\n':
			penX = lineStart
			penY += float64(face.Height()) / scale
			continue
		case '\t':
			penX += 4 * float64(face.Width(" ")) / scale
			continue
		}
		x, y := penX, penY
		if i < len(run.XList) {
			x = run.XList[i] + shift
		}
		if i < len(run.YList) {
			y = run.YList[i]
		}
		at := m.Mul(canvas.Translate(x, y))
		drawRun(cv, face, string(r), c, at)
		if bold {
			drawRun(cv, face, string(r), c, at.Mul(canvas.Translate(1/scale, 0)))
		}
		penX, penY = x+float64(face.Width(string(r)))/scale+spacing, y
	}
}
