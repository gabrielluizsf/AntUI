//go:build android

package antui

func scaleFromDensity(dpi int) float64 {
	if dpi <= 0 {
		return 0
	}
	return float64(dpi) / 160
}
