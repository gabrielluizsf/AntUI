//go:build android

package antui

// setTitle: an Android app's name is in its manifest, and the label in the
// recents list is set through Java. Nothing here can change it.
func (b *androidWindow) setTitle(string) {}
