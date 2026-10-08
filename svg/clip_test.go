package svg

import (
	"strings"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// A clip is a cut made to a picture: the shapes that say what the cut keeps are
// written once, and every element that names them is drawn into a picture of
// its own, cut to the shape of them, and laid over what is behind. What is
// here is what that comes out as. Points are asked for in the drawing's own
// units, which is what a viewBox of the same numbers makes the pixels worth.

// cut is a drawing that painted itself and had nothing to say about how it was
// read, which is what most of these tests want to know before they look at
// what came out. The ones that expect something said read the drawing
// themselves instead.
func cut(t *testing.T, src string, w, h int) *canvas.Canvas {
	t.Helper()
	img, cv := painted(t, src, w, h)
	if len(img.Warnings()) != 0 {
		t.Fatalf("the drawing said %v, want it read whole", img.Warnings())
	}
	return cv
}

// alphaAt is how much of what was painted at a point is there: 255 where the
// drawing covers the pixel whole, 0 where it does not reach, and a number in
// between where a soft edge is over it.
func alphaAt(cv *canvas.Canvas, x, y int) int {
	return int(pixelAt(cv, x, y).A())
}

// kept says the point came out under the drawing and gone says it came out
// clear, which is the whole of what a clip is about: what its shapes name, and
// nowhere else.
func kept(t *testing.T, cv *canvas.Canvas, x, y int) {
	t.Helper()
	if a := alphaAt(cv, x, y); a < 250 {
		t.Errorf("at (%d,%d) the drawing is %d out of 255, want it kept", x, y, a)
	}
}

func gone(t *testing.T, cv *canvas.Canvas, x, y int) {
	t.Helper()
	if a := alphaAt(cv, x, y); a != 0 {
		t.Errorf("at (%d,%d) the drawing is %d out of 255, want it cut away", x, y, a)
	}
}

// TestRenderCutsAnElementToTheShapesItNames is the whole of what a clip-path
// does: the element fills the drawing, so what is missing from the middle of
// the result is the clip's doing and nothing else's.
func TestRenderCutsAnElementToTheShapesItNames(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<defs>
			<clipPath id="c">
				<rect x="2" y="2" width="6" height="6"/>
			</clipPath>
		</defs>
		<rect width="10" height="10" fill="#ff0000" clip-path="url(#c)"/>
	</svg>`, 10, 10)

	for _, p := range [][2]int{{3, 3}, {5, 5}, {7, 5}, {5, 7}, {4, 4}} {
		kept(t, cv, p[0], p[1])
	}
	for _, p := range [][2]int{{1, 5}, {5, 1}, {9, 5}, {5, 9}, {1, 1}, {9, 9}} {
		gone(t, cv, p[0], p[1])
	}
}

// TestRenderKeepsWhatAnyOfTheShapesOfAClipPathCovers: several shapes in one
// clipPath are a union, so a point under any one of them is kept, and a point
// under none of them is nowhere — however much the shapes about it cover.
func TestRenderKeepsWhatAnyOfTheShapesOfAClipPathCovers(t *testing.T) {
	// The clip is written after the element that names it, which is the usual
	// place for one: only the whole drawing can say what an id is.
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<rect width="10" height="10" fill="#ff0000" clip-path="url(#c)"/>
		<clipPath id="c">
			<circle cx="2" cy="5" r="1.5"/>
			<circle cx="8" cy="5" r="1.5"/>
		</clipPath>
	</svg>`, 10, 10)

	for _, p := range [][2]int{{1, 5}, {2, 5}, {2, 4}, {7, 5}, {8, 5}, {7, 4}} {
		kept(t, cv, p[0], p[1])
	}
	// Between the two shapes, which leave a gap of two whole pixels, and
	// anywhere else the circles do not reach.
	for _, p := range [][2]int{{4, 5}, {5, 5}, {5, 0}, {0, 0}, {9, 9}} {
		gone(t, cv, p[0], p[1])
	}
}

