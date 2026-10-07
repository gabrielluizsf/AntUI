package svg

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// A `<use>` puts a shape that was written once somewhere else in the picture,
// which only works if what it names is found, painted where the `<use>` is and
// not where it was written, and given the style of the place it is used rather
// than the style it was written in.

// TestRenderPaintsWhatAUsePointsAt is the whole of it: the shape lives in a
// `<defs>`, is drawn nowhere, and comes out where the `<use>` put it.
func TestRenderPaintsWhatAUsePointsAt(t *testing.T) {
	img, cv := painted(t, `<svg viewBox="0 0 20 20">
		<defs><rect id="r" x="0" y="0" width="4" height="4" fill="red"/></defs>
		<use href="#r" x="10" y="10"/>
	</svg>`, 20, 20)
	if len(img.Warnings()) != 0 {
		t.Errorf("the drawing said %v, want it happy with the use", img.Warnings())
	}
	if got := pixelAt(cv, 12, 12); got != canvas.RGB(255, 0, 0) {
		t.Errorf("where the use put it is %v, want it red", got)
	}
	if got := pixelAt(cv, 2, 2); got.A() != 0 {
		t.Errorf("where it was written is %v, want it clear: a defs draws nothing", got)
	}
	// The same thing asked for in the tree rather than in the picture: the
	// `<use>` is there, and the rectangle is inside it rather than in the defs.
	use := img.Root.find("use")
	if use == nil {
		t.Fatal("the use is not in the drawing")
	}
	if len(use.Kids) != 1 || use.Kids[0].Name != "rect" {
		t.Fatalf("the use holds %d nodes, want the one rectangle it names", len(use.Kids))
	}
	if b := boundsOf(use.Kids[0].Path); b.MinX != 10 || b.MinY != 10 {
		t.Errorf("the rectangle drawn through the use is at %v, want it moved to 10,10", b)
	}
}

// TestRenderDoesNotPaintWhatIsInsideDefs is the half of it that does not need
// a use: a defs is not drawn, so that a shape written there can be pointed at
// as often as it likes without coming out twice.
func TestRenderDoesNotPaintWhatIsInsideDefs(t *testing.T) {
	img, cv := painted(t, `<svg viewBox="0 0 20 20">
		<defs><rect id="r" width="20" height="20" fill="red"/></defs>
	</svg>`, 20, 20)
	if len(img.Warnings()) != 0 {
		t.Errorf("the drawing said %v, want nothing said about a defs", img.Warnings())
	}
	if got := pixelAt(cv, 10, 10); got.A() != 0 {
		t.Errorf("the middle is %v, want it clear: a defs is not part of the picture", got)
	}
}

// TestUseKeepsWhatItNamesWhereItWasWritten is the other case: a `<use>` may
// name a shape that is part of the picture, and then the picture has it twice
// — where it was written and where the `<use>` put it.
func TestUseKeepsWhatItNamesWhereItWasWritten(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 30 20">
		<rect id="r" x="0" y="0" width="4" height="4" fill="red"/>
		<use href="#r" x="10" y="10"/>
	</svg>`, 30, 20)
	if got := pixelAt(cv, 2, 2); got != canvas.RGB(255, 0, 0) {
		t.Errorf("where it was written is %v, want it red", got)
	}
	if got := pixelAt(cv, 12, 12); got != canvas.RGB(255, 0, 0) {
		t.Errorf("where the use put it is %v, want it red", got)
	}
}

// TestUseTakesTheStyleOfWhereItIsUsed is what makes one definition serve every
// place it appears: the shape inherits from the `<use>` and the elements above
// it, never from where it was written. Its own fill still wins over both.
func TestUseTakesTheStyleOfWhereItIsUsed(t *testing.T) {
	for _, tc := range []struct {
		name string
		svg  string
		want canvas.Color
	}{
		{
			name: "the group above the use colours it",
			svg: `<defs><circle id="c" cx="10" cy="10" r="5"/></defs>
				<g fill="red"><use href="#c"/></g>`,
			want: canvas.RGB(255, 0, 0),
		},
		{
			name: "its own fill wins over the group",
			svg: `<defs><circle id="c" cx="10" cy="10" r="5" fill="blue"/></defs>
				<g fill="red"><use href="#c"/></g>`,
			want: canvas.RGB(0, 0, 255),
		},
		{
			name: "the use itself colours it",
			svg: `<defs><circle id="c" cx="10" cy="10" r="5"/></defs>
				<use href="#c" fill="lime"/>`,
			want: canvas.RGB(0, 255, 0),
		},
	} {
		_, cv := painted(t, `<svg viewBox="0 0 20 20">`+tc.svg+`</svg>`, 20, 20)
		if got := pixelAt(cv, 10, 10); got != tc.want {
			t.Errorf("%s: the shape is %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestUsePutsItsXAndYInsideItsTransform is the order of the two: the x and y
// move what is referenced in its own space first, and the transform on the
// `<use>` then moves all of it. The other order puts the rectangle somewhere
// else entirely, and the two places are far enough apart to tell apart.
func TestUsePutsItsXAndYInsideItsTransform(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 40 40">
		<defs><rect id="r" x="0" y="0" width="5" height="5" fill="red"/></defs>
		<use href="#r" x="10" transform="scale(2)"/>
	</svg>`, 40, 40)
	// translate(10,0) first and scale(2) over it: 0..5 becomes 20..30.
	if got := pixelAt(cv, 25, 5); got != canvas.RGB(255, 0, 0) {
		t.Errorf("where the two of them together put it is %v, want it red", got)
	}
	if got := pixelAt(cv, 15, 5); got.A() != 0 {
		t.Errorf("where the other order would have put it is %v, want it clear", got)
	}
}

