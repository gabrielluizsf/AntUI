package canvas

import "testing"

// A mask is a clip cut out of a picture rather than painted into it: what is
// under the shapes stays, the rest of the clip region is emptied, and the edge
// of each shape is soft. A `<clipPath>` leans on all of that.
//
// A mask may also be a picture that says how much of what it covers is there,
// rather than where it starts and stops: that is MaskBy, where every pixel
// comes down by how much of the picture is over it, and an SVG `<mask>` leans
// on that one.

// filledLayer is a layer painted solid across every pixel, so that a mask has
// something to take away from everywhere at once.
func filledLayer(t *testing.T, w, h int) *Canvas {
	t.Helper()
	cv := newLayer(t, w, h)
	cv.FillRect(0, 0, w, h, RGB(255, 0, 0))
	return cv
}

// TestMaskShapesKeepsOnlyWhatIsUnderAShape is the whole of it: inside the
// shape the picture is as it was, outside the shape there is nothing left of
// it, and what is kept keeps its colour rather than coming out faded.
func TestMaskShapesKeepsOnlyWhatIsUnderAShape(t *testing.T) {
	cv := filledLayer(t, 8, 8)
	p := NewPath()
	p.AddRect(2, 2, 4, 4)
	cv.MaskShapes(MaskShape{Path: p, Rule: FillNonZero})

	if got := cv.At(4, 4); got != RGB(255, 0, 0) {
		t.Errorf("inside the shape is %#08x, want the red it was painted", uint32(got))
	}
	for _, at := range [][2]int{{0, 0}, {1, 5}, {7, 7}, {6, 3}, {4, 6}} {
		if a := cv.At(at[0], at[1]).A(); a != 0 {
			t.Errorf("outside the shape at %v is at alpha %d, want nothing left of it", at, a)
		}
	}
}

// TestMaskShapesUnionsItsShapes: two shapes are the inside of the clip
// together, and the place they overlap is inside twice over rather than being
// cut by the second one that reached it. Cutting them one after the other
// would empty the overlap, because the second cut would find nothing there.
func TestMaskShapesUnionsItsShapes(t *testing.T) {
	cv := filledLayer(t, 8, 8)
	a, b := NewPath(), NewPath()
	a.AddRect(0, 0, 4, 4)
	b.AddRect(3, 3, 4, 4)
	cv.MaskShapes(MaskShape{Path: a, Rule: FillNonZero}, MaskShape{Path: b, Rule: FillNonZero})

	for _, at := range [][2]int{{1, 1}, {6, 6}, {3, 3}, {3, 1}, {1, 3}} {
		if got := cv.At(at[0], at[1]); got.A() == 0 {
			t.Errorf("at %v the picture is gone, want it kept: it is under one of the shapes", at)
		}
	}
	if a := cv.At(7, 0).A(); a != 0 {
		t.Errorf("where neither shape reaches is at alpha %d, want nothing", a)
	}
}

// TestMaskShapesOfNoShapesEmptiesTheClipRegion: the union of nothing is
// nothing, so a clip with no shapes behind it keeps nothing at all. It is what
// an empty `<clipPath>` means, and what a reference to one that is not in the
// drawing ends up doing.
func TestMaskShapesOfNoShapesEmptiesTheClipRegion(t *testing.T) {
	for name, shapes := range map[string][]MaskShape{
		"no shapes at all":          nil,
		"no path":                   {{Rule: FillNonZero}},
		"a path with nothing in it": {{Path: NewPath(), Rule: FillNonZero}},
		"a path outside the canvas": {{Path: pathAt(20, 20, 4, 4), Rule: FillNonZero}},
	} {
		cv := filledLayer(t, 6, 6)
		cv.MaskShapes(shapes...)
		for y := range 6 {
			for x := range 6 {
				if a := cv.At(x, y).A(); a != 0 {
					t.Fatalf("%s: pixel (%d,%d) is at alpha %d, want the whole region emptied", name, x, y, a)
				}
			}
		}
	}
}

// pathAt is a rectangle path, which is the shape every one of these tests
// wants to cut with without spelling out the four corners each time.
func pathAt(x, y, w, h float64) *Path {
	p := NewPath()
	p.AddRect(x, y, w, h)
	return p
}

