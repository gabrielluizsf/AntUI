package svg

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// writtenBox is the box a drawing's writing came out in: the pixels anything
// was painted on, as the same four numbers boundsOf answers for a path, so a
// test can say where a piece of writing landed the way it says where a shape
// did. Nothing painted is an empty box rather than a failure, which is what a
// test about what is left out asks for.
func writtenBox(cv *canvas.Canvas) box {
	return coloredBox(cv, func(canvas.Color) bool { return true })
}

// coloredBox is the box of the pixels that keep says yes to, which is how a
// test tells one colour of writing from another in the same picture.
func coloredBox(cv *canvas.Canvas, keep func(canvas.Color) bool) box {
	b := box{Empty: true}
	for y := range cv.Height {
		for x := range cv.Width {
			c := cv.At(x, y)
			if c.A() == 0 || !keep(c) {
				continue
			}
			if b.Empty {
				b = box{MinX: float64(x), MinY: float64(y), MaxX: float64(x), MaxY: float64(y)}
				continue
			}
			b.MinX = min(b.MinX, float64(x))
			b.MinY = min(b.MinY, float64(y))
			b.MaxX = max(b.MaxX, float64(x))
			b.MaxY = max(b.MaxY, float64(y))
		}
	}
	return b
}

// firstPainted is the colour of the first pixel anything was painted on, which
// is a simpler thing to ask of a picture than where its writing is: the corner
// of the box of a piece of writing is the outside of it rather than one of its
// own pixels.
func firstPainted(cv *canvas.Canvas) canvas.Color {
	for y := range cv.Height {
		for x := range cv.Width {
			if c := cv.At(x, y); c.A() > 0 {
				return c
			}
		}
	}
	return canvas.Transparent
}

// boxOf is where one piece of writing landed in a forty by twenty drawing, as
// the box of the pixels it painted on.
func boxOf(t *testing.T, body string) box {
	t.Helper()
	_, cv := painted(t, `<svg viewBox="0 0 40 20">`+body+`</svg>`, 40, 20)
	return writtenBox(cv)
}

// TestRenderWritesAtThePenTheDrawingNames puts the same writing at two places
// and expects the second to be exactly as far along as the drawing asked,
// because where the pen goes is the one thing a file writing its own text has
// to be able to say.
func TestRenderWritesAtThePenTheDrawingNames(t *testing.T) {
	base := boxOf(t, `<text x="4" y="16" font-size="8">Hi</text>`)
	if base.Empty {
		t.Fatal("nothing was written at all")
	}
	here := boxOf(t, `<text x="24" y="16" font-size="8">Hi</text>`)
	if got, want := here.MinX-base.MinX, 20.0; got != want {
		t.Errorf("moving the writing along 20 moved it %g, want %g (boxes %s and %s)", got, want, base, here)
	}
	up := boxOf(t, `<text x="4" y="10" font-size="8">Hi</text>`)
	if got, want := up.MinY-base.MinY, -6.0; got != want {
		t.Errorf("lifting the writing 6 lifted it %g, want %g (boxes %s and %s)", got, want, base, up)
	}
	// A drawing that says nothing about where the writing starts starts at the
	// origin of its own coordinates, which is the left edge of the picture.
	unnamed := boxOf(t, `<text y="16" font-size="8">Hi</text>`)
	if unnamed.MinX != 0 {
		t.Errorf("writing with no x of its own starts at %g, want the origin", unnamed.MinX)
	}
	// The writing sits above its baseline: a baseline at sixteen leaves the
	// letters in the rows before it, none of them hanging below it.
	if base.MaxY > 16 {
		t.Errorf("the writing reaches down to %g, past the baseline of 16", base.MaxY)
	}
}

