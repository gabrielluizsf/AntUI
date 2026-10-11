package canvas

// dashPhase is which stretch of the pattern a line starting at offset begins
// in, and how much of that stretch is already gone by the time its first point
// is reached.
func dashPhase(pattern []float64, offset, period float64) (int, float64) {
	phase, gone := 0, wrap(offset, period)
	for {
		phase %= len(pattern)
		if gone < pattern[phase] {
			break
		}
		gone -= pattern[phase]
		phase++
	}
	return phase, gone
}
