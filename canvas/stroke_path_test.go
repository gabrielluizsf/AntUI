package canvas

import (
	"math"
	"testing"
)

func TestStrokePathHorizontalLine(t *testing.T) {
	cv := newTestCanvas(t, 40, 20)
	p := NewPath()
	p.MoveTo(10, 10)
	p.LineTo(30, 10)
	cv.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: 4})

	// A four-wide line through row 10 covers rows 8, 9, 10 and 11, and nothing
	// above or below them.
	for y := range 20 {
		want := 0
		if y >= 8 && y < 12 {
			want = 255
		}
		if got := coverageAt(cv, 20, y); got != want {
			t.Errorf("row %d = %d, want %d", y, got, want)
		}
	}
	if got := coverageAt(cv, 5, 10); got != 0 {
		t.Errorf("before the line = %d, want 0", got)
	}
}

func TestStrokePathWidthAcross(t *testing.T) {
	cv := newTestCanvas(t, 40, 20)
	p := NewPath()
	p.MoveTo(20, 2)
	p.LineTo(20, 18)
	cv.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: 6})

	// The stroke is centred on the line, so it reaches half a width to either
	// side of it.
	for x := range 40 {
		want := 0
		if x >= 17 && x < 23 {
			want = 255
		}
		if got := coverageAt(cv, x, 10); got != want {
			t.Errorf("column %d = %d, want %d", x, got, want)
		}
	}
}

func TestStrokePathAntialiasedEdge(t *testing.T) {
	cv := newTestCanvas(t, 20, 20)
	p := NewPath()
	p.MoveTo(2, 10.5)
	p.LineTo(18, 10.5)
	cv.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: 4})

	// The line runs down the middle of a row, so the two rows it reaches are
	// each half covered. A stroke with a hard rim here would be the one drawing
	// a staircase on every horizontal line in the picture.
	if got := coverageAt(cv, 10, 8); !near(got, 128, 3) {
		t.Errorf("coverage on the half-covered row = %d, want about 128", got)
	}
	if got := coverageAt(cv, 10, 10); got != 255 {
		t.Errorf("coverage on the covered row = %d, want 255", got)
	}
}

func TestStrokePathZeroWidth(t *testing.T) {
	cv := newTestCanvas(t, 20, 20)
	p := NewPath()
	p.MoveTo(2, 10)
	p.LineTo(18, 10)
	cv.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: 0})
	cv.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: -3})
	if got := coverageAt(cv, 10, 10); got != 0 {
		t.Errorf("a stroke of no width painted %d", got)
	}
}

func TestStrokePathOpenLineStopsAtEnds(t *testing.T) {
	cv := newTestCanvas(t, 40, 20)
	p := NewPath()
	p.MoveTo(10, 10)
	p.LineTo(20, 10)
	cv.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: 4})

	// A butt cap ends on the end point, so nothing is painted past it.
	if got := coverageAt(cv, 9, 10); got != 0 {
		t.Errorf("before a butt cap = %d, want 0", got)
	}
	if got := coverageAt(cv, 20, 10); got != 0 {
		t.Errorf("at a butt cap = %d, want 0", got)
	}
	if got := coverageAt(cv, 15, 10); got != 255 {
		t.Errorf("along the line = %d, want 255", got)
	}
}

func TestStrokePathSquareCap(t *testing.T) {
	cv := newTestCanvas(t, 40, 20)
	p := NewPath()
	p.MoveTo(10, 10)
	p.LineTo(20, 10)
	cv.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: 4, Cap: CapSquare})

	// A square cap carries the stroke half a width past the end, which is what
	// keeps a dashed line the same length however it happens to end.
	if got := coverageAt(cv, 21, 10); got != 255 {
		t.Errorf("past a square cap = %d, want 255", got)
	}
	if got := coverageAt(cv, 22, 10); got != 0 {
		t.Errorf("past the square cap = %d, want 0", got)
	}
	if got := coverageAt(cv, 9, 10); got != 255 {
		t.Errorf("before a square cap = %d, want 255", got)
	}
}

