package svg

import (
	"strings"
	"testing"
)

// The selector side is what turns a sheet into paint: an element is matched by
// its type, a class, an id, an attribute or the position it stands in, and by
// the logical pseudo-classes built on top of those. These check that every
// selector a drawing can answer is answered, that one it cannot never matches
// silently, and that a selector that is not CSS at all is said rather than
// guessed at.

// selNodeFor reads a drawing and answers the node of the element the id names,
// which is what a selector is tested against.
func selNodeFor(t *testing.T, src, id string) *selNode {
	t.Helper()
	root, err := parseDocument(src)
	if err != nil {
		t.Fatalf("parse %q: %v", src, err)
	}
	n := findSelNode(selTree(root, nil, nil), id)
	if n == nil {
		t.Fatalf("no element #%s in %q", id, src)
	}
	return n
}

// findSelNode walks the tree for the element an id names.
func findSelNode(n *selNode, id string) *selNode {
	if n.el.attr("id") == id {
		return n
	}
	for _, k := range n.kids {
		if got := findSelNode(k, id); got != nil {
			return got
		}
	}
	return nil
}

// matchesOne says whether a single-selector list stands on the node, read
// through the namespace scope a `<style>` of a drawing that binds the usual
// prefixes stands in.
func matchesOne(t *testing.T, sel string, n *selNode) bool {
	t.Helper()
	sels, err := parseSelectorList(sel, testNS)
	if err != nil {
		t.Fatalf("selector %q: %v", sel, err)
	}
	if len(sels) != 1 {
		t.Fatalf("selector %q read %d selectors, want 1", sel, len(sels))
	}
	return matchComplex(sels[0], n)
}

// testNS is the namespace scope a `<style>` stands in when the drawing binds
// the prefixes its own elements and attributes are written with: the SVG one,
// and the two XML keeps for a drawing's own names.
var testNS = namespace{"": svgNS, "svg": svgNS, "xlink": xlinkNS, "xml": xmlNS}

const ofTypeSrc = `<svg><g id="g"><rect id="a"/><circle id="b"/><rect id="c"/></g></svg>`

func TestSelectorOfType(t *testing.T) {
	for _, tc := range []struct {
		sel  string
		id   string
		want bool
	}{
		{"rect:first-of-type", "a", true},
		{"rect:first-of-type", "c", false},
		{"rect:last-of-type", "c", true},
		{"rect:last-of-type", "a", false},
		{"circle:only-of-type", "b", true},
		{"rect:only-of-type", "a", false},
		{"rect:nth-of-type(2)", "c", true},
		{"rect:nth-of-type(2)", "a", false},
		{"circle:nth-of-type(1)", "b", true},
		{"rect:nth-last-of-type(1)", "c", true},
		{"rect:nth-last-of-type(1)", "a", false},
	} {
		n := selNodeFor(t, ofTypeSrc, tc.id)
		if got := matchesOne(t, tc.sel, n); got != tc.want {
			t.Errorf("%q on #%s = %v, want %v", tc.sel, tc.id, got, tc.want)
		}
	}
}

const nthSrc = `<svg><g id="g"><rect id="1"/><rect id="2"/><rect id="3"/><rect id="4"/></g></svg>`

func TestSelectorNthChild(t *testing.T) {
	for _, tc := range []struct {
		sel  string
		id   string
		want bool
	}{
		{"rect:nth-child(1)", "1", true},
		{"rect:nth-child(2)", "2", true},
		{"rect:nth-child(2)", "3", false},
		{"rect:nth-child(2n)", "2", true},
		{"rect:nth-child(2n)", "4", true},
		{"rect:nth-child(2n)", "3", false},
		{"rect:nth-child(odd)", "3", true},
		{"rect:nth-child(even)", "3", false},
		{"rect:nth-child(2n+1)", "3", true},
		{"rect:nth-child(-n+2)", "2", true},
		{"rect:nth-child(-n+2)", "3", false},
		{"rect:nth-last-child(1)", "4", true},
		{"rect:nth-last-child(2)", "3", true},
	} {
		n := selNodeFor(t, nthSrc, tc.id)
		if got := matchesOne(t, tc.sel, n); got != tc.want {
			t.Errorf("%q on #%s = %v, want %v", tc.sel, tc.id, got, tc.want)
		}
	}
}

const nthOfSrc = `<svg><g id="g"><rect id="a" class="k"/><rect id="b"/><rect id="c" class="k"/></g></svg>`

func TestSelectorNthChildOfList(t *testing.T) {
	for _, tc := range []struct {
		sel  string
		id   string
		want bool
	}{
		{"rect:nth-child(1 of .k)", "a", true},
		{"rect:nth-child(2 of .k)", "c", true},
		{"rect:nth-child(1 of .k)", "c", false},
		{"rect:nth-child(2 of .k)", "b", false},
		{"rect:nth-last-child(1 of .k)", "c", true},
	} {
		n := selNodeFor(t, nthOfSrc, tc.id)
		if got := matchesOne(t, tc.sel, n); got != tc.want {
			t.Errorf("%q on #%s = %v, want %v", tc.sel, tc.id, got, tc.want)
		}
	}
}

