package svg

// readMarkers takes every `<marker>` that named itself, under the ids the three
// marker properties may point at, and builds the picture each one holds. They
// are named before any of them is built, the way the clip paths, the masks and
// the patterns are: a marker may be written after the element that names it,
// and what is inside one may name anything else the drawing has — a gradient,
// a clip, a mask, a `<use>`, or another marker written further down the file.
//
// Each marker's picture is built from the style that reaches the `<marker>` from
// the elements containing it — not from the element that names it, which is
// often far away and a different colour. A `<g fill="red">` with a marker in it
// is what makes that marker's own picture red, which is how a drawing says how
// its markers are painted in one place. The style is read along the walk here
// with the warning set aside, because the elements it passes are built, and
// complained about, where they are drawn.
func (img *Image) readMarkers(root *element) {
	type declared struct {
		id string
		e  *element
		in Style
	}
	var found []declared
	seen := map[string]bool{}
	var walk func(e *element, inherited Style)
	walk = func(e *element, inherited Style) {
		if e.Name == "marker" {
			if id := e.attr("id"); id != "" && !seen[id] {
				seen[id] = true
				found = append(found, declared{id: id, e: e, in: inherited})
			}
		}
		st := inherited.with(e, discardWarn)
		// A clip, a mask and a filter are spent on the element that writes them
		// and not passed on, the way [Image.build] takes them off before its
		// children are built: a marker is not cut by its container's clip.
		st.clipRef, st.maskRef, st.filters = "", "", nil
		for _, k := range e.Kids {
			walk(k, st)
		}
	}
	walk(root, img.rootStyle())

	img.markers = make(map[string]*markerDef, len(found))
	// Every marker is named with nothing in it first, and only then is any of
	// them read and built: a marker's picture may point at another by an id
	// that has not been reached yet, and a reference to a marker that is not
	// there is a warning rather than a picture.
	for _, d := range found {
		img.markers[d.id] = &markerDef{id: d.id}
	}
	for _, d := range found {
		img.readMarker(img.markers[d.id], d.e)
	}
	for _, d := range found {
		img.markers[d.id].content = img.build(d.e, d.in, true)
	}
	img.cutMarkerCycles()
}
