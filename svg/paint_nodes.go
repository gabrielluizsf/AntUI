package svg

import "github.com/gabrielluizsf/antui/canvas"

// paintNodes walks the tree, drawing each node and then its children inside it.
// The matrix is the one the node's own transform composed onto the one it
// inherits, so a child is drawn where its parent put it.
//
// masking is the masks whose pictures are being drawn on the way to this node,
// and patterning the patterns whose tiles are being drawn on the way to it. A
// mask may name one of them — itself, or two that name each other — and
// following that would draw the same picture while it is already being drawn,
// for ever; the mask is left off instead, which is what keeps a broken drawing
// a picture rather than a program that does not come back. A pattern is left
// off the same way, and the shape that asked for it takes the colour after the
// `url(...)` if it brought one.
func paintNodes(cv *canvas.Canvas, n *Node, m canvas.Matrix, current canvas.Color, masking []*maskDef, patterning []*pattern) {
	if n.Style.Hidden {
		return
	}
	// An element that asked for a coordinate system held still against the host
	// is painted, and has everything inside it painted, through the transform
	// the effect makes of the one the drawing gave it.
	if n.Style.VectorEffect.constraining() {
		m = vectorEffectMatrix(m, n.Style)
	}
	if n.clip != nil || n.mask != nil || len(n.filters) > 0 {
		paintLayered(cv, n, m, current, masking, patterning)
		return
	}
	paintBody(cv, n, m, current, masking, patterning)
}
