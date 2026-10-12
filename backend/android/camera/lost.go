//go:build android

package camera

// Lost reports whether the camera was taken away — by another app, or by
// being unplugged. A preview that has been lost produces no more frames and
// has to be closed and opened again.
func (p *Preview) Lost() bool { return p.device.Lost() }
