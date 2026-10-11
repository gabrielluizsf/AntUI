package antui

// Push folds an event into the input state and queues it for NextEvent.
func (win *Window) Push(ev Event) {
	win.applyEvent(ev)

	if len(win.queue) >= maxEvents {
		win.droppedEvents++
		return
	}
	win.queue = append(win.queue, ev)
}
