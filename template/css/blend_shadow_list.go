package css

// blendShadowList blends two shadow lists, shadow by shadow. Lists of equal
// length blend in step; a list that grew or shrank swaps at the half way
// point, since there is no honest better guess.
func blendShadowList(a, b []Shadow, t float64) []Shadow {
	if len(a) != len(b) {
		if t < 0.5 {
			return a
		}
		return b
	}
	if len(a) == 0 {
		return nil
	}
	out := make([]Shadow, len(a))
	for i := range a {
		out[i] = blendShadow(a[i], b[i], t)
	}
	return out
}
