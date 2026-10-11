package svg

// peek is a short look at the text ahead, for a message about what was found
// where a tag was expected.
func (p *xmlParser) peek(n int) string {
	end := min(p.pos+n, len(p.src))
	return p.src[p.pos:end]
}
