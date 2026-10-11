package antui

// widgetClick is the press-and-release logic every clickable widget shares.
func (win *Window) widgetClick(id uint32, hovered bool) bool {
	if hovered {
		win.uiHot = id
	}
	switch {
	case win.uiActive == id:
		if win.MouseReleased(MouseLeft) {
			win.uiActive = 0
			return hovered
		}
	case hovered && win.MousePressed(MouseLeft):
		win.uiActive = id
		win.uiFocus = 0
		if win.MouseReleased(MouseLeft) {
			win.uiActive = 0
			return true
		}
	}
	return false
}
