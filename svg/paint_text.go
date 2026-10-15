package svg

import "github.com/gabrielluizsf/antui/canvas"

// paintText writes a `<text>`: its runs one after another with the pen carried
// from one to the next, which is what puts the writing of a `<tspan>` after the
// writing before it rather than each of them at the origin.
//
// The pen is in the drawing's own coordinates and the letters are put down as
// pixels, so each piece goes through the transform the node was given: the size
// it is drawn at is its size in those units brought onto the canvas, and where
// it starts is the pen put through the same transform. The pen then moves on by
// however wide the piece came out, so the piece after it starts where this one
// ended rather than where it began.
//
// The writing is walked before it is drawn rather than as it is drawn, because
// the fill may be a colour at every point — a gradient or a pattern — and
// neither can answer which colour a pixel takes until the box the writing
// covers is known. The stroke goes down before the fill the same way a shape's
// does, so the part of it that falls inside the letters is covered by the fill
// and the part that falls outside is the stroke showing.
func paintText(cv *canvas.Canvas, n *Node, m canvas.Matrix, paints palette, patterning []*pattern) {
	pieces := walkRuns(n.Runs, m)
	if len(pieces) == 0 {
		return
	}
	// The box the writing covers, as the drawing wrote it rather than as the
	// transform moved it: the same box a shape's geometry is measured to, and
	// the one a gradient or a pattern in fractions of the box it paints is a
	// fraction of. Writing is the only thing here with no box of its own until
	// a font is asked, and this is the font being asked.
	measured := *n
	measured.box = writingBox(pieces)
	for i := range pieces {
		paintPiece(cv, &measured, &pieces[i], m, paints, patterning)
	}
}
