package svg

import "github.com/gabrielluizsf/antui/canvas"

// paintBody is one node and everything inside it, drawn straight onto the
// canvas it is given: its shape first, then its writing, then the picture an
// `<image>` is, then its children in the order they were written.
func paintBody(cv *canvas.Canvas, n *Node, m canvas.Matrix, paints palette, masking []*maskDef, patterning []*pattern) {
	if n.Path != nil && !n.Path.Empty() {
		paintShape(cv, n, m, paints, masking, patterning)
	}
	if len(n.Runs) > 0 {
		paintText(cv, n, m, paints, patterning)
	}
	if n.Pic != nil {
		paintPicture(cv, n, m)
	}
	for _, k := range n.Kids {
		paintNodes(cv, k, m, paints, masking, patterning)
	}
}
