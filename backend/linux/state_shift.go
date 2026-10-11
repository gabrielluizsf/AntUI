//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// The state mask the server sends with a key or button event.
const (
	stateShift  = 0x1
	stateLock   = 0x2 // caps lock
	stateCtrl   = 0x4
	stateAlt    = 0x8
	stateSuper  = 0x40
	stateAltGr  = 0x80
	stateGroup2 = 0x2000
)
