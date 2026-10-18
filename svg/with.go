package svg

// with reads the attributes written on one element over the style it inherited
// from the one above, and answers the style its children inherit in turn. Only
// what the element actually says changes: an attribute that is not there leaves
// the inherited value alone, which is what lets a `<g fill="red">` colour
// everything inside it and a single `<rect fill="blue">` be the one blue thing.
//
// The warn is called for each attribute that is written but cannot be read, so
// that a drawing with a number that is not one says so rather than being painted
// as though it had never asked.
func (s Style) with(e *element, warn func(string, ...any)) Style {
	attrs := s.withVars(withAttrs(e))

	// pathLength is of the path on the element it is written on, so it does not
	// travel: a group that declares how long its own shape is says nothing
	// about the shapes inside, and each of those starts from none of its own.
	s.pathLength = 0

	// vector-effect does not travel either, and for a plainer reason than the
	// length: the spec says it is not inherited, so what a `<g>` asks for is
	// about the `<g>` itself and not the value every shape inside reads. What
	// does reach them is the coordinate system the effect makes of the `<g>`.
	s.VectorEffect = VectorEffect{}

	s.applyAttrs(attrs, warn)

	// xml:base travels down, resolved against the base already in hand, so a
	// relative address written inside the drawing resolves against the nearest
	// base it was given rather than the one at the top.
	s.base = readBase(attrs["xml:base"], s.base)

	// The dash pattern and its offset are read together once every other
	// attribute has been, because the order a map is walked in is not one to
	// depend on. The offset on its own says nothing, so it is only read where
	// there is a pattern for it to move.
	if _, ok := attrs["stroke-dasharray"]; ok {
		s.Dash = parseDashArray(attrs["stroke-dasharray"])
	}
	if _, ok := attrs["stroke-dashoffset"]; ok && s.Dash != nil {
		if v, ok := parseLength(attrs["stroke-dashoffset"]); ok {
			s.Dash.Offset = v
		} else {
			warn("the dash offset %q is not a length, so the pattern starts where it would have", attrs["stroke-dashoffset"])
		}
	}
	return s
}
