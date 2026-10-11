package template

// setDisabled switches the screen off or back on. Coming back on does not
// rewind the clock: the timeout still counts from when the template was
// made, because that is when the screen first appeared.
func (l *loading) setDisabled(disabled bool) {
	if l != nil {
		l.disabled.Store(disabled)
	}
}
