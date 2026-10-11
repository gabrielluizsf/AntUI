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
