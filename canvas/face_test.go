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
