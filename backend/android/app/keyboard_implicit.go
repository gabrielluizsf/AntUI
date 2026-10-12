//go:build android

package app

// Flags for the on-screen keyboard, as the platform numbers them.
const (
	// KeyboardImplicit asks for the keyboard the way tapping a text field
	// does: the system may decline, on a device with a real keyboard
	// attached.
	KeyboardImplicit = 0x0001
	// KeyboardForced insists. The user cannot dismiss it with the back
	// button, so it is the wrong choice for almost everything.
	KeyboardForced = 0x0002
	// KeyboardImplicitOnly hides the keyboard only if it was shown
	// implicitly, and KeyboardNotAlways leaves it up if the user asked for
	// it by hand.
	KeyboardImplicitOnly = 0x0001
	KeyboardNotAlways    = 0x0002
)
