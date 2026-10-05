package css

import (
	"strings"
	"testing"
)

// TestSupportsEvaluatesTheDeclaration puts the engine's own cascade under the
// condition: a property it has, holding a value it reads, is a test it passes.
func TestSupportsEvaluatesTheDeclaration(t *testing.T) {
	cases := []struct {
		cond string
		want bool
		why  string
	}{
		{"(display: grid)", true, "the engine lays out grids"},
		{"(display: flex)", true, "and flex as well"},
		{"(display: subgrid)", false, "a value no grammar here reads"},
		{"(display: 5)", false, "neither does this one, and it fails in silence"},
		{"(color: #00ff00)", true, "colours parse"},
		{"(color: nonsense)", false, "an unreadable colour is no colour at all"},
		{"(no-such-property: 1px)", false, "a property the engine has no case for"},
		{"(-webkit-box-orient: horizontal)", false, "a vendor prefix is dropped outright"},
		{"(--brand: red)", true, "a custom property is the author's own"},
		{"(float: left)", true, "read, though the flow cannot apply it"},
		{"(color: initial)", true, "a cascade keyword on a property it has"},
		{"(no-such-property: initial)", false, "and no keyword can vouch for a name it has not"},
		{"(width: calc(100% - 20px))", true, "a formula the engine evaluates"},
		{"(min-height: 40vh)", true, "viewport lengths too"},
	}
	for _, c := range cases {
		got, warns := supportsCondition(c.cond)
		if got != c.want {
			t.Errorf("supportsCondition(%s) = %v, want %v (%s)", c.cond, got, c.want, c.why)
		}
		if len(warns) != 0 {
			t.Errorf("supportsCondition(%s) warned %q; a test the engine can read never warns", c.cond, warns)
		}
	}
}

// TestSupportsReadsTheLogic walks the condition's own grammar: and, or, not
// and the groups that carry them.
func TestSupportsReadsTheLogic(t *testing.T) {
	cases := []struct {
		cond string
		want bool
	}{
		{"(display: grid) and (color: #fff)", true},
		{"(display: grid) and (color: nonsense)", false},
		{"(display: subgrid) or (color: #fff)", true},
		{"(display: subgrid) or (color: nonsense)", false},
		{"not (display: subgrid)", true},
		{"not (display: grid)", false},
		{"not not (display: grid)", true},
		{"(display: grid) and not (color: nonsense)", true},
		{"(display: grid) or not (color: nonsense)", true},
		{"(not (display: grid))", false},
		{"((display: grid))", true},
		{"((display: subgrid) or (color: #fff))", true},
		{"(display: GRID)", true},
		{"  ( display : grid )  ", true},
	}
	for _, c := range cases {
		got, warns := supportsCondition(c.cond)
		if got != c.want {
			t.Errorf("supportsCondition(%s) = %v, want %v", c.cond, got, c.want)
		}
		if len(warns) != 0 {
			t.Errorf("supportsCondition(%s) warned %q", c.cond, warns)
		}
	}
}

// TestSupportsReadsSelectorFunction asks selector() with the reader the engine
// runs every rule through: a selector it would keep is one it can evaluate.
func TestSupportsReadsSelectorFunction(t *testing.T) {
	cases := []struct {
		cond string
		want bool
		why  string
	}{
		{"selector(button)", true, "a tag it reads"},
		{"selector(.card)", true, "a class it reads"},
		{"selector(:hover)", true, "a state it tracks"},
		{"selector(*)", true, "the universal selector"},
		{"selector(button.primary:hover)", true, "several at once"},
		{"selector(button, .card)", true, "every selector in the list"},
		{"selector(.a .b)", false, "a combinator the flat model cannot evaluate"},
		{"selector(.a > .b)", false, "nor a child one"},
		{"selector(input[type=text])", false, "nor an attribute"},
		{"selector(p::before)", false, "nor a pseudo-element"},
		{"selector(:nth-child(2))", false, "nor a structural pseudo-class"},
		{"selector(:not(.a))", false, "nor :not, which it records and never fires"},
		{"selector(???)", false, "nothing in it to select a widget with"},
		{"selector(button) and (display: grid)", true, "and the group beside it"},
		{"selector(.a .b) or (display: grid)", true, "one side is enough"},
		{"selector(button) and (display: subgrid)", false, "both sides must hold"},
	}
	for _, c := range cases {
		got, warns := supportsCondition(c.cond)
		if got != c.want {
			t.Errorf("supportsCondition(%s) = %v, want %v (%s)", c.cond, got, c.want, c.why)
		}
		if len(warns) != 0 {
			t.Errorf("supportsCondition(%s) warned %q", c.cond, warns)
		}
	}
}

