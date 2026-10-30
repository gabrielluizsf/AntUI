//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// Native is the window's XID and the display it lives on, as DISPLAY spells
// it, for another renderer to open its own connection and draw into the
// window. The XID is zero once the connection is gone.
func (x *Driver) Native() (window uint32, display string) {
	if !x.alive {
		return 0, x.display
	}
	return x.window, x.display
}
