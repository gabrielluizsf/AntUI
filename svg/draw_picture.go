package svg

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// drawPicture is one tile drawn into a canvas of its own, at the size it lands
// at on the page: the content of the pattern put through the coordinates it
// was written in. What runs out of the tile is cut by the canvas being exactly
// the tile — which is what `overflow` on a pattern says unless it says
// `visible`, and then the picture is painted into the tile from every copy of
// it whose overflow reaches in, so a picture written wider than its tile lands
// in the tiles beside it instead of being cut off along the edge. It answers
// nil where there is nothing to draw it with, which is the same as a picture
// that keeps nothing.
func (pt *pattern) drawPicture(box [4]float64, x, y, w, h, s float64, current canvas.Color, patterning []*pattern) *canvas.Canvas {
	px, py := int(math.Ceil(w*s)), int(math.Ceil(h*s))
	if px < 1 {
		px = 1
	}
	if py < 1 {
		py = 1
	}
	layer, err := canvas.NewLayer(px, py)
	if err != nil {
		return nil
	}
	// The picture's own coordinates onto the tile: a point of the picture goes
	// to where it stands in the tile and from there to a pixel of it. What is
	// written in fractions of the box of the shape goes through that box first,
	// the same as a gradient in objectBoundingBox units does. A copy of the
	// picture reaches into the tile from between the two transforms: the offset
	// that brings it in is counted in the tile's own units, so it has to be
	// laid after the box has had its say or the box scales it away.
	base := canvas.Scale(s, s).Mul(canvas.Translate(-x, -y))
	boxT := canvas.Identity()
	if pt.contentBox {
		boxT = canvas.Translate(box[0], box[1]).Mul(canvas.Scale(box[2], box[3]))
	}
	one := func(at canvas.Matrix) {
		paintNodes(layer, pt.content, at, current, nil, append(patterning, pt))
	}
	if pt.visible {
		// `overflow="visible"`: the copies of the picture at the offsets that
		// bring its overflow into the tile, painted left to right in rows from
		// the top — the order the spec says tiles are laid in, which is what
		// says which of two overlapping copies is on top. The picture's span
		// in the tile's own coordinates comes out of what it painted, through
		// the same box the picture itself went through.
		if b, ok := paintedBox(pt.content); ok {
			x0, y0, x1, y1 := b[0], b[1], b[0]+b[2], b[1]+b[3]
			if pt.contentBox {
				x0, x1 = box[0]+x0*box[2], box[0]+x1*box[2]
				y0, y1 = box[1]+y0*box[3], box[1]+y1*box[3]
			}
			if nx0, nx1, ok := copyRange(x0-x, x1-x, w); ok {
				if ny0, ny1, ok := copyRange(y0-y, y1-y, h); ok {
					for ny := ny0; ny <= ny1; ny++ {
						for nx := nx0; nx <= nx1; nx++ {
							one(base.Mul(canvas.Translate(float64(nx)*w, float64(ny)*h)).Mul(boxT))
						}
					}
					return layer
				}
			}
		}
		// The picture reaches nothing into the tile, or reaches so far past it
		// that drawing it over and over would be dearer than the corner of it
		// that would come out: one copy, cut to the tile, which is the default.
	}
	one(base.Mul(boxT))
	return layer
}
