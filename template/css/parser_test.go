package css

// These tests pin down the parser's contract: a stylesheet written in valid
// CSS3 never fails to parse and never hangs, and the parts of the language the
// flat widget model cannot evaluate are kept as never-matching rules with a
// warning instead of being rejected. Only structural damage is an error.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// writeTemp writes a stylesheet to a temp file and returns its path.
func writeTemp(t *testing.T, name, sheet string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(sheet), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// parseOK is the itch a tolerant parser must scratch: the sheet parses, and
// every selector it carried survives as a rule.
func parseOK(t *testing.T, sheet string) *Sheet {
	t.Helper()
	sh, err := Parse(sheet)
	if err != nil {
		t.Fatalf("Parse(%q) errored on valid CSS3: %v", sheet, err)
	}
	return sh
}

func TestParseSelectorTolerant(t *testing.T) {
	sheets := []string{
		`.x { color: red; }`,
		`BUTTON { color: red; }`,
		`div.container > span.icon:hover { color: red; }`,
		`nav ul li a { color: red; }`,
		`h1 + h2 { color: red; }`,
		`p ~ em { color: red; }`,
		`input[type="text"] { color: red; }`,
		`[disabled] { color: red; }`,
		`.abc[class~="x" i] { color: red; }`,
		`a[href^="https://"]:focus::before { color: red; }`,
		`p::before { color: red; }`,
		`div::-webkit-scrollbar { color: red; }`,
		`li:nth-child(2n+1) { color: red; }`,
		`p:nth-of-type(3) { color: red; }`,
		`:not(.foo, .bar) { color: red; }`,
		`main :is(h1, h2, h3) { color: red; }`,
		`.card:where(.wide, .tall) { color: red; }`,
		`a:has(> img) { color: red; }`,
		`.café { color: red; }`,
		`.\\2014 { color: red; }`,
		`.\\.escaped { color: red; }`,
		`svg|a { color: red; }`,
		`* { color: red; }`,
		`div/* hole */span:hover { color: red; }`,
		`#widget-id { color: red; }`,
		`.a:hover:focus { color: red; }`,
	}
	for _, sheet := range sheets {
		sh := parseOK(t, sheet)
		if len(sh.Rules()) == 0 {
			t.Errorf("%q parsed but produced no rule", sheet)
		}
	}
}

func TestParseAtRuleTolerant(t *testing.T) {
	sheets := []string{
		`@media (min-width: 600px) { body { color: red; } }`,
		`@media screen and (min-width: 600px) { body { color: red; } }`,
		`@supports (display: grid) { .g { display: grid; } }`,
		`@layer components { .btn { color: red; } }`,
		`@layer reset, base, components;`,
		`@container (min-width: 500px) { .card { width: 90%; } }`,
		`@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }`,
		`@font-face { font-family: "Inter"; src: url("inter.woff2") format("woff2"); }`,
		`@page :first { margin: 2cm; }`,
		`@import url("other.css") screen;`,
		`@charset "utf-8";`,
		`@property --brand { syntax: "<color>"; inherits: false; initial-value: red; }`,
		`@namespace svg url(http://www.w3.org/2000/svg);`,
		`@unknown-rule foo bar;`,
		`@unknown-block foo { a { b: c; } }`,
	}
	for _, sheet := range sheets {
		parseOK(t, sheet)
	}
}

func TestParseDeclarationTolerant(t *testing.T) {
	sheets := []string{
		`.b { background-image: url("a(b).png"); }`,
		`.c { content: "a;b"; }`,
		`.d { font: 14px/1.5 "Helvetica Neue", Arial, sans-serif; }`,
		`.e { background: linear-gradient(45deg, #f00, #00f); }`,
		`.f { box-shadow: 0 1px 2px rgba(0, 0, 0, 0.5); }`,
		`.g { grid-template-columns: repeat(auto-fill, minmax(150px, 1fr)); }`,
		`.h { margin: calc(100% - 20px); }`,
		`.i { content: 'it\'s fine'; }`,
		`.j { margin: 1px /* spacer */ 2px; }`,
		`.k { width: 50%; height: 32px; }`,
		`.l { color: rgb(255 0 0 / 50%); }`,
		`.m {
			color: red;
			/* between declarations */
			background-color: blue /* after value no semicolon */
		}`,
	}
	for _, sheet := range sheets {
		parseOK(t, sheet)
	}
}

// The corpus is stylesheets written by real tools and manuals — a CSS reset,
// a Tailwind preflight re-creation and MDN samples — the kind of file a
// browser opens every day and this engine must never refuse.
func TestParseRealWorldCSSZeroErrors(t *testing.T) {
	reset := `
html, body, div, span, applet, object, iframe, h1, h2, h3, h4, h5, h6, p,
blockquote, pre, a, abbr, acronym, address, big, cite, code, del, dfn, em,
img, ins, kbd, q, s, samp, small, strike, strong, sub, sup, tt, var, b, u, i,
center, dl, dt, dd, ol, ul, li, fieldset, form, label, legend, table, caption,
tbody, tfoot, thead, tr, th, td, article, aside, canvas, details, embed,
figure, figcaption, footer, header, hgroup, menu, nav, output, ruby, section,
summary, time, mark, audio, video {
  margin: 0;
  padding: 0;
  border: 0;
  font-size: 100%;
  font: inherit;
  vertical-align: baseline;
}
ol, ul { list-style: none; }
blockquote, q { quotes: none; }
table { border-collapse: collapse; border-spacing: 0; }
`
	preflight := `
*,
::before,
::after {
  box-sizing: border-box;
  border-width: 0;
  border-style: solid;
  border-color: #e5e7eb;
}
html {
  -webkit-text-size-adjust: 100%;
  tab-size: 4;
  font-feature-settings: normal;
  font-variation-settings: normal;
}
body { margin: 0; line-height: inherit; }
hr { height: 0; color: inherit; border-top-width: 1px; }
abbr:where([title]) { text-decoration: underline dotted; }
h1, h2, h3, h4, h5, h6 { font-size: inherit; font-weight: inherit; }
a { color: inherit; text-decoration: inherit; }
b, strong { font-weight: bolder; }
small { font-size: 80%; }
sub, sup { font-size: 75%; line-height: 0; position: relative; vertical-align: baseline; }
button, input, optgroup, select, textarea {
  font-family: inherit;
  font-size: 100%;
  font-weight: inherit;
  line-height: inherit;
  color: inherit;
  margin: 0;
  padding: 0;
}
button, select { text-transform: none; }
button,
input:where([type=button]),
input:where([type=reset]),
input:where([type=submit]) {
  -webkit-appearance: button;
  background-color: transparent;
  background-image: none;
}
@media (prefers-reduced-motion: reduce) {
  * { animation-duration: 0.01ms !important; scroll-behavior: auto !important; }
}
`
	mdn := `
:root { --gap: 12px; --accent: #bada55; }
.hero {
  display: grid;
  gap: var(--gap);
  grid-template-columns: repeat(auto-fill, minmax(min(20rem, 100%), 1fr));
  background:
    linear-gradient(135deg, rgb(255 0 0 / 0.3), hsl(120 50% 50% / 0.3)),
    url("tile.png") no-repeat center / cover;
}
@media (min-width: 700px) and (orientation: landscape) {
  .hero { grid-template-columns: repeat(3, 1fr); }
}
@media (width >= 900px) {
  .hero { gap: calc(var(--gap) + 12px); }
}
.alert {
  color: #fff;
  background-color: var(--accent, dodgerblue);
  border: 1px dashed currentColor;
  padding: 4px 8px;
}
`
	for _, sheet := range []string{reset, preflight, mdn} {
		sh := parseOK(t, sheet)
		// The sheets reach the cascade without dying even where they paint.
		_ = sh.Style("body", nil, StateNone, 800)
		_ = sh.Style("div", []string{"hero"}, StateNone, 400)
		_ = sh.Style("div", []string{"alert"}, StateHover, 800)
	}
}

// TestParseNeverHangs feeds the parser the shapes that used to spin it in a
// loop. Each case must come back, with a timer to prove it: this is the guard
// that a selector or at-rule always advances past itself. Run with
// `go test -timeout` as the outer net.
func TestParseNeverHangs(t *testing.T) {
	inputs := []string{
		`[disabled] { color: red; }`,
		`[type="text" i] { color: red; }`,
		`input[disabled]:hover { color: red; }`,
		`.a .b .c { color: red; }`,
		`p:nth-child(2n+1) { color: red; }`,
		`.!b { color: red; }`,
		`# { color: red; }`,
		`@media { body { color: red; } }`,
		`@media (prefers-color-scheme: dark) { body { color: white; } }`,
		`@media (min-width: 600mm) { span { color: red; } }`,
		`@media (min-width: 10px) { @media (max-width: 5px) { .x { color: red; } } }`,
		`svg|a { fill: orange; }`,
		`.\\2014 em { color: red; }`,
		`------ { color: red; }`,
		`body { background: url("a(b).png"); --x: var(--y; color: red; }`,
		`@unknown { 1 { : } }`,
	}
	for _, in := range inputs {
		timer := time.AfterFunc(10*time.Second, func() {
			panic("parser hung on: " + in)
		})
		_, _ = Parse(in)
		timer.Stop()
	}
}

func TestStyleImportantCascade(t *testing.T) {
	sh := parseOK(t, `
.a   { color: red; }
.a.b { color: blue; }
.a.b { color: green !important; }
`)
	if raw, ok := sh.Property("div", []string{"a", "b"}, "color", 800); !ok || raw != "green" {
		t.Errorf("Property of important vs higher-specificity normal = %q, %v; want green", raw, ok)
	}
	// With only the looser class, the important rule of .a.b cannot reach it.
	if raw, _ := sh.Property("div", []string{"a"}, "color", 800); raw != "red" {
		t.Errorf("Property without class b = %q, want red (the .a rule)", raw)
	}
	st := sh.Style("div", []string{"a", "b"}, StateNone, 800)
	if !st.Has("color") {
		t.Fatal("important rule should set color")
	}
	if st.Color != 0xFF008000 {
		t.Errorf("st.Color = %#x, want green", st.Color)
	}

	// Among important declarations themselves, the last one wins.
	sh = parseOK(t, `
.x { color: red !important; }
.x { color: blue; }
.x { color: green !important; }
`)
	if raw, _ := sh.Property("div", []string{"x"}, "color", 800); raw != "green" {
		t.Errorf("later !important should win, got %q", raw)
	}
}

func TestStyleResolvesCustomProperties(t *testing.T) {
	sh := parseOK(t, `
:root     { --brand: orange; }
.card     { --brand: blue; color: var(--brand, red); }
.fallback { background-color: var(--never, #111111); }
.chain    { --b: 20px; --a: var(--b); width: var(--a); }
`)
	st := sh.Style("div", []string{"card"}, StateNone, 800)
	if !st.Has("color") || st.Color != 0xFF0000FF {
		t.Errorf("var(--brand) resolved to %#x (has=%v), want blue", st.Color, st.Has("color"))
	}
	if st.Custom["--brand"] != "blue" {
		t.Errorf("Custom[--brand] = %q, want blue", st.Custom["--brand"])
	}

	// :root applies to the body element, carrying its own custom properties.
	if st := sh.Style("body", nil, StateNone, 800); st.Custom["--brand"] != "orange" {
		t.Errorf("body Custom[--brand] = %q, want orange", st.Custom["--brand"])
	}

	st = sh.Style("div", []string{"fallback"}, StateNone, 800)
	if !st.Has("background-color") || st.Background != 0xFF111111 {
		t.Errorf("var fallback should paint #111111, got %#x has=%v", st.Background, st.Has("background-color"))
	}

	st = sh.Style("div", []string{"chain"}, StateNone, 800)
	if !st.Has("width") || st.Width.Px(800) != 20 {
		t.Errorf("chained var() should resolve width to 20px, got %d has=%v", st.Width.Px(800), st.Has("width"))
	}
}

func TestApplyWarnsOnVendorAndUnknown(t *testing.T) {
	sh := parseOK(t, `.merk { -webkit-box-shadow: 0 1px 2px; bug-prop: nope; background-color: blue; }`)
	before := len(sh.Warn)
	st := sh.Style("div", []string{"merk"}, StateNone, 800)
	if !st.Has("background-color") {
		t.Error("the known declaration must still apply")
	}
	added := sh.Warn[before:]
	if !strings.Contains(strings.Join(added, "\n"), "vendor-prefixed") ||
		!strings.Contains(strings.Join(added, "\n"), "unknown property") {
		t.Errorf("style should warn about both, got %v", added)
	}
}

func TestSelectorUnsupportedNeverMatches(t *testing.T) {
	sh := parseOK(t, `
input[type="text"]  { background-color: blue; }
p > span            { color: red; }
::before            { background-color: green; }
div:hover span      { background-color: white; }
#widget             { color: purple; }
.ok                 { background-color: lime; }
`)
	for _, tag := range []string{"input", "span", "div", "p"} {
		if st := sh.Style(tag, nil, StateNone, 800); st.Has("background-color") || st.Has("color") {
			t.Errorf("%s: an unsupported selector must never apply, got %v/%v",
				tag, st.Has("background-color"), st.Has("color"))
		}
	}
	st := sh.Style("div", []string{"ok"}, StateNone, 800)
	if !st.Has("background-color") {
		t.Error(".ok must still apply")
	}
}

func TestSelectorCaseInsensitive(t *testing.T) {
	// Tags and class names are stored and matched in lowercase, so HTML-side
	// case differences cannot split a style away from its element.
	sh := parseOK(t, `
BUTTON        { background-color: red; }
.PrimaryCard  { color: blue; }
`)
	if st := sh.Style("button", nil, StateNone, 800); !st.Has("background-color") {
		t.Error("BUTTON selector should match the button tag")
	}
	if st := sh.Style("div", []string{"primarycard"}, StateNone, 800); !st.Has("color") {
		t.Error(".PrimaryCard selector should match a lowercase class")
	}
}

func TestSelectorRootIsBody(t *testing.T) {
	sh := parseOK(t, `:root { background-color: #0f0f0f; }`)
	st := sh.Style("body", nil, StateNone, 800)
	if !st.Has("background-color") || st.Background != 0xFF0F0F0F {
		t.Errorf(":root should style the body element, got %#x has=%v", st.Background, st.Has("background-color"))
	}
}

func TestSupportsAppliesNestedRules(t *testing.T) {
	sh := parseOK(t, `
@supports (display: flex) {
  .grid { width: 90%; }
  @media (min-width: 720px) { .wide { background-color: blue; } }
}
`)
	st := sh.Style("div", []string{"grid"}, StateNone, 800)
	if !st.Has("width") || st.Width.Px(800) != 720 {
		t.Errorf("90%% of 800 should be 720, got %d has=%v", st.Width.Px(800), st.Has("width"))
	}
	if st := sh.Style("div", []string{"wide"}, StateNone, 800); !st.Has("background-color") {
		t.Error("media inside @supports should apply when the width fits")
	}
	if st := sh.Style("div", []string{"wide"}, StateNone, 400); st.Has("background-color") {
		t.Error("media inside @supports must still gate on width")
	}
}

func TestMediaGrammar(t *testing.T) {
	sh := parseOK(t, `
@media (min-width: 720px) and (max-width: 900px) { .a { width: 10px; } }
@media (min-width: 1200px), (max-width: 300px) { .b { width: 20px; } }
@media not (min-width: 720px) { .c { width: 30px; } }
@media only screen { .d { width: 40px; } }
@media (min-width: 720px) or (max-width: 300px) { .e { width: 50px; } }
@media (orientation: portrait) { .f { width: 60px; } }
`)
	widthAt := func(w int, cls string) (int, bool) {
		st := sh.Style("div", []string{cls}, StateNone, w)
		if !st.Has("width") {
			return 0, false
		}
		return st.Width.Px(w), true
	}

	if v, ok := widthAt(800, "a"); !ok || v != 10 {
		t.Errorf("and-window at 800: %d, %v; want 10, true", v, ok)
	}
	if _, ok := widthAt(500, "a"); ok {
		t.Error("and-window should reject 500")
	}
	if _, ok := widthAt(1000, "a"); ok {
		t.Error("and-window should reject 1000")
	}
	if v, ok := widthAt(1500, "b"); !ok || v != 20 {
		t.Errorf("comma-high branch at 1500: %d, %v", v, ok)
	}
	if v, ok := widthAt(200, "b"); !ok || v != 20 {
		t.Errorf("comma-low branch at 200: %d, %v", v, ok)
	}
	if _, ok := widthAt(600, "b"); ok {
		t.Error("comma list should reject 600")
	}
	if v, ok := widthAt(500, "c"); !ok || v != 30 {
		t.Errorf("negated window at 500: %d, %v; want 30, true", v, ok)
	}
	if _, ok := widthAt(800, "c"); ok {
		t.Error("negated window should reject 800")
	}
	// A media type alone leaves no constraint: the rule applies everywhere.
	if _, ok := widthAt(500, "d"); !ok {
		t.Error("only screen should apply at any width")
	}
	if v, ok := widthAt(800, "e"); !ok || v != 50 {
		t.Errorf("or-high branch at 800: %d, %v", v, ok)
	}
	if v, ok := widthAt(200, "e"); !ok || v != 50 {
		t.Errorf("or-low branch at 200: %d, %v", v, ok)
	}
	if _, ok := widthAt(500, "e"); ok {
		t.Error("or list should reject 500")
	}
	// A condition the engine cannot measure is treated as satisfied, so the
	// rule is never dropped for lack of a sensor.
	if _, ok := widthAt(600, "f"); !ok {
		t.Error("unknown condition should not drop the rule")
	}
}

func TestMediaNestedApplies(t *testing.T) {
	sh := parseOK(t, `
@media (min-width: 720px) {
  @media (max-width: 900px) { .box { width: 10px; } }
}
`)
	if st := sh.Style("div", []string{"box"}, StateNone, 800); !st.Has("width") {
		t.Error("nested media inside the window should apply")
	}
	if st := sh.Style("div", []string{"box"}, StateNone, 500); st.Has("width") {
		t.Error("nested media must intersect with the outer window (500 is too narrow)")
	}
	if st := sh.Style("div", []string{"box"}, StateNone, 1000); st.Has("width") {
		t.Error("nested media must intersect with the outer window (1000 is too wide)")
	}
}

func TestSetStyleNeverFailsOnRealCSS(t *testing.T) {
	path := writeTemp(t, "real.css", `.primary { background-color: #3E63DD; }
button:hover { background-color: #2E52C0; }
input[type="text"] { padding: 4px 8px; }
@media (min-width: 600px) { .wide { width: 80%; } }
`)
	classes := NewTable()
	if err := classes.SetStyle(path); err != nil {
		t.Fatalf("SetStyle errored on real CSS: %v", err)
	}
	// The table must answer styles over the file, unsupported selectors and
	// all, instead of dying on the load.
	if st := classes.GetStyle(RoleButton, []string{"primary"}, StateNone, 800); !st.Has("background-color") {
		t.Error("primary button should take its background from the sheet")
	}
	if st := classes.GetStyle(RoleInput, nil, StateNone, 800); st.Has("background-color") {
		t.Error("input[type=text] is unsupported and must never match")
	}
}
