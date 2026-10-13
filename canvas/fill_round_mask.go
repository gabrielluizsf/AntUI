package canvas

// fillRoundMask writes an opaque rounded rectangle into a transparent layer,
// keeping the shape's antialiasing in the alpha channel instead of blending it
// into the colour, which is what makes the shadow's own alpha honest.
func fillRoundMask(layer *Canvas, x, y, w, h, rx, ry int, c Color) {
	if w <= 0 || h <= 0 {
		return
	}
	rgb := c & 0x00FFFFFF
	rx = min(rx, w/2)
	ry = min(ry, h/2)
	if rx <= 0 || ry <= 0 {
		layer.FillRect(Area{x, y, w, h}, rgb|0xFF000000)
		return
	}
	layer.FillRect(Area{x + rx, y, w - 2*rx, h}, rgb|0xFF000000)
	layer.FillRect(Area{x, y + ry, rx, h - 2*ry}, rgb|0xFF000000)
	layer.FillRect(Area{x + w - rx, y + ry, rx, h - 2*ry}, rgb|0xFF000000)
	for iy := y; iy < y+h; iy++ {
		if iy >= y+ry && iy < y+h-ry {
			continue
		}
		var dy int
		if iy < y+ry {
			dy = iy - (y + ry - 1)
		} else {
			dy = iy - (y + h - ry)
		}
		for ix := x; ix < x+w; ix++ {
			if ix >= x+rx && ix < x+w-rx {
				continue
			}
			var dx int
			if ix < x+rx {
				dx = ix - (x + rx - 1)
			} else {
				dx = ix - (x + w - rx)
			}
			if alpha := ellipseCoverage(dx, dy, rx-1, ry-1); alpha != 0 {
				layer.Put(ix, iy, rgb|Color(uint32(alpha)<<24))
			}
		}
	}
}
