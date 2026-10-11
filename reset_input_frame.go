package antui

// resetInputFrame clears the per-frame input state.
func (win *Window) resetInputFrame() {
	win.mouseDownFrame = [mouseCount]bool{}
	win.mouseUpFrame = [mouseCount]bool{}
	win.keyDownFrame = [keyCount]bool{}
	win.keyUpFrame = [keyCount]bool{}
	win.wheel = 0
	win.mouseDX, win.mouseDY = 0, 0
	win.text.Reset()
	win.dropped = win.dropped[:0]
	win.beginTouchFrame()
	win.queue = win.queue[:0]
	win.droppedEvents = 0
	win.uiHot = 0
	win.uiTab = win.uiTab[:0]
	win.uiBlink += win.delta
	win.frontmost = win.frontmost[:0]
}
