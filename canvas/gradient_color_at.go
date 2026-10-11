package canvas

// gradientColorAt interpolates the colour at one point along a gradient line.
// Points outside the stops clamp to the end colours.
func gradientColorAt(stops []GradientStop, t float64) Color {
	if t <= 0 {
		return stops[0].Color
	}
	if t >= 1 {
		return stops[len(stops)-1].Color
	}
	for i := 1; i < len(stops); i++ {
		if t < stops[i].Offset {
			a, b := stops[i-1], stops[i]
			span := b.Offset - a.Offset
			var f float32
			if span > 0 {
				f = float32((t - a.Offset) / span)
			}
			return Mix(a.Color, b.Color, f)
		}
	}
	return stops[len(stops)-1].Color
}
