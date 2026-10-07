package canvas

import "testing"

// AtSize asks for a size rather than for a multiple of the size a face already
// has. The built-in font is one shape cut eight by sixteen, and it answers a
// size by drawing that shape with its cells at the size asked for: a size it
// does not come in is not rounded to one it does, the sizes that are whole
// multiples of its own height are still the faces BuiltinScaled hands out, and
// the same size asked for twice gives the same face back.
func TestFaceAtSizeSizesTheBuiltInFontAtTheSizeAskedFor(t *testing.T) {
	f := BuiltinFace()
	// A size of nothing is nothing asked for, and the face asked of stays.
	if got := f.AtSize(0); got != f {
		t.Error("a size of zero gave a different face")
	}
	// The font's own height is the face it already is rather than a new one of
	// the same size.
	if got := f.AtSize(float64(FontHeight)); got != f {
		t.Errorf("the font's own height gave a face %d tall, want the one it has", got.Height())
	}
	// Three times the height is that shape drawn three times over, and asking
	// again for the same size hands back the face that was made for it.
	three := f.AtSize(float64(3 * FontHeight))
	if three == f || three.Height() != 3*FontHeight {
		t.Fatalf("three times the height gave a face %d tall, want %d", three.Height(), 3*FontHeight)
	}
	if got := f.AtSize(float64(3 * FontHeight)); got != three {
		t.Error("the same size asked for twice gave two faces")
	}
	// A size the shape was not cut at comes out at that size, and everything
	// that measures follows it: twelve pixels of writing is twelve high, not
	// sixteen rounded down to the nearest step.
	twelve := f.AtSize(12)
	if twelve == f || twelve.Height() != 12 {
		t.Fatalf("twelve pixels gave a face %d tall, want 12", twelve.Height())
	}
	if got := twelve.Ascent(); got != 9 {
		t.Errorf("twelve pixels came with an ascent of %d, want 9", got)
	}
	if got := twelve.Width("hello"); got != 30 {
		t.Errorf("five letters at twelve pixels are %d wide, want 30", got)
	}
	if got := f.AtSize(12); got != twelve {
		t.Error("the same size asked for twice gave two faces")
	}
	// A face of nothing to size asks for the built-in rather than answering
	// with nothing to draw with.
	if got := (*Face)(nil).AtSize(float64(FontHeight)); got == nil || got.Height() != FontHeight {
		t.Errorf("sizing no face gave %v, want the built-in at its own height", got)
	}
}

// And the drawing follows the measuring: a size the built-in font does not
// come in writes the letters at that size rather than at the nearest one it
// does, which is the whole point of asking for a size.
func TestBuiltinWritesAtTheSizeItWasGiven(t *testing.T) {
	const baseline = 30
	twelve, own := BuiltinFace().AtSize(12), BuiltinFace()

	got := drawOne(t, twelve, baseline, "H")
	want := drawOne(t, own, baseline, "H")
	if got.Width == 0 || want.Width == 0 {
		t.Fatalf("one of the two painted nothing: %v and %v", got, want)
	}
	// At its own height the shape spans the whole sixteen pixels, at twelve it
	// has to end inside that.
	if got.Height >= want.Height {
		t.Errorf("twelve pixels painted %d high where the font's own size paints %d, want less",
			got.Height, want.Height)
	}
	if got.Height == 0 || got.Height > 12 {
		t.Errorf("twelve pixels painted %d high, want one to twelve", got.Height)
	}
	if got.Y < baseline-twelve.Ascent() {
		t.Errorf("twelve pixels painted at y=%d, above the top of the face at %d",
			got.Y, baseline-twelve.Ascent())
	}
	// Measuring and drawing still agree at a size that is not a step.
	cv, err := NewCanvas(80, 40)
	if err != nil {
		t.Fatal(err)
	}
	if w := twelve.Draw(cv, 4, baseline, "hi", Black); w != twelve.Width("hi") {
		t.Errorf("drawing wrote %d wide where measuring says %d", w, twelve.Width("hi"))
	}
}

// drawOne writes one string in a face and hands back the box the writing came
// out in, empty when nothing was painted.
func drawOne(t *testing.T, f *Face, baseline int, text string) Area {
	t.Helper()
	cv, err := NewCanvas(60, 60)
	if err != nil {
		t.Fatal(err)
	}
	f.Draw(cv, 8, baseline, text, Black)
	var box Area
	seen := false
	for y := range cv.Height {
		for x := range cv.Width {
			if cv.At(x, y).A() == 0 {
				continue
			}
			if !seen {
				box, seen = Area{X: x, Y: y, Width: 1, Height: 1}, true
				continue
			}
			x0, y0 := min(box.X, x), min(box.Y, y)
			x1 := max(box.X+box.Width, x+1)
			y1 := max(box.Y+box.Height, y+1)
			box = Area{X: x0, Y: y0, Width: x1 - x0, Height: y1 - y0}
		}
	}
	if !seen {
		return Area{}
	}
	return box
}

// loadFace is a TrueType file read at a size, for the tests that need a font
// with glyphs the built-in one does not carry.
func loadFace(t *testing.T, path string, pixels float64) *Face {
	t.Helper()
	f, err := LoadFace(path, pixels)
	if err != nil {
		t.Fatalf("LoadFace(%q): %v", path, err)
	}
	return f
}

// TestFaceHandsRunesItLacksToTheFaceBehindIt checks a fallback chain: a rune
// the face has is its own, a rune it does not comes from the face behind it,
// and a chain written round in a circle ends rather than running forever.
func TestFaceHandsRunesItLacksToTheFaceBehindIt(t *testing.T) {
	const dir = "../testdata/fonts/"
	roboto := loadFace(t, dir+"Roboto-Regular.ttf", 16)
	mono := loadFace(t, dir+"DejaVuSansMono.ttf", 16)
	if roboto.file.Cmap.Glyph('→') != 0 {
		t.Fatal("Roboto has an arrow — this test needs a rune it does not")
	}
	if mono.file.Cmap.Glyph('→') == 0 {
		t.Fatal("DejaVu Sans Mono has no arrow to fall back to")
	}

	chained := roboto.WithFallback(mono)
	if got, want := chained.Width("→"), mono.Width("→"); got != want {
		t.Errorf("arrow width = %d, want the fallback's %d", got, want)
	}
	if got, want := chained.Width("A"), roboto.Width("A"); got != want {
		t.Errorf("width of a rune the face has = %d, want its own %d", got, want)
	}
	if chained.glyph('→') != mono.glyph('→') {
		t.Error("the arrow was not drawn from the face behind it")
	}
	// The chain rides along to a size the face is asked for, so a scaled
	// face still falls back at the size it is drawn at.
	if got, want := chained.Scaled(2).Width("→"), mono.Scaled(2).Width("→"); got != want {
		t.Errorf("arrow width at twice the size = %d, want the fallback's %d", got, want)
	}
	// A chain pointing at itself has no face with the rune in it, and so
	// draws what it always did rather than walking round forever.
	round := roboto.WithFallback(roboto)
	if got, want := round.Width("→"), roboto.Width("→"); got != want {
		t.Errorf("arrow in a round chain = %d, want the face's own %d", got, want)
	}
	// Asking for no fallback hands the face back as it is.
	if chained := roboto.WithFallback(nil); chained != roboto {
		t.Error("a face given no fallback was replaced")
	}
}
