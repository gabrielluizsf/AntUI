package svg

// gradientHref is the id another gradient's attributes come from. An href is
// written as a plain `#id`, and files in the wild spell it `url(#id)` the way a
// paint is spelled, since the two say the same thing to a reader of the drawing.
func gradientHref(attrs map[string]string) string {
	for _, name := range []string{"href", "xlink:href"} {
		v, ok := attrs[name]
		if !ok {
			continue
		}
		if id, ok := hashRef(v); ok {
			return id
		}
	}
	return ""
}
