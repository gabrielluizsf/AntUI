package svg

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