// TestMaskShapesKeepsTheSoftEdge: a shape that stops in the middle of a pixel
// leaves that pixel half there. The whole edge of a clip is made of pixels
// like that one, and cutting them whole is what turns a smooth curve into a
// staircase.
func TestMaskShapesKeepsTheSoftEdge(t *testing.T) {
	cv := filledLayer(t, 8, 4)
	cv.MaskShapes(MaskShape{Path: pathAt(0.5, 0, 4, 4), Rule: FillNonZero})

	if got := cv.At(0, 0).A(); !near(int(got), 128, 2) {
		t.Errorf("the half-covered pixel is at alpha %d, want about 128", got)
	}
	if got := cv.At(1, 0).A(); got != 255 {
		t.Errorf("the pixel inside the shape is at alpha %d, want it kept whole", got)
	}
	if got := cv.At(4, 0).A(); !near(int(got), 128, 2) {
		t.Errorf("the half-covered pixel at the far edge is at alpha %d, want about 128", got)
	}
	if got := cv.At(5, 0).A(); got != 0 {
		t.Errorf("the pixel past the shape is at alpha %d, want nothing", got)
	}
	// The colour is untouched by the cut: only how much of it there is moves.
	if got := cv.At(0, 0); got.R() != 255 || got.G() != 0 {
		t.Errorf("the half-kept pixel is %#08x, want it still red", uint32(got))
	}
}

// TestMaskShapesHonoursTheRule is the rule a `<clipPath>` writes with
// `clip-rule`, on the one shape it says which of the insides of: two
// rectangles one inside the other are the whole of both under nonzero, and
// under evenodd the second crossing takes the middle back out again.
//
// The rule is a rule for one shape and not for the clip as a whole: two
// rectangles written as two shapes are two shapes, and their union is both of
// them whatever rule either of them is read with.
func TestMaskShapesHonoursTheRule(t *testing.T) {
	nested := func() *Path {
		p := NewPath()
		p.AddRect(0, 0, 8, 8)
		p.AddRect(2, 2, 4, 4)
		return p
	}
	for _, tc := range []struct {
		name   string
		rule   FillRule
		middle bool
	}{
		{name: "nonzero", rule: FillNonZero, middle: true},
		{name: "evenodd", rule: FillEvenOdd, middle: false},
	} {
		cv := filledLayer(t, 8, 8)
		cv.MaskShapes(MaskShape{Path: nested(), Rule: tc.rule})
		if got := cv.At(4, 4).A() != 0; got != tc.middle {
			t.Errorf("%s: the middle of the two rectangles kept=%t, want %t", tc.name, got, tc.middle)
		}
		if cv.At(1, 1).A() == 0 {
			t.Errorf("%s: the ring around the middle is gone, want it kept", tc.name)
		}
	}

	// Two shapes rather than two subpaths of one: the union keeps the middle
	// whatever rule they are read with, because each is measured on its own.
	for _, tc := range []struct {
		name string
		rule FillRule
	}{
		{name: "nonzero", rule: FillNonZero},
		{name: "evenodd", rule: FillEvenOdd},
	} {
		cv := filledLayer(t, 8, 8)
		cv.MaskShapes(
			MaskShape{Path: pathAt(0, 0, 8, 8), Rule: tc.rule},
			MaskShape{Path: pathAt(2, 2, 4, 4), Rule: tc.rule},
		)
		if got := cv.At(4, 4).A(); got == 0 {
			t.Errorf("two shapes with the %s rule emptied the middle, want the union of both", tc.name)
		}
	}
}

// TestMaskShapesOnlyReachesInsideTheClip: the clip region is what the mask is
// allowed to touch, the same as every other draw call. A pixel beside it keeps
// whatever it had, because a mask is not a way to paint over the rest of the
// canvas.
func TestMaskShapesOnlyReachesInsideTheClip(t *testing.T) {
	cv := filledLayer(t, 8, 8)
	cv.SetClip(0, 0, 4, 4)
	// A shape covering everything: inside the clip there is nowhere it does
	// not reach, and outside the clip the mask has no say at all.
	cv.MaskShapes(MaskShape{Path: pathAt(0, 0, 8, 8), Rule: FillNonZero})

	if got := cv.At(1, 1); got != RGB(255, 0, 0) {
		t.Errorf("inside the clip is %#08x, want the red it was painted", uint32(got))
	}
	if got := cv.At(6, 6); got != RGB(255, 0, 0) {
		t.Errorf("outside the clip is %#08x, want it untouched: the mask cannot reach it", uint32(got))
	}
}

