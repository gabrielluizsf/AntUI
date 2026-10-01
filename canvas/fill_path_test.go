package canvas

import "testing"

// The canvas blends over a framebuffer, and a framebuffer has nothing behind it
// to show through, so a blended pixel always comes back opaque: the coverage of
// a fill shows up in the colour it lands on, not in its alpha. Every test here
// paints over black and reads the red channel back as the share of the pixel the
// shape covered, which is the same number the alpha would have carried.

func TestFillPathRect(t *testing.T) {
	cv := newTestCanvas(t, 16, 16)
	p := NewPath()
	p.AddRect(4, 4, 8, 8)
	cv.FillPath(p, RGB(255, 0, 0), FillNonZero)

	for y := range 16 {
		for x := range 16 {
			in := x >= 4 && x < 12 && y >= 4 && y < 12
			if got := coverageAt(cv, x, y); got != pick(in, 255, 0) {
				t.Fatalf("coverage at (%d,%d) = %d, want %d", x, y, got, pick(in, 255, 0))
			}
		}
	}
}

func TestFillPathFractionalRectCoverage(t *testing.T) {
	// A half-pixel edge: the pixel straddling it comes out at half strength,
	// which is the whole point of cutting a row into lines.
	cv := newTestCanvas(t, 4, 4)
	p := NewPath()
	p.AddRect(0.5, 0, 2, 4)
	cv.FillPath(p, RGB(255, 0, 0), FillNonZero)

	for y := range 4 {
		for x, want := range map[int]int{0: 128, 1: 255, 2: 128, 3: 0} {
			if got := coverageAt(cv, x, y); !near(got, want, 2) {
				t.Errorf("coverage at (%d,%d) = %d, want %d", x, y, got, want)
			}
		}
	}
}

func TestFillPathTriangle(t *testing.T) {
	cv := newTestCanvas(t, 32, 32)
	p := NewPath()
	p.MoveTo(16, 4)
	p.LineTo(28, 26)
	p.LineTo(4, 26)
	p.Close()
	cv.FillPath(p, RGB(255, 0, 0), FillNonZero)

	if got := coverageAt(cv, 16, 20); got != 255 {
		t.Errorf("inside coverage = %d, want 255", got)
	}
	if got := coverageAt(cv, 16, 2); got != 0 {
		t.Errorf("outside coverage = %d, want 0", got)
	}
}

func TestFillPathCircleCoverage(t *testing.T) {
	cv := newTestCanvas(t, 40, 40)
	p := NewPath()
	p.AddCircle(20, 20, 15)
	cv.FillPath(p, RGB(255, 255, 255), FillNonZero)

	if got := coverageAt(cv, 20, 20); got != 255 {
		t.Errorf("centre coverage = %d, want 255", got)
	}
	if got := coverageAt(cv, 0, 0); got != 0 {
		t.Errorf("corner coverage = %d, want 0", got)
	}
	// The rim of a circle has to be soft, not a hard step: a ring of boundary
	// pixels with nothing between them is what a stair-stepped fill looks like.
	soft, hard := 0, 0
	for y := range 40 {
		for x := range 40 {
			switch a := coverageAt(cv, x, y); {
			case a > 0 && a < 255:
				soft++
			case a == 255:
				hard++
			}
		}
	}
	if soft == 0 {
		t.Error("no antialiased edge: the rim of the circle is hard")
	}
	if hard == 0 {
		t.Error("nothing filled")
	}
}

func TestFillPathNonZeroRings(t *testing.T) {
	// Two rings wound the same way, one inside the other. Nonzero sees both
	// wind the same direction, so the middle is covered twice and stays inside:
	// the shape is a disc. Evenodd counts crossings and so leaves the middle
	// covered an even number of times, which is a hole.
	p := NewPath()
	p.AddEllipse(20, 20, 18, 18)
	p.AddEllipse(20, 20, 8, 8)
	if got := fillBands(t, p, 20, FillNonZero); got != 1 {
		t.Errorf("nonzero over two rings wound alike gives %d bands, want 1", got)
	}
	if got := fillBands(t, p, 20, FillEvenOdd); got != 2 {
		t.Errorf("evenodd over two rings wound alike gives %d bands, want 2", got)
	}
}

