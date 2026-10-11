package svg

const (
	// PictureFitMeet keeps the whole picture inside the box, padding what is
	// left over with nothing. It is the default and the usual way a photograph
	// is shown.
	PictureFitMeet PictureFit = iota
	// PictureFitSlice covers the box with the picture, cutting off what falls
	// outside it — the way a full-bleed background is shown.
	PictureFitSlice
	// PictureFitNone stretches the picture to the box with no regard for the
	// shape it was made in.
	PictureFitNone
)
