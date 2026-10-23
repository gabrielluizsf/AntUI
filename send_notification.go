package antui

// SendNotification shows a system notification with a title and a body —
// a desktop toast on Linux, Windows and macOS, a notification in the shade
// on Android. It reports whether the system took it, not whether the user
// saw it: notifications the user turned off are taken and never shown.
func (win *Window) SendNotification(title, body string) bool {
	if win == nil || win.native == nil {
		return false
	}
	return win.native.sendNotification(title, body)
}
