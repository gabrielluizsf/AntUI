package css

// blendTransformList blends two transform lists, function by function. Lists
// of the same length with matching kinds interpolate; anything else swaps
// wholesale at the half way point, which is how CSS treats a transform that
// cannot match up.
func blendTransformList(a, b []TransformFunc, t float64) []TransformFunc {
	if len(a) != len(b) {
		if t < 0.5 {
			return a
		}
		return b
	}
	if len(a) == 0 {
		return nil
	}
	out := make([]TransformFunc, 0, len(a))
	for i := range a {
		fa, fb := a[i], b[i]
		if fa.Kind != fb.Kind {
			if t < 0.5 {
				return a
			}
			return b
		}
		combine := blendTransformFunc(fa, fb, t)
		out = append(out, combine)
	}
	return out
}
