package antui

// stopped reports whether Close has been called, which is what tells a
// backend's goroutine that the socket closing under it was on purpose.
func (a *Audio) stopped() bool {
	if a == nil {
		return true
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.done
}
