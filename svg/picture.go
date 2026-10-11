package svg

import "github.com/gabrielluizsf/antui/canvas"

// Picture is the bitmap an `<image>` draws and the box it is drawn into. It is
// the whole of what was asked for: where the picture came from has already been
// read by the time one of these exists, so painting needs nothing but the
// pixels and the four numbers that say where they land.
type Picture struct {
	// Pixels is the picture itself, straight-alpha the way a canvas holds it,
	// at the size it was encoded at. Two `<image>`s pointing at the same
	// address share these pixels and are never written to: fitting is done with
	// where they are put rather than by changing them.
	Pixels *canvas.Canvas
	// Box is the rectangle of the drawing the picture is fitted into, as x, y,
	// width and height in the drawing's own coordinates, with a percentage
	// already counted as its fraction of the drawing. It is as the file wrote
	// it, before any transform on the element moves it — see [paintPicture],
	// which puts the transform on afterwards the same way writing takes it.
	Box [4]float64
	// Fit is whether the picture keeps its own shape inside the box (meet), is
	// cut down to the box (slice), or is stretched to fill it (none), which is
	// what `preserveAspectRatio` says.
	Fit PictureFit
	// Align is which corner or edge of the box the picture is pinned to when
	// the fit leaves anything over: one of the nine alignments the attribute
	// names, `xMidYMid` being the one a drawing that says nothing gets.
	Align PictureAlign
}
