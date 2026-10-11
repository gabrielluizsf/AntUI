package canvas

import "github.com/gabrielluizsf/antui/font"

// TextWidth is the width in pixels that a string takes at scale 1.
func TextWidth(text string) int {
	if f := DefaultFace(); f != BuiltinFace() {
		return f.Width(text)
	}
	return font.TextWidth(text)
}
