package css

import (
	"math"
	"testing"
)

// flexRule computes the style one rule applies to its own role, matching how
// a declared element reaches the cascade in the template.
func flexRule(t *testing.T, role, rule string) Style {
	t.Helper()
	sh, err := Parse(role + " {" + rule + "}")
	if err != nil {
		t.Fatal(err)
	}
	return sh.Style(role, nil, StateNone, 800)
}

func TestFlexDirectionParsesKeywords(t *testing.T) {
	for raw, want := range map[string]uint8{
		"row":            FlexDirectionRow,
		"row-reverse":    FlexDirectionRowReverse,
		"column":         FlexDirectionColumn,
		"column-reverse": FlexDirectionColumnReverse,
	} {
		if v, ok := parseFlexDirection(raw); !ok || v != want {
			t.Errorf("%q = %d,%v want %d,true", raw, v, ok, want)
		}
	}
	for _, raw := range []string{"", "diagonal", "column "} {
		if _, ok := parseFlexDirection(raw); ok {
			t.Errorf("%q must not parse", raw)
		}
	}
}

func TestFlexWrapParsesKeywords(t *testing.T) {
	for raw, want := range map[string]uint8{
		"nowrap":       FlexWrapNowrap,
		"wrap":         FlexWrapWrap,
		"wrap-reverse": FlexWrapWrapReverse,
	} {
		if v, ok := parseFlexWrap(raw); !ok || v != want {
			t.Errorf("%q = %d,%v want %d,true", raw, v, ok, want)
		}
	}
	if _, ok := parseFlexWrap("around"); ok {
		t.Error(`"around" must not parse`)
	}
}

func TestJustifyContentParsesKeywords(t *testing.T) {
	for raw, want := range map[string]uint8{
		"flex-start":    JustifyFlexStart,
		"start":         JustifyFlexStart,
		"flex-end":      JustifyFlexEnd,
		"end":           JustifyFlexEnd,
		"center":        JustifyCenter,
		"space-between": JustifySpaceBetween,
		"space-around":  JustifySpaceAround,
		"space-evenly":  JustifySpaceEvenly,
	} {
		if v, ok := parseJustifyContent(raw); !ok || v != want {
			t.Errorf("%q = %d,%v want %d,true", raw, v, ok, want)
		}
	}
}

func TestAlignItemsAndSelfParse(t *testing.T) {
	for raw, want := range map[string]uint8{
		"normal":     AlignStretch,
		"stretch":    AlignStretch,
		"flex-start": AlignFlexStart,
		"start":      AlignFlexStart,
		"flex-end":   AlignFlexEnd,
		"center":     AlignCenter,
		"baseline":   AlignBaseline,
	} {
		if v, ok := parseAlignItems(raw); !ok || v != want {
			t.Errorf("align-items %q = %d,%v want %d,true", raw, v, ok, want)
		}
	}
	if v, ok := parseAlignSelf("auto"); !ok || v != AlignAuto {
		t.Errorf("align-self auto = %d,%v want %d,true", v, ok, AlignAuto)
	}
	if _, ok := parseAlignItems("auto"); ok {
		t.Error("align-items must reject auto")
	}
	if _, ok := parseAlignSelf("middle"); ok {
		t.Error("align-self must reject middle")
	}
}

func TestAlignContentParsesKeywords(t *testing.T) {
	for raw, want := range map[string]uint8{
		"normal":        ContentStretch,
		"stretch":       ContentStretch,
		"flex-start":    ContentFlexStart,
		"flex-end":      ContentFlexEnd,
		"center":        ContentCenter,
		"space-between": ContentSpaceBetween,
		"space-around":  ContentSpaceAround,
		"space-evenly":  ContentSpaceEvenly,
	} {
		if v, ok := parseAlignContent(raw); !ok || v != want {
			t.Errorf("%q = %d,%v want %d,true", raw, v, ok, want)
		}
	}
}

