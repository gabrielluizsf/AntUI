package svg

// readFilterPrims reads a filter's primitives in order into def. What a
// primitive may name: the two pictures the element drew, the background this
// package has no way to reach, and every result an earlier primitive in this
// filter left behind. Naming anything else is said now rather than at paint
// time.
func readFilterPrims(def *filterDef, e *element, warn func(string, ...any)) {
	known := map[string]bool{
		"SourceGraphic":   true,
		"SourceAlpha":     true,
		"BackgroundImage": true,
		"BackgroundAlpha": true,
	}
	for _, kid := range e.Kids {
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
}
