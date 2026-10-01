package canvas

import "testing"

func newLayer(t *testing.T, w, h int) *Canvas {
	t.Helper()
	cv, err := NewLayer(w, h)
	if err != nil {
		t.Fatalf("NewLayer(%d, %d): %v", w, h, err)
	}
	return cv
}

// A layer is the one canvas that has nothing behind it yet. Everything else is
// a surface, and a shape with a soft edge painted on a surface is simply opaque
// there, because the surface behind it is. A layer has to answer with the
// transparency the shape actually had, or the edge of anything drawn into one —
// an icon, a shadow, a stroke — comes out as a hard ring of the wrong colour.
func TestLayerKeepsTheAlphaOfWhatIsDrawnOnIt(t *testing.T) {
	// A half-transparent red drawn on a surface is opaque red: there is nothing
	// behind the surface to see. On a layer it is still half transparent.
	surface := newCanvas(t, 4, 4)
	surface.Pixel(1, 1, RGBA(255, 0, 0, 128))
	if got := surface.At(1, 1).A(); got != 255 {
		t.Errorf("a translucent pixel on a surface came out at alpha %d, want 255", got)
	}

	layer := newLayer(t, 4, 4)
	layer.Pixel(1, 1, RGBA(255, 0, 0, 128))
	if got := layer.At(1, 1).A(); got != 128 {
		t.Errorf("a translucent pixel on a layer came out at alpha %d, want 128", got)
	}
	// The colour itself is kept as it was drawn, unweighted by the alpha, so a
	// layer can be scaled down without its colour going dark on the way.
	if r := layer.At(1, 1).R(); r != 255 {
		t.Errorf("the layer's pixel is red %d, want 255", r)
	}
	// A pixel nothing was drawn on is still nothing at all.
	if got := layer.At(3, 3); got.A() != 0 {
		t.Errorf("an untouched pixel of a layer came out at alpha %d, want 0", got.A())
	}
}

// Drawing the same translucent colour twice on a layer is the case the surface
// cannot express: the two draws add up, and the result is more opaque than
// either of them without being fully opaque.
func TestLayerAddsUpWhereItIsDrawnTwice(t *testing.T) {
	red := RGBA(255, 0, 0, 128)

	once := newLayer(t, 2, 2)
	once.Pixel(0, 0, red)

	twice := newLayer(t, 2, 2)
	twice.Pixel(0, 0, red)
	twice.Pixel(0, 0, red)

	first, second := once.At(0, 0).A(), twice.At(0, 0).A()
	if second <= first {
		t.Errorf("the second draw left the layer at alpha %d, want more than the %d of one", second, first)
	}
	if second == 255 {
		t.Errorf("two half-transparent draws came out fully opaque at alpha %d, want the gap to show", second)
	}
	if want := BlendOver(once.At(0, 0), red); twice.At(0, 0) != want {
		t.Errorf("two draws came out %#08x, want the second over the first: %#08x", twice.At(0, 0), want)
	}

	// A surface has nowhere to put that second half, so it saturates instead.
	surface := newCanvas(t, 2, 2)
	surface.Pixel(0, 0, red)
	surface.Pixel(0, 0, red)
	if got := surface.At(0, 0).A(); got != 255 {
		t.Errorf("two draws on a surface came out at alpha %d, want 255", got)
	}
}

// The reason a layer exists is a shape whose edge falls between two pixels. A
// disc covers the pixels inside it whole and the ones outside not at all, and
// the ones along its edge in part; on a layer that part survives.
func TestLayerKeepsTheSoftEdgeOfAShape(t *testing.T) {
	layer := newLayer(t, 20, 20)
	layer.FillCircle(10, 10, 7, White)

	// A disc of an even radius has no edge pixel: either a pixel is inside the
	// radius or it is not, because the radius lands between pixel centres. An
	// odd radius puts the circle's boundary through a column of pixel centres,
	// and that column is where the greys are.
	layer.FillCircle(10, 10, 6, White)

	var partial int
	for y := range layer.Height {
		for x := range layer.Width {
			if a := layer.At(x, y).A(); a > 0 && a < 255 {
				partial++
			}
		}
	}
	if partial == 0 {
		t.Error("an odd radius disc on a layer has no part-covered pixels, want a soft edge")
	}
}

// A layer is only worth having if it can be put back down: blitting one onto a
// surface has to show what was behind it, which is what its alpha is for.
func TestBlittingALayerShowsWhatIsBehindIt(t *testing.T) {
	back := newCanvas(t, 8, 8)
	back.FillRect(0, 0, 8, 8, Blue)

	half := newLayer(t, 8, 8)
	half.FillRect(0, 0, 8, 8, RGBA(255, 255, 255, 128))
	back.BlitOver(0, 0, half)

	// Blue is (0, 0, 255) and white is (255, 255, 255); half of each is a grey
	// that is more blue than red.
	at := back.At(4, 4)
	r, g, b := int(at.R()), int(at.G()), int(at.B())
	if r == 0 || g == 0 || b == 0 {
		t.Fatalf("the layer hid what was behind it: got rgb(%d, %d, %d), want a blend of both", r, g, b)
	}
	if r >= b || g >= b {
		t.Errorf("blending half white over blue gave rgb(%d, %d, %d), want the blue to still be there", r, g, b)
	}
	if got := back.At(4, 4).A(); got != 255 {
		t.Errorf("the surface is at alpha %d after the layer went down, want it solid", got)
	}
}

// Scaling a layer down has to average its colours weighted by what each pixel
// covers, or a thin bright line on a transparent background turns into a dark
// line on a lighter one.
func TestScalingALayerKeepsTheColourOfWhatItCovers(t *testing.T) {
	layer := newLayer(t, 4, 1)
	layer.Pixel(0, 0, White)
	// Three transparent pixels, one fully covered white one.

	scaled := Scaled(layer, 1, 1)
	at := scaled.At(0, 0)
	if got := at.A(); got != 255/4 {
		t.Errorf("a white pixel over three transparent ones came out at alpha %d, want a quarter of it", got)
	}
	if r, g, b := at.R(), at.G(), at.B(); r != 255 || g != 255 || b != 255 {
		t.Errorf("the scaled pixel is rgb(%d, %d, %d), want the white kept rather than mixed towards nothing", r, g, b)
	}
}
