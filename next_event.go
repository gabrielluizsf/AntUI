package antui

// NextEvent pops the next event off this frame's queue, reporting false when
// it is empty.
func (win *Window) NextEvent() (Event, bool) {
	if len(win.queue) == 0 {
		return Event{}, false
	}
	ev := win.queue[0]
	win.queue = win.queue[1:]
	return ev, true
}
