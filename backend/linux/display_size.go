//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// DisplaySize is the size of the display the window is on.
func (x *Driver) DisplaySize() (w, h int, ok bool) {
	if x.screenWidth == 0 || x.screenHeight == 0 {
		return 0, 0, false
	}
	return int(x.screenWidth), int(x.screenHeight), true
}
