package svg

import (
	"math"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// A transform says where a node goes and how big it is once there. These check
// that each of the operations reads as the matrix it stands for, and that a
// list of them composes in the order it was written.

func TestParseTransformReadsEachOperation(t *testing.T) {
	for name, tc := range map[string]struct {
		in   string
		want canvas.Matrix
	}{
		"translate moves and leaves the size alone": {
			"translate(10, 20)", canvas.Translate(10, 20),
		},
		"a translate with one number is along the line of nothing": {
			"translate(10)", canvas.Translate(10, 0),
		},
		"scale on both sides": {
			"scale(2, 3)", canvas.Scale(2, 3),
		},
		"a scale with one number is even": {
			"scale(2)", canvas.Scale(2, 2),
		},
		"a quarter turn": {
			"rotate(90)", canvas.Rotate(math.Pi / 2),
		},
		"a turn written in radians": {
			"rotate(0.5rad)", canvas.Rotate(0.5),
		},
		"a turn written in grad, which is four hundred to the whole": {
			"rotate(100grad)", canvas.Rotate(math.Pi / 2),
		},
		"a skew along one side": {
			"skewX(45)", canvas.Skew(math.Pi/4, 0),
		},
		"a skew along the other": {
			"skewY(45)", canvas.Skew(0, math.Pi/4),
		},
	} {
		got, ok := parseTransform(tc.in)
		if !ok {
			t.Errorf("parseTransform(%q) did not read it", tc.in)
			continue
		}
		if !nearMatrix(got, tc.want) {
			t.Errorf("%s: parseTransform(%q) = %v, want %v", name, tc.in, got, tc.want)
		}
	}
}

func TestParseTransformRotateAboutAPoint(t *testing.T) {
	// Turning about a point is a move to it, a turn, and a move back, so the
	// point itself does not move.
	p := canvas.Point{X: 10, Y: 20}
	got, ok := parseTransform("rotate(180, 10, 20)")
	if !ok {
		t.Fatal("it did not read it")
	}
	mx, my := got.Map(p.X, p.Y)
	moved := canvas.Point{X: mx, Y: my}
	if !nearFloat(moved.X, 10, 0.001) || !nearFloat(moved.Y, 20, 0.001) {
		t.Errorf("the point it turns about moved to %v,%v, want to stay at 10,20", moved.X, moved.Y)
	}
	// A point on the other side of that point ends up the other side again.
	ox, oy := got.Map(20, 20)
	other := canvas.Point{X: ox, Y: oy}
	if !nearFloat(other.X, 0, 0.001) || !nearFloat(other.Y, 20, 0.001) {
		t.Errorf("(20,20) went to %v,%v, want it reflected to 0,20", other.X, other.Y)
	}
}

func TestParseTransformInTheOrderWritten(t *testing.T) {
	// A list of operations is one inside the next, as SVG nests them: the first
	// written is the outermost and so is the last a point goes through. A point
	// is therefore scaled first and moved after, which leaves the move at ten
	// rather than stretched to twenty by the scale.
	got, ok := parseTransform("translate(10, 0) scale(2)")
	if !ok {
		t.Fatal("it did not read it")
	}
	if m := mapX(got, 0, 0); !nearFloat(m, 10, 0.001) {
		t.Errorf("the origin went to x=%v, want it to 10", m)
	}
	if m := mapX(got, 1, 0); !nearFloat(m, 12, 0.001) {
		t.Errorf("(1,0) went to x=%v, want it to 12: a unit is two wide, moved ten along", m)
	}
	// The other way round the move is inside the scale, so the scale is what
	// stretches it and the origin lands on twenty.
	got, ok = parseTransform("scale(2) translate(10, 0)")
	if !ok {
		t.Fatal("it did not read it")
	}
	if m := mapX(got, 0, 0); !nearFloat(m, 20, 0.001) {
		t.Errorf("the origin went to x=%v, want it to 20", m)
	}
}

func TestParseTransformSeparators(t *testing.T) {
	// A list of operations may be written with no separator at all, with spaces,
	// or with commas, and all of them say the same thing.
	want, ok := parseTransform("translate(1, 2) scale(3)")
	if !ok {
		t.Fatal("it did not read it")
	}
	one, ok1 := parseTransform("translate(1,2)scale(3)")
	if !ok1 || !nearMatrix(one, want) {
		t.Errorf("with no separator between them: %v, want %v", one, want)
	}
	two, ok2 := parseTransform("translate(1 2), scale(3)")
	if !ok2 || !nearMatrix(two, want) {
		t.Errorf("with spaces and a comma: %v, want %v", two, want)
	}
}

func TestParseTransformLeavesWhatItCannotRead(t *testing.T) {
	// Nothing to read is not a transform, and a list with an operation missing
	// its numbers is one this cannot make a matrix of, so the rest is used.
	for _, s := range []string{"", "   ", "translate", "translate()", "scale(a, b)", "nonsense(3)"} {
		if m, ok := parseTransform(s); ok {
			t.Errorf("parseTransform(%q) = %v, want it not read", s, m)
		}
	}
	got, ok := parseTransform("scale(2) nonsense(3) translate(5)")
	if !ok {
		t.Fatal("a list with one operation it cannot read came out as no transform at all")
	}
	// The scale is outside the move, so the move of five is stretched to ten.
	if m := mapX(got, 0, 0); !nearFloat(m, 10, 0.001) {
		t.Errorf("the origin went to x=%v, want the operations it could read to be used", m)
	}
}

func TestParseTransformMatrix(t *testing.T) {
	// The matrix form names all six numbers of the transform at once, and is the
	// one the others are written out of.
	for _, tc := range []struct {
		in   string
		want canvas.Matrix
	}{
		{"matrix(1, 0, 0, 1, 0, 0)", canvas.Identity()},
		{"matrix(2, 0, 0, 2, 0, 0)", canvas.Scale(2, 2)},
		{"matrix(1, 0, 0, 1, 10, 20)", canvas.Translate(10, 20)},
	} {
		got, ok := parseTransform(tc.in)
		if !ok {
			t.Errorf("parseTransform(%q) did not read it", tc.in)
			continue
		}
		if !nearMatrix(got, tc.want) {
			t.Errorf("parseTransform(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestScaleOf(t *testing.T) {
	// The width a stroke is drawn at follows how much bigger a transform makes
	// things, and a transform that only turns or moves leaves it as it was.
	for name, tc := range map[string]struct {
		m    canvas.Matrix
		want float64
	}{
		"none at all":           {canvas.Identity(), 1},
		"double":                {canvas.Scale(2, 2), 2},
		"half":                  {canvas.Scale(0.5, 0.5), 0.5},
		"a turn":                {canvas.Rotate(math.Pi / 3), 1},
		"a move":                {canvas.Translate(10, 20), 1},
		"unequal on both sides": {canvas.Scale(2, 8), 4},
	} {
		if got := scaleOf(tc.m); !nearFloat(got, tc.want, 0.001) {
			t.Errorf("%s: scaleOf = %v, want %v", name, got, tc.want)
		}
	}
}

// TestScalesOf is scaleOf read one axis at a time: a transform that stretches
// one axis more than the other says so in the two numbers, where the area's
// square root would give the same number for both and lose which way it was.
func TestScalesOf(t *testing.T) {
	for name, tc := range map[string]struct {
		m     canvas.Matrix
		wantX float64
		wantY float64
	}{
		"none at all":    {canvas.Identity(), 1, 1},
		"double":         {canvas.Scale(2, 2), 2, 2},
		"a turn":         {canvas.Rotate(math.Pi / 3), 1, 1},
		"stretched wide": {canvas.Scale(2, 8), 2, 8},
		"stretched tall": {canvas.Scale(8, 2), 8, 2},
	} {
		got := scalesOf(tc.m)
		if !nearFloat(got.x, tc.wantX, 0.001) || !nearFloat(got.y, tc.wantY, 0.001) {
			t.Errorf("%s: scalesOf = %v,%v, want %v,%v", name, got.x, got.y, tc.wantX, tc.wantY)
		}
	}
}

func TestFitTransform(t *testing.T) { // A viewBox of 24 by 24 drawn onto 48 by 48 is twice as big; onto 100 by 48
	// it is as big as it can be without going outside, and the space left over
	// is shared out around it.
	m := fitTransform([4]float64{0, 0, 24, 24}, 48, 48)
	if p, q := m.Map(0, 0); !nearFloat(p, 0, 0.001) || !nearFloat(q, 0, 0.001) {
		t.Errorf("the top left went to %v,%v, want it at the top left", p, q)
	}
	if p, q := m.Map(24, 24); !nearFloat(p, 48, 0.001) || !nearFloat(q, 48, 0.001) {
		t.Errorf("the far corner went to %v,%v, want it at 48,48", p, q)
	}
	// Wide canvas: the drawing keeps its own shape and sits in the middle.
	m = fitTransform([4]float64{0, 0, 24, 24}, 100, 48)
	if _, q := m.Map(0, 0); !nearFloat(q, 0, 0.001) {
		t.Errorf("the top went to y=%v, want it to stay at the top", q)
	}
	// A viewBox that does not start at the origin is shifted, so both this and
	// one that does fill the same box the same way.
	one := fitTransform([4]float64{0, 0, 24, 24}, 24, 24)
	two := fitTransform([4]float64{-12, -12, 24, 24}, 24, 24)
	for _, pt := range [][2]float64{{0, 0}, {12, 12}, {24, 24}} {
		// The same point in the two drawings, which sit twelve apart because one
		// of them starts its count at the origin and the other twelve before it.
		ax, ay := one.Map(pt[0], pt[1])
		bx, by := two.Map(pt[0]-12, pt[1]-12)
		if !nearFloat(ax, bx, 0.001) || !nearFloat(ay, by, 0.001) {
			t.Errorf("a point at %v lands at %v,%v from one viewBox and %v,%v from another",
				pt, ax, ay, bx, by)
		}
	}
	// A viewBox of no size says nothing about what to fit, so the drawing is
	// stretched across the canvas rather than left out of it.
	m = fitTransform([4]float64{0, 0, 0, 0}, 24, 12)
	if p, q := m.Map(1, 1); !nearFloat(p, 24, 0.001) || !nearFloat(q, 12, 0.001) {
		t.Errorf("(1,1) went to %v,%v, want it at the far corner", p, q)
	}
}

// mapX is where a transform puts a point along the line of nothing, which is
// the number most of these checks are about.
func mapX(m canvas.Matrix, x, y float64) float64 {
	p, _ := m.Map(x, y)
	return p
}

// nearMatrix says whether two transforms are the same one, to the small.
func nearMatrix(a, b canvas.Matrix) bool {
	return nearFloat(a.A, b.A, 0.001) && nearFloat(a.B, b.B, 0.001) &&
		nearFloat(a.C, b.C, 0.001) && nearFloat(a.D, b.D, 0.001) &&
		nearFloat(a.E, b.E, 0.001) && nearFloat(a.F, b.F, 0.001)
}
