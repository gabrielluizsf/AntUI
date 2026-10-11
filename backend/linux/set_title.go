//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// SetTitle names the window in the window manager's decoration.
func (x *Driver) SetTitle(title string) {
	if x.alive {
		x.setTitleOn(title)
	}
}
