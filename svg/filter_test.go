package svg

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// A filter is a list of things done to a picture on its way into the drawing:
// a colour turned grey, then blurred, in the order the file wrote them. What is
// here is what those come out as. Points are asked for in the drawing's own
// units, which is what a viewBox of the same numbers makes the pixels worth,
// and the colour of a whole filter is asked for well inside a shape so that no
// pixel here is one a soft edge had a say in.

// filteredAt says the point came out in the colour the filter left there, and
// the drawing painted it whole: the four things every colour filter promises,
// since none of them may touch how much is there.
func filteredAt(t *testing.T, cv *canvas.Canvas, x, y int, r, g, b uint8) {
	t.Helper()
	c := pixelAt(cv, x, y)
	if !nearByte(c.R(), r) || !nearByte(c.G(), g) || !nearByte(c.B(), b) || c.A() != 255 {
		t.Errorf("at (%d,%d) the pixel is %#08x, want (%d,%d,%d) whole", x, y, uint32(c), r, g, b)
	}
}

// TestFilterGreyscalesThePicture is the plainest of the colour filters: red is
// not a grey, and grayscale(1) says all the way to one — the luminance grey
// that the matrix in the canvas is built from, which for full red is the red's
// share of white and nothing else.
func TestFilterGreyscalesThePicture(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<rect width="10" height="10" fill="#ff0000" filter="grayscale(1)"/>
	</svg>`, 10, 10)

	filteredAt(t, cv, 5, 5, 54, 54, 54)
}

// TestFilterSepiasThePicture is the other matrix with one fixed answer: sepia
// is what a photograph looks like toned brown, and full red toned all the way
// is the first row of that matrix and no more.
func TestFilterSepiasThePicture(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<rect width="10" height="10" fill="#ff0000" filter="sepia(1)"/>
	</svg>`, 10, 10)

	filteredAt(t, cv, 5, 5, 100, 89, 69)
}

// TestFilterRotatesTheHue turns the colour round the wheel: red at 120
// degrees is green's part of the wheel, and the red that was there is gone
// from both ends of it — nothing of the original red is left in the result,
// only the green a third of the way round.
func TestFilterRotatesTheHue(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<rect width="10" height="10" fill="#ff0000" filter="hue-rotate(120deg)"/>
	</svg>`, 10, 10)

	filteredAt(t, cv, 5, 5, 0, 113, 0)
}

// TestFilterBlursPastTheEdgeOfWhatWasDrawn is the filter that reaches beyond
// the element that asked for it: the rectangle stops at 4, and the blurred
// picture carries a little past that and a little less the further out it
// goes, while the middle of the rectangle — far enough from every edge — is
// left whole. The same drawing without the filter says which of that is the
// blur's doing: nothing at all is outside the rectangle, and the edge of it is
// solid where the blur softened it.
func TestFilterBlursPastTheEdgeOfWhatWasDrawn(t *testing.T) {
	const drawing = `<svg viewBox="0 0 15 15">
		<rect x="4" y="4" width="7" height="7" fill="#ff0000"%s/>
	</svg>`
	blurred := cut(t, fmt.Sprintf(drawing, ` filter="blur(1)"`), 15, 15)
	plain := cut(t, fmt.Sprintf(drawing, ""), 15, 15)

	if a := alphaAt(blurred, 1, 7); a == 0 {
		t.Error("outside the rectangle the drawing is clear, want the blur carried past where it stopped")
	}
	if a := alphaAt(plain, 1, 7); a != 0 {
		t.Errorf("without the filter the drawing outside the rectangle is %d out of 255, want nothing there", a)
	}
	if a := alphaAt(blurred, 0, 7); a != 0 {
		t.Errorf("three units out from the edge the drawing is %d out of 255, want the blur's reach spent", a)
	}
	if a, want := alphaAt(blurred, 4, 7), 255; a == want || a == 0 {
		t.Errorf("at the edge of the rectangle the drawing is %d out of 255, want it softened but still there", a)
	}
	if a := alphaAt(plain, 4, 7); a != 255 {
		t.Errorf("without the filter the edge of the rectangle is %d out of 255, want it solid", a)
	}
	if a := alphaAt(blurred, 7, 7); a != 255 {
		t.Errorf("in the middle of the rectangle the drawing is %d out of 255, want the blur to leave it whole", a)
	}
}

// TestFilterKeepsTheOrderTheFunctionsWereWrittenIn is what a list means: the
// sepia browned first and the grey taken of the brown is one colour, and the
// grey taken first and the brown laid over that is another. Neither is wrong
// and neither is the other — only the order the file wrote them in says which
// one it asked for.
func TestFilterKeepsTheOrderTheFunctionsWereWrittenIn(t *testing.T) {
	brownThenGrey := cut(t, `<svg viewBox="0 0 10 10">
		<rect width="10" height="10" fill="#ff0000" filter="sepia(1) grayscale(1)"/>
	</svg>`, 10, 10)
	greyThenBrown := cut(t, `<svg viewBox="0 0 10 10">
		<rect width="10" height="10" fill="#ff0000" filter="grayscale(1) sepia(1)"/>
	</svg>`, 10, 10)

	filteredAt(t, brownThenGrey, 5, 5, 90, 90, 90)
	filteredAt(t, greyThenBrown, 5, 5, 73, 65, 51)
}

// TestFilterOnAGroupIsNotInheritedByItsChildren is the whole of what it means
// that a filter is not inherited: the group's picture is put through the
// brightness once and comes out half as bright, which is what one pass gives —
// a second pass inside the first would have halved it again, and the child
// would have come out a quarter. What is outside the group is not in that
// picture at all and stands as it was.
func TestFilterOnAGroupIsNotInheritedByItsChildren(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<g filter="brightness(0.5)">
			<rect x="1" y="1" width="4" height="4" fill="#ff0000"/>
		</g>
		<rect x="6" y="6" width="4" height="4" fill="#ff0000"/>
	</svg>`, 10, 10)

	filteredAt(t, cv, 2, 2, 128, 0, 0)
	filteredAt(t, cv, 7, 7, 255, 0, 0)
}

