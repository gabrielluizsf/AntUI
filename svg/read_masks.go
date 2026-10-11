package svg

// readMasks takes every `<mask>` that named itself, under the id a `mask` may
// point at, and builds the picture each one holds. They are named before any
// of them is built, the way the clip paths are, because a mask may be written
// after the element that names it and one mask may be named from inside
// another — neither can be followed before every id is known, and a mask whose
// picture points at one written further down the file would otherwise find
// nothing there at all.
func (img *Image) readMasks(root *element) {
	type declared struct {
		id string
		e  *element
	}
	var found []declared
	seen := map[string]bool{}
	var walk func(e *element)
	walk = func(e *element) {
		if e.Name == "mask" {
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

	img.masks = make(map[string]*maskDef, len(found))
	for _, d := range found {
		img.masks[d.id] = &maskDef{}
	}
	// What a mask says about itself comes before any picture is built, because
	// the region is measured the moment an element that names the mask is
	// built — and one mask may name another, so a picture built before every
	// region is known would be measured against a region that was not filled
	// in yet, which is a zero rectangle reaching nothing, and would cut away
	// everything it was put on.
	for _, d := range found {
		img.readMask(d.id, d.e)
	}
	for _, d := range found {
		img.readMaskContent(d.id, d.e)
	}
}
