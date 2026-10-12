package svg

import "github.com/gabrielluizsf/antui/canvas"

// paintShape draws one node's shape: its fill first and its stroke over it,
// which is the order SVG paints, so a stroke half its width over a fill covers
// the fill's edge as it should, and the marker over the pair of them last of
// all — an arrowhead on the line rather than under it.
func paintShape(cv *canvas.Canvas, n *Node, m canvas.Matrix, current canvas.Color, masking []*maskDef, patterning []*pattern) {
	if n.Style.HasFill {
		paintFill(cv, n, m, current, n.Style, patterning)
	}
	if n.Style.HasStroke && n.Style.Width > 0 {
		paintStroke(cv, n, m, current, patterning)
	}
	// The marker comes after the fill and the stroke, which is the order the
	// spec puts them in. It is drawn where the shape has no stroke at all, and
	// where it encloses no fill either: a marker is a picture put on a vertex,
	// and the line beside it is not what puts it there — the vertex is.
	if len(n.markers) > 0 {
		paintMarkers(cv, n, m, current, masking, patterning)
	}
}
