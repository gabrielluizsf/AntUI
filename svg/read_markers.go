package svg

// readMarkers takes every `<marker>` that named itself, under the ids the three
// marker properties may point at, and builds the picture each one holds. They
// are named before any of them is built, the way the clip paths, the masks and
// the patterns are: a marker may be written after the element that names it,
// and what is inside one may name anything else the drawing has — a gradient,
// a clip, a mask, a `<use>`, or another marker written further down the file.
func (img *Image) readMarkers(root *element) {
	type declared struct {
		id string
		e  *element
	}
	var found []declared
	seen := map[string]bool{}
	var walk func(e *element)
	walk = func(e *element) {
		if e.Name == "marker" {
			if id := e.attr("id"); id != "" && !seen[id] {
				seen[id] = true
				found = append(found, declared{id: id, e: e})
			}
		}
		for _, k := range e.Kids {
			walk(k)
		}
	}
	walk(root)

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
		img.markers[d.id].content = img.build(d.e, baseStyle(), true)
	}
	img.cutMarkerCycles()
}
