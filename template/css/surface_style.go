package css

import "github.com/gabrielluizsf/antui/canvas"

// surfaceStyle carries the surface and ink fields of a Style: the colour
// under everything and the layers painted above it.
type surfaceStyle struct {
	Background canvas.Color
	// Background layers, painted above the colour in order. The parallel
	// lists cycle over the layers, so one background-position value
	// positions an entire stack of images.
	BackgroundImages []BackImage
	BackgroundPos    []BackPos
	BackgroundSize   []BackSize
	BackgroundRepeat []BackRepeat
	BackgroundClip   []uint8
	BackgroundOrigin []uint8
	BackgroundAttach []uint8
	Opacity          float64 // 0..1; 1 when unset
	Color            canvas.Color
}
