package css

// originKeywordY reads the y of an origin: a vertical keyword or a length.
func originKeywordY(s string) (Length, bool) {
	switch s {
	case "top":
		return Pct(0), true
	case "center":
		return Pct(50), true
	case "bottom":
		return Pct(100), true
	}
	return Length{}, false
}
