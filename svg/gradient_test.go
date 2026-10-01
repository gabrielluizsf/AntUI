package svg

import (
	"strings"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// painted is a drawing read and painted at the size asked for, which is what
// every test here wants and none of them wants to write the two failures for.
func painted(t *testing.T, src string, w, h int) (*Image, *canvas.Canvas) {
	t.Helper()
	img, err := Parse(src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	cv := img.Render(w, h)
	if cv == nil {
		t.Fatal("the drawing did not paint")
	}
	return img, cv
}

// TestRenderPaintsALinearGradient runs the ramp along the shape it fills: red at
// the start of the line, blue at the far end, and every colour between them in
// the order they were written.
func TestRenderPaintsALinearGradient(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 10 10">
		<defs>
			<linearGradient id="g" x1="0" y1="0" x2="10" y2="0" gradientUnits="userSpaceOnUse">
				<stop offset="0" stop-color="#ff0000"/>
				<stop offset="1" stop-color="#0000ff"/>
			</linearGradient>
		</defs>
		<rect x="0" y="0" width="10" height="10" fill="url(#g)"/>
	</svg>`, 10, 10)

	var prevR, prevB int = 300, -1
	for _, x := range []int{0, 2, 4, 6, 8} {
		c := pixelAt(cv, x, 5)
		if int(c.R()) >= prevR {
			t.Errorf("at x=%d the red is %d, want it falling from left to right", x, c.R())
		}
		if int(c.B()) <= prevB {
			t.Errorf("at x=%d the blue is %d, want it rising from left to right", x, c.B())
		}
		prevR, prevB = int(c.R()), int(c.B())
	}
	if c := pixelAt(cv, 0, 5); c.R() < 200 || c.B() > 60 {
		t.Errorf("the start of the line is %#08x, want it red", uint32(c))
	}
	if c := pixelAt(cv, 8, 5); c.B() < 200 || c.R() > 60 {
		t.Errorf("the end of the line is %#08x, want it blue", uint32(c))
	}
}

// TestRenderPaintsAGradientAcrossTheShapeItFills is the point of fractions of
// the bounding box: the ramp follows the rectangle wherever the rectangle is,
// rather than being laid over the drawing from the corner of it.
func TestRenderPaintsAGradientAcrossTheShapeItFills(t *testing.T) {
	// The rectangle starts a third of the way across the drawing, so a gradient
	// painted against the drawing itself would leave it blue from one end.
	_, cv := painted(t, `<svg viewBox="0 0 30 10">
		<defs>
			<linearGradient id="g">
				<stop offset="0" stop-color="#ff0000"/>
				<stop offset="1" stop-color="#0000ff"/>
			</linearGradient>
		</defs>
		<rect x="10" y="0" width="10" height="10" fill="url(#g)"/>
	</svg>`, 30, 10)

	if c := pixelAt(cv, 10, 5); c.R() < 200 || c.B() > 60 {
		t.Errorf("the left edge of the rectangle is %#08x, want it red", uint32(c))
	}
	if c := pixelAt(cv, 19, 5); c.B() < 200 || c.R() > 60 {
		t.Errorf("the right edge of the rectangle is %#08x, want it blue", uint32(c))
	}
	// The drawing either side of the rectangle was never asked for a colour.
	if c := pixelAt(cv, 0, 5); c.A() != 0 {
		t.Errorf("the empty corner is %#08x, want it clear", uint32(c))
	}
}

// TestRenderPaintsAGradientInTheCoordinatesOfAShapeThatWasMoved is the same
// question asked of a transform: the gradient is written against where the shape
// is before anything moves it, and the transform carries the two of them
// together afterwards.
func TestRenderPaintsAGradientInTheCoordinatesOfAShapeThatWasMoved(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 20 10">
		<defs>
			<linearGradient id="g">
				<stop offset="0" stop-color="#ff0000"/>
				<stop offset="1" stop-color="#0000ff"/>
			</linearGradient>
		</defs>
		<rect x="0" y="0" width="10" height="10" transform="translate(10,0)" fill="url(#g)"/>
	</svg>`, 20, 10)

	if c := pixelAt(cv, 10, 5); c.R() < 200 || c.B() > 60 {
		t.Errorf("the moved rectangle starts at %#08x, want it red", uint32(c))
	}
	if c := pixelAt(cv, 19, 5); c.B() < 200 || c.R() > 60 {
		t.Errorf("the moved rectangle ends at %#08x, want it blue", uint32(c))
	}
}

// TestRenderPaintsARadialGradientFromItsCentreOut is the default shape of one:
// the first stop at the middle of the circle and the last stop at its edge, with
// the colours of the drawing either side of that edge left alone.
func TestRenderPaintsARadialGradientFromItsCentreOut(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 10 10">
		<defs>
			<radialGradient id="g">
				<stop offset="0" stop-color="#ff0000"/>
				<stop offset="1" stop-color="#0000ff"/>
			</radialGradient>
		</defs>
		<rect x="0" y="0" width="10" height="10" fill="url(#g)"/>
	</svg>`, 10, 10)

	if c := pixelAt(cv, 5, 5); c.R() < 200 || c.B() > 60 {
		t.Errorf("the middle of the circle is %#08x, want it red", uint32(c))
	}
	// The corner of the square is outside the circle entirely, and past the end
	// of a gradient with no spread method is where it stops.
	if c := pixelAt(cv, 0, 0); c.B() < 200 || c.R() > 60 {
		t.Errorf("the corner past the circle is %#08x, want it blue", uint32(c))
	}
}

// TestGradientTakesWhatItDoesNotSayFromTheOneItPointsAt is the `href`: the stops
// come from the gradient named, and everything this one says for itself is kept
// over them.
func TestGradientTakesWhatItDoesNotSayFromTheOneItPointsAt(t *testing.T) {
	img, cv := painted(t, `<svg viewBox="0 0 10 10">
		<defs>
			<linearGradient id="stops">
				<stop offset="0" stop-color="#ff0000"/>
				<stop offset="1" stop-color="#0000ff"/>
			</linearGradient>
			<radialGradient id="g" href="#stops"/>
		</defs>
		<rect x="0" y="0" width="10" height="10" fill="url(#g)"/>
	</svg>`, 10, 10)

	if len(img.Warnings()) != 0 {
		t.Errorf("the drawing said %v, want it read whole", img.Warnings())
	}
	// A radial gradient takes its circle from its own defaults, and only the
	// colours are the linear gradient's.
	if c := pixelAt(cv, 5, 5); c.R() < 200 {
		t.Errorf("the middle is %#08x, want the first stop of the pointed-at gradient", uint32(c))
	}
	if c := pixelAt(cv, 9, 9); c.B() < 200 {
		t.Errorf("the far corner is %#08x, want the last stop of the pointed-at gradient", uint32(c))
	}
}

// TestGradientPercentagesAreTakenAgainstTheUnitsTheyEndUpWith is the reason the
// numbers are held until the href has been followed: the units this gradient
// never mentioned are the ones it inherits, and a percentage of them is not a
// percentage of anything the same size.
func TestGradientPercentagesAreTakenAgainstTheUnitsTheyEndUpWith(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 10 10">
		<defs>
			<linearGradient id="base" gradientUnits="userSpaceOnUse" x1="0" x2="10">
				<stop offset="0" stop-color="#ff0000"/>
				<stop offset="1" stop-color="#0000ff"/>
			</linearGradient>
			<linearGradient id="g" href="#base" x2="100%"/>
		</defs>
		<rect x="0" y="0" width="10" height="10" fill="url(#g)"/>
	</svg>`, 10, 10)

	// A hundred percent of the drawing's own ten units is ten, which runs the
	// whole ramp across the rectangle. Read against the unit square it would come
	// out as a hundredth of the width and leave every pixel red.
	if c := pixelAt(cv, 8, 5); c.B() < 200 {
		t.Errorf("the right of the rectangle is %#08x, want the percentage to reach the far side", uint32(c))
	}
}

