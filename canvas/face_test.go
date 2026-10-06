package canvas

import "testing"

// AtSize asks for a size rather than for a multiple of the size a face already
// has, and the built-in font answers in steps of its own height: it is one
// shape drawn n times over, so it has the sizes those steps land on and the
// same size asked for twice gives the same face back.
func TestFaceAtSizeTakesTheBuiltInFontInItsOwnSteps(t *testing.T) {
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
	// A face of nothing to size asks for the built-in rather than answering
	// with nothing to draw with.
	if got := (*Face)(nil).AtSize(float64(FontHeight)); got == nil || got.Height() != FontHeight {
		t.Errorf("sizing no face gave %v, want the built-in at its own height", got)
	}
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
