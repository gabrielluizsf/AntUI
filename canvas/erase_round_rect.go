package canvas

// eraseRoundRect makes a rounded rectangle fully transparent, the hole an
// inset shadow needs.
func eraseRoundRect(layer *Canvas, x, y, w, h, rx, ry int) {
	if layer.Pixels == nil || w <= 0 || h <= 0 {
		return
	}
	x0, y0, x1, y1 := layer.filterBounds(x, y, w, h)
	for py := y0; py < y1; py++ {
		for px := x0; px < x1; px++ {
			if inRoundRect(px, py, x, y, w, h, rx, ry) {
				layer.Pixels[py*layer.Stride+px] = 0
			}
		}
	}
}