// TestRenderRepeatsAGradientPastItsEnds is `spreadMethod`, which says what a
// gradient does with the points past either end of its line rather than leaving
// them at the last colour.
func TestRenderRepeatsAGradientPastItsEnds(t *testing.T) {
	// The line runs halfway across the rectangle, so the second half of it is
	// the ramp again for a repeating gradient and the same ramp backwards for a
	// reflecting one.
	const src = `<svg viewBox="0 0 10 10">
		<defs>
			<linearGradient id="g" x1="0" x2="0.5" spreadMethod="%s">
				<stop offset="0" stop-color="#ff0000"/>
				<stop offset="1" stop-color="#0000ff"/>
			</linearGradient>
		</defs>
		<rect x="0" y="0" width="10" height="10" fill="url(#g)"/>
	</svg>`

	_, repeat := painted(t, substitute(src, "repeat"), 10, 10)
	if c := pixelAt(repeat, 0, 5); c.R() < 200 {
		t.Errorf("the start of the first run is %#08x, want it red", uint32(c))
	}
	if c := pixelAt(repeat, 4, 5); c.B() < 200 {
		t.Errorf("the end of the first run is %#08x, want it blue", uint32(c))
	}
	if c := pixelAt(repeat, 5, 5); c.R() < 200 {
		t.Errorf("the start of the second run is %#08x, want the ramp to have started again", uint32(c))
	}

	_, reflect := painted(t, substitute(src, "reflect"), 10, 10)
	if c := pixelAt(reflect, 5, 5); c.B() < 200 {
		t.Errorf("the reflection starts at %#08x, want it to pick up where the first run ended", uint32(c))
	}
	if c := pixelAt(reflect, 9, 5); c.R() < 200 {
		t.Errorf("the end of the reflection is %#08x, want it back at red", uint32(c))
	}
}

