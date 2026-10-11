package antui

// Close stops the sound. It is safe to call twice.
func (a *Audio) Close() error {
	if a == nil {
		return nil
	}
	a.mu.Lock()
	if a.done {
		a.mu.Unlock()
		return nil
	}
	a.done = true
	drv := a.drv
	a.mu.Unlock()

	if drv == nil {
		return nil
	}
	return drv.close()
}
