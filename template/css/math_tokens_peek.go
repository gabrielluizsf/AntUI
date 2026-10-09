package css

func (t *mathTokens) peek() string {
	ts := t.split()
	if t.i >= len(ts) {
		return ""
	}
	return ts[t.i]
}
