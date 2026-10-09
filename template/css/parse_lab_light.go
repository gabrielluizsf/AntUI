package css

// parseLabLight reads lab/lch lightness: 0-100 as a number or a percentage.
func parseLabLight(f string) (float64, error) {
	return parsePercent(f)
}
