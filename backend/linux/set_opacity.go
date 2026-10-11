//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// SetOpacity fades the whole window with _NET_WM_WINDOW_OPACITY, so what
// sits behind it shows through; 255 is fully opaque. The property tells a
// compositor to paint the window at that opacity, drawing nothing differently
// itself. It reports whether the request could be made: without a compositor
// the property is written, but no one is listening — on plain X11 the window
// stays as it is because there is nothing behind it to show.
func (x *Driver) SetOpacity(alpha uint8) bool {
	if !x.alive || x.atomNetWMOpacity == 0 {
		return false // no EWMH here
	}
	// The property is a 32-bit ARGB mask; 0xffffffff is fully opaque and
	// each step of alpha is one of those bytes.
	value := uint32(alpha) * 0x01010101
	x.changeProperty(x.atomNetWMOpacity, atomCardinal, 32, le32(value), 1)
	return x.alive
}
