package canvas

// MaskShapes keeps only what lies under one of the shapes and empties the rest
// of the clip region, which is a clip: the shapes are the inside of it and
// everything else is drawn nowhere. The edge of every shape is soft — a pixel
// it covers by half keeps half of what it had — which is the antialiasing a
// rectangle cannot have and the reason this is not [Canvas.MaskRect] with a
// path in it.
//
// The shapes are a union and not a run of cuts one after the other: two of
// them overlapping leave their overlap standing, which is what a clip written
// as several shapes means. The union of no shapes at all is nothing, so no
// shapes at all empties the clip region — that is the clip an empty
// `<clipPath>` is, where nothing of what it cuts is drawn. A shape that
// carries cuts of its own (see [MaskShape.And]) has to be under them as well:
// what its own outline covers and what they cover is where it keeps.
//
// What is cut is the alpha alone: the colours stay as they are and only how
// much of each there is changes, so this is a cut to make on a layer, where
// emptiness is a colour a pixel can have. The narrow formats keep no alpha —
// what the mask empties there becomes the colour of nothing, and the soft edge
// along the way is lost, the same as it is in [Canvas.MaskRect].
func (cv *Canvas) MaskShapes(shapes ...MaskShape) {
	cx0, cy0 := cv.Clip.X, cv.Clip.Y
	cx1, cy1 := cx0+cv.Clip.Width, cy0+cv.Clip.Height
	if cx0 >= cx1 || cy0 >= cy1 {
		return
	}
	// The band is the box the shapes can reach, cut down to the clip region:
	// what is beside or beyond it is under no shape at all, and only inside it
	// is there anything to measure.
	bx0, by0, bx1, by1, ok := maskBand(shapes, cx0, cy0, cx1, cy1)
	if !ok {
		// Nothing reaches the clip region, so nothing in it is under a shape.
		for y := cy0; y < cy1; y++ {
			cv.cutRow(y, cx0, cx1, nil)
		}
		return
	}
	// Every shape is measured over the whole band at once, along with the
	// clips each one carries and the clips those carry, all the way down.
	all := make([]maskMeasured, 0, len(shapes))
	for _, s := range shapes {
		if s.Path == nil {
			continue
		}
		all = append(all, measureShape(s, bx0, by0, bx1, by1))
	}
	a := maskApply{
		cv: cv, all: all, row: make([]float64, bx1-bx0),
		runs: &maskRuns{box: bx0, width: bx1 - bx0},
		cx0:  cx0, cx1: cx1, bx0: bx0, bx1: bx1, by0: by0, by1: by1,
	}
	for y := cy0; y < cy1; y++ {
		a.paint(y)
	}
}
