package svg

import "github.com/gabrielluizsf/antui/canvas"

// paintMarkers draws the markers an element asked for at the vertices of its
// shape, in the order the vertices come: each marker's own picture taken to the
// canvas by the matrix that was worked out where the shape was built, and then
// by the one that puts the drawing there. The pictures are drawn the same way
// everything else is — their own fill and stroke, their own groups, their own
// clip and mask — and where the room of the marker cuts its picture, the way
// any element with a clip-path on it is drawn: through a picture of its own,
// cut to the room, laid down over what was there. The masks and patterns
// already being drawn on the way to them go with them, so that a marker inside
// a pattern inside a marker finds the pattern already there rather than drawing
// it again for ever. A marker the cycle cut left off has no marker to draw and
// is passed over.
func paintMarkers(cv *canvas.Canvas, n *Node, m canvas.Matrix, current canvas.Color, masking []*maskDef, patterning []*pattern) {
	for i := range n.markers {
		place := &n.markers[i]
		if place.def == nil || place.def.content == nil {
			continue
		}
		content := place.def.content
		if place.room != nil {
			wrap := &Node{clip: place.room}
			wrap.Kids = []*Node{content}
			content = wrap
		}
		paintNodes(cv, content, m.Mul(place.at), current, masking, patterning)
	}
}