// TestFilterNoneSaysThereIsNothingToDo is the one spelling of no filter, and
// like everywhere else in a drawing it says so without complaint and without
// changing what is there.
func TestFilterNoneSaysThereIsNothingToDo(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<rect width="10" height="10" fill="#ff0000" filter="none"/>
	</svg>`, 10, 10)

	filteredAt(t, cv, 5, 5, 255, 0, 0)
}

// TestBlurIsMeasuredInTheDrawingsOwnUnits is what a size in a filter is a size
// of: the same drawing read at two sizes blurs as far past its rectangle both
// times, because one unit of the drawing is a different number of pixels each
// time and the blur is a number of the drawing's units. A blur counted in raw
// pixels instead would reach further past the shape the bigger the picture got,
// and the probe far out from the rectangle would be the difference.
func TestBlurIsMeasuredInTheDrawingsOwnUnits(t *testing.T) {
	const drawing = `<svg viewBox="0 0 10 10">
		<rect width="5" height="5" fill="#ff0000" filter="blur(1)"/>
	</svg>`
	_, small := painted(t, drawing, 10, 10)
	_, big := painted(t, drawing, 20, 20)
	_, smallPlain := painted(t, `<svg viewBox="0 0 10 10">
		<rect width="5" height="5" fill="#ff0000"/>
	</svg>`, 10, 10)
	_, bigPlain := painted(t, `<svg viewBox="0 0 10 10">
		<rect width="5" height="5" fill="#ff0000"/>
	</svg>`, 20, 20)

	if a := alphaAt(small, 7, 2); a == 0 {
		t.Error("in the small drawing the blur does not reach two units past the rectangle")
	}
	if a := alphaAt(big, 14, 5); a == 0 {
		t.Error("in the big drawing the blur does not reach two units past the rectangle")
	}
	if a := alphaAt(smallPlain, 7, 2); a != 0 {
		t.Errorf("without the filter the small drawing reaches out at (%d,%d): %d out of 255", 7, 2, a)
	}
	if a := alphaAt(bigPlain, 14, 5); a != 0 {
		t.Errorf("without the filter the big drawing reaches out at (%d,%d): %d out of 255", 14, 5, a)
	}
}

// TestParseSaysWhatItCannotReadInTheFilter is every way a filter can ask for
// something this package does not do: a function that is not one of the ones
// it can apply, an argument that is not a number, a reference to a `<filter>`
// element whose insides are not read here, and nothing at all. Each says what
// is wrong and once, and the drawing still paints what it can of itself.
func TestParseSaysWhatItCannotReadInTheFilter(t *testing.T) {
	for _, tc := range []struct{ filter, want string }{
		{"blur(zzz)", "is not one this package can read"},
		{"saturate(1)", "is not one this package can read"},
		{"drop-shadow(0 1px 2px black)", "is not one this package can read"},
		{"blur(-1)", "is not one this package can read"},
		{"url(#f)", "is not read here"},
		{"", "is not one this package can read"},
	} {
		img, err := Parse(`<svg viewBox="0 0 10 10"><rect width="10" height="10" filter="` + tc.filter + `"/></svg>`)
		if err != nil {
			t.Fatalf("filter %q: parse: %v", tc.filter, err)
		}
		ws := img.Warnings()
		if len(ws) != 1 {
			t.Errorf("filter %q said %v, want it said once", tc.filter, ws)
		} else if !strings.Contains(ws.String(), tc.want) {
			t.Errorf("filter %q said %v, want it to say %q", tc.filter, ws, tc.want)
		}
		if cv := img.Render(10, 10); cv == nil {
			t.Errorf("filter %q: the drawing did not paint", tc.filter)
		}
	}
}

// TestAPartialFilterListStillDoesThePartItCan is what one function left out
// costs: the saturate this package cannot read is said and left off, and the
// blur beside it in the same list still runs, because a drawing that asks for
// two things and gets the one that was readable is nearer what it asked for
// than one that gets neither.
func TestAPartialFilterListStillDoesThePartItCan(t *testing.T) {
	img, cv := painted(t, `<svg viewBox="0 0 15 15">
		<rect x="4" y="4" width="7" height="7" fill="#ff0000" filter="blur(1) saturate(1)"/>
	</svg>`, 15, 15)

	ws := img.Warnings()
	if len(ws) != 1 {
		t.Fatalf("the drawing said %v, want the saturate said once", ws)
	}
	if !strings.Contains(ws.String(), "saturate") {
		t.Errorf("the drawing said %v, want it to name the saturate", ws)
	}
	if a := alphaAt(cv, 1, 7); a == 0 {
		t.Error("the blur in the list did not run, want the part of the list that could be read still done")
	}
}

// TestFilterComesFromTheStyleAttribute is the same filter written the way a
// stylesheet writes it: `filter` is a property like any other, so the
// declaration says what the attribute would have said.
func TestFilterComesFromTheStyleAttribute(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<rect width="10" height="10" fill="#ff0000" style="filter: grayscale(1)"/>
	</svg>`, 10, 10)

	filteredAt(t, cv, 5, 5, 54, 54, 54)
}

