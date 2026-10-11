//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"errors"
	"net"
	"time"
)

// fill reads whatever the server has sent. A zero timeout polls: it takes
// what has already arrived and returns rather than waiting for more.
func (x *Driver) fill(timeout time.Duration) bool {
	if !x.alive {
		return false
	}
	x.makeRoom()
	if timeout <= 0 {
		return x.fillNow()
	}

	x.conn.SetReadDeadline(time.Now().Add(timeout))
	at := len(x.rbuf)
	n, err := x.conn.Read(x.rbuf[at:cap(x.rbuf)])
	x.rbuf = x.rbuf[:at+n]

	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return n > 0
		}
		// Anything else — EOF included — means the server is gone.
		x.die("blocking read: " + err.Error())
		return false
	}
	return n > 0
}
