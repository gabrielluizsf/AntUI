package svg

import (
	"strings"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// A mask is a picture that says how much of an element is there: it is
// written once, every element that names it is drawn into a picture of its
// own, that picture is taken down by how much of the mask is over each of its
// pixels, and the result is laid over what is behind. What is here is what
// that comes out as. Points are asked for in the drawing's own units, which is
// what a viewBox of the same numbers makes the pixels worth.

// keptAmount is how much of the drawing came out at a point, which is what a
// mask leaves rather than whether it left it: 255 where the mask keeps
// everything, 0 where it keeps nothing, and a number in between where it keeps
// some of it. The two out of 255 is the rounding a picture at twice the size
// and halved again does on the way out.
func keptAmount(t *testing.T, cv *canvas.Canvas, x, y, want int) {
	t.Helper()
	if got := alphaAt(cv, x, y); got > want+2 || got < want-2 {
		t.Errorf("at (%d,%d) the drawing is %d out of 255, want %d", x, y, got, want)
	}
}

// TestRenderMasksAnElementWithWhatTheMaskKeeps is the whole of what a mask
// does: the element fills the drawing, the mask is white over a third of it,
// grey over the next and black over the last, and what comes out is the
// element whole, the element half there and the element nowhere — with its
// colour where any of it is left, because what a mask takes away is how much
// of a pixel there is and not what colour it is.
func TestRenderMasksAnElementWithWhatTheMaskKeeps(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 9 10">
		<rect width="9" height="10" fill="#ff0000" mask="url(#m)"/>
		<mask id="m">
			<rect x="0" y="0" width="3" height="10" fill="#ffffff"/>
			<rect x="3" y="0" width="3" height="10" fill="#808080"/>
			<rect x="6" y="0" width="3" height="10" fill="#000000"/>
		</mask>
	</svg>`, 9, 10)

	kept(t, cv, 1, 5)
	if got := pixelAt(cv, 1, 5); got != canvas.RGB(255, 0, 0) {
		t.Errorf("the kept pixel is %#08x, want the red the element was painted", uint32(got))
	}
	keptAmount(t, cv, 4, 5, 128)
	gone(t, cv, 7, 5)
	// The `<mask>` is not part of the picture where it stands: had it been,
	// its black band would have been painted over the element rather than
	// measured against it, and there would be a colour there to read.
	for _, p := range [][2]int{{1, 1}, {4, 8}, {7, 1}, {7, 8}} {
		if got := pixelAt(cv, p[0], p[1]); got.B() != 0 {
			t.Errorf("at %v the drawing is %#08x, want no blue: nothing but the element is painted", p, uint32(got))
		}
	}
}

// TestRenderMeasuresAMaskByWhatItsMaskTypeSays is the difference between the
// two measures, and it is a black band: measured by luminance a black mask is
// a hole, measured by alpha the same opaque black keeps everything and the
// only thing that takes anything away is what is clear. The grey band is that
// difference the other way round — by luminance it keeps only as much of the
// element as it is bright, by alpha it keeps all of it, being opaque — while
// half of a white band comes out the same either way, since both measures ask
// how much white is there rather than how bright it is. So does the part of
// the element no mask reaches at all — which is nowhere, since what no mask
// covers keeps nothing.
func TestRenderMeasuresAMaskByWhatItsMaskTypeSays(t *testing.T) {
	for _, tc := range []struct {
		name     string
		maskType string
		want     [5]int
	}{
		{name: "luminance where nothing says", maskType: "", want: [5]int{255, 128, 0, 128, 0}},
		{name: "luminance said as an attribute", maskType: ` mask-type="luminance"`, want: [5]int{255, 128, 0, 128, 0}},
		{name: "alpha said as an attribute", maskType: ` mask-type="alpha"`, want: [5]int{255, 255, 255, 128, 0}},
		{name: "alpha in the style", maskType: ` style="mask-type: alpha"`, want: [5]int{255, 255, 255, 128, 0}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cv := cut(t, `<svg viewBox="0 0 15 10">
				<rect width="15" height="10" fill="#ff0000" mask="url(#m)"/>
				<mask id="m"`+tc.maskType+`>
					<rect x="0" y="0" width="3" height="10" fill="#ffffff"/>
					<rect x="3" y="0" width="3" height="10" fill="#808080"/>
					<rect x="6" y="0" width="3" height="10" fill="#000000"/>
					<rect x="9" y="0" width="3" height="10" fill="#ffffff" fill-opacity="0.5"/>
				</mask>
			</svg>`, 15, 10)

			for i, x := range []int{1, 4, 7, 10, 13} {
				keptAmount(t, cv, x, 5, tc.want[i])
			}
		})
	}
}

// TestRenderSaysSoWhenAMaskTypeIsNotOneItKnows: the two measures disagree
// about the whole of a black mask, so an author who asked for one of them is
// not given the other in silence — the drawing says what it could not read,
// and the mask measures how bright it is, which is what it means anyway.
func TestRenderSaysSoWhenAMaskTypeIsNotOneItKnows(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 10 10">
		<rect width="10" height="10" fill="#ff0000" mask="url(#m)"/>
		<mask id="m" mask-type="shining">
			<rect width="10" height="10" fill="#000000"/>
		</mask>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	ws := img.Warnings()
	if len(ws) != 1 || !strings.Contains(ws.String(), "mask-type") {
		t.Fatalf("the drawing said %v, want it to say what mask-type it could not read", ws)
	}
	gone(t, img.Render(10, 10), 5, 5)
}

// TestRenderOfAMaskThatIsNotThereKeepsNothingAndSaysSoOnce is the rest of
// what a reference to nowhere means here: what is masked cannot be drawn, and
// the group says it once for itself rather than once for every shape inside
// it — the same as a clip that names a clipPath the drawing does not have.
func TestRenderOfAMaskThatIsNotThereKeepsNothingAndSaysSoOnce(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 10 10">
		<g mask="url(#nowhere)">
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

// TestRenderMasksAGroupAsOnePicture: the shapes under a group's mask are one
// picture with one measure over it, so where two half-transparent rectangles
// overlap they overlap first and the whole of that is measured — rather than
// each being measured on its own and the two halves meeting somewhere the mask
// does not cover at all.
func TestRenderMasksAGroupAsOnePicture(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 12 10">
		<g mask="url(#m)">
			<rect x="0" y="0" width="6" height="10" fill="#ff0000" fill-opacity="0.5"/>
			<rect x="4" y="0" width="6" height="10" fill="#ff0000" fill-opacity="0.5"/>
		</g>
		<mask id="m"><rect x="0" y="0" width="8" height="10" fill="#ffffff"/></mask>
	</svg>`, 12, 10)

	// Where only one rectangle reaches, the half of one: the mask keeps all of
	// it, so what comes out is what the two rectangles put there.
	if a := alphaAt(cv, 1, 5); a < 124 || a > 132 {
		t.Errorf("at (1,5) the drawing is %d out of 255, want the half of one rectangle", a)
	}
	// Where both reach, one over the other: three quarters of the whole.
	if a := alphaAt(cv, 5, 5); a < 188 || a > 195 {
		t.Errorf("at (5,5) the drawing is %d out of 255, want one rectangle over the other", a)
	}
	// The mask stops at 8, so past it nothing of the group is there however
	// much of it painted there.
	for _, p := range [][2]int{{9, 1}, {11, 5}, {9, 9}} {
		gone(t, cv, p[0], p[1])
	}
}

