package svg

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// drawRun puts one piece of writing on the canvas. The matrix is the whole of
// where it goes: the pen's own transform with the pen's place already in it,
// and for a letter along an outline, the turn of the outline at that letter
// too. A matrix that only moves and sizes the writing puts each letter down
// where the pen landed, which is the whole of what has to happen.
//
// A matrix that turns the writing is the other case: the letters are pixels
// and a pixel cannot lean, so the piece is written as it stands into a picture
// of its own and that picture is laid over through the turn — the same thing
// the shapes do, since a shape is a path taken through the transform before it
// is rasterised. The picture is made at the resolution the face was asked for
// and the matrix carries it the rest of the way, so the letters are rasterised
// once at their own size and stretched by the turn rather than drawn twice.
func drawRun(cv *canvas.Canvas, face *canvas.Face, text string, c canvas.Color, m canvas.Matrix) {
	if m.B == 0 && m.C == 0 {
		px, py := m.Map(0, 0)
		face.Draw(cv, int(px+0.5), int(py+0.5), text, c)
		return
	}
	const pad = 2
	scale := scaleOf(m)
	// One drawing unit is `scale` of the face's pixels, so the face's own
	// measurement is already the size the picture has to be. A piece with a
	// line break in it is as many lines tall as it has breaks.
	lines := 1 + strings.Count(text, "\n")
	w := face.Width(text) + 2*pad
	h := face.Ascent() + face.Descent() + (lines-1)*face.Height() + 2*pad
	if w <= 0 || h <= 0 {
		return
	}
	layer, err := canvas.NewLayer(w, h)
	if err != nil {
		return
	}
	face.Draw(layer, pad, pad+face.Ascent(), text, c)
	// Layer pixel (lx,ly) is the drawing's own (lx-pad)/scale and
	// (ly-pad-Ascent)/scale, and the matrix takes those from there onwards.
	to := m.Mul(canvas.Translate(-float64(pad)/scale,
		-float64(pad+face.Ascent())/scale).Mul(canvas.Scale(1/scale, 1/scale)))
	cv.BlitMatrix(layer, to, canvas.Area{Width: w, Height: h})
}
