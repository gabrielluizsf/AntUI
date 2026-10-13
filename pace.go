package antui

import "time"

// pace holds the loop to the budget. The frames land on an absolute grid, so
// the time a frame took to draw does not add to the wait of the next one: a
// frame that overran its budget starts the next frame now and the loop drops
// frames rather than building a debt it can never pay back.
//
// The platform is told when the next frame is due, so a backend that waits on
// its own event queue can stay asleep until then rather than being woken to be
// told there is nothing to do.
func (win *Window) pace() {
	budget := win.frameBudget()
	if budget <= 0 {
		win.handOver(time.Time{})
		return
	}
	now := time.Now()
	if win.nextFrame.IsZero() {
		win.nextFrame = now
	}
	win.nextFrame = win.nextFrame.Add(budget)
	if !win.nextFrame.After(now) {
		win.nextFrame = now
		win.handOver(now)
		return
	}
	win.handOver(win.nextFrame)
	win.sleepUntil(win.nextFrame.Sub(now))
}
