//go:build android

package antui

import (
	"github.com/gabrielluizsf/antui/backend"
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/ndk"
	"time"
)

// The Android backend. A phone has no window manager to ask for anything: the
// system decides the size, when there is a surface at all, and when the app
// may draw. So most of the backend interface is a polite refusal, and the
// two calls that do the work are pump — which waits when there is nothing to
// draw on — and present, which is a blit into a buffer the compositor lends
// out one frame at a time.
//
// The Canvas above this is the same Canvas as everywhere else, so a program
// written for a desktop draws on a phone without knowing it moved.
type androidWindow struct {
	activity *app.App
	surface  *ndk.Window
	format   ndk.Format
	scale    float64

	// primary is the touch that is also reported as the left mouse button,
	// or -1 when no finger is down. buttons is the last mouse button state
	// seen, which is what a change is worked out against.
	primary int
	buttons int32
	// refresh is the display's rate, asked once. Zero means not yet asked
	// and a negative means asked and not answered.
	refresh int
	// dark is whether the device is in night mode. It comes from the same
	// configuration the density does, so it is read once at open and again
	// on every configuration change rather than asked per frame.
	dark bool
}
