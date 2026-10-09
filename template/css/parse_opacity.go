package css

func parseOpacity(raw string) (float64, bool) {
	var v float64
	if _, err := scanFloat(raw, &v); err != nil {
		return 0, false
	}
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	return v, true
}
