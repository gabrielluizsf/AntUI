package css

// frameSpan finds the two frames straddling a progress t in [0,1], with the
// local progress between them. t before the first frame pins to the first,
// after the last pins to the last.
func frameSpan(frames []Keyframe, t float64) (lo, hi Keyframe, local float64, ok bool) {
	if len(frames) == 0 {
		return Keyframe{}, Keyframe{}, 0, false
	}
	if t <= frames[0].Offset || len(frames) == 1 {
		return frames[0], frames[0], 0, true
	}
	last := frames[len(frames)-1]
	if t >= last.Offset {
		return last, last, 1, true
	}
	for i := 0; i < len(frames)-1; i++ {
		a, b := frames[i], frames[i+1]
		if t >= a.Offset && t <= b.Offset {
			local := 0.0
			if span := b.Offset - a.Offset; span > 0 {
				local = (t - a.Offset) / span
			}
			return a, b, local, true
		}
	}
	return frames[0], frames[0], 0, true
}
