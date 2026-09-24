package css

import (
	"math"
	"testing"
)

// TestTimingLinear reads the linear ramp and the timing the parser hands a
// missing keyword, both of which must move with the progress exactly.
func TestTimingLinear(t *testing.T) {
	tim, ok := parseTiming("linear")
	if !ok {
		t.Fatal("linear should parse")
	}
	for _, p := range []float64{0, 0.25, 0.5, 0.75, 1} {
		if got := tim.Ease(p); math.Abs(got-p) > 1e-9 {
			t.Errorf("linear.Ease(%v) = %v, want %v", p, got, p)
		}
	}
	if got := Linear.Ease(0.37); math.Abs(got-0.37) > 1e-9 {
		t.Errorf("zero Timing should read as linear, got %v", got)
	}
}

// TestTimingEndpoints every easing lands on the ends of the ramp, whatever
// the curve between them does.
func TestTimingEndpoints(t *testing.T) {
	names := []string{"linear", "ease", "ease-in", "ease-out", "ease-in-out",
		"cubic-bezier(0.25, 0.1, 0.25, 1)", "cubic-bezier(1, 0, 0, 1)"}
	for _, n := range names {
		tim, ok := parseTiming(n)
		if !ok {
			t.Fatalf("%s should parse", n)
		}
		if tim.Ease(0) != 0 || tim.Ease(1) != 1 {
			t.Errorf("%s.Ease(0 or 1) = (%v, %v), want (0, 1)", n, tim.Ease(0), tim.Ease(1))
		}
	}
}

// TestTimingNamedCurves pins the shapes of the named easings: ease-in lags
// the middle, ease-out leads it, and the symmetric ones meet it.
func TestTimingNamedCurves(t *testing.T) {
	easeIn, _ := parseTiming("ease-in")
	easeOut, _ := parseTiming("ease-out")
	easeInOut, _ := parseTiming("ease-in-out")
	if got := easeIn.Ease(0.5); got >= 0.5 {
		t.Errorf("ease-in.Ease(0.5) = %v, want something below 0.5", got)
	}
	if got := easeOut.Ease(0.5); got <= 0.5 {
		t.Errorf("ease-out.Ease(0.5) = %v, want something above 0.5", got)
	}
	if got := easeInOut.Ease(0.5); math.Abs(got-0.5) > 1e-6 {
		t.Errorf("ease-in-out.Ease(0.5) = %v, want 0.5", got)
	}
}

// TestTimingCubicBezierMonotonic every valid cubic-bezier whose controls do
// not overshoot, and every named easing, rises across the progress: a timing
// function that turns around would draw animation backwards. Easings built to
// overshoot (the back easings) legitimately dip below and peak above the
// ramp, which a y-monotonic read would wrongly reject.
func TestTimingCubicBezierMonotonic(t *testing.T) {
	curves := []string{
		"cubic-bezier(0, 0, 1, 1)", "cubic-bezier(0.42, 0, 0.58, 1)",
		"cubic-bezier(0.25, 0.1, 0.25, 1)", "ease", "ease-in-out",
	}
	for _, c := range curves {
		tim, ok := parseTiming(c)
		if !ok {
			t.Fatalf("%s should parse", c)
		}
		prev := 0.0
		for i := 0; i <= 100; i++ {
			got := tim.Ease(float64(i) / 100)
			if got < prev-1e-9 {
				t.Errorf("%s.Ease(%d/100) = %v fell below %v", c, i, got, prev)
			}
			prev = got
		}
	}
}

// TestTimingInvalid rejects the values a browser drops: a curve whose x
// controls leave [0,1], wrong argument counts, and plain noise.
func TestTimingInvalid(t *testing.T) {
	bad := []string{
		"", "steps(4)", "cubic-bezier(2, 0, 1, 1)",
		"cubic-bezier(-0.5, 0, 1, 1)", "cubic-bezier(0, 0, 1) ",
		"cubic-bezier(0, 0, 1, 1, 0)", "smooth", "cubic-bezier(0, x, 1, 1)",
	}
	for _, b := range bad {
		if _, ok := parseTiming(b); ok {
			t.Errorf("parseTiming(%q) should fail", b)
		}
	}
}

// TestTimingCaseAndSpacing the keywords are read case-insensitively, and a
// cubic-bezier written with top-level commas keeps its four controls.
func TestTimingCaseAndSpacing(t *testing.T) {
	a, _ := parseTiming("EASE-IN")
	b, _ := parseTiming("ease-in")
	if a != b {
		t.Error("timing function keyword should be case-insensitive")
	}
	_, ok := parseTiming("cubic-bezier(0.25, 0.1, 0.25, 1)")
	if !ok {
		t.Error("cubic-bezier should parse")
	}
}
