package antui

import (
	"strings"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/font"
)

// systemFace is the face these tests draw with, or a skip. A machine with no
// fonts is a real machine — a build container is usually one — and the right
// answer there is that these tests do not run, not that they fail.
func systemFace(t *testing.T, pixels float64) *canvas.Face {
	t.Helper()
	f, err := SystemFace(pixels)
	if err != nil {
		t.Skipf("no system font here: %v", err)
	}
	return f
}

// The font the system named has to be readable and make sense, or everything
// below it is measuring nonsense.
func TestSystemFontIsAFont(t *testing.T) {
	path, data, err := SystemFont()
	if err != nil {
		t.Skipf("no system font here: %v", err)
	}
	t.Logf("the system's font is %s (%d bytes)", path, len(data))

	f, err := font.ParseTTF(data)
	if err != nil {
		t.Fatalf("the chosen font does not parse: %v", err)
	}
	if f.UnitsPerEm < 16 || f.UnitsPerEm > 16384 {
		t.Errorf("unitsPerEm is %d", f.UnitsPerEm)
	}
	if f.NumGlyphs < 32 {
		t.Errorf("%d glyphs is not a text font", f.NumGlyphs)
	}
	if f.Ascent <= 0 || f.Descent >= 0 {
		t.Errorf("ascent %d, descent %d — a descent is negative", f.Ascent, f.Descent)
	}
	// Every text font has these, and a cmap that cannot find them is one
	// this library would draw as a row of empty boxes.
	for _, r := range "AZaz09 ." {
		if g := f.Cmap.Glyph(r); g == 0 && r != ' ' {
			t.Errorf("the cmap has no glyph for %q", r)
		}
	}
	if a := f.Advance(f.Cmap.Glyph('M')); a <= 0 {
		t.Errorf("M advances %d units", a)
	}
}

// A face's numbers have to be consistent with each other, or a layout built
// on them comes apart.
func TestFaceMetrics(t *testing.T) {
	f := systemFace(t, 16)
	if f.Size() != 16 {
		t.Errorf("Size = %v", f.Size())
	}
	if f.Ascent() <= 0 || f.Descent() < 0 {
		t.Errorf("ascent %d, descent %d", f.Ascent(), f.Descent())
	}
	if f.Height() < f.Ascent()+f.Descent() {
		t.Errorf("a line is %d tall and the glyphs are %d",
			f.Height(), f.Ascent()+f.Descent())
	}
	if f.Fixed() {
		t.Error("a system font is not fixed width")
	}
	// The built-in is, which is what lets anything lay out by counting.
	if !canvas.BuiltinFace().Fixed() {
		t.Error("the built-in font is fixed width")
	}
}

// Width has to grow with the string and with the size, and a newline has to
// measure the widest line rather than the total.
func TestFaceWidth(t *testing.T) {
	f := systemFace(t, 16)
	if w := f.Width(""); w != 0 {
		t.Errorf("the empty string is %d wide", w)
	}
	short, long := f.Width("i"), f.Width("iiii")
	if long <= short {
		t.Errorf("four i's (%d) are not wider than one (%d)", long, short)
	}
	if narrow, wide := f.Width("i"), f.Width("M"); wide <= narrow {
		t.Errorf("M (%d) is not wider than i (%d) in a proportional font", wide, narrow)
	}
	if small, big := f.Width("Hello"), systemFace(t, 32).Width("Hello"); big <= small {
		t.Errorf("Hello is %d at 16px and %d at 32px", small, big)
	}
	// Two lines measure the wider one, not the two of them end to end.
	one, two := f.Width("wwwwwwww"), f.Width("i\nwwwwwwww")
	if two != one {
		t.Errorf("two lines measured %d, the widest alone is %d", two, one)
	}
}

