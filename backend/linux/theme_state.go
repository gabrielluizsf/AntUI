//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"time"
)

// themeState is what the desktop said the last time it was asked, and when
// it is due to be asked again.
type themeState struct {
	dark, known bool
	setting     darkAnswer // what the settings tool said, as of stamp
	stamp       fileStamp  // the file that answer came from
	next        time.Time  // when the cheap sources are read again
}
