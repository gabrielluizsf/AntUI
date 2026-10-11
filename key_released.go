package antui

// KeyReleased reports whether a key came up during this frame.
func (win *Window) KeyReleased(k Key) bool {
	return k > 0 && k < keyCount && win.keyUpFrame[k]
}
