//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "time"

// fillNow takes whatever has already arrived, without waiting for more.
//
// It reads the socket directly rather than through a deadline, because a
// deadline in the past does not mean "do not wait" to Go's poller — it means
// "fail immediately", and the read never happens. Getting that wrong is
// silent: drawing still works, because writing is a separate path, and only
// input disappears.
func (x *Driver) fillNow() bool {
	at := len(x.rbuf)
	into := x.rbuf[at:cap(x.rbuf)]

	// Clear any deadline left over from a blocking read. Both the poller and
	// RawConn.Read honour it, so one set during the handshake and never
	// cleared makes every later read fail the moment it passes — and the
	// connection would look dead five seconds in for no reason at all.
	x.conn.SetReadDeadline(time.Time{})

	if x.raw == nil {
		return x.fillBlocking(into, at)
	}
	return x.fillRaw(into, at)
}
