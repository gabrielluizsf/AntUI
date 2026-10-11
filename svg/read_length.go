package svg

// readLength is the length written on one attribute, answering zero where there
// is none, and saying so where there is one that cannot be read — a number that
// is not a number is something the drawing asked for and this could not do,
// which is a warning, while an attribute that is not there at all is just a
// default.
func readLength(warn func(string, ...any), e *element, attr, what string) float64 {
	if !e.hasAttr(attr) {
		return 0
	}
	v, ok := parseLength(e.attr(attr))
	if !ok {
		warn("the %s %q is not a length, so it is taken as none", what, e.attr(attr))
	}
	return v
}
