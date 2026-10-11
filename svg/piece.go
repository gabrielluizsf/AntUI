package svg

import "github.com/gabrielluizsf/antui/canvas"

// piece is one run of writing carried to where it goes: the face it is written
// in at the size it comes out at on the canvas, how far it comes out in the
// drawing's own units, where the pen stood when it started, and the matrix that
// takes it from there — the whole of what putting it down needs.
type piece struct {
	run        TextRun
	face       *canvas.Face
	scale      float64
	width      float64
	penX, penY float64
	at         canvas.Matrix
}
