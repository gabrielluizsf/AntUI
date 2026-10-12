//go:build android

package app

// Scale is the factor a size in points is multiplied by, or 0 when the
// platform will not say.
func (c Config) Scale() float64 {
	if c.Density <= 0 {
		return 0
	}
	return float64(c.Density) / 160
}
