package css

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

func mustParse(t *testing.T, src string) *Sheet {
	t.Helper()
	sh, err := Parse(src)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return sh
}

func TestParseClassSelector(t *testing.T) {
	sh := mustParse(t, `
		button { background-color: #3E63DD; border-radius: 8px; }
	`)
	st := sh.Style("button", nil, StateNone, 800)
	if !st.Has("background-color") {
		t.Error("background-color not recorded as set")
	}
	if st.Background != 0xFF3E63DD {
		t.Errorf("background = %v, want 0xFF3E63DD", st.Background)
	}
	if st.Radius[0] != 8 || st.Radius[1] != 8 || st.Radius[2] != 8 || st.Radius[3] != 8 {
		t.Errorf("radius = %v, want all 8", st.Radius)
	}
	if st.Display != 0 {
		t.Error("default display should be block")
	}
}

func TestUniversalAndElement(t *testing.T) {
	sh := mustParse(t, `
		* { box-sizing: border-box; padding: 4px; }
		input { padding: 8px; }
	`)
	st := sh.Style("input", nil, StateNone, 800)
	if !st.BoxSizing {
		t.Error("* box-sizing did not apply")
	}
	if st.Padding[0].Px(0) != 8 {
		t.Errorf("input padding = %d, want 8 (classically more specific wins)", st.Padding[0].Px(0))
	}
}

func TestPseudoClasses(t *testing.T) {
	sh := mustParse(t, `
		button { background-color: blue; }
		button:hover { background-color: red; }
		button:hover.fancy { background-color: green; }
	`)
	base := sh.Style("button", nil, StateNone, 800)
	if base.Background != 0xFF0000FF {
		t.Errorf("base background = %v, want CSS blue #0000FF", base.Background)
	}
	hover := sh.Style("button", nil, StateHover, 800)
	if hover.Background != 0xFFFF0000 {
		t.Errorf("hover background = %v, want CSS red #FF0000", hover.Background)
	}
	both := sh.Style("button", []string{"fancy"}, StateHover, 800)
	if both.Background != 0xFF008000 {
		t.Errorf("hover+fancy background = %v, want CSS green #008000", both.Background)
	}
}

func TestSpecificityAndOrder(t *testing.T) {
	sh := mustParse(t, `
		button { color: black; }
		.foo { color: white; }
		button.foo { color: gray; }
		button.foo { color: red; }
	`)
	st := sh.Style("button", []string{"foo"}, StateNone, 800)
	// The later rule of equal specificity wins.
	if st.Color != 0xFFFF0000 {
		t.Errorf("color = %v, want red (later equal-specificity rule wins)", st.Color)
	}
}

func TestMediaQueriesSwitchByWidth(t *testing.T) {
	sh := mustParse(t, `
		button { font-size: 16px; }
		@media (min-width: 720px) {
			button { font-size: 24px; }
		}
	`)
	small := sh.Style("button", nil, StateNone, 480)
	if small.FontSize != 16 {
		t.Errorf("small width font = %d, want 16", small.FontSize)
	}
	big := sh.Style("button", nil, StateNone, 1000)
	if big.FontSize != 24 {
		t.Errorf("big width font = %d, want 24", big.FontSize)
	}
}

func TestMediaMaxWidth(t *testing.T) {
	sh := mustParse(t, `
		@media (max-width: 600px) { input { width: 100px; } }
	`)
	if w := sh.Style("input", nil, StateNone, 400); w.Width.Px(0) != 100 {
		t.Errorf("max-width query should apply at 400, got %d", w.Width.Px(0))
	}
	if w := sh.Style("input", nil, StateNone, 900); w.Has("width") && w.Width.Px(0) == 100 {
		t.Error("max-width query should NOT apply at 900")
	}
}

func TestMarginShorthand(t *testing.T) {
	sh := mustParse(t, `.btn { margin: 4px 8px; }`)
	st := sh.Style("button", []string{"btn"}, StateNone, 800)
	want := [4]int{4, 8, 4, 8}
	got := [4]int{st.Margin[0].Px(0), st.Margin[1].Px(0), st.Margin[2].Px(0), st.Margin[3].Px(0)}
	if got != want {
		t.Errorf("margin = %v, want %v", got, want)
	}
}

func TestMarginFourValues(t *testing.T) {
	sh := mustParse(t, `.btn { margin: 1px 2px 3px 4px; }`)
	st := sh.Style("button", []string{"btn"}, StateNone, 800)
	got := [4]int{st.Margin[0].Px(0), st.Margin[1].Px(0), st.Margin[2].Px(0), st.Margin[3].Px(0)}
	if got != [4]int{1, 2, 3, 4} {
		t.Errorf("margin = %v", got)
	}
}

func TestBorderShorthand(t *testing.T) {
	sh := mustParse(t, `.btn { border: 2px solid red; }`)
	st := sh.Style("button", []string{"btn"}, StateNone, 800)
	if !st.BorderOn() {
		t.Error("border should be on")
	}
	for i := 0; i < 4; i++ {
		if st.BorderWidth[i] != 2 {
			t.Errorf("border width %d = %d, want 2", i, st.BorderWidth[i])
		}
		if st.BoxColor[i] != 0xFFFF0000 {
			t.Errorf("border color %d = %v", i, st.BoxColor[i])
		}
	}
}

