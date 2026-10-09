package css

// group reads one parenthesised group and returns what is inside it, with the
// strings and the nesting inside it left alone.
func (r *supportsReader) group() (string, bool) {
	r.skipSpace()
	if r.eof() || r.src[r.i] != '(' {
		return "", false
	}
	depth := 0
	start := 0
	for !r.eof() {
		c := r.src[r.i]
		if c == '"' || c == '\'' {
			end := quoteEnd(r.src, r.i)
			if end < 0 {
				return "", false
			}
			r.i = end + 1
			continue
		}
		switch c {
		case '(':
			if depth == 0 {
				start = r.i + 1
			}
			depth++
		case ')':
			depth--
			if depth == 0 {
				r.i++
				return r.src[start : r.i-1], true
			}
		}
		r.i++
	}
	return "", false
}
