package antui

// ReleaseModal lets go of the claim [Window.ClaimModal] made, if id is still
// the one holding it — closing a popup that was never the current claim,
// because another one opened since, leaves that one alone.
func (win *Window) ReleaseModal(id uint32) {
	if win.modalID == id {
		win.modalID = 0
	}
}
