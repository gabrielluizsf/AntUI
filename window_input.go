package antui

// windowInput is the keyboard and pointer state for the current frame.
type windowInput struct {
	mouseX, mouseY   int
	mouseDX, mouseDY int
	mouseState       [mouseCount]bool
	mouseDownFrame   [mouseCount]bool
	mouseUpFrame     [mouseCount]bool
	wheel            int
	mods             Mod

	keyState     [keyCount]bool
	keyDownFrame [keyCount]bool
	keyUpFrame   [keyCount]bool
}
