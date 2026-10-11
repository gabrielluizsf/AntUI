//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// Event and reply codes.
const (
	xError            = 0
	xReply            = 1
	xKeyPress         = 2
	xKeyRelease       = 3
	xButtonPress      = 4
	xButtonRelease    = 5
	xMotionNotify     = 6
	xFocusIn          = 9
	xFocusOut         = 10
	xExpose           = 12
	xMapNotify        = 19
	xConfigureNotify  = 22
	xPropertyNotify   = 28
	xSelectionClear   = 29
	xSelectionRequest = 30
	xSelectionNotify  = 31
	xClientMessage    = 33
	xMappingNotify    = 34
)
