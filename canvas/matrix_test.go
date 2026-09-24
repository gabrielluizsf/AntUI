package canvas

import (
	"math"
	"testing"
)

func nearly(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestMatrixRotateMapsClockwise(t *testing.T) {
	m := Rotate(math.Pi / 2)
	x, y := m.Map(1, 0)
	if !nearly(x, 0) || !nearly(y, 1) {
		t.Errorf("rotate(90deg) moved (1,0) to (%v,%v), want (0,1)", x, y)
	}
}

func TestMatrixMulAppliesRightHandSideFirst(t *testing.T) {
	// Translate then rotate, the CSS order of the list.
	m := Rotate(math.Pi / 2).Mul(Translate(10, 0))
	x, y := m.Map(0, 0)
	if !nearly(x, 0) || !nearly(y, 10) {
		t.Errorf("translate then rotate put the origin at (%v,%v), want (0,10)", x, y)
	}

	// Rotate then translate.
	m = Translate(10, 0).Mul(Rotate(math.Pi / 2))
	x, y = m.Map(0, 0)
	if !nearly(x, 10) || !nearly(y, 0) {
		t.Errorf("rotate then translate put the origin at (%v,%v), want (10,0)", x, y)
	}
}

func TestMatrixInverseRoundTrip(t *testing.T) {
	base := Translate(3, -5).Mul(Rotate(0.3)).Mul(Scale(2, 0.5)).Mul(Skew(0.2, -0.1))
	iv, ok := base.Inverse()
	if !ok {
		t.Fatal("a full-rank matrix must invert")
	}
	for _, p := range [][2]float64{{0, 0}, {1, 2}, {-7, 13}} {
		x, y := base.Map(p[0], p[1])
		xi, yi := iv.Map(x, y)
		if !nearly(xi, p[0]) || !nearly(yi, p[1]) {
			t.Errorf("inverse of (%v,%v) returned (%v,%v)", p[0], p[1], xi, yi)
		}
	}
}

func TestMatrixSkewKeepsBoundedSlopes(t *testing.T) {
	// A shear beyond the vertical could not be mapped back; both engines clamp
	// the tangent before it sees infinity.
	m := Skew(math.Pi/2-0.0001, math.Pi/2-0.5)
	if _, ok := m.Inverse(); !ok {
		t.Errorf("shears near the vertical must stay invertible, got degenerate matrix %v", m)
	}
}

func TestMatrixInverseSingular(t *testing.T) {
	if _, ok := Scale(0, 1).Inverse(); ok {
		t.Error("scale(0,1) must not invert")
	}
	if _, ok := Identity().Inverse(); !ok {
		t.Error("identity must invert")
	}
}

// bwin builds a blank canvas to read compositing from.
func bwin(t *testing.T, w, h int) *Canvas {
	t.Helper()
	cv, err := NewCanvas(w, h)
	if err != nil {
		t.Fatal(err)
	}
	cv.FillRect(0, 0, w, h, RGBA(30, 40, 50, 255))
	return cv
}

func TestBlitMatrixTranslate(t *testing.T) {
	src, _ := NewCanvas(5, 5)
	src.FillRect(0, 0, 5, 5, Red)
	dst := bwin(t, 40, 40)
	dst.BlitMatrix(src, Translate(10, 10), Area{0, 0, 5, 5})

	if got := dst.At(12, 12); got != Red {
		t.Errorf("pixel under the blit = %v, want red", got)
	}
	if got := dst.At(9, 9); got == Red {
		t.Error("pixels before the shift must stay untouched")
	}
}

func TestBlitMatrixScaleKeepsCrisp(t *testing.T) {
	src, _ := NewCanvas(5, 5)
	src.FillRect(0, 0, 5, 5, Red)
	dst := bwin(t, 40, 40)
	dst.BlitMatrix(src, Scale(2, 2), Area{0, 0, 5, 5})

	for _, p := range [][2]int{{0, 0}, {4, 4}, {8, 8}} {
		if got := dst.At(p[0], p[1]); got != Red {
			t.Errorf("scaled pixel at %v = %v, want red", p, got)
		}
	}
	if got := dst.At(10, 10); got == Red {
		t.Error("pixels past twice the source edge must stay blank")
	}
}

func TestBlitMatrixRotate(t *testing.T) {
	src, _ := NewCanvas(20, 20)
	src.FillRect(0, 0, 20, 20, Red)
	dst := bwin(t, 40, 40)
	// Rotate the square about its own centre, then slide it until its new
	// centre sits at (15,15): the vertices land at the four cardinal points
	// around it, about 14 pixels out.
	dst.BlitMatrix(src, Translate(15, 15).Mul(Rotate(math.Pi/4)).Mul(Translate(-10, -10)), Area{0, 0, 20, 20})

	if got := dst.At(15, 15); got != Red {
		t.Errorf("centre of the rotated square = %v, want red", got)
	}
	if got := dst.At(20, 18); got != Red {
		t.Errorf("inside the rotated square = %v, want red", got)
	}
	if got := dst.At(15, 27); got != Red {
		t.Errorf("near the south vertex of the rotated square = %v, want red", got)
	}
	// A point a full diagonal away from the centre is outside the square.
	if got := dst.At(7, 7); got != RGBA(30, 40, 50, 255) {
		t.Errorf("outside the rotated square = %v, want blank", got)
	}
	if got := dst.At(35, 15); got != RGBA(30, 40, 50, 255) {
		t.Errorf("east of the rotated square = %v, want blank", got)
	}
}

func TestBlitMatrixRegionRestrictsSource(t *testing.T) {
	src, _ := NewCanvas(10, 10)
	src.FillRect(0, 0, 10, 10, Red)
	dst := bwin(t, 40, 40)
	// Only the left half of the source is in the region, so the composited
	// footprint stops at half the scale.
	dst.BlitMatrix(src, Scale(2, 1), Area{0, 0, 5, 10})

	if got := dst.At(4, 2); got != Red {
		t.Errorf("inside the region = %v, want red", got)
	}
	if got := dst.At(14, 2); got == Red {
		t.Error("outside the region the source must not bleed through")
	}
}

func TestBlitMatrixSingularCompositesNothing(t *testing.T) {
	src, _ := NewCanvas(5, 5)
	src.FillRect(0, 0, 5, 5, Red)
	dst := bwin(t, 40, 40)
	dst.BlitMatrix(src, Scale(0, 0), Area{0, 0, 5, 5})
	if got := dst.At(2, 2); got != RGBA(30, 40, 50, 255) {
		t.Errorf("a collapsed scale must paint nothing, got %v", got)
	}
}

func TestBlitBackdropReadsForward(t *testing.T) {
	dst, _ := NewCanvas(20, 20)
	src, _ := NewCanvas(20, 20)
	src.FillRect(10, 0, 4, 4, Blue)

	dst.BlitBackdrop(src, Translate(5, 0), Area{0, 0, 20, 20})

	// dst point p shows src at m(p) = p+5: the block at 10..13 shifts left.
	if got := dst.At(6, 1); got != Blue {
		t.Errorf("forward-mapped pixel at (6,1) = %v, want blue", got)
	}
	if got := dst.At(12, 1); got != 0 {
		t.Errorf("point 5px away from the block must read blank, got %v", got)
	}
	if got := dst.At(8, 3); got != Blue {
		t.Errorf("the block's far edge at (8,3) = %v, want blue", got)
	}
}

func TestDirtyTracksWrites(t *testing.T) {
	cv, _ := NewCanvas(20, 20)
	if cv.Dirty != (Area{}) {
		t.Errorf("a fresh canvas must be undirty, got %+v", cv.Dirty)
	}

	cv.FillRect(2, 3, 10, 4, Red)
	if got := cv.Dirty; got != (Area{X: 2, Y: 3, Width: 10, Height: 4}) {
		t.Errorf("FillRect dirty = %+v, want the clipped fill", got)
	}

	cv.Pixel(30, 40, Red) // outside the clip: ignored, and not dirtied
	if got := cv.Dirty; got != (Area{X: 2, Y: 3, Width: 10, Height: 4}) {
		t.Errorf("an out-of-clip write must not widen the dirty bounds, got %+v", got)
	}

	cv.Pixel(0, 0, Red)
	if got := cv.Dirty; got != (Area{X: 0, Y: 0, Width: 12, Height: 7}) {
		t.Errorf("dirty after a corner write = %+v, want 0,0..12x7", got)
	}

	cv.ResetDirty()
	if cv.Dirty != (Area{}) {
		t.Errorf("ResetDirty must empty the bounds, got %+v", cv.Dirty)
	}
}
