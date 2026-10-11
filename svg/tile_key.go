package svg

import "github.com/gabrielluizsf/antui/canvas"

// tileKey is everything a pattern's picture is drawn from apart from the
// chain: the tile where it stands and how big it lands, the colour carried
// into it, and the box of the shape — but only where the picture is written
// in fractions of that box, since a picture in the drawing's own coordinates
// never looks at it and shapes of different sizes share one picture.
type tileKey struct {
	box           [4]float64
	x, y, w, h, s float64
	current       canvas.Color
}
