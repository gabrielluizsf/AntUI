package svg

import "github.com/gabrielluizsf/antui/canvas"

// paintMasked takes what the element painted down to what its mask keeps of
// it: the mask's picture is drawn into a picture of its own, at the same size
// and over the same transform so that what it says lines up with what it is
// being said about, and then every pixel of the element is multiplied by how
// much of the mask is over it.
//
// The mask's picture is drawn under the element's own transform, because a
// mask — like the shapes a clip cuts with — is written in the coordinates of
// the element it is put on: a `transform` on that element moves the mask along
// with everything else it moves. The picture of one written in fractions of
// the element's box is laid out over that box instead, with no transform of
// the element over it again: the box already carries the element's transform,
// and the content and the region — both fractions of the same box — would
// otherwise land in two places at once. Its region is measured in the
// drawing's own coordinates and taken to the canvas by the same transform the
// element is, for the same reason.
func paintMasked(layer *canvas.Canvas, n *Node, m canvas.Matrix, paints palette, masking []*maskDef, patterning []*pattern) {
	md := n.mask
	for _, being := range masking {
		if being == md.def {
			// This mask is already being drawn further up, so following it
			// again would draw it while it is being drawn, for ever: a mask
			// that names itself, or two that name each other, is left off and
			// the element stands as it is rather than the drawing never
			// coming back.
			return
		}
	}
	ml, err := canvas.NewLayer(layer.Width, layer.Height)
	if err != nil {
		// With no picture of the mask there is nothing to say how much is
		// kept, and what no mask has anything to say about keeps nothing.
		layer.MaskShapes()
		return
	}
	if md.def.content != nil {
		content := m.Mul(n.Style.Transform)
		if md.def.contentBoxUnits {
			content = m.Mul(canvas.Translate(md.box[0], md.box[1])).Mul(canvas.Scale(md.box[2], md.box[3]))
		}
		paintNodes(ml, md.def.content, content, paints, append(masking, md.def), patterning)
	}
	if md.cuts {
		x0, y0, x1, y1 := regionArea(md.x, md.y, md.w, md.h, m)
		layer.MaskRect(canvas.Area{X: x0, Y: y0, Width: x1 - x0, Height: y1 - y0})
	}
	layer.MaskBy(ml, md.def.measure)
}
