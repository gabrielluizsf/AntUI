package canvas

import "fmt"

// NewCanvas makes an ARGB32 canvas with its pixels zeroed.
func NewCanvas(width, height int) (*Canvas, error) {
	if !sizeFits(width, height) {
		return nil, fmt.Errorf("antui: canvas size %dx%d is not valid", width, height)
	}
	cv := &Canvas{
		Pixels: make([]Color, width*height),
		Width:  width,
		Height: height,
		Stride: width,
		Clip:   Area{0, 0, width, height},
		format: ARGB32,
	}
	return cv, nil
}
