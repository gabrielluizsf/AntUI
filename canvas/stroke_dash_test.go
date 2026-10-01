package canvas

import (
	"math"
	"testing"
)

// A dash pattern is a list of lengths measured in the path's own units, and the
// odd one in the list is the reason half of them need saying twice.

func TestNewDashSplitsOnAndOff(t *testing.T) {
	d, ok := NewDash([]float64{10, 5}, 0)
	if !ok {
		t.Fatal("a pair of lengths is a dash pattern")
	}
	if len(d.On) != 1 || d.On[0] != 10 {
		t.Errorf("on = %v, want [10]", d.On)
	}
	if len(d.Off) != 1 || d.Off[0] != 5 {
		t.Errorf("off = %v, want [5]", d.Off)
	}
	if got := d.period; got != 15 {
		t.Errorf("period = %v, want 15", got)
	}
}

func TestNewDashDoublesAnOddList(t *testing.T) {
	// An odd list is read twice over, so a lone number means dashes and gaps of
	// the same length, and the pattern walks in twenty of them, not ten.
	d, ok := NewDash([]float64{10}, 0)
	if !ok {
		t.Fatal("a lone length is a dash pattern")
	}
	if len(d.On) != 1 || d.On[0] != 10 {
		t.Errorf("on = %v, want [10]", d.On)
	}
	if len(d.Off) != 1 || d.Off[0] != 10 {
		t.Errorf("off = %v, want [10]", d.Off)
	}
	if got := d.period; got != 20 {
		t.Errorf("period = %v, want 20, the doubled pattern", got)
	}
}

func TestNewDashRejectsWhatIsNotAPattern(t *testing.T) {
	for name, lengths := range map[string][]float64{
		"empty":      {},
		"a zero":     {10, 0},
		"a negative": {10, -5},
		"all zero":   {0, 0},
	} {
		t.Run(name, func(t *testing.T) {
			if _, ok := NewDash(lengths, 0); ok {
				t.Error("that is not a dash pattern, but it was taken for one")
			}
		})
	}
}

func TestNewDashWrapsTheOffset(t *testing.T) {
	// An offset past the end of the pattern, or before the start of it, is where
	// it would have landed inside the pattern.
	for _, c := range []struct{ offset, want float64 }{
		{0, 0},
		{20, 5},
		{-5, 10},
		{-20, 10},
	} {
		d, ok := NewDash([]float64{10, 5}, c.offset)
		if !ok {
			t.Fatal("a pair of lengths is a dash pattern")
		}
		if d.Offset != c.want {
			t.Errorf("offset %v wrapped to %v, want %v", c.offset, d.Offset, c.want)
		}
	}
}

// The dashes a line is cut into are found by walking the pattern along it by
// distance, so what a test has to look at is where the drawn runs start and stop.

func TestDashedPathCutsByDistance(t *testing.T) {
	p := NewPath()
	p.MoveTo(0, 10)
	p.LineTo(40, 10)
	d, ok := NewDash([]float64{10, 10}, 0)
	if !ok {
		t.Fatal("a pair of lengths is a dash pattern")
	}
	pts, closed := p.Dashed(d).Points()
	if len(pts) != 2 {
		t.Fatalf("got %d runs, want 2", len(pts))
	}
	for i, run := range pts {
		if closed[i] {
			t.Errorf("run %d is closed, want it open so it can be capped", i)
		}
		if len(run) != 2 {
			t.Fatalf("run %d has %d points, want 2", i, len(run))
		}
		if got, want := run[0].X, float64(20*i); got != want {
			t.Errorf("run %d starts at %v, want %v", i, got, want)
		}
		if got, want := run[1].X, float64(20*i+10); got != want {
			t.Errorf("run %d ends at %v, want %v", i, got, want)
		}
	}
}

func TestDashedPathAnOddList(t *testing.T) {
	// A lone length reads twice over, so the runs are ten on and ten off all
	// the way down a forty-long line.
	p := NewPath()
	p.MoveTo(0, 10)
	p.LineTo(40, 10)
	d, ok := NewDash([]float64{10}, 0)
	if !ok {
		t.Fatal("a lone length is a dash pattern")
	}
	pts, _ := p.Dashed(d).Points()
	if len(pts) != 2 {
		t.Fatalf("got %d runs, want 2", len(pts))
	}
	for i, run := range pts {
		if got, want := run[0].X, float64(20*i); got != want {
			t.Errorf("run %d starts at %v, want %v", i, got, want)
		}
		if got, want := run[1].X, float64(20*i+10); got != want {
			t.Errorf("run %d ends at %v, want %v", i, got, want)
		}
	}
}

func TestDashedPathAnOffsetStartsMidPattern(t *testing.T) {
	// An offset of five begins the line in the middle of a dash, so the first
	// run is the five that are left of it.
	p := NewPath()
	p.MoveTo(0, 10)
	p.LineTo(40, 10)
	d, ok := NewDash([]float64{10, 10}, 5)
	if !ok {
		t.Fatal("a pair of lengths is a dash pattern")
	}
	pts, _ := p.Dashed(d).Points()
	if len(pts) != 2 {
		t.Fatalf("got %d runs, want 2", len(pts))
	}
	// The line opens five along, where the offset put it inside the first dash,
	// so that run is the five left of it and the next one is a whole one.
	if got, want := pts[0][0].X, 5.0; got != want {
		t.Errorf("the first run starts at %v, want %v", got, want)
	}
	if got, want := pts[0][1].X, 10.0; got != want {
		t.Errorf("the first run ends at %v, want %v", got, want)
	}
	if got, want := pts[1][0].X, 20.0; got != want {
		t.Errorf("the second run starts at %v, want %v", got, want)
	}
}