// TestSupportsReportsTheConditionItCannotRead keeps the two answers apart: a
// test that fails is false in silence, a condition that cannot be read is
// false and said so.
func TestSupportsReportsTheConditionItCannotRead(t *testing.T) {
	bad := []string{
		``,
		`display: grid`,
		`display: grid and`,
		`(display: grid) trailing`,
		`(display: grid) andrew`,
		`kaboom(x)`,
		`(display: grid`,
		`(color: "red)`,
		`((display: grid) and)`,
	}
	for _, cond := range bad {
		got, warns := supportsCondition(cond)
		if got {
			t.Errorf("supportsCondition(%q) = true, want false for a condition nothing reads", cond)
		}
		if len(warns) == 0 || !strings.Contains(warns[0], "@supports") {
			t.Errorf("supportsCondition(%q) warns = %q, want one about @supports", cond, warns)
		}
	}
}

// TestSupportsFontFunctionsAnswerNo keeps the two font functions readable and
// false: no font file is ever opened, so there is no format or technology to
// have.
func TestSupportsFontFunctionsAnswerNo(t *testing.T) {
	for _, cond := range []string{
		`(font-format("woff2"))`,
		`(font-tech(color-COLRv1))`,
		`font-format(woff2)`,
	} {
		got, warns := supportsCondition(cond)
		if got || len(warns) != 0 {
			t.Errorf("supportsCondition(%s) = %v, warns %q; want false in silence", cond, got, warns)
		}
	}
}

// TestSupportsGatesTheBlock puts the condition in the sheet: what it passes
// keeps its rules, what it fails drops them, and only the unreadable condition
// leaves a word behind.
func TestSupportsGatesTheBlock(t *testing.T) {
	sh := parseOK(t, `
@supports (display: grid) { .yes { width: 90%; } }
@supports (display: subgrid) { .no { width: 90%; } }
@supports not (display: subgrid) { .negated { width: 90%; } }
@supports (display: grid) and (color: #fff) { .together { width: 90%; } }
@supports (display: grid) and (color: nonsense) { .apart { width: 90%; } }
@supports selector(.a .b) { .comb { width: 90%; } }
@supports ??? { .unreadable { width: 90%; } }
`)
	for _, c := range []struct {
		class string
		want  bool
	}{
		{"yes", true},
		{"no", false},
		{"negated", true},
		{"together", true},
		{"apart", false},
		{"comb", false},
		{"unreadable", false},
	} {
		_, has := sh.Property("div", []string{c.class}, "width", 800)
		if has != c.want {
			t.Errorf("@supports kept .%s = %v, want %v", c.class, has, c.want)
		}
	}
	if len(sh.Warn) != 1 || !strings.Contains(sh.Warn[0], "@supports") {
		t.Errorf("sh.Warn = %q, want one about @supports", sh.Warn)
	}
}

// TestSupportsStatementFormWithNoBlock reads a @supports that gates nothing.
func TestSupportsStatementFormWithNoBlock(t *testing.T) {
	sh := parseOK(t, `@supports (display: grid); .a { color: red; }`)
	if _, has := sh.Property("div", []string{"a"}, "color", 800); !has {
		t.Error("the rule after a @supports statement should still apply")
	}
	if len(sh.Warn) != 0 {
		t.Errorf("sh.Warn = %q, want none", sh.Warn)
	}
}

// TestSupportsMediaInsideAndOutside keeps the two tests in their own places: a
// @media inside a @supports still gates on the window.
func TestSupportsMediaInsideAndOutside(t *testing.T) {
	sh := parseOK(t, `
@supports (display: grid) {
  @media (min-width: 720px) { .wide { background-color: blue; } }
}
@supports (display: subgrid) {
  @media (min-width: 720px) { .dropped { background-color: red; } }
}
`)
	if st := sh.Style("div", []string{"wide"}, StateNone, 800); !st.Has("background-color") {
		t.Error("a media inside a passing @supports should apply when the width fits")
	}
	if st := sh.Style("div", []string{"dropped"}, StateNone, 800); st.Has("background-color") {
		t.Error("a media inside a failing @supports must never reach the sheet")
	}
}
