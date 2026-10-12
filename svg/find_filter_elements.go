package svg

// findFilterElements collects each `<filter>` in the drawing once, in the
// order they are written.
func findFilterElements(root *element) []filterElement {
	var found []filterElement
	seen := map[string]bool{}
	var walk func(e *element)
	walk = func(e *element) {
		if e.Name == "filter" {
			if id := e.attr("id"); id != "" && !seen[id] {
				seen[id] = true
				found = append(found, filterElement{id: id, e: e})
			}
		}
		for _, kid := range e.Kids {
			walk(kid)
		}
	}
	walk(root)
	return found
}
