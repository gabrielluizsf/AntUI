//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "time"

// themeTTL is how long a reading stands before the session's environment and
// the desktop's files are read again. A second is far inside the time it
// takes a person to change a theme and look at the result, and short enough
// that they see it land; reading files on every frame would spend more than
// the answer is worth on a state that changes twice a year.
const themeTTL = time.Second
