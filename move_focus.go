package antui

// moveFocus walks the frame's focus order when Tab or Shift+Tab is pressed.
// It runs at the top of End, no matter how the frame was drawn, so the same
// frame that moved the focus draws the ring.
func (win *Window) moveFocus() {
	if !win.KeyPressed(KeyTab) {
		return
	}
	if len(win.uiTab) == 0 {
		return
	}
	shift := win.mods&ModShift != 0
	for i, id := range win.uiTab {
		if win.uiFocus != id {
			continue
		}
		if shift {
			i = (i - 1 + len(win.uiTab)) % len(win.uiTab)
		} else {
			i = (i + 1) % len(win.uiTab)
		}
		win.uiFocus = win.uiTab[i]
		win.uiCursor = 0
		win.uiBlink = 0
		return
	}
	if shift {
		win.uiFocus = win.uiTab[len(win.uiTab)-1]
	} else {
		win.uiFocus = win.uiTab[0]
	}
	win.uiCursor = 0
	win.uiBlink = 0
}
