//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// makeRoom compacts what has been consumed and leaves space to read into, so
// the buffer does not grow just because the read position walked forward.
func (x *Driver) makeRoom() {
	if x.rpos > 0 {
		x.rbuf = x.rbuf[:copy(x.rbuf, x.rbuf[x.rpos:])]
		x.rpos = 0
	}
	if at := len(x.rbuf); cap(x.rbuf)-at < 8192 {
		grown := make([]byte, at, max(cap(x.rbuf)*2, at+8192))
		copy(grown, x.rbuf)
		x.rbuf = grown
	}
}
