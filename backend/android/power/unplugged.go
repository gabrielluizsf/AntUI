//go:build android

package power

// What the battery is plugged into.
const (
	Unplugged = 0
	AC        = 1
	USB       = 2
	Wireless  = 4
)