const logicalSrc = `<svg><g id="g"><rect id="a" class="a"/><rect id="b" class="b"/><rect id="c"/></g></svg>`

func TestSelectorLogical(t *testing.T) {
	for _, tc := range []struct {
		sel  string
		id   string
		want bool
	}{
		{"rect:not(.a)", "b", true},
		{"rect:not(.a)", "a", false},
		{"rect:not(.a, .b)", "c", true},
		{"rect:not(.a, .b)", "a", false},
		{"rect:is(.a, .b)", "a", true},
		{"rect:is(.a, .b)", "b", true},
		{"rect:is(.a, .b)", "c", false},
		{"rect:where(.a)", "a", true},
		{"rect:where(.a)", "b", false},
	} {
		n := selNodeFor(t, logicalSrc, tc.id)
		if got := matchesOne(t, tc.sel, n); got != tc.want {
			t.Errorf("%q on #%s = %v, want %v", tc.sel, tc.id, got, tc.want)
		}
	}
}

const hasSrc = `<svg><g id="outer"><g id="inner"><rect id="r"/></g><g id="empty"/><rect id="first"/><circle id="sib"/></g></svg>`

func TestSelectorHas(t *testing.T) {
	for _, tc := range []struct {
		sel  string
		id   string
		want bool
	}{
		{"g:has(rect)", "outer", true},
		{"g:has(> rect)", "outer", true},
		{"g:has(rect)", "inner", true},
		{"g:has(> rect)", "inner", true},
		{"g:has(> rect)", "empty", false},
		{"g:has(> g > rect)", "outer", true},
		{"g:has(> g > rect)", "inner", false},
		{"rect:has(+ circle)", "first", true},
		{"rect:has(+ circle)", "sib", false},
		{"g:has(+ g)", "outer", false},
		{"rect:has(~ circle)", "first", true},
	} {
		n := selNodeFor(t, hasSrc, tc.id)
		if got := matchesOne(t, tc.sel, n); got != tc.want {
			t.Errorf("%q on #%s = %v, want %v", tc.sel, tc.id, got, tc.want)
		}
	}
}

func TestSelectorEmpty(t *testing.T) {
	src := `<svg><g id="g"><g id="self"/><g id="kids"><rect/></g><g id="ws"> </g><g id="bare"></g></g></svg>`
	for _, tc := range []struct {
		id   string
		want bool
	}{
		{"self", true},
		{"bare", true},
		{"kids", false},
		{"ws", false},
	} {
		n := selNodeFor(t, src, tc.id)
		if got := matchesOne(t, "g:empty", n); got != tc.want {
			t.Errorf(":empty on #%s = %v, want %v", tc.id, got, tc.want)
		}
	}
}

func TestSelectorScopeAndRoot(t *testing.T) {
	src := `<svg id="root"><g id="g"/></svg>`
	if !matchesOne(t, ":scope", selNodeFor(t, src, "root")) {
		t.Errorf(":scope did not match the root")
	}
	if matchesOne(t, ":scope", selNodeFor(t, src, "g")) {
		t.Errorf(":scope matched a non-root")
	}
	if !matchesOne(t, ":root", selNodeFor(t, src, "root")) {
		t.Errorf(":root did not match the root")
	}
}

func TestSelectorAttributeCaseFlag(t *testing.T) {
	src := `<svg><rect id="a" data-x="ABC"/></svg>`
	n := selNodeFor(t, src, "a")
	if !matchesOne(t, `[data-x="abc" i]`, n) {
		t.Errorf("an i flag did not compare without case")
	}
	if matchesOne(t, `[data-x="abc" s]`, n) {
		t.Errorf("an s flag compared without case")
	}
	if matchesOne(t, `[data-x="abc"]`, n) {
		t.Errorf("a plain comparison ignored the case")
	}
}

func TestSelectorNamespace(t *testing.T) {
	src := `<svg xmlns:xlink="http://www.w3.org/1999/xlink">` +
		`<rect id="a" data-x="1" xlink:href="#b"/></svg>`
	n := selNodeFor(t, src, "a")
	// A type naming the SVG namespace, or naming every one, or naming none at
	// all written as `*|`, stands on the rect the drawing keeps.
	for _, sel := range []string{"rect", "svg|rect", "*|rect", "*|*"} {
		if !matchesOne(t, sel, n) {
			t.Errorf("%q did not match the drawing's rect", sel)
		}
	}
	// No element a drawing paints is in no namespace, so `|rect` never matches.
	if matchesOne(t, "|rect", n) {
		t.Errorf("a no-namespace selector matched an element in the drawing")
	}
	// An attribute of no namespace is written plain or with `|`; `*|` is
	// either; and a prefix naming a namespace a drawing keeps no attribute in
	// stands aside.
	for _, sel := range []string{"[data-x]", "[|data-x]", "[*|data-x]"} {
		if !matchesOne(t, sel, n) {
			t.Errorf("%q did not match the drawing's plain attribute", sel)
		}
	}
	if matchesOne(t, "[svg|data-x]", n) {
		t.Errorf("[svg|data-x] matched a plain attribute, where SVG keeps none")
	}
	// The href XLink writes is reached only by the prefix bound to it.
	if !matchesOne(t, "[xlink|href]", n) {
		t.Errorf("[xlink|href] did not reach the drawing's XLink href")
	}
	if !matchesOne(t, "[*|href]", n) {
		t.Errorf("[*|href] did not reach the drawing's XLink href")
	}
	if matchesOne(t, "[href]", n) {
		t.Errorf("[href] matched an href of no namespace, which the drawing has none of")
	}
	// A type naming a namespace that is not the SVG one never matches.
	if matchesOne(t, "xlink|rect", n) {
		t.Errorf("a type of the XLink namespace matched an SVG element")
	}
}

