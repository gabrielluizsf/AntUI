package canvas

// outerShadow draws the shadow as a rounded rectangle behind the box, blurred
// and offset.
func (cv *Canvas) outerShadow(x, y, w, h, rx, ry, dx, dy, blur, spread int, c Color) {
	sx, sy := x+dx-spread, y+dy-spread
	sw, sh := w+2*spread, h+2*spread
	if sw <= 0 || sh <= 0 {
		return
	}
	pad := blur + 2
	layer, ok := acquireLayer(sw+2*pad, sh+2*pad)
	if !ok {
		return
	}
	defer releaseLayer(layer)
	fillRoundMask(layer, pad, pad, sw, sh, max(rx+spread, 0), max(ry+spread, 0), c)
	layer.Blur(Area{X: 0, Y: 0, Width: layer.Width, Height: layer.Height}, blurRadius(blur))
	ScaleAlpha(layer, float64(c.A())/255)
	cv.BlitOver(sx-pad, sy-pad, layer)
}
