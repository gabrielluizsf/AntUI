package svg

import "github.com/gabrielluizsf/antui/canvas"

// at is the colour the stop is painted in, which for a stop written as
// `currentColor` is the colour of the painting rather than the one the file named.
func (s gradientStop) at(current canvas.Color) canvas.Color {
	if s.current {
		return current
	}
	return s.color
}
