package css

func parseBorderStyles(raw string) ([sideCount]uint8, bool) {
	var out [sideCount]uint8
	parts := splitWords(raw)
	if len(parts) == 0 || len(parts) > sideCount {
		return out, false
	}
	vals := make([]uint8, len(parts))
	for i, p := range parts {
		v, ok := parseBorderStyle(p)
		if !ok {
			return out, false
		}
		vals[i] = v
	}
	fill := expandFour(vals)
	copy(out[:], fill[:])
	return out, true
}
