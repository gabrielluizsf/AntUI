package svg

// saturateMatrix, hueRotateMatrix and luminanceMatrix are the three shortcuts
// of feColorMatrix — `saturate`, `hueRotate` and `luminanceToAlpha` — written
// out as the same twenty numbers `matrix` would take, since that is all a
// matrix ever is here.
func saturateMatrix(s float64) [20]float64 {
	return [20]float64{
		0.213 + 0.787*s, 0.715 - 0.715*s, 0.072 - 0.072*s, 0, 0,
		0.213 - 0.213*s, 0.715 + 0.285*s, 0.072 - 0.072*s, 0, 0,
		0.213 - 0.213*s, 0.715 - 0.715*s, 0.072 + 0.928*s, 0, 0,
		0, 0, 0, 1, 0,
	}
}
