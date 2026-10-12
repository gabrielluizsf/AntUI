//go:build android

package app

// RestoredState is what a previous run of this activity saved, or nil.
func (a *App) RestoredState() []byte {
	if a == nil {
		return nil
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.saved
}
