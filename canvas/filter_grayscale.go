package canvas

// The filter functions. Blur takes a radius in pixels; hue-rotate takes
// degrees; the rest take a fraction, where 1 is the identity for contrast and
// brightness and the full effect for grayscale, sepia and invert.
const (
	FilterGrayscale FilterKind = iota
	FilterSepia
	FilterInvert
	FilterBrightness
	FilterContrast
	FilterHueRotate
	FilterBlur
)
