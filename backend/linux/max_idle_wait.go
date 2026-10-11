//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"time"
)

// maxIdleWait caps how long a Pump will sit on the socket. The program's own
// frame rate is the usual answer; this is the backstop for a program that has
// not said, and it keeps a keystroke from waiting on a slow frame rate.
const maxIdleWait = 16 * time.Millisecond
