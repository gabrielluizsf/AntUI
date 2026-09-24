package css

import (
	"math"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

func TestBackgroundImageParses(t *testing.T) {
	sh := mustParse(t, `button { background-image: url(logo.png), url("face.png"); }`)
	st := sh.Style("button", nil, StateNone, 800)
	if len(st.BackgroundImages) != 2 {
		t.Fatalf("loaded %d images, want 2", len(st.BackgroundImages))
	}
	if st.BackgroundImages[0].URL != "logo.png" {
		t.Errorf("first url = %q, want logo.png", st.BackgroundImages[0].URL)
	}
	if st.BackgroundImages[1].URL != "face.png" {
		t.Errorf("second url = %q, want face.png (quotes stripped)", st.BackgroundImages[1].URL)
	}
	if st.BackgroundImages[0].Grad != nil {
		t.Error("a url() layer must carry no gradient")
	}
}

func TestBackgroundImageNone(t *testing.T) {
	sh := mustParse(t, `button { background-image: none; }`)
	st := sh.Style("button", nil, StateNone, 800)
	if len(st.BackgroundImages) != 0 {
		t.Errorf("none must clear the layers, got %v", st.BackgroundImages)
	}
}

func TestBackgroundLinearGradientParses(t *testing.T) {
	sh := mustParse(t, `button { background-image: linear-gradient(45deg, red, blue); }`)
	st := sh.Style("button", nil, StateNone, 800)
	if len(st.BackgroundImages) != 1 || st.BackgroundImages[0].Grad == nil {
		t.Fatalf("expected one gradient layer, got %+v", st.BackgroundImages)
	}
	g := st.BackgroundImages[0].Grad
	if g.Kind != GradientLinear {
		t.Errorf("kind = %d, want GradientLinear", g.Kind)
	}
	if d := math.Abs(g.Angle - math.Pi/4); d > 1e-6 {
		t.Errorf("angle = %v, want 45 degrees", g.Angle)
	}
	if len(g.Stops) != 2 {
		t.Fatalf("stops = %d, want 2", len(g.Stops))
	}
	if g.Stops[0].Color != red || g.Stops[0].Offset != -1 {
		t.Errorf("first stop = %+v, want CSS red with no offset", g.Stops[0])
	}
}

func TestBackgroundLinearGradientDefaultIsBottom(t *testing.T) {
	sh := mustParse(t, `button { background-image: linear-gradient(red, yellow, blue); }`)
	st := sh.Style("button", nil, StateNone, 800)
	g := st.BackgroundImages[0].Grad
	if g.Kind != GradientLinear {
		t.Fatalf("kind = %d", g.Kind)
	}
	if d := math.Abs(g.Angle - math.Pi); d > 1e-6 {
		t.Errorf("missing direction must be to bottom (pi), got %v", g.Angle)
	}
	if len(g.Stops) != 3 {
		t.Errorf("stops = %d, want 3", len(g.Stops))
	}
}

func TestBackgroundRadialGradientParses(t *testing.T) {
	sh := mustParse(t, `button { background-image: radial-gradient(closest-side circle at 20% 30px, red, yellow 50%, blue); }`)
	st := sh.Style("button", nil, StateNone, 800)
	g := st.BackgroundImages[0].Grad
	if g.Kind != GradientRadial {
		t.Fatalf("kind = %d, want radial", g.Kind)
	}
	if g.Shape != GradientCircle || g.Size != GradientClosestSide {
		t.Errorf("shape/size = %d/%d, want circle closest-side", g.Shape, g.Size)
	}
	if g.Center[0].Edge != BackPosOffset || !g.Center[0].Off.IsPct() {
		t.Errorf("centre x = %+v, want a 20%% offset", g.Center[0])
	}
	if g.Center[1].Edge != BackPosOffset || g.Center[1].Off.Px(0) != 30 {
		t.Errorf("centre y = %+v, want 30px", g.Center[1])
	}
	if len(g.Stops) != 3 || g.Stops[1].Offset != 0.5 {
		t.Fatalf("stops = %+v, want three with 50%% at the middle", g.Stops)
	}
}

func TestBackgroundConicGradientParses(t *testing.T) {
	sh := mustParse(t, `button { background-image: conic-gradient(from 90deg at 25% 75%, red, blue); }`)
	st := sh.Style("button", nil, StateNone, 800)
	g := st.BackgroundImages[0].Grad
	if g.Kind != GradientConic {
		t.Fatalf("kind = %d, want conic", g.Kind)
	}
	if d := math.Abs(g.Angle - math.Pi/2); d > 1e-6 {
		t.Errorf("angle = %v, want 90 degrees", g.Angle)
	}
	if !g.Center[0].Off.IsPct() {
		t.Errorf("centre x = %+v, want 25%%", g.Center[0])
	}
}

func TestBackgroundGradientPrematureEnd(t *testing.T) {
	// "closest-side circle" naming a size with no prelude comma folds into the
	// first stop, so the layer still parses.
	sh := mustParse(t, `button { background-image: radial-gradient(circle red, blue); }`)
	st := sh.Style("button", nil, StateNone, 800)
	g := st.BackgroundImages[0].Grad
	if g == nil {
		t.Fatal("radial gradient should parse")
	}
	if len(g.Stops) != 2 {
		t.Errorf("stops = %d, want 2", len(g.Stops))
	}
}

func TestBackgroundCurrentColorResolves(t *testing.T) {
	sh := mustParse(t, `button { color: #123456; background-image: linear-gradient(red, currentColor); }`)
	st := sh.Style("button", nil, StateNone, 800)
	g := st.BackgroundImages[0].Grad
	if g.Stops[1].Color != 0xFF123456 {
		t.Errorf("currentColor stop = %v, want the button colour", g.Stops[1].Color)
	}
}

func TestBackgroundPositionParses(t *testing.T) {
	sh := mustParse(t, `button { background-position: right 20px top 30px; }`)
	st := sh.Style("button", nil, StateNone, 800)
	if len(st.BackgroundPos) != 1 {
		t.Fatalf("positions = %d, want 1", len(st.BackgroundPos))
	}
	x, y := st.BackgroundPos[0][0], st.BackgroundPos[0][1]
	if x.Edge != BackPosEnd || x.Off.Px(0) != 20 {
		t.Errorf("x = %+v, want right 20px", x)
	}
	if y.Edge != BackPosStart || y.Off.Px(0) != 30 {
		t.Errorf("y = %+v, want top 30px", y)
	}
}

func TestBackgroundPositionKeywords(t *testing.T) {
	for _, tc := range []struct {
		raw  string
		x, y BackPosition
	}{
		{"center", BackPosition{Edge: BackPosMiddle}, BackPosition{Edge: BackPosMiddle}},
		{"left center", BackPosition{Edge: BackPosStart}, BackPosition{Edge: BackPosMiddle}},
		{"bottom right", BackPosition{Edge: BackPosEnd}, BackPosition{Edge: BackPosEnd}},
		{"10% 20px", BackPosition{Edge: BackPosOffset, Off: Pct(10)}, BackPosition{Edge: BackPosOffset, Off: Fixed(20)}},
		{"left 20px", BackPosition{Edge: BackPosStart, Off: Fixed(20)}, BackPosition{Edge: BackPosMiddle}},
	} {
		sh := mustParse(t, `button { background-position: `+tc.raw+`; }`)
		st := sh.Style("button", nil, StateNone, 800)
		p := st.BackgroundPos[0]
		if p[0] != tc.x || p[1] != tc.y {
			t.Errorf("background-position: %s = (%+v, %+v), want (%+v, %+v)",
				tc.raw, p[0], p[1], tc.x, tc.y)
		}
	}
	if !testing.Short() {
		// "left right" overfills an axis, but the corner fallback still lands on
		// both axes.
		sh := mustParse(t, `button { background-position: left right; }`)
		_ = sh.Style("button", nil, StateNone, 800)
	}
}

func TestBackgroundPositionMissingAxisCentres(t *testing.T) {
	sh := mustParse(t, `button { background-position: 20px; }`)
	st := sh.Style("button", nil, StateNone, 800)
	p := st.BackgroundPos[0]
	if p[0].Edge != BackPosOffset || p[0].Off.Px(0) != 20 {
		t.Errorf("x = %+v, want a 20px hand-off", p[0])
	}
	if p[1].Edge != BackPosMiddle {
		t.Errorf("y = %+v, want centred", p[1])
	}
}

func TestBackgroundSizeParses(t *testing.T) {
	for _, tc := range []struct {
		raw string
		sz  BackSize
	}{
		{"cover", BackSize{Cover: true}},
		{"contain", BackSize{Contain: true}},
		{"100px 50px", BackSize{W: Fixed(100), H: Fixed(50)}},
		{"50% auto", BackSize{W: Pct(50), H: Auto()}},
		{"80px", BackSize{W: Fixed(80), H: Auto()}},
	} {
		sh := mustParse(t, `button { background-size: `+tc.raw+`; }`)
		st := sh.Style("button", nil, StateNone, 800)
		if len(st.BackgroundSize) != 1 || st.BackgroundSize[0] != tc.sz {
			t.Errorf("background-size: %s = %+v, want %+v", tc.raw, st.BackgroundSize, tc.sz)
		}
	}
}

func TestBackgroundRepeatParses(t *testing.T) {
	for _, tc := range []struct {
		raw string
		rep BackRepeat
	}{
		{"repeat", BackRepeat{BackRepeatRepeat, BackRepeatRepeat}},
		{"no-repeat", BackRepeat{BackRepeatNoRepeat, BackRepeatNoRepeat}},
		{"repeat no-repeat", BackRepeat{BackRepeatRepeat, BackRepeatNoRepeat}},
		{"repeat-x", BackRepeat{BackRepeatRepeat, BackRepeatNoRepeat}},
		{"repeat-y", BackRepeat{BackRepeatNoRepeat, BackRepeatRepeat}},
		{"space round", BackRepeat{BackRepeatSpace, BackRepeatRound}},
	} {
		sh := mustParse(t, `button { background-repeat: `+tc.raw+`; }`)
		st := sh.Style("button", nil, StateNone, 800)
		if len(st.BackgroundRepeat) != 1 || st.BackgroundRepeat[0] != tc.rep {
			t.Errorf("background-repeat: %s = %+v, want %+v", tc.raw, st.BackgroundRepeat, tc.rep)
		}
	}
}

func TestBackgroundClipOriginAttachment(t *testing.T) {
	sh := mustParse(t, `
		button {
			background-clip: padding-box, border-box;
			background-origin: content-box;
			background-attachment: fixed;
		}
	`)
	st := sh.Style("button", nil, StateNone, 800)
	if len(st.BackgroundClip) != 2 || st.BackgroundClip[0] != BackPadding || st.BackgroundClip[1] != BackBorder {
		t.Errorf("clips = %v, want padding then border", st.BackgroundClip)
	}
	if len(st.BackgroundOrigin) != 1 || st.BackgroundOrigin[0] != BackContent {
		t.Errorf("origins = %v, want content", st.BackgroundOrigin)
	}
	if len(st.BackgroundAttach) != 1 || st.BackgroundAttach[0] != BackAttachFixed {
		t.Errorf("attachments = %v, want fixed", st.BackgroundAttach)
	}
}

func TestBackgroundShorthand(t *testing.T) {
	sh := mustParse(t, `button { background: url(logo.png) no-repeat right / cover, linear-gradient(red, blue) 50% padding-box / 100px 20px; }`)
	st := sh.Style("button", nil, StateNone, 800)
	if len(st.BackgroundImages) != 2 {
		t.Fatalf("layers = %d, want 2", len(st.BackgroundImages))
	}
	if st.BackgroundImages[0].URL != "logo.png" {
		t.Errorf("first layer = %q", st.BackgroundImages[0].URL)
	}
	if st.BackgroundImages[0].Grad != nil || st.BackgroundImages[1].Grad == nil {
		t.Error("layers mis-assigned")
	}
	if len(st.BackgroundPos) != 2 || st.BackgroundPos[0][0].Edge != BackPosEnd {
		t.Errorf("positions = %+v, want the first right-anchored", st.BackgroundPos)
	}
	if len(st.BackgroundSize) != 2 || !st.BackgroundSize[0].Cover || st.BackgroundSize[1].W.Px(0) != 100 {
		t.Errorf("sizes = %+v, want cover and 100px 20px", st.BackgroundSize)
	}
	if len(st.BackgroundOrigin) != 1 || st.BackgroundOrigin[0] != BackPadding {
		t.Errorf("origins = %+v, want the second layer's padding-box", st.BackgroundOrigin)
	}
	if len(st.BackgroundRepeat) != 1 || st.BackgroundRepeat[0][0] != BackRepeatNoRepeat {
		t.Errorf("repeats = %+v, want the first layer no-repeat", st.BackgroundRepeat)
	}
	if !st.Has("background-image") || !st.Has("background-position") {
		t.Error("shorthand must mark its longhands as set")
	}
}

func TestBackgroundShorthandPicksTheColour(t *testing.T) {
	sh := mustParse(t, `button { background: url(logo.png), linear-gradient(red, blue), #123456; }`)
	st := sh.Style("button", nil, StateNone, 800)
	if st.Background != 0xFF123456 {
		t.Errorf("colour = %v, want #123456", st.Background)
	}
	if len(st.BackgroundImages) != 2 {
		t.Errorf("layers = %d, want the two images only", len(st.BackgroundImages))
	}
}

func TestBackgroundInitial(t *testing.T) {
	sh := mustParse(t, `
		button {
			background-color: red;
			background-image: url(a.png);
			background-position: 10px 20px;
			background: initial;
		}
	`)
	st := sh.Style("button", nil, StateNone, 800)
	if st.Background.A() != 0 {
		t.Errorf("background = %v after initial, want transparent", st.Background)
	}
	if len(st.BackgroundImages) != 0 || len(st.BackgroundPos) != 0 {
		t.Errorf("layers must clear on initial, got images %d positions %d", len(st.BackgroundImages), len(st.BackgroundPos))
	}
}

func TestBackgroundMultiLayerCycle(t *testing.T) {
	sh := mustParse(t, `button { background: linear-gradient(red, blue), url(a.png), url(b.png); background-size: 50% auto; }`)
	st := sh.Style("button", nil, StateNone, 800)
	if len(st.BackgroundImages) != 3 {
		t.Fatalf("layers = %d, want 3", len(st.BackgroundImages))
	}
	if len(st.BackgroundSize) != 1 {
		t.Fatalf("sizes = %d, want the single value", len(st.BackgroundSize))
	}
	// One declared value cycles across every layer when painted.
	for i := 0; i < 3; i++ {
		sz := st.BackgroundSize[i%len(st.BackgroundSize)]
		if !sz.W.IsPct() {
			t.Errorf("layer %d size = %+v, want the cycled 50%%", i, sz)
		}
	}
}

func TestBackgroundRejectsGarbage(t *testing.T) {
	for _, raw := range []string{
		`background-image: url(a.png), bogus;`,
		`background-position: left right top bottom top;`,
		`background-size: cover 10px;`,
		`background-repeat: repeat no-repeat space;`,
		`background-clip: nonsense;`,
		`background-attachment: sideways;`,
	} {
		sh := mustParse(t, `button { `+raw+` }`)
		st := sh.Style("button", nil, StateNone, 800)
		if len(st.BackgroundImages) > 0 || len(st.BackgroundPos) > 0 || len(st.BackgroundSize) > 0 ||
			len(st.BackgroundRepeat) > 0 || len(st.BackgroundClip) > 0 || len(st.BackgroundAttach) > 0 {
			t.Errorf("%s must be dropped whole", raw)
		}
	}
}

var red = canvas.RGBA(255, 0, 0, 255)
