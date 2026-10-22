//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// freeCursor gives a cursor back to the server.
func (x *Driver) freeCursor(cursor uint32) {
	x.send(xFreeCursor, 0, le32(cursor))
}
