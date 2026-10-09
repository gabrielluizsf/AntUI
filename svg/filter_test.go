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

// TestRenderFiltersBeforeItCutsAndMasks is the order a picture goes through
// what an element asked for: the filter first, so the blur has done its work
// before the clip says where any of it is kept; then the clip; then the mask,
// which says how much of what is left is there. Reading the clip first would
// let the filter soften the cut it made, which is the one of the three that can
// be seen — the clip and the mask both take away rather than add, so their
// order between themselves leaves the same picture.
func TestRenderFiltersBeforeItCutsAndMasks(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 20 20">
		<defs>
			<clipPath id="c"><rect x="8" y="8" width="6" height="6"/></clipPath>
			<mask id="m"><rect x="8" y="8" width="3" height="6" fill="#ffffff"/></mask>
		</defs>
		<rect x="6" y="6" width="8" height="8" fill="#ff0000" filter="blur(1)" clip-path="url(#c)" mask="url(#m)"/>
	</svg>`, 20, 20)

	// One unit inside the cut, with the mask over it, is solid: had the clip
	// been read before the blur, the blur would have softened the edge the cut
	// left and this pixel would be part covered. Past the mask, and out where
	// only the blur reached, there is nothing.
	kept(t, cv, 9, 11)
	gone(t, cv, 12, 11)
	gone(t, cv, 7, 11)
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

// TestFilterDropShadowLaysTheSilhouetteBehind is what a drop-shadow is: not a
// colour put over the picture but a second one under it — the shape of the
// element in the shadow's colour, moved by the offset. Where only the shadow
// reaches the shadow's colour is what is there; where only the element reaches
// the element's own colour is; and past the shadow there is nothing at all.
func TestFilterDropShadowLaysTheSilhouetteBehind(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<rect x="1" y="1" width="4" height="4" fill="#ff0000" filter="drop-shadow(2 0 0 #0000ff)"/>
	</svg>`, 10, 10)

	filteredAt(t, cv, 2, 3, 255, 0, 0)
	filteredAt(t, cv, 6, 3, 0, 0, 255)
	gone(t, cv, 8, 3)
}

// TestFilterDropShadowBlursTheSilhouetteSoftensIt is the blur of a shadow: a
// shadow with a radius carries past where the shape it was cast from stopped,
// the way a blur does, and the same drawing without the shadow is clear there.
func TestFilterDropShadowBlursTheSilhouetteSoftensIt(t *testing.T) {
	const drawing = `<svg viewBox="0 0 10 10">
		<rect x="3" y="3" width="4" height="4" fill="#ff0000"%s/>
	</svg>`
	shadowed := cut(t, fmt.Sprintf(drawing, ` filter="drop-shadow(0 0 2 #0000ff)"`), 10, 10)
	plain := cut(t, fmt.Sprintf(drawing, ""), 10, 10)

	c := pixelAt(shadowed, 2, 5)
	if c.A() == 0 || c.B() < 128 {
		t.Errorf("a unit outside the rectangle the shadow is %#08x, want blue carried past the shape", uint32(c))
	}
	if a := alphaAt(plain, 2, 5); a != 0 {
		t.Errorf("without the shadow the drawing outside the rectangle is %d out of 255, want nothing there", a)
	}
}

