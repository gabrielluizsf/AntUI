//go:build !windows

package antui

import "time"

// sleepFor is the platform-independent frame wait. Windows supplies a
// high-resolution waitable-timer implementation because its regular sleep can
// round a 60 Hz frame up by several milliseconds.
func sleepFor(d time.Duration) {
	time.Sleep(d)
}
