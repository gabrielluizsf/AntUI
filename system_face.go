package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// SystemFace is the platform's own interface font at a size in pixels.
//
//	face, err := antui.SystemFace(15)
//	if err == nil {
//	    canvas.SetDefaultFace(face)
//	}
//
// A program that wants to look like the machine it is on asks for this and
// carries on with the built-in when it is not there, which is what the two
// lines above do. Refusing to start because a font is missing would be a
// strange thing for a window library to do.
func SystemFace(pixels float64) (*canvas.Face, error) {
	_, data, err := SystemFont()
	if err != nil {
		return nil, err
	}
	return canvas.ParseFace(data, pixels)
}
