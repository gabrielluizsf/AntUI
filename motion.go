//go:build android

package antui

import "github.com/gabrielluizsf/antui/backend/android/ndk"

// motion sorts a pointer event by what produced it. The low bits of every
// pointer source are the same, so this cannot be a switch on the source: a
// mouse has to be recognised by its own bits before anything else, and
// everything left that points is treated as a finger.
func (b *androidWindow) motion(win *Window, e *ndk.InputEvent) bool {
	if e.Source().Is(ndk.SourceMouse) {
		return b.mouse(win, e)
	}
	if e.Source().Is(ndk.SourceJoystick) {
		// A controller's sticks arrive as motion on axes rather than as a
		// position. The buttons already work, through key events; the axes
		// need an API antui does not have yet.
		return false
	}
	return b.touch(win, e)
}
