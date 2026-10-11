package antui

// pushSimple queues an event that carries nothing but its kind.
func (win *Window) PushSimple(t EventType) {
	win.Push(Event{Type: t, Mods: win.mods, X: win.mouseX, Y: win.mouseY})
}