func TestASelectorWearingAnUnboundPrefixIsSaid(t *testing.T) {
	for _, bad := range []string{"foo|rect", "[foo|href]"} {
		if _, err := parseSelectorList(bad, testNS); err == nil {
			t.Errorf("parseSelectorList(%q) = nil error, want one for a prefix nothing binds", bad)
		}
	}
}

func TestSelectorNeverMatchingIsSilent(t *testing.T) {
	for _, sel := range []string{"rect:hover", "rect:focus", "rect::before", "rect:before"} {
		n := selNodeFor(t, `<svg><rect id="a"/></svg>`, "a")
		if matchesOne(t, sel, n) {
			t.Errorf("%q matched, want a selector a still drawing cannot answer to stand aside", sel)
		}
	}
}

func TestSelectorSpecificityOfPseudoClasses(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want int
	}{
		{"rect:first-child", 101},
		{":root", 100},
		{":not(#x)", 10000},
		{":is(rect)", 1},
		{":where(#x)", 0},
		{"rect:nth-child(2)", 101},
		{"rect::before", 2},
		{"rect:lang(en)", 101},
		{"rect:dir(rtl)", 101},
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

const langDirSrc = `<svg>` +
	`<g id="g" xml:lang="en">` +
	`<rect id="plain"/><rect id="us" xml:lang="en-US"/><rect id="fr" xml:lang="fr"/>` +
	`</g>` +
	`<g id="he" dir="rtl"><rect id="inherit"/><rect id="auto" dir="auto">שלום</rect></g>` +
	`</svg>`

func TestSelectorLang(t *testing.T) {
	for _, tc := range []struct {
		sel  string
		id   string
		want bool
	}{
		{":lang(en)", "plain", true},
		{":lang(fr)", "plain", false},
		{":lang(en)", "us", true},
		{":lang(en-US)", "us", true},
		{":lang(en)", "fr", false},
		{":lang(fr, en)", "fr", true},
		{":lang('en')", "plain", true},
	} {
		n := selNodeFor(t, langDirSrc, tc.id)
		if got := matchesOne(t, tc.sel, n); got != tc.want {
			t.Errorf("%q on #%s = %v, want %v", tc.sel, tc.id, got, tc.want)
		}
	}
}

func TestSelectorDir(t *testing.T) {
	for _, tc := range []struct {
		sel  string
		id   string
		want bool
	}{
		{":dir(ltr)", "inherit", false},
		{":dir(rtl)", "inherit", true},
		{":dir(rtl)", "auto", true},
		{":dir(ltr)", "auto", false},
		{":dir(ltr)", "plain", true},
		{":dir(rtl)", "plain", false},
	} {
		n := selNodeFor(t, langDirSrc, tc.id)
		if got := matchesOne(t, tc.sel, n); got != tc.want {
			t.Errorf("%q on #%s = %v, want %v", tc.sel, tc.id, got, tc.want)
		}
	}
}

func TestSheetCascadesThroughLogicalPseudoClasses(t *testing.T) {
	// `:where` counts for nothing, so a class selector beats it however late
	// the `:where` rule is written.
	got := cssFill(t, `<style>.b { fill: red } :where(.a) { fill: blue }</style>`+
		`<rect class="a b" width="10" height="10"/>`)
	if got != cssRed {
		t.Errorf("fill = %v, want :where's zero weight to leave the class's red", got)
	}
	// `:not` carries the weight of its argument, so this beats the class.
	got = cssFill(t, `<style>.a { fill: blue } rect:not(.b) { fill: red }</style>`+
		`<rect class="a" width="10" height="10"/>`)
	if got != cssRed {
		t.Errorf("fill = %v, want :not's weight to win", got)
	}
}

func TestASelectorThatIsNotCSSIsSaid(t *testing.T) {
	img := mustParse(t, `<svg viewBox="0 0 10 10">`+
		`<style>rect:bogus { fill: red }</style><rect width="10" height="10"/></svg>`)
	if got := pixelAt(img.Render(10, 10), 5, 5); got != cssBlack {
		t.Errorf("fill = %v, want an unreadable selector left out", got)
	}
	if w := img.Warnings().String(); !strings.Contains(w, "bogus") {
		t.Errorf("the drawing said %v, want it to name the selector it could not read", w)
	}
}
