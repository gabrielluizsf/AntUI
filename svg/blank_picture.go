package svg

import "github.com/gabrielluizsf/antui/canvas"

// blankPicture is a cleared picture the whole layer big, which is what a
// primitive that paints rather than changes starts from.
func blankPicture(w, h int) *canvas.Canvas {
	out, err := canvas.NewLayer(w, h)
	if err != nil {
		return nil
	}
	return out
}