// TestRenderMovesTheMaskWithTheElementItCuts: the mask is written in the
// coordinates of the element that names it, so a transform on that element
// moves the mask along with it. Here the two only meet because of that — as
// written, the rectangle starts where the mask stops.
func TestRenderMovesTheMaskWithTheElementItCuts(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 15 10">
		<rect width="10" height="10" transform="translate(5,0)" fill="#ff0000" mask="url(#m)"/>
		<mask id="m"><rect x="0" y="0" width="5" height="10" fill="#ffffff"/></mask>
	</svg>`, 15, 10)

	for _, p := range [][2]int{{5, 5}, {7, 5}, {9, 5}} {
		kept(t, cv, p[0], p[1])
	}
	for _, p := range [][2]int{{4, 5}, {11, 5}, {14, 5}} {
		gone(t, cv, p[0], p[1])
	}
}

// TestRenderDoesNotMaskWhatComesAfterTheElementTheMaskIsOn: the mask is not
// inherited, so it is laid over the element that named it and nothing else —
// the shape written next stands whole, even where the masked one was taken
// away.
func TestRenderDoesNotMaskWhatComesAfterTheElementTheMaskIsOn(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<rect width="10" height="10" fill="#ff0000" mask="url(#m)"/>
		<rect x="6" y="0" width="4" height="10" fill="#0000ff"/>
		<mask id="m"><rect x="0" y="0" width="5" height="10" fill="#ffffff"/></mask>
	</svg>`, 10, 10)

	kept(t, cv, 2, 5)
	if got := pixelAt(cv, 7, 5); got != canvas.RGB(0, 0, 255) {
		t.Errorf("the shape after the masked one is %v, want it blue and whole", got)
	}
	gone(t, cv, 5, 7)
}

