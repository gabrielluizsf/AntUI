package canvas

// Blur blurs the rectangle in place. Three box passes stand in for a Gaussian,
// which is what CSS blur radii are measured in; the colour is premultiplied
// while it moves so a transparent pixel cannot drag its black into the edges.
// It works on ARGB32 only.
func (cv *Canvas) Blur(area Area, radius int) {
	cv.BlurXY(area, radius, radius)
}