// TestRenderOfAClipPathWithNoShapesKeepsNothing: the union of nothing is
// nothing, so a clipPath that holds no shapes at all is a cut that keeps
// nowhere rather than no cut at all.
func TestRenderOfAClipPathWithNoShapesKeepsNothing(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<rect width="10" height="10" fill="#ff0000" clip-path="url(#c)"/>
		<clipPath id="c"/>
	</svg>`, 10, 10)

	for _, p := range [][2]int{{1, 1}, {5, 5}, {9, 9}} {
		gone(t, cv, p[0], p[1])
	}
}

// TestRenderOfAClipPathThatIsNotThereKeepsNothingAndSaysSoOnce is the rest of
// what a reference to nowhere means here: what is cut cannot be drawn, and the
// group says it once for itself rather than once for every shape inside it.
func TestRenderOfAClipPathThatIsNotThereKeepsNothingAndSaysSoOnce(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 10 10">
		<g clip-path="url(#nowhere)">
			<rect width="10" height="10" fill="#ff0000"/>
			<rect x="0" y="0" width="5" height="5" fill="#00ff00"/>
			<circle cx="5" cy="5" r="4" fill="#0000ff"/>
		</g>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	ws := img.Warnings()
	if len(ws) != 1 {
		t.Fatalf("the drawing said %v, want it said once for the whole group", ws)
	}
	if !strings.Contains(ws.String(), "nowhere") {
		t.Errorf("the warning is %q, want it to name what the drawing could not find", ws.String())
	}
	for _, p := range [][2]int{{1, 1}, {5, 5}, {9, 9}} {
		gone(t, img.Render(10, 10), p[0], p[1])
	}
}

// TestRenderCutsAGroupAsOnePicture: the shapes under a group's clip are one
// picture with one cut over it, so where two half-transparent rectangles
// overlap they overlap first and the whole of that is cut, rather than each
// being cut on its own and the two cut pieces meeting somewhere the clip does
// not keep at all.
func TestRenderCutsAGroupAsOnePicture(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 12 10">
		<defs>
			<clipPath id="c"><rect x="0" y="0" width="8" height="10"/></clipPath>
		</defs>
		<g clip-path="url(#c)">
			<rect x="0" y="0" width="6" height="10" fill="#ff0000" fill-opacity="0.5"/>
			<rect x="4" y="0" width="8" height="10" fill="#ff0000" fill-opacity="0.5"/>
		</g>
	</svg>`, 12, 10)

	// Where only one rectangle reaches, half of it: 0.5 of 255.
	if a := alphaAt(cv, 1, 5); a < 124 || a > 132 {
		t.Errorf("at (1,5) the drawing is %d out of 255, want the half of one rectangle", a)
	}
	// Where both reach, one over the other: 0.5 + 0.5 of what the first left
	// showing, which is three quarters of the whole.
	if a := alphaAt(cv, 5, 5); a < 188 || a > 195 {
		t.Errorf("at (5,5) the drawing is %d out of 255, want one rectangle over the other", a)
	}
	// Both rectangles run past the clip, and neither of them reaches past it.
	for _, p := range [][2]int{{9, 1}, {11, 5}, {9, 9}} {
		gone(t, cv, p[0], p[1])
	}
}

// TestRenderMovesTheClipWithTheElementItCuts: the clip is written in the
// coordinates of the element that names it, so a transform on that element
// moves the cut along with it. Here the two only meet because of that — as
// written, the rectangle starts where the clip stops.
func TestRenderMovesTheClipWithTheElementItCuts(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 15 10">
		<defs>
			<clipPath id="c"><rect x="0" y="0" width="5" height="10"/></clipPath>
		</defs>
		<rect width="10" height="10" transform="translate(5,0)" fill="#ff0000" clip-path="url(#c)"/>
	</svg>`, 15, 10)

	for _, p := range [][2]int{{5, 5}, {7, 5}, {9, 5}} {
		kept(t, cv, p[0], p[1])
	}
	for _, p := range [][2]int{{4, 5}, {1, 5}, {11, 5}, {14, 5}} {
		gone(t, cv, p[0], p[1])
	}
}

// TestRenderDoesNotCutWhatComesAfterTheElementTheClipIsOn: the clip is not
// inherited, so it cuts the element that named it and nothing else — the shape
// written next stands whole, even where the clipped one was cut away.
func TestRenderDoesNotCutWhatComesAfterTheElementTheClipIsOn(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<defs>
			<clipPath id="c"><rect x="0" y="0" width="5" height="10"/></clipPath>
		</defs>
		<rect width="10" height="10" fill="#ff0000" clip-path="url(#c)"/>
		<rect x="6" y="0" width="4" height="10" fill="#0000ff"/>
	</svg>`, 10, 10)

	kept(t, cv, 2, 5)
	if got := pixelAt(cv, 7, 5); got != canvas.RGB(0, 0, 255) {
		t.Errorf("the shape after the clipped one is %v, want it blue and whole", got)
	}
	// And where the clipped rectangle is cut away, the drawing is clear even
	// though the rectangle alone would have covered this point.
	gone(t, cv, 5, 7)
}

