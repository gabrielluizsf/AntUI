package canvas

import "fmt"

// NewCanvasFormat makes a canvas holding pixels in the given format. ARGB32
// is exactly NewCanvas. A Pal8 canvas starts on a grey ramp, so a program
// that forgets to set a palette gets a grey picture rather than a black one
// and can see what it forgot.
func NewCanvasFormat(width, height int, format Format) (*Canvas, error) {
	if format == ARGB32 {
		return NewCanvas(width, height)
	}
	if format != RGB565 && format != Pal8 {
		return nil, fmt.Errorf("antui: there is no format %d", int(format))
	}
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("antui: canvas size %dx%d is not valid", width, height)
	}
	cv := &Canvas{
		Width:  width,
		Height: height,
		Stride: width,
		Clip:   Area{0, 0, width, height},
		format: format,
		narrow: make([]byte, width*height*format.Bytes()),
		dither: true,
	}
	if format == Pal8 {
		ramp := make([]Color, 256)
		for i := range ramp {
			ramp[i] = RGB(uint8(i), uint8(i), uint8(i))
		}
		cv.palette = ramp
		_ = cv.SetPalette(ramp)
	}
	return cv, nil
}
