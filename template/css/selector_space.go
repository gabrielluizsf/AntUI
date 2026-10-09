package css

// skipSpace consumes a run of whitespace, flagging a combinator when a token
// already stands before it.
func (s *selectorScan) skipSpace() {
	j := s.i
	for j < len(s.t) && isSpace(s.t[j]) {
		j++
	}
	if j < len(s.t) && s.saw {
		s.comb = true
	}
	s.i = j
}
