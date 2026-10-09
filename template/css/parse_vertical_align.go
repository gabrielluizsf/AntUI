package css

// parseVerticalAlign reads a vertical-align keyword or length. A length lands
// in the returned BaselineShift; keywords leave it zero.
func parseVerticalAlign(raw string, ctx Units) (uint8, Length, bool) {
	switch raw {
	case "baseline":
		return VerticalAlignBaseline, Length{}, true
	case "sub":
		return VerticalAlignSub, Length{}, true
	case "super":
		return VerticalAlignSuper, Length{}, true
	case "middle":
		return VerticalAlignMiddle, Length{}, true
	case "top":
		return VerticalAlignTop, Length{}, true
	case "bottom":
		return VerticalAlignBottom, Length{}, true
	case "text-top":
		return VerticalAlignTextTop, Length{}, true
	case "text-bottom":
		return VerticalAlignTextBottom, Length{}, true
	}
	if l, err := parseLengthAt(raw, ctx); err == nil {
		return VerticalAlignBaseline, l, true
	}
	return 0, Length{}, false
}