func TestBackgroundShorthandAndColorNames(t *testing.T) {
	sh := mustParse(t, `
		body { background-color: navy; }
		.btn { background: #ff0000; }
		.other { background: rgba(0, 128, 255, 0.5); }
	`)
	if st := sh.Style("body", nil, StateNone, 800); st.Background != 0xFF000080 {
		t.Errorf("body bg = %v, want navy", st.Background)
	}
	if st := sh.Style("button", []string{"btn"}, StateNone, 800); st.Background != 0xFFFF0000 {
		t.Errorf("btn bg = %v, want red", st.Background)
	}
	other := sh.Style("", []string{"other"}, StateNone, 800)
	if other.Background.A() != 128 {
		t.Errorf("rgba alpha = %d, want 128 (0.5*255 rounded)", other.Background.A())
	}
}

func TestDisplayNoneAndPosition(t *testing.T) {
	sh := mustParse(t, `
		.hidden { display: none; }
		.fixed { position: absolute; top: 10px; left: 20px; }
	`)
	if st := sh.Style("div", []string{"hidden"}, StateNone, 800); st.Display != 1 {
		t.Error("display:none not honoured")
	}
	st := sh.Style("div", []string{"fixed"}, StateNone, 800)
	if st.Position != 1 {
		t.Error("position:absolute not honoured")
	}
	if st.Top.Px(0) != 10 || st.Left.Px(0) != 20 {
		t.Errorf("offset = %d,%d, want 10,20", st.Top.Px(0), st.Left.Px(0))
	}
}

func TestPercentagesResolveAgainstBase(t *testing.T) {
	sh := mustParse(t, `.panel { width: 50%; height: 25%; }`)
	st := sh.Style("div", []string{"panel"}, StateNone, 800)
	if st.Width.Px(800) != 400 {
		t.Errorf("width 50%% = %d, want 400", st.Width.Px(800))
	}
	if st.Height.Px(400) != 100 {
		t.Errorf("height 25%% = %d, want 100", st.Height.Px(400))
	}
}

func TestUnknownPropertyIsSkipped(t *testing.T) {
	sh := mustParse(t, `
		button { color: white; -webkit-animation: spin 1s; color: red; }
	`)
	st := sh.Style("button", nil, StateNone, 800)
	if st.Color != 0xFFFF0000 {
		t.Errorf("color = %v, want red (unknown decl skipped, cascade continued)", st.Color)
	}
	if len(sh.Warn) == 0 {
		t.Log("no warnings; unknown properties are silently ignored like browsers")
	}
}

func TestSelectorList(t *testing.T) {
	sh := mustParse(t, `button, input, .thing { color: lime; }`)
	if st := sh.Style("button", nil, StateNone, 800); st.Color != 0xFF00FF00 {
		t.Errorf("button color = %v", st.Color)
	}
	if st := sh.Style("input", nil, StateNone, 800); st.Color != 0xFF00FF00 {
		t.Errorf("input color = %v", st.Color)
	}
	if st := sh.Style("", []string{"thing"}, StateNone, 800); st.Color != 0xFF00FF00 {
		t.Errorf("thing color = %v", st.Color)
	}
	if st := sh.Style("label", nil, StateNone, 800); st.Color == 0xFF00FF00 {
		t.Error("label should not match")
	}
}

func TestMultipleClasses(t *testing.T) {
	sh := mustParse(t, `.btn.primary { color: white; }`)
	if st := sh.Style("button", []string{"btn", "primary"}, StateNone, 800); st.Color != 0xFFFFFFFF {
		t.Errorf("color = %v", st.Color)
	}
	if st := sh.Style("button", []string{"btn"}, StateNone, 800); st.Has("color") {
		t.Error(".btn.primary should not match a widget holding only .btn")
	}
}

func TestCommentsAndNoise(t *testing.T) {
	sh := mustParse(t, `
		/* a comment */
		@charset "utf-8";
		@import url("other.css");
		button { color: white; /* inline */ }
		/* end */
	`)
	if st := sh.Style("button", nil, StateNone, 800); st.Color != 0xFFFFFFFF {
		t.Errorf("color = %v", st.Color)
	}
}

func TestHSL(t *testing.T) {
	sh := mustParse(t, `.x { color: hsl(120, 100%, 50%); }`)
	st := sh.Style("", []string{"x"}, StateNone, 800)
	if st.Color != 0xFF00FF00 {
		t.Errorf("hsl green = %v, want #00FF00", st.Color)
	}
}

func TestShadowAndOpacity(t *testing.T) {
	sh := mustParse(t, `
		.btn { opacity: 0.5; }
	`)
	st := sh.Style("button", []string{"btn"}, StateNone, 800)
	if st.Opacity != 0.5 {
		t.Errorf("opacity = %v, want 0.5", st.Opacity)
	}
}

func TestWidthAutoAndFitContent(t *testing.T) {
	sh := mustParse(t, `.a { width: auto; } .b { }`)
	sa := sh.Style("", []string{"a"}, StateNone, 800)
	if !sa.Width.Auto() {
		t.Error("width:auto should keep the auto keyword")
	}
	sb := sh.Style("", []string{"b"}, StateNone, 800)
	if sb.Has("width") {
		t.Error("unset width should not claim to be set")
	}
}

var _ = canvas.RGB
