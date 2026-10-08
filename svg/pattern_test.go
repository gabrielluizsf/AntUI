package svg

import (
	"strings"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// A pattern is a picture laid over a shape over and over: the picture is drawn
// once into a tile, and every pixel of the shape asks the tile which part of it
// it is over, wrapped round by the size of the tile. What is here is what that
// comes out as. Points are asked for in the drawing's own units, which is what
// a viewBox of the same numbers makes the pixels worth, and the tests put every
// edge of every tile on a pixel edge so that what a point answers is the whole
// of the pixel rather than a bit of two halves of it.

// greenAt says the pixel came out of the tile's own picture, and clearAt says
// nothing at all was painted there — the two things a tiling of one small
// drawing over a shape is made of: where the picture reaches and where the tile
// lets the shape through.
func greenAt(t *testing.T, cv *canvas.Canvas, x, y int) {
	t.Helper()
	c := pixelAt(cv, x, y)
	if c.A() < 250 || c.G() < 200 || c.R() > 60 || c.B() > 60 {
		t.Errorf("at (%d,%d) the pixel is %#08x, want the green of the tile", x, y, uint32(c))
	}
}

func clearAt(t *testing.T, cv *canvas.Canvas, x, y int) {
	t.Helper()
	if a := alphaAt(cv, x, y); a != 0 {
		c := pixelAt(cv, x, y)
		t.Errorf("at (%d,%d) the pixel is %#08x, want nothing painted there", x, y, uint32(c))
	}
}

// blueAt is the colour after the `url(...)`: what a shape takes where the
// pattern it named cannot paint, which is the second chance the file gave it.
func blueAt(t *testing.T, cv *canvas.Canvas, x, y int) {
	t.Helper()
	c := pixelAt(cv, x, y)
	if c.A() < 250 || c.B() < 200 || c.R() > 60 || c.G() > 60 {
		t.Errorf("at (%d,%d) the pixel is %#08x, want the colour after the pattern", x, y, uint32(c))
	}
}

// nothingPainted says the whole drawing came out clear, which is every way a
// pattern can be pointed at and still leave the shape without a picture: a
// pattern that names itself, two that name each other, one that holds nothing.
func nothingPainted(t *testing.T, cv *canvas.Canvas) {
	t.Helper()
	for y := 0; y < cv.Height; y++ {
		for x := 0; x < cv.Width; x++ {
			if alphaAt(cv, x, y) != 0 {
				c := pixelAt(cv, x, y)
				t.Fatalf("at (%d,%d) the drawing is %#08x, want it clear", x, y, uint32(c))
			}
		}
	}
}

// TestRenderTilesAPatternAcrossTheShapeItFills is the whole of what a pattern
// is: a tile of four by four with a green square in the corner of it, laid over
// a shape of ten by ten so the same square comes back at every corner of every
// tile and the rest of the tile leaves the shape with nothing to show.
func TestRenderTilesAPatternAcrossTheShapeItFills(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<defs>
			<pattern id="p" patternUnits="userSpaceOnUse" width="4" height="4">
				<rect width="2" height="2" fill="#00ff00"/>
			</pattern>
		</defs>
		<rect width="10" height="10" fill="url(#p)"/>
	</svg>`, 10, 10)

	for _, pt := range [][2]int{{1, 1}, {5, 1}, {9, 1}, {1, 5}, {5, 5}, {9, 9}} {
		greenAt(t, cv, pt[0], pt[1])
	}
	for _, pt := range [][2]int{{3, 1}, {7, 1}, {1, 3}, {3, 3}, {7, 7}, {9, 7}} {
		clearAt(t, cv, pt[0], pt[1])
	}
}

// TestRenderTilesAPatternInTheBoxOfEveryShapeItFills is fractions of the box:
// the tile is written as half the width and all of the height of whatever shape
// the pattern paints, and what is inside the tile is written the same way, so
// the two shapes below — one twice as wide as the other — each get their own
// tiling out of the one pattern.
func TestRenderTilesAPatternInTheBoxOfEveryShapeItFills(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 20 10">
		<defs>
			<pattern id="p" width="0.5" height="1" patternContentUnits="objectBoundingBox">
				<rect width="1" height="0.5" fill="#00ff00"/>
			</pattern>
		</defs>
		<rect x="0" y="0" width="8" height="4" fill="url(#p)"/>
		<rect x="10" y="0" width="4" height="4" fill="url(#p)"/>
	</svg>`, 20, 10)

	for _, pt := range [][2]int{{1, 1}, {3, 1}, {5, 1}, {7, 1}, {11, 1}, {13, 1}} {
		greenAt(t, cv, pt[0], pt[1])
	}
	for _, pt := range [][2]int{{1, 3}, {5, 3}, {7, 3}, {11, 3}, {13, 3}, {9, 1}, {19, 9}} {
		clearAt(t, cv, pt[0], pt[1])
	}
}

