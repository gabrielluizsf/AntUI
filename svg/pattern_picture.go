package svg

import (
	"slices"

	"github.com/gabrielluizsf/antui/canvas"
)

// picture is one tile at the size it lands at on the page, drawn into a canvas
// of its own — and kept, so that the shapes asking the same of the pattern
// share it rather than every one of them drawing it again. What it is drawn
// from is the whole of the key it is kept under: the tile where it stands and
// how big it lands, the colour carried into it, the box of the shape where the
// picture is written in fractions of that box, and the chain of patterns that
// was already being drawn — the same tile reached through two chains may come
// out differently, a picture naming a pattern one chain has open being cut off
// in that one. It answers nil where there is nothing to draw it with, which is
// the same as a picture that keeps nothing, and a nil is kept the same as any
// other.
func (pt *pattern) picture(box [4]float64, x, y, w, h, s float64, current canvas.Color, patterning []*pattern) *canvas.Canvas {
	key := tileKey{x: x, y: y, w: w, h: h, s: s, current: current}
	if pt.contentBox {
		key.box = box
	}
	pt.mu.Lock()
	for _, t := range pt.tiles {
		if t.key == key && slices.Equal(t.chain, patterning) {
			pic := t.pic
			pt.mu.Unlock()
			return pic
		}
	}
	pt.drawn++
	pt.mu.Unlock()
	pic := pt.drawPicture(box, x, y, w, h, s, current, patterning)
	pt.keepTile(key, patterning, pic)
	return pic
}
