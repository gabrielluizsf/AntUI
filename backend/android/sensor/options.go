//go:build android

package sensor

import (
	"time"
)

// Options is how a stream is opened.
type Options struct {
	// Rate is how often to report. Zero asks for the sensor's own default.
	// It is a ceiling, not a promise: the platform delivers at that rate or
	// slower, and never faster than the sensor's own minimum.
	Rate time.Duration
	// Latency lets the platform hold readings and deliver them in a burst,
	// which is what lets the processor sleep between them. Zero delivers each
	// as it happens. It needs API 26; below that it is ignored.
	Latency time.Duration
	// Wake asks for a sensor that keeps the device awake to report. Almost
	// nothing should: a step counter that does not wake the device reports
	// everything it collected when the device next wakes, which is the same
	// answer for a fraction of the battery.
	Wake bool
	// Buffer is how many readings [Stream.Events] keeps when nothing is
	// reading them. The oldest are dropped. Zero means no channel at all,
	// and only [Stream.Latest].
	Buffer int
}
