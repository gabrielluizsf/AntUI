package antui

// fail records why the device stopped. The first reason is kept: what follows
// a broken pipe is a hundred more broken pipes.
func (a *Audio) fail(err error) {
	if a == nil || err == nil {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.err == nil {
		a.err = err
	}
}