// TestUseDrawsAGroupWhereItIsAsked is the case a defs is really for: a group
// of shapes written once and put in the picture as many times as it is asked
// for, each time where that use asked for it.
func TestUseDrawsAGroupWhereItIsAsked(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 30 20">
		<defs>
			<g id="pair">
				<rect x="0" y="0" width="3" height="3" fill="red"/>
				<rect x="5" y="0" width="3" height="3" fill="blue"/>
			</g>
		</defs>
		<use href="#pair"/>
		<use href="#pair" x="10"/>
	</svg>`, 30, 20)
	for _, tc := range []struct {
		x, y int
		want canvas.Color
	}{
		{1, 1, canvas.RGB(255, 0, 0)},
		{6, 1, canvas.RGB(0, 0, 255)},
		{11, 1, canvas.RGB(255, 0, 0)},
		{16, 1, canvas.RGB(0, 0, 255)},
	} {
		if got := pixelAt(cv, tc.x, tc.y); got != tc.want {
			t.Errorf("at %d,%d is %v, want %v", tc.x, tc.y, got, tc.want)
		}
	}
}

// TestUseDrawsASymbolWhereItIsAsked: a symbol is drawn nowhere it stands and
// only through a use, and there its own viewBox is mapped onto the width and
// height the use gave it, at the x and y it gave it.
func TestUseDrawsASymbolWhereItIsAsked(t *testing.T) {
	img, cv := painted(t, `<svg viewBox="0 0 40 40">
		<defs>
			<symbol id="s" viewBox="0 0 10 10">
				<rect x="0" y="0" width="10" height="10" fill="red"/>
			</symbol>
		</defs>
		<use href="#s" x="10" y="10" width="10" height="10"/>
	</svg>`, 40, 40)
	if len(img.Warnings()) != 0 {
		t.Errorf("the drawing said %v, want it happy with the symbol", img.Warnings())
	}
	if got := pixelAt(cv, 15, 15); got != canvas.RGB(255, 0, 0) {
		t.Errorf("inside the symbol is %v, want it red", got)
	}
	if got := pixelAt(cv, 5, 5); got.A() != 0 {
		t.Errorf("where the symbol stands is %v, want it clear: a symbol draws nothing alone", got)
	}
}

// TestUseScalesASymbolByTheSizeItIsAsked is the other half of the same thing:
// the use asks for a box twice the size of the viewBox the symbol was drawn
// with, so the symbol comes out twice as big rather than at the size it was
// written at.
func TestUseScalesASymbolByTheSizeItIsAsked(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 40 40">
		<defs>
			<symbol id="s" viewBox="0 0 10 10">
				<rect x="0" y="0" width="10" height="10" fill="red"/>
			</symbol>
		</defs>
		<use href="#s" width="20" height="20"/>
	</svg>`, 40, 40)
	if got := pixelAt(cv, 15, 15); got != canvas.RGB(255, 0, 0) {
		t.Errorf("outside the size it was written at is %v, want it red: the symbol was scaled up", got)
	}
}

