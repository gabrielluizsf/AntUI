package antui

import "math"

// trackMove works out whether the finger has wandered past the tap slop, held
// long enough for a long press, or is dragging.
func (win *Window) trackMove(now float64, t *Touch) {
	g := &win.gesture

	dx := float64(t.X - g.startX)
	dy := float64(t.Y - g.startY)
	if !g.moved && math.Hypot(dx, dy) > win.slop(TapSlop) {
		g.moved = true
	}

	// A long press fires while the finger is still down, once.
	if !g.moved && !g.longFired && now-g.startTime >= LongPressTime {
		g.longFired = true
		g.out.LongPress = true
	}

	if g.moved && (t.DX != 0 || t.DY != 0) {
		g.out.Drag = true
		g.out.DX, g.out.DY = t.DX, t.DY
	}
}
