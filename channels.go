package antui

// Channels is how many there are in a frame. Two, everywhere, for now.
func (a *Audio) Channels() int {
	if a == nil {
		return 0
	}
	return a.channels
}
