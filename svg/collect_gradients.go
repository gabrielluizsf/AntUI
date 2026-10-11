package svg

// collectGradients takes every gradient element in the drawing, once each, under
// the id it answers to. An id nothing can name is still kept: a drawing may
// point at a gradient by an id it reuses, and a gradient declared without one is
// still part of the drawing even though no `url()` can reach it.
func collectGradients(e *element, img *Image, into map[string]*gradient) {
	if e == nil {
		return
	}
	if e.Name == "linearGradient" || e.Name == "radialGradient" {
		warn := func(format string, args ...any) { img.warn(e, format, args...) }
		g := parseGradient(e, warn)
		// A percentage written here is a fraction of the drawing, which is the
		// viewBox rather than the pixels it will be drawn into.
		g.viewport = [2]float64{img.ViewBox[2], img.ViewBox[3]}
		if id := e.attr("id"); id != "" {
			if _, taken := into[id]; !taken {
				into[id] = g
			}
		}
	}
	for _, k := range e.Kids {
		collectGradients(k, img, into)
	}
}
