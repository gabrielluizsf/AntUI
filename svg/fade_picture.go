package svg

import "github.com/gabrielluizsf/antui/canvas"

// fadePicture is a picture at an opacity, in pixels of its own — the same
// multiply of alpha that [fade] does to a colour, done to every pixel of a
// bitmap instead. The picture it is given is left exactly as it is, since what
// is in the cache belongs to every element that points at the same address.
func fadePicture(src *canvas.Canvas, opacity float64) *canvas.Canvas {
	if opacity >= 1 {
		return src
	}
	cv, err := canvas.NewCanvas(src.Width, src.Height)
	if err != nil {
		// There is nowhere to put a faded copy, and the unfaded one is the
		// worse of the two things to paint: it is a canvas that could not be
		// allocated, which does not happen, and saying nothing beats painting
		// a hole where a picture was.
		return src
	}
	for y := range src.Height {
		for x := range src.Width {
			if c := src.At(x, y); c != 0 {
				cv.Put(x, y, fade(c, opacity))
			}
		}
	}
	return cv
}
