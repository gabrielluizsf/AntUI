package css

// parseFilters reads a space-separated filter list. "none" clears it.
func parseFilters(raw string, ctx Units) ([]Filter, bool) {
	if raw == "none" {
		return nil, true
	}
	var out []Filter
	for _, tk := range splitTokens(raw) {
		if f, ok := parseFilter(tk, ctx); ok {
			out = append(out, f)
		}
	}
	return out, len(out) > 0
}