// TestRenderMovesTheTileWithPatternTransform is the pattern's own transform:
// the shape stays where it is and the lattice under it slides, so the green
// square lands one unit to the right of where it would have been and the
// corner it used to fill comes out clear.
func TestRenderMovesTheTileWithPatternTransform(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<defs>
			<pattern id="p" patternUnits="userSpaceOnUse" width="4" height="4"
				patternTransform="translate(1 0)">
				<rect width="2" height="2" fill="#00ff00"/>
			</pattern>
		</defs>
		<rect width="10" height="10" fill="url(#p)"/>
	</svg>`, 10, 10)

	for _, pt := range [][2]int{{1, 1}, {2, 1}, {5, 1}, {6, 1}, {1, 5}} {
		greenAt(t, cv, pt[0], pt[1])
	}
	for _, pt := range [][2]int{{0, 1}, {3, 1}, {4, 1}, {7, 1}, {0, 5}} {
		clearAt(t, cv, pt[0], pt[1])
	}
}

// TestRenderTilesAPatternThroughARotation turns the lattice half round: the
// green square that was in the top left of every tile is in the bottom right
// of it instead, which is a point of asking how far along the tile a pixel is
// when it stands before the first tile rather than inside one.
func TestRenderTilesAPatternThroughARotation(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<defs>
			<pattern id="p" patternUnits="userSpaceOnUse" width="4" height="4"
				patternTransform="rotate(180)">
				<rect width="2" height="2" fill="#00ff00"/>
			</pattern>
		</defs>
		<rect width="10" height="10" fill="url(#p)"/>
	</svg>`, 10, 10)

	for _, pt := range [][2]int{{2, 2}, {3, 3}, {6, 6}, {7, 7}, {2, 6}, {6, 2}} {
		greenAt(t, cv, pt[0], pt[1])
	}
	for _, pt := range [][2]int{{1, 1}, {5, 5}, {1, 5}, {5, 1}, {4, 1}} {
		clearAt(t, cv, pt[0], pt[1])
	}
}

// TestRenderTilesAPatternWhereTheElementWent moves the shape instead: the tile
// is written where the rectangle was written, so the transform that put the
// rectangle on the page carries the tiling along with it rather than leaving
// the picture behind on the drawing board.
func TestRenderTilesAPatternWhereTheElementWent(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<defs>
			<pattern id="p" patternUnits="userSpaceOnUse" width="4" height="4">
				<rect width="2" height="2" fill="#00ff00"/>
			</pattern>
		</defs>
		<rect x="0" y="0" width="8" height="10" transform="translate(2 0)" fill="url(#p)"/>
	</svg>`, 10, 10)

	for _, pt := range [][2]int{{3, 1}, {6, 1}, {3, 5}, {7, 5}} {
		greenAt(t, cv, pt[0], pt[1])
	}
	for _, pt := range [][2]int{{4, 1}, {5, 1}, {11, 1}} {
		clearAt(t, cv, pt[0], pt[1])
	}
}

// TestRenderPaintsAPatternAlongAStroke widens the line into the shape it covers
// and paints that shape with the tile, so the green comes back along the stroke
// the same way it comes back across a fill, and the middle of the rectangle —
// which the stroke never reaches — has nothing at all.
func TestRenderPaintsAPatternAlongAStroke(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<defs>
			<pattern id="p" patternUnits="userSpaceOnUse" width="4" height="4">
				<rect width="2" height="2" fill="#00ff00"/>
			</pattern>
		</defs>
		<rect x="1" y="1" width="8" height="8" fill="none" stroke="url(#p)" stroke-width="2"/>
	</svg>`, 10, 10)

	for _, pt := range [][2]int{{1, 1}, {5, 1}, {1, 5}} {
		greenAt(t, cv, pt[0], pt[1])
	}
	for _, pt := range [][2]int{{3, 1}, {1, 3}, {5, 5}} {
		clearAt(t, cv, pt[0], pt[1])
	}
}

