package svg

import "github.com/gabrielluizsf/antui/canvas"

// pictureFor is the picture an href points at, read once and kept under the
// address it resolves to — the base the drawing sits at counted in — so that
// every other `<image>` naming the same address in the same drawing shares it.
// A picture that cannot be read says so and answers nothing, which is what
// leaves the hole in the picture with a warning beside it.
func (img *Image) pictureFor(base, href string, warn func(string, ...any)) *canvas.Canvas {
	target := resolveHref(base, href)
	if pic := img.pictures[target]; pic != nil {
		return pic
	}
	data, err := fetchPicture(target)
	if err == nil {
		var pic *canvas.Canvas
		pic, err = decodePicture(data)
		if err == nil {
			if img.pictures == nil {
				img.pictures = map[string]*canvas.Canvas{}
			}
			img.pictures[target] = pic
			return pic
		}
	}
	warn("the picture %q could not be read: %v", href, err)
	return nil
}
