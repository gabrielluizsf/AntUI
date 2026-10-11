package antui

// Err is what went wrong after the device was opened — a card unplugged, a
// server that went away — or nil.
//
// A device that dies does not take the game with it. Everything goes on being
// mixed and nothing is heard, which is what happens on a machine with the
// speakers turned off anyway.
func (a *Audio) Err() error {
	if a == nil {
		return nil
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.err
}
