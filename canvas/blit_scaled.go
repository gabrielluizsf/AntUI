package canvas

// BlitScaled composites another canvas over this one, fitted into the
// rectangle given rather than at its own size.
func (cv *Canvas) BlitScaled(x, y, width, height int, src *Canvas) {
	if src == nil || width <= 0 || height <= 0 {
		return
	}
	if width == src.Width && height == src.Height {
		cv.Blit(x, y, src)
		return
	}
	if scaled := Scaled(src, width, height); scaled != nil {
		cv.Blit(x, y, scaled)
	}
}
