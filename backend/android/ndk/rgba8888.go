//go:build android

package ndk

// The formats a window buffer can be asked for.
const (
	RGBA8888 Format = 1
	RGBX8888 Format = 2 // the same bytes, alpha ignored
	RGB888   Format = 3 // three bytes, no padding; rare and slow
	RGB565   Format = 4
)
