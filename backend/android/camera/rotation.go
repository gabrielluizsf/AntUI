//go:build android

package camera

// Rotation is how far each frame is being turned, in degrees clockwise. It
// is 0 unless [Options.Upright] was set.
func (p *Preview) Rotation() int { return p.upright }
