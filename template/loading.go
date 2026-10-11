package template

import (
	"sync/atomic"
	"time"
)

// loading is the start-up screen's state: when the wait began, how long it
// may last, and whether it is still up. Every field but the start time is
// atomic because the screen lives across goroutines — Ready arrives from
// whatever goroutine finished loading while the frames keep painting on the
// one that draws, and the two must not tear.
type loading struct {
	// start is when the template was made, which is when the screen first
	// appeared and the moment the timeout counts from. It is written once,
	// before the template is handed to anyone, and only read afterwards.
	start time.Time

	// timeout is the quit rule and the bar's denominator in nanoseconds, so
	// that both count against exactly the same number: the bar reaches 100%
	// at the instant the app would close. Zero or less is never stored —
	// [loading.setTimeout] leaves such a duration alone.
	timeout atomic.Int64

	// ready says the program called Ready, which ends the screen for good:
	// loading happens once, so nothing brings the screen back afterwards.
	ready atomic.Bool

	// disabled says the developer asked to be without the screen at all.
	// It starts as the negation of [loadingEnabledByDefault] when the
	// template is made.
	disabled atomic.Bool

	// pushed says the screen has been drawn and presented as the template's
	// first frame, which is what turning it off afterwards has to wipe.
	pushed atomic.Bool

	// timedOut says the timeout is what closed the app, kept so a test can
	// tell a program that ran out of time from one whose window was closed
	// some other way.
	timedOut atomic.Bool
}