func TestDashedPathANegativeOffset(t *testing.T) {
	// A negative offset counts back from the end of the pattern, which lands
	// five before the end of it, inside the gap.
	p := NewPath()
	p.MoveTo(0, 10)
	p.LineTo(40, 10)
	d, ok := NewDash([]float64{10, 10}, -5)
	if !ok {
		t.Fatal("a pair of lengths is a dash pattern")
	}
	pts, _ := p.Dashed(d).Points()
	if len(pts) != 2 {
		t.Fatalf("got %d runs, want 2", len(pts))
	}
	// Wrapped, minus five is where five before the end of the period is, which
	// is the tail of the gap, so the line opens at the second dash.
	if got, want := pts[0][0].X, 10.0; got != want {
		t.Errorf("the first run starts at %v, want %v", got, want)
	}
}

func TestDashedPathLongerDashesMakeFewerRuns(t *testing.T) {
	// A dash longer than the line is one run the whole way.
	p := NewPath()
	p.MoveTo(0, 10)
	p.LineTo(40, 10)
	d, ok := NewDash([]float64{100, 10}, 0)
	if !ok {
		t.Fatal("a pair of lengths is a dash pattern")
	}
	pts, _ := p.Dashed(d).Points()
	if len(pts) != 1 {
		t.Fatalf("got %d runs, want 1", len(pts))
	}
	if got, want := pts[0][1].X, 40.0; got != want {
		t.Errorf("the run ends at %v, want %v", got, want)
	}
}

func TestDashedPathCutsAcrossManySegments(t *testing.T) {
	// The walk goes by distance and not by points, so a dash falls where the
	// distances add up however many segments it takes to get there.
	p := NewPath()
	p.MoveTo(0, 0)
	p.LineTo(10, 0)
	p.LineTo(10, 10)
	p.LineTo(20, 10)
	d, ok := NewDash([]float64{10, 10}, 0)
	if !ok {
		t.Fatal("a pair of lengths is a dash pattern")
	}
	pts, _ := p.Dashed(d).Points()
	if len(pts) != 2 {
		t.Fatalf("got %d runs, want 2", len(pts))
	}
	if got, want := pts[0][len(pts[0])-1].X, 10.0; got != want {
		t.Errorf("the first run ends at %v, want %v", got, want)
	}
	if got, want := pts[1][0].X, 10.0; got != want {
		t.Errorf("the second run starts at %v, want %v", got, want)
	}
	if got, want := pts[1][len(pts[1])-1].X, 20.0; got != want {
		t.Errorf("the second run ends at %v, want %v", got, want)
	}
}

func TestDashedPathAClosedShape(t *testing.T) {
	// A closed shape is walked the whole way round, the last point leading back
	// to the first, so a dash can cross the seam.
	p := NewPath()
	p.AddRect(0, 0, 20, 20)
	p.Close()
	d, ok := NewDash([]float64{10, 10}, 0)
	if !ok {
		t.Fatal("a pair of lengths is a dash pattern")
	}
	pts, _ := p.Dashed(d).Points()
	if len(pts) == 0 {
		t.Fatal("a dashed square came back with no runs on it")
	}
	total := 0.0
	for _, run := range pts {
		for i := 1; i < len(run); i++ {
			total += math.Hypot(run[i].X-run[i-1].X, run[i].Y-run[i-1].Y)
		}
	}
	if total != 40 {
		t.Errorf("%v of dash on a square of eighty, want 40", total)
	}
}

func TestDashedPathNothingToDash(t *testing.T) {
	// A path of nothing is still a path of nothing, and a line that is not
	// dashed comes back whole.
	if got := (*Path)(nil).Dashed(Dash{}); got != nil {
		t.Error("dashes of nothing is something")
	}
	empty := NewPath()
	if got := empty.Dashed(Dash{}); got == nil || !got.Empty() {
		t.Error("dashes of an empty path is not empty")
	}
	p := NewPath()
	p.MoveTo(0, 10)
	p.LineTo(40, 10)
	pts, _ := p.Dashed(Dash{}).Points()
	if len(pts) != 1 || len(pts[0]) != 2 {
		t.Errorf("a line that is not dashed came back as %v", pts)
	}
}

func TestDashedPathKeepsALonePoint(t *testing.T) {
	// A single point has no length to be cut along, so it is kept as it was and
	// the cap on it is what draws it.
	p := NewPath()
	p.MoveTo(5, 5)
	d, ok := NewDash([]float64{10, 10}, 0)
	if !ok {
		t.Fatal("a pair of lengths is a dash pattern")
	}
	pts, _ := p.Dashed(d).Points()
	if len(pts) != 1 || len(pts[0]) != 1 {
		t.Fatalf("a lone point came back as %v", pts)
	}
	if pts[0][0] != (Point{5, 5}) {
		t.Errorf("the lone point is at %v, want (5,5)", pts[0][0])
	}
}
