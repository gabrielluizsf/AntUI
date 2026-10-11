package font

import "image"

// ColorGlyph is one glyph drawn as a colour picture. Every number on it is in
// the pixels of the strike the picture was cut at: the top-left corner sits
// BearingX to the right of the pen and BearingY above the baseline, Advance
// is how far the pen moves on, and PPEM is how many pixels per em that strike
// was — a face at another size scales all of them by its own size over PPEM.
type ColorGlyph struct {
	Image              image.Image
	PPEM               int
	BearingX, BearingY int
	Advance            int
}
