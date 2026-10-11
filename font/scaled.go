package font

// Scaled is how many pixels one font unit is at a size in pixels per em.
func (f *TTF) Scaled(pixels float64) float64 {
	return pixels / float64(f.UnitsPerEm)
}
