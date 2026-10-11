package canvas

// insetShadow draws the shadow inside the box: the layer is filled and the
// inner rounded rectangle is punched out, so what is left blurs inward from
// the edges. The result is clipped to the box.
func (cv *Canvas) insetShadow(x, y, w, h, rx, ry, dx, dy, blur, spread int, c Color) {
	pad := blur + 2
	layer, ok := acquireLayer(w+2*pad, h+2*pad)
	if !ok {
		return
	}
	defer releaseLayer(layer)
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
