package antui

// TouchByID finds a touch by the identifier that stays with it.
func (win *Window) TouchByID(id int) (Touch, bool) {
	for _, t := range win.touches {
		if t.ID == id {
			return t, true
		}
	}
	return Touch{}, false
}
