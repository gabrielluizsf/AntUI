package css

// parseBackSizeOne reads one background-size value.
func parseBackSizeOne(part string, ctx Units) (BackSize, bool) {
	tokens := splitTokens(part)
	if len(tokens) == 0 {
		return BackSize{}, false
	}
	if len(tokens) == 1 {
		switch tokens[0] {
		case "cover":
			return BackSize{Cover: true}, true
		case "contain":
			return BackSize{Contain: true}, true
		}
	}
	if len(tokens) > 2 {
		return BackSize{}, false
	}
	var sz BackSize
	var vals []Length
	for _, tk := range tokens {
		l, err := parseLengthAt(tk, ctx)
		if err != nil {
			return BackSize{}, false
		}
		vals = append(vals, l)
	}
	sz.W = vals[0]
	if len(vals) > 1 {
		sz.H = vals[1]
	} else {
		sz.H = Auto()
	}
	return sz, true
}
