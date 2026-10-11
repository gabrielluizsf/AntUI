package svg

// Warnings is everything the drawing asked for that this package could not do,
// read after it is painted. A drawing that came out as it should has none.
func (img *Image) Warnings() Warnings { return img.warnings }
