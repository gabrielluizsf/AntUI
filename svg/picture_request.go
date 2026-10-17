package svg

// pictureRequest is one `<image>` asking for the picture an address points at:
// the slot the node was built with, how the picture is fitted into it, the
// opacity the element draws at, and how a failure is said. They travel
// together because all of it belongs to the element rather than to the
// address, and a fetch still running off the frame has to carry them until
// the picture arrives.
type pictureRequest struct {
	slot    *Picture
	fit     PictureFit
	align   PictureAlign
	opacity float64
	href    string
	warn    func(string, ...any)
	// deferred is what a leading `defer` on the element's
	// preserveAspectRatio asked for: the fit of the content pointed at wins
	// over everything in this request when that content is a drawing with a
	// preserveAspectRatio of its own. See [readPreserveAspect].
	deferred bool
	// chain is the drawings being read through one another to get to this
	// element, which an address that names one of them is stopped by. See
	// [svgFrames].
	chain []string
}

// apply puts a fetched picture where the element asked for it: into its slot,
// fitted the way the two of them finally agree on — a `defer` handing the fit
// to the drawing behind the address when it has a preserveAspectRatio of its
// own, and the element's own attribute standing when it does not, which is
// every photograph — faded frame by frame by the opacity the element draws at,
// and with the canvas the last paint made of the drawing dropped. The picture
// entering the next frame is exactly this, and it is what makes a drawing
// painted while the fetch ran show the photo it was waiting for.
func (r pictureRequest) apply(img *Image, pics *pictureFrames) {
	if r.opacity < 1 {
		pics = pics.faded(r.opacity)
	}
	fit, align := r.fit, r.align
	if r.deferred && pics.hasDocFit {
		fit, align = pics.docFit, pics.docAlign
	}
	r.slot.arrive(pics, fit, align)
	img.pictureArrived(pics.moves())
}

// unmet says the picture could not be read, which is all that is left to do
// when a fetch running off the frame comes back empty-handed: the hole stays,
// and the warning names the address that did not answer.
func (r pictureRequest) unmet(err error) {
	r.warn("the picture %q could not be read: %v", r.href, err)
}
