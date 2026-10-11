package template

// markReady ends the screen. It may be called from any goroutine, which is
// the point: the goroutine that finished loading says so, and the next frame
// notices.
func (l *loading) markReady() {
	if l != nil {
		l.ready.Store(true)
	}
}
