package svg

import (
	"sync"
)

// Picture is the picture an `<image>` draws and the box it is drawn into. It is
// the whole of what was asked for: where the picture came from has already been
// read by the time one of these exists, so painting needs nothing but the
// frames and the four numbers that say where they land.
type Picture struct {
	// frames is the picture itself — every frame of an animated one and the
	// single frame of a picture that does not move — at the size it was
	// encoded at. Two `<image>`s pointing at the same address share these
	// frames and are never written to: fitting is done with where they are put
	// rather than by changing them, and which frame of a moving picture is
	// shown is read from the clock at paint time instead of written here.
	//
	// For a picture read from a server it is nothing until the bytes arrive,
	// which happens while the drawing is being painted: the picture of such an
	// element is put in later, by the goroutine that fetched it, and read here
	// through [Picture.at] so that writing it and reading it do not happen
	// at the same time.
	frames *pictureFrames

	// mu guards frames: an `<image>` fetched over the network has its picture
	// filled in by the goroutine that fetched it, while the drawing it belongs
	// to may be painting at the same time.
	mu sync.Mutex

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
