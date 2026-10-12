//go:build android

package location

import (
	"time"
)

// Options is how a watch is opened.
type Options struct {
	// Provider is [GPS] or [Network]. Empty asks for both, and reports from
	// whichever answers.
	Provider string
	// Interval is the shortest gap between reports. It is a floor the
	// platform is allowed to exceed and, on API 29 and up, to ignore
	// downwards — the system throttles background apps whatever they ask.
	Interval time.Duration
	// Distance is how far the device must move before it is worth another
	// report, in metres. Zero reports on time alone.
	Distance float32
}