// TestRenderWritesAtTheSizeItSays takes a bigger font-size to a bigger box, and
// a size written on the group to the same picture as the one written on the
// writing itself, since a size is the sort of thing a file says once for
// everything under it.
func TestRenderWritesAtTheSizeItSays(t *testing.T) {
	small := boxOf(t, `<text x="4" y="16" font-size="8">Hi</text>`)
	big := boxOf(t, `<text x="4" y="16" font-size="16">Hi</text>`)
	if big.MaxX-big.MinX <= small.MaxX-small.MinX {
		t.Errorf("font-size 16 wrote %g across, no wider than font-size 8 at %g", big.MaxX-big.MinX, small.MaxX-small.MinX)
	}
	if big.MaxY-big.MinY <= small.MaxY-small.MinY {
		t.Errorf("font-size 16 wrote %g tall, no taller than font-size 8 at %g", big.MaxY-big.MinY, small.MaxY-small.MinY)
	}
	// Sixteen is what a size of nothing says: a drawing that does not name one
	// is written at the same size as one that names the default.
	if got := boxOf(t, `<text x="4" y="16">Hi</text>`); got != big {
		t.Errorf("writing with no font-size came out %s, want the default size's %s", got, big)
	}
	// A size on the group is the size of what is inside it.
	if got := boxOf(t, `<g font-size="8"><text x="4" y="16">Hi</text></g>`); got != small {
		t.Errorf("size inherited from the group came out %s, want %s", got, small)
	}
	// And a size on the writing itself beats the one it inherited.
	if got := boxOf(t, `<g font-size="8"><text x="4" y="16" font-size="16">Hi</text></g>`); got != big {
		t.Errorf("size written on the writing came out %s, want %s", got, big)
	}
}

// TestRenderHangsTheWritingOnThePen moves one piece of writing through the
// three places it may hang: at the pen, centred on it, and ending at it. The
// last is what puts a label to the left of the point it names.
func TestRenderHangsTheWritingOnThePen(t *testing.T) {
	start := boxOf(t, `<text x="20" y="16" font-size="8">Hi</text>`)
	middle := boxOf(t, `<text x="20" y="16" font-size="8" text-anchor="middle">Hi</text>`)
	end := boxOf(t, `<text x="20" y="16" font-size="8" text-anchor="end">Hi</text>`)
	if start.Empty || middle.Empty || end.Empty {
		t.Fatalf("one of the three came out empty: %s, %s, %s", start, middle, end)
	}
	if !(end.MinX < middle.MinX && middle.MinX < start.MinX) {
		t.Errorf("the three hang as %s, %s, %s; want the end furthest left and the start furthest right", end, middle, start)
	}
	// The three are the same writing, so all three are as wide as each other.
	width := start.MaxX - start.MinX
	for _, b := range []box{middle, end} {
		if b.MaxX-b.MinX != width {
			t.Errorf("%s is %g across, want the %g of the writing at the start", b, b.MaxX-b.MinX, width)
		}
	}
	// Centred on the pen of twenty, the writing spans ten either side of it,
	// which is the middle of the box at twenty.
	if got, lo, hi := (middle.MinX+middle.MaxX)/2, 19.0, 21.0; got < lo || got > hi {
		t.Errorf("the writing centred on 20 came out centred at %g, want it between %g and %g", got, lo, hi)
	}
	// Ending at the pen, the last pixel of the writing is at the pen.
	if end.MaxX < 19 || end.MaxX > 20 {
		t.Errorf("the writing ending at 20 ends at %g, want it at the pen", end.MaxX)
	}
}

// TestRenderCarriesThePenAcrossThePieces writes in two pieces and expects the
// second to start where the first left off, which is the whole of what a
// `<tspan>` is for: it is not a new piece of writing at the origin, it is the
// writing carried on.
func TestRenderCarriesThePenAcrossThePieces(t *testing.T) {
	_, cv := painted(t, `<svg viewBox="0 0 40 20">`+
		`<text x="4" y="16" font-size="8">A<tspan fill="#ff0000">B</tspan></text>`+
		`</svg>`, 40, 20)
	isRed := func(c canvas.Color) bool { return c.R() > 128 }
	black := coloredBox(cv, func(c canvas.Color) bool { return !isRed(c) })
	red := coloredBox(cv, isRed)
	if black.Empty || red.Empty {
		t.Fatalf("the two pieces did not both come out: black %s, red %s", black, red)
	}
	if red.MinX <= black.MinX {
		t.Errorf("the second piece starts at %g, not after the first at %g", red.MinX, black.MinX)
	}
	if red.MinX <= black.MaxX {
		t.Errorf("the second piece starts at %g, over the first which reaches %g", red.MinX, black.MaxX)
	}
	// A move along from where the writing already is shifts it by exactly that
	// much, because it is added to the place the pen had reached.
	plain := boxOf(t, `<text x="4" y="16" font-size="8">A</text>`)
	moved := boxOf(t, `<text x="4" y="16" font-size="8" dx="6">A</text>`)
	if got, want := moved.MinX-plain.MinX, 6.0; got != want {
		t.Errorf("a dx of 6 moved the writing %g, want %g", got, want)
	}
	lifted := boxOf(t, `<text x="4" y="16" font-size="8" dy="3">A</text>`)
	if got, want := lifted.MinY-plain.MinY, 3.0; got != want {
		t.Errorf("a dy of 3 lifted the writing %g, want %g", got, want)
	}
}

