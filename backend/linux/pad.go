//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// pad is how many bytes of padding follow a field of this length, the
// protocol requiring everything to sit on a four-byte boundary.
func pad(length int) int { return (4 - length%4) % 4 }
