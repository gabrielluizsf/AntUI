//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"errors"
	"net"
	"time"
)

// fillBlocking waits for a moment when there is no raw access to the socket.
// Slower to notice input, but it does notice it.
func (x *Driver) fillBlocking(into []byte, at int) bool {
	x.conn.SetReadDeadline(time.Now().Add(time.Millisecond))
	defer x.conn.SetReadDeadline(time.Time{})
	n, err := x.conn.Read(into)
	x.rbuf = x.rbuf[:at+n]
	if err != nil {
		var netErr net.Error
		if !errors.As(err, &netErr) || !netErr.Timeout() {
			x.die("fallback read: " + err.Error())
			return false
		}
	}
	return n > 0
}
