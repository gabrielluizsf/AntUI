package antui

// Events iterates this frame's queue, draining it as it goes.
func (win *Window) Events(yield func(Event) bool) {
	for {
		ev, ok := win.NextEvent()
		if !ok || !yield(ev) {
			return
		}
	}
}