// Drawing has to put ink where Width says it will, or text runs over what is
// beside it.
func TestFaceDrawsWithinItsWidth(t *testing.T) {
	f := systemFace(t, 20)
	const text = "Handgloves"
	width := f.Width(text)

	cv, err := canvas.NewCanvas(width+40, 40)
	if err != nil {
		t.Fatal(err)
	}
	cv.Clear(canvas.White)
	if got := f.Draw(cv, 10, 10+f.Ascent(), text, canvas.Black); got != width {
		t.Errorf("Draw returned %d and Width says %d", got, width)
	}

	// Where the ink actually is.
	left, right := cv.Width, -1
	ink := 0
	for y := range cv.Height {
		for x := range cv.Width {
			if cv.At(x, y) != canvas.White {
				ink++
				left, right = min(left, x), max(right, x)
			}
		}
	}
	if ink == 0 {
		t.Fatal("nothing was drawn")
	}
	if left < 10-2 {
		t.Errorf("ink starts at %d, left of the pen at 10", left)
	}
	if right > 10+width+2 {
		t.Errorf("ink ends at %d, past the pen (10) plus the width (%d)", right, width)
	}
}

// An accented letter is a composite glyph — a letter and an accent, placed —
// and it is a different path through the parser from a plain one.
func TestFaceDrawsAccents(t *testing.T) {
	f := systemFace(t, 24)
	for _, pair := range []struct{ plain, accented string }{
		{"a", "ã"}, {"c", "ç"}, {"o", "ô"}, {"u", "ü"}, {"n", "ñ"},
	} {
		plain := inkOf(t, f, pair.plain)
		accented := inkOf(t, f, pair.accented)
		if accented <= plain {
			t.Errorf("%q has %d lit pixels and %q has %d — the accent is missing",
				pair.accented, accented, pair.plain, plain)
		}
	}
}

// inkOf is how many pixels one string lights up.
func inkOf(t *testing.T, f *canvas.Face, text string) int {
	t.Helper()
	cv, err := canvas.NewCanvas(f.Width(text)+20, f.Height()*2)
	if err != nil {
		t.Fatal(err)
	}
	cv.Clear(canvas.White)
	f.Draw(cv, 4, f.Height(), text, canvas.Black)
	n := 0
	for y := range cv.Height {
		for x := range cv.Width {
			if cv.At(x, y) != canvas.White {
				n++
			}
		}
	}
	return n
}

// Nothing changes for a program that never mentions faces. Two hundred call
// sites draw through Text and TextWidth, and they have to keep meaning the
// 8x16 font until somebody says otherwise.
func TestBuiltinIsStillTheDefault(t *testing.T) {
	if canvas.DefaultFace() != canvas.BuiltinFace() {
		t.Fatal("something set a default face")
	}
	if w := canvas.TextWidth("hello"); w != 5*canvas.FontWidth {
		t.Errorf("TextWidth said %d, and five cells is %d", w, 5*canvas.FontWidth)
	}
	if h := canvas.TextHeight(); h != canvas.FontHeight {
		t.Errorf("TextHeight said %d, and the font is %d", h, canvas.FontHeight)
	}
}

// And when somebody does say otherwise, measuring follows drawing — which is
// the whole reason the default is process-wide.
func TestDefaultFaceMovesMeasuringToo(t *testing.T) {
	f := systemFace(t, 18)
	canvas.SetDefaultFace(f)
	defer canvas.SetDefaultFace(nil)

	if canvas.TextWidth("hello") != f.Width("hello") {
		t.Errorf("TextWidth said %d and the face says %d",
			canvas.TextWidth("hello"), f.Width("hello"))
	}
	if canvas.TextHeight() != f.Height() {
		t.Errorf("TextHeight said %d and the face says %d", canvas.TextHeight(), f.Height())
	}
	// A proportional font is not the bitmap's arithmetic.
	if canvas.TextWidth("iiiii") == canvas.TextWidth("MMMMM") {
		t.Error("five i's measured the same as five M's")
	}
}

// A canvas may keep its own face whatever the program's is, which is what a
// game's HUD wants while the menus are in the system font.
func TestCanvasFaceOverridesTheDefault(t *testing.T) {
	f := systemFace(t, 18)
	canvas.SetDefaultFace(f)
	defer canvas.SetDefaultFace(nil)

	cv, err := canvas.NewCanvas(200, 60)
	if err != nil {
		t.Fatal(err)
	}
	if cv.Face() != f {
		t.Error("a canvas with no face of its own is not using the default")
	}
	cv.SetFace(canvas.BuiltinFace())
	if cv.Face() != canvas.BuiltinFace() {
		t.Error("SetFace did not take")
	}
	cv.Clear(canvas.White)
	if w := cv.Text(0, 0, "hello", canvas.Black); w != 5*canvas.FontWidth {
		t.Errorf("the canvas drew %d wide and the built-in is %d", w, 5*canvas.FontWidth)
	}
	cv.SetFace(nil)
	if cv.Face() != f {
		t.Error("nil did not go back to the default")
	}
}