func TestStrokePathRoundCap(t *testing.T) {
	cv := newTestCanvas(t, 40, 20)
	p := NewPath()
	p.MoveTo(10, 10)
	p.LineTo(20, 10)
	cv.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: 4, Cap: CapRound})

	// A round cap bulges a half-disc past the end, so the row on the line is
	// covered further out than a square one reaches, and the row above is
	// covered a little as the disc curves away. The disc is round, so the last
	// column of it is only part covered.
	if got := coverageAt(cv, 21, 10); got < 200 {
		t.Errorf("on the line past a round cap = %d, want nearly full", got)
	}
	if got := coverageAt(cv, 22, 10); got != 0 {
		t.Errorf("beyond a round cap = %d, want 0", got)
	}
	if got := coverageAt(cv, 21, 9); got == 0 {
		t.Error("a round cap is flat above the line, so it is not a round cap")
	}
}

func TestStrokePathLonePoint(t *testing.T) {
	cv := newTestCanvas(t, 20, 20)
	p := NewPath()
	p.MoveTo(10, 10)
	cv.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: 6, Cap: CapRound})

	if got := coverageAt(cv, 10, 10); got != 255 {
		t.Errorf("centre of the dot = %d, want 255", got)
	}
	if got := coverageAt(cv, 10, 5); got != 0 {
		t.Errorf("past the dot = %d, want 0", got)
	}
}

func TestStrokePathClosedRing(t *testing.T) {
	cv := newTestCanvas(t, 60, 60)
	p := NewPath()
	p.AddEllipse(30, 30, 20, 20)
	cv.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: 4})

	// On the ring, on it and inside it, but not in the middle of it: a closed
	// shape stroked is an outline, and the middle is left as it was.
	if got := coverageAt(cv, 30, 10); got != 255 {
		t.Errorf("on the ring = %d, want 255", got)
	}
	if got := coverageAt(cv, 30, 30); got != 0 {
		t.Errorf("inside the ring = %d, want 0", got)
	}
	if got := coverageAt(cv, 0, 0); got != 0 {
		t.Errorf("outside the ring = %d, want 0", got)
	}
}

func TestStrokePathMiterJoin(t *testing.T) {
	// A right angle, stroked with a miter: the corner is cut square, so the
	// point of the miter reaches a half-width beyond both outer edges.
	cv := newTestCanvas(t, 60, 60)
	p := NewPath()
	p.MoveTo(10, 40)
	p.LineTo(30, 20)
	p.LineTo(50, 40)
	cv.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: 4, Join: JoinMiter})

	if got := coverageAt(cv, 30, 20); got != 255 {
		t.Errorf("the corner itself = %d, want 255", got)
	}
	// The miter of a right angle reaches √2 half-widths out from the corner
	// along the bisector, which for a four-wide stroke is 2.8 pixels, putting
	// the tip a little above y=17. The triangle it fills is a shallow one, so
	// the row it ends in is only part covered.
	if got := coverageAt(cv, 30, 17); got == 0 {
		t.Error("the miter did not reach past the corner")
	}
	if got := coverageAt(cv, 30, 16); got != 0 {
		t.Errorf("past the miter spike = %d, want 0", got)
	}
}

func TestStrokePathBevelJoin(t *testing.T) {
	// The same corner beveled: the spike is gone, and the corner is cut
	// straight across between the two outer edges.
	cv := newTestCanvas(t, 60, 60)
	p := NewPath()
	p.MoveTo(10, 40)
	p.LineTo(30, 20)
	p.LineTo(50, 40)
	cv.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: 4, Join: JoinBevel})

	if got := coverageAt(cv, 30, 17); got != 0 {
		t.Errorf("above a bevel = %d, want 0", got)
	}
	if got := coverageAt(cv, 30, 20); got != 255 {
		t.Errorf("the corner itself = %d, want 255", got)
	}
}

