package css

// blendFilterList blends two filter lists, filter by filter, with the same
// length rule as the shadows. A drop-shadow pair blends; a pair of the same
// function kind lerps its amount; anything else swaps at the half way point.
func blendFilterList(a, b []Filter, t float64) []Filter {
	if len(a) != len(b) {
		if t < 0.5 {
			return a
		}
		return b
	}
	if len(a) == 0 {
		return nil
	}
	out := make([]Filter, len(a))
	for i := range a {
		out[i] = blendFilter(a[i], b[i], t)
	}
	return out
}
