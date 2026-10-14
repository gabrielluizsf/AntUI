package canvas

// FilterRegion applies one filter to the rectangle in place. It works on an
// ARGB32 canvas; the narrow formats have no per-pixel colour to filter and are
// left untouched. A blur is handled by [Canvas.Blur] and is the one filter
// that reads outside the rectangle, so callers give it room to bleed.
func (cv *Canvas) FilterRegion(area Area, kind FilterKind, amount float64) {
	x, y, w, h := area.X, area.Y, area.Width, area.Height
	if cv.Pixels == nil || w <= 0 || h <= 0 {
		return
	}
	if kind == FilterBlur {
		cv.Blur(area, int(amount+0.5))
		return
	}
	x0, y0, x1, y1 := cv.filterBounds(x, y, w, h)
	region := Area{X: x0, Y: y0, Width: x1 - x0, Height: y1 - y0}
	cv.markDirtyRect(region)
	for py := y0; py < y1; py++ {
		row := cv.Pixels[py*cv.Stride : py*cv.Stride+cv.Width]
		for px := x0; px < x1; px++ {
			row[px] = applyFilter(row[px], kind, amount)
		}
	}
	cv.markRegion(region)
}
