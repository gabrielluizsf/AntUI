package antui

import "time"

// Now is the seconds since the library was first asked the time, on a
// monotonic clock that no change to the system clock can move.
func Now() float64 {
	clockOnce.Do(func() { clockZero = time.Now() })
	return time.Since(clockZero).Seconds()
}
