package canvas

import "testing"

func TestFillEllipseSpansEachRadius(t *testing.T) {
	cv := newCanvas(t, 60, 40)
	cv.Clear(Black)

	cv.FillEllipse(30, 20, 20, 8, White)

	// Along the horizontal axis out to rx and along the vertical out to ry.
	if got := cv.At(30+18, 20); got != White {
		t.Errorf("inside rx = %08X, want white", uint32(got))
	}
	if got := cv.At(30, 20+6); got != White {
		t.Errorf("inside ry = %08X, want white", uint32(got))
	}
	// Beyond ry but well inside rx: empty for a wide ellipse.
	if got := cv.At(30, 20+12); got != Black {
		t.Errorf("outside ry = %08X, want black", uint32(got))
	}
	if got := cv.At(30+28, 20); got != Black {
		t.Errorf("outside rx = %08X, want black", uint32(got))
	}
}

func TestFillEllipseMatchesCircleWhenRadiiEqual(t *testing.T) {
	ellipse := newCanvas(t, 40, 40)
	circle := newCanvas(t, 40, 40)
	ellipse.Clear(Black)
	circle.Clear(Black)

	ellipse.FillEllipse(20, 20, 11, 11, White)
	circle.FillCircle(20, 20, 11, White)

	for i := range ellipse.Pixels {
		if ellipse.Pixels[i] != circle.Pixels[i] {
			t.Fatalf("pixel %d differs: ellipse %08X, circle %08X",
				i, uint32(ellipse.Pixels[i]), uint32(circle.Pixels[i]))
		}
	}
}

func TestFillRoundRectXYFillsAsymmetricCorners(t *testing.T) {
	cv := newCanvas(t, 60, 40)
	cv.Clear(Black)

	cv.FillRoundRectXY(10, 10, 40, 20, 12, 4, White)

	if got := cv.At(30, 20); got != White {
		t.Errorf("centre = %08X, want white", uint32(got))
	}
	// The wide rx corner reaches farther along x than the shallow ry corner
	// does along y.
	if got := cv.At(10+10, 10+2); got != White {
		t.Errorf("shallow top edge = %08X, want white", uint32(got))
	}
	if got := cv.At(10+1, 10+1); got != Black {
		t.Errorf("clipped corner = %08X, want black", uint32(got))
	}
}

func TestRoundRectXYDrawsOutlineOnly(t *testing.T) {
	cv := newCanvas(t, 60, 40)
	cv.Clear(Black)

	cv.RoundRectXY(10, 10, 40, 20, 10, 5, White)

	if got := cv.At(30, 10); got != White {
		t.Errorf("top edge = %08X, want white", uint32(got))
	}
	if got := cv.At(30, 20); got != Black {
		t.Errorf("interior = %08X, want black", uint32(got))
	}
	if got := cv.At(10, 20); got != White {
		t.Errorf("left edge = %08X, want white", uint32(got))
	}
}
