//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"time"
)

// fileStamp is a file's identity without opening it — whether it is there,
// how big it is, when it was last written — which is all it takes to notice
// that the desktop wrote a new answer into one.
type fileStamp struct {
	size int64
	mod  time.Time
	ok   bool
}
