package antui

// pushTouch folds a touch event into the window's touch state. It is called
// from push, so the state is right whether or not the queue overflowed.
func (win *Window) pushTouch(ev Event) {
	win.sawTouch = true
	i := win.touchIndex(ev.TouchID)
	switch ev.Type {
	case EventTouchDown:
		// A second down for a finger already on the screen is the system
		// repeating itself, which happens. Taking it as a move keeps the
		// history rather than starting a new one.
		if i >= 0 {
			break
		}
		win.touches = append(win.touches, Touch{
			ID:       ev.TouchID,
			X:        ev.X,
			Y:        ev.Y,
			StartX:   ev.X,
			StartY:   ev.Y,
			Pressure: ev.Pressure,
			Size:     ev.TouchSize,
			Tool:     ev.Tool,
			Began:    true,
		})
		return
	case EventTouchUp, EventTouchCancel:
		if i < 0 {
			return
		}
		win.touches[i].Ended = true
		win.touches[i].Cancelled = ev.Type == EventTouchCancel
	}
	if i < 0 {
		// A move for a finger that was never reported down. It happens when
		// the app gains a surface mid-gesture; treating it as a landing is
		// better than dropping the finger on the floor.
		win.touches = append(win.touches, Touch{
			ID: ev.TouchID, StartX: ev.X, StartY: ev.Y, Began: true,
			Tool: ev.Tool,
		})
		i = len(win.touches) - 1
	}
	t := &win.touches[i]
	t.DX += ev.X - t.X
	t.DY += ev.Y - t.Y
	t.X, t.Y = ev.X, ev.Y
	if ev.Pressure > 0 {
		t.Pressure = ev.Pressure
	}
	if ev.TouchSize > 0 {
		t.Size = ev.TouchSize
	}
}
