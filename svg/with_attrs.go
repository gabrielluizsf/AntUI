package svg

// withAttrs gathers an element's presentation attributes, with the `style="…"`
// declarations read on top: a declaration wins where both name the same thing.
func withAttrs(e *element) map[string]string {
	attrs := map[string]string{}
	for k, v := range e.Attr {
		if k != "" && k != "style" {
			attrs[k] = v
		}
	}
	for k, v := range parseDeclarations(e.attr("style")) {
		attrs[k] = v
	}
	return attrs
}