// TestRenderOfAnEmptyMaskKeepsNothing: a `<mask>` with no picture in it is a
// measure that says nothing about anything, and what no mask has anything to
// say about keeps nothing — which is the same as a clipPath of no shapes.
func TestRenderOfAnEmptyMaskKeepsNothing(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<rect width="10" height="10" fill="#ff0000" mask="url(#m)"/>
		<mask id="m"/>
	</svg>`, 10, 10)

	for _, p := range [][2]int{{1, 1}, {5, 5}, {9, 9}} {
		gone(t, cv, p[0], p[1])
	}
}

// TestRenderOfAMaskOfNoneDrawsTheElementWhole: `none` is a mask that is not
// there, which is read as no mask and not as a complaint.
func TestRenderOfAMaskOfNoneDrawsTheElementWhole(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<rect width="10" height="10" fill="#ff0000" mask="none"/>
	</svg>`, 10, 10)

	for _, p := range [][2]int{{0, 0}, {5, 5}, {9, 9}} {
		kept(t, cv, p[0], p[1])
	}
}

// TestRenderOfAMaskInFractionsOfABoxDrawsTheElementWhole: a mask whose picture
// is written against the box of the element it is put on is laid out over
// that box — the box measured the way a gradient in fractions of it is — and
// here the picture is white over the whole of it, so what comes out is the
// element whole, with nothing said about it.
func TestRenderOfAMaskInFractionsOfABoxDrawsTheElementWhole(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<rect width="10" height="10" fill="#ff0000" mask="url(#m)"/>
		<mask id="m" maskContentUnits="objectBoundingBox">
			<rect width="1" height="1" fill="#ffffff"/>
		</mask>
	</svg>`, 10, 10)

	for _, p := range [][2]int{{1, 1}, {5, 5}, {9, 9}} {
		kept(t, cv, p[0], p[1])
	}
}

// TestRenderTakesHalfTheElementDownByAMaskInFractionsOfABox: the same picture
// written against the box of the element in halves of it keeps what the left
// half covers and takes down what the right half leaves out. The element is
// asked for once where it stands and once after a transform moved it: the box
// is the one it paints at either way, so content and region stay over the
// element together rather than the content being left where it was written.
func TestRenderTakesHalfTheElementDownByAMaskInFractionsOfABox(t *testing.T) {
	for _, tc := range []struct {
		name string
		tr   string
		kept [][2]int
		gone [][2]int
	}{
		{
			name: "where the element stands",
			kept: [][2]int{{1, 1}, {4, 5}, {4, 8}},
			gone: [][2]int{{6, 5}, {9, 1}, {9, 9}},
		},
		{
			name: "where a transform moved it",
			tr:   ` transform="translate(5,0)"`,
			kept: [][2]int{{6, 1}, {9, 5}, {9, 8}},
			gone: [][2]int{{12, 5}, {14, 1}, {14, 9}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cv := cut(t, `<svg viewBox="0 0 15 10">
				<rect width="10" height="10"`+tc.tr+` fill="#ff0000" mask="url(#m)"/>
				<mask id="m" maskContentUnits="objectBoundingBox">
					<rect width="0.5" height="1" fill="#ffffff"/>
				</mask>
			</svg>`, 15, 10)

			for _, p := range tc.kept {
				kept(t, cv, p[0], p[1])
			}
			for _, p := range tc.gone {
				gone(t, cv, p[0], p[1])
			}
		})
	}
}

// TestRenderOfAMaskInFractionsOfABoxOnAnElementWithNoBoxSaysSoAndDrawsWithoutIt:
// a group with nothing in it that paints has no box to lay fractions of, and
// a mask that cannot be laid out is left off altogether with one warning
// saying so — the element drawn whole is the same answer a clip gives where
// its box cannot be measured.
func TestRenderOfAMaskInFractionsOfABoxOnAnElementWithNoBoxSaysSoAndDrawsWithoutIt(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 10 10">
		<g mask="url(#m)"/>
		<mask id="m" maskContentUnits="objectBoundingBox">
			<rect width="1" height="1"/>
		</mask>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	ws := img.Warnings()
	if len(ws) != 1 || !strings.Contains(ws.String(), "no box to measure") {
		t.Fatalf("the drawing said %v, want it to say why the mask was left off", ws)
	}
	if cv := img.Render(10, 10); cv == nil {
		t.Fatal("the drawing did not paint")
	}
}

// TestRenderCutsAMaskToTheRegionItReaches: the mask is only allowed to say
// anything inside the region it was given, and what it says outside of it is
// beyond the region rather than under it — so a picture of white over the
// whole drawing still takes the element away past where the region stops. The
// region may be written in the drawing's own coordinates or in fractions of
// the box of the element, which is what `maskUnits` chooses between and what
// the default of the box grown by a tenth all round is a fraction of.
func TestRenderCutsAMaskToTheRegionItReaches(t *testing.T) {
	for _, tc := range []struct {
		name   string
		region string
	}{
		{"in the drawing's own coordinates", ` maskUnits="userSpaceOnUse" x="0" y="0" width="5" height="10"`},
		{"in fractions of the box", ` maskUnits="objectBoundingBox" x="0" y="0" width="0.5" height="1"`},
		{"in percentages of the box", ` maskUnits="objectBoundingBox" x="0%" y="0%" width="50%" height="100%"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cv := cut(t, `<svg viewBox="0 0 10 10">
				<rect width="10" height="10" fill="#ff0000" mask="url(#m)"/>
				<mask id="m"`+tc.region+`>
					<rect width="10" height="10" fill="#ffffff"/>
				</mask>
			</svg>`, 10, 10)

			for _, p := range [][2]int{{1, 1}, {4, 5}, {4, 8}} {
				kept(t, cv, p[0], p[1])
			}
			for _, p := range [][2]int{{6, 5}, {9, 1}, {9, 9}} {
				gone(t, cv, p[0], p[1])
			}
		})
	}
}

