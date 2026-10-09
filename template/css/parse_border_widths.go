package css

func parseBorderWidths(raw string) ([4]int, bool) {
	var out [4]int
	parts := splitWords(raw)
	if len(parts) == 0 || len(parts) > 4 {
		return out, false
	}
	vals := make([]int, len(parts))
	for i, p := range parts {
		l, err := parseLength(p)
		if err != nil || l.u != unitPx || l.value < 0 {
			return out, false
		}
		vals[i] = int(l.value + 0.5)
	}
	fill := expandFour(vals)
	copy(out[:], fill[:])
	return out, true
}
