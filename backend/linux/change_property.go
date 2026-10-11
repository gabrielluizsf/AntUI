//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// changeProperty sets a property on the window. count is in items of the
// size named by format, not in bytes.
func (x *Driver) changeProperty(property, typ uint32, format byte, value []byte, count uint32) {
	x.changePropertyOn(x.window, property, typ, format, value, count)
}
