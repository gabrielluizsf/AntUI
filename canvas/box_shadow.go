package canvas

// BoxShadowParams holds parameters for painting a box shadow.
type BoxShadowParams struct {
	X, Y, Width, Height int
	RX, RY              int
	DX, DY              int
	Blur, Spread        int
	Color               Color
	Inset               bool
}

// BoxShadow paints one blurred rounded-rectangle shadow.
func (cv *Canvas) BoxShadow(p BoxShadowParams) {
	x, y, w, h, rx, ry, dx, dy, blur, spread := p.X, p.Y, p.Width, p.Height, p.RX, p.RY, p.DX, p.DY, p.Blur, p.Spread
	c, inset := p.Color, p.Inset
	if cv.Pixels == nil || c.A() == 0 || w <= 0 || h <= 0 {
		return
	}
	if inset {
		cv.insetShadow(x, y, w, h, rx, ry, dx, dy, blur, spread, c)
		return
	}
	cv.outerShadow(x, y, w, h, rx, ry, dx, dy, blur, spread, c)
}