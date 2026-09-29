package css

import "testing"

func TestColumnCountParsesNumbersAndAuto(t *testing.T) {
	for raw, want := range map[string]int{
		"3":    3,
		"1":    1,
		"12":   12,
		"auto": 0,
		" 0  ": 0,
	} {
		if v, ok := parseColumnCount(raw); !ok || v != want {
			t.Errorf("%q = %d,%v want %d,true", raw, v, ok, want)
		}
	}
	for _, raw := range []string{"", "two", "-2", "3.5", "3px"} {
		if _, ok := parseColumnCount(raw); ok {
			t.Errorf("%q must not parse", raw)
		}
	}
}

func TestColumnWidthParsesLengths(t *testing.T) {
	for _, raw := range []string{"200px", "12em", "40%"} {
		l, ok := parseColumnWidth(raw, Units{Font: 16, Width: 800})
		if !ok {
			t.Errorf("%q must parse", raw)
			continue
		}
		if l.value <= 0 {
			t.Errorf("%q = %v want a positive width", raw, l)
		}
	}
	if l, ok := parseColumnWidth("auto", Units{}); !ok || !l.Auto() {
		t.Errorf(`"auto" = %v,%v want auto,true`, l, ok)
	}
	if l, ok := parseColumnWidth("16em", Units{Font: 16}); !ok || l.Resolve(Units{Font: 16}) != 256 {
		t.Errorf("16em = %v,%v want 256", l, ok)
	}
	for _, raw := range []string{"", "-4px", "wide"} {
		if _, ok := parseColumnWidth(raw, Units{}); ok {
			t.Errorf("%q must not parse", raw)
		}
	}
}

func TestColumnsShorthandTakesCountAndWidth(t *testing.T) {
	for _, tc := range []struct {
		raw      string
		count    int
		widthRef int
	}{
		{"3", 3, -1},
		{"200px", 0, 200},
		{"3 200px", 3, 200},
		{"200px 3", 3, 200},
		{"auto", 0, -1},
	} {
		count, l, ok := parseColumns(tc.raw, Units{})
		if !ok || count != tc.count {
			t.Errorf("%q count = %d,%v want %d", tc.raw, count, ok, tc.count)
		}
		switch {
		case tc.widthRef < 0 && !l.Auto():
			t.Errorf("%q width = %v want auto", tc.raw, l)
		case tc.widthRef >= 0 && l.Resolve(Units{}) != tc.widthRef:
			t.Errorf("%q width = %v want %d", tc.raw, l, tc.widthRef)
		}
	}
	for _, raw := range []string{"", "3 sideways", "-2"} {
		if _, _, ok := parseColumns(raw, Units{}); ok {
			t.Errorf("%q must not parse", raw)
		}
	}
}

func TestColumnFillParsesKeywords(t *testing.T) {
	for raw, want := range map[string]uint8{
		"balance": ColumnFillBalance,
		"auto":    ColumnFillAuto,
	} {
		if v, ok := parseColumnFill(raw); !ok || v != want {
			t.Errorf("%q = %d,%v want %d,true", raw, v, ok, want)
		}
	}
	if _, ok := parseColumnFill("fill"); ok {
		t.Error(`"fill" must not parse`)
	}
}

func TestBreakInsideParsesAvoidSpellings(t *testing.T) {
	for _, raw := range []string{"auto"} {
		if v, ok := parseBreakInside(raw); !ok || v != BreakAuto {
			t.Errorf("%q = %d,%v want BreakAuto", raw, v, ok)
		}
	}
	for _, raw := range []string{"avoid", "avoid-column", "avoid-page", "avoid-region"} {
		if v, ok := parseBreakInside(raw); !ok || v != BreakAvoid {
			t.Errorf("%q = %d,%v want BreakAvoid", raw, v, ok)
		}
	}
	for _, raw := range []string{"", "always", "avoid-line"} {
		if _, ok := parseBreakInside(raw); ok {
			t.Errorf("%q must not parse", raw)
		}
	}
}

func TestFloatAndClearParseKeywords(t *testing.T) {
	for raw, want := range map[string]uint8{
		"none":         FloatNone,
		"left":         FloatLeft,
		"right":        FloatRight,
		"inline-start": FloatInlineStart,
		"inline-end":   FloatInlineEnd,
	} {
		if v, ok := parseFloat(raw); !ok || v != want {
			t.Errorf("float %q = %d,%v want %d,true", raw, v, ok, want)
		}
	}
	for raw, want := range map[string]uint8{
		"none":  ClearNone,
		"left":  ClearLeft,
		"right": ClearRight,
		"both":  ClearBoth,
	} {
		if v, ok := parseClear(raw); !ok || v != want {
			t.Errorf("clear %q = %d,%v want %d,true", raw, v, ok, want)
		}
	}
	if _, ok := parseFloat("up"); ok {
		t.Error(`float "up" must not parse`)
	}
	if _, ok := parseClear("up"); ok {
		t.Error(`clear "up" must not parse`)
	}
}