func TestGapParsesOneAndTwoValues(t *testing.T) {
	row, col, ok := parseGap("10px", Units{})
	if !ok || row.Resolve(Units{}) != 10 || col.Resolve(Units{}) != 10 {
		t.Errorf(`10px row=%v col=%v ok=%v want 10/10`, row, col, ok)
	}
	row, col, ok = parseGap("8px 12px", Units{})
	if !ok || row.Resolve(Units{}) != 8 || col.Resolve(Units{}) != 12 {
		t.Errorf(`8px 12px row=%v col=%v ok=%v want 8/12`, row, col, ok)
	}
	for _, raw := range []string{"", "1px 2px 3px", "-5px", "auto", "none"} {
		if _, _, ok := parseGap(raw, Units{}); ok {
			t.Errorf("%q must not parse", raw)
		}
	}
}

func TestOrderParsesInteger(t *testing.T) {
	if v, ok := parseOrder("5"); !ok || v != 5 {
		t.Errorf(`"5" = %d,%v want 5,true`, v, ok)
	}
	if v, ok := parseOrder("-2"); !ok || v != -2 {
		t.Errorf(`"-2" = %d,%v want -2,true`, v, ok)
	}
	if v, ok := parseOrder("auto"); !ok || v != 0 {
		t.Errorf(`"auto" = %d,%v want 0,true`, v, ok)
	}
}

func TestFlexNumberRejectsNegatives(t *testing.T) {
	if v, ok := parseFlexNumber("2.5"); !ok || math.Abs(v-2.5) > 1e-9 {
		t.Errorf(`"2.5" = %v,%v want 2.5,true`, v, ok)
	}
	if v, ok := parseFlexNumber("0"); !ok || v != 0 {
		t.Errorf(`"0" = %v,%v want 0,true`, v, ok)
	}
	for _, raw := range []string{"-1", "big", ""} {
		if _, ok := parseFlexNumber(raw); ok {
			t.Errorf("%q must not parse", raw)
		}
	}
}

func TestFlexBasisParsesLengthAndKeywords(t *testing.T) {
	if b, ok := parseFlexBasis("auto", Units{}); !ok || !b.Auto() {
		t.Errorf(`"auto" = %v,%v want Auto,true`, b, ok)
	}
	if b, ok := parseFlexBasis("content", Units{}); !ok || !b.None() {
		t.Errorf(`"content" = %v,%v want none,true`, b, ok)
	}
	if b, ok := parseFlexBasis("100px", Units{}); !ok || b.Resolve(Units{}) != 100 {
		t.Errorf(`"100px" = %v,%v want 100,true`, b, ok)
	}
	if b, ok := parseFlexBasis("30%", Units{}); !ok || !b.IsPct() {
		t.Errorf(`"30%%" = %v,%v want pct,true`, b, ok)
	}
	if _, ok := parseFlexBasis("thick", Units{}); ok {
		t.Error(`"thick" must not parse`)
	}
}

func TestFlexShorthandForms(t *testing.T) {
	cases := []struct {
		raw    string
		grow   float64
		shrink float64
		basis  int // expected resolved basis; -1 expects the auto keyword,
		//       -2 expects a percentage
	}{
		{"none", 0, 0, -1},
		{"auto", 1, 1, -1},
		{"initial", 0, 1, -1},
		{"2", 2, 1, 0},
		{"1 0 100px", 1, 0, 100},
		{"0 1 30%", 0, 1, -2},
	}
	for _, tc := range cases {
		g, s, b, ok := parseFlex(tc.raw, Units{})
		if !ok || g != tc.grow || s != tc.shrink {
			t.Errorf("%q = %v,%v,%v ok=%v want grow=%v shrink=%v", tc.raw, g, s, b, ok, tc.grow, tc.shrink)
			continue
		}
		if tc.basis == -1 {
			if !b.Auto() {
				t.Errorf("%q basis = %v, want auto", tc.raw, b)
			}
			continue
		}
		got := b.Resolve(Units{})
		if tc.basis == -2 {
			if !b.IsPct() {
				t.Errorf("%q basis = %v, want a percentage", tc.raw, b)
			}
			continue
		}
		if got != tc.basis {
			t.Errorf("%q basis resolves to %d, want %d", tc.raw, got, tc.basis)
		}
	}
	g, s, b, ok := parseFlex("2 0 50px", Units{})
	if !ok || g != 2 || s != 0 || b.Resolve(Units{}) != 50 {
		t.Errorf(`"2 0 50px" = %v,%v,%v,%v want 2,0,50px`, g, s, b, ok)
	}
	g, s, b, ok = parseFlex("1 200px", Units{})
	if !ok || g != 1 || s != 1 || b.Resolve(Units{}) != 200 {
		t.Errorf(`"1 200px" = %v,%v,%v,%v want 1,1,200px`, g, s, b, ok)
	}
	for _, raw := range []string{"grow", "1 2 3 4px", "1 2 3 4"} {
		if _, _, _, ok := parseFlex(raw, Units{}); ok {
			t.Errorf("%q must not parse", raw)
		}
	}
}