// TestReadFiltersTakesWhatItCanAndSaysTheRest reads the list itself rather
// than through a drawing, which is where each spelling is decided: `none` is
// nothing at all and says nothing, an empty list is nothing at all and says
// so, a function with no arguments is its full strength, a percentage is a
// fraction of it, an angle is in any of the units an angle may be written in,
// and one function the package cannot read leaves the rest of the list where
// it was.
func TestReadFiltersTakesWhatItCanAndSaysTheRest(t *testing.T) {
	for _, tc := range []struct {
		raw   string
		want  []filterOp
		warns int
	}{
		{"none", nil, 0},
		{"", nil, 1},
		{"grayscale", nil, 1},
		{"grayscale(50%)", []filterOp{{canvas.FilterGrayscale, 0.5}}, 0},
		{"sepia()", []filterOp{{canvas.FilterSepia, 1}}, 0},
		{"hue-rotate(0.5turn)", []filterOp{{canvas.FilterHueRotate, 180}}, 0},
		{"blur(2px)", []filterOp{{canvas.FilterBlur, 2}}, 0},
		{"blur()", []filterOp{{canvas.FilterBlur, 0}}, 0},
		{"brightness(200%)", []filterOp{{canvas.FilterBrightness, 2}}, 0},
		{"url(#f) invert(1)", []filterOp{{canvas.FilterInvert, 1}}, 1},
		{"drop-shadow(0 1px 2px black) contrast(1)", []filterOp{{canvas.FilterContrast, 1}}, 1},
	} {
		var warns int
		got := readFilters(tc.raw, func(string, ...any) { warns++ })
		if warns != tc.warns {
			t.Errorf("readFilters(%q) warned %d times, want %d", tc.raw, warns, tc.warns)
		}
		if len(got) != len(tc.want) {
			t.Errorf("readFilters(%q) gave %v, want %v", tc.raw, got, tc.want)
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("readFilters(%q) gave %v, want %v", tc.raw, got, tc.want)
				break
			}
		}
	}
}
