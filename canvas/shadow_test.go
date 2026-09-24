package canvas

import "testing"

func TestOuterShadowSitsBehindTheBox(t *testing.T) {
	cv, _ := NewCanvas(40, 40)
	cv.BoxShadow(15, 15, 10, 10, 0, 0, 0, 0, 0, 0, Black, false)
	if a := cv.At(20, 20).A(); a != 255 {
		t.Fatalf("inside the shadow should be opaque: %d", a)
	}
	if a := cv.At(5, 5).A(); a != 0 {
		t.Fatalf("far outside should be untouched: %d", a)
	}
}

func TestOuterShadowOffsetAndSpread(t *testing.T) {
	cv, _ := NewCanvas(40, 40)
	cv.BoxShadow(15, 15, 10, 10, 0, 0, 0, 0, 0, 3, Black, false)
	if a := cv.At(13, 20).A(); a == 0 {
		t.Fatal("spread should grow the shadow past the box edge")
	}
	if a := cv.At(5, 5).A(); a != 0 {
		t.Fatal("shadow should not reach the far corner")
	}
}

func TestOuterShadowOffset(t *testing.T) {
	cv, _ := NewCanvas(40, 40)
	cv.BoxShadow(15, 15, 10, 10, 0, 0, 5, 0, 0, 0, Black, false)
	if a := cv.At(27, 20).A(); a != 255 {
		t.Fatal("offset shadow should cover the shifted area")
	}
	if a := cv.At(16, 20).A(); a != 0 {
		t.Fatal("offset shadow should leave the old area")
	}
}

func TestInsetShadowHugsTheEdges(t *testing.T) {
	cv, _ := NewCanvas(40, 40)
	cv.BoxShadow(10, 10, 20, 20, 0, 0, 0, 0, 4, 0, Black, true)
	edge := cv.At(10, 20).A()
	centre := cv.At(20, 20).A()
	if edge == 0 {
		t.Fatal("inset shadow should darken the edge")
	}
	if centre != 0 {
		t.Fatalf("inset shadow should leave the centre clear: %d", centre)
	}
	if edge <= centre {
		t.Fatalf("edge %d should be stronger than centre %d", edge, centre)
	}
}

func TestShadowAlphaScalesWithColour(t *testing.T) {
	cv, _ := NewCanvas(20, 20)
	cv.FillRect(0, 0, 20, 20, White)
	cv.BoxShadow(5, 5, 10, 10, 0, 0, 0, 0, 0, 0, Fade(Black, 128), false)
	got := cv.At(10, 10)
	if got.R() < 120 || got.R() > 136 {
		t.Fatalf("half-transparent black over white should be mid grey: %v", got)
	}
}

func TestBlurredShadowStaysSoftOnEmptyCanvas(t *testing.T) {
	cv, _ := NewCanvas(60, 60)
	cv.BoxShadow(10, 10, 20, 20, 0, 0, 0, 0, 8, 0, Fade(Black, 255), false)
	// Blit blending for an opaque framebuffer forces every touched pixel
	// opaque, which would turn the whole padded layer into a solid slab —
	// the shadow "bigger than the box". The far canvas edge sits beyond any
	// genuine glow and must stay untouched.
	if a := cv.At(0, 30).A(); a > 20 {
		t.Fatalf("the far edge should stay nearly transparent, got alpha %d", a)
	}
	if a := cv.At(2, 30).A(); a > 120 {
		t.Fatalf("the glow should fade well short of the canvas edge, got alpha %d", a)
	}
	if a := cv.At(15, 20).A(); a < 150 {
		t.Fatalf("the shadow interior should stay strong, got alpha %d", a)
	}
}
