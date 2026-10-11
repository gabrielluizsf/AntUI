package antui

// setTouchFirst is how a backend says there is no other kind of input on
// this platform, so that the first frame already knows rather than waiting
// for a finger.
func (win *Window) SetTouchFirst() { win.touchFirst = true }
