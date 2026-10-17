package svg

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// pictureFrames is what an address answered with: every frame the picture
// draws — one for a picture that does not move — and how long each one is on
// screen, in seconds. An animated GIF comes as several; a PNG or a JPEG as one,
// which has no clock to it and answers the same pixel whenever it is asked.
type pictureFrames struct {
	frames []*canvas.Canvas
	show   []float64

	// ownW and ownH are the size the picture says it is — the pixels it was
	// encoded at, or the size the document behind it says it is — which is
	// what a side of an `<image>` box left open answers with. See [autoBox].
	ownW, ownH float64

	// docFit and docAlign are the preserveAspectRatio of the drawing these
	// frames were painted from, when it wrote one: which is what a `defer` on
	// the `<image>` pointing at it asks to be fitted by. A photograph has
	// none — no PNG carries the attribute — so hasDocFit says whether there
	// is one at all, and without one a deferred element is fitted by its own
	// attribute, 'defer' ignored, as SVG 1.1 has it.
	docFit    PictureFit
	docAlign  PictureAlign
	hasDocFit bool
}

// ownSize is that size, as the fitting wants it: two numbers where a picture
// always has two, and nothing where there is no picture to ask.
func (f *pictureFrames) ownSize() (w, h float64) {
	if f == nil {
		return 0, 0
	}
	return f.ownW, f.ownH
}

// moves says the picture animates: more than one frame going around a loop,
// which is what makes repainting it as the clock runs worth doing.
func (f *pictureFrames) moves() bool { return f != nil && len(f.frames) > 1 }
