package antui

// KeyDown reports whether a key is held.
func (win *Window) KeyDown(k Key) bool {
	return k > 0 && k < keyCount && win.keyState[k]
}
