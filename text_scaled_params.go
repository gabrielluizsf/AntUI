package antui

import "github.com/gabrielluizsf/antui/canvas"

// TextScaledParams groups the text, colour and scale of a scaled text call.
type TextScaledParams struct {
	Text  string
	Color canvas.Color
	Scale int
}