//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// The event-mask bit inside a ChangeWindowAttributes value mask: the same
// number a window of our own is given as 0x800 when it asks for its events.
// An event mask is what one connection asks to be told about, and it counts
// for that connection only.
const cwEventMask = 0x800
