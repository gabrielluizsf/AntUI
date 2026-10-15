package svg

import "github.com/gabrielluizsf/antui/canvas"

// applyFilters runs an element's filter list over the picture it drew, one
// step after another in the order they were written: a grayscale and then a
// blur is not the same picture as a blur and then a grayscale, and the list
// says which one it wants. A step that stands for a `<filter>` element runs
// the primitives inside it over the picture so far, and the steps beside one
// in the same list keep the order they were written in.
//
// Every step goes over the whole picture the layer holds, which is exactly
// what the element asked to have put through it. A smaller rectangle would
// have to be what the element painted, widened before a blur by as far as three
// box passes can carry a picture past where the element stopped — and the
// canvas only remembers where it was written when it is watching for changes,
// which this one is not, so the painted rectangle is not something to be had
// here. The whole of the layer is honest about that, gives a blur room to
// reach in every direction, and costs a walk over transparent pixels that come
// out of a colour filter untouched and of a blur still carrying nothing.
//
// A blur is a size and is written in the drawing's own units, so it is turned
// into pixels here the same as a stroke's width is: one unit of the drawing is
// scale pixels of the canvas, and the scale has the element's transform in it
// because what a `<g transform="scale(2)">` draws is twice as big and its blur
// has to be too. A transform that stretches one axis more than the other makes
// the two scales different, and a blur reads one per axis so it is stretched
// with the shape instead of taking an average of the two. The rest of the
// functions are fractions and angles, which a scale does not touch, and the
// lengths a primitive writes for itself are turned into pixels where the
// primitive is applied.
//
// The picture may come back as a different one than it went in as, because
// every primitive leaves a picture of its own behind, so the caller takes
// what comes out rather than the layer it put in.
func applyFilters(cv *canvas.Canvas, filters []filterOp, s scale2, current canvas.Color) *canvas.Canvas {
	for _, f := range filters {
		if f.def != nil {
			cv = applyFilterDef(cv, f.def, s, current)
			continue
		}
		if f.shadow != nil {
			// A shadow is a picture behind the picture, not a colour put over
			// it, so it is the one step that lays a second picture down and
			// hands the pair on to whatever follows.
			cv = applyDropShadow(cv, f.shadow, s, current)
			continue
		}
		if f.kind == canvas.FilterBlur {
			// A blur is a length, and a length is longer along one axis than
			// the other when the transform stretches it unevenly, so the two
			// radii come from the two axes rather than from their average.
			cv.BlurXY(canvas.Area{X: 0, Y: 0, Width: cv.Width, Height: cv.Height}, radius(f.amount*s.x), radius(f.amount*s.y))
			continue
		}
		cv.FilterRegion(canvas.Area{X: 0, Y: 0, Width: cv.Width, Height: cv.Height}, f.kind, f.amount)
	}
	return cv
}
