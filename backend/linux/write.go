//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"errors"
	"net"
	"time"
)

// write sends a request, listening to the server as it goes.
//
// The listening is the point. A frame of pixels is over a megabyte, and while
// we push it the server is pushing events back; if neither side reads, both
// block on a full socket and the program hangs with a half-drawn window. So
// every chunk is preceded by a drain of whatever has arrived, which is what
// keeps the two of us from waiting on each other.
func (x *Driver) write(data []byte) bool {
	deadline := time.Now().Add(xReplyTimeout)

	for len(data) > 0 && x.alive {
		// Listen before pushing more at it.
		for x.fill(0) {
		}

		chunk := min(len(data), writeChunk)
		x.conn.SetWriteDeadline(time.Now().Add(xReplyTimeout))
		n, err := x.conn.Write(data[:chunk])
		data = data[n:]

		if n > 0 {
			deadline = time.Now().Add(xReplyTimeout)
		}
		if err == nil {
			continue
		}
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() && time.Now().Before(deadline) {
			// The server is not taking it yet. That is only fatal if it
			// never does.
			continue
		}
		x.die("write: " + err.Error())
		return false
	}
	return x.alive
}
