package css

// parseShadow reads one shadow value: an optional inset, two to four lengths,
// and a colour in any position.
func parseShadow(part string, ctx Units, allowInset bool, max int) (Shadow, bool) {
	var lens []int
	sh := Shadow{Color: CurrentColor}
	for _, tk := range splitTokens(part) {
		if allowInset && tk == "inset" {
			sh.Inset = true
			continue
		}
		if c, err := ParseColor(tk); err == nil {
			sh.Color = c
			continue
		}
		l, err := parseLengthAt(tk, ctx)
		if err != nil || l.IsPct() || l.Auto() || l.None() {
			return Shadow{}, false
		}
		lens = append(lens, l.Resolve(ctx))
	}
	if len(lens) < 2 || len(lens) > max {
		return Shadow{}, false
	}
	sh.X, sh.Y = lens[0], lens[1]
	if len(lens) > 2 {
		sh.Blur = lens[2]
	}
	if len(lens) > 3 {
		sh.Spread = lens[3]
	}
	return sh, true
}
