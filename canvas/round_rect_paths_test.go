package canvas

import "testing"

// FillRoundRect and RoundRect used to be two functions each: one for a single
// radius all round and one for a radius each way. They are now one function
// taking a [Rounded], which picks between the two the way the old pair did —
// and the choice is not free, because [discCoverage] and [ellipseCoverage] do
// NOT agree where a pixel sits between an ellipse and a disc of the same
// radius. Equal radii therefore have to keep taking the disc path or every
// rounded corner in the library changes by a shade.
//
// These hold the unified painters to the two they replaced.

// TestFlatRadiiStillUseTheDisc: with the radii equal, a corner must land on the
// disc's coverage, not the ellipse's. The two disagree on thousands of points,
// so this is the difference between a refactor and a visual regression.
func TestFlatRadiiStillUseTheDisc(t *testing.T) {
	const size = 48
	for _, radius := range []int{1, 2, 3, 5, 8, 13, 24} {
		a := newCanvas(t, size, size)
		a.Clear(Black)
		a.FillRoundRect(Round(Area{4, 4, size - 8, size - 8}, radius), White)

		b := newCanvas(t, size, size)
		b.Clear(Black)
		fillRoundRectDiscPath(b, 4, 4, size-8, size-8, radius, White)

		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				if a.At(x, y) != b.At(x, y) {
					t.Fatalf("r=%d: fill differs at %d,%d: %08X vs %08X",
						radius, x, y, uint32(a.At(x, y)), uint32(b.At(x, y)))
				}
			}
		}
	}
}

// fillRoundRectDiscPath is the body of the old single-radius FillRoundRect:
// the rects it laid down and the disc coverage along its four corners.
func fillRoundRectDiscPath(cv *Canvas, x, y, w, h, r int, c Color) {
	r = min(r, min(w, h)/2)
	if r <= 0 {
		cv.FillRect(Area{x, y, w, h}, c)
		return
	}
	baseAlpha := int(c.A())
	if baseAlpha == 0 {
		return
	}
	cv.FillRect(Area{x + r, y, w - 2*r, h}, c)
	cv.FillRect(Area{x, y + r, r, h - 2*r}, c)
	cv.FillRect(Area{x + w - r, y + r, r, h - 2*r}, c)

	for iy := 0; iy < h; iy++ {
		y += 0
		if iy >= r && iy < h-r {
			continue
		}
		var dy int
		if iy < r {
			dy = iy - (r - 1)
		} else {
			dy = iy - (h - r)
		}
		for ix := 0; ix < w; ix++ {
			if ix >= r && ix < w-r {
				continue
			}
			var dx int
			if ix < r {
				dx = ix - (r - 1)
			} else {
				dx = ix - (w - r)
			}
			alpha := discCoverage(dx, dy, r-1)
			if alpha == 0 {
				continue
			}
			cv.Pixel(x+ix, y+iy, Fade(c, alpha*baseAlpha/255))
		}
	}
}

// TestEllipseRadiiStillUseTheEllipse: the other side of the choice. A corner
// wider one way than the other must go through the ellipse.
func TestEllipseRadiiStillUseTheEllipse(t *testing.T) {
	const size = 60
	rx, ry := 14, 4
	a := newCanvas(t, size, size)
	a.Clear(Black)
	a.FillRoundRect(Rounded{Area{4, 4, size - 8, size - 8}, rx, ry}, White)

	b := newCanvas(t, size, size)
	b.Clear(Black)
	fillRoundRectEllipsePath(b, 4, 4, size-8, size-8, rx, ry, White)

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if a.At(x, y) != b.At(x, y) {
				t.Fatalf("rx=%d ry=%d: fill differs at %d,%d: %08X vs %08X",
					rx, ry, x, y, uint32(a.At(x, y)), uint32(b.At(x, y)))
			}
		}
	}
}

// fillRoundRectEllipsePath is the body of the old FillRoundRectXY.
func fillRoundRectEllipsePath(cv *Canvas, x, y, w, h, rx, ry int, c Color) {
	rx, ry = min(rx, w/2), min(ry, h/2)
	if w <= 0 || h <= 0 {
		return
	}
	if rx <= 0 || ry <= 0 {
		cv.FillRect(Area{x, y, w, h}, c)
		return
	}
	baseAlpha := int(c.A())
	if baseAlpha == 0 {
		return
	}
	cv.FillRect(Area{x + rx, y, w - 2*rx, h}, c)
	cv.FillRect(Area{x, y + ry, rx, h - 2*ry}, c)
	cv.FillRect(Area{x + w - rx, y + ry, rx, h - 2*ry}, c)

	for iy := 0; iy < h; iy++ {
		if iy >= ry && iy < h-ry {
			continue
		}
		var dy int
		if iy < ry {
			dy = iy - (ry - 1)
		} else {
			dy = iy - (h - ry)
		}
		for ix := 0; ix < w; ix++ {
			if ix >= rx && ix < w-rx {
				continue
			}
			var dx int
			if ix < rx {
				dx = ix - (rx - 1)
			} else {
				dx = ix - (w - rx)
			}
			alpha := ellipseCoverage(dx, dy, rx-1, ry-1)
			if alpha == 0 {
				continue
			}
			cv.Pixel(x+ix, y+iy, Fade(c, alpha*baseAlpha/255))
		}
	}
}

// TestDiscAndEllipseCoverageDisagree is the reason the two paths stay apart,
// stated as a fact rather than left in a comment: a reader who "simplifies" the
// switch into a single call should see this fail first.
func TestDiscAndEllipseCoverageDisagree(t *testing.T) {
	diff := 0
	for r := 1; r <= 40; r++ {
		for dx := -r - 2; dx <= r+2; dx++ {
			for dy := -r - 2; dy <= r+2; dy++ {
				if discCoverage(dx, dy, r) != ellipseCoverage(dx, dy, r, r) {
					diff++
				}
			}
		}
	}
	if diff == 0 {
		t.Fatal("discCoverage and ellipseCoverage agree everywhere for rx==ry; " +
			"if that ever becomes true the two paths can be one again")
	}
	t.Logf("the two disagree on %d points — the switch is load-bearing", diff)
}
