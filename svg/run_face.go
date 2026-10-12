package svg

import "github.com/gabrielluizsf/antui/canvas"

// runFace is the face one run is drawn with and the width it covers, in the
// drawing's own coordinates. The face is asked for at the size in pixels this
// writing covers on the canvas, so a drawing painted big is written big rather
// than a small piece of writing blown up afterwards. The face itself is the one
// the canvas draws with: a font-family list is read, but a drawing carries no
// files of its own, so the names run into the face the program already had —
// see [Style].
func runFace(run TextRun, scale float64) (*canvas.Face, float64) {
	st := run.Style
	size := st.FontSize
	if size <= 0 {
		size = defaultFontSize
	}
	face := canvas.DefaultFace().AtSize(size * scale)
	width := (float64(face.Width(run.Text)) +
		st.LetterSpacing*float64(spaced(st, run.Text))) / scale
	return face, width
}
