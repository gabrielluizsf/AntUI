package css

import (
	"math"
	"testing"
)

// trule computes the style for one label rule, the same way an element's
// declaration would reach the cascade.
func trule(t *testing.T, rule string) Style {
	t.Helper()
	sh, err := Parse("label {" + rule + "}")
	if err != nil {
		t.Fatal(err)
	}
	return sh.Style("label", nil, StateNone, 800)
}

func TestTransformListParses(t *testing.T) {
	st := trule(t, "transform: translate(10px 20%) rotate(90deg) scale(2);")
	if len(st.Transform) != 3 {
		t.Fatalf("list length = %d, want 3", len(st.Transform))
	}
	tr, ro, sc := st.Transform[0], st.Transform[1], st.Transform[2]
	if tr.Kind != TransformTranslate {
		t.Errorf("first kind = %d, want translate", tr.Kind)
	}
	if got := tr.Dx.Resolve(Units{Width: 400}); got != 10 {
		t.Errorf("translate x = %v, want 10 reference pixels", got)
	}
	if !tr.Dy.IsPct() {
		t.Errorf("translate y should stay a percentage, got %v", tr.Dy)
	}
	if ro.Kind != TransformRotate || math.Abs(ro.Ax.Deg()-90) > 1e-9 {
		t.Errorf("second = kind %d angle %vdeg, want rotate 90deg", ro.Kind, ro.Ax.Deg())
	}
	if sc.Kind != TransformScale || sc.Sx != 2 || sc.Sy != 2 {
		t.Errorf("third = kind %d scale %v/%v, want 2/2", sc.Kind, sc.Sx, sc.Sy)
	}
}

func TestTransformSingleArgumentExpands(t *testing.T) {
	st := trule(t, "transform: scaleX(0.5) translate(15px);")
	if len(st.Transform) != 2 {
		t.Fatalf("list length = %d, want 2", len(st.Transform))
	}
	sx, tr := st.Transform[0], st.Transform[1]
	if sx.Sx != 0.5 || sx.Sy != 1 {
		t.Errorf("scaleX = %v/%v, want 0.5/1", sx.Sx, sx.Sy)
	}
	if tr.Dy.Resolve(Units{}) != 0 {
		t.Errorf("translate without y = %v, want 0", tr.Dy.Resolve(Units{}))
	}
}

func TestTransformSkewAndMatrix(t *testing.T) {
	st := trule(t, "transform: matrix(1, 2, 3, 4, 5, 6) skew(30deg 45deg);")
	if len(st.Transform) != 2 {
		t.Fatalf("list length = %d, want 2", len(st.Transform))
	}
	m := st.Transform[0]
	if m.Kind != TransformMatrix {
		t.Fatalf("first kind = %d, want matrix", m.Kind)
	}
	for i, want := range []float64{1, 2, 3, 4, 5, 6} {
		if m.M[i] != want {
			t.Errorf("matrix[%d] = %v, want %v", i, m.M[i], want)
		}
	}
	k := st.Transform[1]
	if k.Kind != TransformSkew || math.Abs(k.Ax.Deg()-30) > 1e-9 || math.Abs(k.Ay.Deg()-45) > 1e-9 {
		t.Errorf("second = kind %d angles %v/%v, want skew 30deg 45deg", k.Kind, k.Ax.Deg(), k.Ay.Deg())
	}
}

func TestTransformNoneAndGarbageDrop(t *testing.T) {
	if st := trule(t, "transform: none;"); st.Transform != nil {
		t.Errorf("none must clear the list, got %v", st.Transform)
	}
	// An unknown function is a bad value: the property falls to its initial
	// value, exactly like an unknown property would.
	if st := trule(t, "transform: translate(10px) spin(2);"); st.Transform != nil {
		t.Errorf("an unknown function must drop the whole value, got %v", st.Transform)
	}
	if st := trule(t, "transform: rotate(90deg 45deg);"); st.Transform != nil {
		t.Errorf("too many rotate arguments must drop the value, got %v", st.Transform)
	}
	if st := trule(t, "transform: scale(2,);"); st.Transform != nil {
		t.Errorf("a stray comma must drop the value, got %v", st.Transform)
	}
}

