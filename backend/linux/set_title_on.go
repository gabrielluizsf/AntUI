//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// setTitleOn writes the title twice: WM_NAME for window managers that only
// read Latin-1, and _NET_WM_NAME in UTF-8 for everything from this century.
func (x *Driver) setTitleOn(title string) {
	latin1 := make([]byte, 0, len(title))
	for _, r := range title {
		if r <= 0xFF {
			latin1 = append(latin1, byte(r))
		} else {
			latin1 = append(latin1, '?')
		}
	}
	x.changeProperty(atomWMName, atomString, 8, latin1, uint32(len(latin1)))
	if x.atomNetWMName != 0 && x.atomUTF8String != 0 {
		x.changeProperty(x.atomNetWMName, x.atomUTF8String, 8, []byte(title), uint32(len(title)))
	}
}
