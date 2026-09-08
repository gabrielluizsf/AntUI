package antui

import "github.com/gabrielluizsf/antui/canvas"

// SetIcon gives the window an icon, at whatever sizes are handed to it.
//
// It reports whether the window system took it. Pass a square picture of any
// size and the window system scales it, or pass several sizes and let it
// pick — [canvas.IconSet] makes a whole set from one picture.
func (win *Window) SetIcon(images ...*canvas.Canvas) bool {
	if win == nil || win.native == nil {
		return false
	}
	var kept []*canvas.Canvas
	for _, image := range images {
		if image != nil && image.Width > 0 && image.Height > 0 {
			kept = append(kept, image)
		}
	}
	if len(kept) == 0 {
		return false
	}
	return win.native.setIcon(kept)
}