// substitute fills in the one part of a drawing a test changes about it, which
// keeps the two drawings around it the same in every other line.
func substitute(src, what string) string { return strings.ReplaceAll(src, "%s", what) }

// TestRenderAppliesAGradientTransform is the transform on the gradient itself,
// which moves the paint without moving the shape it paints.
func TestRenderAppliesAGradientTransform(t *testing.T) {
	const src = `<svg viewBox="0 0 10 10">
		<defs>
			<linearGradient id="g" x1="0" x2="1"%s>
				<stop offset="0" stop-color="#ff0000"/>
				<stop offset="1" stop-color="#0000ff"/>
			</linearGradient>
		</defs>
		<rect x="0" y="0" width="10" height="10" fill="url(#g)"/>
	</svg>`

	_, plain := painted(t, substitute(src, ""), 10, 10)
	_, moved := painted(t, substitute(src, ` gradientTransform="translate(0.5,0)"`), 10, 10)

	// The gradient is half a width to the right of the shape, so where the plain
	// one is finished the moved one has only come halfway.
	atPlain, atMoved := pixelAt(plain, 9, 5).B(), pixelAt(moved, 9, 5).B()
	if atMoved >= atPlain-60 {
		t.Errorf("the far end of the moved gradient is blue at %d, want it well short of the plain one at %d", atMoved, atPlain)
	}
}

// TestRenderPaintsAStrokeWithAGradient is the stroke asking the same question
// the fill does: the outline is a shape of its own and the gradient runs across
// it in the drawing's own coordinates.
func TestRenderPaintsAStrokeWithAGradient(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 10 10">
		<defs>
			<linearGradient id="g" gradientUnits="userSpaceOnUse" x1="0" x2="10">
				<stop offset="0" stop-color="#ff0000"/>
				<stop offset="1" stop-color="#0000ff"/>
			</linearGradient>
		</defs>
		<rect x="1" y="1" width="8" height="8" fill="none" stroke="url(#g)" stroke-width="2"/>
	</svg>`, 10, 10)

	if c := pixelAt(cv, 1, 5); c.R() < 200 || c.B() > 60 {
		t.Errorf("the left of the stroke is %#08x, want it red", uint32(c))
	}
	if c := pixelAt(cv, 8, 5); c.B() < 200 || c.R() > 60 {
		t.Errorf("the right of the stroke is %#08x, want it blue", uint32(c))
	}
	// The inside of the rectangle has no fill at all.
	if c := pixelAt(cv, 5, 5); c.A() != 0 {
		t.Errorf("the middle of the rectangle is %#08x, want it clear", uint32(c))
	}
}

// TestRenderAppliesTheOpacityOfEachStop is `stop-opacity`, which fades a colour
// into the layer rather than into the colour after it along the line.
func TestRenderAppliesTheOpacityOfEachStop(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 10 10">
		<defs>
			<linearGradient id="g" x1="0" x2="1">
				<stop offset="0" stop-color="#ff0000" stop-opacity="0.5"/>
				<stop offset="1" stop-color="#0000ff"/>
			</linearGradient>
		</defs>
		<rect x="0" y="0" width="10" height="10" fill="url(#g)"/>
	</svg>`, 10, 10)

	c := pixelAt(cv, 0, 5)
	if c.R() < 240 {
		t.Errorf("the half-transparent stop came out %#08x, want it still red", uint32(c))
	}
	if a := int(c.A()); a < 118 || a > 138 {
		t.Errorf("its alpha is %d, want about 128", a)
	}
}

// TestRenderPaintsAGradientInCurrentColor is a stop that names the colour of the
// widget rather than one of its own, which is only known when it is painted.
func TestRenderPaintsAGradientInCurrentColor(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 10 10">
		<defs>
			<linearGradient id="g" x1="0" x2="1">
				<stop offset="0" stop-color="currentColor"/>
				<stop offset="1" stop-color="#0000ff"/>
			</linearGradient>
		</defs>
		<rect x="0" y="0" width="10" height="10" fill="url(#g)"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	cv := img.RenderWith(10, 10, canvas.RGB(0, 255, 0))
	if c := pixelAt(cv, 0, 5); c.G() < 200 {
		t.Errorf("the start of the ramp is %#08x, want it in the colour it was painted with", uint32(c))
	}
}