// TestRenderMovesTheMaskOfAUseWithItsXAndY: a `<use>` is drawn where its x
// and y put it, and the mask goes there too — the mask is followed after those
// have moved the style, so the two start from the same place rather than from
// the corner the drawing was written in, exactly as a clip does.
func TestRenderMovesTheMaskOfAUseWithItsXAndY(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 20 10">
		<defs>
			<mask id="m"><rect x="0" y="0" width="5" height="10" fill="#ffffff"/></mask>
			<rect id="r" width="15" height="10" fill="#ff0000"/>
		</defs>
		<use href="#r" x="5" mask="url(#m)"/>
	</svg>`, 20, 10)

	// The rectangle is drawn from 5 to 20 and the mask from 5 to 10, so what
	// comes out is the two where they meet — not nothing at all, which is what
	// a mask left in the corner the rectangle came from would make of it.
	for _, p := range [][2]int{{5, 5}, {7, 5}, {9, 5}} {
		kept(t, cv, p[0], p[1])
	}
	for _, p := range [][2]int{{1, 5}, {4, 5}, {11, 5}, {17, 5}} {
		gone(t, cv, p[0], p[1])
	}
}

// TestRenderOfAMaskThatNamesItselfStandsStill: a mask whose picture masks
// itself — and one pair of masks that name each other — would draw the same
// picture while it was already being drawn, for ever. The mask in force
// further out is the one that counts and the inner one is left off, so the
// drawing comes out as a picture rather than as a program that does not come
// back.
func TestRenderOfAMaskThatNamesItselfStandsStill(t *testing.T) {
	for _, tc := range []struct {
		name string
		src  string
	}{
		{"one mask naming itself", `<svg viewBox="0 0 10 10">
			<rect width="10" height="10" fill="#ff0000" mask="url(#a)"/>
			<mask id="a">
				<rect width="10" height="10" fill="#ffffff" mask="url(#a)"/>
			</mask>
		</svg>`},
		{"two masks naming each other", `<svg viewBox="0 0 10 10">
			<rect width="10" height="10" fill="#ff0000" mask="url(#a)"/>
			<mask id="a">
				<rect width="10" height="10" fill="#ffffff" mask="url(#b)"/>
			</mask>
			<mask id="b">
				<rect width="10" height="10" fill="#ffffff" mask="url(#a)"/>
			</mask>
		</svg>`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cv := cut(t, tc.src, 10, 10)
			for _, p := range [][2]int{{1, 1}, {5, 5}, {9, 9}} {
				kept(t, cv, p[0], p[1])
			}
		})
	}
}

// TestRenderCutsWritingToTheRegionOfAMaskInFractionsOfABox: writing paints
// with no box of its own until a font is asked, and the box the letters come
// out in is what the region written against it is a fraction of — so the
// region cuts the element the same as it cuts a shape, rather than reaching
// everywhere because there was no box to measure it against. Writing at this
// size comes out with every edge antialiased, so what is asked of the kept
// letter is that anything was painted on its stem rather than how much.
func TestRenderCutsWritingToTheRegionOfAMaskInFractionsOfABox(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 40 20">
		<text x="4" y="16" font-size="8" fill="#ff0000" mask="url(#m)">hi</text>
		<mask id="m" x="0" width="0.5">
			<rect width="40" height="20" fill="#ffffff"/>
		</mask>
	</svg>`, 40, 20)

	// The writing runs from 4 across to about 11, and the region keeps the
	// first half of that box: the `h` stands and the `i`, past the middle,
	// is taken away.
	if alphaAt(cv, 4, 15) == 0 {
		t.Error("at (4,15) the drawing is clear, want the letter at the start of the writing kept")
	}
	gone(t, cv, 10, 15)
}

