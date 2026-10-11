package template

// gridSpread divides what one item wants between the tracks it spans, giving
// each the same share of what is missing.
func gridSpread(values []int, value int) {
	if len(values) == 0 {
		return
	}
	used := 0
	for _, current := range values {
		used += current
	}
	if used < value {
		extra := (value - used + len(values) - 1) / len(values)
		for i := range values {
			values[i] += extra
		}
	}
}