// TestUseDrawsNothingForASizeOfNothing is the one size SVG says draws nothing:
// a use whose width or height is zero has no box to put anything in.
func TestUseDrawsNothingForASizeOfNothing(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 40 40">
		<defs>
			<symbol id="s" viewBox="0 0 10 10">
				<rect x="0" y="0" width="10" height="10" fill="red"/>
			</symbol>
		</defs>
		<use href="#s" width="0" height="20"/>
	</svg>`, 40, 40)
	if got := pixelAt(cv, 10, 10); got.A() != 0 {
		t.Errorf("a use of no size painted %v, want it clear", got)
	}
}

// TestUseReadsTheOlderSpellingOfTheHref: a drawing written for SVG 1.1 says
// xlink:href, one written for SVG 2 says href, and both name the same element.
func TestUseReadsTheOlderSpellingOfTheHref(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 20 20" xmlns:xlink="http://www.w3.org/1999/xlink">
		<defs><rect id="r" x="0" y="0" width="4" height="4" fill="red"/></defs>
		<use xlink:href="#r" x="10"/>
	</svg>`, 20, 20)
	if got := pixelAt(cv, 12, 2); got != canvas.RGB(255, 0, 0) {
		t.Errorf("where the xlink put it is %v, want it red", got)
	}
}

// TestUseSaysWhatItCannotDraw: a reference that cannot be followed is a hole
// in the picture, and a hole is worth saying something about. Nothing is
// painted for any of them, since there is nothing to paint.
func TestUseSaysWhatItCannotDraw(t *testing.T) {
	for _, src := range []string{
		`<svg viewBox="0 0 10 10"><use/></svg>`,
		`<svg viewBox="0 0 10 10"><use href="#nowhere"/></svg>`,
		`<svg viewBox="0 0 10 10"><use href="other.svg#r"/></svg>`,
		`<svg viewBox="0 0 10 10"><use href="#"/></svg>`,
	} {
		img, err := Parse(src)
		if err != nil {
			t.Errorf("parse %s: %v", src, err)
			continue
		}
		if len(img.Warnings()) == 0 {
			t.Errorf("parse %s: it said nothing about the reference it could not follow", src)
			continue
		}
		if got := pixelAt(img.Render(10, 10), 2, 2); got.A() != 0 {
			t.Errorf("parse %s: it painted %v, want nothing where there is nothing to draw", src, got)
		}
	}
}

// TestUseDoesNotGoRoundInCircles: a use that names itself and two that name
// each other are circles rather than shapes. SVG calls them an error; what
// matters here is that the drawing is read rather than read for ever, and that
// it says so instead of hanging.
func TestUseDoesNotGoRoundInCircles(t *testing.T) {
	for _, src := range []string{
		`<svg viewBox="0 0 10 10"><use id="u" href="#u"/></svg>`,
		`<svg viewBox="0 0 10 10">
			<g id="a"><use href="#b"/></g>
			<g id="b"><use href="#a"/></g>
		</svg>`,
	} {
		img, err := Parse(src)
		if err != nil {
			t.Errorf("parse %s: %v", src, err)
			continue
		}
		if len(img.Warnings()) == 0 {
			t.Errorf("parse %s: it said nothing about the circle in the references", src)
		}
		_ = img.Render(10, 10)
	}
}

// TestUseDrawsOneShapeTwiceWithoutComplaining is the other side of the circle:
// two uses of one shape are two shapes, not a circle, and neither of them is
// worth a word.
func TestUseDrawsOneShapeTwiceWithoutComplaining(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 20 10">
		<defs><rect id="r" width="5" height="5" fill="red"/></defs>
		<use href="#r"/><use href="#r"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(img.Warnings()) != 0 {
		t.Errorf("the drawing said %v, want nothing said about using one shape twice", img.Warnings())
	}
	cv := img.Render(20, 10)
	// The second use draws over the first, which is what two of the same shape
	// at the same place look like: one shape, still there.
	if got := pixelAt(cv, 2, 2); got != canvas.RGB(255, 0, 0) {
		t.Errorf("the shape used twice is %v, want it red", got)
	}
}

