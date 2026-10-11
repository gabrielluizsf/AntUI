package svg

import "github.com/gabrielluizsf/antui/canvas"

// Render paints the drawing into a canvas of the given size and answers it. The
// size is in pixels and is independent of the size the drawing says it is: an
// icon drawn at twenty-four and asked for at forty-eight comes out at
// forty-eight, with everything about it twice as big.
//
// A drawing that paints itself in `currentColor` is painted in black, which is
// what the keyword means where nothing else has said what the colour is. Use
// [Image.RenderWith] to name one.
//
// The result is remembered, so painting the same drawing at the same size again
// answers the canvas that was already made rather than drawing it a second time.
// That is what makes an icon that is on screen every frame cost nothing after
// the first one. Painting at a new size draws it again and replaces what was
// remembered, so the cache holds one size and no more.
func (img *Image) Render(width, height int) *canvas.Canvas {
	return img.RenderWith(width, height, canvas.RGB(0, 0, 0))
}
