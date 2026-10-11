package antui

// touchIndex finds a live touch, or -1.
func (win *Window) touchIndex(id int) int {
	for i := range win.touches {
		if win.touches[i].ID == id {
			return i
		}
	}
	return -1
}
