package svg

// skipSpace passes over the spaces, commas and tabs between two things, and any
// minus or plus signs that are written as their own token.
func (s *pathScanner) skipSpace() {
	for s.pos < len(s.src) {
		switch s.src[s.pos] {
		case ' ', ',', '\t', '\n', '\r':
			s.pos++
		default:
			return
		}
	}
}
