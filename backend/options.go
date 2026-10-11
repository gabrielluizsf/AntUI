package backend

// Options is everything a window can be opened with. The zero value plus a
// size is exactly what Open gives: a resizable window with no bounds.
type Options struct {
	Title         string
	Width, Height int
	Limits        Limits

	// Points asks for Width and Height in points rather than pixels. A window
	// asked for at 1280 by 720 opens 2560 by 1440 pixels on a display at 200
	// per cent and appears the same size on both.
	//
	// Without it "fixed size" means a different thing on every machine: a
	// locked 1280 by 720 window is a postage stamp on a 4K laptop. The limits
	// are scaled with it, since a minimum in points is what was meant.
	Points bool

	// Fullscreen opens straight into a full screen, without the frame of a
	// windowed frame appearing first.
	Fullscreen bool
}
