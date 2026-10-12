//go:build android

package display

// Display is one screen: the built-in one, or a television the device is
// casting to, or a monitor it is plugged into.
type Display struct {
	// ID is 0 for the built-in screen and something else for the rest.
	ID   int
	Name string
	// Width and Height are the whole screen in pixels, including whatever
	// the system bars are covering.
	Width, Height int
	// Density is dots per inch as Android buckets it; 160 is one pixel per
	// device-independent pixel.
	Density int
	// Refresh is how many frames a second it shows. It is a float because a
	// panel that says 60 is usually running at 59.94 or 60.0000047, and
	// rounding it before the caller sees it loses the only interesting part.
	Refresh float32
	// HDR is whether it can show more than eight bits a channel. It is
	// false below API 24, where there was no way to ask.
	HDR bool
}
