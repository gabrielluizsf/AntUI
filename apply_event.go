package antui

// applyEvent folds an event into the input state.
func (win *Window) applyEvent(ev Event) {
	switch ev.Type {
	case EventClose:
		win.shouldClose = true
	case EventMouseMove:
		win.mouseDX += ev.X - win.mouseX
		win.mouseDY += ev.Y - win.mouseY
		win.mouseX, win.mouseY = ev.X, ev.Y
	case EventMouseDown:
		win.mouseX, win.mouseY = ev.X, ev.Y
		if ev.Button >= 0 && ev.Button < mouseCount {
			win.mouseState[ev.Button] = true
			win.mouseDownFrame[ev.Button] = true
		}
		win.mods = ev.Mods
	case EventMouseUp:
		win.mouseX, win.mouseY = ev.X, ev.Y
		if ev.Button >= 0 && ev.Button < mouseCount {
			win.mouseState[ev.Button] = false
			win.mouseUpFrame[ev.Button] = true
		}
		win.mods = ev.Mods
	case EventMouseWheel:
		win.wheel += ev.Wheel
	case EventKeyDown:
		if ev.Key > 0 && ev.Key < keyCount {
			win.keyState[ev.Key] = true
			win.keyDownFrame[ev.Key] = true
		}
		win.mods = ev.Mods
	case EventKeyUp:
		if ev.Key > 0 && ev.Key < keyCount {
			win.keyState[ev.Key] = false
			win.keyUpFrame[ev.Key] = true
		}
		win.mods = ev.Mods
	case EventText:
		if win.text.Len()+len(ev.Text) < textBuffer {
			win.text.WriteString(ev.Text)
		}
	case EventDropFiles:
		win.mouseX, win.mouseY = ev.X, ev.Y
		win.dropped = append(win.dropped, ev.Files...)
	case EventExpose:
		win.fullRedraw = true
	case EventFocus:
		if !ev.Focused {
			win.focusLost()
		}
	case EventTouchDown, EventTouchMove, EventTouchUp, EventTouchCancel:
		win.pushTouch(ev)
	}
}