// TestRenderReadsTheClipRuleOfTheShapesItCutsWith: one outline drawn as two
// rectangles that cross. With the rule a fill of the same outline takes, the
// crossing is inside twice and so is inside; with evenodd it counts twice and
// so is nowhere, which puts a hole through the middle of the clip. The rule may
// be written on the clipPath, where everything inside it takes it, or on the
// one shape that needs it, or nowhere at all — where the same rule a fill takes
// is what counts.
func TestRenderReadsTheClipRuleOfTheShapesItCutsWith(t *testing.T) {
	const shape = "M0 0 H10 V10 H0 Z M5 0 H15 V10 H5 Z"
	for _, tc := range []struct {
		name string
		clip string
		hole bool
	}{
		{"written on the clipPath", `<clipPath id="c" clip-rule="evenodd"><path d="` + shape + `"/></clipPath>`, true},
		{"written on the shape", `<clipPath id="c"><path clip-rule="evenodd" d="` + shape + `"/></clipPath>`, true},
		{"written nowhere", `<clipPath id="c"><path d="` + shape + `"/></clipPath>`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cv := cut(t, `<svg viewBox="0 0 15 10">
				<rect width="15" height="10" fill="#ff0000" clip-path="url(#c)"/>
				`+tc.clip+`
			</svg>`, 15, 10)
			if tc.hole {
				gone(t, cv, 7, 5)
			} else {
				kept(t, cv, 7, 5)
			}
			kept(t, cv, 2, 5)
			kept(t, cv, 12, 5)
		})
	}
}

// TestRenderCutsTheStrokeOfAnElementToo: a stroke is part of the element's
// picture like its fill is, so the ring comes out where the clip keeps it and
// is nowhere it does not — including at the line the clip stops on.
func TestRenderCutsTheStrokeOfAnElementToo(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<defs>
			<clipPath id="c"><rect x="0" y="0" width="5" height="10"/></clipPath>
		</defs>
		<circle cx="5" cy="5" r="4" fill="none" stroke="#ff0000" stroke-width="2" clip-path="url(#c)"/>
	</svg>`, 10, 10)

	kept(t, cv, 1, 5)
	kept(t, cv, 1, 7)
	gone(t, cv, 9, 5)
	gone(t, cv, 5, 1)
}

// TestRenderMovesTheClipOfAUseWithItsXAndY: a `<use>` is drawn where its x and
// y put it, and its clip goes there too — the clip is followed after those have
// moved it, so the two start from the same place rather than from the corner
// the drawing was written in.
func TestRenderMovesTheClipOfAUseWithItsXAndY(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 20 10">
		<defs>
			<clipPath id="c"><rect x="0" y="0" width="5" height="10"/></clipPath>
			<rect id="r" width="15" height="10" fill="#ff0000"/>
		</defs>
		<use href="#r" x="5" clip-path="url(#c)"/>
	</svg>`, 20, 10)

	// The rectangle is drawn from 5 to 20 and the clip from 5 to 10, so what
	// comes out is the two where they meet — not nothing at all, which is what
	// the clip cut in the corner the rectangle came from.
	for _, p := range [][2]int{{5, 5}, {7, 5}, {9, 5}} {
		kept(t, cv, p[0], p[1])
	}
	for _, p := range [][2]int{{1, 5}, {4, 5}, {11, 5}, {17, 5}} {
		gone(t, cv, p[0], p[1])
	}
}

