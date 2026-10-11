//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

import (
	x11 "github.com/gabrielluizsf/antui/backend/linux"
)

// x11Native is the Linux backend, in two halves: this struct adapts the
// window core's calls to the platform package, and the package itself lives
// in backend/linux.
type x11Native struct{ d *x11.Driver }
