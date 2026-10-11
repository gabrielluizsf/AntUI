//go:build darwin && cgo

package macos

// ContentScale is pixels per point. On macOS it is 1, and that is the true
// answer rather than a missing one.
//
// The view reports its bounds in points and the framebuffer is that size, so
// a window asked for at 1280 by 720 is 1280 by 720 pixels of Canvas and
// appears at 1280 by 720 points on every Mac — the postage-stamp problem the
// other two systems have simply does not arise here, because AppKit does the
// scaling on the way to the glass.
//
// What it costs is sharpness: on a retina display the frame is resampled up
// rather than drawn at the panel's own resolution. Drawing at
// convertRectToBacking size would fix that and double the pixels pushed
// every frame, which is a decision for a machine that can measure it.
func (d *Driver) ContentScale() float64 { return 1 }
