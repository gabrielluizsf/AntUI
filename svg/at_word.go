package svg

// atWord says the parser stands at a word of the markup — DOCTYPE, ENTITY —
// with the space, the bracket or the `%` that follows it rather than at a
// longer name that only opens with those letters.
func (p *xmlParser) atWord(w string) bool {
	if !p.at(w) {
		return false
	}
	next := p.pos + len(w)
	return next >= len(p.src) || !isNameChar(p.src[next])
}