// TestRenderTakesTheElementDownByAMaskPaintedLikeAPicture: the mask is not a
// cut of shapes but a picture of its own, so anything that can be drawn can
// measure — a mask drawn with a gradient fades the element the same way, and
// a mask with a transform of its own moves where it is measured against.
func TestRenderTakesTheElementDownByAMaskPaintedLikeAPicture(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<rect width="10" height="10" fill="#ff0000" mask="url(#m)"/>
		<defs>
			<linearGradient id="g" x1="0" y1="0" x2="10" y2="0" gradientUnits="userSpaceOnUse">
				<stop offset="0" stop-color="#ffffff"/>
				<stop offset="1" stop-color="#000000"/>
			</linearGradient>
			<mask id="m" maskUnits="userSpaceOnUse" x="0" y="0" width="10" height="10">
				<g transform="translate(10,0) scale(-1,1)">
					<rect width="10" height="10" fill="url(#g)"/>
				</g>
			</mask>
		</defs>
	</svg>`, 10, 10)

	// The gradient runs white to black and the group turns it round, so the
	// dark end of the ramp — which without the group's transform would be at
	// the right — takes the element away at the left, and the bright end
	// keeps it at the right. What comes out climbs from one to the other.
	lo, mid, hi := alphaAt(cv, 1, 5), alphaAt(cv, 5, 5), alphaAt(cv, 9, 5)
	if lo > 80 {
		t.Errorf("at (1,5) the drawing is %d out of 255, want nearly all of it taken away: the mask is dark there", lo)
	}
	if !(lo < mid && mid < hi) {
		t.Errorf("the drawing reads %d, %d, %d from left to right, want it climbing along the ramp", lo, mid, hi)
	}
	if hi < 185 {
		t.Errorf("at (9,5) the drawing is %d out of 255, want most of it kept: the mask is bright there", hi)
	}
}

// TestRenderOfAnUnreadableMaskSaysSoAndDrawsWithoutIt is the rest of what
// `mask` means when it is not a reference: the element is drawn as though
// there were no mask at all, with one warning saying what could not be read,
// rather than being drawn black or not at all because of a word this package
// did not follow.
func TestRenderOfAnUnreadableMaskSaysSoAndDrawsWithoutIt(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 10 10">
		<rect width="10" height="10" fill="#ff0000" mask="fading"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	ws := img.Warnings()
	if len(ws) != 1 || !strings.Contains(ws.String(), "fading") {
		t.Fatalf("the drawing said %v, want it to say what mask it could not follow", ws)
	}
	for _, p := range [][2]int{{1, 1}, {5, 5}, {9, 9}} {
		kept(t, img.Render(10, 10), p[0], p[1])
	}
}
