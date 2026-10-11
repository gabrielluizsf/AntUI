package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// applyDropShadow lays the element's silhouette behind its own picture: the
// shape the picture has, in the shadow's colour, moved by the offset and
// softened by the blur, and then the picture itself over it. The lengths are in
// the drawing's own units, so they are turned into pixels the way a blur in the
// attribute is, with the element's transform in the scale.
//
// The silhouette is the picture's alpha with the shadow's colour, which is what
// makes a `drop-shadow` hug the shape rather than a box the way a `box-shadow`
// does. The result is a picture of the whole layer, like every other filter
// step, so what comes out is what the next step — and the clip and the mask
// after it — are given.
func applyDropShadow(src *canvas.Canvas, sh *dropShadow, s scale2, current canvas.Color) *canvas.Canvas {
	if src == nil {
		return nil
	}
	out := blankPicture(src.Width, src.Height)
	if out == nil {
		return src
	}
	c := sh.colour
	if sh.current {
		c = current
	}
	if shadow := silhouette(src, c); shadow != nil {
		shadow.BlurXY(0, 0, shadow.Width, shadow.Height, radius(sh.blur*s.x), radius(sh.blur*s.y))
		out.BlitOver(int(math.Round(sh.dx*s.x)), int(math.Round(sh.dy*s.y)), shadow)
	}
	out.BlitOver(0, 0, src)
	return out
}