func TestStrokePathRoundJoin(t *testing.T) {
	cv := newTestCanvas(t, 60, 60)
	p := NewPath()
	p.MoveTo(10, 40)
	p.LineTo(30, 20)
	p.LineTo(50, 40)
	cv.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: 4, Join: JoinRound})

	if got := coverageAt(cv, 30, 20); got != 255 {
		t.Errorf("the corner itself = %d, want 255", got)
	}
	// A round join is rounded about the corner, so it bulges out past the flat
	// chord a bevel cuts and covers ground beside the corner a bevel leaves
	// bare, while still stopping short of the miter's spike.
	if got := coverageAt(cv, 29, 18); got == 0 {
		t.Error("a round join is not a round join, the arc never left the corner")
	}
	if got := coverageAt(cv, 30, 17); got != 0 {
		t.Errorf("a round join reached the miter's spike: %d", got)
	}
}

func TestStrokePathMiterLimit(t *testing.T) {
	// A corner sharper than the miter limit allows. The spike a miter would
	// draw here is many times the width of the stroke, so it is beveled
	// instead, and the far end of it is left bare.
	cv := newTestCanvas(t, 60, 60)
	p := NewPath()
	p.MoveTo(10, 30)
	p.LineTo(30, 30)
	p.LineTo(15.9, 44.1)
	cv.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: 2, Join: JoinMiter, MiterLimit: 2})

	// The turn is 135°, whose miter runs 1/cos(67.5°) half-widths, or 2.6
	// stroke widths, out to a spike past x=32. A limit of two does not reach that
	// far, so the corner is cut flat and the spike is never drawn.
	if got := coverageAt(cv, 31, 30); got != 0 {
		t.Errorf("inside a miter that should have been beveled = %d, want 0", got)
	}
	if got := coverageAt(cv, 30, 30); got == 0 {
		t.Error("the corner itself was left bare")
	}
}

func TestStrokePathMiterLimitNotReached(t *testing.T) {
	// The same corner with a limit past what it needs, which is what the default
	// of four gives.
	cv := newTestCanvas(t, 60, 60)
	p := NewPath()
	p.MoveTo(10, 30)
	p.LineTo(30, 30)
	p.LineTo(15.9, 44.1)
	cv.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: 2, Join: JoinMiter, MiterLimit: 4})

	if got := coverageAt(cv, 31, 30); got == 0 {
		t.Error("a miter within the limit left its spike bare")
	}
}

func TestStrokePathCurve(t *testing.T) {
	// A curve stroked keeps the width all the way round, which is the thing a
	// stroke drawn as a series of disconnected rectangles gets wrong at every
	// join.
	cv := newTestCanvas(t, 80, 80)
	p := NewPath()
	p.AddEllipse(40, 40, 30, 20)
	cv.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: 4})

	// The top of the ellipse, its bottom, and the two sides. The ellipse is
	// thirty wide and twenty high about (40,40), so the ring runs from y=20 to
	// y=60 and from x=10 to x=70.
	for _, pt := range [][2]int{{40, 20}, {40, 60}, {10, 40}, {70, 40}} {
		if got := coverageAt(cv, pt[0], pt[1]); got < 200 {
			t.Errorf("the ring at (%d,%d) = %d, want nearly full", pt[0], pt[1], got)
		}
	}
	for _, pt := range [][2]int{{40, 40}, {5, 5}, {75, 75}} {
		if got := coverageAt(cv, pt[0], pt[1]); got != 0 {
			t.Errorf("off the ring at (%d,%d) = %d, want 0", pt[0], pt[1], got)
		}
	}
}

func TestStrokePathNoHolesAtJoins(t *testing.T) {
	// Every piece of the outline is wound the same way round, so where two
	// pieces overlap the winding adds up instead of cancelling. A join built the
	// other way round shows up as a bright pixel next to a bare one, which is
	// what this walks for.
	cv := newTestCanvas(t, 60, 60)
	p := NewPath()
	p.MoveTo(10, 30)
	p.LineTo(30, 10)
	p.LineTo(50, 30)
	p.LineTo(50, 50)
	p.LineTo(30, 50)
	cv.StrokePath(p, RGB(255, 255, 255), StrokeStyle{Width: 6, Join: JoinRound, Cap: CapRound})

	holes := 0
	for y := range 60 {
		for x := range 60 {
			if cv.At(x, y).R() > 0 && cv.At(x, y).R() < 255 {
				holes++
			}
		}
	}
	if holes > 200 {
		t.Errorf("%d pixels are part covered, which is more than a rounded join leaves", holes)
	}
}

