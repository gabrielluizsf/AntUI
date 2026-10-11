//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// announceDrops says the window will take a drag. Without this property no
// source will offer it anything at all.
func (x *Driver) announceDrops() {
	x.dropAtoms()
	if x.atomXdndAware == 0 {
		return
	}
	x.changeProperty(x.atomXdndAware, atomAtom, 32, le32(xdndVersion), 1)
}
