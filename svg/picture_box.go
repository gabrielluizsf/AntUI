package svg

// pictureBox is the box an `<image>` draws its picture into, as the drawing
// wrote it, and which of its sides are the picture's own to say: a side left
// open takes the size the picture arrives with, or its proportion where the
// other side was written. x and y are never open — a picture is put where it
// was asked to be put, which is the origin when nothing asked.
func (img *Image) pictureBox(e *element, warn func(string, ...any)) (box [4]float64, autoW, autoH bool) {
	box[0] = img.pictureLength(warn, e, "x", 0)
	box[1] = img.pictureLength(warn, e, "y", 1)
	box[2], autoW = img.pictureSide(warn, e, "width", 0)
	box[3], autoH = img.pictureSide(warn, e, "height", 1)
	return box, autoW, autoH
}