// TestFilterBlursEachAxisAsTheTransformStretchesIt is why a blur needs two
// radii: a transform that stretches one axis carries the blur with it, so the
// blur reaches further along that axis than across the other, rather than the
// same distance both ways as an average of the two would give.
func TestFilterBlursEachAxisAsTheTransformStretchesIt(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 60 60">
		<rect x="8" y="29" width="1" height="1" fill="#ff0000" filter="blur(1)" transform="scale(4,1)"/>
	</svg>`, 60, 60)

	// The square lands at x=32..36, y=29..30. Along x the blur is stretched
	// four times and reaches five pixels out; along y it is not stretched and
	// is gone by then.
	if c := pixelAt(cv, 27, 30); c.A() == 0 {
		t.Error("five pixels to the left is clear, want the blur carried there by the stretched axis")
	}
	if c := pixelAt(cv, 34, 24); c.A() != 0 {
		t.Errorf("five pixels above is %v, want nothing: the blur is not stretched up the other axis", c)
	}
}

// TestFilterDropShadowWithoutAColourTakesTheCurrentColour is the colour CSS
// gives a shadow it was not told: `currentColor`, which here is the colour the
// drawing is painted in rather than a black guessed at read time.
func TestFilterDropShadowWithoutAColourTakesTheCurrentColour(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 10 10">
		<rect x="1" y="1" width="4" height="4" fill="#ff0000" filter="drop-shadow(2 0 0)"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if ws := img.Warnings(); len(ws) != 0 {
		t.Fatalf("the drawing said %v, want a shadow it reads whole", ws)
	}
	cv := img.RenderWith(10, 10, canvas.RGB(0, 255, 0))
	if cv == nil {
		t.Fatal("the drawing did not paint")
	}

	filteredAt(t, cv, 6, 3, 0, 255, 0)
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
// element the drawing does not have, and nothing at all. Each says what is
// wrong and once, and the drawing still paints what it can of itself.
func TestParseSaysWhatItCannotReadInTheFilter(t *testing.T) {
	for _, tc := range []struct{ filter, want string }{
		{"blur(zzz)", "is not one this package can read"},
		{"saturate(1)", "is not one this package can read"},
		{"drop-shadow(1px zzz)", "is not one this package can read"},
		{"drop-shadow(1 2 3 4 black)", "is not one this package can read"},
		{"blur(-1)", "is not one this package can read"},
		{"url(#f)", "which the drawing does not have"},
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
		{"grayscale(50%)", []filterOp{{kind: canvas.FilterGrayscale, amount: 0.5}}, 0},
		{"sepia()", []filterOp{{kind: canvas.FilterSepia, amount: 1}}, 0},
		{"hue-rotate(0.5turn)", []filterOp{{kind: canvas.FilterHueRotate, amount: 180}}, 0},
		{"blur(2px)", []filterOp{{kind: canvas.FilterBlur, amount: 2}}, 0},
		{"blur()", []filterOp{{kind: canvas.FilterBlur, amount: 0}}, 0},
		{"brightness(200%)", []filterOp{{kind: canvas.FilterBrightness, amount: 2}}, 0},
		{"url(#f) invert(1)", []filterOp{{ref: "f"}, {kind: canvas.FilterInvert, amount: 1}}, 0},
		{"drop-shadow(0 1px 2px black) contrast(1)", []filterOp{
			{shadow: &dropShadow{dy: 1, blur: 2, colour: canvas.Color(0xFF000000)}},
			{kind: canvas.FilterContrast, amount: 1},
		}, 0},
	} {
		var warns int
		got := readFilters(tc.raw, func(string, ...any) { warns++ })
		if warns != tc.warns {
			t.Errorf("readFilters(%q) warned %d times, want %d", tc.raw, warns, tc.warns)
		}
		if !sameOps(got, tc.want) {
			t.Errorf("readFilters(%q) gave %v, want %v", tc.raw, got, tc.want)
		}
	}
}

// sameOps says two filter lists read the same, comparing a step that stands for
// a shadow by the shadow it holds rather than by the pointer to it.
func sameOps(got, want []filterOp) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] == want[i] {
			continue
		}
		if got[i].shadow == nil || want[i].shadow == nil || *got[i].shadow != *want[i].shadow {
			return false
		}
	}
	return true
}

// TestFilterReferenceRunsTheFilterElementItNames is the plainest form of the
// other way a file asks for a filter: the reference follows the `<filter>`
// element it names, and what is inside one runs where the reference stood.
// The primitive here takes no input of its own, which means the picture the
// element drew — so the rectangle comes out moved by the offset, which is
// only something a filter element could have done.
func TestFilterReferenceRunsTheFilterElementItNames(t *testing.T) {
	img, cv := painted(t, `<svg viewBox="0 0 10 10">
		<filter id="f"><feOffset dx="1"/></filter>
		<rect width="10" height="10" fill="#ff0000" filter="url(#f)"/>
	</svg>`, 10, 10)

	if ws := img.Warnings(); len(ws) != 0 {
		t.Errorf("the drawing said %v, want nothing said for a filter it has", ws)
	}
	gone(t, cv, 0, 5)
	filteredAt(t, cv, 5, 5, 255, 0, 0)
}

// TestFilterPrimitiveChainKeepsTheNamesItGaveItsResults is a filter of more
// than one primitive: each leaves its picture under the name it was given,
// and what follows picks it up by that name. The blue the first floods the
// whole picture with is the second's second input, and what comes out of the
// second is the red rectangle where there is one and the blue underneath
// where there is not.
func TestFilterPrimitiveChainKeepsTheNamesItGaveItsResults(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 15 15">
		<filter id="f">
			<feFlood flood-color="#0000ff" result="f"/>
			<feBlend in="SourceGraphic" in2="f"/>
		</filter>
		<rect width="10" height="10" fill="#ff0000" filter="url(#f)"/>
	</svg>`, 15, 15)

	filteredAt(t, cv, 5, 5, 255, 0, 0)
	filteredAt(t, cv, 13, 13, 0, 0, 255)
}

