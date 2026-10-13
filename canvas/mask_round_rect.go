package canvas

// MaskRoundRect zeros everything outside a rounded rectangle, the clip a
// background layer gets so the corner transparency of the box shows through
// to what lies behind. A radius of zero makes it exactly [Canvas.MaskRect].
func (cv *Canvas) MaskRoundRect(b Rounded) {
	x, y, w, h := b.X, b.Y, b.Width, b.Height
	rx, ry := min(b.RX, w/2), min(b.vertical(), h/2)
	if w <= 0 || h <= 0 {
		return
	}
	if rx <= 0 || ry <= 0 {
		cv.MaskRect(b.Area)
		return
	}
	for py := cv.Clip.Y; py < cv.Clip.Y+cv.Clip.Height; py++ {
		// The straight run between the corner bands has no curve: only the
		// columns outside the rectangle need clearing. The top and bottom
		// bands (and the rows beyond the rectangle) run the ellipse test, with
		// the area outside the rectangle cleared before it.
		if py >= y+ry && py < y+h-ry {
			for px := cv.Clip.X; px < cv.Clip.X+cv.Clip.Width; px++ {
				if px < x || px >= x+w {
					cv.Put(px, py, 0)
				}
			}
			continue
		}
		for px := cv.Clip.X; px < cv.Clip.X+cv.Clip.Width; px++ {
			if px < x || px >= x+w || py < y || py >= y+h {
				cv.Put(px, py, 0)
				continue
			}
			if !inRoundRect(px, py, x, y, w, h, rx, ry) {
				cv.Put(px, py, 0)
			}
		}
	}
}
