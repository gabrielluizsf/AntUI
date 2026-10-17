package svg

// hasAttr reports whether the attribute was written at all, which is a
// different question from whether it was written empty: `fill=""` is not the
// same as no fill, and only the first is an instruction.
func (e *element) hasAttr(name string) bool {
	_, ok := e.Attr[name]
	return ok
}
