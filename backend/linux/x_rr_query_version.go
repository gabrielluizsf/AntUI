//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// RandR, for the refresh rate. The core protocol has no idea what one is.
const (
	xRRQueryVersion  = 0
	xRRGetScreenInfo = 5
)
