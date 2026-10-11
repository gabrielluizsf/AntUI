package svg

// readDropShadow reads the inside of a `drop-shadow(...)`: two or three lengths
// — the offset and, when it is there, the blur radius — and, in any position, a
// colour. Writing no colour leaves it `currentColor`, which is what CSS does,
// and the colour is only turned into pixels when the element is painted and its
// colour is known. Anything else is not a shadow this package can read, and the
// whole function is left out rather than half of it guessed.
func readDropShadow(args string) (*dropShadow, bool) {
	sh := &dropShadow{current: true}
	var lens []float64
	for _, tk := range filterTokens(args) {
		if isCurrentColor(tk) {
			sh.current = true
			sh.colour = 0
			continue
		}
		if c, ok := parsePaint(tk); ok {
			sh.current = false
			sh.colour = c
			continue
		}
		v, ok := parseLength(tk)
		if !ok {
			return nil, false
		}
		lens = append(lens, v)
	}
	if len(lens) < 2 || len(lens) > 3 {
		return nil, false
	}
	sh.dx, sh.dy = lens[0], lens[1]
	if len(lens) == 3 {
		if lens[2] < 0 {
			return nil, false
		}
		sh.blur = lens[2]
	}
	return sh, true
}
