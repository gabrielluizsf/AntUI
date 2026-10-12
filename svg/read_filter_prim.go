package svg

// readFilterPrim reads one primitive of a `<filter>` — the six this package
// can apply, with the inputs it takes and the name it leaves its result under.
// Anything else inside a `<filter>`, and a primitive whose own numbers are not
// the ones it needs, is said and left out: the rest of the filter still runs,
// the same way one unreadable function in the attribute's list does.
func readFilterPrim(e *element, warn func(string, ...any)) (filterPrim, bool) {
	p := filterPrim{in: e.attr("in"), in2: e.attr("in2"), result: e.attr("result")}
	switch e.Name {
	case "feGaussianBlur":
		p.kind = primGaussianBlur
		if !p.readGaussianBlur(e, warn) {
			return filterPrim{}, false
		}
	case "feColorMatrix":
		p.kind = primColorMatrix
		if !p.readColorMatrix(e, warn) {
			return filterPrim{}, false
		}
	case "feOffset":
		p.kind = primOffset
		if !p.readOffset(e, warn) {
			return filterPrim{}, false
		}
	case "feFlood":
		p.kind = primFlood
		if !p.readFlood(e, warn) {
			return filterPrim{}, false
		}
	case "feComposite":
		p.kind = primComposite
		if !p.readComposite(e, warn) {
			return filterPrim{}, false
		}
	case "feBlend":
		p.kind = primBlend
		if !p.readBlend(e, warn) {
			return filterPrim{}, false
		}
	default:
		warn("the filter primitive %q is not one this package can read, so it is left out", e.Name)
		return filterPrim{}, false
	}
	return p, true
}