func TestFillPathNonZeroRingsOpposite(t *testing.T) {
	// The same two rings wound the other way round: now nonzero sees the middle
	// wound once forward and once back, which cancels, and evenodd does not
	// care either way. A hole is a hole, and the rule that draws one is the one
	// a shape drawn as a ring actually wants.
	p := NewPath()
	p.AddEllipse(20, 20, 18, 18)
	reverseEllipse(p, 20, 20, 8, 8)
	if got := fillBands(t, p, 20, FillNonZero); got != 2 {
		t.Errorf("nonzero over rings wound opposite gives %d bands, want 2", got)
	}
	if got := fillBands(t, p, 20, FillEvenOdd); got != 2 {
		t.Errorf("evenodd over rings wound opposite gives %d bands, want 2", got)
	}
}

func TestFillPathOverlappingSubpaths(t *testing.T) {
	// Two boxes drawn on top of each other. Nonzero counts the overlap as
	// covered twice, which is still inside, so the two read as one shape;
	// evenodd counts crossings rather than cover, and an even count is out, so
	// the overlap becomes a hole. A row through the overlap has four edges to
	// cross, and the rules part company there.
	over := NewPath()
	over.AddRect(4, 4, 12, 12)
	over.AddRect(10, 10, 12, 12)
	if got := fillBands(t, over, 12, FillNonZero); got != 1 {
		t.Errorf("nonzero over two overlapping boxes gives %d bands, want 1", got)
	}
	if got := fillBands(t, over, 12, FillEvenOdd); got != 2 {
		t.Errorf("evenodd over two overlapping boxes gives %d bands, want 2", got)
	}
}

func TestFillPathSelfCrossing(t *testing.T) {
	// A shape drawn across itself, once each way. The hourglass has one piece
	// on every row whichever rule is asked for, and the middle — where the two
	// lines pass over each other — is inside under both.
	over := NewPath()
	over.MoveTo(4, 4)
	over.LineTo(20, 4)
	over.LineTo(4, 20)
	over.LineTo(20, 20)
	over.Close()
	if got := fillBands(t, over, 8, FillNonZero); got != 1 {
		t.Errorf("nonzero over an hourglass gives %d bands, want 1", got)
	}
	if got := fillBands(t, over, 8, FillEvenOdd); got != 1 {
		t.Errorf("evenodd over an hourglass gives %d bands, want 1", got)
	}
	if got := fillBands(t, over, 16, FillNonZero); got != 1 {
		t.Errorf("nonzero below the crossing gives %d bands, want 1", got)
	}
}

func TestFillPathHonoursClip(t *testing.T) {
	cv := newTestCanvas(t, 16, 16)
	cv.SetClip(4, 4, 8, 8)
	p := NewPath()
	p.AddRect(0, 0, 16, 16)
	cv.FillPath(p, RGB(255, 255, 255), FillNonZero)

	if got := coverageAt(cv, 3, 8); got != 0 {
		t.Errorf("outside the clip = %d, want 0", got)
	}
	if got := coverageAt(cv, 8, 8); got != 255 {
		t.Errorf("inside the clip = %d, want 255", got)
	}
}

func TestFillPathClipsOffCanvas(t *testing.T) {
	cv := newTestCanvas(t, 8, 8)
	p := NewPath()
	p.AddRect(-40, -40, 60, 60)
	cv.FillPath(p, RGB(255, 255, 255), FillNonZero)

	for _, pt := range [][2]int{{0, 0}, {7, 0}, {0, 7}, {7, 7}} {
		if got := coverageAt(cv, pt[0], pt[1]); got != 255 {
			t.Errorf("corner (%d,%d) = %d, want 255", pt[0], pt[1], got)
		}
	}
}

func TestFillPathEmpty(t *testing.T) {
	cv := newTestCanvas(t, 8, 8)
	painted := RGB(255, 255, 255)

	cv.FillPath(nil, painted, FillNonZero)
	cv.FillPath(NewPath(), painted, FillNonZero)
	// A path with only a move in it covers no area, and one with a single point
	// has no area to cross.
	moved := NewPath()
	moved.MoveTo(2, 2)
	cv.FillPath(moved, painted, FillNonZero)
	single := NewPath()
	single.MoveTo(3, 3)
	single.LineTo(3, 3)
	cv.FillPath(single, painted, FillNonZero)

	if got := coverageAt(cv, 4, 4); got != 0 {
		t.Errorf("filling nothing painted %d", got)
	}
}