// TestFilterColorMatrixSwapsTheChannelsOfEveryPixel is feColorMatrix with the
// matrix written out: red's two channels swapped comes out green, which is
// only true if the matrix ran on the colour of every pixel of the picture
// rather than on the shape of it.
func TestFilterColorMatrixSwapsTheChannelsOfEveryPixel(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 10 10">
		<filter id="f">
			<feColorMatrix type="matrix" values="0 1 0 0 0  1 0 0 0 0  0 0 1 0 0  0 0 0 1 0"/>
		</filter>
		<rect width="10" height="10" fill="#ff0000" filter="url(#f)"/>
	</svg>`, 10, 10)

	filteredAt(t, cv, 5, 5, 0, 255, 0)
}

// TestFilterCompositeArithmeticAddsItsOwnNumbers is feComposite's arithmetic
// on itself: every channel of the picture plus half of itself plus the
// constant, which for full red is red with half of everything else added —
// the numbers k2 and k4 said to be.
func TestFilterCompositeArithmeticAddsItsOwnNumbers(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 10 10">
		<filter id="f">
			<feComposite operator="arithmetic" k2="0.5" k4="0.5"/>
		</filter>
		<rect width="10" height="10" fill="#ff0000" filter="url(#f)"/>
	</svg>`, 10, 10)

	filteredAt(t, cv, 5, 5, 255, 128, 128)
}

