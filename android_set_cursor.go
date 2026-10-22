//go:build android

package antui

// setCursor: a touch screen has no pointer to give a picture to.
func (b *androidWindow) setCursor(CursorIcon) bool { return false }
