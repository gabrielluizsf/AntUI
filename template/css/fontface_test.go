package css

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// fixtureDir is where the fonts the tests load live, in the form an
// @font-face in a stylesheet with no file of its own reaches them by: an
// absolute path.
func fixtureDir(t *testing.T) string {
	t.Helper()
	dir, err := filepath.Abs(filepath.Join("..", "..", "testdata", "fonts"))
	if err != nil {
		t.Fatalf("Abs: %v", err)
	}
	return filepath.ToSlash(dir)
}

// faceSheet is a stylesheet declaring the fixture Roboto family: regular,
// bold and italic, each named by the weight and slant it carries.
func faceSheet(t *testing.T) *Sheet {
	t.Helper()
	dir := fixtureDir(t)
	return parseOK(t, `
@font-face {
	font-family: "Roboto";
	src: url("`+dir+`/Roboto-Regular.ttf") format("truetype");
	font-weight: 400;
	font-display: swap;
}
@font-face {
	font-family: "Roboto";
	src: url("`+dir+`/Roboto-Bold.ttf");
	font-weight: 700;
}
@font-face {
	font-family: "Roboto";
	src: url("`+dir+`/Roboto-Italic.ttf");
	font-style: italic;
}
`)
}

// TestFontFaceReadsAFile checks that a @font-face block opens its src at the
// reference size and answers to the family it declared, and that the block
// says nothing the engine has to report: font-display is read as silence.
func TestFontFaceReadsAFile(t *testing.T) {
	sh := faceSheet(t)
	if len(sh.Warn) != 0 {
		t.Fatalf("warnings = %v, want none", sh.Warn)
	}
	ff := sh.Font("roboto", FontWeightNormal, false)
	if ff == nil {
		t.Fatal("Font(roboto) = nil, want the face it declares")
	}
	if ff.Face() == nil {
		t.Fatal("Face() = nil, want the file it read")
	}
	if ff.Face().Size() != DefaultFontSize {
		t.Errorf("face size = %v, want the reference size %d", ff.Face().Size(), DefaultFontSize)
	}
	if ff.Bold() || ff.Slanted() {
		t.Errorf("bold = %v, slanted = %v; want both false for the regular face",
			ff.Bold(), ff.Slanted())
	}
}

// TestFontMatchTakesTheWeightAndSlantAskedFor walks the family the way CSS
// does: the weight decides between the regular and the bold file, a slant
// asks for the italic one, and a weight with no file of its own reaches out
// for the nearest one rather than dropping the text.
func TestFontMatchTakesTheWeightAndSlantAskedFor(t *testing.T) {
	sh := faceSheet(t)
	for _, c := range []struct {
		name    string
		weight  uint16
		slanted bool
		bold    bool
		lean    bool
	}{
		{"the regular face at its own weight", FontWeightNormal, false, false, false},
		{"the bold face at its own weight", FontWeightBold, false, true, false},
		{"a weight between them climbs to the bold", 600, false, true, false},
		{"a light weight falls to the regular", 300, false, false, false},
		{"an italic style takes the italic file", FontWeightNormal, true, false, true},
		{"a bold italic style takes the italic file, since there is no bold italic", FontWeightBold, true, false, true},
	} {
		ff := sh.Font("roboto", c.weight, c.slanted)
		if ff == nil {
			t.Errorf("%s: Font = nil", c.name)
			continue
		}
		if ff.Bold() != c.bold || ff.Slanted() != c.lean {
			t.Errorf("%s: bold = %v, slanted = %v; want %v, %v",
				c.name, ff.Bold(), ff.Slanted(), c.bold, c.lean)
		}
	}
}

// TestFontTakesTheFirstNameItFinds checks the font-family list itself: names
// are walked left to right, a name with no face in the sheet is passed over,
// and a family no name holds draws in the face the canvas already had.
func TestFontTakesTheFirstNameItFinds(t *testing.T) {
	sh := faceSheet(t)
	if ff := sh.Font("serif, roboto", FontWeightNormal, false); ff == nil {
		t.Error("Font(serif, roboto) = nil, want the face behind the name with none")
	}
	if ff := sh.Font("no-such-family", FontWeightNormal, false); ff != nil {
		t.Error("Font(no-such-family) returned a face, want nil")
	}
	if ff := sh.Font("", FontWeightNormal, false); ff != nil {
		t.Error("Font(\"\") returned a face, want nil — no family was asked for")
	}
}

// TestFontFamilyIsReadAndInherited checks the declaration itself: a list is
// kept whole, lowercased and unquoted the way the @font-face names it, and
// a widget that never mentioned one draws with the family its body set.
func TestFontFamilyIsReadAndInherited(t *testing.T) {
	sh := parseOK(t, `
		body  { font-family: "Roboto", "DejaVu Sans Mono", sans-serif; }
		button { color: red; }
	`)
	st := sh.Style("button", nil, StateNone, 800)
	if want := "roboto, dejavu sans mono, sans-serif"; st.FontFamily != want {
		t.Errorf("font-family = %q, want %q", st.FontFamily, want)
	}
	sh2 := parseOK(t, `button { font-family: initial; }`)
	if got := sh2.Style("button", nil, StateNone, 800).FontFamily; got != "" {
		t.Errorf("font-family after initial = %q, want empty", got)
	}
}

