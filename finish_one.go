package antui

import "math"

// finishOne is the lift: it turns a released finger into a fling or a tap,
// counting double and triple taps.
func (win *Window) finishOne(now float64, t *Touch) {
	g := &win.gesture
	g.tracking = false

	// A cancelled touch is the system taking the gesture away — a
	// notification pulled down over the app. It is not a tap and it is not a
	// throw.
	if t.Cancelled {
		return
	}

	speed := math.Hypot(g.vx, g.vy)
	switch {
	case g.moved:
		if speed >= win.slop(FlingSpeed) {
			g.out.Fling = true
			g.out.VX, g.out.VY = g.vx, g.vy
		}
	case now-g.startTime < LongPressTime:
		g.out.Tap = true
		near := math.Hypot(float64(t.X-g.tapX), float64(t.Y-g.tapY)) <= win.slop(DoubleTapSlop)
		if near && now-g.tapTime <= DoubleTapTime {
			g.tapCount++
		} else {
			g.tapCount = 1
		}
		g.out.Count = g.tapCount
		g.tapTime, g.tapX, g.tapY = now, t.X, t.Y
	}
}
