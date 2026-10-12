package svg

// readFilterDefs reads every `<filter>` in the drawing, by the id an element's
// `filter="url(#name)"` names. They are read while the drawing's ids are all
// in place and before anything is built, for the same reason as the masks and
// the gradients: a filter may be written after the element that points at it.
// What is inside a filter is only what its primitives say — no picture painted
// with the rest of the drawing — so the ids being in place is all one needs to
// be read.
func (img *Image) readFilterDefs(root *element) {
	found := findFilterElements(root)
	img.filters = make(map[string]*filterDef, len(found))
	for _, d := range found {
		def := &filterDef{}
		warn := func(format string, args ...any) { img.warn(d.e, format, args...) }
		readFilterPrims(def, d.e, warn)
		img.filters[d.id] = def
	}
}