// TestRenderDoesNotFollowAPatternThatNamesItself is a pattern whose picture is
// the pattern again: following that would draw the tile while the tile is being
// drawn, for ever, so the shape is painted with no picture at all instead — a
// drawing that comes out empty rather than as a program that does not come
// back — and the drawing says which pattern closed the circle, once, the same
// way it says which marker would have gone round for ever.
func TestRenderDoesNotFollowAPatternThatNamesItself(t *testing.T) {
	img, cv := painted(t, `<svg viewBox="0 0 10 10">
		<defs>
			<pattern id="p" patternUnits="userSpaceOnUse" width="4" height="4">
				<rect width="4" height="4" fill="url(#p)"/>
			</pattern>
		</defs>
		<rect width="10" height="10" fill="url(#p)"/>
	</svg>`, 10, 10)

	if ws := img.Warnings(); len(ws) != 1 || !strings.Contains(ws.String(), "points back at the pattern") {
		t.Fatalf("the drawing said %v, want it to say which pattern closes the circle", ws)
	}
	nothingPainted(t, cv)
}

// TestRenderDoesNotFollowTwoPatternsThatNameEachOther is the same round trip
// taken in two steps, which is the way a drawing that never says a pattern
// names itself gets to naming one anyway: the first asks the second, the
// second asks the first, and the first is already being drawn.
func TestRenderDoesNotFollowTwoPatternsThatNameEachOther(t *testing.T) {
	img, cv := painted(t, `<svg viewBox="0 0 10 10">
		<defs>
			<pattern id="a" patternUnits="userSpaceOnUse" width="4" height="4">
				<rect width="4" height="4" fill="url(#b)"/>
			</pattern>
			<pattern id="b" patternUnits="userSpaceOnUse" width="4" height="4">
				<rect width="4" height="4" fill="url(#a)"/>
			</pattern>
		</defs>
		<rect width="10" height="10" fill="url(#a)"/>
	</svg>`, 10, 10)

	if ws := img.Warnings(); len(ws) != 1 || !strings.Contains(ws.String(), "points back at the pattern") {
		t.Fatalf("the drawing said %v, want it to say which pattern closes the circle", ws)
	}
	nothingPainted(t, cv)
}

