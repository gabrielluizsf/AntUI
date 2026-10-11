package svg

// readFilterDefs reads every `<filter>` in the drawing, by the id an element's
// `filter="url(#name)"` names. They are read while the drawing's ids are all
// in place and before anything is built, for the same reason as the masks and
// the gradients: a filter may be written after the element that points at it.
// What is inside a filter is only what its primitives say — no picture painted
// with the rest of the drawing — so the ids being in place is all one needs to
// be read.
func (img *Image) readFilterDefs(root *element) {
	type declared struct {
		id string
		e  *element
	}
	var found []declared
	seen := map[string]bool{}
	var walk func(e *element)
	walk = func(e *element) {
		if e.Name == "filter" {
			if id := e.attr("id"); id != "" && !seen[id] {
				seen[id] = true
				found = append(found, declared{id: id, e: e})
			}
		}
		for _, kid := range e.Kids {
			walk(kid)
		}
	}
	walk(root)
	img.filters = make(map[string]*filterDef, len(found))
	for _, d := range found {
		def := &filterDef{}
		warn := func(format string, args ...any) { img.warn(d.e, format, args...) }
		// What a primitive may name: the two pictures the element drew, the
		// background this package has no way to reach, and every result an
		// earlier primitive in this filter left behind. Naming anything else
		// is said now rather than at paint time.
		known := map[string]bool{
			"SourceGraphic":   true,
			"SourceAlpha":     true,
			"BackgroundImage": true,
			"BackgroundAlpha": true,
		}
		for _, kid := range d.e.Kids {
			p, ok := readFilterPrim(kid, warn)
			if !ok {
				continue
			}
			for _, name := range []string{p.in, p.in2} {
				if name == "" || known[name] {
					if name == "BackgroundImage" || name == "BackgroundAlpha" {
						warn("the background %q a filter reads is not where a filter can see it, so it is taken as nothing", name)
					}
					continue
				}
				warn("the filter input %q names nothing this filter has, so it is taken as nothing", name)
			}
			def.prims = append(def.prims, p)
			if p.result != "" {
				known[p.result] = true
			}
		}
		img.filters[d.id] = def
	}
}
