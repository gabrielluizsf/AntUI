package svg

import "github.com/gabrielluizsf/antui/canvas"

// tilePicture is one kept picture of a pattern: what it was drawn from, the
// chain of patterns being drawn when it was — kept because the same tile
// reached through two different chains may come out differently, a picture
// naming a pattern that one chain has already open is cut off in that one —
// and the picture itself, which is only ever read after it was drawn.
type tilePicture struct {
	key   tileKey
	chain []*pattern
	pic   *canvas.Canvas
}
