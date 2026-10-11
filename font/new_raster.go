package font

func newRaster(w, h int) *raster {
	if w <= 0 || h <= 0 {
		return &raster{}
	}
	return &raster{w: w, h: h, a: make([]float64, (w+1)*h)}
}
