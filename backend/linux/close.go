//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// Close destroys the window and lets go of the connection.
func (x *Driver) Close() {
	if x.conn == nil {
		return
	}
	if x.alive {
		if x.gc != 0 {
			x.send(xFreeGC, 0, le32(x.gc))
		}
		if x.cursor != 0 {
			x.send(xFreeCursor, 0, le32(x.cursor))
		}
		if x.window != 0 {
			x.send(xDestroyWindow, 0, le32(x.window))
		}
	}
	x.conn.Close()
	x.conn = nil
	x.alive = false
	x.keysyms = nil
	x.scratch = nil
	x.rbuf = nil
}
