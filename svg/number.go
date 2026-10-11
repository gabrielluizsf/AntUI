package svg

// number reads one number, which may be written with an exponent as SVG allows.
func (s *pathScanner) number() (float64, bool) {
	s.skipSpace()
	start := s.pos
	if s.pos < len(s.src) && (s.src[s.pos] == '+' || s.src[s.pos] == '-') {
		s.pos++
	}
	seen := false
	for s.pos < len(s.src) && s.src[s.pos] >= '0' && s.src[s.pos] <= '9' {
		s.pos++
		seen = true
	}
	if s.pos < len(s.src) && s.src[s.pos] == '.' {
		s.pos++
		for s.pos < len(s.src) && s.src[s.pos] >= '0' && s.src[s.pos] <= '9' {
			s.pos++
			seen = true
		}
	}
	if seen && s.pos < len(s.src) && (s.src[s.pos] == 'e' || s.src[s.pos] == 'E') {
		mark := s.pos
		s.pos++
		if s.pos < len(s.src) && (s.src[s.pos] == '+' || s.src[s.pos] == '-') {
			s.pos++
		}
		digits := false
		for s.pos < len(s.src) && s.src[s.pos] >= '0' && s.src[s.pos] <= '9' {
			s.pos++
			digits = true
		}
		if !digits {
			s.pos = mark
		}
	}
	if !seen {
		s.pos = start
		return 0, false
	}
	v, err := parseNumber(s.src[start:s.pos])
	if err != nil {
		return 0, false
	}
	return v, true
}
