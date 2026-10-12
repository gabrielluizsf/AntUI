package apk

import (
	_ "image/jpeg"
)

// The launcher icon at each density, in pixels. They are 48 device-
// independent pixels multiplied out, which is the size every launcher asks
// for and has since the beginning.
var densities = []struct {
	name string
	size int
	// adaptive is the same 48dp icon's *adaptive* layer, which is 108dp —
	// larger than what is shown, because the launcher crops it to whatever
	// shape that device uses and may move it about while it does.
	adaptive int
}{
	{"mdpi", 48, 108},
	{"hdpi", 72, 162},
	{"xhdpi", 96, 216},
	{"xxhdpi", 144, 324},
	{"xxxhdpi", 192, 432},
}
