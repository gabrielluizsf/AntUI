package svg

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// currentColorDrawing is the shape of an icon that follows the colour it is
// painted in: one drawing, any colour, which is what lets a single file serve
// every state of a control.
const currentColorDrawing = `<svg viewBox="0 0 20 20">
	<rect x="0" y="0" width="20" height="20" fill="currentColor"/>
</svg>`

func TestRenderPaintsCurrentColorInTheColourItIsGiven(t *testing.T) {
	// The keyword does not name a colour, it stands for the one the drawing is
	// being painted in, so the colour handed to the painting is the one that
	// comes out. The same drawing, read once, in two colours.
	img, err := Parse(currentColorDrawing)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	for _, want := range []canvas.Color{
		canvas.RGB(255, 0, 0),
		canvas.RGB(0, 128, 255),
		canvas.RGB(0, 0, 0),
	} {
		got := pixelAt(img.RenderWith(20, 20, want), 10, 10)
		if got.R() != want.R() || got.G() != want.G() || got.B() != want.B() {
			t.Errorf("painted in %v, got %v", want, got)
		}
	}
}

func TestRenderPaintsCurrentColorBlackByDefault(t *testing.T) {
	// Render says nothing about a colour, and a drawing written in currentColor
	// still has to paint something: black, which is what the keyword means
	// where nothing else has said what the colour is.
	img, err := Parse(currentColorDrawing)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	got := pixelAt(img.Render(20, 20), 10, 10)
	if got.R() != 0 || got.G() != 0 || got.B() != 0 || got.A() != 255 {
		t.Errorf("Render painted %v, want opaque black", got)
	}
}

func TestRenderRemembersOneColourAtATime(t *testing.T) {
	// The colour is part of what is remembered: a canvas painted red is not the
	// answer to a question about blue. Asking again for either is free, and
	// asking in a new one draws the drawing again.
	img, err := Parse(currentColorDrawing)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	red := img.RenderWith(20, 20, canvas.RGB(255, 0, 0))
	again := img.RenderWith(20, 20, canvas.RGB(255, 0, 0))
	if again != red {
		t.Error("asking twice in the same colour drew it twice")
	}
	blue := img.RenderWith(20, 20, canvas.RGB(0, 0, 255))
	if blue == red {
		t.Error("a second colour answered the canvas of the first")
	}
	if got := pixelAt(blue, 10, 10); got.B() != 255 || got.R() != 0 {
		t.Errorf("the second colour painted %v, want blue", got)
	}
	// Back to the first colour: drawn again rather than answered with the wrong
	// one, which is the same rule the size follows.
	back := img.RenderWith(20, 20, canvas.RGB(255, 0, 0))
	if back == blue {
		t.Error("it answered with the canvas of the colour before")
	}
	if got := pixelAt(back, 10, 10); got.R() != 255 {
		t.Errorf("going back to red painted %v, want red", got)
	}
}

func TestRenderCurrentColorOnAStroke(t *testing.T) {
	// A stroke may be painted in the colour too, and a line across the middle of
	// the drawing puts it where one pixel of the answer is that line.
	img, err := Parse(`<svg viewBox="0 0 20 20">
		<line x1="0" y1="10" x2="20" y2="10" stroke="currentColor" stroke-width="4"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	got := pixelAt(img.RenderWith(20, 20, canvas.RGB(255, 255, 0)), 10, 10)
	if got.R() != 255 || got.G() != 255 || got.B() != 0 {
		t.Errorf("the stroke painted %v, want yellow", got)
	}
}

func TestRenderCurrentColorIsInheritedByWhatIsInsideAGroup(t *testing.T) {
	// The keyword is inherited like any other paint, so a group that says
	// currentColor paints everything under it in the colour, and a child that
	// names a colour of its own keeps it.
	img, err := Parse(`<svg viewBox="0 0 20 20">
		<g fill="currentColor">
			<rect x="0" y="0" width="10" height="10"/>
			<rect x="10" y="10" width="10" height="10" fill="green"/>
		</g>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	cv := img.RenderWith(20, 20, canvas.RGB(0, 0, 255))
	if got := pixelAt(cv, 5, 5); got.B() != 255 || got.R() != 0 {
		t.Errorf("the inherited shape painted %v, want the colour it was given", got)
	}
	// `green` is the CSS green, which is a dark one — #008000 — and not the
	// green a screen is usually thought of as.
	if got := pixelAt(cv, 15, 15); got != canvas.RGB(0, 128, 0) {
		t.Errorf("the shape with a colour of its own painted %v, want the green it was given", got)
	}
}

func TestRenderReadsCurrentColorInAnyCase(t *testing.T) {
	// CSS keywords are written in any case, and a file that spells it
	// `currentcolor` is not asking for a colour of that name.
	img, err := Parse(`<svg viewBox="0 0 20 20">
		<rect width="20" height="20" fill="currentcolor"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	got := pixelAt(img.RenderWith(20, 20, canvas.RGB(255, 0, 0)), 10, 10)
	if got.R() != 255 {
		t.Errorf("painted %v, want the colour it was given", got)
	}
}

func TestRenderCurrentColorIsNoWarning(t *testing.T) {
	// A drawing written entirely in currentColor is a drawing this package can
	// paint, so it has nothing to say against it.
	img, err := Parse(currentColorDrawing)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if w := img.Warnings(); len(w) != 0 {
		t.Errorf("a drawing in currentColor warned: %v", w)
	}
}

func TestParseReadsCurrentColorAsAPaint(t *testing.T) {
	// The shape is painted — a keyword that resolved to nothing would leave the
	// node with no fill at all — and it is marked as the colour it follows, so
	// the painter knows not to use the colour that was never read.
	img, err := Parse(currentColorDrawing)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	node := img.Root.find("rect")
	if node == nil {
		t.Fatal("the shape is not in the drawing")
	}
	if !node.Style.HasFill {
		t.Error("a shape painted in currentColor has no fill at all")
	}
	if !node.Style.FillCurrent {
		t.Error("a shape painted in currentColor is not marked as following the colour")
	}
}
