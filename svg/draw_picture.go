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
	px, py := max(int(math.Ceil(w*s)), 1), max(int(math.Ceil(h*s)), 1)
	layer, err := canvas.NewLayer(px, py)
	if err != nil {
		return nil
	}
	// The picture's own coordinates onto the tile: a point of the picture goes
	// to where it stands in the tile and from there to a pixel of it. What is
	// written in fractions of the box of the shape goes through that box first,
	// the same as a gradient in objectBoundingBox units does.
	base := canvas.Scale(s, s).Mul(canvas.Translate(-x, -y))
	boxT := canvas.Identity()
	if pt.contentBox {
		boxT = canvas.Translate(box[0], box[1]).Mul(canvas.Scale(box[2], box[3]))
	}
	one := func(at canvas.Matrix) {
		paintNodes(layer, pt.content, at, palette{current: current}, nil, append(patterning, pt))
	}
	if pt.visible && pt.drawTiledPicture(box, x, y, w, h, base, boxT, one) {
		return layer
	}
	one(base.Mul(boxT))
	return layer
}
