package svg

import "github.com/gabrielluizsf/antui/canvas"

// maskUnder is what a node carries for the mask it was found to be under: the
// region the mask reaches, worked out against the box of everything the node
// paints — which is only known once everything inside it is there, and is why
// the reference is named before the node is built and this is done after.
//
// The region of a mask written in the drawing's own coordinates is moved by
// the element's transform, the same as the shapes of a clip are: what is
// written against this element is written in its coordinates, and a transform
// on it moves the whole of it. The region of one written in fractions of the
// box needs no moving, because the box it is a fraction of is the box the
// element paints at, transform and all — and so is the content of one in
// fractions of the box, which is laid out over that same box when the mask is
// painted. Where that box is asked for and there is none to measure — a group
// with nothing in it that paints anywhere — the element is drawn without the
// mask with one warning saying so: an element drawn whole is one an author
// can see, and a warning is what tells them it was not masked.
func maskUnder(def *maskDef, st Style, n *Node, warn func(string, ...any)) *masked {
	if def == nil {
		return nil
	}
	m := &masked{def: def}
	r := def.region
	var box [4]float64
	if r.boxUnits || def.contentBoxUnits {
		var ok bool
		box, ok = paintedBox(n)
		if def.contentBoxUnits {
			if !ok || box[2] <= 0 || box[3] <= 0 {
				warn("the mask is written in objectBoundingBox units against this element, which has no box to measure, so the element is drawn without it")
				return nil
			}
			m.box = box
		} else if !ok {
			// A region measured against a box that cannot be measured is not
			// one to cut with: a mask reaching everywhere is one an author can
			// see the whole of, and one reaching nowhere is a drawing that
			// vanished with nothing to say about it. Writing is measured the
			// same as a shape is now, so what is left with no box of its own
			// is a node with nothing in it that paints at all.
			return m
		}
	}
	switch {
	case !r.boxUnits:
		m.x, m.y, m.w, m.h = r.x, r.y, r.w, r.h
		if st.HasTransform && st.Transform != (canvas.Matrix{}) {
			m.x, m.y, m.w, m.h = movedBox(m.x, m.y, m.w, m.h, st.Transform)
		}
		m.cuts = true
	default:
		m.x, m.y = box[0]+r.x*box[2], box[1]+r.y*box[3]
		m.w, m.h = r.w*box[2], r.h*box[3]
		m.cuts = true
	}
	return m
}
