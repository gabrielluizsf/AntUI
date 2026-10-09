package svg

import (
	"sync"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// A whole drawing is read from the text it was written in: its size, what is
// inside it, and what could not be read. These check the reading of a drawing
// from end to end, and then that painting it says what it should on the canvas.

func TestParseReadsTheSizeOfTheDrawing(t *testing.T) {
	img, err := Parse(`<svg width="24" height="32" viewBox="0 0 12 16"><rect/></svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if img.Width != 24 || img.Height != 32 {
		t.Errorf("size = %vx%v, want 24x32", img.Width, img.Height)
	}
	if img.ViewBox != [4]float64{0, 0, 12, 16} {
		t.Errorf("viewBox = %v, want 0 0 12 16", img.ViewBox)
	}
	if img.Root == nil {
		t.Fatal("there is no root to paint")
	}
}

func TestParseReadsTheViewBoxOfADrawingWithNoSize(t *testing.T) {
	// A drawing may be given a viewBox alone, and that is what its size is then
	// taken from, since the viewBox says how big it was drawn.
	img, err := Parse(`<svg viewBox="0 0 12 16"><rect/></svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if img.Width != 12 || img.Height != 16 {
		t.Errorf("size = %vx%v, want the viewBox's 12x16", img.Width, img.Height)
	}
}

func TestParseFallsBackToASizeOfItsOwn(t *testing.T) {
	// A drawing with neither a size nor a viewBox says nothing about how big it
	// is, so it is given one to be painted at rather than being left with none.
	img, err := Parse(`<svg><rect/></svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if img.Width <= 0 || img.Height <= 0 {
		t.Errorf("size = %vx%v, want one to fall back on", img.Width, img.Height)
	}
	if img.ViewBox[2] <= 0 || img.ViewBox[3] <= 0 {
		t.Errorf("viewBox = %v, want one to fall back on", img.ViewBox)
	}
}

func TestParsePutsTheShapesWhereTheyAreInTheDrawing(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 100 100">
		<g fill="red"><rect x="10" y="10" width="20" height="20"/><circle cx="50" cy="50" r="10"/></g>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	rect := img.Root.find("rect")
	if rect == nil {
		t.Fatal("the rectangle is not in the drawing")
	}
	if rect.Path == nil {
		t.Fatal("the rectangle has no shape")
	}
	if !rect.Style.HasFill || rect.Style.Fill != canvas.RGB(255, 0, 0) {
		t.Errorf("the rectangle's fill = %v, want the group's red", rect.Style.Fill)
	}
	if b := boundsOf(rect.Path); b.MaxX > 30.001 || b.MaxY > 30.001 {
		t.Errorf("the rectangle reaches %v, want it to stop at 30,30", b)
	}
	// The circle is the second child of the same group and takes the same fill.
	circle := img.Root.find("circle")
	if circle == nil || circle.Style.Fill != canvas.RGB(255, 0, 0) {
		t.Error("the circle did not take the group's fill")
	}
}

func TestParseLeavesOutWhatItCannotDraw(t *testing.T) {
	// Something a drawing asked for that this package cannot do — here a
	// picture behind an address that answers nothing — is a warning rather
	// than an error, and the rest of the drawing still paints.
	img, err := Parse(`<svg viewBox="0 0 10 10"><image href="a.png" width="5" height="5"/><rect x="0" y="0" width="5" height="5"/></svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if img.Root.find("rect") == nil {
		t.Error("the rectangle went out with the image")
	}
	if len(img.Warnings()) == 0 {
		t.Error("nothing said that the picture in the image could not be read")
	}
}

func TestParseWarnsAboutWhatItCannotRead(t *testing.T) {
	// Damage that is not structural is said and then got past, so one bad
	// attribute does not cost the whole drawing.
	for _, src := range []string{
		`<svg viewBox="0 0 10 10"><rect width="abc" height="5"/></svg>`,
		`<svg viewBox="0 0 10 10"><path d="M0 0 X10 10"/></svg>`,
		`<svg viewBox="0 0 10 10"><polygon points="1,2 3"/></svg>`,
		`<svg viewBox="0 0 10 10"><nonsense/></svg>`,
		`<svg viewBox="0 0 10 10"><rect transform="nonsense(3)"/></svg>`,
	} {
		img, err := Parse(src)
		if err != nil {
			t.Errorf("parse %s: %v", src, err)
			continue
		}
		if len(img.Warnings()) == 0 {
			t.Errorf("parse %s: it said nothing about what it could not read", src)
		}
	}
}

func TestParseRefusesWhatItCannotReadAtAll(t *testing.T) {
	// A tag that is never closed is not a drawing, and there is nothing to paint
	// out of it.
	for _, src := range []string{
		``,
		`   `,
		`not a drawing`,
		`<rect/>`,
		`<svg><g></svg>`,
		`<svg><g><rect/></svg>`,
	} {
		if img, err := Parse(src); err == nil {
			t.Errorf("parse %q: it read a drawing out of %+v", src, img)
		}
	}
}

func TestParseAnImageIsSafeToUseFromSeveralGoroutines(t *testing.T) {
	// An icon is drawn once and painted from every frame, which may be from
	// several goroutines at once.
	// A drawing with something in it that could not be drawn, so that painting it
	// over and over also asks for the warnings again and again: everything said
	// about a drawing is said while it is read, which is over before it is
	// shared, and nothing is written to it afterwards.
	img, err := Parse(`<svg viewBox="0 0 24 24">
		<circle cx="12" cy="12" r="10" fill="red"/>
		<image href="a.png"/>
		<path d="M0 0 X1 1"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(img.Warnings()) == 0 {
		t.Fatal("the drawing said nothing about the two things it could not draw")
	}
	var wg sync.WaitGroup
	for i := range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 5 {
				if cv := img.Render(24, 24); cv == nil {
					t.Error("nothing came back")
					return
				}
				_ = i
			}
		}()
	}
	wg.Wait()
}

func TestRenderPaintsWhatTheDrawingSays(t *testing.T) {
	// A red circle in the middle of a drawing of a hundred across, painted at a
	// hundred pixels across, is red in the middle and clear at the edge.
	img, err := Parse(`<svg viewBox="0 0 100 100"><circle cx="50" cy="50" r="50" fill="red"/></svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	cv := img.Render(100, 100)
	if cv == nil {
		t.Fatal("nothing came back")
	}
	if got := pixelAt(cv, 50, 50); got != canvas.RGB(255, 0, 0) {
		t.Errorf("the middle is %v, want it red", got)
	}
	// A corner is outside the circle, and a canvas is transparent where nothing
	// has been painted.
	if got := pixelAt(cv, 1, 1); got.A() != 0 {
		t.Errorf("the corner is %v, want it clear", got)
	}
}

func TestRenderIsAskedForTheSizeItIsGiven(t *testing.T) {
	// A drawing is asked for at whatever size it is wanted, and it is painted
	// filling that size whichever one it is.
	img, err := Parse(`<svg viewBox="0 0 10 10"><rect x="0" y="0" width="10" height="10" fill="blue"/></svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	for _, size := range []int{8, 16, 32, 100} {
		cv := img.Render(size, size)
		if cv == nil {
			t.Fatalf("%d: nothing came back", size)
		}
		if cv.Width != size || cv.Height != size {
			t.Errorf("asked for %d, got a canvas of %dx%d", size, cv.Width, cv.Height)
		}
		if got := pixelAt(cv, size/2, size/2); got != canvas.RGB(0, 0, 255) {
			t.Errorf("at %d the middle is %v, want it blue", size, got)
		}
		// A drawing that fills its own box fills the canvas it is painted on.
		if got := pixelAt(cv, 1, 1); got.A() == 0 {
			t.Errorf("at %d the corner is clear, want the rectangle to reach it", size)
		}
	}
}

func TestRenderStrokesOverFills(t *testing.T) {
	// A shape with both a fill and a stroke has the stroke over the fill, which
	// is the order they are drawn in and the order a browser draws them in.
	img, err := Parse(`<svg viewBox="0 0 20 20">
		<rect x="5" y="5" width="10" height="10" fill="red" stroke="blue" stroke-width="4"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	cv := img.Render(20, 20)
	if got := pixelAt(cv, 5, 5); got != canvas.RGB(0, 0, 255) {
		t.Errorf("the corner where the stroke is = %v, want it blue", got)
	}
	if got := pixelAt(cv, 10, 10); got != canvas.RGB(255, 0, 0) {
		t.Errorf("the middle under the stroke = %v, want it red", got)
	}
}

func TestRenderKeepsTheStrokeWidthAsItIsWritten(t *testing.T) {
	// A drawing scaled up is its shapes that grow, not the lines drawn on them:
	// a stroke of one on a drawing of twenty painted at two hundred is ten
	// pixels wide, and not two hundred, which is what a stroke on the pixels
	// themselves would be.
	img, err := Parse(`<svg viewBox="0 0 20 20">
		<line x1="0" y1="10" x2="20" y2="10" stroke="red" stroke-width="1"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	cv := img.Render(200, 200)
	// The line is at y=10 of 20, which is y=100 of 200, and it is ten pixels
	// thick, so it covers from 95 to 105 and not the whole of the canvas.
	if got := pixelAt(cv, 100, 100); got != canvas.RGB(255, 0, 0) {
		t.Errorf("the middle of the line is %v, want it red", got)
	}
	if got := pixelAt(cv, 100, 130); got.A() != 0 {
		t.Errorf("far below the line is %v, want it clear: the stroke is ten wide, not two hundred", got)
	}
}

func TestRenderStrokesFollowANonUniformScale(t *testing.T) {
	// A transform that stretches one axis stretches the stroke on it too: a
	// vertical line scaled twice across is twice as thick, while the same line
	// left alone is one thick. The width follows the shape rather than an
	// average of the two axes, which is what `scale(2,1)` would otherwise give.
	img, err := Parse(`<svg viewBox="0 0 20 20">
		<line x1="4" y1="0" x2="4" y2="20" stroke="red" stroke-width="2" transform="scale(2,1)"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	cv := img.Render(20, 20)
	// The line is at x=4, taken to x=8, and its width of two is stretched with
	// it, so across it covers x=6 to x=10 — where reading the width off the
	// average of the two axes would reach only from x≈6.6 and leave x=6 part
	// covered instead of solid.
	for _, x := range []int{6, 7, 8, 9} {
		if got := pixelAt(cv, x, 10); got != canvas.RGB(255, 0, 0) {
			t.Errorf("across the stretched stroke at x=%d is %v, want it red", x, got)
		}
	}
	if got := pixelAt(cv, 11, 10); got.A() != 0 {
		t.Errorf("past the stroke at x=11 is %v, want it clear", got)
	}
}

func TestRenderTransformsMoveTheShape(t *testing.T) {
	// A shape moved to one side is painted on that side, and the stroke on it
	// keeps the width it was written with.
	img, err := Parse(`<svg viewBox="0 0 20 20">
		<rect x="2" y="8" width="4" height="4" fill="red" transform="translate(10, 0)"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	cv := img.Render(20, 20)
	if got := pixelAt(cv, 14, 10); got != canvas.RGB(255, 0, 0) {
		t.Errorf("where it was moved to is %v, want it red", got)
	}
	if got := pixelAt(cv, 3, 10); got.A() != 0 {
		t.Errorf("where it was before is %v, want it clear", got)
	}
}

func TestRenderSkipsWhatIsHidden(t *testing.T) {
	// A shape that says it is not to be drawn is not drawn, and the rest of the
	// drawing is.
	img, err := Parse(`<svg viewBox="0 0 20 20">
		<rect x="0" y="0" width="10" height="10" fill="red" display="none"/>
		<rect x="10" y="10" width="10" height="10" fill="blue"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	cv := img.Render(20, 20)
	if got := pixelAt(cv, 5, 5); got.A() != 0 {
		t.Errorf("the hidden shape is %v, want it clear", got)
	}
	if got := pixelAt(cv, 15, 15); got != canvas.RGB(0, 0, 255) {
		t.Errorf("the shape that was not hidden is %v, want it blue", got)
	}
}

func TestRenderOpacities(t *testing.T) {
	// An opacity says how much of what is painted shows through, and a canvas
	// has nothing behind a pixel to show through to — what comes out is opaque
	// either way — so an opacity is seen as a colour mixed with what was painted
	// underneath it. A half-opaque red over a white one is halfway between them.
	for _, tc := range []struct {
		name  string
		attr  string
		red   uint8
		green uint8
	}{
		{"half of the fill", `fill-opacity="0.5"`, 255, 127},
		{"a quarter of everything", `opacity="0.25"`, 255, 191},
		{"all of it", `opacity="1"`, 255, 0},
		{"none of it", `opacity="0"`, 255, 255},
	} {
		img, err := Parse(`<svg viewBox="0 0 20 20">
			<rect x="0" y="0" width="20" height="20" fill="white"/>
			<rect x="0" y="0" width="20" height="20" fill="red" ` + tc.attr + `/>
		</svg>`)
		if err != nil {
			t.Fatalf("%s: parse: %v", tc.name, err)
		}
		got := pixelAt(img.Render(20, 20), 10, 10)
		if !nearByte(got.R(), tc.red) || !nearByte(got.G(), tc.green) {
			t.Errorf("%s: got %v, want the red over the white at #%02X%02X%02X",
				tc.name, got, tc.red, tc.green, tc.green)
		}
	}
}

func TestRenderDashes(t *testing.T) {
	// A dashed stroke has gaps in it, and a line drawn whole does not. A line
	// across the middle of the drawing with a dash of four and a gap of four has
	// a gap somewhere along it, and a line with no pattern does not.
	img, err := Parse(`<svg viewBox="0 0 20 20">
		<line x1="0" y1="10" x2="20" y2="10" stroke="red" stroke-width="2" stroke-dasharray="2 2"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	cv := img.Render(20, 20)
	painted, clear := 0, 0
	for x := range 20 {
		if pixelAt(cv, x, 10).A() != 0 {
			painted++
		} else {
			clear++
		}
	}
	if painted == 0 || clear == 0 {
		t.Errorf("along the line %d are painted and %d are clear, want some of each", painted, clear)
	}
}

func TestRenderRemembersWhatItHasDrawn(t *testing.T) {
	// The same drawing at the same size is drawn once, and asking again answers
	// the canvas that was made, which is what a button that redraws itself every
	// frame relies on.
	img, err := Parse(`<svg viewBox="0 0 10 10"><rect width="10" height="10" fill="red"/></svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	first := img.Render(16, 16)
	second := img.Render(16, 16)
	if first != second {
		t.Error("asking twice for the same size drew it twice")
	}
	// A different size is drawn again, since the one remembered is the wrong shape
	// for it.
	other := img.Render(32, 32)
	if other == first {
		t.Error("a different size answered the canvas that was remembered")
	}
	if other.Width != 32 {
		t.Errorf("asked for 32, got a canvas of %d", other.Width)
	}
	// Going back to the first size draws it once more rather than answering with
	// a canvas of the wrong size.
	back := img.Render(16, 16)
	if back.Width != 16 {
		t.Errorf("asked for 16 again, got a canvas of %d", back.Width)
	}
	if back == other {
		t.Error("it answered with the canvas of the size before")
	}
}

func TestRenderNothingToPaint(t *testing.T) {
	// A size of nothing is not a canvas to draw on, and a drawing that was never
	// read has nothing in it to draw.
	var none *Image
	if cv := none.Render(10, 10); cv != nil {
		t.Error("a drawing that is not there painted something")
	}
	img, err := Parse(`<svg viewBox="0 0 10 10"><rect width="10" height="10"/></svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	for _, size := range [2]int{0, -5} {
		if cv := img.Render(size, size); cv != nil {
			t.Errorf("a size of %d painted something", size)
		}
	}
}

func TestRenderAGroupTransformsWhatIsInsideIt(t *testing.T) {
	// A transform on a group moves everything inside it, and the shapes inside
	// are painted where the group put them.
	img, err := Parse(`<svg viewBox="0 0 20 20">
		<g transform="translate(5, 5)">
			<rect x="0" y="0" width="5" height="5" fill="red"/>
		</g>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	cv := img.Render(20, 20)
	if got := pixelAt(cv, 7, 7); got != canvas.RGB(255, 0, 0) {
		t.Errorf("inside the moved group is %v, want it red", got)
	}
	if got := pixelAt(cv, 2, 2); got.A() != 0 {
		t.Errorf("where it would be without the move is %v, want it clear", got)
	}
}

func TestRenderANestedGroupMultipliesItsTransforms(t *testing.T) {
	// Two moves are two moves, and the shape inside lands where both of them
	// together put it.
	img, err := Parse(`<svg viewBox="0 0 20 20">
		<g transform="translate(5, 0)">
			<g transform="translate(0, 5)">
				<rect x="0" y="0" width="4" height="4" fill="red"/>
			</g>
		</g>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	cv := img.Render(20, 20)
	if got := pixelAt(cv, 6, 6); got != canvas.RGB(255, 0, 0) {
		t.Errorf("inside both moves is %v, want it red", got)
	}
}

// find is the first node of that name in the drawing, and nil where there is
// none. It is the image's own, and the tests use it to be sure of what was read
// rather than to guess from the bounds alone.

// nearByte says whether two of the eight bits of a colour are the same to within
// the rounding a blend of them does.
func nearByte(got, want uint8) bool {
	d := int(got) - int(want)
	return d <= 1 && d >= -1
}

// pixelAt is the colour of one pixel of a canvas, which is what a test says
// about where a shape was painted. A pixel outside the canvas has no colour, and
// so reads as clear.
func pixelAt(cv *canvas.Canvas, x, y int) canvas.Color {
	if x < 0 || y < 0 || x >= cv.Width || y >= cv.Height {
		return canvas.Transparent
	}
	return cv.At(x, y)
}

// A drawing is painted into a picture that is then laid over a widget, so the
// picture has to keep the transparency of the shapes in it. An edge that falls
// between two pixels is covered in part, and rendering it onto a surface rather
// than onto a layer answers that pixel opaque — the colour comes out a quarter
// as bright and fully solid, which is a dark rim round every soft edge in the
// drawing. Rendering onto a layer is what keeps the edge soft all the way to
// the widget behind it.
func TestRenderKeepsTheTransparencyOfASoftEdge(t *testing.T) {
	// A disc of an odd radius in an even box puts its boundary through a column
	// of pixel centres, so the pixels along it are covered in part rather than
	// all or none.
	img, err := Parse(`<svg width="20" height="20"><circle cx="10" cy="10" r="7" fill="#ffffff"/></svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	cv := img.Render(20, 20)
	if cv == nil {
		t.Fatal("the drawing did not paint")
	}

	var covered, partial, blank int
	for y := range cv.Height {
		for x := range cv.Width {
			switch a := cv.At(x, y).A(); {
			case a == 0:
				blank++
			case a == 255:
				covered++
			default:
				partial++
			}
		}
	}
	if covered == 0 {
		t.Error("nothing in the picture is solid, want the middle of the disc to be")
	}
	if blank == 0 {
		t.Error("nothing in the picture is see-through, want the corners to be")
	}
	if partial == 0 {
		t.Error("the disc has no part-covered pixels, want a soft edge rather than a hard one")
	}
	// The corner of the picture is outside the disc in every way.
	if got := cv.At(0, 0).A(); got != 0 {
		t.Errorf("the corner is at alpha %d, want it clear of the disc", got)
	}
	// The middle of the disc is inside it, and solid white.
	if got := cv.At(10, 10); got.A() != 255 || got.R() != 255 {
		t.Errorf("the middle of the disc is %#08x, want solid white", uint32(got))
	}
	// A part-covered pixel keeps the colour of the shape rather than fading
	// towards nothing, so laying it over a widget tints it instead of darkening
	// it. This is the difference a layer makes and the reason for it.
	for y := range cv.Height {
		for x := range cv.Width {
			if c := cv.At(x, y); c.A() > 0 && c.A() < 255 && c.R() != 255 {
				t.Fatalf("the edge pixel at %d,%d is red %d, want the white of the disc kept", x, y, c.R())
			}
		}
	}
}

// Two shapes that overlap on a layer add up where they meet rather than the
// second hiding the first, which is what an opaque surface would do.
func TestRenderAddsUpWhereTwoShapesOverlap(t *testing.T) {
	img, err := Parse(`<svg width="20" height="20">` +
		`<circle cx="8" cy="10" r="6" fill="#ffffff"/>` +
		`<circle cx="12" cy="10" r="6" fill="#ffffff"/>` +
		`</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	cv := img.Render(20, 20)
	if cv == nil {
		t.Fatal("the drawing did not paint")
	}
	// Where the two discs cross, the picture is as solid as either of them; the
	// edge past it is not, and the difference between the two is the soft edge.
	var solid, soft int
	for y := range cv.Height {
		for x := range cv.Width {
			switch a := cv.At(x, y).A(); {
			case a == 255:
				solid++
			case a > 0:
				soft++
			}
		}
	}
	if solid == 0 {
		t.Error("nothing in the picture is solid")
	}
	if soft == 0 {
		t.Error("the overlapping discs have no soft edge")
	}
	// The overlap itself is solid: the picture is a union of the two shapes.
	if got := cv.At(10, 10).A(); got != 255 {
		t.Errorf("where the discs overlap is at alpha %d, want it solid", got)
	}
}