// TestRenderCarriesThePositionOverWhitespaceWritesTheSpaceInBetween takes the
// whitespace at the two ends of the writing for what it is — the file's own
// line breaks, which say nothing about the drawing — and keeps the whitespace
// between two pieces of it, which is a space and is part of the writing.
func TestRenderCarriesThePositionOverWhitespaceWritesTheSpaceInBetween(t *testing.T) {
	// The whitespace in front of the writing is dropped with the position still
	// on it: the place the drawing named is where the writing begins, whether or
	// not a line break sits between the two in the file.
	leading := boxOf(t, `<text x="12" y="16" font-size="8"> <tspan>A</tspan></text>`)
	if leading.Empty {
		t.Fatal("nothing was written")
	}
	if leading.MinX != 12 {
		t.Errorf("the writing after the whitespace starts at %g, want the x of 12 it was given", leading.MinX)
	}
	// The same writing with a piece before it keeps the space between them: the
	// second piece is a space further on than it would be without it.
	apart := boxOf(t, `<text x="4" y="16" font-size="8">A<tspan> B</tspan></text>`)
	together := boxOf(t, `<text x="4" y="16" font-size="8">AB</text>`)
	if apart.MaxX <= together.MaxX {
		t.Errorf("the space between the two pieces was dropped: %s reaches no further than %s", apart, together)
	}
	// And the whitespace behind the writing goes too, taking no space with it.
	trailing := boxOf(t, `<text x="4" y="16" font-size="8">A </text>`)
	if got, want := trailing, boxOf(t, `<text x="4" y="16" font-size="8">A</text>`); got != want {
		t.Errorf("writing with a space after it came out %s, want %s", got, want)
	}
}

// TestRenderPaintsTheFillItIsGiven checks the three things a fill may be: a
// colour of its own, the colour the picture was painted in, and no colour at
// all.
func TestRenderPaintsTheFillItIsGiven(t *testing.T) {
	const body = `<text x="4" y="16" font-size="8"`
	for _, tc := range []struct {
		name, fill string
		want       canvas.Color
	}{
		{"a colour", `fill="#ff0000"`, canvas.RGB(255, 0, 0)},
		{"nothing", `fill="none"`, canvas.Transparent},
		{"currentColor", `fill="currentColor"`, canvas.RGB(0, 0, 255)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			img, err := Parse(`<svg viewBox="0 0 40 20">` + body + ` ` + tc.fill + `>Hi</text></svg>`)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			cv := img.RenderWith(40, 20, canvas.RGB(0, 0, 255))
			if cv == nil {
				t.Fatal("the drawing did not paint")
			}
			b := writtenBox(cv)
			if tc.want.A() == 0 {
				if !b.Empty {
					t.Errorf("writing filled with nothing painted %s, want no picture at all", b)
				}
				return
			}
			if b.Empty {
				t.Fatal("nothing was written")
			}
			// A letter is painted a pixel at a time and a letter only half
			// covered keeps its colour rather than fading towards the picture
			// behind it, which is what makes it a letter rather than a smudge.
			if got := firstPainted(cv); got.R() != tc.want.R() || got.G() != tc.want.G() || got.B() != tc.want.B() {
				t.Errorf("the writing came out %#08x, want the colour %#08x", uint32(got), uint32(tc.want))
			}
		})
	}
}

// TestRenderKeepsHiddenWritingOffThePicture takes `display="none"` for what it
// says: the writing is in the drawing and not on the picture, and that is
// ordinary enough to say nothing about.
func TestRenderKeepsHiddenWritingOffThePicture(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 40 20"><text x="4" y="16" font-size="8" display="none">Hi</text></svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if b := writtenBox(img.Render(40, 20)); !b.Empty {
		t.Errorf("hidden writing painted %s, want no picture at all", b)
	}
	if n := img.Warnings(); len(n) != 0 {
		t.Errorf("hidden writing was said out loud: %v", n)
	}
}

