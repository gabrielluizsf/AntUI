package template

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// TestClipPathChamfersTheShape checks the site's signature button shape: a
// clip-path cones the top-right and bottom-left corners, so those pixels are
// left blank while the middle is filled and the two uncut corners are square.
func TestClipPathChamfersTheShape(t *testing.T) {
	win := blankWin(t, 40, 40)
	// The site's cut, on a 20x20 box at 10,10: 6px off the top-right and
	// bottom-left corners.
	st := bgRule(t, `background-color: red;
		clip-path: polygon(0 0, calc(100% - 6px) 0, 100% 6px, 100% 100%, 6px 100%, 0 calc(100% - 6px));`)
	cs := newCSSStyle(win)
	cs.paintBackground(win, st, 10, 10, 20, 20, 0, 0)

	red := canvas.RGBA(canvas.RGBAParams{R: 255, A: 255})
	blank := canvas.RGBA(canvas.RGBAParams{R: 30, G: 40, B: 50, A: 255})

	if got := win.Canvas().At(20, 20); got != red {
		t.Errorf("centre = %v, want the fill", got)
	}
	if got := win.Canvas().At(11, 11); got != red {
		t.Errorf("top-left corner must stay square, got %v", got)
	}
	if got := win.Canvas().At(28, 11); got != blank {
		t.Errorf("top-right chamfer must be cut, got %v", got)
	}
	if got := win.Canvas().At(11, 28); got != blank {
		t.Errorf("bottom-left chamfer must be cut, got %v", got)
	}
}

// TestClipPathFallsBackWhenUnparsable checks a value the engine cannot draw
// leaves the box alone rather than dropping the widget.
func TestClipPathFallsBackWhenUnparsable(t *testing.T) {
	sh, err := css.Parse("label { clip-path: circle(50%); }")
	if err != nil {
		t.Fatal(err)
	}
	if p := sh.Style(css.StyleQuery{Tag: "label"}).ClipPath; len(p) != 0 {
		t.Errorf("a circle must not become a polygon, got %v", p)
	}
}