// TestMaskShapesCutsOnlyTheAlphaOfAPixelThatWasNeverDrawn is the other half of
// "only the alpha moves": a pixel with nothing on it is not given a colour by
// being cut, and one that keeps everything it had is not written back at all.
func TestMaskShapesCutsOnlyTheAlphaOfAPixelThatWasNeverDrawn(t *testing.T) {
	cv := newLayer(t, 4, 4)
	cv.Pixel(1, 1, RGBA(0, 0, 255, 255))
	cv.MaskShapes(MaskShape{Path: pathAt(0, 0, 4, 4), Rule: FillNonZero})

	if got := cv.At(0, 0); got != 0 {
		t.Errorf("a pixel nothing was drawn on came out %#08x, want it clear", uint32(got))
	}
	if got := cv.At(1, 1); got != RGBA(0, 0, 255, 255) {
		t.Errorf("the pixel that was drawn on came out %#08x, want it as it was", uint32(got))
	}
}

// maskedLayers is a picture cut into three bands — one that keeps everything,
// one that keeps about half and one that keeps nothing — and the same picture
// of it returned in both measures, which is what every MaskBy test below wants
// to hold against the same drawing.
func maskedLayers(t *testing.T, w, h int) (cv, pic *Canvas) {
	t.Helper()
	cv = filledLayer(t, w, h)
	pic = newLayer(t, w, h)
	band := w / 3
	pic.FillRect(0, 0, band, h, RGB(255, 255, 255))
	pic.FillRect(band, 0, band, h, RGB(128, 128, 128))
	pic.FillRect(2*band, 0, w-2*band, h, RGB(0, 0, 0))
	return cv, pic
}

// TestMaskByTakesTheAlphaDownToWhatThePictureKeeps is the whole of a mask
// measured by luminance: bright keeps, dark takes away, and what is between
// the two leaves that much of the picture — a multiply of the alpha and not a
// second decision about whether the pixel is there at all.
func TestMaskByTakesTheAlphaDownToWhatThePictureKeeps(t *testing.T) {
	cv, pic := maskedLayers(t, 9, 3)
	cv.MaskBy(pic, MaskLuminance)

	if got := cv.At(1, 1).A(); got != 255 {
		t.Errorf("under the bright band the picture is at alpha %d, want it kept whole", got)
	}
	if got := cv.At(4, 1).A(); !near(int(got), 128, 2) {
		t.Errorf("under the grey band the picture is at alpha %d, want about half of it", got)
	}
	if got := cv.At(7, 1).A(); got != 0 {
		t.Errorf("under the dark band the picture is at alpha %d, want nothing left of it", got)
	}
}

// TestMaskByKeepsByHowMuchOfThePictureIsThere is the other measure: under
// MaskAlpha brightness is not a word at all, so an opaque black band keeps
// everything the luminance one took away, and what is clear takes away
// whatever it is over.
func TestMaskByKeepsByHowMuchOfThePictureIsThere(t *testing.T) {
	cv := filledLayer(t, 6, 2)
	pic := newLayer(t, 6, 2)
	pic.FillRect(0, 0, 2, 2, RGB(0, 0, 0))
	pic.FillRect(2, 0, 2, 2, RGB(255, 255, 255))
	pic.FillRect(4, 0, 2, 2, RGBA(255, 255, 255, 128))
	cv.MaskBy(pic, MaskAlpha)

	if got := cv.At(0, 0).A(); got != 255 {
		t.Errorf("under the opaque black band the picture is at alpha %d, want it kept whole", got)
	}
	if got := cv.At(2, 0).A(); got != 255 {
		t.Errorf("under the opaque white band the picture is at alpha %d, want it kept whole", got)
	}
	if got := cv.At(4, 0).A(); !near(int(got), 128, 2) {
		t.Errorf("under the half-clear band the picture is at alpha %d, want about half of it", got)
	}
	// The same black band, measured by luminance, is the one that keeps
	// nothing — which is the whole difference between the two measures.
	cv2, pic2 := maskedLayers(t, 9, 1)
	cv2.MaskBy(pic2, MaskLuminance)
	if got := cv2.At(7, 0).A(); got != 0 {
		t.Errorf("measured by luminance the dark band leaves alpha %d, want nothing", got)
	}
}