// TestFilterBlendMultiplies is feBlend with its second input named: the grey
// the first primitive flooded the picture with is what the red rectangle is
// multiplied against, and where the rectangle stopped the grey is what is
// left whole underneath.
func TestFilterBlendMultiplies(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 15 15">
		<filter id="f">
			<feFlood flood-color="#808080" result="f"/>
			<feBlend in="SourceGraphic" in2="f" mode="multiply"/>
		</filter>
		<rect width="10" height="10" fill="#ff0000" filter="url(#f)"/>
	</svg>`, 15, 15)

	filteredAt(t, cv, 5, 5, 128, 0, 0)
	filteredAt(t, cv, 13, 13, 128, 128, 128)
}

// TestFilterSaysWhatItCannotReadInsideTheFilterElement is one primitive this
// package has no way to run inside a filter the rest of it can: the one is
// named in the warning and left out, and what comes after it in the same
// filter still runs, the same way one unreadable function in the attribute's
// list does.
func TestFilterSaysWhatItCannotReadInsideTheFilterElement(t *testing.T) {
	img, cv := painted(t, `<svg viewBox="0 0 10 10">
		<filter id="f">
			<feTurbulence baseFrequency="0.1"/>
			<feOffset dx="1"/>
		</filter>
		<rect width="10" height="10" fill="#ff0000" filter="url(#f)"/>
	</svg>`, 10, 10)

	ws := img.Warnings()
	if len(ws) != 1 {
		t.Fatalf("the drawing said %v, want the turbulence said once", ws)
	}
	if !strings.Contains(ws.String(), "feTurbulence") {
		t.Errorf("the drawing said %v, want it to name the feTurbulence", ws)
	}
	gone(t, cv, 0, 5)
	filteredAt(t, cv, 5, 5, 255, 0, 0)
}

// TestAPartialFilterListKeepsTheReferenceAndTheFunction is both ways of asking
// for a filter in one list: the reference and the function beside it run in
// the order they were written, so the rectangle moves and then goes grey —
// neither half of the list knows or cares what the other half was written as.
func TestAPartialFilterListKeepsTheReferenceAndTheFunction(t *testing.T) {
	img, cv := painted(t, `<svg viewBox="0 0 10 10">
		<filter id="f"><feOffset dx="1"/></filter>
		<rect width="10" height="10" fill="#ff0000" filter="url(#f) grayscale(1)"/>
	</svg>`, 10, 10)

	if ws := img.Warnings(); len(ws) != 0 {
		t.Errorf("the drawing said %v, want nothing said for a list it can read whole", ws)
	}
	gone(t, cv, 0, 5)
	filteredAt(t, cv, 5, 5, 54, 54, 54)
}

// TestFilterCompositeOperatorsEachTakeTheirOwnWayOfPuttingTwoPicturesTogether
// is feComposite's operators over the same two pictures — the red rectangle
// and the blue the first primitive flooded the whole picture with — where each
// one answers differently: what it keeps of the source, what it keeps of what
// was underneath, and what comes out where only one of the two is there.
func TestFilterCompositeOperatorsEachTakeTheirOwnWayOfPuttingTwoPicturesTogether(t *testing.T) {
	for _, tc := range []struct {
		operator string
		in, out  [3]int // -1 in a channel means nothing is there at all
	}{
		{"", [3]int{255, 0, 0}, [3]int{0, 0, 255}},
		{"in", [3]int{255, 0, 0}, [3]int{-1, -1, -1}},
		{"out", [3]int{-1, -1, -1}, [3]int{-1, -1, -1}},
		{"atop", [3]int{255, 0, 0}, [3]int{0, 0, 255}},
		{"xor", [3]int{-1, -1, -1}, [3]int{0, 0, 255}},
		{"lighter", [3]int{255, 0, 255}, [3]int{0, 0, 255}},
	} {
		img, cv := painted(t, fmt.Sprintf(`<svg viewBox="0 0 15 15">
			<filter id="f">
				<feFlood flood-color="#0000ff" result="b"/>
				<feComposite in="SourceGraphic" in2="b" operator=%q/>
			</filter>
			<rect width="10" height="10" fill="#ff0000" filter="url(#f)"/>
		</svg>`, tc.operator), 15, 15)

		if ws := img.Warnings(); len(ws) != 0 {
			t.Errorf("operator %q said %v, want nothing said for one it reads", tc.operator, ws)
		}
		for _, probe := range []struct {
			x, y int
			want [3]int
		}{{5, 5, tc.in}, {13, 13, tc.out}} {
			if probe.want[0] < 0 {
				gone(t, cv, probe.x, probe.y)
				continue
			}
			filteredAt(t, cv, probe.x, probe.y, uint8(probe.want[0]), uint8(probe.want[1]), uint8(probe.want[2]))
		}
	}
}

// TestFilterBlendModesEachTakeTheirOwnWayOfMixingTwoPictures is feBlend's
// modes over the same two pictures — the orange rectangle and the grey the
// first primitive flooded the whole picture with — where each mode answers
// differently, which is only visible on colours the two actually share. The
// grey outside the rectangle is what every one of them leaves there.
func TestFilterBlendModesEachTakeTheirOwnWayOfMixingTwoPictures(t *testing.T) {
	for _, tc := range []struct {
		mode string
		in   [3]int
	}{
		{"", [3]int{255, 128, 0}},
		{"multiply", [3]int{128, 64, 0}},
		{"screen", [3]int{255, 192, 128}},
		{"darken", [3]int{128, 128, 0}},
		{"lighten", [3]int{255, 128, 128}},
	} {
		img, cv := painted(t, fmt.Sprintf(`<svg viewBox="0 0 15 15">
			<filter id="f">
				<feFlood flood-color="#808080" result="b"/>
				<feBlend in="SourceGraphic" in2="b" mode=%q/>
			</filter>
			<rect width="10" height="10" fill="#ff8000" filter="url(#f)"/>
		</svg>`, tc.mode), 15, 15)

		if ws := img.Warnings(); len(ws) != 0 {
			t.Errorf("mode %q said %v, want nothing said for one it reads", tc.mode, ws)
		}
		filteredAt(t, cv, 5, 5, uint8(tc.in[0]), uint8(tc.in[1]), uint8(tc.in[2]))
		filteredAt(t, cv, 13, 13, 128, 128, 128)
	}
}

// TestFilterColorMatrixShortcutsWriteTheMatrixTheSpecGivesThem is the three
// kinds of feColorMatrix that are not the matrix itself: each says what it
// means in its own words and is turned into the twenty numbers the matrix
// takes, and the answer is the one those numbers give.
func TestFilterColorMatrixShortcutsWriteTheMatrixTheSpecGivesThem(t *testing.T) {
	for _, tc := range []struct {
		prim  string
		want  [3]uint8
		wantA uint8
	}{{
		`<feColorMatrix type="saturate" values="0"/>`,
		[3]uint8{54, 54, 54}, 255,
	}, {
		`<feColorMatrix type="hueRotate" values="180"/>`,
		[3]uint8{0, 109, 109}, 255,
	}, {
		`<feColorMatrix type="luminanceToAlpha"/>`,
		[3]uint8{0, 0, 0}, 54,
	}} {
		_, cv := painted(t, `<svg viewBox="0 0 10 10">
			<filter id="f">`+tc.prim+`</filter>
			<rect width="10" height="10" fill="#ff0000" filter="url(#f)"/>
		</svg>`, 10, 10)

		c := pixelAt(cv, 5, 5)
		if !nearByte(c.R(), tc.want[0]) || !nearByte(c.G(), tc.want[1]) ||
			!nearByte(c.B(), tc.want[2]) || !nearByte(c.A(), tc.wantA) {
			t.Errorf("%s at (5,5) is %#08x, want (%d,%d,%d) of %d", tc.prim, uint32(c), tc.want[0], tc.want[1], tc.want[2], tc.wantA)
		}
	}
}

// TestFilterTakesTheShapeWithoutTheColourForSourceAlpha is what SourceAlpha
// names: the picture as only where something is, with the colour taken off —
// so a matrix that changes nothing leaves the rectangle's shape in black
// rather than the red that was there.
func TestFilterTakesTheShapeWithoutTheColourForSourceAlpha(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 15 15">
		<filter id="f">
			<feColorMatrix in="SourceAlpha" values="1 0 0 0 0  0 1 0 0 0  0 0 1 0 0  0 0 0 1 0"/>
		</filter>
		<rect width="10" height="10" fill="#ff0000" filter="url(#f)"/>
	</svg>`, 15, 15)

	filteredAt(t, cv, 5, 5, 0, 0, 0)
	gone(t, cv, 13, 13)
}

