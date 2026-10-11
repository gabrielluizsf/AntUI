package antui

// recognizeOne watches the first finger down for a tap, a long press, a drag
// and a fling.
func (win *Window) recognizeOne(now float64) {
	g := &win.gesture
	t := win.trackFinger(now)
	if t == nil {
		return
	}
	g.out.X, g.out.Y = t.X, t.Y

	win.trackVelocity(now, t)
	win.trackMove(now, t)

	if !t.Ended {
		return
	}
	win.finishOne(now, t)
}
