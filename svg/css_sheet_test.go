package svg

import (
	"strings"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// A drawing may say what it looks like in a `<style>` sheet of its own, as well
// as in the attributes on the elements. These check that the sheet is read, that
// its selectors find the elements they name, and that the cascade settles a
// thing written more than once the way CSS says it should — by importance, by
// the weight of the selector, by the order it was written, and under the
// element's own `style` and its attributes.

var (
	cssRed   = canvas.RGB(255, 0, 0)
	cssBlue  = canvas.RGB(0, 0, 255)
	cssBlack = canvas.RGB(0, 0, 0)
)

// cssFill paints one drawing of a single square with the sheet and shape given
// and answers the colour at its middle, which is the fill the cascade settled
// on.
func cssFill(t *testing.T, body string) canvas.Color {
	t.Helper()
	_, cv := painted(t, `<svg viewBox="0 0 10 10">`+body+`</svg>`, 10, 10)
	return pixelAt(cv, 5, 5)
}

func TestSheetMatchesAType(t *testing.T) {
	got := cssFill(t, `<style>rect { fill: red }</style><rect width="10" height="10"/>`)
	if got != cssRed {
		t.Errorf("fill = %v, want the red the sheet named the type with", got)
	}
}

func TestSheetMatchesAClass(t *testing.T) {
	got := cssFill(t, `<style>.a { fill: red }</style><rect class="x a y" width="10" height="10"/>`)
	if got != cssRed {
		t.Errorf("fill = %v, want the red the sheet named the class with", got)
	}
}

func TestSheetMatchesAnID(t *testing.T) {
	got := cssFill(t, `<style>#me { fill: red }</style><rect id="me" width="10" height="10"/>`)
	if got != cssRed {
		t.Errorf("fill = %v, want the red the sheet named the id with", got)
	}
}

func TestSheetMatchesTheUniversalSelector(t *testing.T) {
	got := cssFill(t, `<style>* { fill: red }</style><rect width="10" height="10"/>`)
	if got != cssRed {
		t.Errorf("fill = %v, want the red the universal selector gave everything", got)
	}
}

func TestSheetMatchesAnAttribute(t *testing.T) {
	got := cssFill(t, `<style>[width="10"] { fill: red }</style><rect width="10" height="10"/>`)
	if got != cssRed {
		t.Errorf("fill = %v, want the red the attribute selector named", got)
	}
}

func TestSheetDescendantAndChildCombinators(t *testing.T) {
	desc := cssFill(t, `<style>svg rect { fill: red }</style><rect width="10" height="10"/>`)
	if desc != cssRed {
		t.Errorf("fill = %v, want a descendant combinator to reach the rect", desc)
	}
	child := cssFill(t, `<style>g > rect { fill: red }</style><g><rect width="10" height="10"/></g>`)
	if child != cssRed {
		t.Errorf("fill = %v, want a child combinator to reach the rect", child)
	}
	// A child combinator is the direct child and not anything below it.
	grand := cssFill(t, `<style>svg > rect { fill: red }</style><g><rect width="10" height="10"/></g>`)
	if grand != cssBlack {
		t.Errorf("fill = %v, want the grandchild left black", grand)
	}
}

func TestSheetSiblingCombinator(t *testing.T) {
	got := cssFill(t, `<style>rect + circle { fill: red }</style>`+
		`<rect width="1" height="1"/><circle cx="5" cy="5" r="5"/>`)
	if got != cssRed {
		t.Errorf("fill = %v, want the next sibling to be reached", got)
	}
}

func TestSheetSpecificityAndOrder(t *testing.T) {
	// The id outweighs the class, the class the type, however they are ordered.
	id := cssFill(t, `<style>.a { fill: blue } rect { fill: blue } #me { fill: red }</style>`+
		`<rect id="me" class="a" width="10" height="10"/>`)
	if id != cssRed {
		t.Errorf("fill = %v, want the id's red to win", id)
	}
	class := cssFill(t, `<style>rect { fill: blue } .a { fill: red }</style>`+
		`<rect class="a" width="10" height="10"/>`)
	if class != cssRed {
		t.Errorf("fill = %v, want the class's red to win over the type", class)
	}
	// Between two of the same weight the one written later wins.
	order := cssFill(t, `<style>.a { fill: blue } .a { fill: red }</style>`+
		`<rect class="a" width="10" height="10"/>`)
	if order != cssRed {
		t.Errorf("fill = %v, want the later rule's red to win", order)
	}
}

func TestSheetCascadeAgainstAttributesAndStyle(t *testing.T) {
	// A presentation attribute is the weakest of all: any rule beats it.
	attr := cssFill(t, `<style>rect { fill: red }</style><rect fill="blue" width="10" height="10"/>`)
	if attr != cssRed {
		t.Errorf("fill = %v, want the sheet's red over the attribute's blue", attr)
	}
	// The `style` attribute beats a normal rule.
	inline := cssFill(t, `<style>rect { fill: red }</style>`+
		`<rect style="fill: blue" width="10" height="10"/>`)
	if inline != cssBlue {
		t.Errorf("fill = %v, want the style attribute's blue over the rule's red", inline)
	}
	// An important rule beats the `style` attribute.
	important := cssFill(t, `<style>rect { fill: red !important }</style>`+
		`<rect style="fill: blue" width="10" height="10"/>`)
	if important != cssRed {
		t.Errorf("fill = %v, want the important red over the style attribute", important)
	}
	// And an important declaration in the `style` attribute beats everything.
	first := cssFill(t, `<style>rect { fill: red !important }</style>`+
		`<rect style="fill: blue !important" width="10" height="10"/>`)
	if first != cssBlue {
		t.Errorf("fill = %v, want the important style attribute to be last", first)
	}
}

func TestSheetInherits(t *testing.T) {
	got := cssFill(t, `<style>g { fill: red }</style><g><rect width="10" height="10"/></g>`)
	if got != cssRed {
		t.Errorf("fill = %v, want the group's red inherited by the rect", got)
	}
}

func TestSheetSelectorList(t *testing.T) {
	got := cssFill(t, `<style>rect, circle { fill: red }</style><rect width="10" height="10"/>`)
	if got != cssRed {
		t.Errorf("fill = %v, want the rect named in the list to take the fill", got)
	}
}

func TestSheetStripsComments(t *testing.T) {
	got := cssFill(t, `<style>/* a } comment */ rect { fill: red }</style><rect width="10" height="10"/>`)
	if got != cssRed {
		t.Errorf("fill = %v, want the comment skipped and the rule read", got)
	}
}

func TestSheetLeavesOutWhatItCannotRead(t *testing.T) {
	// A rule asking for a state a still drawing has none of stands aside: it
	// is a valid selector, so nothing is said, and the drawing is painted as
	// if the rule had not named it.
	img := mustParse(t, `<svg viewBox="0 0 10 10">`+
		`<style>rect:hover { fill: red }</style><rect width="10" height="10"/></svg>`)
	if got := pixelAt(img.Render(10, 10), 5, 5); got != cssBlack {
		t.Errorf("fill = %v, want the rect left black rather than :hover applied", got)
	}
	if w := img.Warnings().String(); strings.Contains(w, "hover") {
		t.Errorf("the drawing said %v, want a valid :hover left unremarked", w)
	}
	// An at-rule a still drawing does not use is skipped whole, and said so.
	img = mustParse(t, `<svg viewBox="0 0 10 10">`+
		`<style>@media screen { rect { fill: red } }</style><rect width="10" height="10"/></svg>`)
	if got := pixelAt(img.Render(10, 10), 5, 5); got != cssBlack {
		t.Errorf("fill = %v, want the at-rule left out rather than half-applied", got)
	}
	if !strings.Contains(img.Warnings().String(), "media") {
		t.Errorf("the drawing said %v, want it to name the at-rule it left out", img.Warnings())
	}
}

func TestParseSelectorList(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want int
	}{
		{"rect", 1},
		{"rect, circle.foo", 2},
		{"g > rect", 1},
		{"#a .b[c] *", 1},
		{"rect:hover", 1},
		{":nth-child(2)", 1},
		{"li:nth-child(2n+1)", 1},
		{"a:not(.x, .y)", 1},
		{":is(h1, h2, h3)", 1},
		{":has(> rect)", 1},
		{"::before", 1},
		{"svg|rect", 1},
		{"[fill=\"red\" i]", 1},
		{"", 0},
	} {
		got, err := parseSelectorList(tc.in, testNS)
		if err != nil {
			t.Errorf("parseSelectorList(%q) = %v, want no error", tc.in, err)
			continue
		}
		if len(got) != tc.want {
			t.Errorf("parseSelectorList(%q) read %d selectors, want %d", tc.in, len(got), tc.want)
		}
	}
	for _, bad := range []string{"rect:bogus", "rect:not()", ":nth-child(x)", "rect >", "[", "rect..", ":bogus(1)"} {
		if _, err := parseSelectorList(bad, testNS); err == nil {
			t.Errorf("parseSelectorList(%q) = nil error, want one", bad)
		}
	}
}

func TestSpecificity(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want int
	}{
		{"*", 0},
		{"rect", 1},
		{".a", 100},
		{"[fill]", 100},
		{"#x", 10000},
		{"#x.a[b]", 10200},
		{"#x rect.a", 10101},
	} {
		sel, err := parseSelectorList(tc.in, testNS)
		if err != nil {
			t.Fatalf("parseSelectorList(%q): %v", tc.in, err)
		}
		if got := specificity(sel[0]); got != tc.want {
			t.Errorf("specificity(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}
