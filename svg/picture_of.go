package svg

import "github.com/gabrielluizsf/antui/canvas"

// pictureOf is a picture of the one it is given: the primitives draw into a
// copy of their input so that what they took may be named again by what
// follows and still be as it was.
func pictureOf(src *canvas.Canvas) *canvas.Canvas {
	if src == nil {
		return nil
	}
	out := blankPicture(src.Width, src.Height)
	if out == nil {
		return nil
	}
	if src.Pixels != nil && out.Pixels != nil {
		for y := 0; y < src.Height; y++ {
			copy(out.Pixels[y*out.Stride:y*out.Stride+src.Width], src.Pixels[y*src.Stride:y*src.Stride+src.Width])
		}
		return out
	}
	out.BlitOver(0, 0, src)
	return out
}
