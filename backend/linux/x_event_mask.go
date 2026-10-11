//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// The events the window asks the server for: key press/release, button
// press/release, pointer motion, exposure, structure notify and focus change.
const xEventMask = 0x1 | 0x2 | 0x4 | 0x8 | 0x40 | 0x8000 | 0x20000 | 0x200000
