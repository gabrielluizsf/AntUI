package canvas

// View returns a canvas onto part of this one, sharing its pixels rather
// than copying them. It keeps the parent's stride, so its rows step over the
// rest of the parent — which is what makes a sprite sheet need no special
// case: a frame of it is just a canvas.
//
// Narrow-format canvases have no view: their pixels are not Colors.
func (cv *Canvas) View(x, y, w, h int) *Canvas {
	if cv.Pixels == nil {
		return nil
	}
	x0 := clampInt(x, 0, cv.Width)
	y0 := clampInt(y, 0, cv.Height)
	x1 := clampInt(x+w, x0, cv.Width)
	y1 := clampInt(y+h, y0, cv.Height)

	width, height := x1-x0, y1-y0
	view := &Canvas{
		Width:  width,
		Height: height,
		Stride: cv.Stride,
		Clip:   Area{0, 0, width, height},
		format: ARGB32,
	}
	if width > 0 && height > 0 {
		view.Pixels = cv.Pixels[y0*cv.Stride+x0:]
	}
	return view
}