// TestUseDrawsNothingWhenItIsHidden: display:none on the use hides what it
// would have drawn, the same as it hides anything else, and nothing behind it
// is built to be drawn.
func TestUseDrawsNothingWhenItIsHidden(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 20 20">
		<defs><rect id="r" width="20" height="20" fill="red"/></defs>
		<use href="#r" display="none"/>
	</svg>`, 20, 20)
	if got := pixelAt(cv, 10, 10); got.A() != 0 {
		t.Errorf("a hidden use painted %v, want it clear", got)
	}
}

// TestUseDrawsASymbolWithNoViewBox: a symbol with no viewBox has no size of
// its own to map onto the box the use gave it, so it keeps the coordinates it
// was written in and is only moved.
func TestUseDrawsASymbolWithNoViewBox(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 20 20">
		<defs><symbol id="s"><rect width="4" height="4" fill="red"/></symbol></defs>
		<use href="#s" x="10" y="10" width="4" height="4"/>
	</svg>`, 20, 20)
	if got := pixelAt(cv, 12, 12); got != canvas.RGB(255, 0, 0) {
		t.Errorf("where the use put it is %v, want it red", got)
	}
	if got := pixelAt(cv, 2, 2); got.A() != 0 {
		t.Errorf("where it was written is %v, want it clear", got)
	}
}

// TestUseFitsASvgLikeASymbol: a drawing named by a use is a viewport of its
// own, so its viewBox is mapped onto the box the use gave it the same way a
// symbol's is — the use is what says how big it comes out.
func TestUseFitsASvgLikeASymbol(t *testing.T) {
	img, cv := painted(t, `<svg viewBox="0 0 40 40">
		<defs><svg id="v" viewBox="0 0 10 10"><rect width="10" height="10" fill="red"/></svg></defs>
		<use href="#v" width="20" height="20"/>
	</svg>`, 40, 40)
	if len(img.Warnings()) != 0 {
		t.Errorf("the drawing said %v, want it happy with the svg", img.Warnings())
	}
	if got := pixelAt(cv, 15, 15); got != canvas.RGB(255, 0, 0) {
		t.Errorf("inside the box the use asked for is %v, want it red: the viewBox was fitted", got)
	}
	if got := pixelAt(cv, 25, 25); got.A() != 0 {
		t.Errorf("outside the box the use asked for is %v, want it clear", got)
	}
}

// TestUseTakesTheSizeOfTheSvgItNames is the size of the drawing itself: a use
// that asks for no box leaves the nested drawing at the width and height it was
// written with, rather than blowing it up to the whole picture.
func TestUseTakesTheSizeOfTheSvgItNames(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 40 40">
		<defs><svg id="v" width="10" height="10" viewBox="0 0 10 10"><rect width="10" height="10" fill="red"/></svg></defs>
		<use href="#v"/>
	</svg>`, 40, 40)
	if got := pixelAt(cv, 5, 5); got != canvas.RGB(255, 0, 0) {
		t.Errorf("inside the size it was written at is %v, want it red", got)
	}
	if got := pixelAt(cv, 15, 5); got.A() != 0 {
		t.Errorf("outside the size it was written at is %v, want it clear", got)
	}
}

// TestUseSizeWinsOverTheSizeOfTheSvgItNames: the use asks, so the use decides,
// and the box the drawing was written with is only what happens when it does not.
func TestUseSizeWinsOverTheSizeOfTheSvgItNames(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 40 40">
		<defs><svg id="v" width="10" height="10" viewBox="0 0 10 10"><rect width="10" height="10" fill="red"/></svg></defs>
		<use href="#v" width="30" height="30"/>
	</svg>`, 40, 40)
	if got := pixelAt(cv, 25, 25); got != canvas.RGB(255, 0, 0) {
		t.Errorf("inside the box the use asked for is %v, want it red", got)
	}
	if got := pixelAt(cv, 35, 35); got.A() != 0 {
		t.Errorf("outside the box the use asked for is %v, want it clear", got)
	}
}
