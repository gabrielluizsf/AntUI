//go:build android

package antui

// setOpacity: a phone has no window manager to ask for a fade, and the
// window covers the screen anyway — there is nothing to see behind it.
func (b *androidWindow) setOpacity(uint8) bool { return false }
