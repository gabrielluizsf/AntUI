//go:build windows

package antui

// A guid is what COM calls everything by.
type guid struct {
	a uint32
	b uint16
	c uint16
	d [8]byte
}
