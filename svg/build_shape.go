package svg

import "github.com/gabrielluizsf/antui/canvas"

// buildShape builds the shape this node draws, where it was written, and then
// moves it and its dashed outline by the element's transform.
func (img *Image) buildShape(n *Node, e *element, st Style) {
	n.Path = img.shape(e, st)
	if n.Path == nil {
		return
	}
	measureShape(n, st)
	placeShape(n, st)
}

// measureShape measures the shape as the drawing wrote it: the box a gradient
// in objectBoundingBox units is a fraction of, the length the drawing declared,
// and the outline cut into its dash pattern.
func measureShape(n *Node, st Style) {
	if minX, minY, maxX, maxY, ok := n.Path.Bounds(); ok {
		// The box before the transform moved the shape: the transform moves the
		// paint over the shape along with it rather than stretching it again.
		n.box = [4]float64{minX, minY, maxX - minX, maxY - minY}
	}
	// The declared length is spent on this node's own style, before the
	// transform moves the path, because it is a length of the shape as it was
	// written; the children are built from a style that still holds the dash
	// pattern as it was written, for the same reason.
	applyPathLength(&n.Style, n.Path)
	// Only what is going to be stroked is cut at all, by the same two
	// conditions the painting checks before it widens anything. A pattern that
	// cuts nothing hands the path back as it is, and that is the shape itself
	// rather than a second copy of it.
	if st.HasStroke && st.Width > 0 && n.Style.Dash != nil {
		if cut := n.Path.Dashed(*n.Style.Dash); cut != n.Path {
			n.dashed = cut
		}
	}
}

// placeShape puts the transform into the shape, not into the pixels around it:
// the shape goes where the transform puts it and the stroke is widened
// afterwards, in the coordinates it landed in. A group has no shape of its own
// to put it on, and carries the transform to its children instead.
func placeShape(n *Node, st Style) {
	// The vertices the markers land on and the directions the path is going
	// there are numbers of the shape as it was written, and the transform below
	// takes them along with the vertices the same way it takes the shape.
	n.markers = placeMarkers(n.Path, n.Style)
	if st.HasTransform && st.Transform != (canvas.Matrix{}) {
		n.Path.Transform(st.Transform)
		if n.dashed != nil {
			n.dashed.Transform(st.Transform)
		}
	}
}