func TestTransformScalePercent(t *testing.T) {
	st := trule(t, "transform: scale(150%);")
	if len(st.Transform) != 1 || st.Transform[0].Sx != 1.5 {
		t.Errorf("scale(150%%) = %v, want 1.5", st.Transform)
	}
}

func TestTransformOriginKeywords(t *testing.T) {
	for _, tc := range []struct {
		rule string
		x, y int
	}{
		{"transform-origin: top left;", 0, 0},
		{"transform-origin: bottom right;", 100, 100},
		{"transform-origin: right;", 100, 50},
		{"transform-origin: top;", 50, 0},
		{"transform-origin: 20% 30%;", 20, 30},
		{"transform-origin: center;", 50, 50},
	} {
		st := trule(t, "transform: rotate(1deg); "+tc.rule)
		gx := st.TransformOrigin[0].Resolve(Units{Width: 100})
		gy := st.TransformOrigin[1].Resolve(Units{Width: 100})
		if gx != tc.x || gy != tc.y {
			t.Errorf("%q = (%v,%v), want (%v,%v)", tc.rule, gx, gy, tc.x, tc.y)
		}
	}
}

func TestTransformOriginTwoValuesSwapAxes(t *testing.T) {
	// "top left" and "left top" are the same pivot; a vertical keyword in
	// front is the axis reading, not the order.
	st := trule(t, "transform: rotate(1deg); transform-origin: top left;")
	left := st.TransformOrigin[0].Resolve(Units{Width: 100})
	top := st.TransformOrigin[1].Resolve(Units{Width: 100})
	if left != 0 || top != 0 {
		t.Errorf("top left = (%v,%v), want (0,0)", left, top)
	}
}

func TestTransformOriginThirdValueDropped(t *testing.T) {
	st := trule(t, "transform: translate(1px); transform-origin: top left 12px;")
	x := st.TransformOrigin[0].Resolve(Units{Width: 100})
	y := st.TransformOrigin[1].Resolve(Units{Width: 100})
	if x != 0 || y != 0 {
		t.Errorf("top left with a z = (%v,%v), want (0,0)", x, y)
	}
}

func TestTransformOriginLengths(t *testing.T) {
	st := trule(t, "transform: translate(1px); transform-origin: 10px 50%;")
	gx := st.TransformOrigin[0].Resolve(Units{})
	if gx != 10 {
		t.Errorf("x length = %v, want 10 reference pixels", gx)
	}
	gy := st.TransformOrigin[1].Resolve(Units{Width: 200})
	if gy != 100 {
		t.Errorf("y percentage = %v, want 100 (50%% of 200)", gy)
	}
}

func TestTransformOriginUnsetStaysInitial(t *testing.T) {
	st := trule(t, "color: red;")
	// An untouched origin keeps its initial value, which the painter reads as
	// the CSS default of the centre: the computed style never stores it.
	if st.Has("transform-origin") {
		t.Error("an undeclared origin must not be marked set")
	}
	if st.TransformOrigin[0] != (Length{}) || st.TransformOrigin[1] != (Length{}) {
		t.Errorf("an unset origin = %v, want the untouched zero store", st.TransformOrigin)
	}
}

func TestTransformOriginBadValuesRejected(t *testing.T) {
	for _, rule := range []string{
		"transform-origin: top top;",
		"transform-origin: 1 2 3 4;",
		"transform-origin: center auto;",
		"transform-origin: sideways;",
		"transform-origin: left, top;",
	} {
		if st := trule(t, rule); st.Has("transform-origin") {
			t.Errorf("%q must drop the property, but it was stored", rule)
		}
	}
}
