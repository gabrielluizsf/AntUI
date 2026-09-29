package template

// The two ways a border can paint itself: filled over a surface the sheet
// declared, and stroked over one it did not.

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// TestBorderOverItsOwnBackground: a ring drawn the way it always has been, the
// whole box in the border colour and the background over the middle of it.
func TestBorderOverItsOwnBackground(t *testing.T) {
	win := blankWin(t, 40, 40)
	st := bgRule(t, "background-color: #0000FF; border: 2px solid #FF0000; border-radius: 6px;")
	cs := newCSSStyle(win)
	cs.paintBox(win, st, 10, 10, 20, 20)

	if got := win.Canvas().At(20, 20); got != canvas.RGB(0, 0, 0xFF) {
		t.Errorf("the middle of the box is %v, want the background it declared", got)
	}
	if got := win.Canvas().At(10, 20); got != canvas.RGB(0xFF, 0, 0) {
		t.Errorf("the left edge of the ring is %v, want red", got)
	}
}

// TestBorderWithoutABackgroundLeavesThePageBehindIt: a sheet that gives a box a
// border and no surface has nothing to paint over the inside of the ring, so the
// page shows through the middle of it. The theme's surface is not the box's to
// paint, and a light theme's is white — a border with no background is not a
// white box.
func TestBorderWithoutABackgroundLeavesThePageBehindIt(t *testing.T) {
	win := blankWin(t, 40, 40)
	st := bgRule(t, "border: 2px solid #FF0000; border-radius: 6px;")
	cs := newCSSStyle(win)
	cs.paintBox(win, st, 10, 10, 20, 20)

	if got := win.Canvas().At(20, 20); got != canvas.RGBA(30, 40, 50, 255) {
		t.Errorf("the middle of the ring is %v, want the page behind it, untouched", got)
	}
	if got := win.Canvas().At(10, 20); got != canvas.RGB(0xFF, 0, 0) {
		t.Errorf("the left edge of the ring is %v, want red", got)
	}
}

// TestNoPaintWithoutADeclaredFill: a fill the sheet never gave stays
// transparent all the way to the canvas. Fading it would replace its missing
// alpha with the style's opacity and hand the canvas the black its colour
// carries underneath, which paints a box the sheet said nothing about.
func TestNoPaintWithoutADeclaredFill(t *testing.T) {
	win := blankWin(t, 40, 40)
	cs := newCSSStyle(win)

	off := false
	cs.Switch(win, State{}, 0, 0, "pinned", off)
	cs.SelectOption(win, State{}, 0, 30, 40, 10, "stable", false)

	for y := 0; y < 40; y++ {
		for x := 0; x < 40; x++ {
			switch got := win.Canvas().At(x, y); got {
			case canvas.RGB(0, 0, 0):
				t.Fatalf("an opaque black pixel at %d,%d: a transparent fill faded into a black box", x, y)
			case canvas.RGB(255, 255, 255):
				t.Fatalf("a white pixel at %d,%d: a surface painted where the sheet declared none", x, y)
			}
		}
	}
}

// TestToggleSurfaceStaysOnTheControl: a checkbox is a box of eighteen pixels
// with its label beside it, and the background its style declares is the box's,
// not a surface running under the label. The cell it was given is the label's,
// and a cell with no background of its own must not be painted.
func TestToggleSurfaceStaysOnTheControl(t *testing.T) {
	win := blankWin(t, 200, 40)
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, "checkbox { background-color: rgb(0,0,255); }")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	win.Begin()
	tpl.Checkbox("Enable", new(bool))
	win.End()

	// The box itself: an eighteen pixel control at the cell's left, so the
	// surface is the control's and the page is still behind it.
	cv := win.Canvas()
	if got := cv.At(72, 9); got != canvas.RGB(0, 0, 255) {
		t.Errorf("inside the box at 72,9 is %v, want the declared background", got)
	}
	// Under the label, beside the box: the cell carries no background of its
	// own, so nothing of the surface runs out to here.
	if got := cv.At(110, 9); got != canvas.RGBA(30, 40, 50, 255) {
		t.Errorf("under the label at 110,9 is %v, want the page behind the cell", got)
	}
}
