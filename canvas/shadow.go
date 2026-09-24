package canvas

// BoxShadow paints one blurred rounded-rectangle shadow for a box at x, y, w,
// h with corner radii rx and ry. dx and dy offset the shadow the way a
// box-shadow length pair does, spread grows the shadow before it blurs and
// blur is the CSS blur radius. An outer shadow is cast around the box; an
// inset shadow is cast inside it. It works on ARGB32 only.
func (cv *Canvas) BoxShadow(x, y, w, h, rx, ry, dx, dy, blur, spread int, c Color, inset bool) {
	if cv.Pixels == nil || c.A() == 0 || w <= 0 || h <= 0 {
		return
	}
	if inset {
		cv.insetShadow(x, y, w, h, rx, ry, dx, dy, blur, spread, c)
		return
	}
	cv.outerShadow(x, y, w, h, rx, ry, dx, dy, blur, spread, c)
}

// outerShadow draws the shadow as a rounded rectangle behind the box, blurred
// and offset.
func (cv *Canvas) outerShadow(x, y, w, h, rx, ry, dx, dy, blur, spread int, c Color) {
	sx, sy := x+dx-spread, y+dy-spread
	sw, sh := w+2*spread, h+2*spread
	if sw <= 0 || sh <= 0 {
		return
	}
	pad := blur + 2
	layer, err := NewCanvas(sw+2*pad, sh+2*pad)
	if err != nil {
		return
	}
	fillRoundMask(layer, pad, pad, sw, sh, max(rx+spread, 0), max(ry+spread, 0), c)
	layer.Blur(0, 0, layer.Width, layer.Height, blurRadius(blur))
	ScaleAlpha(layer, float64(c.A())/255)
	cv.BlitOver(sx-pad, sy-pad, layer)
}

// insetShadow draws the shadow inside the box: the layer is filled and the
// inner rounded rectangle is punched out, so what is left blurs inward from
// the edges. The result is clipped to the box.
func (cv *Canvas) insetShadow(x, y, w, h, rx, ry, dx, dy, blur, spread int, c Color) {
	pad := blur + 2
	layer, err := NewCanvas(w+2*pad, h+2*pad)
	if err != nil {
		return
	}
	layer.FillRect(0, 0, layer.Width, layer.Height, c|0xFF000000)
	holeX, holeY := pad+dx+spread, pad+dy+spread
	holeW, holeH := w-2*spread, h-2*spread
	if holeW > 0 && holeH > 0 {
		eraseRoundRect(layer, holeX, holeY, holeW, holeH, max(rx-spread, 0), max(ry-spread, 0))
	}
	layer.Blur(0, 0, layer.Width, layer.Height, blurRadius(blur))
	ScaleAlpha(layer, float64(c.A())/255)
	saved := cv.Clip
	cv.SetClip(x, y, w, h)
	cv.BlitOver(x-pad, y-pad, layer)
	cv.Clip = saved
}

// blurRadius turns a CSS blur radius into the box radius three passes want.
// A Gaussian stands in for the CSS blur with a standard deviation of half the
// radius, but three box passes carry heavier shoulders than a Gaussian and
// bleed noticeably wider than a browser at the same σ. Taking a third of the
// radius brings the visible reach back to roughly a Gaussian of that σ, so a
// shadow hugs its box the way it does on the web.
func blurRadius(blur int) int {
	if blur <= 0 {
		return 0
	}
	return max(blur/3, 1)
}

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
		layer.FillRect(x, y, w, h, rgb|0xFF000000)
		return
	}
	layer.FillRect(x+rx, y, w-2*rx, h, rgb|0xFF000000)
	layer.FillRect(x, y+ry, rx, h-2*ry, rgb|0xFF000000)
	layer.FillRect(x+w-rx, y+ry, rx, h-2*ry, rgb|0xFF000000)
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

// ScaleAlpha multiplies every pixel's alpha, leaving its colour alone. It is
// how a layer painted at full strength is brought down to a colour's own
// transparency before it is composited.
func ScaleAlpha(cv *Canvas, factor float64) {
	if factor >= 1 || factor < 0 {
		return
	}
	for i, c := range cv.Pixels {
		cv.Pixels[i] = Fade(c, int(float64(c.A())*factor+0.5))
	}
}
