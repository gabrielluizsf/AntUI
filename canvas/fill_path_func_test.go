package canvas

import "testing"

// A flat fill is one colour asked of a shape once. A gradient is a colour at
// every point, and these tests are about the seam between the two: the coverage
// comes from the shape as it always has, and the colour is the shade's answer
// for that one pixel.

// TestFillPathFuncPaintsAColourDecidedPerPixel is the whole difference from a
// flat fill — the shade is asked per pixel and what it says is what lands there.
func TestFillPathFuncPaintsAColourDecidedPerPixel(t *testing.T) {
	cv := newTestCanvas(t, 8, 4)
	p := NewPath()
	p.AddRect(0, 0, 8, 4)
	cv.FillPathFunc(p, FillNonZero, func(x, y int) Color {
		return RGB(uint8(x*255/7), uint8(y*255/3), 0)
	})

	for y := range 4 {
		for x := range 8 {
			want := RGB(uint8(x*255/7), uint8(y*255/3), 0)
			if got := cv.At(x, y); got != want {
				t.Errorf("pixel (%d,%d) = %#08x, want %#08x", x, y, uint32(got), uint32(want))
			}
		}
	}
}

// TestFillPathFuncFoldsCoverageIntoTheShade is the antialiased edge of a fill,
// which a gradient gets for nothing: the colour of a part-covered pixel is only
// as strong as the shape covers it.
func TestFillPathFuncFoldsCoverageIntoTheShade(t *testing.T) {
	cv := newTestCanvas(t, 4, 4)
	p := NewPath()
	p.AddRect(0.5, 0, 4, 4)
	cv.FillPathFunc(p, FillNonZero, func(x, y int) Color { return RGB(255, 0, 0) })

	for x, want := range map[int]int{0: 128, 1: 255, 2: 255, 3: 255} {
		if got := coverageAt(cv, x, 0); !near(got, want, 2) {
			t.Errorf("coverage at x=%d = %d, want %d", x, got, want)
		}
	}
}

// TestFillPathFuncMatchesAFlatFillOfTheSameColour is the check that nothing about
// the coverage changed on the way through the shade: asked for one colour, this
// fills exactly what [Canvas.FillPath] fills, pixel for pixel.
func TestFillPathFuncMatchesAFlatFillOfTheSameColour(t *testing.T) {
	red := RGB(255, 0, 0)
	flat := newTestCanvas(t, 16, 16)
	shaded := newTestCanvas(t, 16, 16)

	p := NewPath()
	p.MoveTo(2, 14)
	p.LineTo(8, 1)
	p.LineTo(13, 9)
	p.LineTo(4, 12)
	p.Close()

	flat.FillPath(p, red, FillNonZero)
	shaded.FillPathFunc(p, FillNonZero, func(x, y int) Color { return red })

	for y := range 16 {
		for x := range 16 {
			if got, want := shaded.At(x, y), flat.At(x, y); got != want {
				t.Fatalf("pixel (%d,%d) = %#08x, want %#08x", x, y, uint32(got), uint32(want))
			}
		}
	}
}

// TestFillPathFuncAsksOnlyForThePixelsTheShapeCovers is what keeps a gradient off
// the rest of the box around it: a shade may be expensive — it may be running a
// whole other coordinate system backwards — and it is not asked about a pixel the
// shape has nothing to say about.
func TestFillPathFuncAsksOnlyForThePixelsTheShapeCovers(t *testing.T) {
	cv := newTestCanvas(t, 8, 8)
	p := NewPath()
	p.AddRect(2, 3, 3, 2)

	asked := map[[2]int]bool{}
	cv.FillPathFunc(p, FillNonZero, func(x, y int) Color {
		asked[[2]int{x, y}] = true
		return RGB(255, 0, 0)
	})

	if len(asked) != 6 {
		t.Errorf("the shade was asked about %d pixels, want the 6 the rectangle covers", len(asked))
	}
	for x := 2; x < 5; x++ {
		for y := 3; y < 5; y++ {
			if !asked[[2]int{x, y}] {
				t.Errorf("pixel (%d,%d) was not asked about, want it to be", x, y)
			}
		}
	}
}

// TestFillPathFuncPaintsNothingThereIsNothingToAsk is the shape of the two ways
// this can have nothing to do: no path to fill, and no shade to ask.
func TestFillPathFuncPaintsNothingThereIsNothingToAsk(t *testing.T) {
	cv := newTestCanvas(t, 4, 4)
	p := NewPath()
	p.AddRect(0, 0, 4, 4)

	cv.FillPathFunc(nil, FillNonZero, func(x, y int) Color { return RGB(255, 0, 0) })
	cv.FillPathFunc(p, FillNonZero, nil)
	cv.FillPathFunc(NewPath(), FillNonZero, func(x, y int) Color { return RGB(255, 0, 0) })

	for y := range 4 {
		for x := range 4 {
			if a := cv.At(x, y).A(); a != 0 {
				t.Fatalf("pixel (%d,%d) came out at alpha %d, want the canvas clear", x, y, a)
			}
		}
	}
}