func TestColumnRuleShorthandSplitsComponents(t *testing.T) {
	st := Style{}
	set := map[string]bool{}
	applyColumnRule(&st, set, "2px solid red", func() {})
	if st.ColumnRuleWidth != 2 || st.ColumnRuleStyle != BorderSolid {
		t.Errorf("width/style = %d/%d want 2/solid", st.ColumnRuleWidth, st.ColumnRuleStyle)
	}
	if st.ColumnRuleColor != namedColors["red"] {
		t.Errorf("colour = %v want red", st.ColumnRuleColor)
	}
	if !set["column-rule-width"] || !set["column-rule-style"] || !set["column-rule-color"] {
		t.Errorf("longhands not recorded: %v", set)
	}
}

func TestColumnRuleShorthandRejectsUnrelatedValues(t *testing.T) {
	st := Style{}
	if applyColumnRule(&st, map[string]bool{}, "wobbly", func() {}) {
		t.Error("a value naming no component must not be taken")
	}
	if applyColumnRule(&st, map[string]bool{}, "", func() {}) {
		t.Error("an empty value must not be taken")
	}
}

// TestMultiColumnRuleFoldsIntoTheStyle walks a declaration through the cascade
// the way a stylesheet reaches a template, including the width in em and the
// colour keyword currentColor.
func TestMultiColumnRuleFoldsIntoTheStyle(t *testing.T) {
	sh, err := Parse("body { color: rgb(1, 2, 3); column-count: 2; column-width: 20em; column-gap: 12px; column-rule: 1px solid currentColor; column-fill: auto; }")
	if err != nil {
		t.Fatal(err)
	}
	st := sh.Style("body", nil, StateNone, 800)
	if st.ColumnCount != 2 {
		t.Errorf("column-count = %d want 2", st.ColumnCount)
	}
	if got := st.ColumnWidth.Resolve(Units{Font: 16}); got != 320 {
		t.Errorf("column-width = %d want 320", got)
	}
	if got := st.ColumnGap.Resolve(Units{}); got != 12 {
		t.Errorf("column-gap = %d want 12", got)
	}
	if st.ColumnRuleWidth != 1 || st.ColumnRuleStyle != BorderSolid {
		t.Errorf("rule = %d/%d want 1/solid", st.ColumnRuleWidth, st.ColumnRuleStyle)
	}
	if st.ColumnRuleColor != st.Color {
		t.Errorf("currentColor rule = %v want the colour %v", st.ColumnRuleColor, st.Color)
	}
	if st.ColumnFill != ColumnFillAuto {
		t.Errorf("column-fill = %d want auto", st.ColumnFill)
	}
}

func TestBreakInsideLonghandAndInitial(t *testing.T) {
	st := flexRule(t, "flex", "break-inside: avoid;")
	if st.BreakInside != BreakAvoid {
		t.Errorf("break-inside = %d want avoid", st.BreakInside)
	}
	if !st.Has("break-inside") {
		t.Error("break-inside not recorded as set")
	}
	st = flexRule(t, "flex", "page-break-inside: avoid;")
	if st.BreakInside != BreakAvoid || !st.Has("break-inside") {
		t.Error("the older page-break-inside name must reach the same field")
	}
	st = flexRule(t, "flex", "break-inside: revert;")
	if st.BreakInside != BreakAuto {
		t.Error("revert must put the initial value back")
	}
}

// TestFloatAndClearWarnWithoutBeingSpelledWrong keeps a stylesheet that says
// float from being reported as a property the engine has never heard of: the
// value is read, and the warning says what this flow cannot do with it.
func TestFloatAndClearWarnWithoutBeingSpelledWrong(t *testing.T) {
	sh, err := Parse("label { float: right; clear: both; }")
	if err != nil {
		t.Fatal(err)
	}
	st := sh.Style("label", nil, StateNone, 800)
	if st.Float != FloatRight || st.Clear != ClearBoth {
		t.Errorf("float/clear = %d/%d want right/both", st.Float, st.Clear)
	}
	if len(sh.Warn) != 2 {
		t.Fatalf("warnings = %v want one per declaration", sh.Warn)
	}
	for _, w := range sh.Warn {
		if w == "" {
			t.Error("empty warning")
		}
	}
}

// TestFloatAndClearNoneIsTheInitialAndSaysNothing: a stylesheet that asks for
// no floating at all is asking for what the flow already does, so it is not a
// warning — only a value that asks for something is.
func TestFloatAndClearNoneIsTheInitialAndSaysNothing(t *testing.T) {
	sh, err := Parse("label { float: none; clear: none; }")
	if err != nil {
		t.Fatal(err)
	}
	st := sh.Style("label", nil, StateNone, 800)
	if st.Float != FloatNone || st.Clear != ClearNone {
		t.Errorf("float/clear = %d/%d want none/none", st.Float, st.Clear)
	}
	if len(sh.Warn) != 0 {
		t.Errorf("warnings = %v want none for the initial values", sh.Warn)
	}
}
