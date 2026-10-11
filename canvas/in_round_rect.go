package canvas

// inRoundRect reports whether a point is inside a rounded rectangle.
func inRoundRect(px, py, x, y, w, h, rx, ry int) bool {
	if px < x || py < y || px >= x+w || py >= y+h {
		return false
	}
	rx = min(rx, w/2)
	ry = min(ry, h/2)
	if rx <= 0 || ry <= 0 {
		return true
	}
	cx := clampInt(px, x+rx, x+w-rx)
	cy := clampInt(py, y+ry, y+h-ry)
	dx, dy := px-cx, py-cy
	if dx == 0 && dy == 0 {
		return true
	}
	return float64(dx*dx)/float64(rx*rx)+float64(dy*dy)/float64(ry*ry) <= 1
}
