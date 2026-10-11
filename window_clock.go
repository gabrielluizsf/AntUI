package antui

import "time"

// windowClock is the frame pacing state.
type windowClock struct {
	delta     float64
	nextFrame time.Time // when the next frame is due; see pace
	targetFPS int       // 0 follows the display, see SetFPS
	refreshHz int       // the display's rate, 0 until asked, -1 when unknown
	sleep     func(time.Duration)
}
