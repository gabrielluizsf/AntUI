package svg

import "github.com/gabrielluizsf/antui/canvas"

// RenderWith is [Image.Render] with the colour that a `currentColor` paint
// takes. It is what lets one drawing be written once and painted in as many
// colours as the interface it belongs to has states: the same tick in the
// normal, hovered and pressed colour of a checkbox, without being written three
// times.
//
// The colour is part of what is remembered, so asking for the same drawing at
// the same size in two colours draws it twice, and asking again in either of
// them costs nothing. A caller that wants an icon to hold a while and be
// repainted often is better off asking for one colour at a time.
func (img *Image) RenderWith(width, height int, current canvas.Color) *canvas.Canvas {
	if img == nil || img.Root == nil || width <= 0 || height <= 0 {
		return nil
	}
	// The one size and colour already drawn, and the canvas it was drawn on. A
	// call for the same one hands the same canvas back; a caller that draws on
	// it is drawing on the icon itself, which is the point.
	img.mu.Lock()
	if img.refreshBuild() {
		img.drawn = nil
	}
	if img.drawn != nil && img.drawnWidth == width && img.drawnHeight == height && img.drawnColour == current {
		cv := img.drawn
		img.mu.Unlock()
		return cv
	}
	img.mu.Unlock()

	cv := img.paint(width, height, current)

	img.mu.Lock()
	img.drawn, img.drawnWidth, img.drawnHeight, img.drawnColour = cv, width, height, current
	img.mu.Unlock()
	return cv
}
