//go:build android

package antui

// displayScale is the density over 160, which is the number Android itself
// multiplies a size in device-independent pixels by. Android hands a new one
// over on a configuration change, so reading it every frame is what makes
// that reach the window.
func (b *androidWindow) displayScale() (float64, bool) { return scaleOrNone(b.scale) }
