package canvas

// BlitScaled composites another canvas over this one, fitted into the
// rectangle given rather than at its own size.
func (cv *Canvas) BlitScaled(area Area, src *Canvas) {
	x, y, width, height := area.X, area.Y, area.Width, area.Height
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