func TestFlexPropertiesReachStyle(t *testing.T) {
	st := flexRule(t, "button", `
		flex-direction: column;
		flex-wrap: wrap;
		justify-content: space-between;
		align-items: center;
		align-content: space-around;
		gap: 8px 12px;
		order: 3;
		flex: 2 1 100px;
	`)
	if st.FlexDirection != FlexDirectionColumn {
		t.Errorf("direction = %d want column", st.FlexDirection)
	}
	if st.FlexWrap != FlexWrapWrap {
		t.Errorf("wrap = %d want wrap", st.FlexWrap)
	}
	if st.JustifyContent != JustifySpaceBetween {
		t.Errorf("justify = %d want space-between", st.JustifyContent)
	}
	if st.AlignItems != AlignCenter {
		t.Errorf("align-items = %d want center", st.AlignItems)
	}
	if st.AlignContent != ContentSpaceAround {
		t.Errorf("align-content = %d want space-around", st.AlignContent)
	}
	if st.RowGap.Resolve(Units{}) != 8 || st.ColumnGap.Resolve(Units{}) != 12 {
		t.Errorf("gap row=%v col=%v want 8/12", st.RowGap, st.ColumnGap)
	}
	if st.Order != 3 {
		t.Errorf("order = %d want 3", st.Order)
	}
	if st.FlexGrow != 2 || st.FlexShrink != 1 || st.FlexBasis.Resolve(Units{}) != 100 {
		t.Errorf("flex = %v %v %v want 2 1 100px", st.FlexGrow, st.FlexShrink, st.FlexBasis)
	}
}

func TestFlexInitials(t *testing.T) {
	st := flexRule(t, "button", `
		flex-direction: initial;
		flex-wrap: initial;
		justify-content: initial;
		align-items: initial;
		align-self: initial;
		align-content: initial;
		gap: initial;
		order: initial;
		flex: initial;
	`)
	if st.FlexDirection != FlexDirectionRow ||
		st.FlexWrap != FlexWrapNowrap ||
		st.JustifyContent != JustifyFlexStart ||
		st.AlignItems != AlignStretch ||
		st.AlignSelf != AlignAuto ||
		st.AlignContent != ContentStretch ||
		st.Order != 0 ||
		st.FlexGrow != 0 || st.FlexShrink != 1 || !st.FlexBasis.Auto() {
		t.Errorf("initial style = %+v", st)
	}
}

func TestDisplayFlexParses(t *testing.T) {
	st := flexRule(t, "div", "display: flex;")
	if st.Display != DisplayFlex {
		t.Errorf("display = %d want flex", st.Display)
	}
}

func TestFlexFlowShorthand(t *testing.T) {
	st := flexRule(t, "div", "flex-flow: column wrap;")
	if st.FlexDirection != FlexDirectionColumn || st.FlexWrap != FlexWrapWrap {
		t.Errorf("flex-flow = dir %d wrap %d want column/wrap", st.FlexDirection, st.FlexWrap)
	}
	st = flexRule(t, "div", "flex-flow: row;")
	if st.FlexDirection != FlexDirectionRow || st.FlexWrap != FlexWrapNowrap {
		t.Errorf("flex-flow row = dir %d wrap %d want row/nowrap", st.FlexDirection, st.FlexWrap)
	}
}
