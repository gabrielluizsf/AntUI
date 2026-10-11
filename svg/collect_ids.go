package svg

// collectIds walks the drawing once, taking every element that wrote an id.
// An id nothing points at is still kept: it costs nothing and a drawing may
// point at one later, from a `<use>` this has not read yet.
func collectIds(e *element, into map[string]*element) {
	if e == nil {
		return
	}
	if id := e.attr("id"); id != "" {
		if _, taken := into[id]; !taken {
			into[id] = e
		}
	}
	for _, k := range e.Kids {
		collectIds(k, into)
	}
}