func TestStrokePathHonoursClip(t *testing.T) {
	cv := newTestCanvas(t, 40, 40)
	cv.SetClip(10, 10, 10, 10)
	p := NewPath()
	p.MoveTo(0, 15)
	p.LineTo(40, 15)
	cv.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: 6})

	if got := coverageAt(cv, 5, 15); got != 0 {
		t.Errorf("outside the clip = %d, want 0", got)
	}
	if got := coverageAt(cv, 15, 15); got != 255 {
		t.Errorf("inside the clip = %d, want 255", got)
	}
}

func TestStrokePathDegenerate(t *testing.T) {
	cv := newTestCanvas(t, 20, 20)
	cv.StrokePath(nil, RGB(255, 0, 0), StrokeStyle{Width: 4})
	cv.StrokePath(NewPath(), RGB(255, 0, 0), StrokeStyle{Width: 4})
	// Two points in the same place make a segment with no direction to be
	// widened along, and nothing to draw.
	same := NewPath()
	same.MoveTo(10, 10)
	same.LineTo(10, 10)
	same.LineTo(10, 10)
	cv.StrokePath(same, RGB(255, 0, 0), StrokeStyle{Width: 4})
	if got := coverageAt(cv, 10, 10); got != 0 {
		t.Errorf("a path of identical points painted %d", got)
	}
}

func TestStrokePathAlpha(t *testing.T) {
	cv := newTestCanvas(t, 20, 20)
	p := NewPath()
	p.MoveTo(2, 10)
	p.LineTo(18, 10)
	cv.StrokePath(p, RGBA(255, 255, 255, 128), StrokeStyle{Width: 4})
	if got := coverageAt(cv, 10, 10); !near(got, 128, 2) {
		t.Errorf("coverage = %d, want about 128", got)
	}
}

func TestStrokePathNarrowFormat(t *testing.T) {
	cv, err := NewCanvasFormat(20, 20, RGB565)
	if err != nil {
		t.Fatalf("NewCanvasFormat: %v", err)
	}
	p := NewPath()
	p.MoveTo(2, 10)
	p.LineTo(18, 10)
	cv.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: 4})
	if got := cv.At(10, 10); got != RGB(255, 0, 0) {
		t.Errorf("on the line = %v, want red", got)
	}
	if got := cv.At(10, 2); got.A() == 0 {
		t.Error("off the line was painted")
	}
}

func TestStrokeOutlineIsFillable(t *testing.T) {
	// The outline is what a shadow under a stroke is drawn from, so it has to
	// be a path on its own rather than something only the stroke can use.
	p := NewPath()
	p.MoveTo(4, 20)
	p.LineTo(36, 20)
	outline := StrokeOutline(p, StrokeStyle{Width: 8})
	if outline == nil || outline.Empty() {
		t.Fatal("no outline for a line")
	}
	cv := newTestCanvas(t, 40, 40)
	cv.FillPath(outline, RGB(255, 0, 0), FillNonZero)
	if got := coverageAt(cv, 20, 20); got != 255 {
		t.Errorf("inside the outline = %d, want 255", got)
	}
	if got := coverageAt(cv, 20, 12); got != 0 {
		t.Errorf("outside the outline = %d, want 0", got)
	}
}

func TestStrokeOutlineNothingToDraw(t *testing.T) {
	if got := StrokeOutline(nil, StrokeStyle{Width: 4}); got != nil {
		t.Error("an outline of nothing is something")
	}
	if got := StrokeOutline(NewPath(), StrokeStyle{Width: 4}); got != nil {
		t.Error("an outline of an empty path is something")
	}
	p := NewPath()
	p.MoveTo(0, 0)
	p.LineTo(10, 0)
	if got := StrokeOutline(p, StrokeStyle{Width: 0}); got != nil {
		t.Error("an outline of no width is something")
	}
}

