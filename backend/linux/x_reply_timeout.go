//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "time"

// How long to wait for a reply while opening. A server that has not answered
// in five seconds is not going to.
const xReplyTimeout = 5 * time.Second
