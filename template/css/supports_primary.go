package css

func (r *supportsReader) primary() bool {
	r.skipSpace()
	if r.eof() {
		r.bad = true
		return false
	}
	if r.src[r.i] == '(' {
		inner, ok := r.group()
		if !ok {
			r.bad = true
			return false
		}
		return r.inside(inner)
	}
	name, n := readName(r.src[r.i:])
	if name == "" {
		r.bad = true
		return false
	}
	r.i += n
	inner, ok := r.group()
	if !ok {
		r.bad = true
		return false
	}
	return r.function(name, inner)
}

// inside answers what a pair of parentheses holds: a declaration when a colon
// splits it at the top level, a nested condition when it does not.
func (r *supportsReader) inside(inner string) bool {
	if d, ok := splitSupportDecl(inner); ok {
		return supportsDecl(d.Prop, d.Raw)
	}
	sub := &supportsReader{src: inner}
	ok := sub.condition()
	if sub.bad || !sub.atEnd() {
		r.bad = true
		return false
	}
	return ok
}
