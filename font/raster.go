package font

// raster accumulates edge contributions and then reads them off as
// coverage.
type raster struct {
	w, h int
	a    []float64 // one more column than w, so the last delta has somewhere
}
