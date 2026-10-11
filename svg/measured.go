package svg

import "github.com/gabrielluizsf/antui/canvas"

// measured is the region of a node against the canvas it is being painted on:
// the shapes carry the element's own transform already, and the matrix that
// maps the drawing onto the canvas is what is left to put on them.
func (r *clipRegion) measured(m canvas.Matrix) []canvas.MaskShape {
	return mapMaskShapes(r.shapes, m)
}
