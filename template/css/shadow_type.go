package css

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Shadow is one box-shadow or text-shadow value: an offset, a blur radius, an
// optional spread and colour, and whether it falls inside the box. Lengths are
// reference pixels; the template scales them with the window. A colour that
// was never written stays [CurrentColor] until the cascade resolves it.
type Shadow struct {
	X, Y   int
	Blur   int
	Spread int
	Color  canvas.Color
	Inset  bool
}
