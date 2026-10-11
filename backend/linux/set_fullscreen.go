//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// SetFullscreen asks the window manager for a full-screen window, reporting
// whether the request could be made at all.
func (x *Driver) SetFullscreen(on bool) bool {
	if !x.alive || x.atomNetWMStateFullscreen == 0 {
		return false // no EWMH here
	}
	action := uint32(0)
	if on {
		action = 1
	}
	return x.wmState(action, x.atomNetWMStateFullscreen)
}