// TestParseKeepsTheWritingInTheOrderItWasWritten reads a `<text>` back out of
// the drawing: the pieces in the order they were written, each with the place
// and the style it was given, and the style of the writing around it where it
// was given none.
func TestParseKeepsTheWritingInTheOrderItWasWritten(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 40 20"><text x="4" y="16" font-size="8" text-anchor="end">Hi ` +
		`<tspan fill="#ff0000" x="20">there</tspan> !</text></svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	text := img.Root.find("text")
	if text == nil {
		t.Fatal("the drawing has no text in it")
	}
	if len(text.Kids) != 0 {
		t.Errorf("the writing built %d nodes inside the text, want it kept as pieces instead", len(text.Kids))
	}
	runs := text.Runs
	if len(runs) != 3 {
		t.Fatalf("got %d pieces of writing, want three: %+v", len(runs), runs)
	}
	// The first piece is the writing up to the `<tspan>`, with the style the
	// `<text>` gave it: the size, where it hangs, and the place it starts at.
	if got := runs[0].Text; got != "Hi " {
		t.Errorf("the first piece is %q, want \"Hi \" — the space before the tspan is part of the writing", got)
	}
	if !runs[0].HasX || runs[0].X != 4 {
		t.Errorf("the first piece starts at %v (said %v), want the x of 4", runs[0].X, runs[0].HasX)
	}
	if runs[0].Style.FontSize != 8 || runs[0].Style.Anchor != AnchorEnd {
		t.Errorf("the first piece was read with size %g hanging at %d, want size 8 ending at the pen", runs[0].Style.FontSize, runs[0].Style.Anchor)
	}
	// The second is the `<tspan>`'s own writing, at the place the `<tspan>`
	// named rather than the one on the `<text>`, in the colour it asked for.
	if got := runs[1].Text; got != "there" {
		t.Errorf("the second piece is %q, want \"there\"", got)
	}
	if !runs[1].HasX || runs[1].X != 20 {
		t.Errorf("the second piece starts at %v (said %v), want the x of 20 on the tspan", runs[1].X, runs[1].HasX)
	}
	if runs[1].Style.Fill != canvas.RGB(255, 0, 0) {
		t.Errorf("the second piece is filled %#08x, want the red it asked for", uint32(runs[1].Style.Fill))
	}
	// The third carries on where the second left off, with no place of its own.
	if got := runs[2].Text; got != " !" {
		t.Errorf("the third piece is %q, want \" !\"", got)
	}
	if runs[2].HasX {
		t.Errorf("the third piece was given an x of %g, want none: it starts where the one before it ended", runs[2].X)
	}
	// A `<tspan>` with nothing written in it is not a piece of the writing at
	// all rather than an empty piece taking up a place in it.
	img, err = Parse(`<svg viewBox="0 0 40 20"><text>Hi<tspan font-size="8"/></text></svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	empty := img.Root.find("text")
	if empty == nil {
		t.Fatal("the drawing has no text in it")
	}
	if got := len(empty.Runs); got != 1 {
		t.Errorf("a tspan with nothing in it gave %d pieces, want the one piece of writing", got)
	}
}

// TestParseReadsTheSizeOfTheWriting reads a font-size the way CSS writes it: a
// length, a share of the size already held, a name, and a step up or down from
// what it has.
func TestParseReadsTheSizeOfTheWriting(t *testing.T) {
	for _, tc := range []struct {
		raw  string
		want float64
	}{
		{"12", 12},
		{"12px", 12},
		{"0.5in", 48},
		{"medium", 16},
		{"large", 18},
		{"50%", 8}, // half of the size an element with nothing says
		{"larger", 19.2},
		{"smaller", 16 / 1.2},
	} {
		e := mustElement(t, `<text font-size="`+tc.raw+`"/>`)
		got := Style{}.with(e, quiet)
		if got.FontSize != tc.want {
			t.Errorf("font-size %q came out %g, want %g", tc.raw, got.FontSize, tc.want)
		}
	}
	// A step up or down is from the size the writing already has, not from the
	// default: a group that sized its contents sizes the step too.
	e := mustElement(t, `<text font-size="larger"/>`)
	if got := (Style{FontSize: 20}).with(e, quiet); got.FontSize != 24 {
		t.Errorf("larger than 20 came out %g, want 24", got.FontSize)
	}
	// A size that is not one says so and keeps the size it had rather than
	// writing at nothing.
	var said string
	e = mustElement(t, `<text font-size="nope"/>`)
	if got := (Style{FontSize: 12}).with(e, func(f string, a ...any) { said = fmt.Sprintf(f, a...) }); got.FontSize != 12 {
		t.Errorf("an unreadable size came out %g, want the 12 it had", got.FontSize)
	}
	if !strings.Contains(said, "nope") {
		t.Errorf("the warning %q does not say what could not be read", said)
	}
}

