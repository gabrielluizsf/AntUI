//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// The one the root window is asked for instead: property change. The window
// does not want it — it would be told about every hint the window manager
// wrote back — but the root is where a desktop keeps the resource database
// and rewrites it when the display's scale changes, and that write is the
// only signal there is that the number did.
const xPropertyChangeMask = 1 << 22