// TestRenderPaintsNothingWithAPatternThatHoldsNothing is a pattern with no
// picture inside it, which tiles the shape with nothing the way the file asked:
// an empty tile over every part of it.
func TestRenderPaintsNothingWithAPatternThatHoldsNothing(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<defs>
			<pattern id="p" patternUnits="userSpaceOnUse" width="4" height="4"/>
		</defs>
		<rect width="10" height="10" fill="url(#p)"/>
	</svg>`, 10, 10)

	nothingPainted(t, cv)
}

// TestRenderFallsBackToTheColourAfterAPatternItCannotTile is the second chance
// a shape gets: a tile of no width tiles nothing, the pattern says so, and the
// blue written after the `url(...)` is what the shape is filled with instead.
func TestRenderFallsBackToTheColourAfterAPatternItCannotTile(t *testing.T) {
	img, cv := painted(t, `<svg viewBox="0 0 10 10">
		<defs>
			<pattern id="p" patternUnits="userSpaceOnUse" width="0" height="4">
				<rect width="2" height="2" fill="#00ff00"/>
			</pattern>
		</defs>
		<rect width="10" height="10" fill="url(#p) #0000ff"/>
	</svg>`, 10, 10)

	if ws := img.Warnings(); len(ws) != 1 || !strings.Contains(ws.String(), "tiles nothing") {
		t.Errorf("the warnings are %v, want one about a tile that tiles nothing", ws)
	}
	for _, pt := range [][2]int{{1, 1}, {5, 5}, {9, 9}} {
		blueAt(t, cv, pt[0], pt[1])
	}
}

// TestParseSaysWhatItCannotTileWithAPatternOfNoSize is the warning without the
// second chance: a tile of no size in a drawing that gave no colour after the
// `url(...)`, so the shape comes out with nothing on it and the file is told
// which of the two is the reason.
func TestParseSaysWhatItCannotTileWithAPatternOfNoSize(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 10 10">
		<defs>
			<pattern id="p" patternUnits="userSpaceOnUse" width="0" height="4">
				<rect width="2" height="2" fill="#00ff00"/>
			</pattern>
		</defs>
		<rect width="10" height="10" fill="url(#p)"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	ws := img.Warnings()
	if len(ws) != 1 || !strings.Contains(ws.String(), "tiles nothing") {
		t.Errorf("the warnings are %v, want one about a tile that tiles nothing", ws)
	}
	nothingPainted(t, img.Render(10, 10))
}

// TestParseSaysWhatItCannotReadInAPattern is everything a pattern may say that
// this package does not follow: a viewBox it does not fit, a transform it
// cannot work out, a pattern it was to take what it left out from. Each is said
// once, and what the pattern can still do — be a tile where it stands — it
// still does.
func TestParseSaysWhatItCannotReadInAPattern(t *testing.T) {
	const head = `<svg viewBox="0 0 10 10"><defs><pattern id="p" patternUnits="userSpaceOnUse"`
	for _, tc := range []struct {
		name, extra, want string
	}{
		{"a width that is not a length", ` width="nope" height="4"`, "is not a length"},
		{"a viewBox", ` width="4" height="4" viewBox="0 0 4 4"`, "viewBox"},
		{"a transform", ` width="4" height="4" patternTransform="spin(3)"`, "not one this package can read"},
		{"what it does not say", ` width="4" height="4" href="#q"`, "the pattern references"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			src := head + tc.extra + `><rect width="2" height="2" fill="#00ff00"/></pattern></defs>` +
				`<rect width="10" height="10" fill="url(#p)"/></svg>`
			img, err := Parse(src)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			ws := img.Warnings()
			if len(ws) != 1 {
				t.Fatalf("the warnings are %v, want one", ws)
			}
			if !strings.Contains(ws.String(), tc.want) {
				t.Errorf("the warnings %q do not mention %q", ws.String(), tc.want)
			}
		})
	}
}

// TestRenderPaintsAPatternInTheCurrentColour carries the colour the drawing is
// painted in into the picture inside the tile, so the same pattern paints in
// the blue of one widget and the grey of another without being written twice —
// the whole reason `currentColor` is a colour at all.
func TestRenderPaintsAPatternInTheCurrentColour(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 10 10">
		<defs>
			<pattern id="p" patternUnits="userSpaceOnUse" width="4" height="4">
				<rect width="2" height="2" fill="currentColor"/>
			</pattern>
		</defs>
		<rect width="10" height="10" fill="url(#p)"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	cv := img.RenderWith(10, 10, canvas.RGB(0, 0, 255))
	if cv == nil {
		t.Fatal("the drawing did not paint")
	}
	for _, pt := range [][2]int{{1, 1}, {5, 5}, {9, 1}} {
		blueAt(t, cv, pt[0], pt[1])
	}
	for _, pt := range [][2]int{{3, 1}, {1, 3}, {7, 7}} {
		clearAt(t, cv, pt[0], pt[1])
	}
}

// TestRenderTakesAPatternFromTheGroupItIsInside is a fill on a group reaching
// the shapes inside it, which is how a file writes one pattern for a whole
// drawing rather than once per shape: the rectangle inherits the paint the same
// way it inherits any other colour.
func TestRenderTakesAPatternFromTheGroupItIsInside(t *testing.T) {
	cv := cut(t, `<svg viewBox="0 0 10 10">
		<defs>
			<pattern id="p" patternUnits="userSpaceOnUse" width="4" height="4">
				<rect width="2" height="2" fill="#00ff00"/>
			</pattern>
		</defs>
		<g fill="url(#p)">
			<rect width="10" height="10"/>
		</g>
	</svg>`, 10, 10)

	for _, pt := range [][2]int{{1, 1}, {5, 1}, {9, 1}, {1, 5}, {9, 9}} {
		greenAt(t, cv, pt[0], pt[1])
	}
	for _, pt := range [][2]int{{3, 3}, {7, 7}, {3, 7}} {
		clearAt(t, cv, pt[0], pt[1])
	}
}

func TestRenderLeavesOutWhatAPatternTileIsWrittenOutside(t *testing.T) {
	// A pattern tile is the window its picture is seen through: what is
	// written outside it comes out nowhere, which is the `overflow` a pattern
	// stands for unless it says `visible`. The picture here sits wholly to the
	// left of the tile it was written for, so the tile keeps none of it and
	// the shape it paints is left with nothing to show at all.
	cv := cut(t, `<svg viewBox="0 0 8 8">
		<defs>
			<pattern id="p" patternUnits="userSpaceOnUse" width="4" height="4">
				<rect x="-2" y="1" width="1.5" height="2" fill="#ff0000"/>
			</pattern>
		</defs>
		<rect width="8" height="8" fill="url(#p)"/>
	</svg>`, 8, 8)

	nothingPainted(t, cv)
}

func TestRenderPaintsWhatOverflowsAPatternTileWhereItLands(t *testing.T) {
	// `overflow="visible"` turns the cut off: the picture written off to the
	// side of its tile is painted into the tile from the copy of it that
	// lands there, so what runs out of one tile shows in the one beside it
	// rather than being lost. Here the whole picture lies outside the tile to
	// its left, and comes back two units to the right of where it was written,
	// at every copy of the tile across the shape.
	cv := cut(t, `<svg viewBox="0 0 8 8">
		<defs>
			<pattern id="p" patternUnits="userSpaceOnUse" width="4" height="4" overflow="visible">
				<rect x="-2" y="1" width="1.5" height="2" fill="#ff0000"/>
			</pattern>
		</defs>
		<rect width="8" height="8" fill="url(#p)"/>
	</svg>`, 8, 8)

	kept(t, cv, 2, 1) // the copy that landed in the first tile
	kept(t, cv, 6, 5) // and the same in the tiles beside it
	gone(t, cv, 0, 1) // where the picture does not reach
	gone(t, cv, 4, 1) // nor here
}

func TestRenderPaintsAPatternWrittenInFractionsOfItsBoxPastTheTileToo(t *testing.T) {
	// The picture in fractions of the box of the shape — the copy it lands
	// in has to be worked out through that box the same way the picture is
	// drawn through it, or the copy lands a whole box out of place. Here the
	// picture is written left of the tile in fractions, comes back inside the
	// tile, and covers the first two columns of it wherever the shape asks.
	cv := cut(t, `<svg viewBox="0 0 8 8">
		<defs>
			<pattern id="p" patternUnits="userSpaceOnUse" patternContentUnits="objectBoundingBox"
				width="4" height="4" overflow="visible">
				<rect x="-0.5" y="0.25" width="0.25" height="0.5" fill="#ff0000"/>
			</pattern>
		</defs>
		<rect width="8" height="8" fill="url(#p)"/>
	</svg>`, 8, 8)

	kept(t, cv, 0, 1) // the copy that landed at the tile's own edge
	kept(t, cv, 1, 3) // covering the second column whole
	gone(t, cv, 2, 1) // and nothing beside it
}

func TestRenderDoesNotFollowAPatternOverflowingFurtherThanAHandfulOfTiles(t *testing.T) {
	// A picture reaching past its tile is drawn from the copies of it that
	// land in the tile — but only a handful of them: a picture spanning more
	// than about five tiles would be drawn dozens of times over to show a
	// corner of it, so it takes the plain cut instead and the tile keeps only
	// what was written inside it. The picture here is written as a stripe in
	// every fifth unit from far left to far right, and only the stripe written
	// over the tile survives that.
	cv := cut(t, `<svg viewBox="0 0 8 8">
		<defs>
			<pattern id="p" patternUnits="userSpaceOnUse" width="4" height="4" overflow="visible">
				<rect x="-12" y="1" width="1" height="2" fill="#ff0000"/>
				<rect x="-7" y="1" width="1" height="2" fill="#ff0000"/>
				<rect x="-2" y="1" width="1" height="2" fill="#ff0000"/>
				<rect x="3" y="1" width="1" height="2" fill="#ff0000"/>
				<rect x="8" y="1" width="1" height="2" fill="#ff0000"/>
				<rect x="13" y="1" width="1" height="2" fill="#ff0000"/>
			</pattern>
		</defs>
		<rect width="8" height="8" fill="url(#p)"/>
	</svg>`, 8, 8)

	kept(t, cv, 3, 1) // the stripe written over the tile
	gone(t, cv, 0, 1) // every other one, which is what the plain cut leaves out
	gone(t, cv, 2, 1)
}
