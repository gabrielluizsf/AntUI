package svg

import "github.com/gabrielluizsf/antui/canvas"

// pictureRequest is one `<image>` asking for the picture an address points at:
// the slot the node was built with, the opacity the element draws at, and how a
// failure is said. They travel together because all three belong to the element
// rather than to the address, and a fetch still running off the frame has to
// carry them until the picture arrives.
type pictureRequest struct {
	slot    *Picture
	opacity float64
	href    string
	warn    func(string, ...any)
}

// apply puts a fetched picture where the element asked for it: into its slot,
// faded by the opacity the element draws at, and with the canvas the last paint
// made of the drawing dropped. The picture entering the next frame is exactly
// this, and it is what makes a drawing painted while the fetch ran show the
// photo it was waiting for.
func (r pictureRequest) apply(img *Image, pic *canvas.Canvas) {
	if r.opacity < 1 {
		pic = fadePicture(pic, r.opacity)
	}
	r.slot.arrive(pic)
	img.pictureArrived()
}

// unmet says the picture could not be read, which is all that is left to do
// when a fetch running off the frame comes back empty-handed: the hole stays,
// and the warning names the address that did not answer.
func (r pictureRequest) unmet(err error) {
	r.warn("the picture %q could not be read: %v", r.href, err)
}

// watcher is how this element hears back from a fetch running off the frame.
func (r pictureRequest) watcher(img *Image) func(*canvas.Canvas, error) {
	return func(pic *canvas.Canvas, err error) {
		if err != nil {
			r.unmet(err)
			return
		}
		r.apply(img, pic)
	}
}