// TestRenderOfAClipPathInObjectBoundingBox: a clip written in objectBoundingBox
// units is measured against the element's painted box and applied correctly.
// The element is not drawn whole and no warning is emitted for measurable boxes.
func TestRenderOfAClipPathInObjectBoundingBox(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 10 10">
		<rect width="10" height="10" fill="#ff0000" clip-path="url(#c)"/>
		<clipPath id="c" clipPathUnits="objectBoundingBox">
			<rect width="1" height="1"/>
		</clipPath>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	ws := img.Warnings()
	if len(ws) != 0 {
		t.Fatalf("the drawing said %v, want no warning for measurable box", ws)
	}
	// The clip should be applied without warning; rendering should succeed.
	cv := img.Render(10, 10)
	if cv == nil {
		t.Fatal("render returned nil")
	}
	// Verify the clip covers the element (1×1 in objectBoundingBox = 100% of the box).
	// All pixels should be kept since the clip covers the full rect.
	for _, p := range [][2]int{{0, 0}, {5, 5}, {9, 9}} {
		kept(t, cv, p[0], p[1])
	}
}

// TestRenderMeasuresWritingUnderAClipInFractionsOfABox: the shapes of a
// clipPath written in objectBoundingBox units are fractions of the box of the
// element they cut, and the box a piece of writing comes out in is measured
// the same as any other — so the clip is followed instead of being left off
// with a warning that there was no box of writing to measure it against. The
// clip here covers the whole of the box, so what comes out is the writing
// whole; where the fractions of the box put the cut itself is the cut's own
// business, under any clip and not only this one.
func TestRenderMeasuresWritingUnderAClipInFractionsOfABox(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 40 20">
		<text x="4" y="16" font-size="8" fill="#ff0000" clip-path="url(#c)">hi</text>
		<clipPath id="c" clipPathUnits="objectBoundingBox">
			<rect width="1" height="1"/>
		</clipPath>
	</svg>`, 40, 20)

	// Writing at this size comes out with every edge antialiased, so what is
	// asked is that anything was painted rather than how much: both letters
	// are inside the box the clip covers.
	if alphaAt(cv, 4, 15) == 0 {
		t.Error("at (4,15) the drawing is clear, want the letter at the start of the writing kept")
	}
	if alphaAt(cv, 10, 15) == 0 {
		t.Error("at (10,15) the drawing is clear, want the letter past the middle kept")
	}
}

// TestRenderOfAClipPathOfNoneDrawsTheElementWhole: `none` is a clip that is
// not there, which is read as no clip and not as a complaint.
func TestRenderOfAClipPathOfNoneDrawsTheElementWhole(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<rect width="10" height="10" fill="#ff0000" clip-path="none"/>
	</svg>`, 10, 10)

	for _, p := range [][2]int{{0, 0}, {5, 5}, {9, 9}} {
		kept(t, cv, p[0], p[1])
	}
}

// TestRenderOfWritingInAClipPathSaysSoAndKeepsNothing: writing has no outline
// to cut with here, so it is left out of the clip — and a clip of nothing keeps
// nothing, with the drawing saying what it could not do rather than cutting
// something arbitrary instead.
func TestRenderOfWritingInAClipPathSaysSoAndKeepsNothing(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 10 10">
		<rect width="10" height="10" fill="#ff0000" clip-path="url(#c)"/>
		<clipPath id="c"><text x="0" y="5">hi</text></clipPath>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if !strings.Contains(img.Warnings().String(), "writing inside a <clipPath>") {
		t.Errorf("the drawing said %v, want it to say that writing cannot cut", img.Warnings())
	}
	for _, p := range [][2]int{{1, 1}, {5, 5}, {9, 9}} {
		gone(t, img.Render(10, 10), p[0], p[1])
	}
}
