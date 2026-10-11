package canvas

// Scaled is a copy of a canvas at a different size.
//
// Every pixel of the result is the average of the block of the original that
// lands on it, weighted by alpha so that the colour of a transparent region
// does not bleed into the edge of what is beside it.
func Scaled(src *Canvas, width, height int) *Canvas {
	if src == nil || width <= 0 || height <= 0 || src.Width <= 0 || src.Height <= 0 {
		return nil
	}
	out, err := NewCanvas(width, height)
	if err != nil {
		return nil
	}
	scaleInto(out, src)
	return out
}
