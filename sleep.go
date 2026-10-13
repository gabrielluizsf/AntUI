package antui

import "time"

// Sleep pauses for a number of seconds.
func Sleep(seconds float64) {
	if seconds > 0 {
		time.Sleep(time.Duration(seconds * float64(time.Second)))
	}
}