func TestFillPathFullyOutside(t *testing.T) {
	cv := newTestCanvas(t, 8, 8)
	p := NewPath()
	p.AddRect(100, 100, 10, 10)
	cv.FillPath(p, RGB(255, 255, 255), FillNonZero)
	if got := coverageAt(cv, 4, 4); got != 0 {
		t.Errorf("path off the canvas painted %d", got)
	}
}

func TestFillPathZeroAlpha(t *testing.T) {
	cv := newTestCanvas(t, 8, 8)
	p := NewPath()
	p.AddRect(0, 0, 8, 8)
	cv.FillPath(p, RGBA(255, 0, 0, 0), FillNonZero)
	if got := coverageAt(cv, 4, 4); got != 0 {
		t.Errorf("a colour with no alpha painted %d", got)
	}
}

func TestFillPathAlphaScalesCoverage(t *testing.T) {
	// A half-alpha colour over an edge is half as strong again, so the pixel
	// the shape half covers comes out at a quarter: the two coverages multiply.
	cv := newTestCanvas(t, 4, 4)
	p := NewPath()
	p.AddRect(0.5, 0, 2, 4)
	cv.FillPath(p, RGBA(255, 255, 255, 128), FillNonZero)

	if got := coverageAt(cv, 0, 0); !near(got, 64, 3) {
		t.Errorf("coverage = %d, want about 64", got)
	}
}

func TestFillPathNarrowFormat(t *testing.T) {
	// A canvas in a narrow format has no per-pixel alpha to carry a partial
	// blend, so the fill has to go through the same path every other draw does.
	cv, err := NewCanvasFormat(16, 16, RGB565)
	if err != nil {
		t.Fatalf("NewCanvasFormat: %v", err)
	}
	p := NewPath()
	p.AddRect(4, 4, 8, 8)
	cv.FillPath(p, RGB(255, 0, 0), FillNonZero)
	if got := cv.At(8, 8); got != RGB(255, 0, 0) {
		t.Errorf("inside = %v, want red", got)
	}
	if got := cv.At(1, 1); got.A() == 0 {
		t.Error("outside was painted")
	}
}

// reverseEllipse adds an ellipse to a path the other way round, by moving
// anticlockwise through the same four corners.
func reverseEllipse(p *Path, cx, cy, rx, ry float64) {
	p.MoveTo(cx+rx, cy)
	p.CubicTo(cx+rx, cy-kappa*ry, cx+kappa*rx, cy-ry, cx, cy-ry)
	p.CubicTo(cx-kappa*rx, cy-ry, cx-rx, cy-kappa*ry, cx-rx, cy)
	p.CubicTo(cx-rx, cy+kappa*ry, cx-kappa*rx, cy+ry, cx, cy+ry)
	p.CubicTo(cx+kappa*rx, cy+ry, cx+rx, cy+kappa*ry, cx+rx, cy)
	p.Close()
}

// fillBands fills a path and counts the runs of covered pixels along the row y,
// which is how many separate pieces of the shape a line across it crosses.
func fillBands(t *testing.T, p *Path, y int, rule FillRule) int {
	t.Helper()
	cv := newTestCanvas(t, 40, 40)
	cv.FillPath(p, RGB(255, 255, 255), rule)

	bands, inside := 0, false
	for x := range 40 {
		if now := coverageAt(cv, x, y) > 0; now != inside {
			if now {
				bands++
			}
			inside = now
		}
	}
	return bands
}

// coverageAt reads how much of a pixel a fill covered, as the red channel over
// black: a pixel fully covered comes back at 255 and one the shape never
// reached at 0.
func coverageAt(cv *Canvas, x, y int) int { return int(cv.At(x, y).R()) }

// pick is one value or the other, for spelling out what a test wants at a spot
// that is either inside or outside the shape.
func pick(cond bool, yes, no int) int {
	if cond {
		return yes
	}
	return no
}

// near reports whether a coverage is within a step or two of what the geometry
// says, because a pixel's share is a fraction and the byte it lands in is
// rounded.
func near(got, want, step int) bool {
	d := got - want
	if d < 0 {
		d = -d
	}
	return d <= step
}

// newTestCanvas is the canvas a test draws on, failing the test rather than
// panicking when it cannot be made.
func newTestCanvas(t *testing.T, w, h int) *Canvas {
	t.Helper()
	cv, err := NewCanvas(w, h)
	if err != nil {
		t.Fatalf("NewCanvas(%d, %d): %v", w, h, err)
	}
	return cv
}