// TestParseSaysWhatItCannotDoWithWriting is every sort of writing this package
// cannot put down: a stroke it does not draw, a gradient it cannot fill with, a
// turn it does not follow, and anything inside the writing that is not writing.
func TestParseSaysWhatItCannotDoWithWriting(t *testing.T) {
	const head = `<svg viewBox="0 0 40 20"><defs><linearGradient id="g">` +
		`<stop offset="0" stop-color="#ff0000"/><stop offset="1" stop-color="#0000ff"/>` +
		`</linearGradient></defs>`
	for _, tc := range []struct {
		name, src, want string
	}{
		{"a stroke", head + `<text x="4" y="16" stroke="red">Hi</text></svg>`, "stroke"},
		{"a gradient", head + `<text x="4" y="16" fill="url(#g)">Hi</text></svg>`, "gradient"},
		{"a turn", head + `<g transform="rotate(30)"><text x="4" y="16">Hi</text></g></svg>`, "turned"},
		{"a place for each letter", head + `<text x="4 9" y="16">Hi</text></svg>`, "each letter"},
		{"a position that is not a number", head + `<text x="nope" y="16">Hi</text></svg>`, "not a number"},
		{"an anchor that is not one", head + `<text text-anchor="beside">Hi</text></svg>`, "text-anchor"},
		{"a child that is not writing", head + `<text x="4" y="16">Hi<rect width="1" height="1"/></text></svg>`, "<rect>"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			img, err := Parse(tc.src)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			ws := img.Warnings()
			if len(ws) == 0 {
				t.Fatal("nothing was said about the writing that could not be put down")
			}
			if !strings.Contains(ws.String(), tc.want) {
				t.Errorf("the warnings %q do not mention %q", ws.String(), tc.want)
			}
		})
	}
}

// TestRenderWritesAnOrphanTspan paints a `<tspan>` that is not inside a
// `<text>` at all, which no file should have and which is still writing with a
// place and a size rather than a shape with no outline.
func TestRenderWritesAnOrphanTspan(t *testing.T) {
	if b := boxOf(t, `<g><tspan x="4" y="16" font-size="8">A</tspan></g>`); b.Empty {
		t.Error("a tspan outside any text painted nothing")
	}
}

// TestTrimEdgesTakesTheWhitespaceOffTheEndsAndKeepsThePlace is the piece of
// reading the writing above rests on: what is dropped with the whitespace at
// the two ends of it, and what of it goes onto the writing that is left.
func TestTrimEdgesTakesTheWhitespaceOffTheEndsAndKeepsThePlace(t *testing.T) {
	runs := trimEdges([]TextRun{
		{Text: " ", HasX: true, X: 12, DX: 3},
		{Text: "hello "},
		{Text: " ", HasY: true, Y: 9},
		{Text: "there"},
		{Text: "   "},
	})
	if len(runs) != 3 {
		t.Fatalf("got %d pieces, want three: %+v", len(runs), runs)
	}
	if runs[0].Text != "hello " {
		t.Errorf("the first piece is %q, want \"hello \" with its space kept", runs[0].Text)
	}
	// The place the dropped whitespace held is where the writing after it goes.
	if !runs[0].HasX || runs[0].X != 12 || runs[0].DX != 3 {
		t.Errorf("the first piece starts at %g with dx %g (said %v), want the x of 12 and dx of 3 the dropped piece held", runs[0].X, runs[0].DX, runs[0].HasX)
	}
	if runs[1].Text != " " {
		t.Errorf("the piece in the middle is %q, want the single space it was written as", runs[1].Text)
	}
	if runs[2].Text != "there" {
		t.Errorf("the last piece is %q, want \"there\"", runs[2].Text)
	}
	// A position written on the piece that is already there is not moved by
	// one that was dropped in front of it: the nearest of them counts.
	runs = trimEdges([]TextRun{{Text: " "}, {Text: "hello", HasY: true, Y: 4}})
	if len(runs) != 1 || runs[0].Y != 4 {
		t.Errorf("got %+v, want one piece at y 4", runs)
	}
	// Nor is it moved by one dropped in front of it that named a place the
	// writing after it had named for itself.
	runs = trimEdges([]TextRun{{Text: " ", HasX: true, X: 2}, {Text: "hello", HasX: true, X: 9}})
	if len(runs) != 1 || runs[0].X != 9 {
		t.Errorf("got %+v, want one piece still at x 9", runs)
	}
}