// TestMaskByLeavesTheColourAlone: only how much of a pixel there is moves,
// the same as every other mask here — a picture taken down to half comes out
// half as strong rather than mixed towards black.
func TestMaskByLeavesTheColourAlone(t *testing.T) {
	cv, pic := maskedLayers(t, 9, 1)
	cv.MaskBy(pic, MaskLuminance)

	if got := cv.At(1, 0); got != RGB(255, 0, 0) {
		t.Errorf("the kept pixel is %#08x, want the red it was painted", uint32(got))
	}
	if got := cv.At(4, 0); got.R() != 255 || got.G() != 0 || got.B() != 0 {
		t.Errorf("the half-kept pixel is %#08x, want it still red", uint32(got))
	}
}

// TestMaskByOnlyReachesInsideTheClip: the clip region is what a mask is
// allowed to touch, as it is for every other call that changes pixels. A pixel
// beside the clip keeps whatever it had, because a mask is not a way to paint
// over the rest of the canvas either.
func TestMaskByOnlyReachesInsideTheClip(t *testing.T) {
	cv, pic := maskedLayers(t, 9, 3)
	cv.SetClip(0, 0, 4, 3)
	cv.MaskBy(pic, MaskLuminance)

	if got := cv.At(7, 1); got != RGB(255, 0, 0) {
		t.Errorf("outside the clip is %#08x, want it untouched: the mask cannot reach it", uint32(got))
	}
	if got := cv.At(1, 1).A(); got != 255 {
		t.Errorf("inside the clip under the bright band is at alpha %d, want it kept", got)
	}
}

// TestMaskByOfAPixelThePictureDoesNotCoverKeepsNothing is the union of no
// shapes again: a picture smaller than what it masks says nothing about the
// pixels beside it, and a pixel no mask has anything to say about keeps
// nothing rather than everything.
func TestMaskByOfAPixelThePictureDoesNotCoverKeepsNothing(t *testing.T) {
	cv := filledLayer(t, 8, 8)
	pic := newLayer(t, 4, 4)
	pic.FillRect(0, 0, 4, 4, RGB(255, 255, 255))
	cv.MaskBy(pic, MaskLuminance)

	if got := cv.At(1, 1).A(); got != 255 {
		t.Errorf("under the picture the alpha is %d, want it kept", got)
	}
	if got := cv.At(5, 5).A(); got != 0 {
		t.Errorf("past the picture the alpha is %d, want nothing: nothing said anything about it", got)
	}
}

// TestMaskByOnNarrowCanvasIsNoOp: the narrow formats keep no alpha per pixel
// to take down, so the call has nowhere to write what it decided and leaves
// the canvas as it was rather than reading colours it cannot store.
func TestMaskByOnNarrowCanvasIsNoOp(t *testing.T) {
	cv, err := NewCanvasFormat(2, 2, RGB565)
	if err != nil {
		t.Fatal(err)
	}
	pic, err := NewCanvasFormat(2, 2, RGB565)
	if err != nil {
		t.Fatal(err)
	}
	cv.MaskBy(pic, MaskLuminance)
	if cv.Pixels != nil {
		t.Fatal("narrow canvas should not have ARGB pixels")
	}
	// And a mask with no picture behind it is not a mask that empties
	// everything: there is nothing to measure by, so nothing is changed.
	layer := newLayer(t, 2, 2)
	layer.FillRect(0, 0, 2, 2, RGB(0, 255, 0))
	layer.MaskBy(nil, MaskLuminance)
	if got := layer.At(0, 0); got != RGB(0, 255, 0) {
		t.Errorf("with no picture to measure by the pixel came out %#08x, want it as it was", uint32(got))
	}
}
