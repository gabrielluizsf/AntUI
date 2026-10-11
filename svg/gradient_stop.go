package svg

import "github.com/gabrielluizsf/antui/canvas"

// gradientStop is one colour change: a colour and where along the line it
// happens, from 0 to 1. A stop written as `currentColor` keeps the keyword
// rather than a colour, because the colour of the widget is only known at
// painting time.
type gradientStop struct {
	offset    float64
	hasOffset bool
	color     canvas.Color
	current   bool
}
