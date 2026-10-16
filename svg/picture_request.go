package svg

// pictureRequest is one `<image>` asking for the picture an address points at:
// the slot the node was built with, the opacity the element draws at, and how
// a failure is said. They travel together because all three belong to the
// element rather than to the address, and a fetch still running off the frame
// has to carry them until the picture arrives.
type pictureRequest struct {
	slot    *Picture
	opacity float64
	href    string
	warn    func(string, ...any)
	// chain is the drawings being read through one another to get to this
	// element, which an address that names one of them is stopped by. See
	// [svgFrames].
	chain []string
}

// apply puts a fetched picture where the element asked for it: into its slot,
// faded frame by frame by the opacity the element draws at, and with the canvas
// the last paint made of the drawing dropped. The picture entering the next
// frame is exactly this, and it is what makes a drawing painted while the fetch
// ran show the photo it was waiting for.
func (r pictureRequest) apply(img *Image, pics *pictureFrames) {
	if r.opacity < 1 {
		pics = pics.faded(r.opacity)
	}
	r.slot.arrive(pics)
	img.pictureArrived(pics.moves())
}

// unmet says the picture could not be read, which is all that is left to do
// when a fetch running off the frame comes back empty-handed: the hole stays,
// and the warning names the address that did not answer.
func (r pictureRequest) unmet(err error) {
	r.warn("the picture %q could not be read: %v", r.href, err)
}

// watcher is how this element hears back from a fetch running off the frame.
func (r pictureRequest) watcher(img *Image) func(*pictureFrames, error) {
	return func(pic *pictureFrames, err error) {
		if err != nil {
			r.unmet(err)
			return
		}
		r.apply(img, pic)
	}
}
