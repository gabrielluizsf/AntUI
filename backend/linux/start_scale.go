//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// startScale reads the display's scale once and watches the resource database
// for the write that means the user changed it.
func (x *Driver) startScale() {
	x.watchScale()
	x.scale, x.scaleKnown = x.askScale()
}
