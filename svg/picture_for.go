package svg

import "github.com/gabrielluizsf/antui/canvas"

// pictureFor is the picture an href points at, read once and kept under that
// address for every other `<image>` in the same drawing. A picture that cannot
// be read says so and answers nothing, which is what leaves the hole in the
// picture with a warning beside it.
func (img *Image) pictureFor(href string, warn func(string, ...any)) *canvas.Canvas {
	if pic := img.pictures[href]; pic != nil {
		return pic
	}
	data, err := fetchPicture(href)
	if err == nil {
		var pic *canvas.Canvas
		pic, err = decodePicture(data)
		if err == nil {
			if img.pictures == nil {
				img.pictures = map[string]*canvas.Canvas{}
			}
			img.pictures[href] = pic
			return pic
		}
	}
	warn("the picture %q could not be read: %v", href, err)
	return nil
}
