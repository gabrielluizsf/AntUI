package antui

// KeyPressed reports whether a key went down during this frame.
func (win *Window) KeyPressed(k Key) bool {
	return k > 0 && k < keyCount && win.keyDownFrame[k]
}
