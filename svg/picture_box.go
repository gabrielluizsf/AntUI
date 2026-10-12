package svg

// pictureBox is the box an `<image>` draws its picture into, as the drawing
// wrote it.
func (img *Image) pictureBox(e *element, warn func(string, ...any)) [4]float64 {
	return [4]float64{
		img.pictureLength(warn, e, "x", 0),
		img.pictureLength(warn, e, "y", 1),
		img.pictureLength(warn, e, "width", 0),
		img.pictureLength(warn, e, "height", 1),
	}
}