// TestFontFaceThatWillNotOpenIsReported checks the reasons a @font-face
// leaves a family with no face: a file that is not there, a source this
// engine does not read, a name it cannot match. Each is a warning, and none
// of them is a compile error.
func TestFontFaceThatWillNotOpenIsReported(t *testing.T) {
	dir := fixtureDir(t)
	for _, c := range []struct {
		name   string
		sheet  string
		warns  string
		family string
	}{
		{
			name:   "a file that is not there",
			sheet:  `@font-face { font-family: "Roboto"; src: url("` + dir + `/nope.ttf"); }`,
			warns:  "nope.ttf",
			family: "roboto",
		},
		{
			name:   "a woff file, which is a container this engine does not unpack",
			sheet:  `@font-face { font-family: "Roboto"; src: url("r.woff2") format("woff2"); }`,
			warns:  "woff is not read",
			family: "roboto",
		},
		{
			name:   "a local() source, which is a font this machine already has",
			sheet:  `@font-face { font-family: "Roboto"; src: local("Roboto"); }`,
			warns:  "local()",
			family: "roboto",
		},
		{
			name:   "a source that names neither function",
			sheet:  `@font-face { font-family: "Roboto"; src: roboto.ttf; }`,
			warns:  "not a url() or local()",
			family: "roboto",
		},
		{
			name:   "a declaration the block has no use for",
			sheet:  `@font-face { font-family: "Roboto"; src: url("x.ttf"); unicode-range: U+0-7F; }`,
			warns:  "unicode-range",
			family: "roboto",
		},
		{
			name:   "a block with no family to answer to",
			sheet:  `@font-face { src: url("x.ttf"); }`,
			warns:  "no font-family",
			family: "",
		},
		{
			name:   "a family with no src to read",
			sheet:  `@font-face { font-family: "Roboto"; }`,
			warns:  "no src",
			family: "roboto",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			sh, err := Parse(c.sheet)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			joined := strings.Join(sh.Warn, "\n")
			if !strings.Contains(joined, c.warns) {
				t.Errorf("warnings = %q, want one about %q", sh.Warn, c.warns)
			}
			if ff := sh.Font(c.family, FontWeightNormal, false); ff != nil {
				t.Errorf("Font(%q) = %v, want no face", c.family, ff)
			}
		})
	}
}

// TestFontFaceReadsItsFileBesideTheSheet checks that a @font-face written in
// a file opens a relative src beside that file, the way an @import does.
func TestFontFaceReadsItsFileBesideTheSheet(t *testing.T) {
	dir := t.TempDir()
	body, err := os.ReadFile(filepath.Join(fixtureDir(t), "Roboto-Regular.ttf"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "Roboto-Regular.ttf"), body, 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	main := filepath.Join(dir, "app.css")
	if err := os.WriteFile(main, []byte(
		`@font-face { font-family: "Roboto"; src: url("Roboto-Regular.ttf"); }`+"\n"+
			`.a { font-family: Roboto; }`), 0o644); err != nil {
		t.Fatalf("write sheet: %v", err)
	}
	sh, err := ParseFile(main)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(sh.Warn) != 0 {
		t.Fatalf("warnings = %v, want none", sh.Warn)
	}
	if sh.Font("roboto", FontWeightNormal, false) == nil {
		t.Error("Font(roboto) = nil, want the face read beside the sheet")
	}
	if got := sh.Style("div", []string{"a"}, StateNone, 800).FontFamily; got != "roboto" {
		t.Errorf("font-family = %q, want roboto", got)
	}
}

// fixtureFace is one of the fixture files read at the reference size, for
// saying what a face drawn out of it should measure.
func fixtureFace(t *testing.T, name string) *canvas.Face {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(fixtureDir(t), name))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	f, err := canvas.ParseFace(data, float64(DefaultFontSize))
	if err != nil {
		t.Fatalf("ParseFace(%q): %v", name, err)
	}
	return f
}

// TestFontFamilyHandsRunesDownTheList checks that a font-family list is not
// only a chooser of one file: the faces of the names behind the first are
// chained after it, so a rune one font has and another does not still comes
// out in the right shape — and the answer is kept, so a frame asks once.
func TestFontFamilyHandsRunesDownTheList(t *testing.T) {
	pin := canvas.DefaultFace()
	canvas.SetDefaultFace(nil) // the built-in, so the size the faces take is known
	defer canvas.SetDefaultFace(pin)

	dir := fixtureDir(t)
	sh := parseOK(t, `
@font-face { font-family: "Roboto"; src: url("`+dir+`/Roboto-Regular.ttf"); }
@font-face { font-family: "DejaVu Sans Mono"; src: url("`+dir+`/DejaVuSansMono.ttf"); }
`)
	roboto := fixtureFace(t, "Roboto-Regular.ttf")
	mono := fixtureFace(t, "DejaVuSansMono.ttf")
	if roboto.Width("→") == mono.Width("→") {
		t.Fatal("the two fixtures measure an arrow alike — this test needs a rune only one of them draws")
	}

	alone := sh.Font("roboto", FontWeightNormal, false)
	list := sh.Font("roboto, dejavu sans mono", FontWeightNormal, false)
	if alone == nil || list == nil {
		t.Fatalf("Font = %v and %v, want a face for both questions", alone, list)
	}
	if got, want := list.Face().Width("A"), roboto.Width("A"); got != want {
		t.Errorf("a rune the first font has measured %d, want that font's %d", got, want)
	}
	if got, want := list.Face().Width("→"), mono.Width("→"); got != want {
		t.Errorf("a rune only the second font has measured %d, want that font's %d", got, want)
	}
	if got, want := alone.Face().Width("→"), roboto.Width("→"); got != want {
		t.Errorf("a list of one name drew from another font: %d, want %d", got, want)
	}
	if again := sh.Font("roboto, dejavu sans mono", FontWeightNormal, false); again != list {
		t.Error("the same question asked twice gave two answers")
	}
	if sh.Font("serif", FontWeightNormal, false) != nil {
		t.Error("a name the sheet holds no face for returned one")
	}
}
