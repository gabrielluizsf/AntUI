package canvas

import (
	"math"
	"testing"
)

// TestRingArcFullRing paints a whole ring and checks the annulus is filled
// while the hole and the area around it stay clear.
func TestRingArcFullRing(t *testing.T) {
	cv, _ := NewCanvas(40, 40)
	cx, cy, r, thick := 20, 20, 12, 4
	cv.RingArc(cx, cy, r, thick, 0, 2*math.Pi, Red)
	for _, p := range [][2]int{{cx, cy - r + 1}, {cx + r - 1, cy}, {cx, cy + r - 1}} {
		if got := cv.At(p[0], p[1]); got.A() == 0 {
			t.Errorf("ring at %v should be painted, got %v", p, got)
		}
	}
	for _, p := range [][2]int{{cx, cy}, {cx - r + 1, cy - r + 1}, {2, 2}} {
		if got := cv.At(p[0], p[1]); got.A() != 0 {
			t.Errorf("outside the annulus at %v should stay clear, got %v", p, got)
		}
	}
}

// TestRingArcSweepClipsPaintsAroundAngle paints a quarter ring and checks the
// pixels on the covered side are filled and the uncovered side is not.
func TestRingArcSweepClipped(t *testing.T) {
	cv, _ := NewCanvas(40, 40)
	cx, cy, r, thick := 20, 20, 12, 4
	// Start at the top and sweep clockwise a quarter turn: the top-right
	// diagonal is inside the arc; the bottom-right diagonal stays empty.
	cv.RingArc(cx, cy, r, thick, -math.Pi/2, math.Pi/2, Red)
	// In the sweep: one pixel in on the covered diagonal.
	if got := cv.At(cx+8, cy-8); got.A() == 0 {
		t.Errorf("inside the arc at (28,12) should be painted, got %v", got)
	}
	// Outside the sweep: one pixel in on the uncovered diagonal, and the hole.
	for _, p := range [][2]int{{cx + 8, cy + 8}, {cx, cy}} {
		if got := cv.At(p[0], p[1]); got.A() != 0 {
			t.Errorf("outside the arc at %v should stay clear, got %v", p, got)
		}
	}
}
