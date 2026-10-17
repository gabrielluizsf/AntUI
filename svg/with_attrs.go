package svg

// withAttrs gathers everything that was said about one element's paint, in the
// order the cascade weighs it: the presentation attributes first, then the
// drawing's own `<style>` sheets by weight, then the element's `style="…"` on
// top of them, and at the top an important declaration, whether a sheet or the
// `style` attribute wrote it. A declaration wins where two name the same thing.
func withAttrs(e *element) map[string]string {
	attrs := map[string]string{}
	for k, v := range e.Attr {
		if k != "" && k != "style" {
			attrs[k] = v
		}
	}
	if e.sheet != nil {
		for k, v := range e.sheet.normal {
			attrs[k] = v
		}
	}
	var important []cssDecl
	for _, d := range parseDecls(e.attr("style")) {
		if d.important {
			important = append(important, d)
			continue
		}
		attrs[d.name] = d.value
	}
	if e.sheet != nil {
		for k, v := range e.sheet.important {
			attrs[k] = v
		}
	}
	for _, d := range important {
		attrs[d.name] = d.value
	}
	return attrs
}
