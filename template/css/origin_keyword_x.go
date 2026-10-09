package css

// originKeywordX reads the x of an origin: a horizontal keyword or a length.
func originKeywordX(s string, ctx Units) (Length, bool) {
	switch s {
	case "left":
		return Pct(0), true
	case "center":
		return Pct(50), true
	case "right":
		return Pct(100), true
	}
	return originAxisLength(s, ctx)
}
