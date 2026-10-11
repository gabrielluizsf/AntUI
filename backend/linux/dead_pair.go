//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// deadPair keys the compose table by the accent and the letter it lands on.
type deadPair struct {
	dead uint32
	base rune
}
