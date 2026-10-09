package css

// condition reads the or-list at the top. Both sides are always read, even
// when the left one already settled the answer, so a malformed right-hand side
// is still reported instead of left behind for atEnd to call trailing garbage.
func (r *supportsReader) condition() bool {
	ok := r.conjunction()
	for r.word("or") {
		if r.conjunction() {
			ok = true
		}
	}
	return ok
}

func (r *supportsReader) conjunction() bool {
	ok := r.negation()
	for r.word("and") {
		if !r.negation() {
			ok = false
		}
	}
	return ok
}

// negation reads a "not" — which binds to what follows it, tighter than and
// and or — and lets a second one in, so `not not (x)` is the same test as
// `(x)`.
func (r *supportsReader) negation() bool {
	if r.word("not") {
		return !r.negation()
	}
	return r.primary()
}
