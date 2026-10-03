package canvas

import (
	"math"
	"testing"
)

// A path measures as the line a stroke is walked along it: every segment of
// every subpath, with a closed one counted back to the point it started from.

func TestPathLength(t *testing.T) {
	// A path of nothing is nothing long, a point leads nowhere, and the length
	// of a line is the distance between its ends whichever way they are written.
	var none *Path
	if got := none.Length(); got != 0 {
		t.Errorf("a path that is not there measures %v, want 0", got)
	}
	if got := NewPath().Length(); got != 0 {
		t.Errorf("an empty path measures %v, want 0", got)
	}
	point := NewPath()
	point.MoveTo(4, 4)
	if got := point.Length(); got != 0 {
		t.Errorf("a lone point measures %v, want 0", got)
	}
	line := NewPath()
	line.MoveTo(0, 0)
	line.LineTo(3, 4)
	if got := line.Length(); got != 5 {
		t.Errorf("a line of three and four measures %v, want 5", got)
	}
}

func TestPathLengthOfEverySubpath(t *testing.T) {
	// Runs of the path are measured one after another and added up: a move
	// costs nothing of its own, and two lines apart are their two lengths.
	p := NewPath()
	p.MoveTo(0, 0)
	p.LineTo(10, 0)
	p.MoveTo(0, 10)
	p.LineTo(10, 10)
	if got := p.Length(); got != 20 {
		t.Errorf("two runs of ten measure %v, want 20", got)
	}
}

func TestPathLengthOfAClosedShapeIsAllOfIt(t *testing.T) {
	// A closed shape is walked the whole way round, the side that leads back to
	// the first point included: a square of ten by ten goes forty round, and
	// the same points left open are only thirty.
	open := NewPath()
	open.AddPolyline([]Point{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, false)
	if got := open.Length(); got != 30 {
		t.Errorf("the open square measures %v, want 30", got)
	}
	closed := NewPath()
	closed.AddRect(0, 0, 10, 10)
	if got := closed.Length(); got != 40 {
		t.Errorf("the closed square measures %v, want 40", got)
	}
}

func TestPathLengthOfACurveIsTheWalkAlongIt(t *testing.T) {
	// A curve is a run of points by the time anything reads it, so its length
	// is the length of those points — the same walk a dash is cut into — and
	// that comes out as the circumference of the circle to within the straight
	// lines the flattening puts in its place.
	p := NewPath()
	p.AddCircle(0, 0, 10)
	want := 2 * math.Pi * 10
	if got := p.Length(); math.Abs(got-want) > 0.5 {
		t.Errorf("a circle of radius 10 measures %v, want about %v", got, want)
	}
}