// TestParseSaysNothingAboutADrawingMadeOfGradients is the test that a drawing
// full of the new tags is one this package reads whole rather than one it reads
// with a warning per line: the stops and the gradients are not shapes, and they
// are not shapes it fails to draw either.
func TestParseSaysNothingAboutADrawingMadeOfGradients(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 10 10">
		<defs>
			<linearGradient id="a">
				<stop offset="0" stop-color="red"/>
				<stop offset="50%" stop-color="gold"/>
				<stop offset="1" stop-color="blue"/>
			</linearGradient>
			<radialGradient id="b" href="#a"/>
		</defs>
		<rect x="0" y="0" width="5" height="10" fill="url(#a)"/>
		<rect x="5" y="0" width="5" height="10" fill="url(#b)"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(img.Warnings()) != 0 {
		t.Errorf("the drawing said %v, want nothing read wrong", img.Warnings())
	}
}

// TestParseWarnsAboutAGradientThatIsNotThere is a fill naming something the
// drawing does not have, with no second choice after it: the shape comes out
// unpainted and the drawing says so, once.
func TestParseWarnsAboutAGradientThatIsNotThere(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 10 10">
		<rect x="0" y="0" width="10" height="10" fill="url(#nowhere)"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got := len(img.Warnings()); got != 1 {
		t.Errorf("the drawing said %d things, want exactly the one missing gradient", got)
	}
	if n := img.Root.find("rect"); n == nil || n.Style.HasFill {
		t.Error("the rectangle was filled anyway")
	}
	if c := pixelAt(img.Render(10, 10), 5, 5); c.A() != 0 {
		t.Errorf("the rectangle came out %#08x, want it clear", uint32(c))
	}
}

// TestRenderFallsBackToTheColourAfterTheReference is the second choice a paint
// may carry: `url(#g) red` is the gradient where there is one and red where
// there is not, and a shape that brought a fallback with it has nothing to warn
// about.
func TestRenderFallsBackToTheColourAfterTheReference(t *testing.T) {
	img, cv := painted(t, `<svg viewBox="0 0 10 10">
		<rect x="0" y="0" width="10" height="10" fill="url(#nowhere) red"/>
	</svg>`, 10, 10)
	if len(img.Warnings()) != 0 {
		t.Errorf("the drawing said %v, want it happy with the colour it fell back on", img.Warnings())
	}
	if c := pixelAt(cv, 5, 5); c != canvas.RGB(255, 0, 0) {
		t.Errorf("the rectangle is %#08x, want it red", uint32(c))
	}
}

// TestRenderUsesTheGradientWhereThereIsOneAndTheColourWhereThereIsNot is the
// other half of the fallback: the reference is followed first, and the colour
// after it is only asked for when it is not there.
func TestRenderUsesTheGradientWhereThereIsOneAndTheColourWhereThereIsNot(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 10 10">
		<defs>
			<linearGradient id="g" x1="0" x2="1">
				<stop offset="0" stop-color="#ff0000"/>
				<stop offset="1" stop-color="#0000ff"/>
			</linearGradient>
		</defs>
		<rect x="0" y="0" width="10" height="10" fill="url(#g) yellow"/>
	</svg>`, 10, 10)
	if c := pixelAt(cv, 0, 5); c.R() < 200 {
		t.Errorf("the start of the ramp is %#08x, want the gradient rather than the colour after it", uint32(c))
	}
}

// TestRenderPaintsNothingForAGradientWithNoStops is one of the four ways SVG
// says a gradient does not paint: a ramp with no colours on it has nothing to
// say about any point of the shape.
func TestRenderPaintsNothingForAGradientWithNoStops(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 10 10">
		<defs>
			<linearGradient id="g" x1="0" x2="1"/>
		</defs>
		<rect x="0" y="0" width="10" height="10" fill="url(#g)"/>
	</svg>`, 10, 10)
	if c := pixelAt(cv, 5, 5); c.A() != 0 {
		t.Errorf("the rectangle came out %#08x, want it clear", uint32(c))
	}
}

// TestStyleKeepsTheReferenceItCannotFollowYet is the shape of what `with` does
// with a `url()`: it keeps the id rather than dropping the fill, and it says
// nothing about it, because only the whole drawing can.
func TestStyleKeepsTheReferenceItCannotFollowYet(t *testing.T) {
	var warned int
	warn := func(format string, args ...any) { warned++ }
	st := Style{}.with(mustElement(t, `<rect fill="url(#g)"/>`), warn)
	if st.fillRef != "g" {
		t.Errorf("the fill names %q, want the id it was given", st.fillRef)
	}
	if st.HasFill {
		t.Error("the fill was decided before the drawing was read")
	}
	if warned != 0 {
		t.Errorf("it warned %d times, want none until the drawing is known", warned)
	}
}
