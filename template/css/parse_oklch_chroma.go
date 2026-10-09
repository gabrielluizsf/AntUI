package css

// parseOKLCHChroma reads oklch chroma: a number, or a percentage where 100%
// is 0.4.
func parseOKLCHChroma(f string) (float64, error) {
	return parseOKLabChroma(f)
}
