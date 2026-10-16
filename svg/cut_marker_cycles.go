package svg

// cutMarkerCycles takes the marker references inside markers that would go
// round for ever. A marker whose picture has a shape with a marker of its own
// is ordinary — a marker made of markers — and is drawn as it says; one whose
// picture points back at the marker it is in the middle of is a circle, and
// drawing it would never stop. The reference that closes the circle is left off
// with one warning about it, which draws that one vertex without a marker
// rather than drawing nothing at all. It does not matter which of the three
// properties made the circle: all of them are marks on a vertex and all of them
// are walked here.
//
// A circle through a mask or a pattern needs no cutting here: those are already
// followed with everything that is being followed on the way to them, and the
// second time round the one that is already being drawn is left off — see
// [paintMasked] and [patternShade], and [Image.warnMaskPatternCycles] for
// saying which reference closes one of those. What is left over once they have
// had their say is the markers that point only at markers, which is what this
// walks.
func (img *Image) cutMarkerCycles() {
	gray := map[*markerDef]bool{}
	done := map[*markerDef]bool{}
	var walk func(n *Node)
	var visit func(def *markerDef)
	walk = func(n *Node) {
		if n == nil {
			return
		}
		for i := range n.markers {
			place := &n.markers[i]
			if place.def == nil {
				continue
			}
			switch {
			case gray[place.def]:
				img.warnNamed(n.Name,
					"the %s here points back at the marker this picture is inside, which would draw for ever, so no marker is drawn at this vertex",
					place.prop)
				place.def = nil
			case !done[place.def]:
				visit(place.def)
			}
		}
		for _, k := range n.Kids {
			walk(k)
		}
	}
	visit = func(def *markerDef) {
		if done[def] {
			return
		}
		gray[def] = true
		walk(def.content)
		gray[def] = false
		done[def] = true
	}
	for _, def := range img.markers {
		visit(def)
	}
}