// TestParseSaysWhatItCannotReadInAFilterPrimitive is every way a primitive can
// ask for something this package does not do: a number that is not one, a
// matrix that is not twenty numbers, a kind that is not one it knows, an
// input naming nothing the filter has, and the background no filter can see
// here. Each says what is wrong and once, and the drawing still paints.
func TestParseSaysWhatItCannotReadInAFilterPrimitive(t *testing.T) {
	for _, tc := range []struct{ prim, want string }{
		{`<feGaussianBlur stdDeviation="zzz"/>`, "stdDeviation"},
		{`<feGaussianBlur stdDeviation="-1"/>`, "negative"},
		{`<feGaussianBlur stdDeviation="1 -2"/>`, "negative"},
		{`<feColorMatrix type="matrix" values="1 2 3"/>`, "twenty numbers"},
		{`<feColorMatrix type="nonsense"/>`, `"nonsense"`},
		{`<feOffset dx="zzz"/>`, "feOffset dx"},
		{`<feOffset in="nope"/>`, "names nothing"},
		{`<feOffset in="BackgroundImage"/>`, "not where a filter can see"},
		{`<feFlood flood-color="zzz"/>`, "black is flooded instead"},
		{`<feFlood flood-opacity="zzz"/>`, "flood-opacity"},
		{`<feComposite operator="fold"/>`, `"fold"`},
		{`<feComposite k2="zzz"/>`, "k2"},
		{`<feBlend mode="peel"/>`, `"peel"`},
	} {
		img, err := Parse(`<svg viewBox="0 0 10 10"><filter id="f">` + tc.prim + `</filter><rect width="10" height="10" filter="url(#f)"/></svg>`)
		if err != nil {
			t.Fatalf("primitive %s: parse: %v", tc.prim, err)
		}
		ws := img.Warnings()
		if len(ws) != 1 {
			t.Errorf("primitive %s said %v, want it said once", tc.prim, ws)
		} else if !strings.Contains(ws.String(), tc.want) {
			t.Errorf("primitive %s said %v, want it to say %q", tc.prim, ws, tc.want)
		}
		if cv := img.Render(10, 10); cv == nil {
			t.Errorf("primitive %s: the drawing did not paint", tc.prim)
		}
	}
}