func TestStrokePathMarksDirty(t *testing.T) {
	// A stroke has to say it wrote, or the frame on the screen is never told to
	// be looked at again.
	again := newTestCanvas(t, 40, 40)
	again.Compare(make([]Color, 40*40))
	p := NewPath()
	p.MoveTo(5, 20)
	p.LineTo(35, 20)
	again.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: 4})
	area := again.PresentChanges()
	if area.Width == 0 {
		t.Fatal("a stroke onto a frame of black changed nothing")
	}
	if area.X > 10 || area.X+area.Width < 30 {
		t.Errorf("changed area %v does not cover the stroke", area)
	}
}

// A stroke of a square is a square outline: the sides are as wide as asked and
// the corners meet them, which is what the miter and bevel joins between them
// are for.
func TestStrokePathRectCorners(t *testing.T) {
	for name, join := range map[string]LineJoin{"miter": JoinMiter, "round": JoinRound, "bevel": JoinBevel} {
		t.Run(name, func(t *testing.T) {
			cv := newTestCanvas(t, 60, 60)
			p := NewPath()
			p.AddRect(10, 10, 40, 40)
			cv.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: 6, Join: join})

			// The middle of each side.
			if got := coverageAt(cv, 30, 10); got != 255 {
				t.Errorf("the top side = %d, want 255", got)
			}
			if got := coverageAt(cv, 10, 30); got != 255 {
				t.Errorf("the left side = %d, want 255", got)
			}
			// The corner, whichever join drew it.
			if got := coverageAt(cv, 10, 10); got != 255 {
				t.Errorf("the corner = %d, want 255", got)
			}
			// The middle, and well outside.
			if got := coverageAt(cv, 30, 30); got != 0 {
				t.Errorf("the middle = %d, want 0", got)
			}
			if got := coverageAt(cv, 55, 55); got != 0 {
				t.Errorf("outside = %d, want 0", got)
			}
		})
	}
}

func TestStrokePathWidthScaledByMatrix(t *testing.T) {
	// A path transformed before it is stroked comes out its own width, because
	// the width belongs to the stroke and not to the line. This is the one thing
	// a path drawn with the whole canvas under a transform cannot do, and it is
	// why the transform is applied to the path instead.
	cv := newTestCanvas(t, 60, 60)
	p := NewPath()
	p.MoveTo(10, 15)
	p.LineTo(20, 15)
	p.Transform(Scale(2, 2))
	cv.StrokePath(p, RGB(255, 0, 0), StrokeStyle{Width: 4})

	// The line now runs from (20,30) to (40,30), and the stroke is still four
	// pixels wide about it rather than eight.
	for y := range 60 {
		want := 0
		if y >= 28 && y < 32 {
			want = 255
		}
		if got := coverageAt(cv, 30, y); got != want {
			t.Errorf("row %d = %d, want %d", y, got, want)
		}
	}
	if got := coverageAt(cv, 30, 50); got != 0 {
		t.Errorf("the line was not widened by the transform: %d", got)
	}
}

func TestStrokePathMathHelpers(t *testing.T) {
	// The arc a round join or cap is made of has to be fine enough that the
	// rounding does not show as facets, and coarse enough not to cost more than
	// it is worth.
	big := arcPoints(Point{}, 10, 0, math.Pi, 1)
	if len(big) < 8 {
		t.Errorf("a half circle of radius 10 is %d points, which facets", len(big))
	}
	small := arcPoints(Point{}, 1, 0, math.Pi/2, 1)
	if len(small) > 16 {
		t.Errorf("a quarter of a small circle is %d points, which is more than it is worth", len(small))
	}
	if got := unit(3, 4); got != (Point{0.6, 0.8}) {
		t.Errorf("unit(3, 4) = %v", got)
	}
	if got := unit(0, 0); got != (Point{}) {
		t.Errorf("unit of nothing = %v, want the zero direction", got)
	}
}
