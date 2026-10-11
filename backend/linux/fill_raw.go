//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"errors"
	"os"
	"syscall"
)

// fillRaw reads straight from the socket once, never waiting for readability.
func (x *Driver) fillRaw(into []byte, at int) bool {
	var n int
	var readErr error
	controlErr := x.raw.Read(func(fd uintptr) bool {
		for {
			var err error
			n, err = syscall.Read(int(fd), into)
			if err == syscall.EINTR {
				continue
			}
			readErr = err
			return true // never wait for readability
		}
	})
	if controlErr != nil {
		// A deadline that slipped past is not a dead server, only a read
		// that did not happen.
		if errors.Is(controlErr, os.ErrDeadlineExceeded) {
			return false
		}
		x.die("raw control: " + controlErr.Error())
		return false
	}

	switch {
	case readErr == syscall.EAGAIN || readErr == syscall.EWOULDBLOCK:
		return false // nothing waiting, and nothing wrong
	case readErr != nil:
		x.die("raw read: " + readErr.Error())
		return false
	case n == 0:
		x.die("the server closed the connection")
		return false
	}
	x.rbuf = x.rbuf[:at+n]
	return true
}
