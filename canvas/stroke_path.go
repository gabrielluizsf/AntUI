package canvas

// StrokePath draws the outline of a path: a band of the given width centred on
// the line, looking the same width whichever way the line runs because every
// piece of the band is filled the same way. A width of zero or less draws
// nothing, which is what SVG says and the only sensible reading of it.
func (cv *Canvas) StrokePath(path *Path, c Color, style StrokeStyle) {
	cv.FillPath(StrokeOutline(path, style), c, FillNonZero)
}
