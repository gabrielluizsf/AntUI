package svg

import (
	"fmt"
	"math"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// The shapes are read from the attributes a drawing gives them, and the points
// that come out are the points that were written down.

func TestRectShape(t *testing.T) {
	p := rectShape(mustElement(t, `<rect x="10" y="20" width="100" height="40"/>`), t.Fatalf)
	if p == nil {
		t.Fatal("a rectangle came out as no shape at all")
	}
	// The four corners, and nothing else, since there is nothing to round.
	if got := countPoints(p); got != 5 {
		t.Errorf("got %d points, want the four corners and the one the close went back to", got)
	}
	if !hasPoint(p, canvas.Point{X: 10, Y: 20}) {
		t.Error("the corner at x=10 y=20 is not among the points")
	}
	if !hasPoint(p, canvas.Point{X: 110, Y: 60}) {
		t.Error("the corner at x=110 y=60 is not among the points")
	}
	if boundsOf(p).MaxX > 110.001 {
		t.Errorf("the shape reaches x=%v, want it to stop at 110", boundsOf(p).MaxX)
	}
}

func TestRectShapeRoundsItsCorners(t *testing.T) {
	// A radius adds points along each corner, so a rounded rectangle has more
	// points than it has corners and still reaches no further than it did.
	plain := rectShape(mustElement(t, `<rect x="0" y="0" width="100" height="40"/>`), t.Fatalf)
	round := rectShape(mustElement(t, `<rect x="0" y="0" width="100" height="40" rx="5"/>`), t.Fatalf)
	if round == nil {
		t.Fatal("a rounded rectangle came out as no shape at all")
	}
	if countPoints(round) <= countPoints(plain) {
		t.Errorf("a rounded rectangle has %d points, want more than the %d of a plain one",
			countPoints(round), countPoints(plain))
	}
	if b := boundsOf(round); b.MaxX > 100.001 || b.MaxY > 40.001 {
		t.Errorf("it reaches %v, want no further than the 100x40 it was given", b)
	}
	// A radius given once is the radius of both corners; given twice, each is
	// the one it was given, and the two spell the same rounded rectangle.
	both := rectShape(mustElement(t, `<rect x="0" y="0" width="100" height="40" rx="5" ry="5"/>`), t.Fatalf)
	if !samePoints(round, both) {
		t.Error("rx alone and rx with ry are not the same shape")
	}
}

func TestRectShapePullsBackTooLargeARadius(t *testing.T) {
	// A radius bigger than half the side it curves along would turn the corners
	// out past each other, so it is pulled back to half.
	for _, rx := range []string{"1000", "50", "60", "99999"} {
		p := rectShape(mustElement(t, `<rect x="0" y="0" width="100" height="40" rx="`+rx+`"/>`), t.Fatalf)
		if p == nil {
			t.Fatalf("rx=%s: the rectangle came out as no shape at all", rx)
		}
		b := boundsOf(p)
		if b.MinX < -0.001 || b.MinY < -0.001 || b.MaxX > 100.001 || b.MaxY > 40.001 {
			t.Errorf("rx=%s: it reaches %v, which is outside the rectangle itself", rx, b)
		}
	}
}

func TestRectShapeWithNoArea(t *testing.T) {
	// A rectangle with no width or height is not a shape, and there is nothing
	// to draw.
	for _, src := range []string{
		`<rect x="0" y="0" width="0" height="40"/>`,
		`<rect x="0" y="0" width="40" height="0"/>`,
		`<rect x="0" y="0" width="-40" height="40"/>`,
		`<rect x="0" y="0"/>`,
	} {
		if p := rectShape(mustElement(t, src), t.Fatalf); p != nil {
			t.Errorf("%s: got a shape, want none", src)
		}
	}
}

func TestCircleShape(t *testing.T) {
	p := circleShape(mustElement(t, `<circle cx="50" cy="50" r="20"/>`), t.Fatalf)
	if p == nil {
		t.Fatal("a circle came out as no shape at all")
	}
	b := boundsOf(p)
	if b.MinX < 29.9 || b.MaxX > 70.1 || b.MinY < 29.9 || b.MaxY > 70.1 {
		t.Errorf("it reaches %v, want it to be the 20 about 50,50 it was given", b)
	}
}

func TestEllipseShape(t *testing.T) {
	p := ellipseShape(mustElement(t, `<ellipse cx="50" cy="50" rx="30" ry="10"/>`), t.Fatalf)
	if p == nil {
		t.Fatal("an ellipse came out as no shape at all")
	}
	b := boundsOf(p)
	if b.MaxX-b.MinX < 59 || b.MaxY-b.MinY < 19 {
		t.Errorf("it is %v across, want the 60 by 20 it was given", b)
	}
	if b.MaxY-b.MinY > 21 {
		t.Errorf("it is %v tall, want 20", b)
	}
}

func TestCircleAndEllipseWithNoArea(t *testing.T) {
	for _, src := range []string{
		`<circle cx="0" cy="0" r="0"/>`,
		`<circle cx="0" cy="0" r="-5"/>`,
		`<circle cx="0" cy="0"/>`,
		`<ellipse cx="0" cy="0" rx="0" ry="5"/>`,
		`<ellipse cx="0" cy="0" rx="5" ry="0"/>`,
		`<ellipse cx="0" cy="0"/>`,
	} {
		e := mustElement(t, src)
		if p := circleShape(e, t.Fatalf); p != nil {
			t.Errorf("%s: circleShape got a shape, want none", src)
		}
		if p := ellipseShape(e, t.Fatalf); p != nil {
			t.Errorf("%s: ellipseShape got a shape, want none", src)
		}
	}
}

func TestLineShape(t *testing.T) {
	p := lineShape(mustElement(t, `<line x1="0" y1="0" x2="10" y2="10"/>`))
	if p == nil {
		t.Fatal("a line came out as no shape at all")
	}
	if got := countPoints(p); got != 2 {
		t.Errorf("got %d points, want the two ends of the line", got)
	}
	if !hasPoint(p, canvas.Point{X: 10, Y: 10}) {
		t.Error("the far end of the line is not among the points")
	}
}

func TestPolylineAndPolygon(t *testing.T) {
	line := pointsShape("0,0 10,0 10,10", false, t.Fatalf)
	if line == nil {
		t.Fatal("a polyline came out as no shape at all")
	}
	shape := pointsShape("0,0 10,0 10,10", true, t.Fatalf)
	if shape == nil {
		t.Fatal("a polygon came out as no shape at all")
	}
	// The same three points, the polygon carrying the one it was closed back to
	// and the polyline left open.
	if countPoints(line) != 3 {
		t.Errorf("a polyline has %d points, want the three it was given", countPoints(line))
	}
	if countPoints(shape) != 4 {
		t.Errorf("a polygon has %d points, want the three it was given and the one the close went back to",
			countPoints(shape))
	}
	if _, closed := pathClosed(line); closed {
		t.Error("a polyline was closed, want it left open")
	}
	closed := countClosed(shape)
	if closed == 0 {
		t.Error("a polygon was not closed, want it joined back to where it began")
	}
}

func TestPointsShapeReadsEitherSpelling(t *testing.T) {
	// The points may be written with commas, with spaces, or with one of each,
	// and the shape is the same either way.
	want := pointsShape("0,0 10,0 10,10", true, t.Fatalf)
	for _, s := range []string{
		"0 0, 10 0, 10 10",
		"0, 0 10, 0 10, 10",
		"  0,0   10,0   10,10  ",
		"0,-0 10,0 10,10",
	} {
		if got := pointsShape(s, true, t.Fatalf); !samePoints(got, want) {
			t.Errorf("%q gave a different shape", s)
		}
	}
}

func TestPointsShapeWithNothingToDraw(t *testing.T) {
	// A list that is not a list of points, and one that is not a whole number
	// of points, are both warnings rather than errors: the rest of the drawing
	// still paints.
	for _, s := range []string{"", "   ", "1", "0,0", "0,0 10,0 10"} {
		if p := pointsShape(s, true, quiet); p != nil {
			t.Errorf("%q: got a shape, want none", s)
		}
	}
}

func TestPathData(t *testing.T) {
	// The straight commands and the shorthand of a run of them.
	for _, tc := range []struct {
		name  string
		d     string
		point canvas.Point
	}{
		{"move and line", "M10 20L30 40", canvas.Point{X: 30, Y: 40}},
		{"pairs written together", "M10,20 30,40", canvas.Point{X: 30, Y: 40}},
		{"horizontal", "M10 20H50", canvas.Point{X: 50, Y: 20}},
		{"vertical", "M10 20V50", canvas.Point{X: 10, Y: 50}},
		{"relative", "M10 20l20 20", canvas.Point{X: 30, Y: 40}},
		{"relative horizontal", "M10 20h40", canvas.Point{X: 50, Y: 20}},
		{"relative vertical", "M10 20v30", canvas.Point{X: 10, Y: 50}},
		{"exponents", "M1e1 2E1L3 4", canvas.Point{X: 3, Y: 4}},
		{"leading dot", "M.5 .5L1.5 1.5", canvas.Point{X: 1.5, Y: 1.5}},
	} {
		p := pathData(tc.d, t.Fatalf)
		if p == nil {
			t.Errorf("%s: the path data came out as no shape at all", tc.name)
			continue
		}
		if !hasPoint(p, tc.point) {
			b := boundsOf(p)
			t.Errorf("%s: the path reaches %v, want it to go to %v", tc.name, b, tc.point)
		}
	}
}

func TestPathDataRepeatsItsCommand(t *testing.T) {
	// After a moveto a number is another moveto, and after anything else it is
	// a lineto, so a run of pairs goes on the way its letter says it does.
	again := pathData("M0 0L10 10 20 20", t.Fatalf)
	if !hasPoint(again, canvas.Point{X: 20, Y: 20}) {
		t.Errorf("a repeated L did not carry on to 20,20: %v", boundsOf(again))
	}
	// A pair of numbers after a move is a line rather than another move, so
	// this is one line and not two that start over.
	moves := pathData("M0 0 10 10", t.Fatalf)
	if got := countSubpaths(moves); got != 1 {
		t.Errorf("a pair after a move made %d subpaths, want the one it began", got)
	}
}

func TestPathDataRelative(t *testing.T) {
	// A relative command is measured from where the line is, not from the
	// origin, and the first one is measured from the origin because that is
	// where the line is.
	p := pathData("m10 10l5 5l5 5", t.Fatalf)
	if p == nil {
		t.Fatal("the path came out as no shape at all")
	}
	if !hasPoint(p, canvas.Point{X: 20, Y: 20}) {
		t.Errorf("the path reaches %v, want it to go to 20,20", boundsOf(p))
	}
}

func TestPathDataClose(t *testing.T) {
	// A close joins the line back to where it began, and whatever comes after it
	// is measured from there rather than from where the line had got to.
	p := pathData("M0 0L10 0L10 10ZL5 5", t.Fatalf)
	if p == nil {
		t.Fatal("the path came out as no shape at all")
	}
	if !hasPoint(p, canvas.Point{X: 5, Y: 5}) {
		t.Errorf("after a close, the next point is not 5,5: %v", boundsOf(p))
	}
}

func TestPathDataCurves(t *testing.T) {
	// Every curve command takes a point as its control or its end, and a curve
	// bulges outside the line between its ends, which is how one is told apart
	// from having been read as a straight one.
	for _, d := range []string{
		"M0 0C10 0 10 10 0 10",
		"M0 0Q10 10 0 20",
		"M0 0L10 0Z",
		"M0 0A5 5 0 0 1 10 0",
	} {
		p := pathData(d, t.Fatalf)
		if p == nil {
			t.Errorf("%s: the path came out as no shape at all", d)
		}
	}
	// A quadratic reaches halfway to its control point and no further, so a
	// control at x=10 carries it out to x=5 — well past the straight line
	// between its two ends, which stays on the line through them.
	curve := pathData("M0 0Q10 10 0 20", t.Fatalf)
	// The curve is flattened into short lines to be filled, which lands a hair
	// inside the true peak of five, so the bar is a bulge and not the exact tip.
	if b := boundsOf(curve); b.MaxX < 4.5 {
		t.Errorf("a curve with its control at (10,10) reaches x=%v, want it out to about 5", b.MaxX)
	}
}

func TestPathDataSmoothCurves(t *testing.T) {
	// A smooth curve's first control is the mirror of the last one of the curve
	// before it, so it goes on the other side and the line stays smooth through
	// the join. Written after a straight line there is no control to mirror, and
	// the first control is the point the line was at.
	through := pathData("M0 0C10 0 20 0 30 0S50 0 60 0", t.Fatalf)
	if through == nil {
		t.Fatal("the path came out as no shape at all")
	}
	after := pathData("M0 0L10 0S30 0 40 0", t.Fatalf)
	if after == nil {
		t.Fatal("the path came out as no shape at all")
	}
	if countPoints(through) < 3 || countPoints(after) < 3 {
		t.Errorf("got %d and %d points, want a curve with control points in each",
			countPoints(through), countPoints(after))
	}
	quad := pathData("M0 0Q10 10 20 0T40 0", t.Fatalf)
	if quad == nil {
		t.Fatal("the path came out as no shape at all")
	}
	if countPoints(quad) < 3 {
		t.Errorf("got %d points, want a smooth quadratic curve", countPoints(quad))
	}
}

func TestPathDataArcFlagsWrittenTogether(t *testing.T) {
	// The two flags of an arc may be written with nothing between them at all,
	// and the numbers after them run straight on from the second, which is what
	// makes `a1 1 0 011 1` mean a radius of one, no turn, and an end at 1,1.
	joined := pathData("M0 0a1 1 0 011 1", t.Fatalf)
	spaced := pathData("M0 0a1 1 0 0 1 1 1", t.Fatalf)
	if joined == nil || spaced == nil {
		t.Fatal("the path data came out as no shape at all")
	}
	if !samePoints(joined, spaced) {
		t.Errorf("the flags written together gave a different arc:\n%v\n%v",
			boundsOf(joined), boundsOf(spaced))
	}
}

func TestPathDataWithNothingToDraw(t *testing.T) {
	// A path with no data is not a shape. One whose data is not a list of
	// commands is a warning rather than an error, because the rest of the
	// drawing still paints.
	if p := pathData("", t.Fatalf); p != nil {
		t.Error("an empty d gave a shape")
	}
	if p := pathData("   ", t.Fatalf); p != nil {
		t.Error("a d of spaces gave a shape")
	}
	for _, d := range []string{
		"10 20",          // numbers with no command in front of them
		"M10",            // a move with half a point
		"M10 20 30",      // a move then one number and no more
		"M10 20L30",      // a line with half a point
		"M10 20X30 40",   // a command that a path does not have
		"M10 20A5 5 0 0", // an arc with too few numbers
		"M10 20C1 2 3 4", // a curve with half a control point
	} {
		if p := pathData(d, quiet); p != nil {
			t.Errorf("%q: got a shape, want none", d)
		}
	}
}

func TestPathDataWarnsAboutWhatItCannotRead(t *testing.T) {
	// A command it cannot read says so, so a caller can tell the difference
	// between a path with nothing in it and one that could not be read.
	var warnings int
	warn := func(string, ...any) { warnings++ }
	if pathData("M0 0X10 10", warn) != nil {
		t.Error("an unknown command gave a shape")
	}
	if warnings == 0 {
		t.Error("an unknown command drew nothing and said nothing")
	}
}

// box is the rectangle a shape reaches. Bounds answers its four edges as five
// numbers, one of which says whether there were any points at all, and a shape
// with no points has no box, which is an empty one here.
type box struct {
	MinX, MinY, MaxX, MaxY float64
	Empty                  bool
}

func (b box) String() string {
	if b.Empty {
		return "nothing"
	}
	return fmt.Sprintf("x %g to %g, y %g to %g", b.MinX, b.MaxX, b.MinY, b.MaxY)
}

// boundsOf is the box a path reaches, as a rectangle rather than as the five
// numbers Bounds answers them in.
func boundsOf(p *canvas.Path) box {
	minX, minY, maxX, maxY, ok := p.Bounds()
	return box{MinX: minX, MinY: minY, MaxX: maxX, MaxY: maxY, Empty: !ok}
}

// quiet is a warning that a test is expecting, so it is not worth a line of its
// own; the test that cares what was said counts them instead.
func quiet(string, ...any) {}

// Helpers for reading a path back into the points it was built from.

func mustElement(t *testing.T, src string) *element {
	t.Helper()
	e, err := parseDocument(src)
	if err != nil {
		t.Fatalf("parse %s: %v", src, err)
	}
	return e
}

func pathClosed(p *canvas.Path) (int, bool) {
	_, closed := p.Points()
	for _, c := range closed {
		if c {
			return 0, true
		}
	}
	return 0, false
}

func countClosed(p *canvas.Path) int {
	_, closed := p.Points()
	n := 0
	for _, c := range closed {
		if c {
			n++
		}
	}
	return n
}

func countPoints(p *canvas.Path) int {
	pts, _ := p.Points()
	n := 0
	for _, sub := range pts {
		n += len(sub)
	}
	return n
}

func countSubpaths(p *canvas.Path) int {
	pts, _ := p.Points()
	return len(pts)
}

func hasPoint(p *canvas.Path, want canvas.Point) bool {
	pts, _ := p.Points()
	for _, sub := range pts {
		for _, q := range sub {
			if math.Abs(q.X-want.X) < 0.01 && math.Abs(q.Y-want.Y) < 0.01 {
				return true
			}
		}
	}
	return false
}

func samePoints(a, b *canvas.Path) bool {
	if a == nil || b == nil {
		return a == b
	}
	ap, _ := a.Points()
	bp, _ := b.Points()
	if len(ap) != len(bp) {
		return false
	}
	for i := range ap {
		if len(ap[i]) != len(bp[i]) {
			return false
		}
		for j := range ap[i] {
			if math.Abs(ap[i][j].X-bp[i][j].X) > 0.001 || math.Abs(ap[i][j].Y-bp[i][j].Y) > 0.001 {
				return false
			}
		}
	}
	return true
}
