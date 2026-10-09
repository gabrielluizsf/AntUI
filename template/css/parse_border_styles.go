package css

func parseBorderStyles(raw string) ([4]uint8, bool) {
	var out [4]uint8
	parts := splitWords(raw)
	if len(parts) == 0 || len(parts) > 4 {
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
	fill := expandFourInt8(vals)
	copy(out[:], fill[:])
	return out, true
}
