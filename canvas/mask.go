package canvas

// MaskRect zeros everything outside a rectangle, keeping only the pixels
// inside it. It is how a background painted into a layer is clipped to part of
// the box — the padding or content box — before the layer is composited.
// It honours the clip on every canvas format.
func (cv *Canvas) MaskRect(x, y, w, h int) {
	for py := cv.Clip.Y; py < cv.Clip.Y+cv.Clip.Height; py++ {
		for px := cv.Clip.X; px < cv.Clip.X+cv.Clip.Width; px++ {
			if px < x || px >= x+w || py < y || py >= y+h {
				cv.Put(px, py, 0)
			}
		}
	}
}

// MaskRoundRect zeros everything outside a rounded rectangle, the clip a
// background layer gets so the corner transparency of the box shows through
// to what lies behind. A radius of zero makes it exactly [Canvas.MaskRect].
func (cv *Canvas) MaskRoundRect(x, y, w, h, rx, ry int) {
	if w <= 0 || h <= 0 {
		return
	}
	rx = min(rx, w/2)
	ry = min(ry, h/2)
	if rx <= 0 || ry <= 0 {
		cv.MaskRect(x, y, w, h)
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
