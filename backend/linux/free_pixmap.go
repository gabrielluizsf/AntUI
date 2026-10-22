//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// freePixmap gives a pixmap back to the server.
func (x *Driver) freePixmap(pid uint32) {
	x.send(xFreePixmap, 0, le32(pid))
}
