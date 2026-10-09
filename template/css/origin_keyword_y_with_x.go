package css

// originKeywordYWithX reads the y of an origin in two-value form, where the
// horizontal keywords are not allowed and a length is.
func originKeywordYWithX(s string, ctx Units) (Length, bool) {
	switch s {
	case "top":
		return Pct(0), true
	case "bottom":
		return Pct(100), true
	case "center":
		return Pct(50), true
	}
	return originAxisLength(s, ctx)
}
