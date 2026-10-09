package css

func roundLerp(a, b int, t float64) int {
	return int(float64(a) + float64(b-a)*t + 0.5)
}
