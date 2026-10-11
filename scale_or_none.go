package antui

// scaleOrNone is a platform's own scale as the raw answer: a positive
// number is what the display said, and zero is that it did not — the
// difference DisplayScale keeps and ContentScale stands a one in for.
func scaleOrNone(scale float64) (float64, bool) {
	if scale > 0 {
		return scale, true
	}
	return 0, false
}
