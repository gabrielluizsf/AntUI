package canvas

// maskKeep is how much of one pixel of cv the pixel of the mask over it keeps.
//
// The brightness is taken on the sRGB values the canvas holds rather than on
// the linear ones the spec's own colour space would ask for, which is the same
// choice grayscale() makes and the same numbers it uses: a mask and a graying
// filter that disagreed about what a colour's brightness was would disagree
// about how much of a picture a grey mask keeps.
func maskKeep(c Color, mode MaskMode) float64 {
	alpha := float64(c.A()) / 255
	if mode == MaskAlpha {
		return alpha
	}
	bright := (0.2126*float64(c.R()) + 0.7152*float64(c.G()) + 0.0722*float64(c.B())) / 255
	return bright * alpha
}