// Bytes that are not a font are an error and never a crash: a font file is
// something a program may be handed.
func TestParseFaceRefusesRubbish(t *testing.T) {
	for name, data := range map[string][]byte{
		"empty":     {},
		"short":     []byte("no"),
		"wrong tag": append([]byte("junk"), make([]byte, 200)...),
		"truncated": append([]byte{0, 1, 0, 0, 0, 4}, make([]byte, 10)...),
	} {
		if _, err := canvas.ParseFace(data, 16); err == nil {
			t.Errorf("%s was accepted as a font", name)
		}
	}
	// And a real font at an impossible size.
	_, data, err := SystemFont()
	if err == nil {
		if _, err := canvas.ParseFace(data, 0); err == nil {
			t.Error("a size of zero was accepted")
		}
	}
}

// A CFF font has outlines this library cannot read, and saying so by name is
// better than drawing nothing.
func TestOTTOIsRefusedByName(t *testing.T) {
	otto := append([]byte("OTTO"), make([]byte, 200)...)
	_, err := canvas.ParseFace(otto, 16)
	if err == nil {
		t.Fatal("a CFF font was accepted")
	}
	if !strings.Contains(err.Error(), "CFF") {
		t.Errorf("the error does not say what is wrong: %v", err)
	}
}

// The built-in font at a scale is still the built-in font, and everything
// that measures has to follow — on a 420dpi phone sixteen pixels is about a
// millimetre and a half, so this is what makes "the built-in font" a real
// choice there rather than a joke.
func TestBuiltinScaled(t *testing.T) {
	if canvas.BuiltinScaled(1) != canvas.BuiltinFace() {
		t.Error("scale 1 is not the plain built-in")
	}
	if canvas.BuiltinScaled(0) != canvas.BuiltinFace() {
		t.Error("scale 0 is not the plain built-in")
	}
	big := canvas.BuiltinScaled(3)
	if !big.Fixed() {
		t.Error("a scaled bitmap font is still fixed width")
	}
	if got, want := big.Height(), canvas.FontHeight*3; got != want {
		t.Errorf("height %d, want %d", got, want)
	}
	if got, want := big.Width("hello"), 5*canvas.FontWidth*3; got != want {
		t.Errorf("width %d, want %d", got, want)
	}
	if canvas.BuiltinScaled(3) != big {
		t.Error("asking twice made two faces")
	}

	// And drawing agrees with measuring.
	canvas.SetDefaultFace(big)
	defer canvas.SetDefaultFace(nil)
	if canvas.TextWidth("hello") != big.Width("hello") {
		t.Errorf("TextWidth %d, face %d", canvas.TextWidth("hello"), big.Width("hello"))
	}
	if canvas.TextHeight() != canvas.FontHeight*3 {
		t.Errorf("TextHeight %d, want %d", canvas.TextHeight(), canvas.FontHeight*3)
	}

	cv, err := canvas.NewCanvas(200, 80)
	if err != nil {
		t.Fatal(err)
	}
	cv.Clear(canvas.White)
	if w := cv.Text(4, 4, "hi", canvas.Black); w != 2*canvas.FontWidth*3 {
		t.Errorf("drew %d wide, want %d", w, 2*canvas.FontWidth*3)
	}
	// The ink has to be inside the box it claimed.
	for y := range cv.Height {
		for x := range cv.Width {
			if cv.At(x, y) == canvas.White {
				continue
			}
			if x < 4 || x >= 4+2*canvas.FontWidth*3 || y < 4 || y >= 4+canvas.FontHeight*3 {
				t.Fatalf("ink at (%d,%d) is outside the %dx%d it claimed",
					x, y, 2*canvas.FontWidth*3, canvas.FontHeight*3)
			}
		}
	}
}