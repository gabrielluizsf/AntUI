package css

// blendLength interpolates two lengths. Same-unit lengths lerp their values,
// so a width that goes 50% to 75% stays a percentage all the way. Keywords
// and units that cannot mix resolve to their reference pixels and lerp those.
func blendLength(a, b Length, t float64, ctx Units) Length {
	if a.u == b.u && a.u != unitAuto && a.u != unitNone {
		return Length{u: a.u, value: a.value + (b.value-a.value)*t}
	}
	if a.u == unitAuto || a.u == unitNone || b.u == unitAuto || b.u == unitNone {
		if t < 0.5 {
			return a
		}
		return b
	}
	ar, br := a.Resolve(ctx), b.Resolve(ctx)
	return Fixed(float64(ar) + float64(br-ar)*t)
}
