package css

// expr parses a sum of terms.
func (p *mathParser) expr() (float64, error) {
	v, err := p.term()
	if err != nil {
		return 0, err
	}
	for {
		switch p.t.peek() {
		case "+", "-":
			op := p.t.next()
			w, err := p.term()
			if err != nil {
				return 0, err
			}
			if op == "+" {
				v += w
			} else {
				v -= w
			}
		default:
			return v, nil
		}
	}
}
