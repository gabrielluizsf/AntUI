//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// freeGC gives a graphics context back to the server.
func (x *Driver) freeGC(gc uint32) {
	x.send(xFreeGC, 0, le32(gc))
}
