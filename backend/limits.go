package backend

// Limits are what a window may be resized to. The zero value is a window
// with no bounds at all, which is what a window is without them — so
// Limits{MinWidth: 640, MinHeight: 480} says only what it looks like it says
// and leaves everything else alone.
type Limits struct {
	// The smallest and largest drawable size, in pixels. Zero means no bound
	// on that side. The common case is a window resizable but not below the
	// point where the interface stops fitting, which is a minimum and nothing
	// else.
	//
	// A maximum below a minimum is taken at its word and the maximum wins: a
	// maximum is what a display or a lock imposes, and a minimum is only what
	// the interface would like.
	MinWidth, MinHeight int
	MaxWidth, MaxHeight int

	// Aspect is the width-to-height ratio the window keeps while it is
	// dragged, or zero to keep none. 16.0/9.0 for a widescreen composition.
	//
	// It cannot always be kept — every full screen on a display of another
	// shape is a case where it cannot — and the answer there is bars at the
	// sides rather than a stretched picture. Stretching is the failure people
	// notice and cannot name.
	Aspect float64

	// Fixed forbids resizing outright: the window is both the smallest and
	// the largest it may be. Pixel art at a fixed multiple, a layout composed
	// for one shape, a game meant to look identical on every machine — all of
	// them want this, and a maximise button that silently breaks the
	// composition is a bug the player finds first.
	Fixed bool

	// NoMaximize takes the maximise button away from a window that can still
	// be dragged to a new size.
	//
	// On X11 there is no way to ask for that on its own. A client cannot set
	// _NET_WM_ALLOWED_ACTIONS — the window manager owns it — and the only
	// lever ICCCM gives is an equal minimum and maximum, which is Fixed. So
	// this is honoured on Win32 and Cocoa and quietly does nothing on X11.
	NoMaximize bool
}
