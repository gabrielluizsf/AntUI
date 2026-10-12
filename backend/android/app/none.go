//go:build android

package app

// The events. The ones that carry something say so.
const (
	// None is no event: what [App.Poll] returns when nothing is waiting.
	None Event = iota

	// Start, Resume, Pause and Stop are the activity lifecycle. An app is
	// between Resume and Pause when the user is looking at it.
	Start
	Resume
	Pause
	Stop
	// Destroy is the last event. Next reports false after it.
	Destroy

	// WindowUp says there is a surface to draw on: [App.Window] is not nil
	// from here until WindowDown.
	WindowUp
	// WindowResized and Redraw both mean draw a whole frame now. Redraw in
	// particular is the system saying the old contents are gone.
	WindowResized
	Redraw
	// WindowDown says the surface is going away. **The UI thread is blocked
	// until the next call to Next**, so stop drawing and come back promptly;
	// [App.Window] is already nil.
	WindowDown

	// FocusGained and FocusLost are the window's keyboard focus.
	FocusGained
	FocusLost

	// ContentRect is the area not covered by system bars; [App.ContentRect]
	// has it.
	ContentRect
	// ConfigChanged is a rotation, a locale, a density, a theme.
	ConfigChanged
	// LowMemory is the system asking for caches back.
	LowMemory

	// InputUp and InputDown bracket the input queue's life. Reading it is
	// the input layer's business, not this package's.
	InputUp
	InputDown
)
