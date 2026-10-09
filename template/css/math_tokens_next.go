package css

func (t *mathTokens) next() string {
	ts := t.split()
	if t.i >= len(ts) {
		return ""
	}
	v := ts[t.i]
	t.i++
	return v
}
