package svg

import "github.com/gabrielluizsf/antui/canvas"

// paintLayered draws a node into a picture of its own, puts that picture
// through the filter list the element asked for, cuts the result down to the
// shapes its clip-path named, takes what its mask says of what is left over,
// and lays it over what is already there. The picture has to be made before it
// is touched, whole: a group is one picture with one filter, one clip and one
// mask over it, and cutting each of its children as it was drawn would cut
// them apart from each other, so that two translucent shapes overlapping under
// the edge of the clip would meet somewhere the clip does not keep at all.
//
// The order the picture goes through them is the order the spec puts them in:
// the filter first, because it changes the picture itself and reaches a little
// past where the element stopped; then the clip, which says where of that
// picture is kept; then the mask, which says how much of what is left is
// there. A clip or a mask measured before the filter would cut away the reach
// of a blur that was still coming, and the filter after them would soften the
// cut they made instead of the picture they were cut from.
//
// The picture is the size of the whole canvas rather than of the node, because
// a clip may reach anywhere and the cut has to be of the picture where the clip
// says and nowhere else. What a `<clipPath>` keeps nothing of comes out of this
// the same way: a picture drawn, cut to nothing, and laid over what was there
// as nothing at all — and so does a mask whose reference the drawing could not
// follow, which keeps nothing rather than keeping everything.
func paintLayered(cv *canvas.Canvas, n *Node, m canvas.Matrix, paints palette, masking []*maskDef, patterning []*pattern) {
	if n.mask != nil && n.mask.def.content == nil {
		// There is no picture to make: the mask keeps nothing, so nothing of
		// what it masks is drawn.
		return
	}
	layer, err := canvas.NewLayer(cv.Width, cv.Height)
	if err != nil {
		// There is no picture to make — a canvas of no size has nothing to
		// draw on and nothing to draw into either.
		return
	}
	paintBody(layer, n, m, paints, masking, patterning)
	if len(n.filters) > 0 {
		// The scale the element is drawn at goes with the list because what
		// the group's transform drew is bigger or smaller than the writing
		// says: a blur written as 1 is 1 of the drawing's units wherever those
		// units landed on the canvas. The transform is composed the same way
		// the children were drawn with it. What comes back may be another
		// picture than the one that went in, since every filter primitive
		// leaves a picture of its own behind.
		layer = applyFilters(layer, n.filters, scalesOf(m.Mul(n.Style.Transform)), paints.current)
	}
	if n.clip != nil {
		layer.MaskShapes(n.clip.measured(m)...)
	}
	if n.mask != nil {
		paintMasked(layer, n, m, paints, masking, patterning)
	}
	cv.BlitOver(0, 0, layer)
}
