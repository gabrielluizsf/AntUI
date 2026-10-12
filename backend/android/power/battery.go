//go:build android

package power

// Battery is what the device says about its own power.
type Battery struct {
	// Level is 0 to 1, or -1 when the device will not say.
	Level float32
	// Charging is whether it is going up.
	Charging bool
	// Full is charging that has finished.
	Full bool
	// Plugged is what it is plugged into: [Unplugged], [AC], [USB] or
	// [Wireless].
	Plugged int
}
