package antui

// TouchScreen reports whether this window is being driven by a touchscreen,
// which is what decides whether a game has to put controls on the screen.
//
// It is not a question about the hardware, it is a question about the
// player: a laptop with a touchscreen nobody uses should be given a keyboard
// game, and a tablet with a keyboard attached still has a screen to press.
// So the answer is "a finger has touched this window", which is false until
// one does — and on a platform where there is nothing else, true from the
// start.
//
// A game that wants to decide for itself should ignore this and set its own
// flag; this is the default, not the rule.
func (win *Window) TouchScreen() bool {
	if win == nil {
		return false
	}
	return win.touchFirst || win.sawTouch
}
